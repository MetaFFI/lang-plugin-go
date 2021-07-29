package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"text/template"

	compiler "github.com/OpenFFI/plugin-sdk/compiler/go"
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
	def *compiler.IDLDefinition
	outputDir string
	outputFilename string
}
//--------------------------------------------------------------------
func NewGuestCompiler(definition *compiler.IDLDefinition, outputDir string, outputFilename string) *GuestCompiler{

	return &GuestCompiler{def: definition, outputDir: outputDir, outputFilename: outputFilename}
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
	outputFullFileName := fmt.Sprintf("%v%v%v_OpenFFIGuest%v", this.outputDir, string(os.PathSeparator), this.outputFilename, getDynamicLibSuffix())
	err = ioutil.WriteFile(outputFullFileName, file, 0700)
	if err != nil{
		return "", fmt.Errorf("Failed to write dynamic library to %v. Error: %v", this.outputDir+this.outputFilename, err)
	}

	return outputFullFileName, nil

}
//--------------------------------------------------------------------
func (this *GuestCompiler) parseHeader() (string, error){
	tmp, err := template.New("guest").Parse(GuestHeaderTemplate)
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
		Modules []*compiler.ModuleDefinition
	}{
		Imports: make([]string, 0),
		Modules: this.def.Modules,
	}

	set := make(map[string]bool)

	for _, m := range this.def.Modules{
		for _, f := range m.Functions{
			if pack, found := f.PathToForeignFunction["package"]; found{

				if pack != `main`{
					set[os.ExpandEnv(pack)] = true
				}
			}

			if mod, found := f.PathToForeignFunction["module"]; found{

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

	tmp, err := template.New("guest").Funcs(templatesFuncMap).Parse(GuestImportsTemplate)
	if err != nil{
		return "", fmt.Errorf("Failed to parse GuestImportsTemplate: %v", err)
	}

	buf := strings.Builder{}
	err = tmp.Execute(&buf, imports)
	importsCode := buf.String()

	return importsCode, err
}
//--------------------------------------------------------------------
func (this *GuestCompiler) parseForeignFunctions() (string, error){

	tmpEntryPoint, err := template.New("guest").Funcs(templatesFuncMap).Parse(GuestFunctionXLLRTemplate)
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

	functionStubs, err := this.parseForeignFunctions()
	if err != nil{ return "", err }

	res := header + imports + cimports + functionStubs + GuestHelperFunctions + GuestMainFunction

	return res, nil
}
//--------------------------------------------------------------------
func (this *GuestCompiler) buildDynamicLibrary(code string)([]byte, error){

	dir, err := os.MkdirTemp("", "openffi_go_compiler*")
	if err != nil{
		return nil, fmt.Errorf("Failed to create temp dir to build code: %v", err)
	}
	defer func(){ if err == nil{ _ = os.RemoveAll(dir) } }()

	dir = dir+string(os.PathSeparator)

	err = ioutil.WriteFile(dir+"openffi_guest.go", []byte(code), 0700)
	if err != nil{
		return nil, fmt.Errorf("Failed to write guest go code: %v", err)
	}

	// TODO: This should move to "generate code" that need to return a map of files
	cgoCode, err := this.parseCImportsCGoFile()
	if err != nil{
		return nil, fmt.Errorf("Failed to generate CGo guest go code: %v", err)
	}

	err = ioutil.WriteFile(dir+"openffi_guest_cgo.go", []byte(cgoCode), 0700)
	if err != nil{
		return nil, fmt.Errorf("Failed to write guest go code: %v", err)
	}


	fmt.Println("Building Go foreign functions")

	// add go.mod
	modInitCmd := exec.Command("go", "mod", "init", "main")
	modInitCmd.Dir = dir
	fmt.Printf("%v\n", strings.Join(modInitCmd.Args, " "))
	output, err := modInitCmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("Failed building Go foreign function with error: %v.\nOutput:\n%v", err, string(output))
	}

	// add "replace"s if there are local imports
	addedLocalModules := make(map[string]bool)
	for _, m := range this.def.Modules{
		for _, f := range m.Functions {
			for k, v := range f.PathToForeignFunction {
				if k == "module" {
					v = os.ExpandEnv(v)
					if fi, _ := os.Stat(v); fi != nil && fi.IsDir() { // if module is local dir
						if _, alreadyAdded := addedLocalModules[v]; !alreadyAdded {
							// point module to
							getCmd := exec.Command("go", "mod", "edit", "-replace", fmt.Sprintf("%v=%v", os.ExpandEnv(f.PathToForeignFunction["package"]), v))
							getCmd.Dir = dir
							fmt.Printf("%v\n", strings.Join(getCmd.Args, " "))
							output, err = getCmd.CombinedOutput()
							if err != nil {
								println(string(output))
								return nil, fmt.Errorf("Failed building Go foreign function with error: %v.\nOutput:\n%v", err, string(output))
							}

							addedLocalModules[v] = true
						}
					}
				}
			}
		}
	}

	gomod, err := ioutil.ReadFile(dir+"go.mod")
	if err != nil{
		println("Failed to find go.mod in "+dir+"go.mod")
	}
	println(string(gomod))

	// build dynamic library
	getCmd := exec.Command("go", "get", "-v")
	getCmd.Dir = dir
	fmt.Printf("%v\n", strings.Join(getCmd.Args, " "))
	output, err = getCmd.CombinedOutput()
	if err != nil{
		println(string(output))
		return nil, fmt.Errorf("Failed building Go foreign function in \"%v\" with error: %v.\nOutput:\n%v", dir, err, string(output))
	}

	cleanCmd := exec.Command("go", "clean", "-cache")
	cleanCmd.Dir = dir
	fmt.Printf("%v\n", strings.Join(cleanCmd.Args, " "))
	output, err = cleanCmd.CombinedOutput()
	if err != nil{
		return nil, fmt.Errorf("Failed building Go foreign function in \"%v\" with error: %v.\nOutput:\n%v", dir, err, string(output))
	}

	buildCmd := exec.Command("go", "build", "-v", "-tags=guest" , "-buildmode=c-shared", "-gcflags=-shared", "-o", dir+this.outputFilename+getDynamicLibSuffix())
	buildCmd.Dir = dir
	fmt.Printf("%v\n", strings.Join(buildCmd.Args, " "))
	output, err = buildCmd.CombinedOutput()
	if err != nil{
		return nil, fmt.Errorf("Failed building Go foreign function in \"%v\" with error: %v.\nOutput:\n%v", dir, err, string(output))
	}

	// copy to output dir
	fileData, err := ioutil.ReadFile(dir+this.outputFilename+getDynamicLibSuffix())
	if err != nil{
		return nil, fmt.Errorf("Failed to read created dynamic library. Error: %v", err)
	}

	return fileData, nil
}
//--------------------------------------------------------------------

