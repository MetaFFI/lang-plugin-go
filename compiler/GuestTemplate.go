package main

const GuestHeaderTemplate = `
// Code generated by MetaFFI. Modify only in marked places.
// Guest code for {{.IDLFilenameWithExtension}}

package main
`

const GuestImportsTemplate = `
import "fmt"
import "unsafe"
import "reflect"
import "github.com/pkg/profile"
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
#cgo CFLAGS: -I{{GetEnvVar "METAFFI_HOME" true}}


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
#cgo CFLAGS: -I{{GetEnvVar "METAFFI_HOME" true}}

#include <include/cdt_structs.h>
#include <include/cdt_capi_loader.h>

metaffi_size get_int_item(metaffi_size* array, int index);
void* convert_union_to_ptr(void* p);
void set_cdt_type(struct cdt* p, metaffi_type t);
metaffi_type get_cdt_type(struct cdt* p);
void set_go_runtime_flag();
struct cdt* get_cdt_element(struct cdts* pdata, int cdts_index);

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
{{/* TODO: Make function for each type */}}
func fromCDTToGo(pdata *C.struct_cdts, cdtsIndex int, i int) interface{}{

    data := C.get_cdt_element(pdata, C.int(cdtsIndex))

	var res interface{}
	index := C.int(i)
	in_res_cdt := C.get_cdt(data, index)
	res_type := C.get_cdt_type(in_res_cdt)
	switch res_type{

		case {{GetMetaFFIType "handle"}}: // handle
			pcdt_in_handle_res := ((*C.struct_cdt_metaffi_handle)(C.convert_union_to_ptr(unsafe.Pointer(&in_res_cdt.cdt_val))))
			var in_res C.metaffi_handle = pcdt_in_handle_res.val

			res = GetObject(Handle(in_res))
			if res == nil{ // handle belongs to another language 
				res = Handle(in_res)
			}

		case {{GetMetaFFIArrayType "handle"}}: // []Handle
			pcdt_in_handle_res := ((*C.struct_cdt_metaffi_handle_array)(C.convert_union_to_ptr(unsafe.Pointer(&in_res_cdt.cdt_val))))
			var in_res *C.metaffi_handle = pcdt_in_handle_res.vals
			var in_res_dimensions_lengths *C.metaffi_size = pcdt_in_handle_res.dimensions_lengths
			// var in_res_dimensions C.metaffi_size = pcdt_in_handle_res.dimensions - TODO: not used until multi-dimensions support!

			res_typed := make([]interface{}, 0)
			for i:=C.int(0) ; i<C.int(C.int(C.get_int_item(in_res_dimensions_lengths, 0))) ; i++{
				val := C.get_metaffi_handle_element(in_res, C.int(i))

				val_obj := GetObject(Handle(val))
				if val_obj == nil{ // handle belongs to
					res_typed = append(res_typed, Handle(val))
				} else {
					res_typed = append(res_typed, val_obj)
				}
			}
			res = res_typed

		{{range $numTypeEnumIndex, $numType := GetNumericTypes}}{{if ne $numType "Handle"}}
		case {{GetMetaFFIType $numType}}: // {{$numType}}
			pcdt_in_{{$numType}}_res := ((*C.struct_cdt_{{ MakeMetaFFIType $numType}})(C.convert_union_to_ptr(unsafe.Pointer(&in_res_cdt.cdt_val))))
			var in_res C.{{ MakeMetaFFIType $numType}} = pcdt_in_{{$numType}}_res.val
			
			res = {{$numType}}(in_res)

		{{end}}{{end}}

		{{range $numTypeEnumIndex, $numType := GetNumericTypes}}{{if ne $numType "Handle"}}
		case {{GetMetaFFIArrayType $numType}}: // []{{$numType}}
			pcdt_in_{{$numType}}_res := ((*C.struct_cdt_{{ MakeMetaFFIType $numType}}_array)(C.convert_union_to_ptr(unsafe.Pointer(&in_res_cdt.cdt_val))))
			var in_res *C.{{ MakeMetaFFIType $numType}} = pcdt_in_{{$numType}}_res.vals
			var in_res_dimensions_lengths *C.metaffi_size = pcdt_in_{{$numType}}_res.dimensions_lengths
			// var in_res_dimensions C.metaffi_size = pcdt_in_{{$numType}}_res.dimensions - TODO: not used until multi-dimensions support!
					
			res_typed := make([]{{$numType}}, 0)
			for i:=C.int(0) ; i<C.int(C.int(C.get_int_item(in_res_dimensions_lengths, 0))) ; i++{
				val := C.get_{{ MakeMetaFFIType $numType}}_element(in_res, C.int(i))
				res_typed = append(res_typed, {{$numType}}(val))
			}
			res = res_typed
		{{end}}{{end}}

		{{range $numTypeEnumIndex, $stringType := GetMetaFFIStringTypes}}
		case {{GetMetaFFIType $stringType}}: // {{$stringType}}
			in_res_cdt := C.get_cdt(data, index)
			pcdt_in_{{$stringType}}_res := ((*C.struct_cdt_metaffi_{{$stringType}})(C.convert_union_to_ptr(unsafe.Pointer(&in_res_cdt.cdt_val))))
			var in_res_len C.metaffi_size = pcdt_in_{{$stringType}}_res.length
			var in_res C.metaffi_{{$stringType}} = pcdt_in_{{$stringType}}_res.val
		
			res = C.GoStringN(in_res, C.int(in_res_len))
		{{end}}

		{{range $numTypeEnumIndex, $stringType := GetMetaFFIStringTypes}}
		case {{GetMetaFFIArrayType $stringType}}: // []{{$stringType}}
			in_res_cdt := C.get_cdt(data, index)
			pcdt_in_{{$stringType}}_res := ((*C.struct_cdt_metaffi_{{$stringType}}_array)(C.convert_union_to_ptr(unsafe.Pointer(&in_res_cdt.cdt_val))))
		
			var in_res *C.metaffi_{{$stringType}} = pcdt_in_{{$stringType}}_res.vals
			var in_res_sizes *C.metaffi_size = pcdt_in_{{$stringType}}_res.vals_sizes
			var in_res_dimensions_lengths *C.metaffi_size = pcdt_in_{{$stringType}}_res.dimensions_lengths
			//var in_res_dimensions C.metaffi_size = pcdt_in_{{$stringType}}_res.dimensions - TODO: not used until multi-dimensions support!
		
			res_typed := make([]string, 0, int(C.get_int_item(in_res_dimensions_lengths, 0)))
			for i:=C.int(0) ; i<C.int(C.get_int_item(in_res_dimensions_lengths, 0)) ; i++{
				var str_size C.metaffi_size
				str := C.get_metaffi_{{$stringType}}_element(in_res, C.int(i), in_res_sizes, &str_size)
				res_typed = append(res_typed, C.GoStringN(str, C.int(str_size)))
			}
			res = res_typed
		{{end}}


		case {{GetMetaFFIType "bool"}}: // bool
			in_res_cdt := C.get_cdt(data, index)
			pcdt_in_bool_res := ((*C.struct_cdt_metaffi_bool)(C.convert_union_to_ptr(unsafe.Pointer(&in_res_cdt.cdt_val))))
			var in_res C.metaffi_bool = pcdt_in_bool_res.val
			
			res = in_res != C.metaffi_bool(0)

		case {{GetMetaFFIArrayType "bool"}}: // []bool
			in_res_cdt := C.get_cdt(data, index)
			pcdt_in_bool_res := ((*C.struct_cdt_metaffi_bool_array)(C.convert_union_to_ptr(unsafe.Pointer(&in_res_cdt.cdt_val))))
			var in_res *C.metaffi_bool = pcdt_in_bool_res.vals
			var in_res_dimensions_lengths *C.metaffi_size = pcdt_in_bool_res.dimensions_lengths
			// var in_res_dimensions C.metaffi_size = pcdt_in_bool_res.dimensions - TODO: not used until multi-dimensions support!
					
			res_typed := make([]bool, 0)
			for i:=C.int(0) ; i<C.int(C.int(C.get_int_item(in_res_dimensions_lengths, 0))) ; i++{
				val := C.get_metaffi_bool_element(in_res, C.int(i))
				var bval bool
				if val != 0 { bval = true } else { bval = false }
				res_typed = append(res_typed, bval)
			}

			res = res_typed

		default:
			panic(fmt.Errorf("Return value %v is not of a supported type, but of type: %v", "res", res_type))
	}

	return res
}
{{/* TODO: Make function for each type */}}
func fromGoToCDT(input interface{}, pdata *C.struct_cdts, cdtsIndex int, i int){

	data := C.get_cdt_element(pdata, C.int(cdtsIndex))

	index := C.int(i)
	switch input.(type) {

		{{ range $numTypeIndex, $numType := GetNumericTypes }}
		case {{$numType}}:
			out_input := C.{{ MakeMetaFFIType $numType}}(input.({{$numType}}))
			out_input_cdt := C.get_cdt(data, index)
			C.set_cdt_type(out_input_cdt, C.{{MakeMetaFFIType $numType}}_type)
			out_input_cdt.free_required = 1
			pcdt_out_{{$numType}}_input := ((*C.struct_cdt_{{ MakeMetaFFIType $numType}})(C.convert_union_to_ptr(unsafe.Pointer(&out_input_cdt.cdt_val))))
			pcdt_out_{{$numType}}_input.val = out_input

		{{end}}
		

		{{ range $numTypeIndex, $numType := GetNumericTypes }}
		case []{{$numType}}:
			out_input_dimensions := C.metaffi_size(1)
			out_input_dimensions_lengths := (*C.metaffi_size)(C.malloc(C.sizeof_metaffi_size))
			*out_input_dimensions_lengths = C.len_to_metaffi_size(C.longlong(len(input.([]{{$numType}}))))
		
			out_input := (*C.{{MakeMetaFFIType $numType}})(C.malloc(C.len_to_metaffi_size(C.longlong(len(input.([]{{$numType}}))))*C.sizeof_{{ MakeMetaFFIType $numType}}))
			for i, val := range input.([]{{$numType}}){
				C.set_{{MakeMetaFFIType $numType}}_element(out_input, C.int(i), C.{{ MakeMetaFFIType $numType}}(val))
			}
		
			out_input_cdt := C.get_cdt(data, index)
			C.set_cdt_type(out_input_cdt, C.{{MakeMetaFFIType $numType}}_array_type)
			out_input_cdt.free_required = 1
			pcdt_out_{{$numType}}_input := ((*C.struct_cdt_{{ MakeMetaFFIType $numType}}_array)(C.convert_union_to_ptr(unsafe.Pointer(&out_input_cdt.cdt_val))))
			pcdt_out_{{$numType}}_input.vals = out_input
			pcdt_out_{{$numType}}_input.dimensions_lengths = out_input_dimensions_lengths
			pcdt_out_{{$numType}}_input.dimensions = out_input_dimensions

		{{end}}

		case int:
			out_input := C.metaffi_int64(int64(input.(int)))
			out_input_cdt := C.get_cdt(data, index)
			C.set_cdt_type(out_input_cdt, C.metaffi_int64_type)
			out_input_cdt.free_required = 1
			pcdt_out_int64_input := ((*C.struct_cdt_metaffi_int64)(C.convert_union_to_ptr(unsafe.Pointer(&out_input_cdt.cdt_val))))
			pcdt_out_int64_input.val = out_input

		case []int:
			out_input_dimensions := C.metaffi_size(1)
			out_input_dimensions_lengths := (*C.metaffi_size)(C.malloc(C.sizeof_metaffi_size))
			*out_input_dimensions_lengths = C.len_to_metaffi_size(C.longlong(len(input.([]int))))
		
			out_input := (*C.metaffi_int64)(C.malloc(C.len_to_metaffi_size(C.longlong(len(input.([]int))))*C.sizeof_metaffi_int64))
			for i, val := range input.([]int){
				C.set_metaffi_int64_element(out_input, C.int(i), C.metaffi_int64(val))
			}
		
			out_input_cdt := C.get_cdt(data, index)
			C.set_cdt_type(out_input_cdt, C.metaffi_int64_array_type)
			out_input_cdt.free_required = 1
			pcdt_out_int64_input := ((*C.struct_cdt_metaffi_int64_array)(C.convert_union_to_ptr(unsafe.Pointer(&out_input_cdt.cdt_val))))
			pcdt_out_int64_input.vals = out_input
			pcdt_out_int64_input.dimensions_lengths = out_input_dimensions_lengths
			pcdt_out_int64_input.dimensions = out_input_dimensions

		case bool:
			var out_input C.metaffi_bool
			if input.(bool) { out_input = C.metaffi_bool(1) } else { out_input = C.metaffi_bool(0) }
			out_input_cdt := C.get_cdt(data, index)
			C.set_cdt_type(out_input_cdt, C.metaffi_bool_type)
			out_input_cdt.free_required = 1
			pcdt_out_bool_input := ((*C.struct_cdt_metaffi_bool)(C.convert_union_to_ptr(unsafe.Pointer(&out_input_cdt.cdt_val))))
			pcdt_out_bool_input.val = out_input

		case string:
			out_input_len := C.metaffi_size(C.len_to_metaffi_size(C.longlong(len(input.(string)))))
			out_input := C.CString(input.(string))
			out_input_cdt := C.get_cdt(data, index)
			C.set_cdt_type(out_input_cdt, C.metaffi_string8_type)
			out_input_cdt.free_required = 1
			pcdt_out_string8_input := ((*C.struct_cdt_metaffi_string8)(C.convert_union_to_ptr(unsafe.Pointer(&out_input_cdt.cdt_val))))
			pcdt_out_string8_input.val = out_input
			pcdt_out_string8_input.length = out_input_len

		case []bool:
			out_input_dimensions := C.metaffi_size(1)
			out_input_dimensions_lengths := (*C.metaffi_size)(C.malloc(C.sizeof_metaffi_size))
			*out_input_dimensions_lengths = C.metaffi_size(len(input.([]bool)))
		
			out_input := (*C.metaffi_bool)(C.malloc(C.metaffi_size(len(input.([]bool)))*C.sizeof_metaffi_bool))
			for i, val := range input.([]bool){
				var bval C.metaffi_bool
				if val { bval = C.metaffi_bool(1) } else { bval = C.metaffi_bool(0) }
				C.set_metaffi_bool_element(out_input, C.int(i), C.metaffi_bool(bval))
			}
		
			out_input_cdt := C.get_cdt(data, index)
			C.set_cdt_type(out_input_cdt, C.metaffi_bool_array_type)
			out_input_cdt.free_required = 1
			pcdt_out_bool_input := ((*C.struct_cdt_metaffi_bool_array)(C.convert_union_to_ptr(unsafe.Pointer(&out_input_cdt.cdt_val))))
			pcdt_out_bool_input.vals = out_input
			pcdt_out_bool_input.dimensions_lengths = out_input_dimensions_lengths
			pcdt_out_bool_input.dimensions = out_input_dimensions

		case []string:
			out_input := (*C.metaffi_string8)(C.malloc(C.len_to_metaffi_size(C.longlong(len(input.([]string))))*C.sizeof_metaffi_string8))
			out_input_sizes := (*C.metaffi_size)(C.malloc(C.len_to_metaffi_size(C.longlong(len(input.([]string))))*C.sizeof_metaffi_size))
			out_input_dimensions := C.metaffi_size(1)
			out_input_dimensions_lengths := (*C.metaffi_size)(C.malloc(C.sizeof_metaffi_size * (out_input_dimensions)))
			*out_input_dimensions_lengths = C.metaffi_size(len(input.([]string)))
			
			for i, val := range input.([]string){
				C.set_metaffi_string8_element(out_input, out_input_sizes, C.int(i), C.metaffi_string8(C.CString(val)), C.metaffi_size(len(val)))
			}
			
			out_input_cdt := C.get_cdt(data, index)
			C.set_cdt_type(out_input_cdt, C.metaffi_string8_array_type)
			out_input_cdt.free_required = 1
			pcdt_out_string8_input := ((*C.struct_cdt_metaffi_string8_array)(C.convert_union_to_ptr(unsafe.Pointer(&out_input_cdt.cdt_val))))
			pcdt_out_string8_input.vals = out_input
			pcdt_out_string8_input.vals_sizes = out_input_sizes
			pcdt_out_string8_input.dimensions_lengths = out_input_dimensions_lengths
			pcdt_out_string8_input.dimensions = out_input_dimensions
			
		default:
			
			if input == nil{ // return handle "0"
				out_input := C.metaffi_handle(uintptr(0))
				out_input_cdt := C.get_cdt(data, index)
				C.set_cdt_type(out_input_cdt, C.metaffi_handle_type)
				out_input_cdt.free_required = 0
				pcdt_out_handle_input := ((*C.struct_cdt_metaffi_handle)(C.convert_union_to_ptr(unsafe.Pointer(&out_input_cdt.cdt_val))))
				pcdt_out_handle_input.val = out_input
				return
			}

			// check if the object is type of a primitive
			inputVal := reflect.ValueOf(input)
			inputType := reflect.TypeOf(input)
			switch inputType.Kind(){
				case reflect.Bool: fromGoToCDT(bool(inputVal.Bool()), pdata, cdtsIndex, i); return

				case reflect.Float32: fromGoToCDT(float32(inputVal.Float()), pdata, cdtsIndex, i); return
				case reflect.Float64: fromGoToCDT(float64(inputVal.Float()), pdata, cdtsIndex, i); return
				
				case reflect.Int8: fromGoToCDT(int8(inputVal.Int()), pdata, cdtsIndex, i); return
				case reflect.Int16: fromGoToCDT(int16(inputVal.Int()), pdata, cdtsIndex, i); return
				case reflect.Int32: fromGoToCDT(int32(inputVal.Int()), pdata, cdtsIndex, i); return
				case reflect.Int: fallthrough
				case reflect.Int64: fromGoToCDT(int64(inputVal.Int()), pdata, cdtsIndex, i); return

				case reflect.Uint8: fromGoToCDT(uint8(inputVal.Uint()), pdata, cdtsIndex, i); return
				case reflect.Uint16: fromGoToCDT(uint16(inputVal.Uint()), pdata, cdtsIndex, i); return
				case reflect.Uint32: fromGoToCDT(uint32(inputVal.Uint()), pdata, cdtsIndex, i); return
				case reflect.Uint: fallthrough
				case reflect.Uint64: fromGoToCDT(uint64(inputVal.Uint()), pdata, cdtsIndex, i); return

				case reflect.Uintptr: fromGoToCDT(uint64(inputVal.UnsafeAddr()), pdata, cdtsIndex, i); return

				case reflect.String: fromGoToCDT(string(inputVal.String()), pdata, cdtsIndex, i); return

				case reflect.Slice:
					switch inputType.Elem().Kind(){
						case reflect.Float32:
							dstSlice := make([]float32, inputVal.Len(), inputVal.Cap())
							for i:=0 ; i < inputVal.Len() ; i++{ dstSlice[i] = float32(inputVal.Index(i).Float()) }
							fromGoToCDT(dstSlice, pdata, cdtsIndex, i)
							return

						case reflect.Float64:
							dstSlice := make([]float64, inputVal.Len(), inputVal.Cap())
							for i:=0 ; i < inputVal.Len() ; i++{ dstSlice[i] = float64(inputVal.Index(i).Float()) }
							fromGoToCDT(dstSlice, pdata, cdtsIndex, i)
							return

						case reflect.Bool:
							dstSlice := make([]bool, inputVal.Len(), inputVal.Cap())
							for i:=0 ; i < inputVal.Len() ; i++{ dstSlice[i] = inputVal.Index(i).Bool() }
							fromGoToCDT(dstSlice, pdata, cdtsIndex, i)
							return
				
						case reflect.Int8:
							dstSlice := make([]int8, inputVal.Len(), inputVal.Cap())
							for i:=0 ; i < inputVal.Len() ; i++{ dstSlice[i] = int8(inputVal.Index(i).Int()) }
							fromGoToCDT(dstSlice, pdata, cdtsIndex, i)
							return
							
						case reflect.Int16:
							dstSlice := make([]int16, inputVal.Len(), inputVal.Cap())
							for i:=0 ; i < inputVal.Len() ; i++{ dstSlice[i] = int16(inputVal.Index(i).Int()) }
							fromGoToCDT(dstSlice, pdata, cdtsIndex, i)
							return

						case reflect.Int32:
							dstSlice := make([]int32, inputVal.Len(), inputVal.Cap())
							for i:=0 ; i < inputVal.Len() ; i++{ dstSlice[i] = int32(inputVal.Index(i).Int()) }
							fromGoToCDT(dstSlice, pdata, cdtsIndex, i)
							return

						case reflect.Int: fallthrough
						case reflect.Int64:
							dstSlice := make([]int64, inputVal.Len(), inputVal.Cap())
							for i:=0 ; i < inputVal.Len() ; i++{ dstSlice[i] = int64(inputVal.Index(i).Int()) }
							fromGoToCDT(dstSlice, pdata, cdtsIndex, i)
							return
		
						case reflect.Uint8: fromGoToCDT(uint8(inputVal.Uint()), pdata, cdtsIndex, i)
							dstSlice := make([]uint8, inputVal.Len(), inputVal.Cap())
							for i:=0 ; i < inputVal.Len() ; i++{ dstSlice[i] = uint8(inputVal.Index(i).Uint()) }
							fromGoToCDT(dstSlice, pdata, cdtsIndex, i)
							return

						case reflect.Uint16: fromGoToCDT(uint16(inputVal.Uint()), pdata, cdtsIndex, i)
							dstSlice := make([]uint16, inputVal.Len(), inputVal.Cap())
							for i:=0 ; i < inputVal.Len() ; i++{ dstSlice[i] = uint16(inputVal.Index(i).Uint()) }
							fromGoToCDT(dstSlice, pdata, cdtsIndex, i)
							return

						case reflect.Uint32:
							dstSlice := make([]uint16, inputVal.Len(), inputVal.Cap())
							for i:=0 ; i < inputVal.Len() ; i++{ dstSlice[i] = uint16(inputVal.Index(i).Uint()) }
							fromGoToCDT(dstSlice, pdata, cdtsIndex, i)
							return

						case reflect.Uint: fallthrough
						case reflect.Uint64:
							dstSlice := make([]uint64, inputVal.Len(), inputVal.Cap())
							for i:=0 ; i < inputVal.Len() ; i++{ dstSlice[i] = uint64(inputVal.Index(i).Uint()) }
							fromGoToCDT(dstSlice, pdata, cdtsIndex, i)
							return
		
						case reflect.Uintptr: 
							dstSlice := make([]uint64, inputVal.Len(), inputVal.Cap())
							for i:=0 ; i < inputVal.Len() ; i++{ dstSlice[i] = uint64(inputVal.Index(i).UnsafeAddr()) }
							fromGoToCDT(dstSlice, pdata, cdtsIndex, i)
							return
		
						case reflect.String:
							dstSlice := make([]string, inputVal.Len(), inputVal.Cap())
							for i:=0 ; i < inputVal.Len() ; i++{ dstSlice[i] = string(inputVal.Index(i).String()) }
							fromGoToCDT(dstSlice, pdata, cdtsIndex, i)
							return
					}

					fallthrough // if no kind matched, treat as handle

				default:
					input_handle := SetObject(input) // if already in table, return existing handle			
					
					out_input := C.metaffi_handle(input_handle)
					out_input_cdt := C.get_cdt(data, index)
					C.set_cdt_type(out_input_cdt, C.metaffi_handle_type)
					out_input_cdt.free_required = 1
					pcdt_out_handle_input := ((*C.struct_cdt_metaffi_handle)(C.convert_union_to_ptr(unsafe.Pointer(&out_input_cdt.cdt_val))))
					pcdt_out_handle_input.val = out_input
			}
	}
}
`

