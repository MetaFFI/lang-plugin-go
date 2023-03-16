package main

import (
	compiler "github.com/MetaFFI/plugin-sdk/compiler/go"
	"unicode"
)

import "C"

//--------------------------------------------------------------------
func IsPublic(name string) bool{
	if name == ""{
		return false
	}

	return unicode.IsUpper(rune(name[0]))
}
//--------------------------------------------------------------------
//export init_plugin
func init_plugin(){
	compiler.CreateIDLPluginInterfaceHandler(NewGoIDLCompiler())
}
//--------------------------------------------------------------------

func main(){}
