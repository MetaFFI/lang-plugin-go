package main

import "C"
import . "github.com/MetaFFI/plugin-sdk/compiler/go"

//export init_plugin
func init_plugin() {
	CreateIDLPluginInterfaceHandler(NewGoIDLCompiler())
}
func main() {}
