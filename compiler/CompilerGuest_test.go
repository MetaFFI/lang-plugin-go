package main

import (
	idl "github.com/MetaFFI/lang-plugin-go/idl/IDLCompiler"
	"github.com/MetaFFI/plugin-sdk/compiler/go/IDL"
	"os"
	"testing"
)

const runtime_test_json = `{"idl_source":"TestRuntime","idl_extension":".go","idl_filename_with_extension":"TestRuntime.go","idl_full_path":"/src/github.com/MetaFFI/lang-plugin-go/runtime/test/TestRuntime.go","metaffi_guest_lib":"TestRuntime_MetaFFIGuest","target_language":"go","modules":[{"name":"TestRuntime","comment":"","tags":{},"functions":[{"name":"HelloWorld","comment":"","tags":{},"function_path":{"entrypoint_function":"EntryPoint_HelloWorld","metaffi_guest_lib":"TestRuntime_MetaFFIGuest","module":"/src/github.com/MetaFFI/lang-plugin-go/runtime/test/","package":"TestRuntime"},"parameters":[],"return_values":[],"overload_index":0},{"name":"ReturnsAnError","comment":"","tags":{},"function_path":{"entrypoint_function":"EntryPoint_ReturnsAnError","metaffi_guest_lib":"TestRuntime_MetaFFIGuest","module":"C:/src/github.com/MetaFFI/lang-plugin-go/runtime/test/","package":"TestRuntime"},"parameters":[],"return_values":[],"overload_index":0},{"name":"DivIntegers","comment":"","tags":{},"function_path":{"entrypoint_function":"EntryPoint_DivIntegers","metaffi_guest_lib":"TestRuntime_MetaFFIGuest","module":"C:/src/github.com/MetaFFI/lang-plugin-go/runtime/test/","package":"TestRuntime"},"parameters":[{"name":"x","type":"int64","type_alias":"int","comment":"","tags":{},"dimensions":0,"is_optional":false},{"name":"y","type":"int64","type_alias":"int","comment":"","tags":{},"dimensions":0,"is_optional":false}],"return_values":[{"name":"r0","type":"float32","type_alias":"","comment":"","tags":{},"dimensions":0,"is_optional":false}],"overload_index":0},{"name":"JoinStrings","comment":"","tags":{},"function_path":{"entrypoint_function":"EntryPoint_JoinStrings","metaffi_guest_lib":"TestRuntime_MetaFFIGuest","module":"C:/src/github.com/MetaFFI/lang-plugin-go/runtime/test/","package":"TestRuntime"},"parameters":[{"name":"arrs","type":"string8_array","type_alias":"string","comment":"","tags":{},"dimensions":1,"is_optional":false}],"return_values":[{"name":"r0","type":"string8","type_alias":"","comment":"","tags":{},"dimensions":0,"is_optional":false}],"overload_index":0},{"name":"WaitABit","comment":"","tags":{},"function_path":{"entrypoint_function":"EntryPoint_WaitABit","metaffi_guest_lib":"TestRuntime_MetaFFIGuest","module":"C:/src/github.com/MetaFFI/lang-plugin-go/runtime/test/","package":"TestRuntime"},"parameters":[{"name":"d","type":"int64","type_alias":"time.Duration","comment":"","tags":{},"dimensions":0,"is_optional":false}],"return_values":[{"name":"r0","type":"handle","type_alias":"","comment":"","tags":{},"dimensions":0,"is_optional":false}],"overload_index":0},{"name":"GetSomeClasses","comment":"","tags":{},"function_path":{"entrypoint_function":"EntryPoint_GetSomeClasses","metaffi_guest_lib":"TestRuntime_MetaFFIGuest","module":"C:/src/github.com/MetaFFI/lang-plugin-go/runtime/test/","package":"TestRuntime"},"parameters":[],"return_values":[{"name":"r0","type":"handle_array","type_alias":"","comment":"","tags":{},"dimensions":1,"is_optional":false}],"overload_index":0},{"name":"ExpectThreeSomeClasses","comment":"","tags":{},"function_path":{"entrypoint_function":"EntryPoint_ExpectThreeSomeClasses","metaffi_guest_lib":"TestRuntime_MetaFFIGuest","module":"C:/src/github.com/MetaFFI/lang-plugin-go/runtime/test/","package":"TestRuntime"},"parameters":[{"name":"arr","type":"handle_array","type_alias":"SomeClass","comment":"","tags":{},"dimensions":1,"is_optional":false}],"return_values":[],"overload_index":0},{"name":"ExpectThreeBuffers","comment":"","tags":{},"function_path":{"entrypoint_function":"EntryPoint_ExpectThreeBuffers","metaffi_guest_lib":"TestRuntime_MetaFFIGuest","module":"C:/src/github.com/MetaFFI/lang-plugin-go/runtime/test/","package":"TestRuntime"},"parameters":[{"name":"buffers","type":"uint8_array","type_alias":"byte","comment":"","tags":{},"dimensions":2,"is_optional":false}],"return_values":[],"overload_index":0},{"name":"GetThreeBuffers","comment":"","tags":{},"function_path":{"entrypoint_function":"EntryPoint_GetThreeBuffers","metaffi_guest_lib":"TestRuntime_MetaFFIGuest","module":"C:/src/github.com/MetaFFI/lang-plugin-go/runtime/test/","package":"TestRuntime"},"parameters":[],"return_values":[{"name":"r0","type":"uint8_array","type_alias":"","comment":"","tags":{},"dimensions":2,"is_optional":false}],"overload_index":0},{"name":"NewTestMap","comment":"","tags":{},"function_path":{"entrypoint_function":"EntryPoint_NewTestMap","metaffi_guest_lib":"TestRuntime_MetaFFIGuest","module":"C:/src/github.com/MetaFFI/lang-plugin-go/runtime/test/","package":"TestRuntime"},"parameters":[],"return_values":[{"name":"r0","type":"handle","type_alias":"","comment":"","tags":{},"dimensions":0,"is_optional":false}],"overload_index":0}],"classes":[{"name":"SomeClass","comment":"","tags":{},"function_path":{"module":"C:/src/github.com/MetaFFI/lang-plugin-go/runtime/test/","package":"TestRuntime"},"constructors":[],"release":{"name":"ReleaseSomeClass","comment":"Releases object","tags":{},"function_path":{"entrypoint_function":"EntryPoint_SomeClass_ReleaseSomeClass","metaffi_guest_lib":"TestRuntime_MetaFFIGuest","module":"C:/src/github.com/MetaFFI/lang-plugin-go/runtime/test/","package":"TestRuntime"},"parameters":[{"name":"this_instance","type":"handle","type_alias":"","comment":"","tags":{},"dimensions":0,"is_optional":false}],"return_values":[],"overload_index":0,"instance_required":true},"methods":[{"name":"Print","comment":"","tags":{},"function_path":{"entrypoint_function":"EntryPoint_SomeClass_Print","metaffi_guest_lib":"TestRuntime_MetaFFIGuest","module":"C:/src/github.com/MetaFFI/lang-plugin-go/runtime/test/","package":"TestRuntime"},"parameters":[{"name":"this_instance","type":"handle","type_alias":"","comment":"","tags":{},"dimensions":0,"is_optional":false}],"return_values":[],"overload_index":0,"instance_required":true}],"fields":[]},{"name":"TestMap","comment":"","tags":{},"function_path":{"module":"C:/src/github.com/MetaFFI/lang-plugin-go/runtime/test/","package":"TestRuntime"},"constructors":[],"release":{"name":"ReleaseTestMap","comment":"Releases object","tags":{},"function_path":{"entrypoint_function":"EntryPoint_TestMap_ReleaseTestMap","metaffi_guest_lib":"TestRuntime_MetaFFIGuest","module":"C:/src/github.com/MetaFFI/lang-plugin-go/runtime/test/","package":"TestRuntime"},"parameters":[{"name":"this_instance","type":"handle","type_alias":"","comment":"","tags":{},"dimensions":0,"is_optional":false}],"return_values":[],"overload_index":0,"instance_required":true},"methods":[{"name":"Set","comment":"","tags":{"receiver_pointer":"true"},"function_path":{"entrypoint_function":"EntryPoint_TestMap_Set","metaffi_guest_lib":"TestRuntime_MetaFFIGuest","module":"C:/src/github.com/MetaFFI/lang-plugin-go/runtime/test/","package":"TestRuntime"},"parameters":[{"name":"this_instance","type":"handle","type_alias":"","comment":"","tags":{},"dimensions":0,"is_optional":false},{"name":"k","type":"string8","type_alias":"","comment":"","tags":{},"dimensions":0,"is_optional":false},{"name":"v","type":"any","type_alias":"interface{}","comment":"","tags":{},"dimensions":0,"is_optional":false}],"return_values":[],"overload_index":0,"instance_required":true},{"name":"Get","comment":"","tags":{"receiver_pointer":"true"},"function_path":{"entrypoint_function":"EntryPoint_TestMap_Get","metaffi_guest_lib":"TestRuntime_MetaFFIGuest","module":"C:/src/github.com/MetaFFI/lang-plugin-go/runtime/test/","package":"TestRuntime"},"parameters":[{"name":"this_instance","type":"handle","type_alias":"","comment":"","tags":{},"dimensions":0,"is_optional":false},{"name":"k","type":"string8","type_alias":"","comment":"","tags":{},"dimensions":0,"is_optional":false}],"return_values":[{"name":"r0","type":"any","type_alias":"","comment":"","tags":{},"dimensions":0,"is_optional":false}],"overload_index":0,"instance_required":true},{"name":"Contains","comment":"","tags":{"receiver_pointer":"true"},"function_path":{"entrypoint_function":"EntryPoint_TestMap_Contains","metaffi_guest_lib":"TestRuntime_MetaFFIGuest","module":"C:/src/github.com/MetaFFI/lang-plugin-go/runtime/test/","package":"TestRuntime"},"parameters":[{"name":"this_instance","type":"handle","type_alias":"","comment":"","tags":{},"dimensions":0,"is_optional":false},{"name":"k","type":"string8","type_alias":"","comment":"","tags":{},"dimensions":0,"is_optional":false}],"return_values":[{"name":"r0","type":"bool","type_alias":"","comment":"","tags":{},"dimensions":0,"is_optional":false}],"overload_index":0,"instance_required":true}],"fields":[{"name":"Name","type":"string8","type_alias":"string","comment":"","tags":{},"dimensions":0,"is_optional":false,"getter":{"name":"GetName","comment":"","tags":{"receiver_pointer":"true"},"function_path":{"entrypoint_function":"EntryPoint_TestMap_GetName","metaffi_guest_lib":"TestRuntime_MetaFFIGuest","module":"C:/src/github.com/MetaFFI/lang-plugin-go/runtime/test/","package":"TestRuntime"},"parameters":[{"name":"this_instance","type":"handle","type_alias":"","comment":"","tags":{},"dimensions":0,"is_optional":false}],"return_values":[{"name":"Name","type":"string8","type_alias":"string","comment":"","tags":{},"dimensions":0,"is_optional":false}],"overload_index":0,"instance_required":true},"setter":{"name":"SetName","comment":"","tags":{"receiver_pointer":"true"},"function_path":{"entrypoint_function":"EntryPoint_TestMap_SetName","metaffi_guest_lib":"TestRuntime_MetaFFIGuest","module":"C:/src/github.com/MetaFFI/lang-plugin-go/runtime/test/","package":"TestRuntime"},"parameters":[{"name":"this_instance","type":"handle","type_alias":"","comment":"","tags":{},"dimensions":0,"is_optional":false},{"name":"Name","type":"string8","type_alias":"string","comment":"","tags":{},"dimensions":0,"is_optional":false}],"return_values":[],"overload_index":0,"instance_required":true}}]}],"globals":[{"name":"FiveSeconds","type":"int64","type_alias":"time.Duration","comment":"","tags":{},"dimensions":0,"is_optional":false,"getter":{"name":"GetFiveSeconds","comment":"","tags":{},"function_path":{"entrypoint_function":"EntryPoint_GetFiveSeconds","metaffi_guest_lib":"TestRuntime_MetaFFIGuest","module":"C:/src/github.com/MetaFFI/lang-plugin-go/runtime/test/","package":"TestRuntime"},"parameters":[],"return_values":[{"name":"FiveSeconds","type":"int64","type_alias":"time.Duration","comment":"","tags":{},"dimensions":0,"is_optional":false}],"overload_index":0},"setter":null}],"external_resources":["time"]}]}`

