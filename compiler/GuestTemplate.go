package main

const GuestHeaderTemplate = `
// Code generated by MetaFFI. Modify only in marked places.
// Guest code for {{.IDLFilenameWithExtension}}

package main
`

const GuestImportsTemplate = `
import "fmt"
import "unsafe"
import "github.com/MetaFFI/plugin-sdk/compiler/go/IDL"
import . "github.com/MetaFFI/lang-plugin-go/go-runtime"
{{range $mindex, $i := .Imports}}
import . "{{$i}}"{{end}}

{{range $mindex, $m := .Modules}}
{{range $eindex, $e := $m.ExternalResources}}
import "{{$e}}"{{end}}{{end}}

`

const GuestCImportCGoFileTemplate = `
package main

/*
#cgo !windows LDFLAGS: -L. -ldl
#cgo CFLAGS: -I"{{GetEnvVar "METAFFI_HOME" true}}"


#include <include/cdt_structs.h>
#include <include/cdt_capi_loader.c>

{{/* TODO: Do this without item CGo https://stackoverflow.com/questions/53238602/accessing-c-array-in-golang*/}}
metaffi_size get_int_item(metaffi_size* array, int index)
{
	return array[index];
}

{{/* TODO: Do this without item CGo https://stackoverflow.com/questions/53238602/accessing-c-array-in-golang*/}}
void* convert_union_to_ptr(void* p)
{
	return p;
}

{{/* TODO: Do this without item CGo https://stackoverflow.com/questions/53238602/accessing-c-array-in-golang*/}}
void set_cdt_type(struct cdt* p, metaffi_type t)
{
	p->type = t;
}

{{/* TODO: Do this without item CGo https://stackoverflow.com/questions/53238602/accessing-c-array-in-golang*/}}
metaffi_type get_cdt_type(struct cdt* p)
{
	return p->type;
}

{{/* TODO: Do this without item CGo https://stackoverflow.com/questions/53238602/accessing-c-array-in-golang*/}}
struct cdt* get_cdt_element(struct cdts* pdata, int cdts_index)
{
	return pdata[cdts_index].pcdt;
}

void set_go_runtime_flag()
{
	xllr_set_runtime_flag("go_runtime", 10);
}

metaffi_handle get_null_handle()
{
	return METAFFI_NULL_HANDLE;
}

#ifdef _WIN32
metaffi_size len_to_metaffi_size(long long i)
#else
metaffi_size len_to_metaffi_size(long long i)
#endif
{
	return (metaffi_size)i;
}

*/
import "C"
`

const GuestCImportTemplate = `
/*
#cgo !windows LDFLAGS: -L. -ldl
#cgo CFLAGS: -I"{{GetEnvVar "METAFFI_HOME" true}}"

#include <include/cdt_structs.h>
#include <include/cdt_capi_loader.h>

metaffi_size get_int_item(metaffi_size* array, int index);
void* convert_union_to_ptr(void* p);
void set_cdt_type(struct cdt* p, metaffi_type t);
metaffi_type get_cdt_type(struct cdt* p);
void set_go_runtime_flag();
struct cdt* get_cdt_element(struct cdts* pdata, int cdts_index);
metaffi_handle get_null_handle();

#ifdef _WIN32
	metaffi_size len_to_metaffi_size(long long i);
#else
	metaffi_size len_to_metaffi_size(long long i);
#endif
*/
import "C"
`

const GuestMainFunction = `
func main(){} // main function must be declared to create dynamic library
func init(){
	err := C.load_cdt_capi()
	if err != nil{
		panic("Failed to load MetaFFI XLLR functions: "+C.GoString(err))
	}
	C.set_go_runtime_flag()
}
`

const GuestHelperFunctionsTemplate = `

func errToOutError(out_err **C.char, out_err_len *C.uint64_t, customText string, err error){
	txt := customText
	if err != nil { txt += err.Error() }
	*out_err = C.CString(txt)
	*out_err_len = C.uint64_t(len(txt))
}

func panicHandler(out_err **C.char, out_err_len *C.uint64_t){
	if rec := recover(); rec != nil{
		msg := "Panic in Go function. Panic Data: "
		switch recType := rec.(type){
			case error: msg += (rec.(error)).Error()
			case string: msg += rec.(string)
			default: msg += fmt.Sprintf("Panic with type: %v - %v", recType, rec)
		}

		*out_err = C.CString(msg)
		*out_err_len = C.uint64_t(len(msg))
	}
}

`

