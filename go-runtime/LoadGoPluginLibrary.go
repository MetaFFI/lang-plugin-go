package main

// #include <stdint.h>
import "C"
import (
	"fmt"
	"plugin"
)

var libs []*plugin.Plugin

func init(){
	libs = make([]*plugin.Plugin, 0)
}

//export LoadGoPluginLibrary
func LoadGoPluginLibrary(libPath *C.char, libPathLength C.uint32_t){

	lib := C.GoStringN(libPath, C.int(libPathLength))

	fmt.Printf("Trying to load: "+lib+"\n")

	p, err := plugin.Open(lib)
	if err != nil{ panic(err)}

	libs = append(libs, p)
}