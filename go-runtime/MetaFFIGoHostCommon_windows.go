package metaffi

/*
#cgo !windows LDFLAGS: -L. -ldl
#cgo LDFLAGS: -Wl,--allow-multiple-definition

#include <stdlib.h>
#include <stdint.h>
#include <include/cdt_structs.h>
#include <include/cdt_capi_loader.h>
#include <include/cdt_capi_loader.c>

metaffi_handle get_null_handle()
{
	return METAFFI_NULL_HANDLE;
}

metaffi_size get_int_item(metaffi_size* array, int index)
{
	return array[index];
}

void* convert_union_to_ptr(void* p)
{
	return p;
}

struct cdts* cast_to_cdts(void* p)
{
	return (cdts*)p;
}

struct cdt* cast_to_cdt(void* p)
{
	return (cdt*)p;
}

struct cdt* get_cdts_index_pcdt(struct cdts* p, int index)
{
	return p[index].pcdt;
}

void set_cdt_type(struct cdt* p, metaffi_type t)
{
	p->type = t;
}

metaffi_type get_cdt_type(struct cdt* p)
{
	return p->type;
}

struct cdt* get_cdt_index(struct cdt* p, int index)
{
	return &p[index];
}

void call_plugin_xcall_no_params_no_ret(void** ppv, char** err, uint64_t* out_err)
{
	void* pvoidxcall = ppv[0];
	void* pctxt = ppv[1];

	(((void(*)(void*,char**,uint64_t*))pvoidxcall)(pctxt, err, out_err));
}

void call_plugin_xcall_no_params_ret(void** ppv, struct cdts* cdts, char** err, uint64_t* out_err)
{
	void* pvoidxcall = ppv[0];
	void* pctxt = ppv[1];

	(((void(*)(void*,void*,char**,uint64_t*))pvoidxcall)(pctxt, cdts, err, out_err));
}

void call_plugin_xcall_params_no_ret(void** ppv, struct cdts* cdts, char** err, uint64_t* out_err)
{
	void* pvoidxcall = ppv[0];
	void* pctxt = ppv[1];

	(((void(*)(void*,void*,char**,uint64_t*))pvoidxcall)(pctxt, cdts, err, out_err));
}

void call_plugin_xcall_params_ret(void** ppv, struct cdts* cdts, char** err, uint64_t* out_err)
{
	void* pvoidxcall = ppv[0];
	void* pctxt = ppv[1];

	(((void(*)(void*,void*,char**,uint64_t*))pvoidxcall)(pctxt, cdts, err, out_err));
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
import (
	"fmt"
	"github.com/MetaFFI/plugin-sdk/compiler/go/IDL"
	"reflect"
	"unsafe"
)

func init() {
	err := C.load_cdt_capi()
	if err != nil {
		panic("Failed to load MetaFFI XLLR functions: " + C.GoString(err))
	}
}

func GetCacheSize() int {
	return int(C.get_cache_size())
}

func createMetaFFITypeInfoArray(paramsTypes []IDL.MetaFFITypeInfo) *C.struct_metaffi_type_info {
	size := len(paramsTypes)
	metaffiArray := C.malloc(C.size_t(size) * C.size_t(unsafe.Sizeof(C.struct_metaffi_type_info{})))

	for i, v := range paramsTypes {
		metaffi := (*C.struct_metaffi_type_info)(unsafe.Pointer(uintptr(metaffiArray) + uintptr(i)*unsafe.Sizeof(C.struct_metaffi_type_info{})))
		metaffi._type = C.metaffi_type(v.Type)

		if v.Alias != "" {
			metaffi.alias = C.CString(v.Alias)
			metaffi.alias_length = C.ulonglong(len(v.Alias))
		} else {
			metaffi.alias = nil
			metaffi.alias_length = 0
		}
	}

	return (*C.struct_metaffi_type_info)(metaffiArray)
}

func freeMetaFFITypeInfoArray(metaffiArray *C.struct_metaffi_type_info, size int) {
	for i := 0; i < size; i++ {
		metaffi := (*C.struct_metaffi_type_info)(unsafe.Pointer(uintptr(unsafe.Pointer(metaffiArray)) + uintptr(i)*unsafe.Sizeof(C.struct_metaffi_type_info{})))
		if metaffi.alias != nil {
			C.free(unsafe.Pointer(metaffi.alias))
		}
	}
	C.free(unsafe.Pointer(metaffiArray))
}

func XLLRLoadFunction(runtimePlugin string, modulePath string, functionPath string, paramsTypes []uint64, retvalsTypes []uint64) (*unsafe.Pointer, error) {

	var params []IDL.MetaFFITypeInfo
	if paramsTypes != nil {
		params = make([]IDL.MetaFFITypeInfo, 0)
		for _, p := range paramsTypes {
			params = append(params, IDL.MetaFFITypeInfo{Type: p})
		}
	}

	var retvals []IDL.MetaFFITypeInfo
	if retvalsTypes != nil {
		retvals = make([]IDL.MetaFFITypeInfo, 0)
		for _, r := range retvalsTypes {
			retvals = append(retvals, IDL.MetaFFITypeInfo{Type: r})
		}
	}

	return XLLRLoadFunctionWithAliases(runtimePlugin, modulePath, functionPath, params, retvals)
}

func XLLRLoadFunctionWithAliases(runtimePlugin string, modulePath string, functionPath string, paramsTypes []IDL.MetaFFITypeInfo, retvalsTypes []IDL.MetaFFITypeInfo) (*unsafe.Pointer, error) {

	pruntimePlugin := C.CString(runtimePlugin)
	defer CFree(unsafe.Pointer(pruntimePlugin))

	pmodulePath := C.CString(modulePath)
	defer CFree(unsafe.Pointer(pmodulePath))

	ppath := C.CString(functionPath)
	defer CFree(unsafe.Pointer(ppath))

	var out_err *C.char
	var out_err_len C.uint32_t
	out_err_len = C.uint32_t(0)

	var pparamTypes *C.struct_metaffi_type_info
	if paramsTypes != nil {
		pparamTypes = createMetaFFITypeInfoArray(paramsTypes)
		defer freeMetaFFITypeInfoArray(pparamTypes, len(paramsTypes))
	}

	pparamTypesLen := (C.uint8_t)(len(paramsTypes))

	var ppretvalsTypes *C.struct_metaffi_type_info
	if retvalsTypes != nil {
		ppretvalsTypes = createMetaFFITypeInfoArray(retvalsTypes)
		defer freeMetaFFITypeInfoArray(ppretvalsTypes, len(retvalsTypes))
	}
	pretvalsTypesLen := (C.uint8_t)(len(retvalsTypes))

	id := C.xllr_load_function(pruntimePlugin, C.uint(len(runtimePlugin)),
		pmodulePath, C.uint(len(modulePath)),
		ppath, C.uint(len(functionPath)),
		pparamTypes, ppretvalsTypes,
		pparamTypesLen, pretvalsTypesLen,
		&out_err, &out_err_len)

	if id == nil {
		return nil, fmt.Errorf("Failed to load foreign entity entrypoint \"%v\": %v", functionPath, string(C.GoBytes(unsafe.Pointer(out_err), C.int(out_err_len))))
	}

	return id, nil
}

func XLLRXCallParamsRet(pff *unsafe.Pointer, parameters unsafe.Pointer) error {

	// TODO: Free error message, in case of returned error
	// 		 The problem is that some plugins return strings that cannot be freed - FIX THIS!

	var out_err *C.char
	var out_err_len C.uint64_t
	out_err_len = C.uint64_t(0)

	C.call_plugin_xcall_params_ret(pff, C.cast_to_cdts(parameters), &out_err, &out_err_len)

	if out_err_len != C.uint64_t(0) {
		return fmt.Errorf("%v", string(C.GoBytes(unsafe.Pointer(out_err), C.int(out_err_len))))
	}

	return nil
}

func XLLRXCallNoParamsRet(pff *unsafe.Pointer, return_values unsafe.Pointer) error {

	var out_err *C.char
	var out_err_len C.uint64_t
	out_err_len = C.uint64_t(0)

	C.call_plugin_xcall_no_params_ret(pff, C.cast_to_cdts(return_values), &out_err, &out_err_len)

	if out_err_len != C.uint64_t(0) {
		return fmt.Errorf("%v", string(C.GoBytes(unsafe.Pointer(out_err), C.int(out_err_len))))
	}

	return nil
}

func XLLRXCallParamsNoRet(pff *unsafe.Pointer, parameters unsafe.Pointer) error {

	var out_err *C.char
	var out_err_len C.uint64_t
	out_err_len = C.uint64_t(0)

	C.call_plugin_xcall_params_no_ret(pff, C.cast_to_cdts(parameters), &out_err, &out_err_len)

	if out_err_len != C.uint64_t(0) {
		return fmt.Errorf("%v", string(C.GoBytes(unsafe.Pointer(out_err), C.int(out_err_len))))
	}

	return nil
}

func XLLRXCallNoParamsNoRet(pff *unsafe.Pointer) error {

	var out_err *C.char
	var out_err_len C.uint64_t
	out_err_len = C.uint64_t(0)

	C.call_plugin_xcall_no_params_no_ret(pff, &out_err, &out_err_len)

	if out_err_len != C.uint64_t(0) {
		return fmt.Errorf("%v", string(C.GoBytes(unsafe.Pointer(out_err), C.int(out_err_len))))
	}

	return nil
}

func XLLRLoadRuntimePlugin(runtimePlugin string) error {

	pruntime_plugin := C.CString(runtimePlugin)
	defer CFree(unsafe.Pointer(pruntime_plugin))

	// load foreign runtime
	var out_err *C.char
	var out_err_len C.uint32_t
	out_err_len = C.uint32_t(0)

	C.xllr_load_runtime_plugin(pruntime_plugin, C.uint(len(runtimePlugin)), &out_err, &out_err_len)

	if out_err_len != C.uint32_t(0) {
		return fmt.Errorf("Failed to load runtime %v: %v", runtimePlugin, string(C.GoBytes(unsafe.Pointer(out_err), C.int(out_err_len))))
	}

	return nil
}

func XLLRFreeRuntimePlugin(runtimePlugin string) error {

	pruntime_plugin := C.CString(runtimePlugin)
	defer CFree(unsafe.Pointer(pruntime_plugin))

	var out_err *C.char
	var out_err_len C.uint32_t
	out_err_len = C.uint32_t(0)

	C.xllr_free_runtime_plugin(pruntime_plugin, C.uint(len(runtimePlugin)), &out_err, &out_err_len)

	if out_err_len != C.uint32_t(0) {
		return fmt.Errorf("Failed to free runtime %v: %v", runtimePlugin, string(C.GoBytes(unsafe.Pointer(out_err), C.int(out_err_len))))
	}

	return nil
}

func CFree(p unsafe.Pointer) {
	C.free(p)
}

func GetPCDTFromCDTSIndex(pcdts unsafe.Pointer, index int) unsafe.Pointer {
	return unsafe.Pointer(C.get_cdts_index_pcdt(C.cast_to_cdts(pcdts), 0))
}

func XLLRAllocCDTSBuffer(params C.metaffi_size, rets C.metaffi_size) (pcdts unsafe.Pointer, parametersCDTS unsafe.Pointer, return_valuesCDTS unsafe.Pointer) {
	res := C.xllr_alloc_cdts_buffer(params, rets)
	pcdts = unsafe.Pointer(res)

	if res != nil {
		parametersCDTS = unsafe.Pointer(C.get_cdts_index_pcdt(res, 0))
		return_valuesCDTS = unsafe.Pointer(C.get_cdts_index_pcdt(res, 1))
	}

	return
}

func GetNullHandle() C.metaffi_handle {
	return C.get_null_handle()
}

func GetIntItem(array *C.metaffi_size, index C.int) C.metaffi_size {
	return C.get_int_item(array, index)
}

func ConvertUnionToPtr(p unsafe.Pointer) unsafe.Pointer {
	return C.convert_union_to_ptr(p)
}

func SetCDTType(p *C.cdt, t C.metaffi_type) {
	C.set_cdt_type(p, t)
}

func GetCDTType(p *C.cdt) C.metaffi_type {
	return C.get_cdt_type(p)
}

func LenToMetaFFISize(i C.longlong) C.metaffi_size {
	return C.len_to_metaffi_size(i)
}

func IntToMetaFFISize(i int) C.metaffi_size {
	return LenToMetaFFISize(C.longlong(i))
}

func LoadCDTCAPI() {
	err := C.load_cdt_capi()
	if err != nil {
		panic("Failed to load MetaFFI XLLR functions: " + C.GoString(err))
	}
}

func FromCDTToGo(pdata unsafe.Pointer, i int) interface{} {

	data := C.cast_to_cdt(pdata)
	var res interface{}
	index := C.int(i)
	in_res_cdt := C.get_cdt_index(data, index)
	res_type := C.get_cdt_type(in_res_cdt)
	switch res_type {

	case 32768: // handle
		pcdt_in_handle_res := ((*C.struct_cdt_metaffi_handle)(C.convert_union_to_ptr(unsafe.Pointer(&in_res_cdt.cdt_val))))
		var in_res C.metaffi_handle = pcdt_in_handle_res.val

		if in_res == C.get_null_handle() {
			return nil
		}

		res = GetObject(Handle(in_res))
		if res == nil { // handle belongs to another language
			res = MetaFFIHandle{
				Val:       Handle(in_res),
				RuntimeID: uint64(pcdt_in_handle_res.runtime_id),
			}
		}

	case 98304: // []Handle
		pcdt_in_handle_res := ((*C.struct_cdt_metaffi_handle_array)(C.convert_union_to_ptr(unsafe.Pointer(&in_res_cdt.cdt_val))))
		var in_res_dimensions_lengths *C.metaffi_size = pcdt_in_handle_res.dimensions_lengths
		// var in_res_dimensions C.metaffi_size = pcdt_in_handle_res.dimensions - TODO: not used until multi-dimensions support!

		length := C.int(C.int(C.get_int_item(in_res_dimensions_lengths, 0)))
		res_typed := make([]interface{}, 0, int(length))
		for i := C.int(0); i < length; i++ {
			val := C.get_metaffi_handle_element(pcdt_in_handle_res.vals, C.int(i))

			val_obj := GetObject(Handle(val.val))
			if val_obj == nil { // handle belongs to
				item := MetaFFIHandle{
					Val:       Handle(val.val),
					RuntimeID: uint64(val.runtime_id),
				}
				res_typed = append(res_typed, item)
			} else {
				res_typed = append(res_typed, val_obj)
			}
		}
		res = res_typed

	case 1: // float64
		pcdt_in_float64_res := ((*C.struct_cdt_metaffi_float64)(C.convert_union_to_ptr(unsafe.Pointer(&in_res_cdt.cdt_val))))
		var in_res C.metaffi_float64 = pcdt_in_float64_res.val

		res = float64(in_res)

	case 2: // float32
		pcdt_in_float32_res := ((*C.struct_cdt_metaffi_float32)(C.convert_union_to_ptr(unsafe.Pointer(&in_res_cdt.cdt_val))))
		var in_res C.metaffi_float32 = pcdt_in_float32_res.val

		res = float32(in_res)

	case 4: // int8
		pcdt_in_int8_res := ((*C.struct_cdt_metaffi_int8)(C.convert_union_to_ptr(unsafe.Pointer(&in_res_cdt.cdt_val))))
		var in_res C.metaffi_int8 = pcdt_in_int8_res.val

		res = int8(in_res)

	case 8: // int16
		pcdt_in_int16_res := ((*C.struct_cdt_metaffi_int16)(C.convert_union_to_ptr(unsafe.Pointer(&in_res_cdt.cdt_val))))
		var in_res C.metaffi_int16 = pcdt_in_int16_res.val

		res = int16(in_res)

	case 16: // int32
		pcdt_in_int32_res := ((*C.struct_cdt_metaffi_int32)(C.convert_union_to_ptr(unsafe.Pointer(&in_res_cdt.cdt_val))))
		var in_res C.metaffi_int32 = pcdt_in_int32_res.val

		res = int32(in_res)

	case 32: // int64
		pcdt_in_int64_res := ((*C.struct_cdt_metaffi_int64)(C.convert_union_to_ptr(unsafe.Pointer(&in_res_cdt.cdt_val))))
		var in_res C.metaffi_int64 = pcdt_in_int64_res.val

		res = int64(in_res)

	case 64: // uint8
		pcdt_in_uint8_res := ((*C.struct_cdt_metaffi_uint8)(C.convert_union_to_ptr(unsafe.Pointer(&in_res_cdt.cdt_val))))
		var in_res C.metaffi_uint8 = pcdt_in_uint8_res.val

		res = uint8(in_res)

	case 128: // uint16
		pcdt_in_uint16_res := ((*C.struct_cdt_metaffi_uint16)(C.convert_union_to_ptr(unsafe.Pointer(&in_res_cdt.cdt_val))))
		var in_res C.metaffi_uint16 = pcdt_in_uint16_res.val

		res = uint16(in_res)

	case 256: // uint32
		pcdt_in_uint32_res := ((*C.struct_cdt_metaffi_uint32)(C.convert_union_to_ptr(unsafe.Pointer(&in_res_cdt.cdt_val))))
		var in_res C.metaffi_uint32 = pcdt_in_uint32_res.val

		res = uint32(in_res)

	case 512: // uint64
		pcdt_in_uint64_res := ((*C.struct_cdt_metaffi_uint64)(C.convert_union_to_ptr(unsafe.Pointer(&in_res_cdt.cdt_val))))
		var in_res C.metaffi_uint64 = pcdt_in_uint64_res.val

		res = uint64(in_res)

	case 65537: // []float64
		pcdt_in_float64_res := ((*C.struct_cdt_metaffi_float64_array)(C.convert_union_to_ptr(unsafe.Pointer(&in_res_cdt.cdt_val))))
		var in_res *C.metaffi_float64 = pcdt_in_float64_res.vals
		var in_res_dimensions_lengths *C.metaffi_size = pcdt_in_float64_res.dimensions_lengths
		// var in_res_dimensions C.metaffi_size = pcdt_in_float64_res.dimensions - TODO: not used until multi-dimensions support!

		res_typed := make([]float64, 0)
		for i := C.int(0); i < C.int(C.int(C.get_int_item(in_res_dimensions_lengths, 0))); i++ {
			val := C.get_metaffi_float64_element(in_res, C.int(i))
			res_typed = append(res_typed, float64(val))
		}
		res = res_typed

	case 65538: // []float32
		pcdt_in_float32_res := ((*C.struct_cdt_metaffi_float32_array)(C.convert_union_to_ptr(unsafe.Pointer(&in_res_cdt.cdt_val))))
		var in_res *C.metaffi_float32 = pcdt_in_float32_res.vals
		var in_res_dimensions_lengths *C.metaffi_size = pcdt_in_float32_res.dimensions_lengths
		// var in_res_dimensions C.metaffi_size = pcdt_in_float32_res.dimensions - TODO: not used until multi-dimensions support!

		res_typed := make([]float32, 0)
		for i := C.int(0); i < C.int(C.int(C.get_int_item(in_res_dimensions_lengths, 0))); i++ {
			val := C.get_metaffi_float32_element(in_res, C.int(i))
			res_typed = append(res_typed, float32(val))
		}
		res = res_typed

	case 65540: // []int8
		pcdt_in_int8_res := ((*C.struct_cdt_metaffi_int8_array)(C.convert_union_to_ptr(unsafe.Pointer(&in_res_cdt.cdt_val))))
		var in_res *C.metaffi_int8 = pcdt_in_int8_res.vals
		var in_res_dimensions_lengths *C.metaffi_size = pcdt_in_int8_res.dimensions_lengths
		// var in_res_dimensions C.metaffi_size = pcdt_in_int8_res.dimensions - TODO: not used until multi-dimensions support!

		res_typed := make([]int8, 0)
		for i := C.int(0); i < C.int(C.int(C.get_int_item(in_res_dimensions_lengths, 0))); i++ {
			val := C.get_metaffi_int8_element(in_res, C.int(i))
			res_typed = append(res_typed, int8(val))
		}
		res = res_typed

	case 65544: // []int16
		pcdt_in_int16_res := ((*C.struct_cdt_metaffi_int16_array)(C.convert_union_to_ptr(unsafe.Pointer(&in_res_cdt.cdt_val))))
		var in_res *C.metaffi_int16 = pcdt_in_int16_res.vals
		var in_res_dimensions_lengths *C.metaffi_size = pcdt_in_int16_res.dimensions_lengths
		// var in_res_dimensions C.metaffi_size = pcdt_in_int16_res.dimensions - TODO: not used until multi-dimensions support!

		res_typed := make([]int16, 0)
		for i := C.int(0); i < C.int(C.int(C.get_int_item(in_res_dimensions_lengths, 0))); i++ {
			val := C.get_metaffi_int16_element(in_res, C.int(i))
			res_typed = append(res_typed, int16(val))
		}
		res = res_typed

	case 65552: // []int32
		pcdt_in_int32_res := ((*C.struct_cdt_metaffi_int32_array)(C.convert_union_to_ptr(unsafe.Pointer(&in_res_cdt.cdt_val))))
		var in_res *C.metaffi_int32 = pcdt_in_int32_res.vals
		var in_res_dimensions_lengths *C.metaffi_size = pcdt_in_int32_res.dimensions_lengths
		// var in_res_dimensions C.metaffi_size = pcdt_in_int32_res.dimensions - TODO: not used until multi-dimensions support!

		res_typed := make([]int32, 0)
		for i := C.int(0); i < C.int(C.int(C.get_int_item(in_res_dimensions_lengths, 0))); i++ {
			val := C.get_metaffi_int32_element(in_res, C.int(i))
			res_typed = append(res_typed, int32(val))
		}
		res = res_typed

	case 65568: // []int64
		pcdt_in_int64_res := ((*C.struct_cdt_metaffi_int64_array)(C.convert_union_to_ptr(unsafe.Pointer(&in_res_cdt.cdt_val))))
		var in_res *C.metaffi_int64 = pcdt_in_int64_res.vals
		var in_res_dimensions_lengths *C.metaffi_size = pcdt_in_int64_res.dimensions_lengths
		// var in_res_dimensions C.metaffi_size = pcdt_in_int64_res.dimensions - TODO: not used until multi-dimensions support!

		res_typed := make([]int64, 0)
		for i := C.int(0); i < C.int(C.int(C.get_int_item(in_res_dimensions_lengths, 0))); i++ {
			val := C.get_metaffi_int64_element(in_res, C.int(i))
			res_typed = append(res_typed, int64(val))
		}
		res = res_typed

	case 65600: // []uint8
		pcdt_in_uint8_res := ((*C.struct_cdt_metaffi_uint8_array)(C.convert_union_to_ptr(unsafe.Pointer(&in_res_cdt.cdt_val))))
		var in_res *C.metaffi_uint8 = pcdt_in_uint8_res.vals
		var in_res_dimensions_lengths *C.metaffi_size = pcdt_in_uint8_res.dimensions_lengths
		// var in_res_dimensions C.metaffi_size = pcdt_in_uint8_res.dimensions - TODO: not used until multi-dimensions support!

		res_typed := make([]uint8, 0)
		for i := C.int(0); i < C.int(C.int(C.get_int_item(in_res_dimensions_lengths, 0))); i++ {
			val := C.get_metaffi_uint8_element(in_res, C.int(i))
			res_typed = append(res_typed, uint8(val))
		}
		res = res_typed

	case 65664: // []uint16
		pcdt_in_uint16_res := ((*C.struct_cdt_metaffi_uint16_array)(C.convert_union_to_ptr(unsafe.Pointer(&in_res_cdt.cdt_val))))
		var in_res *C.metaffi_uint16 = pcdt_in_uint16_res.vals
		var in_res_dimensions_lengths *C.metaffi_size = pcdt_in_uint16_res.dimensions_lengths
		// var in_res_dimensions C.metaffi_size = pcdt_in_uint16_res.dimensions - TODO: not used until multi-dimensions support!

		res_typed := make([]uint16, 0)
		for i := C.int(0); i < C.int(C.int(C.get_int_item(in_res_dimensions_lengths, 0))); i++ {
			val := C.get_metaffi_uint16_element(in_res, C.int(i))
			res_typed = append(res_typed, uint16(val))
		}
		res = res_typed

	case 65792: // []uint32
		pcdt_in_uint32_res := ((*C.struct_cdt_metaffi_uint32_array)(C.convert_union_to_ptr(unsafe.Pointer(&in_res_cdt.cdt_val))))
		var in_res *C.metaffi_uint32 = pcdt_in_uint32_res.vals
		var in_res_dimensions_lengths *C.metaffi_size = pcdt_in_uint32_res.dimensions_lengths
		// var in_res_dimensions C.metaffi_size = pcdt_in_uint32_res.dimensions - TODO: not used until multi-dimensions support!

		res_typed := make([]uint32, 0)
		for i := C.int(0); i < C.int(C.int(C.get_int_item(in_res_dimensions_lengths, 0))); i++ {
			val := C.get_metaffi_uint32_element(in_res, C.int(i))
			res_typed = append(res_typed, uint32(val))
		}
		res = res_typed

	case 66048: // []uint64
		pcdt_in_uint64_res := ((*C.struct_cdt_metaffi_uint64_array)(C.convert_union_to_ptr(unsafe.Pointer(&in_res_cdt.cdt_val))))
		var in_res *C.metaffi_uint64 = pcdt_in_uint64_res.vals
		var in_res_dimensions_lengths *C.metaffi_size = pcdt_in_uint64_res.dimensions_lengths
		// var in_res_dimensions C.metaffi_size = pcdt_in_uint64_res.dimensions - TODO: not used until multi-dimensions support!

		res_typed := make([]uint64, 0)
		for i := C.int(0); i < C.int(C.int(C.get_int_item(in_res_dimensions_lengths, 0))); i++ {
			val := C.get_metaffi_uint64_element(in_res, C.int(i))
			res_typed = append(res_typed, uint64(val))
		}
		res = res_typed

	case 4096: // string8
		in_res_cdt := C.get_cdt_index(data, index)
		pcdt_in_string8_res := ((*C.struct_cdt_metaffi_string8)(C.convert_union_to_ptr(unsafe.Pointer(&in_res_cdt.cdt_val))))
		var in_res_len C.metaffi_size = pcdt_in_string8_res.length
		var in_res C.metaffi_string8 = pcdt_in_string8_res.val

		res = C.GoStringN(in_res, C.int(in_res_len))

	case 69632: // []string8
		in_res_cdt := C.get_cdt_index(data, index)
		pcdt_in_string8_res := ((*C.struct_cdt_metaffi_string8_array)(C.convert_union_to_ptr(unsafe.Pointer(&in_res_cdt.cdt_val))))
		var in_res *C.metaffi_string8 = pcdt_in_string8_res.vals
		var in_res_sizes *C.metaffi_size = pcdt_in_string8_res.vals_sizes
		var in_res_dimensions_lengths *C.metaffi_size = pcdt_in_string8_res.dimensions_lengths
		//var in_res_dimensions C.metaffi_size = pcdt_in_string8_res.dimensions - TODO: not used until multi-dimensions support!

		res_typed := make([]string, 0, int(C.get_int_item(in_res_dimensions_lengths, 0)))
		for i := C.int(0); i < C.int(C.get_int_item(in_res_dimensions_lengths, 0)); i++ {
			var str_size C.metaffi_size
			str := C.get_metaffi_string8_element(in_res, C.int(i), in_res_sizes, &str_size)
			res_typed = append(res_typed, C.GoStringN(str, C.int(str_size)))
		}

		res = res_typed

	case 1024: // bool
		in_res_cdt := C.get_cdt_index(data, index)
		pcdt_in_bool_res := ((*C.struct_cdt_metaffi_bool)(C.convert_union_to_ptr(unsafe.Pointer(&in_res_cdt.cdt_val))))
		var in_res C.metaffi_bool = pcdt_in_bool_res.val

		res = in_res != C.metaffi_bool(0)

	case 66560: // []bool
		in_res_cdt := C.get_cdt_index(data, index)
		pcdt_in_bool_res := ((*C.struct_cdt_metaffi_bool_array)(C.convert_union_to_ptr(unsafe.Pointer(&in_res_cdt.cdt_val))))
		var in_res *C.metaffi_bool = pcdt_in_bool_res.vals
		var in_res_dimensions_lengths *C.metaffi_size = pcdt_in_bool_res.dimensions_lengths
		// var in_res_dimensions C.metaffi_size = pcdt_in_bool_res.dimensions - TODO: not used until multi-dimensions support!

		res_typed := make([]bool, 0)
		for i := C.int(0); i < C.int(C.int(C.get_int_item(in_res_dimensions_lengths, 0))); i++ {
			val := C.get_metaffi_bool_element(in_res, C.int(i))
			var bval bool
			if val != 0 {
				bval = true
			} else {
				bval = false
			}
			res_typed = append(res_typed, bval)
		}

		res = res_typed

	default:
		panic(fmt.Errorf("Return value %v is not of a supported type, but of type: %v", "res", res_type))
	}

	return res
}

func FromGoToCDT(input interface{}, pdata unsafe.Pointer, i int) {

	data := C.cast_to_cdt(pdata)

	index := C.int(i)
	switch input.(type) {

	case Handle:
		panic("Expected MetaFFIHandle, not Handle")

	case MetaFFIHandle:
		out_input := C.metaffi_handle((input.(MetaFFIHandle)).Val)
		out_input_runtime_id := C.metaffi_size((input.(MetaFFIHandle)).RuntimeID)
		out_input_cdt := C.get_cdt_index(data, index)
		C.set_cdt_type(out_input_cdt, C.metaffi_handle_type)
		out_input_cdt.free_required = 1
		pcdt_out_Handle_input := ((*C.struct_cdt_metaffi_handle)(C.convert_union_to_ptr(unsafe.Pointer(&out_input_cdt.cdt_val))))
		pcdt_out_Handle_input.val = out_input
		pcdt_out_Handle_input.runtime_id = out_input_runtime_id

	case float64:
		out_input := C.metaffi_float64(input.(float64))
		out_input_cdt := C.get_cdt_index(data, index)
		C.set_cdt_type(out_input_cdt, C.metaffi_float64_type)
		out_input_cdt.free_required = 1
		pcdt_out_float64_input := ((*C.struct_cdt_metaffi_float64)(C.convert_union_to_ptr(unsafe.Pointer(&out_input_cdt.cdt_val))))
		pcdt_out_float64_input.val = out_input

	case float32:
		out_input := C.metaffi_float32(input.(float32))
		out_input_cdt := C.get_cdt_index(data, index)
		C.set_cdt_type(out_input_cdt, C.metaffi_float32_type)
		out_input_cdt.free_required = 1
		pcdt_out_float32_input := ((*C.struct_cdt_metaffi_float32)(C.convert_union_to_ptr(unsafe.Pointer(&out_input_cdt.cdt_val))))
		pcdt_out_float32_input.val = out_input

	case int8:
		out_input := C.metaffi_int8(input.(int8))
		out_input_cdt := C.get_cdt_index(data, index)
		C.set_cdt_type(out_input_cdt, C.metaffi_int8_type)
		out_input_cdt.free_required = 1
		pcdt_out_int8_input := ((*C.struct_cdt_metaffi_int8)(C.convert_union_to_ptr(unsafe.Pointer(&out_input_cdt.cdt_val))))
		pcdt_out_int8_input.val = out_input

	case int16:
		out_input := C.metaffi_int16(input.(int16))
		out_input_cdt := C.get_cdt_index(data, index)
		C.set_cdt_type(out_input_cdt, C.metaffi_int16_type)
		out_input_cdt.free_required = 1
		pcdt_out_int16_input := ((*C.struct_cdt_metaffi_int16)(C.convert_union_to_ptr(unsafe.Pointer(&out_input_cdt.cdt_val))))
		pcdt_out_int16_input.val = out_input

	case int32:
		out_input := C.metaffi_int32(input.(int32))
		out_input_cdt := C.get_cdt_index(data, index)
		C.set_cdt_type(out_input_cdt, C.metaffi_int32_type)
		out_input_cdt.free_required = 1
		pcdt_out_int32_input := ((*C.struct_cdt_metaffi_int32)(C.convert_union_to_ptr(unsafe.Pointer(&out_input_cdt.cdt_val))))
		pcdt_out_int32_input.val = out_input

	case int64:
		out_input := C.metaffi_int64(input.(int64))
		out_input_cdt := C.get_cdt_index(data, index)
		C.set_cdt_type(out_input_cdt, C.metaffi_int64_type)
		out_input_cdt.free_required = 1
		pcdt_out_int64_input := ((*C.struct_cdt_metaffi_int64)(C.convert_union_to_ptr(unsafe.Pointer(&out_input_cdt.cdt_val))))
		pcdt_out_int64_input.val = out_input

	case uint8:
		out_input := C.metaffi_uint8(input.(uint8))
		out_input_cdt := C.get_cdt_index(data, index)
		C.set_cdt_type(out_input_cdt, C.metaffi_uint8_type)
		out_input_cdt.free_required = 1
		pcdt_out_uint8_input := ((*C.struct_cdt_metaffi_uint8)(C.convert_union_to_ptr(unsafe.Pointer(&out_input_cdt.cdt_val))))
		pcdt_out_uint8_input.val = out_input

	case uint16:
		out_input := C.metaffi_uint16(input.(uint16))
		out_input_cdt := C.get_cdt_index(data, index)
		C.set_cdt_type(out_input_cdt, C.metaffi_uint16_type)
		out_input_cdt.free_required = 1
		pcdt_out_uint16_input := ((*C.struct_cdt_metaffi_uint16)(C.convert_union_to_ptr(unsafe.Pointer(&out_input_cdt.cdt_val))))
		pcdt_out_uint16_input.val = out_input

	case uint32:
		out_input := C.metaffi_uint32(input.(uint32))
		out_input_cdt := C.get_cdt_index(data, index)
		C.set_cdt_type(out_input_cdt, C.metaffi_uint32_type)
		out_input_cdt.free_required = 1
		pcdt_out_uint32_input := ((*C.struct_cdt_metaffi_uint32)(C.convert_union_to_ptr(unsafe.Pointer(&out_input_cdt.cdt_val))))
		pcdt_out_uint32_input.val = out_input

	case uint64:
		out_input := C.metaffi_uint64(input.(uint64))
		out_input_cdt := C.get_cdt_index(data, index)
		C.set_cdt_type(out_input_cdt, C.metaffi_uint64_type)
		out_input_cdt.free_required = 1
		pcdt_out_uint64_input := ((*C.struct_cdt_metaffi_uint64)(C.convert_union_to_ptr(unsafe.Pointer(&out_input_cdt.cdt_val))))
		pcdt_out_uint64_input.val = out_input

	case []Handle:
		panic("Expected []MetaFFIHandle, not []Handle")

	case []MetaFFIHandle:
		out_input_dimensions := C.metaffi_size(1)
		out_input_dimensions_lengths := (*C.metaffi_size)(C.malloc(C.sizeof_metaffi_size))
		*out_input_dimensions_lengths = C.ulonglong(len(input.([]MetaFFIHandle)))

		out_input := (*C.struct_cdt_metaffi_handle)(C.malloc(C.ulonglong(len(input.([]MetaFFIHandle))) * (C.sizeof_struct_cdt_metaffi_handle)))
		for i, val := range input.([]MetaFFIHandle) {
			C.set_metaffi_handle_element(out_input, C.int(i), C.metaffi_handle(val.Val), C.metaffi_size(val.RuntimeID))
		}

		out_input_cdt := C.get_cdt_index(data, index)
		C.set_cdt_type(out_input_cdt, C.metaffi_handle_array_type)
		out_input_cdt.free_required = 1
		pcdt_out_Handle_input := ((*C.struct_cdt_metaffi_handle_array)(C.convert_union_to_ptr(unsafe.Pointer(&out_input_cdt.cdt_val))))
		pcdt_out_Handle_input.vals = out_input
		pcdt_out_Handle_input.dimensions_lengths = out_input_dimensions_lengths
		pcdt_out_Handle_input.dimensions = out_input_dimensions

	case []float64:
		out_input_dimensions := C.metaffi_size(1)
		out_input_dimensions_lengths := (*C.metaffi_size)(C.malloc(C.sizeof_metaffi_size))
		*out_input_dimensions_lengths = C.ulonglong(len(input.([]float64)))

		out_input := (*C.metaffi_float64)(C.malloc(C.ulonglong(len(input.([]float64))) * C.sizeof_metaffi_float64))
		for i, val := range input.([]float64) {
			C.set_metaffi_float64_element(out_input, C.int(i), C.metaffi_float64(val))
		}

		out_input_cdt := C.get_cdt_index(data, index)
		C.set_cdt_type(out_input_cdt, C.metaffi_float64_array_type)
		out_input_cdt.free_required = 1
		pcdt_out_float64_input := ((*C.struct_cdt_metaffi_float64_array)(C.convert_union_to_ptr(unsafe.Pointer(&out_input_cdt.cdt_val))))
		pcdt_out_float64_input.vals = out_input
		pcdt_out_float64_input.dimensions_lengths = out_input_dimensions_lengths
		pcdt_out_float64_input.dimensions = out_input_dimensions

	case []float32:
		out_input_dimensions := C.metaffi_size(1)
		out_input_dimensions_lengths := (*C.metaffi_size)(C.malloc(C.sizeof_metaffi_size))
		*out_input_dimensions_lengths = C.ulonglong(len(input.([]float32)))

		out_input := (*C.metaffi_float32)(C.malloc(C.ulonglong(len(input.([]float32))) * C.sizeof_metaffi_float32))
		for i, val := range input.([]float32) {
			C.set_metaffi_float32_element(out_input, C.int(i), C.metaffi_float32(val))
		}

		out_input_cdt := C.get_cdt_index(data, index)
		C.set_cdt_type(out_input_cdt, C.metaffi_float32_array_type)
		out_input_cdt.free_required = 1
		pcdt_out_float32_input := ((*C.struct_cdt_metaffi_float32_array)(C.convert_union_to_ptr(unsafe.Pointer(&out_input_cdt.cdt_val))))
		pcdt_out_float32_input.vals = out_input
		pcdt_out_float32_input.dimensions_lengths = out_input_dimensions_lengths
		pcdt_out_float32_input.dimensions = out_input_dimensions

	case []int8:
		out_input_dimensions := C.metaffi_size(1)
		out_input_dimensions_lengths := (*C.metaffi_size)(C.malloc(C.sizeof_metaffi_size))
		*out_input_dimensions_lengths = C.ulonglong(len(input.([]int8)))

		out_input := (*C.metaffi_int8)(C.malloc(C.ulonglong(len(input.([]int8))) * C.sizeof_metaffi_int8))
		for i, val := range input.([]int8) {
			C.set_metaffi_int8_element(out_input, C.int(i), C.metaffi_int8(val))
		}

		out_input_cdt := C.get_cdt_index(data, index)
		C.set_cdt_type(out_input_cdt, C.metaffi_int8_array_type)
		out_input_cdt.free_required = 1
		pcdt_out_int8_input := ((*C.struct_cdt_metaffi_int8_array)(C.convert_union_to_ptr(unsafe.Pointer(&out_input_cdt.cdt_val))))
		pcdt_out_int8_input.vals = out_input
		pcdt_out_int8_input.dimensions_lengths = out_input_dimensions_lengths
		pcdt_out_int8_input.dimensions = out_input_dimensions

	case []int16:
		out_input_dimensions := C.metaffi_size(1)
		out_input_dimensions_lengths := (*C.metaffi_size)(C.malloc(C.sizeof_metaffi_size))
		*out_input_dimensions_lengths = C.ulonglong(len(input.([]int16)))

		out_input := (*C.metaffi_int16)(C.malloc(C.ulonglong(len(input.([]int16))) * C.sizeof_metaffi_int16))
		for i, val := range input.([]int16) {
			C.set_metaffi_int16_element(out_input, C.int(i), C.metaffi_int16(val))
		}

		out_input_cdt := C.get_cdt_index(data, index)
		C.set_cdt_type(out_input_cdt, C.metaffi_int16_array_type)
		out_input_cdt.free_required = 1
		pcdt_out_int16_input := ((*C.struct_cdt_metaffi_int16_array)(C.convert_union_to_ptr(unsafe.Pointer(&out_input_cdt.cdt_val))))
		pcdt_out_int16_input.vals = out_input
		pcdt_out_int16_input.dimensions_lengths = out_input_dimensions_lengths
		pcdt_out_int16_input.dimensions = out_input_dimensions

	case []int32:
		out_input_dimensions := C.metaffi_size(1)
		out_input_dimensions_lengths := (*C.metaffi_size)(C.malloc(C.sizeof_metaffi_size))
		*out_input_dimensions_lengths = C.ulonglong(len(input.([]int32)))

		out_input := (*C.metaffi_int32)(C.malloc(C.ulonglong(len(input.([]int32))) * C.sizeof_metaffi_int32))
		for i, val := range input.([]int32) {
			C.set_metaffi_int32_element(out_input, C.int(i), C.metaffi_int32(val))
		}

		out_input_cdt := C.get_cdt_index(data, index)
		C.set_cdt_type(out_input_cdt, C.metaffi_int32_array_type)
		out_input_cdt.free_required = 1
		pcdt_out_int32_input := ((*C.struct_cdt_metaffi_int32_array)(C.convert_union_to_ptr(unsafe.Pointer(&out_input_cdt.cdt_val))))
		pcdt_out_int32_input.vals = out_input
		pcdt_out_int32_input.dimensions_lengths = out_input_dimensions_lengths
		pcdt_out_int32_input.dimensions = out_input_dimensions

	case []int64:
		out_input_dimensions := C.metaffi_size(1)
		out_input_dimensions_lengths := (*C.metaffi_size)(C.malloc(C.sizeof_metaffi_size))
		*out_input_dimensions_lengths = C.ulonglong(len(input.([]int64)))

		out_input := (*C.metaffi_int64)(C.malloc(C.ulonglong(len(input.([]int64))) * C.sizeof_metaffi_int64))
		for i, val := range input.([]int64) {
			C.set_metaffi_int64_element(out_input, C.int(i), C.metaffi_int64(val))
		}

		out_input_cdt := C.get_cdt_index(data, index)
		C.set_cdt_type(out_input_cdt, C.metaffi_int64_array_type)
		out_input_cdt.free_required = 1
		pcdt_out_int64_input := ((*C.struct_cdt_metaffi_int64_array)(C.convert_union_to_ptr(unsafe.Pointer(&out_input_cdt.cdt_val))))
		pcdt_out_int64_input.vals = out_input
		pcdt_out_int64_input.dimensions_lengths = out_input_dimensions_lengths
		pcdt_out_int64_input.dimensions = out_input_dimensions

	case []uint8:
		out_input_dimensions := C.metaffi_size(1)
		out_input_dimensions_lengths := (*C.metaffi_size)(C.malloc(C.sizeof_metaffi_size))
		*out_input_dimensions_lengths = C.ulonglong(len(input.([]uint8)))

		out_input := (*C.metaffi_uint8)(C.malloc(C.ulonglong(len(input.([]uint8))) * C.sizeof_metaffi_uint8))
		for i, val := range input.([]uint8) {
			C.set_metaffi_uint8_element(out_input, C.int(i), C.metaffi_uint8(val))
		}

		out_input_cdt := C.get_cdt_index(data, index)
		C.set_cdt_type(out_input_cdt, C.metaffi_uint8_array_type)
		out_input_cdt.free_required = 1
		pcdt_out_uint8_input := ((*C.struct_cdt_metaffi_uint8_array)(C.convert_union_to_ptr(unsafe.Pointer(&out_input_cdt.cdt_val))))
		pcdt_out_uint8_input.vals = out_input
		pcdt_out_uint8_input.dimensions_lengths = out_input_dimensions_lengths
		pcdt_out_uint8_input.dimensions = out_input_dimensions

	case []uint16:
		out_input_dimensions := C.metaffi_size(1)
		out_input_dimensions_lengths := (*C.metaffi_size)(C.malloc(C.sizeof_metaffi_size))
		*out_input_dimensions_lengths = C.ulonglong(len(input.([]uint16)))

		out_input := (*C.metaffi_uint16)(C.malloc(C.ulonglong(len(input.([]uint16))) * C.sizeof_metaffi_uint16))
		for i, val := range input.([]uint16) {
			C.set_metaffi_uint16_element(out_input, C.int(i), C.metaffi_uint16(val))
		}

		out_input_cdt := C.get_cdt_index(data, index)
		C.set_cdt_type(out_input_cdt, C.metaffi_uint16_array_type)
		out_input_cdt.free_required = 1
		pcdt_out_uint16_input := ((*C.struct_cdt_metaffi_uint16_array)(C.convert_union_to_ptr(unsafe.Pointer(&out_input_cdt.cdt_val))))
		pcdt_out_uint16_input.vals = out_input
		pcdt_out_uint16_input.dimensions_lengths = out_input_dimensions_lengths
		pcdt_out_uint16_input.dimensions = out_input_dimensions

	case []uint32:
		out_input_dimensions := C.metaffi_size(1)
		out_input_dimensions_lengths := (*C.metaffi_size)(C.malloc(C.sizeof_metaffi_size))
		*out_input_dimensions_lengths = C.ulonglong(len(input.([]uint32)))

		out_input := (*C.metaffi_uint32)(C.malloc(C.ulonglong(len(input.([]uint32))) * C.sizeof_metaffi_uint32))
		for i, val := range input.([]uint32) {
			C.set_metaffi_uint32_element(out_input, C.int(i), C.metaffi_uint32(val))
		}

		out_input_cdt := C.get_cdt_index(data, index)
		C.set_cdt_type(out_input_cdt, C.metaffi_uint32_array_type)
		out_input_cdt.free_required = 1
		pcdt_out_uint32_input := ((*C.struct_cdt_metaffi_uint32_array)(C.convert_union_to_ptr(unsafe.Pointer(&out_input_cdt.cdt_val))))
		pcdt_out_uint32_input.vals = out_input
		pcdt_out_uint32_input.dimensions_lengths = out_input_dimensions_lengths
		pcdt_out_uint32_input.dimensions = out_input_dimensions

	case []uint64:
		out_input_dimensions := C.metaffi_size(1)
		out_input_dimensions_lengths := (*C.metaffi_size)(C.malloc(C.sizeof_metaffi_size))
		*out_input_dimensions_lengths = C.ulonglong(len(input.([]uint64)))

		out_input := (*C.metaffi_uint64)(C.malloc(C.ulonglong(len(input.([]uint64))) * C.sizeof_metaffi_uint64))
		for i, val := range input.([]uint64) {
			C.set_metaffi_uint64_element(out_input, C.int(i), C.metaffi_uint64(val))
		}

		out_input_cdt := C.get_cdt_index(data, index)
		C.set_cdt_type(out_input_cdt, C.metaffi_uint64_array_type)
		out_input_cdt.free_required = 1
		pcdt_out_uint64_input := ((*C.struct_cdt_metaffi_uint64_array)(C.convert_union_to_ptr(unsafe.Pointer(&out_input_cdt.cdt_val))))
		pcdt_out_uint64_input.vals = out_input
		pcdt_out_uint64_input.dimensions_lengths = out_input_dimensions_lengths
		pcdt_out_uint64_input.dimensions = out_input_dimensions

	case int:
		out_input := C.metaffi_int64(int64(input.(int)))
		out_input_cdt := C.get_cdt_index(data, index)
		C.set_cdt_type(out_input_cdt, C.metaffi_int64_type)
		out_input_cdt.free_required = 1
		pcdt_out_int64_input := ((*C.struct_cdt_metaffi_int64)(C.convert_union_to_ptr(unsafe.Pointer(&out_input_cdt.cdt_val))))
		pcdt_out_int64_input.val = out_input

	case []int:
		out_input_dimensions := C.metaffi_size(1)
		out_input_dimensions_lengths := (*C.metaffi_size)(C.malloc(C.sizeof_metaffi_size))
		*out_input_dimensions_lengths = C.ulonglong(len(input.([]int)))

		out_input := (*C.metaffi_int64)(C.malloc(C.ulonglong(len(input.([]int))) * C.sizeof_metaffi_int64))
		for i, val := range input.([]int) {
			C.set_metaffi_int64_element(out_input, C.int(i), C.metaffi_int64(val))
		}

		out_input_cdt := C.get_cdt_index(data, index)
		C.set_cdt_type(out_input_cdt, C.metaffi_int64_array_type)
		out_input_cdt.free_required = 1
		pcdt_out_int64_input := ((*C.struct_cdt_metaffi_int64_array)(C.convert_union_to_ptr(unsafe.Pointer(&out_input_cdt.cdt_val))))
		pcdt_out_int64_input.vals = out_input
		pcdt_out_int64_input.dimensions_lengths = out_input_dimensions_lengths
		pcdt_out_int64_input.dimensions = out_input_dimensions

	case bool:
		var out_input C.metaffi_bool
		if input.(bool) {
			out_input = C.metaffi_bool(1)
		} else {
			out_input = C.metaffi_bool(0)
		}
		out_input_cdt := C.get_cdt_index(data, index)
		C.set_cdt_type(out_input_cdt, C.metaffi_bool_type)
		out_input_cdt.free_required = 1
		pcdt_out_bool_input := ((*C.struct_cdt_metaffi_bool)(C.convert_union_to_ptr(unsafe.Pointer(&out_input_cdt.cdt_val))))
		pcdt_out_bool_input.val = out_input

	case string:
		out_input_len := C.metaffi_size(C.ulonglong(len(input.(string))))
		out_input := C.CString(input.(string))
		out_input_cdt := C.get_cdt_index(data, index)
		C.set_cdt_type(out_input_cdt, C.metaffi_string8_type)
		out_input_cdt.free_required = 1
		pcdt_out_string8_input := ((*C.struct_cdt_metaffi_string8)(C.convert_union_to_ptr(unsafe.Pointer(&out_input_cdt.cdt_val))))
		pcdt_out_string8_input.val = out_input
		pcdt_out_string8_input.length = out_input_len

	case []bool:
		out_input_dimensions := C.metaffi_size(1)
		out_input_dimensions_lengths := (*C.metaffi_size)(C.malloc(C.sizeof_metaffi_size))
		*out_input_dimensions_lengths = C.metaffi_size(len(input.([]bool)))

		out_input := (*C.metaffi_bool)(C.malloc(C.metaffi_size(len(input.([]bool))) * C.sizeof_metaffi_bool))
		for i, val := range input.([]bool) {
			var bval C.metaffi_bool
			if val {
				bval = C.metaffi_bool(1)
			} else {
				bval = C.metaffi_bool(0)
			}
			C.set_metaffi_bool_element(out_input, C.int(i), C.metaffi_bool(bval))
		}

		out_input_cdt := C.get_cdt_index(data, index)
		C.set_cdt_type(out_input_cdt, C.metaffi_bool_array_type)
		out_input_cdt.free_required = 1
		pcdt_out_bool_input := ((*C.struct_cdt_metaffi_bool_array)(C.convert_union_to_ptr(unsafe.Pointer(&out_input_cdt.cdt_val))))
		pcdt_out_bool_input.vals = out_input
		pcdt_out_bool_input.dimensions_lengths = out_input_dimensions_lengths
		pcdt_out_bool_input.dimensions = out_input_dimensions

	case []string:
		out_input := (*C.metaffi_string8)(C.malloc(C.ulonglong(len(input.([]string))) * C.sizeof_metaffi_string8))
		out_input_sizes := (*C.metaffi_size)(C.malloc(C.ulonglong(len(input.([]string))) * C.sizeof_metaffi_size))
		out_input_dimensions := C.metaffi_size(1)
		out_input_dimensions_lengths := (*C.metaffi_size)(C.malloc(C.sizeof_metaffi_size * (out_input_dimensions)))
		*out_input_dimensions_lengths = C.metaffi_size(len(input.([]string)))

		for i, val := range input.([]string) {
			C.set_metaffi_string8_element(out_input, out_input_sizes, C.int(i), C.metaffi_string8(C.CString(val)), C.metaffi_size(len(val)))
		}

		out_input_cdt := C.get_cdt_index(data, index)
		C.set_cdt_type(out_input_cdt, C.metaffi_string8_array_type)
		out_input_cdt.free_required = 1
		pcdt_out_string8_input := ((*C.struct_cdt_metaffi_string8_array)(C.convert_union_to_ptr(unsafe.Pointer(&out_input_cdt.cdt_val))))
		pcdt_out_string8_input.vals = out_input
		pcdt_out_string8_input.vals_sizes = out_input_sizes
		pcdt_out_string8_input.dimensions_lengths = out_input_dimensions_lengths
		pcdt_out_string8_input.dimensions = out_input_dimensions

	default:

		if input == nil { // return handle "0"
			out_input := C.metaffi_handle(uintptr(0))
			out_input_cdt := C.get_cdt_index(data, index)
			C.set_cdt_type(out_input_cdt, C.metaffi_handle_type)
			out_input_cdt.free_required = 0
			pcdt_out_handle_input := ((*C.struct_cdt_metaffi_handle)(C.convert_union_to_ptr(unsafe.Pointer(&out_input_cdt.cdt_val))))
			pcdt_out_handle_input.val = out_input
			return
		}

		// check if the object is type of a primitive
		inputVal := reflect.ValueOf(input)
		inputType := reflect.TypeOf(input)
		switch inputType.Kind() {
		case reflect.Bool:
			FromGoToCDT(bool(inputVal.Bool()), pdata, i)
			return

		case reflect.Float32:
			FromGoToCDT(float32(inputVal.Float()), pdata, i)
			return
		case reflect.Float64:
			FromGoToCDT(float64(inputVal.Float()), pdata, i)
			return

		case reflect.Int8:
			FromGoToCDT(int8(inputVal.Int()), pdata, i)
			return
		case reflect.Int16:
			FromGoToCDT(int16(inputVal.Int()), pdata, i)
			return
		case reflect.Int32:
			FromGoToCDT(int32(inputVal.Int()), pdata, i)
			return
		case reflect.Int:
			fallthrough
		case reflect.Int64:
			FromGoToCDT(int64(inputVal.Int()), pdata, i)
			return

		case reflect.Uint8:
			FromGoToCDT(uint8(inputVal.Uint()), pdata, i)
			return
		case reflect.Uint16:
			FromGoToCDT(uint16(inputVal.Uint()), pdata, i)
			return
		case reflect.Uint32:
			FromGoToCDT(uint32(inputVal.Uint()), pdata, i)
			return
		case reflect.Uint:
			fallthrough
		case reflect.Uint64:
			FromGoToCDT(uint64(inputVal.Uint()), pdata, i)
			return

		case reflect.Uintptr:
			FromGoToCDT(uint64(inputVal.UnsafeAddr()), pdata, i)
			return

		case reflect.String:
			FromGoToCDT(string(inputVal.String()), pdata, i)
			return

		case reflect.Slice:
			switch inputType.Elem().Kind() {
			case reflect.Float32:
				dstSlice := make([]float32, inputVal.Len(), inputVal.Cap())
				for i := 0; i < inputVal.Len(); i++ {
					dstSlice[i] = float32(inputVal.Index(i).Float())
				}
				FromGoToCDT(dstSlice, pdata, i)
				return

			case reflect.Float64:
				dstSlice := make([]float64, inputVal.Len(), inputVal.Cap())
				for i := 0; i < inputVal.Len(); i++ {
					dstSlice[i] = float64(inputVal.Index(i).Float())
				}
				FromGoToCDT(dstSlice, pdata, i)
				return

			case reflect.Bool:
				dstSlice := make([]bool, inputVal.Len(), inputVal.Cap())
				for i := 0; i < inputVal.Len(); i++ {
					dstSlice[i] = inputVal.Index(i).Bool()
				}
				FromGoToCDT(dstSlice, pdata, i)
				return

			case reflect.Int8:
				dstSlice := make([]int8, inputVal.Len(), inputVal.Cap())
				for i := 0; i < inputVal.Len(); i++ {
					dstSlice[i] = int8(inputVal.Index(i).Int())
				}
				FromGoToCDT(dstSlice, pdata, i)
				return

			case reflect.Int16:
				dstSlice := make([]int16, inputVal.Len(), inputVal.Cap())
				for i := 0; i < inputVal.Len(); i++ {
					dstSlice[i] = int16(inputVal.Index(i).Int())
				}
				FromGoToCDT(dstSlice, pdata, i)
				return

			case reflect.Int32:
				dstSlice := make([]int32, inputVal.Len(), inputVal.Cap())
				for i := 0; i < inputVal.Len(); i++ {
					dstSlice[i] = int32(inputVal.Index(i).Int())
				}
				FromGoToCDT(dstSlice, pdata, i)
				return

			case reflect.Int:
				fallthrough
			case reflect.Int64:
				dstSlice := make([]int64, inputVal.Len(), inputVal.Cap())
				for i := 0; i < inputVal.Len(); i++ {
					dstSlice[i] = int64(inputVal.Index(i).Int())
				}
				FromGoToCDT(dstSlice, pdata, i)
				return

			case reflect.Uint8:
				FromGoToCDT(uint8(inputVal.Uint()), pdata, i)
				dstSlice := make([]uint8, inputVal.Len(), inputVal.Cap())
				for i := 0; i < inputVal.Len(); i++ {
					dstSlice[i] = uint8(inputVal.Index(i).Uint())
				}
				FromGoToCDT(dstSlice, pdata, i)
				return

			case reflect.Uint16:
				FromGoToCDT(uint16(inputVal.Uint()), pdata, i)
				dstSlice := make([]uint16, inputVal.Len(), inputVal.Cap())
				for i := 0; i < inputVal.Len(); i++ {
					dstSlice[i] = uint16(inputVal.Index(i).Uint())
				}
				FromGoToCDT(dstSlice, pdata, i)
				return

			case reflect.Uint32:
				dstSlice := make([]uint16, inputVal.Len(), inputVal.Cap())
				for i := 0; i < inputVal.Len(); i++ {
					dstSlice[i] = uint16(inputVal.Index(i).Uint())
				}
				FromGoToCDT(dstSlice, pdata, i)
				return

			case reflect.Uint:
				fallthrough
			case reflect.Uint64:
				dstSlice := make([]uint64, inputVal.Len(), inputVal.Cap())
				for i := 0; i < inputVal.Len(); i++ {
					dstSlice[i] = uint64(inputVal.Index(i).Uint())
				}
				FromGoToCDT(dstSlice, pdata, i)
				return

			case reflect.Uintptr:
				dstSlice := make([]uint64, inputVal.Len(), inputVal.Cap())
				for i := 0; i < inputVal.Len(); i++ {
					dstSlice[i] = uint64(inputVal.Index(i).UnsafeAddr())
				}
				FromGoToCDT(dstSlice, pdata, i)
				return

			case reflect.String:
				dstSlice := make([]string, inputVal.Len(), inputVal.Cap())
				for i := 0; i < inputVal.Len(); i++ {
					dstSlice[i] = string(inputVal.Index(i).String())
				}
				FromGoToCDT(dstSlice, pdata, i)
				return
			}

			fallthrough // if no kind matched, treat as handle

		default:
			input_handle := SetObject(input) // if already in table, return existing handle

			out_input := C.metaffi_handle(input_handle)
			out_input_cdt := C.get_cdt_index(data, index)
			C.set_cdt_type(out_input_cdt, C.metaffi_handle_type)
			out_input_cdt.free_required = 1
			pcdt_out_handle_input := ((*C.struct_cdt_metaffi_handle)(C.convert_union_to_ptr(unsafe.Pointer(&out_input_cdt.cdt_val))))
			pcdt_out_handle_input.val = out_input
		}
	}
}
