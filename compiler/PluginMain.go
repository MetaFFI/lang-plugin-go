package main

import (
	compiler "github.com/MetaFFI/sdk/compiler/go"
	"github.com/MetaFFI/sdk/compiler/go/plugin"
)

import "C"

//export init_plugin
func init_plugin() {
	plugin.PluginMain = plugin.NewLanguagePluginMain(compiler.NewHostCompiler(), compiler.NewGuestCompiler())
}
func main() {}
