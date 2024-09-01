package main

import (
	"github.com/MetaFFI/lang-plugin-go/api"
	metaffi "github.com/MetaFFI/lang-plugin-go/go-runtime"
	"github.com/MetaFFI/plugin-sdk/compiler/go/IDL"
)

var runtime *api.MetaFFIRuntime

//--------------------------------------------------------------------

type Log4j struct {
	instance metaffi.MetaFFIHandle
	perror   func(...interface{}) ([]interface{}, error)
}

func NewLog4j(name string) (*Log4j, error) {
	this := &Log4j{}

	mod, err := runtime.LoadModule("log4j-api-2.21.1.jar;log4j-core-2.21.1.jar")
	if err != nil {
		return nil, err
	}

	params := []IDL.MetaFFITypeInfo{IDL.MetaFFITypeInfo{StringType: IDL.STRING8}}
	retvals := []IDL.MetaFFITypeInfo{IDL.MetaFFITypeInfo{StringType: IDL.HANDLE, Alias: "org.apache.logging.log4j.Logger"}}
	constructor, err := mod.LoadWithInfo("class=org.apache.logging.log4j.LogManager,callable=getLogger", params, retvals)
	if err != nil {
		return nil, err
	}

	instance, err := constructor(name)
	if err != nil {
		return nil, err
	}

	this.instance = instance[0].(metaffi.MetaFFIHandle)

	this.perror, err = mod.Load("class=org.apache.logging.log4j.Logger,callable=error,instance_required", []IDL.MetaFFIType{IDL.HANDLE, IDL.STRING8}, nil)
	if err != nil {
		return nil, err
	}

	return this, nil
}

func (this *Log4j) Error(msg string) error {
	_, err := this.perror(this.instance, msg)
	return err
}

//--------------------------------------------------------------------

func main() {

	// load runtime
	runtime = api.NewMetaFFIRuntime("openjdk")
	err := runtime.LoadRuntimePlugin()
	if err != nil {
		panic(err)
	}

	logger, err := NewLog4j("MyLogger!")
	if err != nil {
		panic(err)
	}

	err = logger.Error("This is an error from Go!")
	if err != nil {
		panic(err)
	}

}
