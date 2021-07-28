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

openffi_type get_cdt_type(struct cdt* p)
{
	return p->type;
}

void set_go_runtime_flag()
{
	xllr_set_runtime_flag("go_runtime", 10);
}

void* xllr_go_handle;
openffi_handle(*pset_object)(void*);
void*(*pget_object)(openffi_handle);

const char* load_xllr_go_lib_api()
{
	const char* openffi_home = getenv("OPENFFI_HOME");
	if(!openffi_home)
	{
		return "OPENFFI_HOME is not set";
	}

	#ifdef _WIN32
	const char* ext = ".dll";
#elif __APPLE__
	const char* ext = ".dylib";
#else
	const char* ext = ".so";
#endif
	
	char xllr_go_full_path[300] = {0};
	sprintf(xllr_go_full_path, "%s/xllr.go%s", openffi_home, ext);

	char* out_err;
	xllr_go_handle = load_library(xllr_go_full_path, &out_err);
	if(!xllr_go_handle)
	{
		// error has occurred
		printf("Failed to load XLLR Go: %s\n", out_err);
		return "Failed to load XLLR Go";
	}

	pset_object = (openffi_handle(*)(void*))load_symbol(xllr_go_handle, "set_object", &out_err);
	if(!pset_object)
	{
		// error has occurred
		printf("Failed to load set_object: %s\n", out_err);
		return "Failed to load set_object";
	}

	pget_object = (void*(*)(openffi_handle))load_symbol(xllr_go_handle, "get_object", &out_err);
	if(!pget_object)
	{
		// error has occurred
		printf("Failed to load get_object: %s\n", out_err);
		return "Failed to load get_object";
	}

	return 0;
}

openffi_handle set_object(void* obj)
{
	return pset_object(obj);
}

