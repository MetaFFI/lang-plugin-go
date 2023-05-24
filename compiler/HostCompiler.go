package main

import (
	_ "embed"
	"fmt"
	IDL "github.com/MetaFFI/plugin-sdk/compiler/go/IDL"
	"github.com/MetaFFI/plugin-sdk/compiler/go"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"io/ioutil"
	"os"
	"runtime"
	"strings"
	"text/template"
)

//go:embed MetaFFIGoHostCommon.gotpl
var metaFFIGoHostCommon string

var goKeywords = map[string]bool{
    "break": true, "default": true, "func": true, "interface": true, "select": true,
    "case": true, "defer": true, "go": true, "map": true, "struct": true,
    "chan": true, "else": true, "goto": true, "package": true, "switch": true,
    "const": true, "fallthrough": true, "if": true, "range": true, "type": true,
    "continue": true, "for": true, "import": true, "return": true, "var": true,
    "string": true, "int8": true, "int16": true, "int32": true, "int64": true,
    "uint8": true, "uint16": true, "uint32": true, "uint64": true,
    "float32": true, "float64": true, "bool": true,
}

// --------------------------------------------------------------------
type HostCompiler struct {
	def            *IDL.IDLDefinition
	outputDir      string
	hostOptions    map[string]string
	outputFilename string
}

// --------------------------------------------------------------------
func NewHostCompiler() *HostCompiler {
	return &HostCompiler{}
}

// --------------------------------------------------------------------
func (this *HostCompiler) getMetaFFIGoHostCommon(commonPackageName string) string {

	metaffiHome := os.Getenv("METAFFI_HOME")
	if metaffiHome == "" {
		panic("METAFFI_HOME environment variable is not set")
	}
	metaffiHome = strings.ReplaceAll(metaffiHome, "\\", "/")

	os := runtime.GOOS
	var longtype string
	switch os {
	case "windows":
		longtype = "ulonglong"
	default:
		longtype = "ulong"
	}

	p := struct {
		Package     string
		MetaFFIHome string
		LongType    string
	}{
		Package:     commonPackageName,
		MetaFFIHome: metaffiHome,
		LongType:    longtype,
	}

	tmp, err := template.New("metaFFIGoHostCommon").Parse(metaFFIGoHostCommon)
	if err != nil {
		panic(fmt.Errorf("Failed to parse HostHeaderTemplate: %v", err))
	}

	buf := strings.Builder{}
	err = tmp.Execute(&buf, p)
	if err != nil {
		panic(err)
	}

	return buf.String()

}