const (
	GuestFunctionXLLRTemplate = `
{{$def := .}}

{{range $mindex, $m := .Modules}}

{{range $findex, $f := $m.Globals}}
{{if $f.Getter}}
// getter for {{$f.Name}}
//export EntryPoint_{{$f.Getter.Name}} 
func EntryPoint_{{GenerateCodeEntryPointSignature "" $f.Getter.Name $f.Getter.Parameters $f.Getter.ReturnValues}}{

	// catch panics and return them as errors
	defer panicHandler(out_err, out_err_len)

	{{ $paramsLength := len $f.Getter.Parameters }}{{ $returnLength := len $f.Getter.ReturnValues }}

	retvals_CDTS := C.get_cdt_element(xcall_params, 1)
	t0 := IDL.MetaFFITypeInfo{ {{$t := index $f.Getter.ReturnValues 0}}
	    StringType: "{{$t.Type}}",
	    Alias:"{{$t.TypeAlias}}",
	    Dimensions: {{$t.Dimensions}},
		Type: {{GetMetaFFINumericType $t.Type}},
	}
	FromGoToCDT({{AssertAndConvert $f.Name $t $m}}, unsafe.Pointer(retvals_CDTS), t0, 0)

}
{{end}} {{/* end $f.Get */}}

{{if $f.Setter}}
// setter for {{$f.Name}}
//export EntryPoint_{{$f.Setter.Name}}
func EntryPoint_{{GenerateCodeEntryPointSignature "" $f.Setter.Name $f.Setter.Parameters $f.Setter.ReturnValues}}{

	// catch panics and return them as errors
	defer panicHandler(out_err, out_err_len)

	{{ $paramsLength := len $f.Setter.Parameters }}{{ $returnLength := len $f.Setter.ReturnValues }}

	parameters_CDTS := C.get_cdt_element(xcall_params, 0)
	{{$elem := index $f.Setter.Parameters 0}}
	{{ConvertEmptyInterfaceFromCDTSToCorrectType $elem $m true}}
	
{{end}} {{/* end $f.Set */}}

{{end}} {{/* end range globals */}}


// functions
{{range $findex, $f := $m.Functions}}
// Call to foreign {{$f.Name}}
//export EntryPoint_{{$f.Name}}
func EntryPoint_{{GenerateCodeEntryPointSignature "" $f.Name $f.Parameters $f.ReturnValues}}{
	// catch panics and return them as errors
	defer panicHandler(out_err, out_err_len)

	{{ $paramsLength := len $f.Parameters }}{{ $returnLength := len $f.ReturnValues }}

	{{ if gt $paramsLength 0 }}
	parameters_CDTS := C.get_cdt_element(xcall_params, 0)
	{{end}}
	{{ if gt $returnLength 0 }}
	retvals_CDTS := C.get_cdt_element(xcall_params, 1)
	{{end}}

	// parameters from C to Go
	{{range $index, $elem := $f.Parameters}}
	{{$elem.Name}}AsInterface := FromCDTToGo(unsafe.Pointer(parameters_CDTS), {{$index}})
	{{ConvertEmptyInterfaceFromCDTSToCorrectType $elem $m false}}
	{{end}} {{/* end range params */}}
	
	// call original function
	{{range $index, $elem := $f.ReturnValues}}{{if $index}},{{end}}{{$elem.Name}}{{end}}{{if $f.ReturnValues}} := {{end}}{{$f.Name}}({{CallParameters $f 0}})
	
	// return values
	{{range $index, $elem := $f.ReturnValues}}
	if err, isError := interface{}({{$elem.Name}}).(error); isError{ // in case of error
		errToOutError(out_err, out_err_len, "Error returned", err)
		return
	} else { // Convert return values from Go to C
		t{{$index}} := IDL.MetaFFITypeInfo{   {{$t := index $f.ReturnValues $index}}
			StringType: "{{$t.Type}}",
			Alias:"{{$t.TypeAlias}}",
			Dimensions: {{$t.Dimensions}},
			Type: {{GetMetaFFINumericType $t.Type}},
		}
		FromGoToCDT({{AssertAndConvert $elem.Name $t $m}}, unsafe.Pointer(retvals_CDTS), t{{$index}}, {{$index}})
	}	
	{{end}} {{/* end range return vals */}}
}
{{end}} {{/* end range functions */}}


{{range $cindex, $c := $m.Classes}}
// class {{$c.Name}}
{{$className := $c.Name}}

// return empty struct
//export EntryPoint_{{$c.Name}}_EmptyStruct_MetaFFI
func EntryPoint_{{GenerateCodeEntryPointEmptyStructSignature $c.Name}}{
	instance := &{{$c.Name}}{}
	FromGoToCDT(instance, unsafe.Pointer(C.get_cdt_element(xcall_params, 1)), IDL.MetaFFITypeInfo{ StringType: IDL.HANDLE, Dimensions: 0, Type: {{GetMetaFFINumericType "handle"}} }, 0)
}

// constructors
{{range $i, $f := $c.Constructors}}
// Call to foreign {{$f.Name}}
//export EntryPoint_{{$c.Name}}_{{$f.Name}}
func EntryPoint_{{GenerateCodeEntryPointSignature $c.Name $f.Name $f.Parameters $f.ReturnValues}}{
	// catch panics and return them as errors
	defer panicHandler(out_err, out_err_len)

	{{ $paramsLength := len $f.Parameters }}{{ $returnLength := len $f.ReturnValues }}
	
	{{ if gt $paramsLength 0 }}
	parameters_CDTS := C.get_cdt_element(xcall_params, 0)
	{{end}}
	{{ if gt $returnLength 0 }}
	retvals_CDTS := C.get_cdt_element(xcall_params, 1)
	{{end}}

	// parameters from C to Go
	{{range $index, $elem := $f.Parameters}}	
	{{$elem.Name}}AsInterface := FromCDTToGo(unsafe.Pointer(parameters_CDTS), {{$index}})
	{{ConvertEmptyInterfaceFromCDTSToCorrectType $elem $m false}}
	{{end}} {{/* end range params */}}
	
	// call original function
	{{range $index, $elem := $f.ReturnValues}}{{if $index}},{{end}}{{$elem.Name}}{{end}}{{if $f.ReturnValues}} := {{end}}{{$f.Name}}({{CallParameters $f.FunctionDefinition 0}})
	
	// return values
	{{range $index, $elem := $f.ReturnValues}}
	if err, isError := interface{}({{$elem.Name}}).(error); isError{ // in case of error
		errToOutError(out_err, out_err_len, "Error returned", err)
		return
	} else { // Convert return values from Go to C
		t{{$index}} := IDL.MetaFFITypeInfo{ {{$t := index $f.ReturnValues $index}}
			StringType: "{{$t.Type}}",
			Alias:"{{$t.TypeAlias}}",
			Dimensions: {{$t.Dimensions}},
			Type: {{GetMetaFFINumericType $t.Type}},
		}
		FromGoToCDT({{AssertAndConvert $elem.Name $t $m}}, unsafe.Pointer(retvals_CDTS), t{{$index}}, {{$index}})
	}	
	{{end}} {{/* end range return vals */}}
}
{{end}} {{/* end range constructors */}}

// methods
{{range $i, $f := $c.Methods}}
// Call to foreign {{$f.Name}}
//export EntryPoint_{{$c.Name}}_{{$f.Name}}
func EntryPoint_{{GenerateCodeEntryPointSignature $c.Name $f.Name $f.Parameters $f.ReturnValues}}{
	// catch panics and return them as errors
	defer panicHandler(out_err, out_err_len)

	{{ $paramsLength := len $f.Parameters }}{{ $returnLength := len $f.ReturnValues }}
	
	{{ if gt $paramsLength 0 }}
	parameters_CDTS := C.get_cdt_element(xcall_params, 0)
	{{end}}
	{{ if gt $returnLength 0 }}
	retvals_CDTS := C.get_cdt_element(xcall_params, 1)
	{{end}}

	// parameters from C to Go
	{{range $index, $elem := $f.Parameters}}	
	{{$elem.Name}}AsInterface := FromCDTToGo(unsafe.Pointer(parameters_CDTS), {{$index}})
	{{ConvertEmptyInterfaceFromCDTSToCorrectType $elem $m false}}
	{{end}} {{/* end range params */}}
	
	// call original function
	{{ $receiver_pointer := index $f.Tags "receiver_pointer"}}
	{{range $index, $elem := $f.ReturnValues}}{{if $index}},{{end}}{{$elem.Name}}{{end}}{{if $f.ReturnValues}} := {{end}}({{(index $f.Parameters 0).Name }}.({{if eq $receiver_pointer "true"}}*{{end}}{{$className}})).{{$f.Name}}({{ CallParameters $f.FunctionDefinition 1}})

	// return values
	{{range $index, $elem := $f.ReturnValues}}
	if err, isError := interface{}({{$elem.Name}}).(error); isError{ // in case of error
		errToOutError(out_err, out_err_len, "Error returned", err)
		return
	} else { // Convert return values from Go to C
		t{{$index}} := IDL.MetaFFITypeInfo{   {{$t := index $f.ReturnValues $index}}
			StringType: "{{$t.Type}}",
			Alias:"{{$t.TypeAlias}}",
			Dimensions: {{$t.Dimensions}},
			Type: {{GetMetaFFINumericType $t.Type}},
		}
		FromGoToCDT({{AssertAndConvert $elem.Name $t $m}}, unsafe.Pointer(retvals_CDTS), t{{$index}}, {{$index}})
	}	
	{{end}} {{/* end range return values */}}
}
{{end}} {{/* end range methods */}}

// Fields
{{range $i, $f := $c.Fields}}
{{if $f.Getter}}
// getter for {{$f.Name}}
//export EntryPoint_{{$c.Name}}_{{$f.Getter.Name}}
func EntryPoint_{{GenerateCodeEntryPointSignature $c.Name $f.Getter.Name $f.Getter.Parameters $f.Getter.ReturnValues}}{

	// catch panics and return them as errors
	defer panicHandler(out_err, out_err_len)
	{{ $paramsLength := len $f.Getter.Parameters }}{{ $returnLength := len $f.Getter.ReturnValues }}

	{{ if gt $paramsLength 0 }}
	parameters_CDTS := C.get_cdt_element(xcall_params, 0)
	{{end}}
	{{ if gt $returnLength 0 }}
	retvals_CDTS := C.get_cdt_element(xcall_params, 1)
	{{end}}

	// get object
	{{ $elem := index $f.Getter.Parameters 0 }}
	objAsInterface := FromCDTToGo(unsafe.Pointer(parameters_CDTS), 0)
	obj := {{if not $elem.IsAny}}{{if $elem.IsTypeAlias}}{{$elem.GetTypeOrAlias}}{{else}}{{ConvertToGoType $elem $m}}{{end}}({{end}}objAsInterface{{if not $elem.IsAny}}.({{ConvertToGoType $elem $m}})){{end}}

	{{ $receiver_pointer := index $f.Getter.Tags "receiver_pointer"}}
	{{$f.Name}} := obj.({{if eq $receiver_pointer "true"}}*{{end}}{{$className}}).{{$f.Name}}

	{{ if gt $returnLength 0 }}
	t0 := IDL.MetaFFITypeInfo{ {{$t := index $f.Getter.ReturnValues 0}}
		StringType: "{{$t.Type}}",
		Alias:"{{$t.TypeAlias}}",
		Dimensions: {{$t.Dimensions}},
		Type: {{GetMetaFFINumericType $t.Type}},
	}
	FromGoToCDT({{AssertAndConvert $f.Name $t $m}}, unsafe.Pointer(retvals_CDTS), t0, 0)
	{{end}}
}
{{end}} {{/* end $f.Getter */}}

{{if $f.Setter}}
// setter for {{$f.Name}}
//export EntryPoint_{{$c.Name}}_{{$f.Setter.Name}}
func EntryPoint_{{GenerateCodeEntryPointSignature $c.Name $f.Setter.Name $f.Setter.Parameters $f.Setter.ReturnValues}}{

	// catch panics and return them as errors
	defer panicHandler(out_err, out_err_len)

	{{ $paramsLength := len $f.Setter.Parameters }}{{ $returnLength := len $f.Setter.ReturnValues }}

	{{ if gt $paramsLength 0 }}
	parameters_CDTS := C.get_cdt_element(xcall_params, 0)
	{{end}}
	{{ if gt $returnLength 0 }}
	retvals_CDTS := C.get_cdt_element(xcall_params, 1)
	{{end}}

	// get object
	{{ $elem := index $f.Setter.Parameters 0 }}
	thisAsInterface := FromCDTToGo(unsafe.Pointer(parameters_CDTS), 0)
	this := {{if not $elem.IsAny}}{{if $elem.IsTypeAlias}}{{$elem.GetTypeOrAlias}}{{else}}{{ConvertToGoType $elem $m}}{{end}}({{end}}thisAsInterface{{if not $elem.IsAny}}.({{ConvertToGoType $elem $m}})){{end}}

	// get val
	{{ $elem = index $f.Setter.Parameters 1 }}
	{{$elem.Name}}AsInterface := FromCDTToGo(unsafe.Pointer(parameters_CDTS), 1)
	{{ConvertEmptyInterfaceFromCDTSToCorrectType $elem $m false}}
	
	// set new data
	{{ $receiver_pointer := index $f.Setter.Tags "receiver_pointer"}}
	this.({{if eq $receiver_pointer "true"}}*{{end}}{{$className}}).{{$f.Name}} = {{$elem.Name}}
	
}
{{end}}{{/* end $f.Setter */}}

{{end}} {{/* end range fields */}}

// end class {{$c.Name}}
{{end}} {{/* end range classes */}}

{{end}} {{/* end range modules */}}

`
)
