package main

import (
	"fmt"
	compiler "github.com/OpenFFI/plugin-sdk/compiler/go"
	"go/parser"
	"go/token"
	"html/template"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
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
	serializationCode map[string]string
	outputFilename string
}
//--------------------------------------------------------------------
func NewGuestCompiler(definition *compiler.IDLDefinition, outputDir string, outputFilename string, serializationCode map[string]string) *GuestCompiler{

	serializationCodeCopy := make(map[string]string)
	for k, v := range serializationCode{
		serializationCodeCopy[k] = v
	}

	return &GuestCompiler{def: definition, outputDir: outputDir, serializationCode: serializationCodeCopy, outputFilename: outputFilename}
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
	}{
		Imports: make([]string, 0),
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

	tmp, err := template.New("guest").Parse(GuestImportsTemplate)
	if err != nil{
		return "", fmt.Errorf("Failed to parse GuestImportsTemplate: %v", err)
	}

	buf := strings.Builder{}
	err = tmp.Execute(&buf, imports)
	importsCode := buf.String()

	// get all imports from the serialization code

	for filename, code := range this.serializationCode{

		if strings.ToLower(filepath.Ext(filename)) != ".go"{
			continue
		}

		fset := token.NewFileSet()
		ast, err := parser.ParseFile(fset, "", code, parser.ImportsOnly)
		if err != nil{
			return "", fmt.Errorf("Failed to parse serialization code of file %v. Error: %v", filename, err)
		}

		// add to imports code + remove import from serialization code (so it can be appended to the rest of the file)
		for _, i := range ast.Imports{

			// if import equals "main" - skip (as it is the package name of the generated code)
			if i.Path.Value != `"main"`{
				importsCode += "import "+i.Path.Value+"\n"
			}
		}

		// remove imports from serializationCode
		removeImportRegex, err := regexp.Compile("\\n[ ]*import[^(\\\"|\\()]+\\\"[^\\\"]+\\\"|\\n[ ]*import[^(\\(|\")]+\\([^\\)]+\\)")
		if err != nil{
			return "", fmt.Errorf("Failed to create regex to remove imports from serialization code. Error: %v", err)
		}
		this.serializationCode[filename] = removeImportRegex.ReplaceAllString(this.serializationCode[filename], "")
	}

	return importsCode, err
}
//--------------------------------------------------------------------
func (this *GuestCompiler) parseForeignFunctions() (string, error){

	asPublic := func(elem string) string {
		if len(elem) == 0 {
			return ""
		} else if len(elem) == 1 {
			return strings.ToUpper(elem)
		} else {
			return strings.ToUpper(elem[0:1]) + elem[1:]
		}
	}

	toGoNameConv := func(elem string) string{
		elem = strings.Replace(elem, "_", " ", -1)
		elem = strings.Title(elem)
		return strings.Replace(elem, " ", "", -1)
	}

	passParameter := func(p interface{}) string{
		param := p.(*compiler.FieldDefinition)
		res := "req."+asPublic(param.Name)

		if param.PassMethod == "by_pointer"{
			res = "&"+res
		}

		if strings.Contains(param.Type, "int"){
			res = "int("+res+")"
		}

		return res
	}

	funcMap := map[string]interface{}{
		"AsPublic": asPublic,
		"ToGoNameConv": toGoNameConv,
		"PassParameter": passParameter,
	}

	tmpEntryPoint, err := template.New("guest").Funcs(funcMap).Parse(GuestFunctionXLLRTemplate)
	if err != nil{
		return "", fmt.Errorf("Failed to parse GuestFunctionXLLRTemplate: %v", err)
	}

	bufEntryPoint := strings.Builder{}
	err = tmpEntryPoint.Execute(&bufEntryPoint, this.def)

	return bufEntryPoint.String(), err
}
//--------------------------------------------------------------------
func (this *GuestCompiler) generateCode() (string, error){

	header, err := this.parseHeader()
	if err != nil{ return "", err }

	imports, err := this.parseImports()
	if err != nil{ return "", err }

	functionStubs, err := this.parseForeignFunctions()
	if err != nil{ return "", err }

	res := header + imports + GuestCImport + functionStubs + GuestHelperFunctions + GuestMainFunction

	// append serialization code in the same file
	for filename, serializationCode := range this.serializationCode{

		if strings.ToLower(filepath.Ext(filename)) != ".go"{
			continue
		}

		// remove "package" lines
		serializationCode = regexp.MustCompile("\npackage [^\n]+").ReplaceAllString(serializationCode, "")

		res += serializationCode
	}

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

