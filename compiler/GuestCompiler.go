package main

import (
	"fmt"
	"github.com/MetaFFI/plugin-sdk/compiler/go/IDL"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"text/template"
)

//--------------------------------------------------------------------
func getDynamicLibSuffix() string{
	switch runtime.GOOS{
	case "windows": return ".dll"
	case "darwin": return ".dylib"
	default: // We might need to make this more specific in the future
		return ".so"
	}
}
//--------------------------------------------------------------------
type GuestCompiler struct{
	def *IDL.IDLDefinition
	outputDir string
	outputFilename string
	blockName string
	blockCode string
}
//--------------------------------------------------------------------
func NewGuestCompiler(definition *IDL.IDLDefinition, outputDir string, outputFilename string, blockName string, blockCode string) *GuestCompiler{

	if strings.Contains(outputFilename, "#"){
		toRemove := outputFilename[strings.LastIndex(outputFilename, string(os.PathSeparator))+1:strings.Index(outputFilename, "#")+1]
		outputFilename = strings.ReplaceAll(outputFilename, toRemove, "")
	}

	outputFilename = strings.ReplaceAll(outputFilename, filepath.Ext(outputFilename), "")

	return &GuestCompiler{def: definition, outputDir: outputDir, outputFilename: outputFilename, blockName: blockName, blockCode: blockCode}
}
//--------------------------------------------------------------------
func (this *GuestCompiler) Compile() (outputFileName string, err error){

	// generate code
	code, err := this.generateCode()
	if err != nil{
		return "", fmt.Errorf("Failed to generate guest code: %v", err)
	}

	file, err := this.buildDynamicLibrary(code)
	if err != nil{
		return "", fmt.Errorf("Failed to generate guest code: %v", err)
	}

	// write to output
	outputFullFileName := fmt.Sprintf("%v%v%v_MetaFFIGuest%v", this.outputDir, string(os.PathSeparator), this.outputFilename, getDynamicLibSuffix())
	err = ioutil.WriteFile(outputFullFileName, file, 0700)
	if err != nil{
		return "", fmt.Errorf("Failed to write dynamic library to %v. Error: %v", this.outputDir+this.outputFilename, err)
	}

	return outputFullFileName, nil

}
//--------------------------------------------------------------------
func (this *GuestCompiler) parseHeader() (string, error){
	tmp, err := template.New("headers").Parse(GuestHeaderTemplate)
	if err != nil{
		return "", fmt.Errorf("Failed to parse GuestHeaderTemplate: %v", err)
	}

	buf := strings.Builder{}
	err = tmp.Execute(&buf, this.def)

	return buf.String(), err
}
//--------------------------------------------------------------------
func (this *GuestCompiler) parseImports() (string, error){

	// get all imports from the def file
	imports := struct {
		Imports []string
		Modules []*IDL.ModuleDefinition
	}{
		Imports: make([]string, 0),
		Modules: this.def.Modules,
	}

	set := make(map[string]bool)

	for _, m := range this.def.Modules{

		for _, f := range m.Functions{
			if pack, found := f.FunctionPath["package"]; found{

				if pack != `main`{
					set[os.ExpandEnv(pack)] = true
				}
			}

			if mod, found := f.FunctionPath["module"]; found{

				mod = os.ExpandEnv(mod)
				if fi, _ := os.Stat(mod); fi == nil { // ignore if module is local item
					set[os.ExpandEnv(mod)] = true
				}
			}
		}
	}

	for k, _ := range set{
		imports.Imports = append(imports.Imports, k)
	}

	tmp, err := template.New("imports").Funcs(templatesFuncMap).Parse(GuestImportsTemplate)
	if err != nil{
		return "", fmt.Errorf("Failed to parse GuestImportsTemplate: %v", err)
	}

	buf := strings.Builder{}
	err = tmp.Execute(&buf, imports)
	importsCode := buf.String()

	return importsCode, err
}
//--------------------------------------------------------------------
func (this *GuestCompiler) parseGuestHelperFunctions() (string, error){

	tmpEntryPoint, err := template.New("helper").Funcs(templatesFuncMap).Parse(GuestHelperFunctionsTemplate)
	if err != nil{
		return "", fmt.Errorf("Failed to parse GuestFunctionXLLRTemplate: %v", err)
	}

	bufEntryPoint := strings.Builder{}
	err = tmpEntryPoint.Execute(&bufEntryPoint, this.def)

	return bufEntryPoint.String(), err
}
//--------------------------------------------------------------------
func (this *GuestCompiler) parseForeignFunctions() (string, error){

	tmpEntryPoint, err := template.New("foreignfuncs").Funcs(templatesFuncMap).Parse(GuestFunctionXLLRTemplate)
	if err != nil{
		return "", fmt.Errorf("Failed to parse GuestFunctionXLLRTemplate: %v", err)
	}

	bufEntryPoint := strings.Builder{}
	err = tmpEntryPoint.Execute(&bufEntryPoint, this.def)

	return bufEntryPoint.String(), err
}
//--------------------------------------------------------------------
func (this *GuestCompiler) parseCImportsCGoFile() (string, error){

	tmp, err := template.New("guest").Funcs(templatesFuncMap).Parse(GuestCImportCGoFileTemplate)
	if err != nil{
		return "", fmt.Errorf("Failed to parse GuestCImportCGoFileTemplate: %v", err)
	}

	buf := strings.Builder{}
	err = tmp.Execute(&buf, nil)

	return buf.String(), err
}
//--------------------------------------------------------------------
func (this *GuestCompiler) parseCImports() (string, error){

	tmp, err := template.New("guest").Funcs(templatesFuncMap).Parse(GuestCImportTemplate)
	if err != nil{
		return "", fmt.Errorf("Failed to parse GuestCImportTemplate: %v", err)
	}

	buf := strings.Builder{}
	err = tmp.Execute(&buf, nil)

	return buf.String(), err
}
//--------------------------------------------------------------------
func (this *GuestCompiler) generateCode() (string, error){


	header, err := this.parseHeader()
	if err != nil{ return "", err }

	imports, err := this.parseImports()
	if err != nil{ return "", err }

	cimports, err := this.parseCImports()
	if err != nil{ return "", err }

	guestHelpers, err := this.parseGuestHelperFunctions()
	if err != nil{ return "", err }

	functionStubs, err := this.parseForeignFunctions()
	if err != nil{ return "", err }

	res := header + imports + cimports + functionStubs + guestHelpers + GuestMainFunction

	return res, nil
}
//--------------------------------------------------------------------
func (this *GuestCompiler) buildDynamicLibrary(code string)([]byte, error){

	dir, err := os.MkdirTemp("", "metaffi_go_compiler*")
	if err != nil{
		return nil, fmt.Errorf("Failed to create temp dir to build code: %v", err)
	}
	defer func(){ if err == nil{ _ = os.RemoveAll(dir) } }()

	dir = dir+string(os.PathSeparator)

	err = ioutil.WriteFile(dir+"metaffi_guest.go", []byte(code), 0700)
	if err != nil{
		return nil, fmt.Errorf("Failed to write guest go code: %v", err)
	}

	// TODO: This should move to "generate code" that need to return a map of files
	cgoCode, err := this.parseCImportsCGoFile()
	if err != nil{
		return nil, fmt.Errorf("Failed to generate CGo guest go code: %v", err)
	}

	err = ioutil.WriteFile(dir+"metaffi_guest_cgo.go", []byte(cgoCode), 0700)
	if err != nil{
		return nil, fmt.Errorf("Failed to write guest go code: %v", err)
	}


	fmt.Println("Building Go foreign functions")

	// add go.mod
	_, err = this.goModInit(dir, "main")
	if err != nil{
		return nil, err
	}

	// add "replace"s if there are local imports
	addedLocalModules := make(map[string]bool)
	for _, m := range this.def.Modules{
		for _, f := range m.Functions {
			for k, v := range f.FunctionPath {
				if k == "module" {
					v = os.ExpandEnv(v)
					if fi, _ := os.Stat(v); fi != nil && fi.IsDir() { // if module is local dir
						if _, alreadyAdded := addedLocalModules[v]; !alreadyAdded {

							// if embedded code, write the source code into a Package folder and skip "-replace"
							if this.blockCode != ""{


								packageDir := dir+os.ExpandEnv(f.FunctionPath["package"])+string(os.PathSeparator)
								err = os.Mkdir(packageDir, 0777)
								if err != nil {
									return nil, fmt.Errorf("Failed creating directory for embedded code: %v.\nError:\n%v", packageDir, err)
								}

								err = ioutil.WriteFile(packageDir+f.FunctionPath["package"]+".go", []byte(this.blockCode), 0700)
								if err != nil{
									return nil, fmt.Errorf("Failed to embedded block go code: %v", err)
								}

								_, err := this.goModInit(packageDir, f.FunctionPath["package"])
								if err != nil{ return nil, err }

								_, err = this.goGet(packageDir)
								if err != nil{ return nil, err }

								_, err = this.goReplace(dir, os.ExpandEnv(f.FunctionPath["package"]), "./"+f.FunctionPath["package"])
								if err != nil{ return nil, err }

							} else {
								// point module to
								_, err = this.goReplace(dir, os.ExpandEnv(f.FunctionPath["package"]), v)
								if err != nil{ return nil, err }
							}

							addedLocalModules[v] = true
						}
					}
				}
			}
		}
	}

	_, err = ioutil.ReadFile(dir+"go.mod")
	if err != nil{
		println("Failed to find go.mod in "+dir+"go.mod")
	}

	// build dynamic library
	_, err = this.goGet(dir)
	if err != nil{ return nil, err	}

	_, err = this.goClean(dir)
	if err != nil{ return nil, err	}

	_, err = this.goBuild(dir)
	if err != nil{ return nil, err	}

	// copy to output dir
	fileData, err := ioutil.ReadFile(dir+this.outputFilename+getDynamicLibSuffix())
	if err != nil{
		return nil, fmt.Errorf("Failed to read created dynamic library. Error: %v", err)
	}

	return fileData, nil
}
//--------------------------------------------------------------------
func (this *GuestCompiler) goModInit(dir string, packageName string) (string, error) {
	modInitCmd := exec.Command("go", "mod", "init", packageName)
	modInitCmd.Dir = dir
	fmt.Printf("%v\n", strings.Join(modInitCmd.Args, " "))
	output, err := modInitCmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("Failed building Go foreign function with error: %v.\nOutput:\n%v", err, string(output))
	}

	return string(output), err
}
//--------------------------------------------------------------------
func (this *GuestCompiler) goGet(dir string) (string, error) {
	getCmd := exec.Command("go", "get", "-v")
	getCmd.Dir = dir
	fmt.Printf("%v\n", strings.Join(getCmd.Args, " "))
	output, err := getCmd.CombinedOutput()
	if err != nil{
		println(string(output))
		return "", fmt.Errorf("Failed building Go foreign function in \"%v\" with error: %v.\nOutput:\n%v", dir, err, string(output))
	}

	return string(output), err
}
//--------------------------------------------------------------------
func (this *GuestCompiler) goClean(dir string) (string, error) {
	cleanCmd := exec.Command("go", "clean", "-cache")
	cleanCmd.Dir = dir
	fmt.Printf("%v\n", strings.Join(cleanCmd.Args, " "))
	output, err := cleanCmd.CombinedOutput()
	if err != nil{
		return "", fmt.Errorf("Failed building Go foreign function in \"%v\" with error: %v.\nOutput:\n%v", dir, err, string(output))
	}

	return string(output), err
}
//--------------------------------------------------------------------
func (this *GuestCompiler) goBuild(dir string) (string, error) {
	buildCmd := exec.Command("go", "build", "-v", "-tags=guest" , "-buildmode=c-shared", "-gcflags=-shared", "-o", dir+this.outputFilename+getDynamicLibSuffix())
	buildCmd.Dir = dir
	fmt.Printf("%v\n", strings.Join(buildCmd.Args, " "))
	output, err := buildCmd.CombinedOutput()
	if err != nil{
		return "", fmt.Errorf("Failed building Go foreign function in \"%v\" with error: %v.\nOutput:\n%v", dir, err, string(output))
	}

	return string(output), err
}
//--------------------------------------------------------------------
func (this *GuestCompiler) goReplace(dir string, packageName string, packagePath string) (string, error) {
	getCmd := exec.Command("go", "mod", "edit", "-replace", fmt.Sprintf("%v=%v", packageName, packagePath))
	getCmd.Dir = dir
	fmt.Printf("%v\n", strings.Join(getCmd.Args, " "))
	output, err := getCmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("Failed building Go foreign function with error: %v.\nOutput:\n%v", err, string(output))
	}

	return string(output), err
}
//--------------------------------------------------------------------

