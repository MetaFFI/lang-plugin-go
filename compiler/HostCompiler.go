package main

import (
	"fmt"
	compiler "github.com/OpenFFI/plugin-sdk/compiler/go"
	"go/parser"
	"go/token"
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

//--------------------------------------------------------------------
type HostCompiler struct{
	def *compiler.IDLDefinition
	outputDir string
	serializationCode map[string]string
	hostOptions map[string]string
	outputFilename string
}
//--------------------------------------------------------------------
func NewHostCompiler(definition *compiler.IDLDefinition, outputDir string, outputFilename string, serializationCode map[string]string, hostOptions map[string]string) *HostCompiler{

	serializationCodeCopy := make(map[string]string)
	for k, v := range serializationCode{
		serializationCodeCopy[k] = v
	}

	return &HostCompiler{def: definition,
		outputDir: outputDir,
		serializationCode: serializationCodeCopy,
		outputFilename: outputFilename,
		hostOptions: hostOptions}
}
//--------------------------------------------------------------------
func (this *HostCompiler) Compile() (outputFileName string, err error){

	// generate code
	code, err := this.generateCode()
	if err != nil{
		return "", fmt.Errorf("Failed to generate guest code: %v", err)
	}

	// write to output
	outputFileName = this.outputDir+string(os.PathSeparator)+this.outputFilename+"_OpenFFIHost.go"
	err = ioutil.WriteFile( outputFileName, []byte(code), 0600)
	if err != nil{
		return "", fmt.Errorf("Failed to write host code to %v. Error: %v", this.outputDir+this.outputFilename, err)
	}

	return outputFileName, nil

}
//--------------------------------------------------------------------
func (this *HostCompiler) parseHeader() (string, error){
	tmp, err := template.New("host").Parse(HostHeaderTemplate)
	if err != nil{
		return "", fmt.Errorf("Failed to parse HostHeaderTemplate: %v", err)
	}

	buf := strings.Builder{}
	err = tmp.Execute(&buf, this.def)

	return buf.String(), err
}
//--------------------------------------------------------------------
func (this *HostCompiler) parseImports() (string, error){


	importsCode := HostImports

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

		// add to imports code + remove from serialization code (so it can be appended to the rest of the file)
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

	return importsCode, nil
}
//--------------------------------------------------------------------
func (this *HostCompiler) parseForeignStubs() (string, error){

	tmp, err := template.New("host").Funcs(templatesFuncMap).Parse(HostFunctionStubsTemplate)
	if err != nil{
		return "", fmt.Errorf("Failed to parse HostFunctionStubsTemplate: %v", err)
	}

	buf := strings.Builder{}
	err = tmp.Execute(&buf, this.def)

	return buf.String(), err
}
//--------------------------------------------------------------------
func (this *HostCompiler) parsePackage() (string, error){
	tmp, err := template.New("host").Parse(HostPackageTemplate)
	if err != nil{
		return "", fmt.Errorf("Failed to parse HostFunctionStubsTemplate: %v", err)
	}

	PackageName := struct {
		Package string
	}{
		Package: "main",
	}

	if pckName, found := this.hostOptions["package"]; found{
		PackageName.Package = pckName
	}

	buf := strings.Builder{}
	err = tmp.Execute(&buf, &PackageName)

	return buf.String(), err
}
//--------------------------------------------------------------------
func (this *HostCompiler) generateCode() (string, error){

	header, err := this.parseHeader()
	if err != nil{ return "", err }

	packageDeclaration, err := this.parsePackage()
	if err != nil{
		return "", err
	}

	imports, err := this.parseImports()
	if err != nil{ return "", err }

	functionStubs, err := this.parseForeignStubs()
	if err != nil{ return "", err }

	res := header + packageDeclaration + imports + HostCImport + HostHelperFunctions + functionStubs

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

