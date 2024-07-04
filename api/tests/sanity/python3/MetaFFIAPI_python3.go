package main

import (
	"fmt"
	"os"

	"github.com/MetaFFI/lang-plugin-go/api"
	"github.com/MetaFFI/plugin-sdk/compiler/go/IDL"
)

var runtime *api.MetaFFIRuntime
var mod *api.MetaFFIModule

func main() {
	fmt.Println("Loading python311 runtime")
	runtime = api.NewMetaFFIRuntime("python311")
	err := runtime.LoadRuntimePlugin()
	if err != nil {
		panic(err)
	}

	defer func() {
		fmt.Println("Going to releasing runtime")
		err = runtime.ReleaseRuntimePlugin()
		fmt.Println("Released runtime")
		if err != nil {
			panic(err)
		}
	}()

	fmt.Println("Loading test_target.py")
	mod, err = runtime.LoadModule("test_target.py")
	if err != nil {
		panic(err)
	}

	fmt.Println("Running HelloWorld")
	TestHelloWorld()

	fmt.Println("Running ReturnsAnError")
	TestReturnsAnError()

	fmt.Println("Running DivIntegers")
	TestDivIntegers()

	fmt.Println("Running JoinStrings")
	TestJoinStrings()

	fmt.Println("Running WaitABit")
	TestWaitABit()

	fmt.Println("Running TestMapGetSet")
	TestTestMapGetSet()

	fmt.Println("Running TestmapName")
	TestTestmapName()
}

func TestHelloWorld() {

	hellowWorld, err := mod.Load(`callable=hello_world`, nil, nil)
	if err != nil {
		panic(err)
	}

	_, err = hellowWorld()
	if err != nil {
		panic(err)
	}
}

func TestReturnsAnError() {

	returnsAnError, err := mod.Load(`callable=returns_an_error`, nil, nil)
	if err != nil {
		panic(err)
	}

	_, err = returnsAnError()
	if err == nil {
		panic("Expected an error")
	}
}

func TestDivIntegers() {
	divIntegers, err := mod.Load(`callable=div_integers`, []IDL.MetaFFIType{IDL.INT64, IDL.INT64}, []IDL.MetaFFIType{IDL.FLOAT32})
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

	if res[0].(float64) != float64(2) {
		fmt.Printf("Expected 2, got: %v", res[0].(float32))
		os.Exit(1)
	}
}

func TestJoinStrings() {
	divIntegers, err := mod.Load(`callable=join_strings`, []IDL.MetaFFIType{IDL.STRING8_ARRAY}, []IDL.MetaFFIType{IDL.STRING8})
	if err != nil {
		panic(err)
	}

	res, err := divIntegers([]string{"one", "two", "three"})
	if err != nil {
		panic(err)
	}

	if len(res) != 1 {
		fmt.Printf("Expected 1 result. Got %v results", len(res))
	}

	if res[0].(string) != "one,two,three" {
		fmt.Printf("Expected \"one,two,three\", got: %v", res[0].(string))
	}
}

func TestWaitABit() {

	getFiveSeconds, err := mod.Load(`attribute=five_seconds,getter`, nil, []IDL.MetaFFIType{IDL.INT64})
	if err != nil {
		panic(err)
	}

	waitABit, err := mod.Load(`callable=wait_a_bit`, []IDL.MetaFFIType{IDL.INT64}, nil)
	if err != nil {
		panic(err)
	}

	fiveSeconds, err := getFiveSeconds()
	if err != nil {
		panic(err)
	}

	_, err = waitABit(fiveSeconds[0].(int64))
	if err != nil {
		panic(err)
	}

}

func TestTestMapGetSet() {
	newTestMap, err := mod.Load(`callable=testmap`, nil, []IDL.MetaFFIType{IDL.HANDLE})
	if err != nil {
		panic(err)
	}

	res, err := newTestMap()
	if err != nil {
		panic(err)
	}

	testMap := res[0]

	testmapSet, err := mod.Load(`callable=testmap.set,instance_required`, []IDL.MetaFFIType{IDL.HANDLE, IDL.STRING8, IDL.ANY}, nil)
	if err != nil {
		panic(err)
	}

	_, err = testmapSet(testMap, "key1", int64(42))
	if err != nil {
		panic(err)
	}

	testmapGet, err := mod.Load(`callable=testmap.get,instance_required`, []IDL.MetaFFIType{IDL.HANDLE, IDL.STRING8}, []IDL.MetaFFIType{IDL.ANY})
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
	newTestMap, err := mod.Load(`callable=testmap`, nil, []IDL.MetaFFIType{IDL.HANDLE})
	if err != nil {
		panic(err)
	}

	res, err := newTestMap()
	if err != nil {
		panic(err)
	}

	testMap := res[0]

	testmapGet, err := mod.Load(`attribute=name,getter,instance_required`, []IDL.MetaFFIType{IDL.HANDLE}, []IDL.MetaFFIType{IDL.STRING8})
	if err != nil {
		panic(err)
	}

	testmapSet, err := mod.Load(`attribute=name,setter,instance_required`, []IDL.MetaFFIType{IDL.HANDLE, IDL.STRING8}, nil)
	if err != nil {
		panic(err)
	}

	name, err := testmapGet(testMap)
	if err != nil {
		panic(err)
	}

	if name[0].(string) != "name1" {
		fmt.Printf("Expected name1 ; Received: %v", name[0].(string))
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
