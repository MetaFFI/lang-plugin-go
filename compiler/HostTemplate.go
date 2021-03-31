package main

const HostHeaderTemplate = `
// Code generated by OpenFFI. DO NOT EDIT.
// Guest code for {{.IDLFilenameWithExtension}}
{{ $pfn := .IDLFilename}}

package main
`

const HostImports = `
import "fmt"
import "unsafe"
import "github.com/golang/protobuf/proto"
import "runtime"
`

const HostMainFunction = `
func main(){} // main function must be declared to create dynamic library
`

const HostCImport = `
/*
#cgo !windows LDFLAGS: -L. -ldl

#include <stdlib.h>
#include <stdint.h>
void* xllr_handle = NULL;
void (*pcall)(const char*, uint32_t,
			 const char*, uint32_t,
			 const char*, uint32_t,
			 unsigned char*, uint64_t,
			 unsigned char**, uint64_t*,
			 unsigned char**, uint64_t*,
			 uint8_t*) = NULL;

#ifdef _WIN32 //// --- START WINDOWS ---
#include <Windows.h>
void get_last_error_string(DWORD err, char** out_err_str)
{
    DWORD bufLen = FormatMessage(FORMAT_MESSAGE_ALLOCATE_BUFFER |
                                 FORMAT_MESSAGE_FROM_SYSTEM |
							     FORMAT_MESSAGE_IGNORE_INSERTS,
							     NULL,
							     err,
						         MAKELANGID(LANG_NEUTRAL, SUBLANG_DEFAULT),
						         (LPTSTR) out_err_str,
						         0,
						         NULL );

    // TODO: out_err_str should get cleaned up!
}

void* load_library(const char* name, char** out_err)
{
	void* handle = LoadLibraryA(name);
	if(!handle)
	{
		get_last_error_string(GetLastError(), out_err);
	}

	return handle;
}

const char* free_library(void* lib) // return error string. null if no error.
{
	if(!lib)
	{
		return NULL;
	}

	if(!FreeLibrary(lib))
	{
		char* out_err;
		get_last_error_string(GetLastError(), &out_err);
		return out_err;
	}

	return NULL;
}

void* load_symbol(void* handle, const char* name, char** out_err)
{
	void* res = GetProcAddress(handle, name);
	if(!res)
	{
		get_last_error_string(GetLastError(), out_err);
		return NULL;
	}

	return res;
}

#else // ------ START POSIX ----
#include <dlfcn.h>
void* load_library(const char* name, char** out_err)
{
	void* handle = dlopen(name, RTLD_NOW);
	if(!handle)
	{
		*out_err = dlerror();
	}

	return handle;
}

const char* free_library(void* lib)
{
	if(dlclose(lib))
	{
		return dlerror();
	}

	return NULL;
}

void* load_symbol(void* handle, const char* name, char** out_err)
{
	void* res = dlsym(handle, name);
	if(!res)
	{
		*out_err = dlerror();
		return NULL;
	}

	return res;
}

#endif // ------- END POSIX -----

void call(
		const char* runtime_plugin, uint32_t runtime_plugin_len,
		const char* module_name, uint32_t module_name_len,
		const char* func_name, uint32_t func_name_len,
		unsigned char* in_params, uint64_t in_params_len,
		unsigned char** out_params, uint64_t* out_params_len,
		unsigned char** out_ret, uint64_t* out_ret_len,
		uint8_t* is_error
)
{
	pcall(runtime_plugin, runtime_plugin_len,
			module_name, module_name_len,
			func_name, func_name_len,
			in_params, in_params_len,
			out_params, out_params_len,
			out_ret, out_ret_len,
			is_error);
}

const char* load_xllr_api()
{
	char* out_err = NULL;
	pcall = load_symbol(xllr_handle, "call", &out_err);
	if(!pcall)
	{
		return out_err;
	}

	return NULL;
}

*/
import "C"
`

