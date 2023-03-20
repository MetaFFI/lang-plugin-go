package main

import (
	"github.com/GreenFuze/go-parser"
	"github.com/MetaFFI/plugin-sdk/compiler/go/IDL"
	"strings"
)

var Imports map[string]bool

//--------------------------------------------------------------------
type GoIDLCompiler struct {
	goSourceCode         string
	goSourceCodeFilePath string
	gofile               *parser.GoFile
	
	idl *IDL.IDLDefinition
}

//--------------------------------------------------------------------
func NewGoIDLCompiler() *GoIDLCompiler {
	
	Imports = make(map[string]bool)
	
	// get Go file AST
	return &GoIDLCompiler{}
}

//--------------------------------------------------------------------
func (this *GoIDLCompiler) ParseIDL(goSourceCode string, gofilepath string, isEmbeddedCode bool) (*IDL.IDLDefinition, bool, error) {
	
	this.goSourceCode = goSourceCode
	this.goSourceCodeFilePath = gofilepath
	
	var err error
	this.gofile, err = parser.ParseSource(this.goSourceCode, this.goSourceCodeFilePath, true)
	if err != nil {
		return nil, true, err
	}
	this.idl = IDL.NewIDLDefinition(this.goSourceCodeFilePath, "go")
	
	globals := ExtractGlobals(this.gofile, this.idl.MetaFFIGuestLib)
	classes := ExtractClasses(this.gofile, this.idl.MetaFFIGuestLib)
	functions := ExtractFunctions(this.gofile, this.idl.MetaFFIGuestLib)
	
	// parse AST and build IDLDefinition
	
	module := IDL.NewModuleDefinition("go")
	
	for imp, _ := range Imports {
		module.AddExternalResource(imp)
	}
	
	importPath, _, err := this.gofile.ImportPath()
	if err != nil {
		return nil, true, err
	}

	importPath = strings.Replace(importPath, "\\", "/", -1)

	//module.AddExternalResource(importPath)

	
	module.AddGlobals(globals)
	module.AddFunctions(functions)
	
	for _, c := range classes {
		module.AddClass(c)
	}
	
	this.idl.AddModule(module)
	
	module.SetFunctionPath("package", this.gofile.Package)
	module.SetFunctionPath("module", importPath)
	
	this.idl.FinalizeConstruction()

	_, _ = this.idl.ToJSON()

	return this.idl, true, nil
}

//--------------------------------------------------------------------
