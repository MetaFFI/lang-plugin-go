package main

import "C"
import "runtime"

const HostHeaderTemplate = `
// Code generated by MetaFFI. DO NOT EDIT.
// Host code for {{.IDLFilenameWithExtension}}
`

const HostPackageTemplate = `package {{.Package}}
`

const HostImportsTemplate = `
import "fmt"
import "unsafe"
import . "github.com/MetaFFI/lang-plugin-go/go-runtime"
`

const HostCImportTemplate = `
/*
#cgo !windows LDFLAGS: -L. -ldl
#cgo CFLAGS: -I{{GetEnvVar "METAFFI_HOME" true}}

#include <stdlib.h>
#include <stdint.h>
#include <include/cdt_structs.h>
#include <include/cdt_capi_loader.h>

metaffi_handle get_null_handle();
metaffi_size get_int_item(metaffi_size* array, int index);
void* convert_union_to_ptr(void* p);
void set_cdt_type(struct cdt* p, metaffi_type t);
metaffi_type get_cdt_type(struct cdt* p);
metaffi_size len_to_metaffi_size(long long i);
*/
import "C"
`

func GetHostHelperFunctions() string {
	if runtime.GOOS == "windows" {
		return HostHelperFunctionsWindows
	} else {
		return HostHelperFunctionsNonWindows
	}
}

func GetHostHelperFunctionsName() string {
	if runtime.GOOS == "windows" {
		return "HostHelperFunctionsWindows"
	} else {
		return "HostHelperFunctionsNonWindows"
	}
}

const HostHelperFunctionsWindows = `
{{ $idl := . }}

// function IDs
{{range $mindex, $m := .Modules}}

{{range $findex, $f := $m.Globals}}
{{if $f.Getter}}var {{$f.Getter.Name}}_id unsafe.Pointer{{end}}
{{if $f.Setter}}var {{$f.Setter.Name}}_id unsafe.Pointer{{end}}
{{end}}{{/* End globals */}}

{{range $findex, $f := $m.Functions}}
var {{$f.Name}}_id unsafe.Pointer
{{end}}{{/* End Functions */}}

{{range $cindex, $c := $m.Classes}}
{{range $findex, $f := $c.Fields}}
{{if $f.Getter}}var {{$c.Name}}_{{$f.Getter.Name}}_id unsafe.Pointer{{end}}
{{if $f.Setter}}var {{$c.Name}}_{{$f.Setter.Name}}_id unsafe.Pointer{{end}}
{{end}}{{/* End Fields */}}
{{range $findex, $f := $c.Methods}}
var {{$c.Name}}_{{$f.Name}}_id unsafe.Pointer
{{end}}{{/* End Methods */}}
{{range $findex, $f := $c.Constructors}}
var {{$c.Name}}_{{$f.Name}}_id unsafe.Pointer
{{end}}{{/* End Constructor */}}
{{if $c.Releaser}}
var {{$c.Name}}_{{$c.Releaser.Name}}_id unsafe.Pointer
{{end}}{{/* End Releaser */}}
{{end}}{{/* End Classes */}}
{{end}}{{/* End modules */}}

func Load(modulePath string){
	loadCDTCAPI()

	runtime_plugin := "xllr.{{.TargetLanguage}}"
	pruntime_plugin := C.CString(runtime_plugin)
	runtime_plugin_length := C.uint32_t(len(runtime_plugin))

	// load foreign runtime
	var out_err *C.char
    var out_err_len C.uint32_t
    out_err_len = C.uint32_t(0)
	C.xllr_load_runtime_plugin(pruntime_plugin, runtime_plugin_length, &out_err, &out_err_len)
	if out_err_len != C.uint32_t(0){
		panic(fmt.Errorf("Failed to load runtime %v: %v", runtime_plugin, string(C.GoBytes(unsafe.Pointer(out_err), C.int(out_err_len)))))
	}

	// load functions
	loadFF := func(modulePath string, fpath string, params_count int, retval_count int) unsafe.Pointer{
		ppath := C.CString(fpath)
		defer C.free(unsafe.Pointer(ppath))

		pmodulePath := C.CString(modulePath)
		defer C.free(unsafe.Pointer(pmodulePath))

		var out_err *C.char
		var out_err_len C.uint32_t
		out_err_len = C.uint32_t(0)
		id := C.xllr_load_function(pruntime_plugin, runtime_plugin_length, pmodulePath, C.uint(len(modulePath)), ppath, C.uint(len(fpath)), nil,  C.schar(params_count), C.schar(params_count), &out_err, &out_err_len)

		if id == nil{ // failed
			panic(fmt.Errorf("Failed to load foreign entity entrypoint \"%v\": %v", fpath, string(C.GoBytes(unsafe.Pointer(out_err), C.int(out_err_len)))))
		}

		return id
	}

	{{range $mindex, $m := .Modules}}
	{{range $findex, $f := $m.Globals}}
	{{if $f.Getter}}{{$f.Getter.Name}}_id = loadFF(modulePath, `+"`"+`{{$f.Getter.FunctionPathAsString $idl}}`+"`"+`, {{len $f.Getter.Parameters}}, {{len $f.Getter.ReturnValues}}){{end}}
	{{if $f.Setter}}{{$f.Setter.Name}}_id = loadFF(modulePath, `+"`"+`{{$f.Setter.FunctionPathAsString $idl}}`+"`"+`, {{len $f.Setter.Parameters}}, {{len $f.Setter.ReturnValues}}){{end}}
	{{end}}{{/* End globals */}}

	{{range $findex, $f := $m.Functions}}
	{{$f.Name}}_id = loadFF(modulePath, `+"`"+`{{$f.FunctionPathAsString $idl}}`+"`"+`, {{len $f.Parameters}}, {{len $f.ReturnValues}})
	{{end}}{{/* End Functions */}}

	{{range $cindex, $c := $m.Classes}}
	{{range $findex, $f := $c.Fields}}
	{{if $f.Getter}}{{$c.Name}}_{{$f.Getter.Name}}_id = loadFF(modulePath, `+"`"+`{{$f.Getter.FunctionPathAsString $idl}}`+"`"+`, {{len $f.Getter.Parameters}}, {{len $f.Getter.ReturnValues}}){{end}}
	{{if $f.Setter}}{{$c.Name}}_{{$f.Setter.Name}}_id = loadFF(modulePath, `+"`"+`{{$f.Setter.FunctionPathAsString $idl}}`+"`"+`, {{len $f.Setter.Parameters}}, {{len $f.Setter.ReturnValues}}){{end}}
	{{end}}{{/* End Fields */}}
	{{range $findex, $f := $c.Methods}}
	{{$c.Name}}_{{$f.Name}}_id = loadFF(modulePath, `+"`"+`{{$f.FunctionPathAsString $idl}}`+"`"+`, {{len $f.Parameters}}, {{len $f.ReturnValues}})
	{{end}}{{/* End Methods */}}
	{{range $findex, $f := $c.Constructors}}
	{{$c.Name}}_{{$f.Name}}_id = loadFF(modulePath, `+"`"+`{{$f.FunctionPathAsString $idl}}`+"`"+`, {{len $f.Parameters}}, {{len $f.ReturnValues}})
	{{end}}{{/* End Constructor */}}
	{{if $c.Releaser}}
	{{$c.Name}}_{{$c.Releaser.Name}}_id = loadFF(modulePath, `+"`"+`{{$c.Releaser.FunctionPathAsString $idl}}`+"`"+`, {{len $c.Releaser.Parameters}}, {{len $c.Releaser.ReturnValues}})
	{{end}}{{/* End Releaser */}}
	{{end}}{{/* End Classes */}}
	{{end}}{{/* End modules */}}
}

`