void* get_object(openffi_handle h)
{
	return pget_object(h);
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
openffi_type get_cdt_type(struct cdt* p);
void set_go_runtime_flag();
const char* load_xllr_go_lib_api();
openffi_handle set_object(void* obj);
void* get_object(openffi_handle h);
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
	
	C.set_go_runtime_flag()

	createObjectsTable()
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

type handle unsafe.Pointer

var pointers []unsafe.Pointer
var objects []interface{}

func createObjectsTable(){
	
	err := C.load_xllr_go_lib_api()
	if err != nil{ panic(C.GoString(err)) }
	
	pointers = make([]unsafe.Pointer, 0)
	objects = make([]interface{}, 0)
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

	{{if $elem.IsAny}}
	// any
	var {{$elem.Name}} interface{}
	in_{{$elem.Name}}_cdt := C.get_cdt(parameters, {{$index}})
	{{$elem.Name}}_type := C.get_cdt_type(in_{{$elem.Name}}_cdt)
	switch {{$elem.Name}}_type{

		case {{GetOpenFFIType "handle"}}: // handle
			pcdt_in_handle_{{$elem.Name}} := ((*C.struct_cdt_openffi_handle)(C.convert_union_to_ptr(unsafe.Pointer(&in_{{$elem.Name}}_cdt.cdt_val))))
			var in_{{$elem.Name}} C.openffi_handle = pcdt_in_handle_{{$elem.Name}}.val			

			{{$elem.Name}} = C.get_object(in_{{$elem.Name}})
			if {{$elem.Name}} == nil{ // handle belongs to another language 
				{{$elem.Name}} = handle(in_{{$elem.Name}})
			}

		case {{GetOpenFFIArrayType "handle"}}: // []handle
			pcdt_in_handle_{{$elem.Name}} := ((*C.struct_cdt_openffi_handle_array)(C.convert_union_to_ptr(unsafe.Pointer(&in_{{$elem.Name}}_cdt.cdt_val))))
			var in_{{$elem.Name}} *C.openffi_handle = pcdt_in_handle_{{$elem.Name}}.vals
			var in_{{$elem.Name}}_dimensions_lengths *C.openffi_size = pcdt_in_handle_{{$elem.Name}}.dimensions_lengths
			// var in_{{$elem.Name}}_dimensions C.openffi_size = pcdt_in_handle_{{$elem.Name}}.dimensions - TODO: not used until multi-dimensions support!

			{{$elem.Name}}_typed := make([]interface{}, 0)
			for i:=C.int(0) ; i<C.int(C.int(C.get_int_item(in_{{$elem.Name}}_dimensions_lengths, 0))) ; i++{
				val := C.get_openffi_handle_element(in_{{$elem.Name}}, C.int(i))

				val_obj := C.get_object(val)
				if val_obj == nil{ // handle belongs to
					{{$elem.Name}}_typed = append({{$elem.Name}}_typed, val)
				} else {
					{{$elem.Name}}_typed = append({{$elem.Name}}_typed, val_obj)
				}
			}
			{{$elem.Name}} = {{$elem.Name}}_typed

		{{range $numTypeEnumIndex, $numType := GetNumericTypes}}{{if ne $numType "handle"}}
		case {{GetOpenFFIType $numType}}: // {{$numType}}
			pcdt_in_{{$numType}}_{{$elem.Name}} := ((*C.struct_cdt_openffi_{{$numType}})(C.convert_union_to_ptr(unsafe.Pointer(&in_{{$elem.Name}}_cdt.cdt_val))))
			var in_{{$elem.Name}} C.openffi_{{$numType}} = pcdt_in_{{$numType}}_{{$elem.Name}}.val
			
			{{$elem.Name}} = {{$numType}}(in_{{$elem.Name}})

		{{end}}{{end}}

		{{range $numTypeEnumIndex, $numType := GetNumericTypes}}{{if ne $numType "handle"}}
		case {{GetOpenFFIArrayType $numType}}: // []{{$numType}}
			pcdt_in_{{$numType}}_{{$elem.Name}} := ((*C.struct_cdt_openffi_{{$numType}}_array)(C.convert_union_to_ptr(unsafe.Pointer(&in_{{$elem.Name}}_cdt.cdt_val))))
			var in_{{$elem.Name}} *C.openffi_{{$numType}} = pcdt_in_{{$numType}}_{{$elem.Name}}.vals
			var in_{{$elem.Name}}_dimensions_lengths *C.openffi_size = pcdt_in_{{$numType}}_{{$elem.Name}}.dimensions_lengths
			// var in_{{$elem.Name}}_dimensions C.openffi_size = pcdt_in_{{$numType}}_{{$elem.Name}}.dimensions - TODO: not used until multi-dimensions support!
					
			{{$elem.Name}}_typed := make([]{{$numType}}, 0)
			for i:=C.int(0) ; i<C.int(C.int(C.get_int_item(in_{{$elem.Name}}_dimensions_lengths, 0))) ; i++{
				val := C.get_openffi_{{$numType}}_element(in_{{$elem.Name}}, C.int(i))
				{{$elem.Name}}_typed = append({{$elem.Name}}_typed, {{$numType}}(val))
			}
			{{$elem.Name}} = {{$elem.Name}}_typed
		{{end}}{{end}}

		{{range $numTypeEnumIndex, $stringType := GetOpenFFIStringTypes}}
		case {{GetOpenFFIType $stringType}}: // {{$stringType}}
			in_{{$elem.Name}}_cdt := C.get_cdt(parameters, {{$index}})
			pcdt_in_{{$stringType}}_{{$elem.Name}} := ((*C.struct_cdt_openffi_{{$stringType}})(C.convert_union_to_ptr(unsafe.Pointer(&in_{{$elem.Name}}_cdt.cdt_val))))
			var in_{{$elem.Name}}_len C.openffi_size = pcdt_in_{{$stringType}}_{{$elem.Name}}.length
			var in_{{$elem.Name}} C.openffi_{{$stringType}} = pcdt_in_{{$stringType}}_{{$elem.Name}}.val
		
			{{$elem.Name}} = C.GoStringN(in_{{$elem.Name}}, C.int(in_{{$elem.Name}}_len))
		{{end}}

		{{range $numTypeEnumIndex, $stringType := GetOpenFFIStringTypes}}
		case {{GetOpenFFIArrayType $stringType}}: // []{{$stringType}}
			in_{{$elem.Name}}_cdt := C.get_cdt(parameters, {{$index}})
			pcdt_in_{{$stringType}}_{{$elem.Name}} := ((*C.struct_cdt_openffi_{{$stringType}}_array)(C.convert_union_to_ptr(unsafe.Pointer(&in_{{$elem.Name}}_cdt.cdt_val))))
		
			var in_{{$elem.Name}} *C.openffi_{{$stringType}} = pcdt_in_{{$stringType}}_{{$elem.Name}}.vals
			var in_{{$elem.Name}}_sizes *C.openffi_size = pcdt_in_{{$stringType}}_{{$elem.Name}}.vals_sizes
			var in_{{$elem.Name}}_dimensions_lengths *C.openffi_size = pcdt_in_{{$stringType}}_{{$elem.Name}}.dimensions_lengths
			//var in_{{$elem.Name}}_dimensions C.openffi_size = pcdt_in_{{$stringType}}_{{$elem.Name}}.dimensions - TODO: not used until multi-dimensions support!
		
			{{$elem.Name}}_typed := make([]string, 0, int(C.get_int_item(in_{{$elem.Name}}_dimensions_lengths, 0)))
			for i:=C.int(0) ; i<C.int(C.get_int_item(in_{{$elem.Name}}_dimensions_lengths, 0)) ; i++{
				var str_size C.openffi_size
				str := C.get_openffi_{{$stringType}}_element(in_{{$elem.Name}}, C.int(i), in_{{$elem.Name}}_sizes, &str_size)
				{{$elem.Name}}_typed = append({{$elem.Name}}_typed, C.GoStringN(str, C.int(str_size)))
			}
			{{$elem.Name}} = {{$elem.Name}}_typed
		{{end}}


		case {{GetOpenFFIType "bool"}}: // bool
			in_{{$elem.Name}}_cdt := C.get_cdt(parameters, {{$index}})
			pcdt_in_bool_{{$elem.Name}} := ((*C.struct_cdt_openffi_bool)(C.convert_union_to_ptr(unsafe.Pointer(&in_{{$elem.Name}}_cdt.cdt_val))))
			var in_{{$elem.Name}} C.openffi_bool = pcdt_in_bool_{{$elem.Name}}.val
			
			{{$elem.Name}} = in_{{$elem.Name}} != C.openffi_bool(0)

		case {{GetOpenFFIArrayType "bool"}}: // []bool
			in_{{$elem.Name}}_cdt := C.get_cdt(parameters, {{$index}})
			pcdt_in_bool_{{$elem.Name}} := ((*C.struct_cdt_openffi_bool_array)(C.convert_union_to_ptr(unsafe.Pointer(&in_{{$elem.Name}}_cdt.cdt_val))))
			var in_{{$elem.Name}} *C.openffi_bool = pcdt_in_bool_{{$elem.Name}}.vals
			var in_{{$elem.Name}}_dimensions_lengths *C.openffi_size = pcdt_in_bool_{{$elem.Name}}.dimensions_lengths
			// var in_{{$elem.Name}}_dimensions C.openffi_size = pcdt_in_bool_{{$elem.Name}}.dimensions - TODO: not used until multi-dimensions support!
					
			{{$elem.Name}}_typed := make([]bool, 0)
			for i:=C.int(0) ; i<C.int(C.int(C.get_int_item(in_{{$elem.Name}}_dimensions_lengths, 0))) ; i++{
				val := C.get_openffi_bool_element(in_{{$elem.Name}}, C.int(i))
				var bval bool
				if val != 0 { bval = true } else { bval = false }
				{{$elem.Name}}_typed = append({{$elem.Name}}_typed, bval)
			}

			{{$elem.Name}} = {{$elem.Name}}_typed

		default:
			panic(fmt.Errorf("Return value %v is not of a supported type, but of type: %v", "{{$elem.Name}}", {{$elem.Name}}_type))
	}
	{{else if $elem.IsString}}
	
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
	
	{{if $elem.IsHandle}}
	//handle_{{$elem.Name}} := handle((in_{{$elem.Name}}))
	{{$elem.Name}}_obj := C.get_object(C.openffi_handle(in_{{$elem.Name}}))
	if {{$elem.Name}}_obj == nil{ panic(fmt.Errorf("Failed to find object")) }
	{{$elem.Name}} := *((*{{$f.PathToForeignFunction.class}})({{$elem.Name}}_obj))
	{{else}}
	{{$elem.Name}} := {{if eq $elem.Type "bool"}}in_{{$elem.Name}} != C.openffi_bool(0){{else}}{{$elem.Type}}(in_{{$elem.Name}}){{end}}
	{{end}}

	{{end}}
	{{end}}
	{{end}}
	
	// call original function
	{{if $f.IsMethod}}
	{{range $index, $elem := $f.ReturnValues}}{{if $index}},{{end}}{{$elem.Name}}{{end}}{{if $f.ReturnValues}} := {{end}}{{(index $f.Parameters 0).Name }}.{{$f.PathToForeignFunction.function}}({{range $index, $elem := $f.Parameters}}{{if $index}}{{if gt $index 1}},{{end}}{{if $elem.IsTypeAlias}}{{$elem.GetTypeOrAlias}}({{$elem.Name}}){{else}}{{$elem.Name}}{{end}}{{end}}{{end}})
	{{else}}
	{{range $index, $elem := $f.ReturnValues}}{{if $index}},{{end}}{{$elem.Name}}{{end}}{{if $f.ReturnValues}} := {{end}}{{$f.PathToForeignFunction.function}}({{range $index, $elem := $f.Parameters}}{{if $index}},{{end}}{{if $elem.IsTypeAlias}}{{$elem.GetTypeOrAlias}}({{$elem.Name}}){{else}}{{$elem.Name}}{{end}}{{end}})
	{{end}}

	// return values
	{{range $index, $elem := $f.ReturnValues}}
	if err, isError := interface{}({{$elem.Name}}).(error); isError{ // in case of error
		errToOutError(out_err, out_err_len, "Error returned", err)
		return
	} else { // Convert return values from Go to C

		{{if $elem.IsAny}}
		switch {{$elem.Name}}.(type) {

		{{ range $numTypeIndex, $numType := GetNumericTypes }}
		case {{$numType}}:
			out_{{$elem.Name}} := C.openffi_{{$numType}}({{$elem.Name}}.({{$numType}}))
			out_{{$elem.Name}}_cdt := C.get_cdt(return_values, {{$index}})
			C.set_cdt_type(out_{{$elem.Name}}_cdt, C.openffi_{{$numType}}_type)
			out_{{$elem.Name}}_cdt.free_required = 1
			pcdt_out_{{$numType}}_{{$elem.Name}} := ((*C.struct_cdt_openffi_{{$numType}})(C.convert_union_to_ptr(unsafe.Pointer(&out_{{$elem.Name}}_cdt.cdt_val))))
			pcdt_out_{{$numType}}_{{$elem.Name}}.val = out_{{$elem.Name}}

		{{end}}
		

		{{ range $index, $numType := GetNumericTypes }}
		case []{{$numType}}:
			out_{{$elem.Name}}_dimensions := C.openffi_size(1)
			out_{{$elem.Name}}_dimensions_lengths := (*C.openffi_size)(C.malloc(C.sizeof_openffi_size))
			*out_{{$elem.Name}}_dimensions_lengths = C.ulong(len({{$elem.Name}}.([]{{$numType}})))
		
			out_{{$elem.Name}} := (*C.openffi_{{$numType}})(C.malloc(C.ulong(len({{$elem.Name}}.([]{{$numType}})))*C.sizeof_openffi_{{$numType}}))
			for i, val := range {{$elem.Name}}.([]{{$numType}}){
				C.set_openffi_{{$numType}}_element(out_{{$elem.Name}}, C.int(i), C.openffi_{{$numType}}(val))
			}
		
			out_{{$elem.Name}}_cdt := C.get_cdt(return_values, {{$index}})
			C.set_cdt_type(out_{{$elem.Name}}_cdt, C.openffi_{{$numType}}_array_type)
			out_{{$elem.Name}}_cdt.free_required = 1
			pcdt_out_{{$numType}}_{{$elem.Name}} := ((*C.struct_cdt_openffi_{{$numType}}_array)(C.convert_union_to_ptr(unsafe.Pointer(&out_{{$elem.Name}}_cdt.cdt_val))))
			pcdt_out_{{$numType}}_{{$elem.Name}}.vals = out_{{$elem.Name}}
			pcdt_out_{{$numType}}_{{$elem.Name}}.dimensions_lengths = out_{{$elem.Name}}_dimensions_lengths
			pcdt_out_{{$numType}}_{{$elem.Name}}.dimensions = out_{{$elem.Name}}_dimensions

		{{end}}

		case int:
			out_{{$elem.Name}} := C.openffi_int64(int64({{$elem.Name}}.(int)))
			out_{{$elem.Name}}_cdt := C.get_cdt(return_values, {{$index}})
			C.set_cdt_type(out_{{$elem.Name}}_cdt, C.openffi_int64_type)
			out_{{$elem.Name}}_cdt.free_required = 1
			pcdt_out_int64_{{$elem.Name}} := ((*C.struct_cdt_openffi_int64)(C.convert_union_to_ptr(unsafe.Pointer(&out_{{$elem.Name}}_cdt.cdt_val))))
			pcdt_out_int64_{{$elem.Name}}.val = out_{{$elem.Name}}

		case []int:
			out_{{$elem.Name}}_dimensions := C.openffi_size(1)
			out_{{$elem.Name}}_dimensions_lengths := (*C.openffi_size)(C.malloc(C.sizeof_openffi_size))
			*out_{{$elem.Name}}_dimensions_lengths = C.ulong(len({{$elem.Name}}.([]int)))
		
			out_{{$elem.Name}} := (*C.openffi_int64)(C.malloc(C.ulong(len({{$elem.Name}}.([]int)))*C.sizeof_openffi_int64))
			for i, val := range {{$elem.Name}}.([]int){
				C.set_openffi_int64_element(out_{{$elem.Name}}, C.int(i), C.openffi_int64(val))
			}
		
			out_{{$elem.Name}}_cdt := C.get_cdt(return_values, {{$index}})
			C.set_cdt_type(out_{{$elem.Name}}_cdt, C.openffi_int64_array_type)
			out_{{$elem.Name}}_cdt.free_required = 1
			pcdt_out_int64_{{$elem.Name}} := ((*C.struct_cdt_openffi_int64_array)(C.convert_union_to_ptr(unsafe.Pointer(&out_{{$elem.Name}}_cdt.cdt_val))))
			pcdt_out_int64_{{$elem.Name}}.vals = out_{{$elem.Name}}
			pcdt_out_int64_{{$elem.Name}}.dimensions_lengths = out_{{$elem.Name}}_dimensions_lengths
			pcdt_out_int64_{{$elem.Name}}.dimensions = out_{{$elem.Name}}_dimensions

		case bool:
			var out_{{$elem.Name}} C.openffi_bool
			if {{$elem.Name}}.(bool) { out_{{$elem.Name}} = C.openffi_bool(1) } else { out_{{$elem.Name}} = C.openffi_bool(0) }
			out_{{$elem.Name}}_cdt := C.get_cdt(return_values, {{$index}})
			C.set_cdt_type(out_{{$elem.Name}}_cdt, C.openffi_bool_type)
			out_{{$elem.Name}}_cdt.free_required = 1
			pcdt_out_bool_{{$elem.Name}} := ((*C.struct_cdt_openffi_bool)(C.convert_union_to_ptr(unsafe.Pointer(&out_{{$elem.Name}}_cdt.cdt_val))))
			pcdt_out_bool_{{$elem.Name}}.val = out_{{$elem.Name}}

		case string:
			out_{{$elem.Name}}_len := C.openffi_size(C.ulong(len({{$elem.Name}}.(string))))
			out_{{$elem.Name}} := C.CString({{$elem.Name}}.(string))
			out_{{$elem.Name}}_cdt := C.get_cdt(return_values, {{$index}})
			C.set_cdt_type(out_{{$elem.Name}}_cdt, C.openffi_string8_type)
			out_{{$elem.Name}}_cdt.free_required = 1
			pcdt_out_string8_{{$elem.Name}} := ((*C.struct_cdt_openffi_string8)(C.convert_union_to_ptr(unsafe.Pointer(&out_{{$elem.Name}}_cdt.cdt_val))))
			pcdt_out_string8_{{$elem.Name}}.val = out_{{$elem.Name}}
			pcdt_out_string8_{{$elem.Name}}.length = out_{{$elem.Name}}_len

		case []bool:
			out_{{$elem.Name}}_dimensions := C.openffi_size(1)
			out_{{$elem.Name}}_dimensions_lengths := (*C.openffi_size)(C.malloc(C.sizeof_openffi_size))
			*out_{{$elem.Name}}_dimensions_lengths = C.openffi_size(len({{$elem.Name}}.([]bool)))
		
			out_{{$elem.Name}} := (*C.openffi_bool)(C.malloc(C.openffi_size(len({{$elem.Name}}.([]bool)))*C.sizeof_openffi_bool))
			for i, val := range {{$elem.Name}}.([]bool){
				var bval C.openffi_bool
				if val { bval = C.openffi_bool(1) } else { bval = C.openffi_bool(0) }
				C.set_openffi_bool_element(out_{{$elem.Name}}, C.int(i), C.openffi_bool(bval))
			}
		
			out_{{$elem.Name}}_cdt := C.get_cdt(return_values, {{$index}})
			C.set_cdt_type(out_{{$elem.Name}}_cdt, C.openffi_bool_array_type)
			out_{{$elem.Name}}_cdt.free_required = 1
			pcdt_out_bool_{{$elem.Name}} := ((*C.struct_cdt_openffi_bool_array)(C.convert_union_to_ptr(unsafe.Pointer(&out_{{$elem.Name}}_cdt.cdt_val))))
			pcdt_out_bool_{{$elem.Name}}.vals = out_{{$elem.Name}}
			pcdt_out_bool_{{$elem.Name}}.dimensions_lengths = out_{{$elem.Name}}_dimensions_lengths
			pcdt_out_bool_{{$elem.Name}}.dimensions = out_{{$elem.Name}}_dimensions

		case []string:
			out_{{$elem.Name}} := (*C.openffi_string8)(C.malloc(C.ulong(len({{$elem.Name}}.([]string)))*C.sizeof_openffi_string8))
			out_{{$elem.Name}}_sizes := (*C.openffi_size)(C.malloc(C.ulong(len({{$elem.Name}}.([]string)))*C.sizeof_openffi_size))
			out_{{$elem.Name}}_dimensions := C.openffi_size(1)
			out_{{$elem.Name}}_dimensions_lengths := (*C.openffi_size)(C.malloc(C.sizeof_openffi_size * (out_{{$elem.Name}}_dimensions)))
			*out_{{$elem.Name}}_dimensions_lengths = C.openffi_size(len({{$elem.Name}}.([]string)))
			
			for i, val := range {{$elem.Name}}.([]string){
				C.set_openffi_string8_element(out_{{$elem.Name}}, out_{{$elem.Name}}_sizes, C.int(i), C.openffi_string8(C.CString(val)), C.openffi_size(len(val)))
			}
			
			out_{{$elem.Name}}_cdt := C.get_cdt(return_values, {{$index}})
			C.set_cdt_type(out_{{$elem.Name}}_cdt, C.openffi_string8_array_type)
			out_{{$elem.Name}}_cdt.free_required = 1
			pcdt_out_string8_{{$elem.Name}} := ((*C.struct_cdt_openffi_string8_array)(C.convert_union_to_ptr(unsafe.Pointer(&out_{{$elem.Name}}_cdt.cdt_val))))
			pcdt_out_string8_{{$elem.Name}}.vals = out_{{$elem.Name}}
			pcdt_out_string8_{{$elem.Name}}.vals_sizes = out_{{$elem.Name}}_sizes
			pcdt_out_string8_{{$elem.Name}}.dimensions_lengths = out_{{$elem.Name}}_dimensions_lengths
			pcdt_out_string8_{{$elem.Name}}.dimensions = out_{{$elem.Name}}_dimensions
			
		default:
			// Turn object to a handle
			p{{$elem.Name}} := &{{$elem.Name}}
			unsafe_p{{$elem.Name}} := unsafe.Pointer(p{{$elem.Name}})
			objects = append(objects, p{{$elem.Name}})
			pointers = append(pointers, unsafe_p{{$elem.Name}})
	
			{{$elem.Name}}_handle := C.set_object(unsafe_p{{$elem.Name}})
			
			out_{{$elem.Name}} := C.openffi_handle({{$elem.Name}}_handle)
			out_{{$elem.Name}}_cdt := C.get_cdt(return_values, {{$index}})
			C.set_cdt_type(out_{{$elem.Name}}_cdt, C.openffi_handle_type)
			out_{{$elem.Name}}_cdt.free_required = 1
			pcdt_out_handle_{{$elem.Name}} := ((*C.struct_cdt_openffi_handle)(C.convert_union_to_ptr(unsafe.Pointer(&out_{{$elem.Name}}_cdt.cdt_val))))
			pcdt_out_handle_{{$elem.Name}}.val = out_{{$elem.Name}}
		}
		
		{{else if $elem.IsString}}
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
