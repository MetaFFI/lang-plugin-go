package main

import (
	"github.com/MetaFFI/lang-plugin-go/idl/IDLCompiler"
	"github.com/MetaFFI/sdk/compiler/go/plugin"
)

import "C"

// --------------------------------------------------------------------
//
//export init_plugin
func init_plugin() {
	plugin.CreateIDLPluginInterfaceHandler(IDLCompiler.NewGoIDLCompiler())
}

//--------------------------------------------------------------------

func main() {}
