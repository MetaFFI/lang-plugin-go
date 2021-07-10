package main

const GuestHeaderTemplate = `
// Code generated by OpenFFI. Modify only in marked places.
// Guest code for {{.IDLFilenameWithExtension}}

package main
`

const GuestImportsTemplate = `
import "fmt"
import "unsafe"
{{range $mindex, $i := .Imports}}
import . "{{$i}}"{{end}}
`

const GuestCImportCGoFileTemplate = `
package main

/*
#cgo !windows LDFLAGS: -L. -ldl
#cgo CFLAGS: -I{{GetEnvVar "OPENFFI_HOME"}}


#include <include/cdt_structs.h>
#include <include/cdt_capi_loader.c>

openffi_size get_int_item(openffi_size* array, int index)
{
	return array[index];
}

void* convert_union_to_ptr(void* p)
{
	return p;
}

void set_cdt_type(struct cdt* p, openffi_type t)
{
	p->type = t;
}

*/
import "C"
`

const GuestCImportTemplate = `
/*
#cgo !windows LDFLAGS: -L. -ldl
#cgo CFLAGS: -I{{GetEnvVar "OPENFFI_HOME"}}

#include <include/cdt_structs.h>
#include <include/cdt_capi_loader.h>

openffi_size get_int_item(openffi_size* array, int index);
void* convert_union_to_ptr(void* p);
void set_cdt_type(struct cdt* p, openffi_type t);
*/
import "C"
`

const GuestMainFunction = `
func main(){} // main function must be declared to create dynamic library
func init(){
	err := C.load_cdt_capi()
	if err != nil{
		panic("Failed to load OpenFFI XLLR functions: "+C.GoString(err))
	}
}
`

const GuestHelperFunctions = `
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

const GuestFunctionXLLRTemplate = `
// add functions
{{range $mindex, $m := .Modules}}

{{range $findex, $f := $m.Functions}}

