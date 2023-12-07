package python3

import (
	"github.com/MetaFFI/lang-plugin-go/api"
	"github.com/MetaFFI/plugin-sdk/compiler/go/IDL"
	"os"
	"testing"
)

var runtime *api.MetaFFIRuntime
var mod *api.MetaFFIModule

func TestMain(m *testing.M) {
	runtime = api.NewMetaFFIRuntime("python3")
	err := runtime.LoadRuntimePlugin()
	if err != nil {
		panic(err)
	}

	mod, err = runtime.LoadModule("test_target.py")
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

	hellowWorld, err := mod.LoadCallable(`callable=hello_world`, nil, nil)
	if err != nil {
		t.Fatal(err)
	}

	_, err = hellowWorld()
	if err != nil {
		t.Fatal(err)
	}
}

func TestReturnsAnError(t *testing.T) {

	returnsAnError, err := mod.LoadCallable(`callable=returns_an_error`, nil, nil)
	if err != nil {
		t.Fatal(err)
	}

	_, err = returnsAnError()
	if err == nil {
		t.Fatal("Expected an error")
	}
}

func TestDivIntegers(t *testing.T) {
	divIntegers, err := mod.LoadCallable(`callable=div_integers`, []IDL.MetaFFIType{IDL.INT64, IDL.INT64}, []IDL.MetaFFIType{IDL.FLOAT32})
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
	divIntegers, err := mod.LoadCallable(`callable=join_strings`, []IDL.MetaFFIType{IDL.STRING8_ARRAY}, []IDL.MetaFFIType{IDL.STRING8})
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

	getFiveSeconds, err := mod.LoadCallable(`attribute=five_seconds,getter`, nil, []IDL.MetaFFIType{IDL.INT64})
	if err != nil {
		t.Fatal(err)
	}

	waitABit, err := mod.LoadCallable(`callable=wait_a_bit`, []IDL.MetaFFIType{IDL.INT64}, nil)
	if err != nil {
		t.Fatal(err)
	}

	fiveSeconds, err := getFiveSeconds()
	if err != nil {
		t.Fatal(err)
	}

	_, err = waitABit(fiveSeconds[0].(int64))
	if err != nil {
		t.Fatal(err)
	}

}

func TestTestMapGetSet(t *testing.T) {
	newTestMap, err := mod.LoadCallable(`callable=testmap`, nil, []IDL.MetaFFIType{IDL.HANDLE})
	if err != nil {
		t.Fatal(err)
	}

	res, err := newTestMap()
	if err != nil {
		t.Fatal(err)
	}

	testMap := res[0]

	testmapSet, err := mod.LoadCallable(`callable=testmap.set,instance_required`, []IDL.MetaFFIType{IDL.HANDLE, IDL.STRING8, IDL.ANY}, nil)
	if err != nil {
		t.Fatal(err)
	}

	_, err = testmapSet(testMap, "key1", int64(42))
	if err != nil {
		t.Fatal(err)
	}

	testmapGet, err := mod.LoadCallable(`callable=testmap.get,instance_required`, []IDL.MetaFFIType{IDL.HANDLE, IDL.STRING8}, []IDL.MetaFFIType{IDL.ANY})
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
	newTestMap, err := mod.LoadCallable(`callable=testmap`, nil, []IDL.MetaFFIType{IDL.HANDLE})
	if err != nil {
		t.Fatal(err)
	}

	res, err := newTestMap()
	if err != nil {
		t.Fatal(err)
	}

	testMap := res[0]

	testmapGet, err := mod.LoadCallable(`attribute=name,getter,instance_required`, []IDL.MetaFFIType{IDL.HANDLE}, []IDL.MetaFFIType{IDL.STRING8})
	if err != nil {
		t.Fatal(err)
	}

	testmapSet, err := mod.LoadCallable(`attribute=name,setter,instance_required`, []IDL.MetaFFIType{IDL.HANDLE, IDL.STRING8}, nil)
	if err != nil {
		t.Fatal(err)
	}

	name, err := testmapGet(testMap)
	if err != nil {
		t.Fatal(err)
	}

	if name[0].(string) != "name1" {
		t.Fatalf("Expected name1 ; Received: %v", name[0].(string))
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