const (
	GuestFunctionXLLRTemplate = `
{{range $mindex, $m := .Modules}}

{{range $findex, $f := $m.Globals}}
{{if $f.Getter}}
// getter for {{$f.Name}}
//export EntryPoint_{{$f.Getter.Name}} 
func EntryPoint_{{GenerateCodeEntryPointSignature "" $f.Getter.Name $f.Getter.Parameters $f.Getter.ReturnValues}}{

	// catch panics and return them as errors
	defer panicHandler(out_err, out_err_len)

	fromGoToCDT({{$f.Name}}, xcall_params, {{GetCDTReturnValueIndex $f.Getter.Parameters $f.Getter.ReturnValues}}, 0)
}
{{end}} {{/* end $f.Get */}}

{{if $f.Setter}}
// setter for {{$f.Name}}
//export EntryPoint_{{$f.Setter.Name}}
func EntryPoint_{{GenerateCodeEntryPointSignature "" $f.Setter.Name $f.Setter.Parameters $f.Setter.ReturnValues}}{

	// catch panics and return them as errors
	defer panicHandler(out_err, out_err_len)

	{{$f.Name}} = fromCDTToGo(xcall_params, {{GetCDTParametersIndex $f.Setter.Parameters}}, 0){{if not $f.IsAny}}.({{ConvertToGoType $f.ArgDefinition $m}}){{end}}
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
	
	// parameters from C to Go
	{{range $index, $elem := $f.Parameters}}	
	{{$elem.Name}}AsInterface := fromCDTToGo(xcall_params, {{GetCDTParametersIndex $f.Parameters}}, {{$index}})
	{{$elem.Name}} := {{if not $elem.IsAny}}{{if $elem.IsTypeAlias}}{{$elem.GetTypeOrAlias}}{{else}}{{ConvertToGoType $elem $m}}{{end}}({{end}}{{$elem.Name}}AsInterface{{if not $elem.IsAny}}.({{ConvertToGoType $elem $m}})){{end}}
	{{end}} {{/* end range params */}}
	
	// call original function
	{{range $index, $elem := $f.ReturnValues}}{{if $index}},{{end}}{{$elem.Name}}{{end}}{{if $f.ReturnValues}} := {{end}}{{$f.Name}}({{range $index, $elem := $f.Parameters}}{{if $index}},{{end}}{{if $elem.IsTypeAlias}}{{$elem.GetTypeOrAlias}}({{$elem.Name}}){{else}}{{$elem.Name}}{{end}}{{end}})
	
	// return values
	{{range $index, $elem := $f.ReturnValues}}
	if err, isError := interface{}({{$elem.Name}}).(error); isError{ // in case of error
		errToOutError(out_err, out_err_len, "Error returned", err)
		return
	} else { // Convert return values from Go to C
		fromGoToCDT({{$elem.Name}}, xcall_params, {{GetCDTReturnValueIndex $f.Parameters $f.ReturnValues}}, {{$index}})
	}	
	{{end}} {{/* end range return vals */}}
}
{{end}} {{/* end range functions */}}


{{range $cindex, $c := $m.Classes}}
// class {{$c.Name}}
{{$className := $c.Name}}

// constructors
{{range $i, $f := $c.Constructors}}
// Call to foreign {{$f.Name}}
//export EntryPoint_{{$c.Name}}_{{$f.Name}}
func EntryPoint_{{GenerateCodeEntryPointSignature $c.Name $f.Name $f.Parameters $f.ReturnValues}}{
	// catch panics and return them as errors
	defer panicHandler(out_err, out_err_len)

	{{ $paramsLength := len $f.Parameters }}{{ $returnLength := len $f.ReturnValues }}
	
	// parameters from C to Go
	{{range $index, $elem := $f.Parameters}}	
	{{$elem.Name}}AsInterface := fromCDTToGo(xcall_params, {{GetCDTParametersIndex $f.Parameters}}, {{$index}})
	{{$elem.Name}} := {{if not $elem.IsAny}}{{if $elem.IsTypeAlias}}{{$elem.GetTypeOrAlias}}{{else}}{{ConvertToGoType $elem $m}}{{end}}({{end}}{{$elem.Name}}AsInterface{{if not $elem.IsAny}}.({{ConvertToGoType $elem $m}})){{end}}
	{{end}} {{/* end range params */}}
	
	// call original function
	{{range $index, $elem := $f.ReturnValues}}{{if $index}},{{end}}{{$elem.Name}}{{end}}{{if $f.ReturnValues}} := {{end}}{{$f.Name}}({{range $index, $elem := $f.Parameters}}{{if $index}},{{end}}{{if $elem.IsTypeAlias}}{{$elem.GetTypeOrAlias}}({{$elem.Name}}){{else}}{{$elem.Name}}{{end}}{{end}})
	
	// return values
	{{range $index, $elem := $f.ReturnValues}}
	if err, isError := interface{}({{$elem.Name}}).(error); isError{ // in case of error
		errToOutError(out_err, out_err_len, "Error returned", err)
		return
	} else { // Convert return values from Go to C
		fromGoToCDT({{$elem.Name}}, xcall_params, {{GetCDTReturnValueIndex $f.Parameters $f.ReturnValues}}, {{$index}})
	}	
	{{end}} {{/* end range return vals */}}
}
{{end}} {{/* end range constructors */}}

{{if $c.Releaser}}// releaser
//export EntryPoint_{{$c.Name}}_{{$c.Releaser.Name}}
func EntryPoint_{{GenerateCodeEntryPointSignature $c.Name $c.Releaser.Name $c.Releaser.Parameters $c.Releaser.ReturnValues}}{

	in_handle_cdt := C.get_cdt(C.get_cdt_element(xcall_params, C.int({{GetCDTParametersIndex $c.Releaser.Parameters}})), 0)

	// first parameter is expected to be the handle
	pcdt_in_handle := ((*C.struct_cdt_metaffi_handle)(C.convert_union_to_ptr(unsafe.Pointer(&in_handle_cdt.cdt_val))))
	var in_handle C.metaffi_handle = pcdt_in_handle.val			

	err := ReleaseObject(Handle(in_handle))
	if err != nil{
		errToOutError(out_err, out_err_len, "Error returned", err)
		return
	}
}
{{end}} {{/* end releaser */}}

// methods
{{range $i, $f := $c.Methods}}
// Call to foreign {{$f.Name}}
//export EntryPoint_{{$c.Name}}_{{$f.Name}}
func EntryPoint_{{GenerateCodeEntryPointSignature $c.Name $f.Name $f.Parameters $f.ReturnValues}}{
	// catch panics and return them as errors
	defer panicHandler(out_err, out_err_len)

	{{ $paramsLength := len $f.Parameters }}{{ $returnLength := len $f.ReturnValues }}
	
	// parameters from C to Go
	{{range $index, $elem := $f.Parameters}}	
	{{$elem.Name}}AsInterface := fromCDTToGo(xcall_params, {{GetCDTParametersIndex $f.Parameters}}, {{$index}})
	{{$elem.Name}} := {{if not $elem.IsAny}}{{if $elem.IsTypeAlias}}{{$elem.GetTypeOrAlias}}{{else}}{{ConvertToGoType $elem $m}}{{end}}({{end}}{{$elem.Name}}AsInterface{{if not $elem.IsAny}}.({{ConvertToGoType $elem $m}})){{end}}
	{{end}} {{/* end range params */}}
	
	// call original function
	{{ $receiver_pointer := index $f.Tags "receiver_pointer"}}
	{{range $index, $elem := $f.ReturnValues}}{{if $index}},{{end}}{{$elem.Name}}{{end}}{{if $f.ReturnValues}} := {{end}}({{(index $f.Parameters 0).Name }}.({{if eq $receiver_pointer "true"}}*{{end}}{{$className}})).{{$f.Name}}({{range $index, $elem := $f.Parameters}}{{if $index}}{{if gt $index 1}},{{end}}{{if $elem.IsTypeAlias}}{{$elem.GetTypeOrAlias}}({{$elem.Name}}){{else}}{{$elem.Name}}{{end}}{{end}}{{end}})

	// return values
	{{range $index, $elem := $f.ReturnValues}}
	if err, isError := interface{}({{$elem.Name}}).(error); isError{ // in case of error
		errToOutError(out_err, out_err_len, "Error returned", err)
		return
	} else { // Convert return values from Go to C
		fromGoToCDT({{$elem.Name}}, xcall_params, {{GetCDTReturnValueIndex $f.Parameters $f.ReturnValues}}, {{$index}})
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

	// get object
	{{ $elem := index $f.Setter.Parameters 0 }}
	objAsInterface := fromCDTToGo(xcall_params, {{GetCDTParametersIndex $f.Getter.Parameters}}, 0)
	obj := {{if not $elem.IsAny}}{{if $elem.IsTypeAlias}}{{$elem.GetTypeOrAlias}}{{else}}{{ConvertToGoType $elem $m}}{{end}}({{end}}objAsInterface{{if not $elem.IsAny}}.({{ConvertToGoType $elem $m}})){{end}}

	{{ $receiver_pointer := index $f.Getter.Tags "receiver_pointer"}}
	{{$f.Name}}_res := obj.({{if eq $receiver_pointer "true"}}*{{end}}{{$className}}).{{$f.Name}}
	
	fromGoToCDT({{$f.Name}}_res, xcall_params, {{GetCDTReturnValueIndex $f.Getter.Parameters $f.Getter.ReturnValues}}, 0)
}
{{end}} {{/* end $f.Getter */}}

{{if $f.Setter}}
// setter for {{$f.Name}}
//export EntryPoint_{{$c.Name}}_{{$f.Setter.Name}}
func EntryPoint_{{GenerateCodeEntryPointSignature $c.Name $f.Setter.Name $f.Setter.Parameters $f.Setter.ReturnValues}}{

	// catch panics and return them as errors
	defer panicHandler(out_err, out_err_len)

	// get object
	{{ $elem := index $f.Setter.Parameters 0 }}
	objAsInterface := fromCDTToGo(xcall_params, {{GetCDTParametersIndex $f.Setter.Parameters}}, 0)
	obj := {{if not $elem.IsAny}}{{if $elem.IsTypeAlias}}{{$elem.GetTypeOrAlias}}{{else}}{{ConvertToGoType $elem $m}}{{end}}({{end}}objAsInterface{{if not $elem.IsAny}}.({{ConvertToGoType $elem $m}})){{end}}

	// get val
	{{ $elem = index $f.Setter.Parameters 1 }}
	valAsInterface := fromCDTToGo(xcall_params, {{GetCDTParametersIndex $f.Setter.Parameters}}, 1)
	val := {{if not $elem.IsAny}}{{if $elem.IsTypeAlias}}{{$elem.GetTypeOrAlias}}{{else}}{{ConvertToGoType $elem $m}}{{end}}({{end}}valAsInterface{{if not $elem.IsAny}}.({{ConvertToGoType $elem $m}})){{end}}

	// get new data
	{{ $receiver_pointer := index $f.Setter.Tags "receiver_pointer"}}
	obj.({{if eq $receiver_pointer "true"}}*{{end}}{{$className}}).{{$f.Name}} = val
	
}
{{end}}{{/* end $f.Setter */}}

{{end}} {{/* end range fields */}}

// end class {{$c.Name}}
{{end}} {{/* end range classes */}}

{{end}} {{/* end range modules */}}

var p *profile.Profile

//export StartProfiler
func StartProfiler(){
	p = profile.Start(profile.ProfilePath(".")).(*profile.Profile)
}

//export EndProfiler
func EndProfiler(){
	p.Stop()
}

`
)