// --------------------------------------------------------------------
func (this *HostCompiler) Compile(definition *IDL.IDLDefinition, outputDir string, outputFilename string, hostOptions map[string]string) (err error) {

	// make sure definition does not use "go syntax-keywords" as names. If so, change the names a bit...
	compiler.ModifyKeywords(definition, goKeywords, func(keyword string)string{ return keyword+"__" })

	if outputFilename == "" {
		outputFilename = definition.IDLFilename
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
	code, packageName, err := this.generateCode()
	if err != nil {
		return fmt.Errorf("Failed to generate host code: %v", err)
	}

	// TODO: handle multiple modules

	_ = os.Mkdir(this.outputDir+string(os.PathSeparator)+strings.ToLower(this.def.Modules[0].Name), 0777)

	// write MetaFFIGoHostCommon
	err = ioutil.WriteFile(this.outputDir+string(os.PathSeparator)+strings.ToLower(this.def.Modules[0].Name)+string(os.PathSeparator)+"MetaFFIGoHostCommon.go", []byte(this.getMetaFFIGoHostCommon(packageName)), 0600)
	if err != nil {
		return fmt.Errorf("Failed to write host code to %v. Error: %v", this.outputDir+this.outputFilename, err)
	}

	// write to output
	genOutputFilename := this.outputDir + string(os.PathSeparator) + strings.ToLower(this.def.Modules[0].Name) + string(os.PathSeparator) + this.outputFilename + "_MetaFFIHost.go"
	err = ioutil.WriteFile(genOutputFilename, []byte(code), 0600)
	if err != nil {
		return fmt.Errorf("Failed to write host code to %v. Error: %v", this.outputDir+this.outputFilename, err)
	}

	return nil
}

// --------------------------------------------------------------------
func (this *HostCompiler) parseHeader() (string, error) {
	tmp, err := template.New("HostHeaderTemplate").Parse(HostHeaderTemplate)
	if err != nil {
		return "", fmt.Errorf("Failed to parse HostHeaderTemplate: %v", err)
	}

	buf := strings.Builder{}
	err = tmp.Execute(&buf, this.def)

	return buf.String(), err
}

// --------------------------------------------------------------------
func (this *HostCompiler) parseImports() (string, error) {

	tmp, err := template.New("HostImportsTemplate").Funcs(templatesFuncMap).Parse(HostImportsTemplate)
	if err != nil {
		return "", fmt.Errorf("Failed to parse Go HostImportsTemplate: %v", err)
	}

	buf := strings.Builder{}
	err = tmp.Execute(&buf, this.def)

	return buf.String(), err
}

// --------------------------------------------------------------------
func (this *HostCompiler) parseCImports() (string, error) {

	tmp, err := template.New("HostCImportTemplate").Funcs(templatesFuncMap).Parse(HostCImportTemplate)
	if err != nil {
		return "", fmt.Errorf("Failed to parse Go HostCImportTemplate: %v", err)
	}

	buf := strings.Builder{}
	err = tmp.Execute(&buf, this.def)

	return buf.String(), err
}

// --------------------------------------------------------------------
func (this *HostCompiler) parseForeignStubs() (string, error) {

	tmp, err := template.New("Go HostFunctionStubsTemplate").Funcs(templatesFuncMap).Parse(HostFunctionStubsTemplate)
	if err != nil {
		return "", fmt.Errorf("Failed to parse Go HostFunctionStubsTemplate: %v", err)
	}

	buf := strings.Builder{}
	err = tmp.Execute(&buf, this.def)

	return buf.String(), err
}

// --------------------------------------------------------------------
func (this *HostCompiler) parsePackage() (code string, packageName string, err error) {
	tmp, err := template.New("HostPackageTemplate").Funcs(templatesFuncMap).Parse(HostPackageTemplate)
	if err != nil {
		return "", "", fmt.Errorf("Failed to parse Go HostPackageTemplate: %v", err)
	}

	PackageName := struct {
		Package string
	}{
		Package: this.def.Modules[0].Name, // TODO: support multiple modules
	}

	if pckName, found := this.hostOptions["package"]; found {
		PackageName.Package = pckName
	}

	buf := strings.Builder{}
	err = tmp.Execute(&buf, &PackageName)

	return buf.String(), PackageName.Package, err
}

// --------------------------------------------------------------------
func (this *HostCompiler) parseHelper() (string, error) {
	tmp, err := template.New(GetHostHelperFunctionsName()).Funcs(templatesFuncMap).Parse(GetHostHelperFunctions())
	if err != nil {
		return "", fmt.Errorf("Failed to parse Go %v: %v", GetHostHelperFunctionsName(), err)
	}

	buf := strings.Builder{}
	err = tmp.Execute(&buf, this.def)

	return buf.String(), err
}

// --------------------------------------------------------------------
func (this *HostCompiler) generateCode() (code string, packageName string, err error) {

	header, err := this.parseHeader()
	if err != nil {
		return "", "", err
	}

	packageDeclaration, packageName, err := this.parsePackage()
	if err != nil {
		return "", "", err
	}

	imports, err := this.parseImports()
	if err != nil {
		return "", "", err
	}

	cimports, err := this.parseCImports()
	if err != nil {
		return "", "", err
	}

	helper, err := this.parseHelper()
	if err != nil {
		return "", "", err
	}

	functionStubs, err := this.parseForeignStubs()
	if err != nil {
		return "", "", err
	}

	res := header + packageDeclaration + imports + cimports + helper + functionStubs

	return res, packageName, err
}

//--------------------------------------------------------------------
