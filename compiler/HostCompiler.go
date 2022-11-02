package main

import (
	"fmt"
	IDL "github.com/MetaFFI/plugin-sdk/compiler/go/IDL"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"io/ioutil"
	"os"
	"strings"
	"text/template"
)

//--------------------------------------------------------------------
type HostCompiler struct {
	def            *IDL.IDLDefinition
	outputDir      string
	hostOptions    map[string]string
	outputFilename string
}

//--------------------------------------------------------------------
func NewHostCompiler() *HostCompiler {
	return &HostCompiler{}
}

//--------------------------------------------------------------------
func (this *HostCompiler) Compile(definition *IDL.IDLDefinition, outputDir string, outputFilename string, hostOptions map[string]string) (err error) {

	if outputFilename == ""{
        outputFilename = definition.IDLFilename
    }

	if strings.Contains(outputFilename, "#") {
		toRemove := outputFilename[strings.LastIndex(outputFilename, string(os.PathSeparator))+1 : strings.Index(outputFilename, "#")+1]
		outputFilename = strings.ReplaceAll(outputFilename, toRemove, "")
	}
	
	this.def = definition
	this.outputDir = outputDir
	this.hostOptions = hostOptions
	this.outputFilename = outputFilename

	caser := cases.Title(language.Und, cases.NoLower)
	this.def.ReplaceKeywords(map[string]string{
		"type":  caser.String("type"),
		"class": caser.String("class"),
		"func":  caser.String("func"),
		"var":   caser.String("var"),
		"const": caser.String("const"),
	})
	
	// generate code
	code, err := this.generateCode()
	if err != nil {
		return fmt.Errorf("Failed to generate host code: %v", err)
	}

	// write to output
	genOutputFilename := this.outputDir + string(os.PathSeparator) + this.outputFilename + "_MetaFFIHost.go"
	err = ioutil.WriteFile(genOutputFilename, []byte(code), 0600)
	if err != nil {
		return fmt.Errorf("Failed to write host code to %v. Error: %v", this.outputDir+this.outputFilename, err)
	}
	
	return nil
}

//--------------------------------------------------------------------
func (this *HostCompiler) parseHeader() (string, error) {
	tmp, err := template.New("HostHeaderTemplate").Parse(HostHeaderTemplate)
	if err != nil {
		return "", fmt.Errorf("Failed to parse HostHeaderTemplate: %v", err)
	}
	
	buf := strings.Builder{}
	err = tmp.Execute(&buf, this.def)
	
	return buf.String(), err
}

//--------------------------------------------------------------------
func (this *HostCompiler) parseImports() (string, error) {
	
	tmp, err := template.New("HostImportsTemplate").Funcs(templatesFuncMap).Parse(HostImportsTemplate)
	if err != nil {
		return "", fmt.Errorf("Failed to parse Go HostImportsTemplate: %v", err)
	}
	
	buf := strings.Builder{}
	err = tmp.Execute(&buf, this.def)
	
	return buf.String(), err
}

//--------------------------------------------------------------------
func (this *HostCompiler) parseCImports() (string, error) {
	
	tmp, err := template.New("HostCImportTemplate").Funcs(templatesFuncMap).Parse(HostCImportTemplate)
	if err != nil {
		return "", fmt.Errorf("Failed to parse Go HostCImportTemplate: %v", err)
	}
	
	buf := strings.Builder{}
	err = tmp.Execute(&buf, this.def)
	
	return buf.String(), err
}

//--------------------------------------------------------------------
func (this *HostCompiler) parseForeignStubs() (string, error) {
	
	tmp, err := template.New("Go HostFunctionStubsTemplate").Funcs(templatesFuncMap).Parse(HostFunctionStubsTemplate)
	if err != nil {
		return "", fmt.Errorf("Failed to parse Go HostFunctionStubsTemplate: %v", err)
	}
	
	buf := strings.Builder{}
	err = tmp.Execute(&buf, this.def)
	
	return buf.String(), err
}

//--------------------------------------------------------------------
func (this *HostCompiler) parsePackage() (string, error) {
	tmp, err := template.New("HostPackageTemplate").Funcs(templatesFuncMap).Parse(HostPackageTemplate)
	if err != nil {
		return "", fmt.Errorf("Failed to parse Go HostPackageTemplate: %v", err)
	}
	
	PackageName := struct {
		Package string
	}{
		Package: "main",
	}
	
	if pckName, found := this.hostOptions["package"]; found {
		PackageName.Package = pckName
	}
	
	buf := strings.Builder{}
	err = tmp.Execute(&buf, &PackageName)
	
	return buf.String(), err
}

//--------------------------------------------------------------------
func (this *HostCompiler) parseHelper() (string, error) {
	tmp, err := template.New("HostHelperFunctions").Funcs(templatesFuncMap).Parse(GetHostHelperFunctions())
	if err != nil {
		return "", fmt.Errorf("Failed to parse Go HostHelperFunctions: %v", err)
	}
	
	buf := strings.Builder{}
	err = tmp.Execute(&buf, this.def)
	
	return buf.String(), err
}

//--------------------------------------------------------------------
func (this *HostCompiler) generateCode() (string, error) {
	
	header, err := this.parseHeader()
	if err != nil {
		return "", err
	}
	
	packageDeclaration, err := this.parsePackage()
	if err != nil {
		return "", err
	}
	
	imports, err := this.parseImports()
	if err != nil {
		return "", err
	}
	
	cimports, err := this.parseCImports()
	if err != nil {
		return "", err
	}
	
	helper, err := this.parseHelper()
	if err != nil {
		return "", err
	}
	
	functionStubs, err := this.parseForeignStubs()
	if err != nil {
		return "", err
	}
	
	res := header + packageDeclaration + imports + cimports + helper + functionStubs
	
	return res, nil
}

//--------------------------------------------------------------------
