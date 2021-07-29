package main

import (
	"fmt"
	compiler "github.com/OpenFFI/plugin-sdk/compiler/go"
	"io/ioutil"
	"os"
	"strings"
	"text/template"
)

//--------------------------------------------------------------------
type HostCompiler struct{
	def *compiler.IDLDefinition
	outputDir string
	hostOptions map[string]string
	outputFilename string
}
//--------------------------------------------------------------------
func NewHostCompiler(definition *compiler.IDLDefinition, outputDir string, outputFilename string, hostOptions map[string]string) *HostCompiler{

	return &HostCompiler{def: definition,
		outputDir: outputDir,
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


	tmp, err := template.New("host").Funcs(templatesFuncMap).Parse(HostImportsTemplate)
	if err != nil{
		return "", fmt.Errorf("Failed to parse HostFunctionStubsTemplate: %v", err)
	}

	buf := strings.Builder{}
	err = tmp.Execute(&buf, this.def)

	return buf.String(), err
}
//--------------------------------------------------------------------
func (this *HostCompiler) parseCImports() (string, error){

	tmp, err := template.New("host").Funcs(templatesFuncMap).Parse(HostCImportTemplate)
	if err != nil{
		return "", fmt.Errorf("Failed to parse HostFunctionStubsTemplate: %v", err)
	}

	buf := strings.Builder{}
	err = tmp.Execute(&buf, this.def)

	return buf.String(), err
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
	tmp, err := template.New("host").Funcs(templatesFuncMap).Parse(HostPackageTemplate)
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

	cimports, err := this.parseCImports()
	if err != nil{ return "", err }

	functionStubs, err := this.parseForeignStubs()
	if err != nil{ return "", err }

	res := header + packageDeclaration + imports + cimports + HostHelperFunctions + functionStubs

	return res, nil
}
//--------------------------------------------------------------------