// Call to foreign {{$f.PathToForeignFunction.function}}
//export EntryPoint_{{$f.PathToForeignFunction.function}}
func EntryPoint_{{$f.PathToForeignFunction.function}}(parameters *C.struct_cdt, parameters_length C.uint64_t, return_values *C.struct_cdt, return_values_length C.uint64_t, out_err **C.char, out_err_len *C.uint64_t){

	// catch panics and return them as errors
	defer panicHandler(out_err, out_err_len)

	{{ $paramsLength := len $f.Parameters }}{{ $returnLength := len $f.ReturnValues }}
	
	// parameters from C to Go
	{{range $index, $elem := $f.Parameters}}

	{{if $elem.IsString}}
	
	{{if gt $elem.Dimensions 0}}
	// string[] // TODO: handle multi-dimensional arrays
	
	in_{{$elem.Name}}_cdt := C.get_cdt(parameters, {{$index}})
	pcdt_{{$elem.Type}}_{{$elem.Name}} := ((*C.struct_cdt_openffi_{{$elem.Type}}_array)(C.convert_union_to_ptr(unsafe.Pointer(&in_{{$elem.Name}}_cdt.cdt_val))))
	
	var in_{{$elem.Name}} *C.openffi_{{$elem.Type}} = pcdt_{{$elem.Type}}_{{$elem.Name}}.vals
	var in_{{$elem.Name}}_sizes *C.openffi_size = pcdt_{{$elem.Type}}_{{$elem.Name}}.vals_sizes
	var in_{{$elem.Name}}_dimensions_lengths *C.openffi_size = pcdt_{{$elem.Type}}_{{$elem.Name}}.dimensions_lengths
	// var in_{{$elem.Name}}_dimensions C.openffi_size = pcdt_{{$elem.Type}}_{{$elem.Name}}.dimensions - TODO: not used until multi-dimensions support!
		
	{{$elem.Name}} := make([]string, 0, int(C.get_int_item(in_{{$elem.Name}}_dimensions_lengths, 0)))
	for i:=C.int(0) ; i<C.int(C.get_int_item(in_{{$elem.Name}}_dimensions_lengths, 0)) ; i++{
		var str_size C.openffi_size
		str := C.get_openffi_{{$elem.Type}}_element(in_{{$elem.Name}}, i, in_{{$elem.Name}}_sizes, &str_size)
		{{$elem.Name}} = append({{$elem.Name}}, C.GoStringN(str, C.int(str_size)))
	}
	{{else}}
	// string
	in_{{$elem.Name}}_cdt := C.get_cdt(parameters, {{$index}})
	pcdt_{{$elem.Type}}_{{$elem.Name}} := ((*C.struct_cdt_openffi_{{$elem.Type}})(C.convert_union_to_ptr(unsafe.Pointer(&in_{{$elem.Name}}_cdt.cdt_val))))

	var in_{{$elem.Name}}_len C.openffi_size = pcdt_{{$elem.Type}}_{{$elem.Name}}.length
	var in_{{$elem.Name}} C.openffi_{{$elem.Type}} = pcdt_{{$elem.Type}}_{{$elem.Name}}.val

	{{$elem.Name}} := C.GoStringN(in_{{$elem.Name}}, C.int(in_{{$elem.Name}}_len))
	{{end}}{{else}}{{if $elem.IsArray}}

	// non-string array
	
	in_{{$elem.Name}}_cdt := C.get_cdt(parameters, {{$index}})
	pcdt_{{$elem.Type}}_{{$elem.Name}} := ((*C.struct_cdt_openffi_{{$elem.Type}}_array)(C.convert_union_to_ptr(unsafe.Pointer(&in_{{$elem.Name}}_cdt.cdt_val))))

	var in_{{$elem.Name}} *C.openffi_{{$elem.Type}} = pcdt_{{$elem.Type}}_{{$elem.Name}}.vals
	var in_{{$elem.Name}}_dimensions_lengths *C.openffi_size = pcdt_{{$elem.Type}}_{{$elem.Name}}.dimensions_lengths
	// var in_{{$elem.Name}}_dimensions C.openffi_size = pcdt_{{$elem.Type}}_{{$elem.Name}}.dimensions - TODO: not used until multi-dimensions support!
	
	{{$elem.Name}} := make([]{{$elem.Type}}, 0)
	for i:=C.int(0) ; i<C.int(C.int(C.get_int_item(in_{{$elem.Name}}_dimensions_lengths, 0))) ; i++{
		val := C.get_openffi_{{$elem.Type}}_element(in_{{$elem.Name}}, C.int(i))
		{{$elem.Name}} = append({{$elem.Name}}, {{$elem.Type}}(val))
	}
	{{else}}

	// non-string
	
	in_{{$elem.Name}}_cdt := C.get_cdt(parameters, {{$index}})
	
	pcdt_{{$elem.Type}}_{{$elem.Name}} := ((*C.struct_cdt_openffi_{{$elem.Type}})(C.convert_union_to_ptr(unsafe.Pointer(&in_{{$elem.Name}}_cdt.cdt_val))))
	
	var in_{{$elem.Name}} C.openffi_{{$elem.Type}} = pcdt_{{$elem.Type}}_{{$elem.Name}}.val
	{{$elem.Name}} := {{if eq $elem.Type "bool"}}in_{{$elem.Name}} != C.openffi_bool(0){{else}}{{$elem.Type}}(in_{{$elem.Name}}){{end}}
	
	{{end}}
	{{end}}
	{{end}}
	
	// call original function
	
	{{range $index, $elem := $f.ReturnValues}}{{if $index}},{{end}}{{$elem.Name}}{{end}}{{if $f.ReturnValues}} := {{end}}{{$f.PathToForeignFunction.function}}({{range $index, $elem := $f.Parameters}}{{if $index}},{{end}}{{if $elem.IsTypeAlias}}{{$elem.GetTypeOrAlias}}({{$elem.Name}}){{else}}{{$elem.Name}}{{end}}{{end}})

	println("After call")

	// return values
	{{range $index, $elem := $f.ReturnValues}}
	if err, isError := interface{}({{$elem.Name}}).(error); isError{ // in case of error
		errToOutError(out_err, out_err_len, "Error returned", err)
		return
	} else { // Convert return values from Go to C

		{{if $elem.IsString}}
		{{if gt $elem.Dimensions 0}}
		// string array

		out_{{$elem.Name}} := (*C.openffi_{{$elem.Type}})(C.malloc(C.ulong(len({{$elem.Name}}))*{{Sizeof $elem}}))
		out_{{$elem.Name}}_sizes := (*C.openffi_size)(C.malloc(C.ulong(len({{$elem.Name}}))*C.sizeof_openffi_size))
		out_{{$elem.Name}}_dimensions := C.openffi_size( 1 )
		out_{{$elem.Name}}_dimensions_lengths := (*C.openffi_size)(C.malloc(C.sizeof_openffi_size * (out_{{$elem.Name}}_dimensions)))
		*out_{{$elem.Name}}_dimensions_lengths = C.ulong(len({{$elem.Name}}))
		
		for i, val := range {{$elem.Name}}{
			C.set_openffi_{{$elem.Type}}_element(out_{{$elem.Name}}, out_{{$elem.Name}}_sizes, C.int(i), C.openffi_{{$elem.Type}}(C.CString(val)), C.openffi_size(len(val)))
		}

		out_{{$elem.Name}}_cdt := C.get_cdt(return_values, {{$index}})

		C.set_cdt_type(out_{{$elem.Name}}_cdt, C.openffi_{{$elem.Type}}_array_type)
		out_{{$elem.Name}}_cdt.free_required = 1

		pcdt_out_{{$elem.Type}}_{{$elem.Name}} := ((*C.struct_cdt_openffi_{{$elem.Type}}_array)(C.convert_union_to_ptr(unsafe.Pointer(&out_{{$elem.Name}}_cdt.cdt_val))))
		pcdt_out_{{$elem.Type}}_{{$elem.Name}}.vals = out_{{$elem.Name}}
		pcdt_out_{{$elem.Type}}_{{$elem.Name}}.vals_sizes = out_{{$elem.Name}}_sizes
		pcdt_out_{{$elem.Type}}_{{$elem.Name}}.dimensions_lengths = out_{{$elem.Name}}_dimensions_lengths
		pcdt_out_{{$elem.Type}}_{{$elem.Name}}.dimensions = out_{{$elem.Name}}_dimensions
		
		{{else}}
		// string
		out_{{$elem.Name}}_len := C.openffi_size(C.ulong(len({{$elem.Name}})))
		out_{{$elem.Name}} := C.CString({{$elem.Name}})
		
		out_{{$elem.Name}}_cdt := C.get_cdt(return_values, {{$index}})
		C.set_cdt_type(out_{{$elem.Name}}_cdt, C.openffi_{{$elem.Type}}_type)
		out_{{$elem.Name}}_cdt.free_required = 1

		pcdt_out_{{$elem.Type}}_{{$elem.Name}} := ((*C.struct_cdt_openffi_{{$elem.Type}})(C.convert_union_to_ptr(unsafe.Pointer(&out_{{$elem.Name}}_cdt.cdt_val))))
		pcdt_out_{{$elem.Type}}_{{$elem.Name}}.val = out_{{$elem.Name}}
		pcdt_out_{{$elem.Type}}_{{$elem.Name}}.length = out_{{$elem.Name}}_len

		{{end}}{{else}}{{if gt $elem.Dimensions 0}}
		// non-string array
		
		out_{{$elem.Name}}_dimensions := C.openffi_size(1)
		out_{{$elem.Name}}_dimensions_lengths := (*C.openffi_size)(C.malloc(out_{{$elem.Name}}_dimensions*C.sizeof_openffi_{{$elem.Type}}))
		*out_{{$elem.Name}}_dimensions_lengths = C.ulong(len({{$elem.Name}}))

		out_{{$elem.Name}} := (*C.openffi_{{$elem.Type}})(C.malloc(C.ulong(len({{$elem.Name}}))*{{Sizeof $elem}}))
		for i, val := range {{$elem.Name}}{
			C.set_openffi_{{$elem.Type}}_element(out_{{$elem.Name}}, C.int(i), C.openffi_{{$elem.Type}}(val))
		}

		out_{{$elem.Name}}_cdt := C.get_cdt(return_values, {{$index}})
		C.set_cdt_type(out_{{$elem.Name}}_cdt, C.openffi_{{$elem.Type}}_array_type)
		out_{{$elem.Name}}_cdt.free_required = 1

		pcdt_out_{{$elem.Type}}_{{$elem.Name}} := ((*C.struct_cdt_openffi_{{$elem.Type}}_array)(C.convert_union_to_ptr(unsafe.Pointer(&out_{{$elem.Name}}_cdt.cdt_val))))

		pcdt_out_{{$elem.Type}}_{{$elem.Name}}.vals = out_{{$elem.Name}}
		pcdt_out_{{$elem.Type}}_{{$elem.Name}}.dimensions_lengths = out_{{$elem.Name}}_dimensions_lengths
		pcdt_out_{{$elem.Type}}_{{$elem.Name}}.dimensions = out_{{$elem.Name}}_dimensions

		{{else}}
		// non-string
		
		{{if $elem.IsBool}}
		var out_{{$elem.Name}} C.openffi_bool
		if {{$elem.Name}} { 
			out_{{$elem.Name}} = C.openffi_bool(1)
		} else { 
			out_{{$elem.Name}} = C.openffi_bool(0)
		}
		{{else}}
		out_{{$elem.Name}} := C.openffi_{{$elem.Type}}({{$elem.Name}})
		{{end}}
		
		out_{{$elem.Name}}_cdt := C.get_cdt(return_values, {{$index}})
		C.set_cdt_type(out_{{$elem.Name}}_cdt, C.openffi_{{$elem.Type}}_type)
		out_{{$elem.Name}}_cdt.free_required = C.openffi_bool(0)
		pcdt_out_{{$elem.Type}}_{{$elem.Name}} := ((*C.struct_cdt_openffi_{{$elem.Type}})(C.convert_union_to_ptr(unsafe.Pointer(&out_{{$elem.Name}}_cdt.cdt_val))))

		pcdt_out_{{$elem.Type}}_{{$elem.Name}}.val = out_{{$elem.Name}}
		{{end}}
		{{end}}
	}	
	{{end}}	

}
{{end}}{{end}}
`
