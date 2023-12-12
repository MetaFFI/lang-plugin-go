package openjdk

import (
	"github.com/MetaFFI/lang-plugin-go/api"
	"github.com/MetaFFI/plugin-sdk/compiler/go/IDL"
	"os"
	"testing"
)

var runtime *api.MetaFFIRuntime
var testRuntimeModule *api.MetaFFIModule
var testMapModule *api.MetaFFIModule

func TestMain(m *testing.M) {
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

	code := m.Run()

	err = runtime.ReleaseRuntimePlugin()
	if err != nil {
		panic(err)
	}

	os.Exit(code)
}

func TestHelloWorld(t *testing.T) {

	hellowWorld, err := testRuntimeModule.Load(`class=sanity.TestRuntime,callable=helloWorld`, nil, nil)
	if err != nil {
		t.Fatal(err)
	}

	_, err = hellowWorld()
	if err != nil {
		t.Fatal(err)
	}
}

func TestReturnsAnError(t *testing.T) {

	returnsAnError, err := testRuntimeModule.Load(`class=sanity.TestRuntime,callable=returnsAnError`, nil, nil)
	if err != nil {
		t.Fatal(err)
	}

	_, err = returnsAnError()
	if err == nil {
		t.Fatal("Expected an error")
	}
}

func TestDivIntegers(t *testing.T) {
	divIntegers, err := testRuntimeModule.Load(`class=sanity.TestRuntime,callable=divIntegers`, []IDL.MetaFFIType{IDL.INT32, IDL.INT32}, []IDL.MetaFFIType{IDL.FLOAT32})
	if err != nil {
		t.Fatal(err)
	}

	res, err := divIntegers(10, 5)
	if err != nil {
		t.Fatal(err)
	}

	if len(res) != 1 {
		t.Fatalf("Expected 1 result. Got %v results", len(res))
	}

	if res[0].(float32) != float32(2) {
		t.Fatalf("Expected 2, got: %v", res[0].(float32))
	}
}

func TestJoinStrings(t *testing.T) {
	divIntegers, err := testRuntimeModule.Load(`class=sanity.TestRuntime,callable=joinStrings`, []IDL.MetaFFIType{IDL.STRING8_ARRAY}, []IDL.MetaFFIType{IDL.STRING8})
	if err != nil {
		t.Fatal(err)
	}

	res, err := divIntegers([]string{"one", "two", "three"})
	if err != nil {
		t.Fatal(err)
	}

	if len(res) != 1 {
		t.Fatalf("Expected 1 result. Got %v results", len(res))
	}

	if res[0].(string) != "one,two,three" {
		t.Fatalf("Expected \"one,two,three\", got: %v", res[0].(string))
	}
}

func TestWaitABit(t *testing.T) {

	getFiveSeconds, err := testRuntimeModule.Load(`class=sanity.TestRuntime,field=fiveSeconds,getter`, nil, []IDL.MetaFFIType{IDL.INT32})
	if err != nil {
		t.Fatal(err)
	}

	waitABit, err := testRuntimeModule.Load(`class=sanity.TestRuntime,callable=waitABit`, []IDL.MetaFFIType{IDL.INT32}, nil)
	if err != nil {
		t.Fatal(err)
	}

	fiveSeconds, err := getFiveSeconds()
	if err != nil {
		t.Fatal(err)
	}

	_, err = waitABit(fiveSeconds[0].(int32))
	if err != nil {
		t.Fatal(err)
	}

}

func TestTestMapGetSet(t *testing.T) {
	newTestMap, err := testRuntimeModule.Load(`class=sanity.TestMap,callable=<init>`, nil, []IDL.MetaFFIType{IDL.HANDLE})
	if err != nil {
		t.Fatal(err)
	}

	res, err := newTestMap()
	if err != nil {
		t.Fatal(err)
	}

	testMap := res[0]

	testmapSet, err := testRuntimeModule.Load(`class=sanity.TestMap,callable=set,instance_required`, []IDL.MetaFFIType{IDL.HANDLE, IDL.STRING8, IDL.ANY}, nil)
	if err != nil {
		t.Fatal(err)
	}

	_, err = testmapSet(testMap, "key1", int64(42))
	if err != nil {
		t.Fatal(err)
	}

	testmapContains, err := testRuntimeModule.Load(`class=sanity.TestMap,callable=contains,instance_required`, []IDL.MetaFFIType{IDL.HANDLE, IDL.STRING8}, []IDL.MetaFFIType{IDL.BOOL})
	if err != nil {
		t.Fatal(err)
	}

	isContained, err := testmapContains(testMap, "key1")
	if err != nil {
		t.Fatal(err)
	}

	if !isContained[0].(bool) {
		t.Fatalf("Expected to return true")
	}

	testmapGet, err := testRuntimeModule.Load(`class=sanity.TestMap,callable=get,instance_required`, []IDL.MetaFFIType{IDL.HANDLE, IDL.STRING8}, []IDL.MetaFFIType{IDL.ANY})
	if err != nil {
		t.Fatal(err)
	}

	ret, err := testmapGet(testMap, "key1")
	if err != nil {
		t.Fatal(err)
	}

	if len(ret) != 1 {
		t.Fatalf("Expected 1 result. Got %v results", len(ret))
	}

	if ret[0].(int64) != 42 {
		t.Fatalf("Expected 42. Got %v", ret[0].(int64))
	}
}

func TestTestmapName(t *testing.T) {
	newTestMap, err := testRuntimeModule.Load(`class=sanity.TestMap,callable=<init>`, nil, []IDL.MetaFFIType{IDL.HANDLE})
	if err != nil {
		t.Fatal(err)
	}

	res, err := newTestMap()
	if err != nil {
		t.Fatal(err)
	}

	testMap := res[0]

	testmapGet, err := testRuntimeModule.Load(`class=sanity.TestMap,field=name,instance_required,getter`, []IDL.MetaFFIType{IDL.HANDLE}, []IDL.MetaFFIType{IDL.STRING8})
	if err != nil {
		t.Fatal(err)
	}

	testmapSet, err := testRuntimeModule.Load(`class=sanity.TestMap,field=name,instance_required,setter`, []IDL.MetaFFIType{IDL.HANDLE, IDL.STRING8}, nil)
	if err != nil {
		t.Fatal(err)
	}

	name, err := testmapGet(testMap)
	if err != nil {
		t.Fatal(err)
	}

	if name[0].(string) != "" {
		t.Fatalf("Expected empty ; Received: %v", name[0].(string))
	}

	_, err = testmapSet(testMap, "name is my name")
	if err != nil {
		t.Fatal(err)
	}

	name1, err := testmapGet(testMap)
	if err != nil {
		t.Fatal(err)
	}

	if name1[0].(string) != "name is my name" {
		t.Fatalf("Expected \"name is my name\" ; Received: %v", name1[0].(string))
	}

}