// --------------------------------------------------------------------
func TestRuntimeTestFileGuest(t *testing.T) {

	idlDef, err := IDL.NewIDLDefinitionFromJSON(runtime_test_json)
	if err != nil {
		t.Fatal(err)
		return
	}

	_ = os.RemoveAll("temp_runtime_test_guest")
	err = os.Mkdir("temp_runtime_test_guest", 0700)
	if err != nil {
		t.Fatal(err)
		return
	}

	defer func() {
		err = os.RemoveAll("temp_runtime_test_guest")
		if err != nil {
			t.Fatal(err)
			return
		}
	}()

	cmp := NewGuestCompiler()
	err = cmp.Compile(idlDef, "temp_runtime_test_guest", "", nil)
	if err != nil {
		t.Fatal(err)
		return
	}
}

// --------------------------------------------------------------------
func TestTextTemplateGuest(t *testing.T) {
	idlCompiler := idl.NewGoIDLCompiler()
	idlDef, _, err := idlCompiler.ParseIDL("", "text/template")
	if err != nil {
		t.Fatal(err)
		return
	}

	_ = os.RemoveAll("temp_guest")
	err = os.Mkdir("temp_guest", 0700)
	if err != nil {
		t.Fatal(err)
		return
	}

	defer func() {
		err = os.RemoveAll("temp_guest")
		if err != nil {
			t.Fatal(err)
			return
		}
	}()

	cmp := NewGuestCompiler()
	err = cmp.Compile(idlDef, "temp_guest", "", nil)
	if err != nil {
		t.Fatal(err)
		return
	}
}

//--------------------------------------------------------------------
