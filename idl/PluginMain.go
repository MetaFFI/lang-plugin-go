package main

import (
	"github.com/MetaFFI/lang-plugin-go/idl/IDLCompiler"
	compiler "github.com/MetaFFI/plugin-sdk/compiler/go"
)

import "C"

// --------------------------------------------------------------------
//
//export init_plugin
func init_plugin() {
	compiler.CreateIDLPluginInterfaceHandler(IDLCompiler.NewGoIDLCompiler())
}

//--------------------------------------------------------------------

func main() {}
