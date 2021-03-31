package main

import (
	"fmt"
	compiler "github.com/OpenFFI/plugin-sdk/compiler/go"
	"go/parser"
	"go/token"
	"html/template"
	"io/ioutil"
	"regexp"
	"strings"
)

//--------------------------------------------------------------------
type HostCompiler struct{
	def *compiler.IDLDefinition
	outputDir string
	serializationCode map[string]string
	outputFilename string
}
//--------------------------------------------------------------------
func NewHostCompiler(definition *compiler.IDLDefinition, outputDir string, outputFilename string, serializationCode map[string]string) *HostCompiler{

	serializationCodeCopy := make(map[string]string)
	for k, v := range serializationCode{
		serializationCodeCopy[k] = v
	}

	return &HostCompiler{def: definition, outputDir: outputDir, serializationCode: serializationCodeCopy, outputFilename: outputFilename}
}
//--------------------------------------------------------------------
func (this *HostCompiler) Compile() (outputFileName string, err error){

	// generate code
	code, err := this.generateCode()
	if err != nil{
		return "", fmt.Errorf("Failed to generate guest code: %v", err)
	}

	// write to output
	err = ioutil.WriteFile(this.outputDir+this.outputFilename, []byte(code), 0600)
	if err != nil{
		return "", fmt.Errorf("Failed to write host code to %v. Error: %v", this.outputDir+this.outputFilename, err)
	}

	return this.outputFilename, nil

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

		fset := token.NewFileSet()
		ast, err := parser.ParseFile(fset, "", code, parser.ImportsOnly)
		if err != nil{
			return "", fmt.Errorf("Failed to parse serialization code of file %v. Error: %v", filename, err)
		}

		// add to imports code + remove from serialization code (so it can be appended to the rest of the file)
		for _, i := range ast.Imports{
			importsCode += "import \""+i.Path.Value+"\"\n"

			removeImportRegex, err := regexp.Compile(fmt.Sprintf("import[ ]+\"%v\"", i.Path.Value))
			if err != nil{
				return "", fmt.Errorf("Failed to create regex to remove imports from serialization code: \"%v\". Error: %v", fmt.Sprintf("import[ ]+\"%v\"", i.Path.Value), err)
			}

			this.serializationCode[filename] = removeImportRegex.ReplaceAllString(this.serializationCode[filename], "")
		}
	}

	return importsCode, nil
}
//--------------------------------------------------------------------
func (this *HostCompiler) parseForeignStubs() (string, error){
	tmp, err := template.New("host").Parse(HostFunctionStubsTemplate)
	if err != nil{
		return "", fmt.Errorf("Failed to parse HostFunctionStubsTemplate: %v", err)
	}

	buf := strings.Builder{}
	err = tmp.Execute(&buf, this.def)

	return buf.String(), err
}
//--------------------------------------------------------------------
func (this *HostCompiler) generateCode() (string, error){

	header, err := this.parseHeader()
	if err != nil{ return "", err }

	imports, err := this.parseImports()
	if err != nil{ return "", err }

	functionStubs, err := this.parseForeignStubs()
	if err != nil{ return "", err }

	res := header + imports + HostCImport + HostMainFunction + HostHelperFunctions + functionStubs

	// append serialization code in the same file
	for _, serializationCode := range this.serializationCode{
		res += serializationCode
	}

	return res, nil
}
//--------------------------------------------------------------------

