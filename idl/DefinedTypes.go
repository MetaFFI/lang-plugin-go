package main

import (
	"go/ast"
	"go/token"
)

var DefinedTypes *definedTypes

type definedTypes struct {
	types map[string]bool
	file *ast.File
}
//--------------------------------------------------------------------
func InitDefinedTypes(file *ast.File){
	DefinedTypes = NewDefinedTypes(file)
}
//--------------------------------------------------------------------
func NewDefinedTypes(file *ast.File) *definedTypes {

	this := &definedTypes{ types: make(map[string]bool), file: file }
	this.fillDefinedTypes()

	return this
}
//--------------------------------------------------------------------
func (this *definedTypes) fillDefinedTypes(){
	ast.Inspect(this.file, this.visitor)
}
//--------------------------------------------------------------------
func (this *definedTypes) visitor(node ast.Node) bool{

	if gendecl, ok := node.(*ast.GenDecl); ok{
		if gendecl.Tok == token.TYPE{
			for _, s := range gendecl.Specs{
				if typeSpec, ok := s.(*ast.TypeSpec); ok{
					if typeSpec.Name != nil{
						this.types[typeSpec.Name.Name] = true
					}
				}
			}
		}
	}

	return true
}
//--------------------------------------------------------------------
func (this *definedTypes) isDefined(t string) bool{
	_, found := this.types[t]
	return found
}
//--------------------------------------------------------------------