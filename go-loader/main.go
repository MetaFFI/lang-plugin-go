package main
import "C"
import (
	"os"
	"plugin"
	"runtime"
)

//import "C"

// stub
func main(){}

func init(){
	println("in loader init")
	loadOpenFFIGoRuntime()
	println("loaded XLLR go runtime")
}

var goRuntimeLib *plugin.Plugin

func loadOpenFFIGoRuntime(){
	openffiHome := os.Getenv("OPENFFI_HOME")
	if openffiHome == ""{
		panic("Cannot find OPENFFI_HOME environment variable")
	}

	var extension string
	switch runtime.GOOS{
	case "windows": extension = ".dll"
	case "darwin": extension = ".dylib"
	default: extension = ".so"
	}

	var err error
	goRuntimeLib, err = plugin.Open(openffiHome+string(os.PathSeparator)+"xllr.go.runtime"+extension)
	if err != nil{ panic(err) }
}