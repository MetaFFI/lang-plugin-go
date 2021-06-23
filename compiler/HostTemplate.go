package main
import "C"

const HostHeaderTemplate = `
// Code generated by OpenFFI. DO NOT EDIT.
// Host code for {{.IDLFilenameWithExtension}}
`

const HostPackageTemplate = `package {{.Package}}
`


const HostImports = `
import "fmt"
import "unsafe"
`

const HostCImportTemplate = `
/*
#cgo !windows LDFLAGS: -L. -ldl
#cgo CFLAGS: -I{{GetEnvVar "OPENFFI_HOME"}}

#include <stdlib.h>
#include <stdint.h>
#include <include/language_plugin_helpers.cpp>
*/
import "C"
`

const HostHelperFunctions = `
func init(){
	C.xllr_handle = nil

	err := C.load_args_helpers()
	if err != nil{
		panic("Failed to load OpenFFI XLLR functions: "+C.GoString(err))
	}
}
`

const HostFunctionStubsTemplate = `
{{ $pfn := .IDLFilename}}
{{range $mindex, $m := .Modules}}

// Code to call foreign functions in module {{$m.Name}} via XLLR
{{range $findex, $f := $m.Functions}}
// Call to foreign {{$f.PathToForeignFunction.function}}
{{if $f.Comment}}/*
{{$f.Comment}}
*/{{end}}
{{range $index, $elem := $f.Parameters}}
{{if $elem.Comment}}// {{$elem.Name}} - {{$elem.Comment}}{{end}}{{end}}
var {{$f.PathToForeignFunction.function}}_id int64 = -1
func {{AsPublic $f.PathToForeignFunction.function}}({{range $index, $elem := $f.Parameters}}{{if $index}},{{end}} {{$elem.Name}} {{if $elem.IsArray}}[]{{end}}{{if $elem.InnerTypes}}*{{end}}{{$elem.Type}}{{if eq $elem.Type "map"}}[{{$elem.MapKeyType}}]{{$elem.MapValueType}}{{end}}{{end}}) ({{range $index, $elem := $f.ReturnValues}}{{if $index}},{{end}}{{$elem.Name}} {{if $elem.IsArray}}[]{{end}}{{if $elem.InnerTypes}}*{{end}}{{$elem.Type}}{{if eq $elem.Type "map"}}[{{$elem.MapKeyType}}]{{$elem.MapValueType}}{{end}}{{end}}{{if $f.ReturnValues}},{{end}} err error){

	if {{$f.PathToForeignFunction.function}}_id == -1{

		// load function (no need to use a lock)
		runtime_plugin := "xllr.{{$m.TargetLanguage}}"
		pruntime_plugin := C.CString(runtime_plugin)
		defer C.free(unsafe.Pointer(pruntime_plugin))

		path := "{{$f.PathToForeignFunctionAsString}}"
		ppath := C.CString(path)
		defer C.free(unsafe.Pointer(ppath))

		var out_err *C.char
		var out_err_len C.uint32_t
		out_err_len = C.uint32_t(0)
		{{$f.PathToForeignFunction.function}}_id = int64(C.load_function(pruntime_plugin, C.uint(len(runtime_plugin)), ppath, C.uint(len(path)), C.int64_t(-1), &out_err, &out_err_len))
		
		if {{$f.PathToForeignFunction.function}}_id == -1{ // failed
			err = fmt.Errorf("Failed to load function %v: %v", "{{$f.PathToForeignFunction.function}}", string(C.GoBytes(unsafe.Pointer(out_err), C.int(out_err_len))))
			return
		}
	}

	paramsBufferLength := C.uint64_t({{CalculateArgsLength $f.Parameters}})
	paramsBuffer := C.alloc_args_buffer(C.int(paramsBufferLength))

	returnValuesBufferLength := C.uint64_t({{CalculateArgsLength $f.ReturnValues}})
	returnValuesBuffer := C.alloc_args_buffer(C.int(returnValuesBufferLength))
	
	// convert parameters to C
	{{$paramIndex := 0}}
	{{range $index, $elem := $f.Parameters}}
	{{ConvertToCHost $elem "in" $paramIndex "paramsBuffer"}}
	{{$fieldSize := CalculateArgLength $elem }}{{$paramIndex = Add $paramIndex $fieldSize}}
	{{end}}

	var out_err *C.char
	var out_err_len C.uint64_t
	out_err_len = C.uint64_t(0)

	C.call(pruntime_plugin, C.uint(len(runtime_plugin)),
			C.int64_t({{$f.PathToForeignFunction.function}}_id),
			paramsBuffer, paramsBufferLength,
			returnValuesBuffer, returnValuesBufferLength,
			&out_err, &out_err_len)

	// check errors
	if out_err_len != 0{
		err = fmt.Errorf("Function failed. Error: %v", string(C.GoBytes(unsafe.Pointer(out_err), C.int(out_err_len))))
		return
	}

	// convert from C to Go
	{{$retIndex := 0}}
	{{range $index, $elem := $f.ReturnValues}}
	{{ConvertToGo $elem "out" "ret" $retIndex "returnValuesBuffer"}}
	{{$fieldSize := CalculateArgLength $elem }}{{$retIndex = Add $retIndex $fieldSize}}
	{{end}}
	
	return {{range $index, $elem := $f.ReturnValues}}{{if $index}},{{end}}ret_{{$elem.Name}},{{end}} nil
}
{{end}}
{{end}}

`