const HostHelperFunctionsNonWindows = `

// function IDs
{{range $mindex, $m := .Modules}}

{{range $findex, $f := $m.Globals}}
{{if $f.Getter}}var {{$f.Getter.Name}}_id unsafe.Pointer{{end}}
{{if $f.Setter}}var {{$f.Setter.Name}}_id unsafe.Pointer{{end}}
{{end}}{{/* End globals */}}

{{range $findex, $f := $m.Functions}}
var {{$f.Name}}_id unsafe.Pointer
{{end}}{{/* End Functions */}}

{{range $cindex, $c := $m.Classes}}
{{range $findex, $f := $c.Fields}}
{{if $f.Getter}}var {{$c.Name}}_{{$f.Getter.Name}}_id unsafe.Pointer{{end}}
{{if $f.Setter}}var {{$c.Name}}_{{$f.Setter.Name}}_id unsafe.Pointer{{end}}
{{end}}{{/* End Fields */}}
{{range $findex, $f := $c.Methods}}
var {{$c.Name}}_{{$f.Name}}_id unsafe.Pointer
{{end}}{{/* End Methods */}}
{{range $findex, $f := $c.Constructors}}
var {{$c.Name}}_{{$f.Name}}_id unsafe.Pointer
{{end}}{{/* End Constructor */}}
{{if $c.Releaser}}
var {{$c.Name}}_{{$c.Releaser.Name}}_id unsafe.Pointer
{{end}}{{/* End Releaser */}}
{{end}}{{/* End Classes */}}
{{end}}{{/* End modules */}}

func Load(modulePath string){
	err := C.load_cdt_capi()
	if err != nil{
		panic("Failed to load MetaFFI XLLR functions: "+C.GoString(err))
	}

	runtime_plugin := "xllr.{{.TargetLanguage}}"
	pruntime_plugin := C.CString(runtime_plugin)
	runtime_plugin_length := C.uint32_t(len(runtime_plugin))

	// load foreign runtime
	var out_err *C.char
    var out_err_len C.uint32_t
    out_err_len = C.uint32_t(0)
	C.xllr_load_runtime_plugin(pruntime_plugin, runtime_plugin_length, &out_err, &out_err_len)
	if out_err_len != C.uint32_t(0){
		panic(fmt.Errorf("Failed to load runtime %v: %v", runtime_plugin, string(C.GoBytes(unsafe.Pointer(out_err), C.int(out_err_len)))))
	}

	// load functions
	loadFF := func(modulePath string, fpath string, params_count int, retval_count int) unsafe.Pointer{
		ppath := C.CString(fpath)
		defer C.free(unsafe.Pointer(ppath))

		pmodulePath := C.CString(modulePath)
        defer C.free(unsafe.Pointer(pmodulePath))

		var out_err *C.char
		var out_err_len C.uint32_t
		out_err_len = C.uint32_t(0)
		id := C.xllr_load_function(pruntime_plugin, runtime_plugin_length, pmodulePath, C.uint(len(modulePath)), ppath, C.uint(len(fpath)), nil, C.schar(params_count), C.schar(params_count), &out_err, &out_err_len)
		
		if id == nil{ // failed
			panic(fmt.Errorf("Failed to load foreign entity entrypoint \"%v\": %v", fpath, string(C.GoBytes(unsafe.Pointer(out_err), C.int(out_err_len)))))
		}

		return id
	}

	{{ $idl := . }}
	{{range $mindex, $m := .Modules}}
	{{range $findex, $f := $m.Globals}}
	{{if $f.Getter}}{{$f.Getter.Name}}_id = loadFF(modulePath, "{{$f.Getter.FunctionPathAsString $idl}}", {{len $f.Getter.Parameters}}, {{len $f.Getter.ReturnValues}} ){{end}}
	{{if $f.Setter}}{{$f.Setter.Name}}_id = loadFF(modulePath, "{{$f.Setter.FunctionPathAsString $idl}}", {{len $f.Setter.Parameters}}, {{len $f.Setter.ReturnValues}} ){{end}}
	{{end}}{{/* End globals */}}
	
	{{range $findex, $f := $m.Functions}}
	{{$f.Name}}_id = loadFF(modulePath, "{{$f.FunctionPathAsString $idl}}", {{len $f.Parameters}}, {{len $f.ReturnValues}})
	{{end}}{{/* End Functions */}}

	{{range $cindex, $c := $m.Classes}}
	{{range $findex, $f := $c.Fields}}
	{{if $f.Getter}}{{$c.Name}}_{{$f.Getter.Name}}_id = loadFF(modulePath, "{{$f.Getter.FunctionPathAsString $idl}}", {{len $f.Getter.Parameters}}, {{len $f.Getter.ReturnValues}}){{end}}
	{{if $f.Setter}}{{$c.Name}}_{{$f.Setter.Name}}_id = loadFF(modulePath, "{{$f.Setter.FunctionPathAsString $idl}}", {{len $f.Setter.Parameters}}, {{len $f.Setter.ReturnValues}}){{end}}
	{{end}}{{/* End Fields */}}
	{{range $findex, $f := $c.Methods}}
	{{$c.Name}}_{{$f.Name}}_id = loadFF(modulePath, "{{$f.FunctionPathAsString $idl}}", {{len $f.Parameters}}, {{len $f.ReturnValues}})
	{{end}}{{/* End Methods */}}
	{{range $findex, $f := $c.Constructors}}
	{{$c.Name}}_{{$f.Name}}_id = loadFF(modulePath, "{{$f.FunctionPathAsString $idl}}", {{len $f.Parameters}}, {{len $f.ReturnValues}})
	{{end}}{{/* End Constructor */}}
	{{if $c.Releaser}}
	{{$c.Name}}_{{$c.Releaser.Name}}_id = loadFF(modulePath, "{{$c.Releaser.FunctionPathAsString $idl}}", {{len $c.Releaser.Parameters}}, {{len $c.Releaser.ReturnValues}})
	{{end}}{{/* End Releaser */}}
	{{end}}{{/* End Classes */}}
	{{end}}{{/* End modules */}}
}


`