const HostHelperFunctions = `
func freeXLLR() error{
	errstr := C.free_library(C.xllr_handle)

	if errstr != nil{
		return fmt.Errorf("Failed to free XLLR: %v", C.GoString(errstr))
	}

	return nil
}

// TODO: make sure it is called only once!
func loadXLLR() error{

	if C.xllr_handle != nil && C.pcall != nil{
        return nil
    }

	var name *C.char
	if runtime.GOOS == "darwin" {
		name = C.CString("xllr.dylib")
	}else if runtime.GOOS == "windows"{
		name = C.CString("xllr.dll")
	} else {
		name = C.CString("xllr.so")
	}

	defer C.free(unsafe.Pointer(name))

	// TODO: load all other exported functions from XLLR
	var out_err *C.char
	if C.xllr_handle = C.load_library(name, &out_err)
	C.xllr_handle == nil{ // error has occurred
		return fmt.Errorf("Failed to load XLLR: %v", C.GoString(out_err))
	}

	callstr := C.CString("call")
    defer C.free(unsafe.Pointer(callstr))
    if cerr := C.load_xllr_api(); cerr != nil{
        return fmt.Errorf("Failed to load call function: %v", C.GoString(cerr))
    }

	return nil
}
`

const HostFunctionStubsTemplate = `
{{range $mindex, $m := .Modules}}

// Code to call foreign functions in module {{$m.Name}} via XLLR
{{range $findex, $f := $m.Functions}}
// Call to foreign {{$f.PathToForeignFunction.function}}
{{if $f.Comment}}// {{$f.Comment}}{{end}}
{{range $index, $elem := $f.Parameters}}
{{if $elem.Comment}}// $elem.Name - $elem.Comment{{end}}{{end}}
func {{title $f.PathToForeignFunction.function}}({{range $index, $elem := $f.Parameters}}{{if $index}},{{end}} {{$elem.Name}} {{$elem.Type}}{{end}}) ({{range $index, $elem := $f.ReturnValues}}{{if $index}},{{end}}{{$elem.Name}} {{$elem.Type}}{{end}}{{if $f.ReturnValues}},{{end}} err error){

	// serialize parameters
	req := {{$f.ParametersType}}{}
	{{range $index, $elem := $f.Parameters}}
	req.{{$elem.Name}} = {{$elem.Name}}
	{{end}}

	// load XLLR
	err = loadXLLR()
	if err != nil{
		err = fmt.Errorf("Failed to marshal return values into protobuf. Error: %v", err)
		return
	}
	
	// call function
	runtime_plugin := "xllr.{{$m.TargetLanguage}}"
	pruntime_plugin := C.CString(runtime_plugin)
	defer C.free(unsafe.Pointer(pruntime_plugin))

	module_name := "{{$pfn}}OpenFFIGuest"
	pmodule_name := C.CString(module_name)
	defer C.free(unsafe.Pointer(pmodule_name))

	func_name := "Foreign{{$f.Name}}"
	pfunc_name := C.CString(func_name)
	defer C.free(unsafe.Pointer(pfunc_name))

	// in parameters
	in_params, err := proto.Marshal(&req)
	if err != nil{
		err = fmt.Errorf("Failed to marshal return values into protobuf. Error: %v", err)
		return
	}

	var pin_params *C.uchar
	var in_params_len C.uint64_t
	if len(in_params) > 0{
		pin_params = (*C.uchar)(unsafe.Pointer(&in_params[0]))
		in_params_len = C.uint64_t(len(in_params))
	} else {
		in_params_len = C.uint64_t(0)
	}

	var out_ret *C.uchar
	var out_ret_len C.uint64_t
	out_ret_len = C.uint64_t(0)

	var out_params *C.uchar
	var out_params_len C.uint64_t
	out_params_len = C.uint64_t(0)

	var out_is_error C.uchar
	out_is_error = C.uchar(0)

	C.call(pruntime_plugin, C.uint(len(runtime_plugin)),
			pmodule_name, C.uint(len(module_name)),
			pfunc_name, C.uint(len(func_name)),
			pin_params, in_params_len,
			&out_params, &out_params_len,
			&out_ret, &out_ret_len,
			&out_is_error)

	// check errors
	if out_is_error != 0{
		err = fmt.Errorf("Function failed. Error: %v", string(C.GoBytes(unsafe.Pointer(out_ret), C.int(out_ret_len))))
		return
	}

	// deserialize result	
	ret := {{$f.ReturnValues}}{}
	out_ret_buf := C.GoBytes(unsafe.Pointer(out_ret), C.int(out_ret_len))
	err = proto.Unmarshal(out_ret_buf, &ret)
	if err != nil{
		err = fmt.Errorf("Failed to unmarshal return values into protobuf. Error: %v", err)
		return
	}

	return {{range $index, $elem := $f.ReturnValues}}{{if $index}},{{end}}ret.{{$elem.Name}}{{end}}{{if $f.ReturnValues}},{{end}} nil

}
{{end}}

{{end}}

`

