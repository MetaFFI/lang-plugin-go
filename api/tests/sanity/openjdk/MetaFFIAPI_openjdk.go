package main

import (
	"fmt"
	"os"

	"github.com/MetaFFI/lang-plugin-go/api"
	"github.com/MetaFFI/plugin-sdk/compiler/go/IDL"
)

var runtime *api.MetaFFIRuntime
var testRuntimeModule *api.MetaFFIModule
var testMapModule *api.MetaFFIModule

func main() {

	runtime = api.NewMetaFFIRuntime("openjdk")
	err := runtime.LoadRuntimePlugin()
	if err != nil {
		panic(err)
	}

	testRuntimeModule, err = runtime.LoadModule("./sanity/TestRuntime.class")
	if err != nil {
		panic(err)
	}

	testMapModule, err = runtime.LoadModule("./sanity/TestMap.class")
	if err != nil {
		panic(err)
	}

	TestHelloWorld()
	TestDivIntegers()
	TestJoinStrings()
	TestTestMapGetSet()
	TestTestmapName()
	TestWaitABit()

	err = runtime.ReleaseRuntimePlugin()
	if err != nil {
		panic(err)
	}

}

func TestHelloWorld() {
	hellowWorld, err := testRuntimeModule.Load(`class=sanity.TestRuntime,callable=helloWorld`, nil, nil)

	if err != nil {
		panic(err)
	}

	_, err = hellowWorld()
	if err != nil {
		panic(err)
	}
}

func TestReturnsAnError() {

	returnsAnError, err := testRuntimeModule.Load(`class=sanity.TestRuntime,callable=returnsAnError`, nil, nil)
	if err != nil {
		panic(err)
	}

	_, err = returnsAnError()
	if err == nil {
		panic("Expected an error")
	}
}

func TestDivIntegers() {

	divIntegers, err := testRuntimeModule.Load(`class=sanity.TestRuntime,callable=divIntegers`, []IDL.MetaFFIType{IDL.INT32, IDL.INT32}, []IDL.MetaFFIType{IDL.FLOAT32})

	if err != nil {
		panic(err)
	}

	res, err := divIntegers(10, 5)
	if err != nil {
		panic(err)
	}

	if len(res) != 1 {
		fmt.Printf("Expected 1 result. Got %v results", len(res))
		os.Exit(1)
	}

	if res[0].(float32) != float32(2) {
		fmt.Printf("Expected 2, got: %v", res[0].(float32))
		os.Exit(1)
	}

}

func TestJoinStrings() {

	joinStrings, err := testRuntimeModule.Load(`class=sanity.TestRuntime,callable=joinStrings`, []IDL.MetaFFIType{IDL.STRING8_ARRAY}, []IDL.MetaFFIType{IDL.STRING8})

	if err != nil {
		panic(err)
	}

	res, err := joinStrings([]string{"one", "two", "three"})
	if err != nil {
		panic(err)
	}

	if len(res) != 1 {
		fmt.Printf("Expected 1 result. Got %v results", len(res))
		os.Exit(1)
	}

	if res[0].(string) != "one,two,three" {
		fmt.Printf("Expected \"one,two,three\", got: %v", res[0].(string))
		os.Exit(1)
	}

}

func TestWaitABit() {

	getFiveSeconds, err := testRuntimeModule.Load(`class=sanity.TestRuntime,field=fiveSeconds,getter`, nil, []IDL.MetaFFIType{IDL.INT32})
	if err != nil {
		panic(err)
	}

	waitABit, err := testRuntimeModule.Load(`class=sanity.TestRuntime,callable=waitABit`, []IDL.MetaFFIType{IDL.INT32}, nil)
	if err != nil {
		panic(err)
	}

	fiveSeconds, err := getFiveSeconds()
	if err != nil {
		panic(err)
	}

	_, err = waitABit(fiveSeconds[0].(int32))
	if err != nil {
		panic(err)
	}

}

func TestTestMapGetSet() {

	newTestMap, err := testRuntimeModule.Load(`class=sanity.TestMap,callable=<init>`, nil, []IDL.MetaFFIType{IDL.HANDLE})

	if err != nil {
		panic(err)
	}

	res, err := newTestMap()
	if err != nil {
		panic(err)
	}

	testMap := res[0]

	testmapSet, err := testRuntimeModule.Load(`class=sanity.TestMap,callable=set,instance_required`, []IDL.MetaFFIType{IDL.HANDLE, IDL.STRING8, IDL.ANY}, nil)

	if err != nil {
		panic(err)
	}

	_, err = testmapSet(testMap, "key1", int64(42))
	if err != nil {
		panic(err)
	}

	testmapContains, err := testRuntimeModule.Load(`class=sanity.TestMap,callable=contains,instance_required`, []IDL.MetaFFIType{IDL.HANDLE, IDL.STRING8}, []IDL.MetaFFIType{IDL.BOOL})

	if err != nil {
		panic(err)
	}

	isContained, err := testmapContains(testMap, "key1")
	if err != nil {
		panic(err)
	}

	if !isContained[0].(bool) {
		fmt.Printf("Expected to return true")
		os.Exit(1)
	}

	testmapGet, err := testRuntimeModule.Load(`class=sanity.TestMap,callable=get,instance_required`, []IDL.MetaFFIType{IDL.HANDLE, IDL.STRING8}, []IDL.MetaFFIType{IDL.ANY})

	if err != nil {
		panic(err)
	}

	ret, err := testmapGet(testMap, "key1")
	if err != nil {
		panic(err)
	}

	if len(ret) != 1 {
		fmt.Printf("Expected 1 result. Got %v results", len(ret))
		os.Exit(1)
	}

	if ret[0].(int64) != 42 {
		fmt.Printf("Expected 42. Got %v", ret[0].(int64))
		os.Exit(1)
	}

}

func TestTestmapName() {
	newTestMap, err := testRuntimeModule.Load(`class=sanity.TestMap,callable=<init>`, nil, []IDL.MetaFFIType{IDL.HANDLE})
	if err != nil {
		panic(err)
	}

	res, err := newTestMap()
	if err != nil {
		panic(err)
	}

	testMap := res[0]

	testmapGet, err := testRuntimeModule.Load(`class=sanity.TestMap,field=name,instance_required,getter`, []IDL.MetaFFIType{IDL.HANDLE}, []IDL.MetaFFIType{IDL.STRING8})
	if err != nil {
		panic(err)
	}

	testmapSet, err := testRuntimeModule.Load(`class=sanity.TestMap,field=name,instance_required,setter`, []IDL.MetaFFIType{IDL.HANDLE, IDL.STRING8}, nil)
	if err != nil {
		panic(err)
	}

	name, err := testmapGet(testMap)
	if err != nil {
		panic(err)
	}

	if name[0].(string) != "" {
		fmt.Printf("Expected empty ; Received: %v", name[0].(string))
		os.Exit(1)
	}

	_, err = testmapSet(testMap, "name is my name")
	if err != nil {
		panic(err)
	}

	name1, err := testmapGet(testMap)
	if err != nil {
		panic(err)
	}

	if name1[0].(string) != "name is my name" {
		fmt.Printf("Expected \"name is my name\" ; Received: %v", name1[0].(string))
		os.Exit(1)
	}

}