const HostFunctionStubsTemplate = `
{{ $pfn := .IDLFilename}}
{{ $idl := . }}
{{range $mindex, $m := .Modules}}

{{range $findex, $f := $m.Globals}}
{{if $f.Getter}}
{{if $f.Comment}}/*
{{$f.Comment}}
*/{{end}}
func {{ToGoNameConv $f.Getter.Name}}() (instance {{ConvertToGoType $f.ArgDefinition $m}}, err error){
	{{ $paramsLength := len $f.Getter.Parameters }}{{ $returnLength := len $f.Getter.ReturnValues }}
	{{GenerateCodeAllocateCDTS $f.Getter.Parameters $f.Getter.ReturnValues}}

	// parameters
	{{range $index, $elem := $f.Getter.Parameters}}
	fromGoToCDT({{$elem.Name}}, xcall_params, {{$index}})
	{{end}}{{/* End Parameters */}}

	{{GenerateCodeXCall "" $f.Getter.Name $f.Getter.Parameters $f.Getter.ReturnValues}}
	
	{{range $index, $elem := $f.Getter.ReturnValues}}
	{{$elem.Name}}AsInterface := fromCDTToGo(return_valuesCDTS, {{$index}})
	var {{$elem.Name}} {{ConvertToGoType $elem $m}}
	
	if {{$elem.Name}}AsInterface != nil{
		{{if $elem.IsAny}}
		{{$elem.Name}} = {{$elem.Name}}AsInterface
		
		{{else if not $elem.IsHandle}}
		{{$elem.Name}} = {{if $elem.IsTypeAlias}}{{if $elem.IsArray}}[]{{end}}{{GetTypeOrAlias $elem $m}}{{else}}{{ConvertToGoType $elem $m}}{{end}}({{$elem.Name}}AsInterface.({{ConvertToGoType $elem $m}}))
		{{else}} {{/* Handle */}}		
		{{/* Go object */}}
		
		{{if not $elem.IsArray}}
		if obj, ok := {{$elem.Name}}AsInterface.(Handle); ok{ // None Go object			
			{{$elem.Name}} = {{HandleNoneGoObject $elem $m}}
		} else {
			{{$elem.Name}} = {{if $elem.IsTypeAlias}}{{if $elem.IsArray}}[]{{end}}{{GetTypeOrAlias $elem $m}}{{else}}{{ConvertToGoType $elem $m}}{{end}}({{$elem.Name}}AsInterface.({{ConvertToGoType $elem $m}}))
		}
		{{else}}
		{{if $elem.IsTypeAlias}} {{/* a type is specified */}}
		if len({{$elem.Name}}AsInterface.([]interface{})) > 0{
			{{$elem.Name}} = make([]{{GetTypeOrAlias $elem $m}}, len({{$elem.Name}}AsInterface.([]interface{})))
			if _, ok := {{$elem.Name}}AsInterface.([]interface{})[0].(Handle); ok{
				for i, h := range {{$elem.Name}}AsInterface.([]interface{}){
					{{$elem.Name}}[i] = {{GetTypeOrAlias $elem $m}}{ h: h.(Handle) }
				}
			} else {
				for i, obj := range {{$elem.Name}}AsInterface.([]interface{}){
					{{$elem.Name}}[i] = obj.({{GetTypeOrAlias $elem $m}})
				}
			}
		}
		{{else}}
		{{$elem.Name}} = {{$elem.Name}}AsInterface
		{{end}}
		{{end}}

		{{end}}{{/* end handling types */}}
	}

	{{end}}{{/* End return values */}}

	return {{range $index, $elem := $f.Getter.ReturnValues}}{{if $index}},{{end}}{{$elem.Name}}{{end}}{{if gt $returnLength 0}},{{end}} nil
}
{{end}}{{/* End Getter */}}
{{if $f.Setter}}
func {{ToGoNameConv $f.Setter.Name}}({{$f.Name}} {{ConvertToGoType $f.ArgDefinition $m}}) (err error){
	{{ $paramsLength := len $f.Setter.Parameters }}{{ $returnLength := len $f.Setter.ReturnValues }}

	{{GenerateCodeAllocateCDTS $f.Setter.Parameters $f.Setter.ReturnValues}}
	
	// parameters
	{{range $index, $elem := $f.Setter.Parameters}}
	fromGoToCDT({{$elem.Name}}, parametersCDTS, {{$index}})
	{{end}}{{/* End Parameters */}}

	{{GenerateCodeXCall "" $f.Getter.Name $f.Setter.Parameters $f.Setter.ReturnValues}}
	
	{{range $index, $elem := $f.Setter.ReturnValues}}
	{{$elem.Name}}AsInterface := fromCDTToGo(return_valuesCDTS, {{$index}})
	if {{$elem.Name}}AsInterface != nil{
		{{if $elem.IsAny}}
		{{$elem.Name}} = {{$elem.Name}}AsInterface
		
		{{else if not $elem.IsHandle}}
		{{$elem.Name}} = {{if $elem.IsTypeAlias}}{{if $elem.IsArray}}[]{{end}}{{GetTypeOrAlias $elem $m}}{{else}}{{ConvertToGoType $elem $m}}{{end}}({{$elem.Name}}AsInterface.({{ConvertToGoType $elem $m}}))
		{{else}} {{/* Handle */}}		
		{{/* Go object */}}
		
		{{if not $elem.IsArray}}
		if obj, ok := {{$elem.Name}}AsInterface.(Handle); ok{ // None Go object			
			{{$elem.Name}} = {{HandleNoneGoObject $elem $m}}
		} else {
			{{$elem.Name}} = {{if $elem.IsTypeAlias}}{{if $elem.IsArray}}[]{{end}}{{GetTypeOrAlias $elem $m}}{{else}}{{ConvertToGoType $elem $m}}{{end}}({{$elem.Name}}AsInterface.({{ConvertToGoType $elem $m}}))
		}
		{{else}}
		{{if $elem.IsTypeAlias}} {{/* a type is specified */}}
		if len({{$elem.Name}}AsInterface.([]interface{})) > 0{
			{{$elem.Name}} = make([]{{GetTypeOrAlias $elem $m}}, len({{$elem.Name}}AsInterface.([]interface{})))
			if _, ok := {{$elem.Name}}AsInterface.([]interface{})[0].(Handle); ok{
				for i, h := range {{$elem.Name}}AsInterface.([]interface{}){
					{{$elem.Name}}[i] = {{GetTypeOrAlias $elem $m}}{ h: h.(Handle) }
				}
			} else {
				for i, obj := range {{$elem.Name}}AsInterface.([]interface{}){
					{{$elem.Name}}[i] = obj.({{GetTypeOrAlias $elem $m}})
				}
			}
		}
		{{else}}
		{{$elem.Name}} = {{$elem.Name}}AsInterface
		{{end}}
		{{end}}

		{{end}}{{/* end handling types */}}
	}
	
	{{end}}{{/* End return values */}}

	return {{range $index, $elem := $f.Setter.ReturnValues}}{{if $index}},{{end}}{{$elem.Name}}{{end}}{{if gt $returnLength 0}},{{end}} nil
}
{{end}}{{/* End Setter */}}
{{end}}{{/* End Global */}}


// Code to call foreign functions in module {{$m.Name}} via XLLR
{{range $findex, $f := $m.Functions}}
// Call to foreign {{$f.Name}}
{{if $f.Comment}}/*
{{$f.Comment}}
*/{{end}}
{{range $index, $elem := $f.Parameters}}
{{if $elem.Comment}}// {{$elem.Name}} - {{$elem.Comment}}{{end}}{{end}}{{/* End Parameters comments */}}
func {{ToGoNameConv $f.Name}}({{range $index, $elem := $f.Parameters}}{{if $index}},{{end}} {{$elem.Name}} {{ConvertToGoType $elem $m}}{{end}}) ({{range $index, $elem := $f.ReturnValues}}{{if $index}},{{end}}{{$elem.Name}} {{ConvertToGoType $elem $m}}{{end}}{{if $f.ReturnValues}},{{end}} err error){

	{{ $paramsLength := len $f.Parameters }}{{ $returnLength := len $f.ReturnValues }}

	{{GenerateCodeAllocateCDTS $f.Parameters $f.ReturnValues}}
	
	// parameters
	{{range $index, $elem := $f.Parameters}}
	fromGoToCDT({{$elem.Name}}, parametersCDTS, {{$index}})
	{{end}}{{/* End Parameters */}}

	{{GenerateCodeXCall "" $f.Name $f.Parameters $f.ReturnValues}}
	
	{{range $index, $elem := $f.ReturnValues}}
	{{$elem.Name}}AsInterface := fromCDTToGo(return_valuesCDTS, {{$index}})
	if {{$elem.Name}}AsInterface != nil{
		{{if $elem.IsAny}}
		{{$elem.Name}} = {{$elem.Name}}AsInterface
		
		{{else if not $elem.IsHandle}}
		{{$elem.Name}} = {{if $elem.IsTypeAlias}}{{if $elem.IsArray}}[]{{end}}{{GetTypeOrAlias $elem $m}}{{else}}{{ConvertToGoType $elem $m}}{{end}}({{$elem.Name}}AsInterface.({{ConvertToGoType $elem $m}}))
		{{else}} {{/* Handle */}}		
		{{/* Go object */}}
		
		{{if not $elem.IsArray}}
		if obj, ok := {{$elem.Name}}AsInterface.(Handle); ok{ // None Go object			
			{{$elem.Name}} = {{HandleNoneGoObject $elem $m}}
		} else {
			{{$elem.Name}} = {{if $elem.IsTypeAlias}}{{if $elem.IsArray}}[]{{end}}{{GetTypeOrAlias $elem $m}}{{else}}{{ConvertToGoType $elem $m}}{{end}}({{$elem.Name}}AsInterface.({{ConvertToGoType $elem $m}}))
		}
		{{else}}
		{{if $elem.IsTypeAlias}} {{/* a type is specified */}}
		if len({{$elem.Name}}AsInterface.([]interface{})) > 0{
			{{$elem.Name}} = make([]{{GetTypeOrAlias $elem $m}}, len({{$elem.Name}}AsInterface.([]interface{})))
			if _, ok := {{$elem.Name}}AsInterface.([]interface{})[0].(Handle); ok{
				for i, h := range {{$elem.Name}}AsInterface.([]interface{}){
					{{$elem.Name}}[i] = {{GetTypeOrAlias $elem $m}}{ h: h.(Handle) }
				}
			} else {
				for i, obj := range {{$elem.Name}}AsInterface.([]interface{}){
					{{$elem.Name}}[i] = obj.({{GetTypeOrAlias $elem $m}})
				}
			}
		}
		{{else}}
		{{$elem.Name}} = {{$elem.Name}}AsInterface
		{{end}}
		{{end}}

		{{end}}{{/* end handling types */}}
	}
	
	{{end}}{{/* End return values */}}

	return {{range $index, $elem := $f.ReturnValues}}{{if $index}},{{end}}{{$elem.Name}}{{end}}{{if gt $returnLength 0}},{{end}} nil
}
{{end}}{{/* End Function */}}

{{range $cindex, $c := $m.Classes}}
type {{AsPublic $c.Name}} struct{
	h Handle
}
{{range $findex, $f := $c.Constructors}}
func New{{ToGoNameConv $f.Name}}({{range $index, $elem := $f.Parameters}}{{if $index}},{{end}} {{$elem.Name}} {{ConvertToGoType $elem $m}}{{end}}) (instance *{{AsPublic $c.Name}}, err error){
	{{ $paramsLength := len $f.Parameters }}{{ $returnLength := len $f.ReturnValues }}

	{{GenerateCodeAllocateCDTS $f.Parameters $f.ReturnValues}}
	
	// parameters
	{{range $index, $elem := $f.Parameters}}
	fromGoToCDT({{$elem.Name}}, parametersCDTS, {{$index}})
	{{end}}{{/* End Parameters */}}

	{{GenerateCodeXCall $c.Name $f.Name $f.Parameters $f.ReturnValues}}
	
	inst := &{{AsPublic $c.Name}}{}

	{{range $index, $elem := $f.ReturnValues}}
	{{$elem.Name}}AsInterface := fromCDTToGo(return_valuesCDTS, {{$index}})
	if {{$elem.Name}}AsInterface != nil{
		inst.h = {{$elem.Name}}AsInterface.(Handle)
	} else {
		return nil, fmt.Errorf("Object creation returned nil")
	}
		
	{{end}}{{/* End return values */}}

	return inst, nil	
}
{{end}}{{/* End Constructor */}}

{{range $findex, $f := $c.Fields}}
{{if $f.Getter}}
func {{GenerateMethodReceiverCode $f.Getter}} {{GenerateMethodName $f.Getter}}({{GenerateMethodParams $f.Getter $m}}) ({{range $index, $elem := $f.Getter.ReturnValues}}{{if $index}},{{end}}{{$elem.Name}} {{ConvertToGoType $elem $m}}{{end}}{{if $f.Getter.ReturnValues}},{{end}} err error){
	
	{{ $paramsLength := len $f.Getter.Parameters }}{{ $returnLength := len $f.Getter.ReturnValues }}

	{{GenerateCodeAllocateCDTS $f.Getter.Parameters $f.Getter.ReturnValues}}
	
	// get parameters
	{{if $f.Getter.InstanceRequired}}
	fromGoToCDT(this.h, parametersCDTS, 0)
	{{range $index, $elem := $f.Getter.Parameters}}{{if gt $index 0}}
	fromGoToCDT(this.h, parametersCDTS, {{$index}})
	{{end}}{{end}}{{/* End Parameters */}}
	{{else}}
	{{range $index, $elem := $f.Getter.Parameters}}
	fromGoToCDT(this.h, parametersCDTS, {{$index}})
	{{end}} {{/* End Parameters */}}
	{{end}} {{/* End InstanceRequired */}}


	{{GenerateCodeXCall $c.Name $f.Getter.Name $f.Getter.Parameters $f.Getter.ReturnValues}}
	
	{{range $index, $elem := $f.Getter.ReturnValues}}
	{{$elem.Name}}AsInterface := fromCDTToGo(return_valuesCDTS, {{$index}})
	if {{$elem.Name}}AsInterface != nil{
		{{if $elem.IsAny}}
		// any
		{{$elem.Name}} = {{$elem.Name}}AsInterface
		
		{{else if not $elem.IsHandle}}
		// not handle
		{{$elem.Name}} = {{if $elem.IsTypeAlias}}{{if $elem.IsArray}}[]{{end}}{{GetTypeOrAlias $elem $m}}{{else}}{{ConvertToGoType $elem $m}}{{end}}({{$elem.Name}}AsInterface.({{ConvertToGoType $elem $m}}))
		
		{{else}} {{/* Handle */}}
		// handle
		{{if not $elem.IsArray}}
		if obj, ok := {{$elem.Name}}AsInterface.(Handle); ok{ // None Go object			
			{{$elem.Name}} = {{HandleNoneGoObject $elem $m}}
		} else {
			{{$elem.Name}} = {{if $elem.IsTypeAlias}}{{if $elem.IsArray}}[]{{end}}{{GetTypeOrAlias $elem $m}}{{else}}{{ConvertToGoType $elem $m}}{{end}}({{$elem.Name}}AsInterface.({{ConvertToGoType $elem $m}}))
		}
		{{else}}
		{{if $elem.IsTypeAlias}} {{/* a type is specified */}}
		if len({{$elem.Name}}AsInterface.([]interface{})) > 0{
			{{$elem.Name}} = make([]{{GetTypeOrAlias $elem $m}}, len({{$elem.Name}}AsInterface.([]interface{})))
			if _, ok := {{$elem.Name}}AsInterface.([]interface{})[0].(Handle); ok{
				for i, h := range {{$elem.Name}}AsInterface.([]interface{}){
					{{$elem.Name}}[i] = {{GetTypeOrAlias $elem $m}}{ h: h.(Handle) }
				}
			} else {
				for i, obj := range {{$elem.Name}}AsInterface.([]interface{}){
					{{$elem.Name}}[i] = obj.({{GetTypeOrAlias $elem $m}})
				}
			}
		}
		{{else}}
		{{$elem.Name}} = {{$elem.Name}}AsInterface
		{{end}}
		{{end}}

		{{end}}{{/* end handling types */}}
	}

	{{end}}{{/* End return values */}}

	return {{range $index, $elem := $f.Getter.ReturnValues}}{{if $index}},{{end}}{{$elem.Name}}{{end}}{{if gt $returnLength 0}},{{end}} nil	
}
{{end}}{{/* End Getter */}}
{{if $f.Setter}}
func {{GenerateMethodReceiverCode $f.Setter}} {{GenerateMethodName $f.Setter}}({{GenerateMethodParams $f.Setter $m}}) ({{range $index, $elem := $f.Setter.ReturnValues}}{{if $index}},{{end}}{{$elem.Name}} {{ConvertToGoType $elem $m}}{{end}}{{if $f.Setter.ReturnValues}},{{end}} err error){
	
	{{ $paramsLength := len $f.Setter.Parameters }}{{ $returnLength := len $f.Setter.ReturnValues }}

	{{GenerateCodeAllocateCDTS $f.Setter.Parameters $f.Setter.ReturnValues}}
	
	// parameters
	fromGoToCDT(this.h, parametersCDTS, 0) // object
	{{range $index, $elem := $f.Setter.Parameters}}{{if gt $index 0}}
	fromGoToCDT({{$elem.Name}}, parametersCDTS, {{$index}})
	{{end}}{{end}}{{/* End Parameters */}}

	{{GenerateCodeXCall $c.Name $f.Setter.Name $f.Setter.Parameters $f.Setter.ReturnValues}}
	
	{{range $index, $elem := $f.Setter.ReturnValues}}
	{{$elem.Name}}AsInterface := fromCDTToGo(return_valuesCDTS, {{$index}})
	if {{$elem.Name}}AsInterface != nil{
		{{if $elem.IsAny}}
		{{$elem.Name}} = {{$elem.Name}}AsInterface
		
		{{else if not $elem.IsHandle}}
		{{$elem.Name}} = {{if $elem.IsTypeAlias}}{{if $elem.IsArray}}[]{{end}}{{GetTypeOrAlias $elem $m}}{{else}}{{ConvertToGoType $elem $m}}{{end}}({{$elem.Name}}AsInterface.({{ConvertToGoType $elem $m}}))
		{{else}} {{/* Handle */}}		
		{{/* Go object */}}
		
		{{if not $elem.IsArray}}
		if obj, ok := {{$elem.Name}}AsInterface.(Handle); ok{ // None Go object			
			{{$elem.Name}} = {{HandleNoneGoObject $elem $m}}
		} else {
			{{$elem.Name}} = {{if $elem.IsTypeAlias}}{{if $elem.IsArray}}[]{{end}}{{GetTypeOrAlias $elem $m}}{{else}}{{ConvertToGoType $elem $m}}{{end}}({{$elem.Name}}AsInterface.({{ConvertToGoType $elem $m}}))
		}
		{{else}}
		{{if $elem.IsTypeAlias}} {{/* a type is specified */}}
		if len({{$elem.Name}}AsInterface.([]interface{})) > 0{
			{{$elem.Name}} = make([]{{GetTypeOrAlias $elem $m}}, len({{$elem.Name}}AsInterface.([]interface{})))
			if _, ok := {{$elem.Name}}AsInterface.([]interface{})[0].(Handle); ok{
				for i, h := range {{$elem.Name}}AsInterface.([]interface{}){
					{{$elem.Name}}[i] = {{GetTypeOrAlias $elem $m}}{ h: h.(Handle) }
				}
			} else {
				for i, obj := range {{$elem.Name}}AsInterface.([]interface{}){
					{{$elem.Name}}[i] = obj.({{GetTypeOrAlias $elem $m}})
				}
			}
		}
		{{else}}
		{{$elem.Name}} = {{$elem.Name}}AsInterface
		{{end}}
		{{end}}
		
		{{end}}{{/* end handling types */}}
	}
	
	{{end}}{{/* End return values */}}

	return {{range $index, $elem := $f.Setter.ReturnValues}}{{if $index}},{{end}}{{$elem.Name}}{{end}}{{if gt $returnLength 0}},{{end}} nil	
}
{{end}}{{/* End Setter */}}
{{end}}{{/* End Fields */}}
{{range $findex, $f := $c.Methods}}
func {{GenerateMethodReceiverCode $f}} {{GenerateMethodName $f}}({{GenerateMethodParams $f $m}}) ({{range $index, $elem := $f.ReturnValues}}{{if $index}},{{end}}{{$elem.Name}} {{ConvertToGoType $elem $m}}{{end}}{{if $f.ReturnValues}},{{end}} err error){
	
	{{ $paramsLength := len $f.Parameters }}{{ $returnLength := len $f.ReturnValues }}

	{{GenerateCodeAllocateCDTS $f.Parameters $f.ReturnValues}}
	
	// parameters
	{{if $f.InstanceRequired}}
	fromGoToCDT(this.h, parametersCDTS, 0) // object
	{{range $index, $elem := $f.Parameters}}{{if gt $index 0}}
	fromGoToCDT({{$elem.Name}}, parametersCDTS, {{$index}})
	{{end}}{{end}}
	{{else}}
	{{range $index, $elem := $f.Parameters}}
	fromGoToCDT({{$elem.Name}}, parametersCDTS, {{$index}})
	{{end}}
	{{end}}

	{{GenerateCodeXCall $c.Name $f.Name $f.Parameters $f.ReturnValues}}
	
	{{range $index, $elem := $f.ReturnValues}}
	{{$elem.Name}}AsInterface := fromCDTToGo(return_valuesCDTS, {{$index}})
	if {{$elem.Name}}AsInterface != nil{
		{{if $elem.IsAny}}
		{{$elem.Name}} = {{$elem.Name}}AsInterface
		
		{{else if not $elem.IsHandle}}
		{{$elem.Name}} = {{if $elem.IsTypeAlias}}{{if $elem.IsArray}}[]{{end}}{{GetTypeOrAlias $elem $m}}{{else}}{{ConvertToGoType $elem $m}}{{end}}({{$elem.Name}}AsInterface.({{ConvertToGoType $elem $m}}))
		{{else}} {{/* Handle */}}		
		{{/* Go object */}}
		
		{{if not $elem.IsArray}}
		if obj, ok := {{$elem.Name}}AsInterface.(Handle); ok{ // None Go object			
			{{$elem.Name}} = {{HandleNoneGoObject $elem $m}}
		} else {
			{{$elem.Name}} = {{if $elem.IsTypeAlias}}{{if $elem.IsArray}}[]{{end}}{{GetTypeOrAlias $elem $m}}{{else}}{{ConvertToGoType $elem $m}}{{end}}({{$elem.Name}}AsInterface.({{ConvertToGoType $elem $m}}))
		}
		{{else}}
		{{if $elem.IsTypeAlias}} {{/* a type is specified */}}
		if len({{$elem.Name}}AsInterface.([]interface{})) > 0{
			{{$elem.Name}} = make([]{{GetTypeOrAlias $elem $m}}, len({{$elem.Name}}AsInterface.([]interface{})))
			if _, ok := {{$elem.Name}}AsInterface.([]interface{})[0].(Handle); ok{
				for i, h := range {{$elem.Name}}AsInterface.([]interface{}){
					{{$elem.Name}}[i] = {{GetTypeOrAlias $elem $m}}{ h: h.(Handle) }
				}
			} else {
				for i, obj := range {{$elem.Name}}AsInterface.([]interface{}){
					{{$elem.Name}}[i] = obj.({{GetTypeOrAlias $elem $m}})
				}
			}
		}
		{{else}}
		{{$elem.Name}} = {{$elem.Name}}AsInterface
		{{end}}
		{{end}}

		{{end}}{{/* end handling types */}}
	}
	
	{{end}}{{/* End return values */}}

	return {{range $index, $elem := $f.ReturnValues}}{{if $index}},{{end}}{{$elem.Name}}{{end}}{{if gt $returnLength 0}},{{end}} nil
}
{{end}}{{/* End Methods */}}
{{if $c.Releaser}}{{ $f := $c.Releaser}}
func (this *{{AsPublic $c.Name}}) {{ToGoNameConv $f.Name}}({{range $index, $elem := $f.Parameters}}{{if $index}},{{end}} {{$elem.Name}} {{ConvertToGoType $elem $m}}{{end}}) ({{range $index, $elem := $f.ReturnValues}}{{if $index}},{{end}}{{$elem.Name}} {{ConvertToGoType $elem $m}}{{end}}{{if $f.ReturnValues}},{{end}} err error){
	
	{{ $paramsLength := len $f.Parameters }}{{ $returnLength := len $f.ReturnValues }}

	{{GenerateCodeAllocateCDTS $f.Parameters $f.ReturnValues}}
	
	// parameters
	fromGoToCDT(this.h, parametersCDTS, 0) // object
	{{range $index, $elem := $f.Parameters}}{{if gt $index 0}}
	fromGoToCDT({{$elem.Name}}, parametersCDTS, {{$index}})
	{{end}}{{end}}{{/* End Parameters */}}

	{{GenerateCodeXCall $c.Name $f.Name $f.Parameters $f.ReturnValues}}
	
	{{range $index, $elem := $f.ReturnValues}}
	{{$elem.Name}}AsInterface := fromCDTToGo(return_valuesCDTS, {{$index}})
	if {{$elem.Name}}AsInterface != nil{
		{{if $elem.IsAny}}
		{{$elem.Name}} = {{$elem.Name}}AsInterface
		
		{{else if not $elem.IsHandle}}
		{{$elem.Name}} = {{if $elem.IsTypeAlias}}{{if $elem.IsArray}}[]{{end}}{{GetTypeOrAlias $elem $m}}{{else}}{{ConvertToGoType $elem $m}}{{end}}({{$elem.Name}}AsInterface.({{ConvertToGoType $elem $m}}))
		{{else}} {{/* Handle */}}		
		{{/* Go object */}}
		
		{{if not $elem.IsArray}}
		if obj, ok := {{$elem.Name}}AsInterface.(Handle); ok{ // None Go object			
			{{$elem.Name}} = {{HandleNoneGoObject $elem $m}}
		} else {
			{{$elem.Name}} = {{if $elem.IsTypeAlias}}{{if $elem.IsArray}}[]{{end}}{{GetTypeOrAlias $elem $m}}{{else}}{{ConvertToGoType $elem $m}}{{end}}({{$elem.Name}}AsInterface.({{ConvertToGoType $elem $m}}))
		}
		{{else}}
		{{if $elem.IsTypeAlias}} {{/* a type is specified */}}
		if len({{$elem.Name}}AsInterface.([]interface{})) > 0{
			{{$elem.Name}} = make([]{{GetTypeOrAlias $elem $m}}, len({{$elem.Name}}AsInterface.([]interface{})))
			if _, ok := {{$elem.Name}}AsInterface.([]interface{})[0].(Handle); ok{
				for i, h := range {{$elem.Name}}AsInterface.([]interface{}){
					{{$elem.Name}}[i] = {{GetTypeOrAlias $elem $m}}{ h: h.(Handle) }
				}
			} else {
				for i, obj := range {{$elem.Name}}AsInterface.([]interface{}){
					{{$elem.Name}}[i] = obj.({{GetTypeOrAlias $elem $m}})
				}
			}
		}
		{{else}}
		{{$elem.Name}} = {{$elem.Name}}AsInterface
		{{end}}
		{{end}}

		{{end}}{{/* end handling types */}}
	}
	
	{{end}}{{/* End return values */}}

	return {{range $index, $elem := $f.ReturnValues}}{{if $index}},{{end}}{{$elem.Name}}{{end}}{{if gt $returnLength 0}},{{end}} nil
}
{{end}}{{/* End Releaser */}}
{{end}}{{/* End Classes */}}
{{end}}{{/* End modules */}}

`
