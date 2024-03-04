package metaffi

/*
#cgo !windows LDFLAGS: -L. -ldl
#cgo LDFLAGS: -Wl,--allow-multiple-definition

#include <stdlib.h>
#include <stdint.h>
#include <include/cdt_structs.h>
#include <include/cdt_capi_loader.h>
#include <include/cdt_capi_loader.c>

struct cdt_metaffi_handle get_null_handle()
{
	struct cdt_metaffi_handle res;
	res.val = NULL;
	res.runtime_id = 0;
	res.release = NULL;
	return res;
}

metaffi_size get_int_item(metaffi_size* array, int index)
{
	return array[index];
}

void set_int_item(metaffi_size* array, int index, metaffi_size value)
{
	array[index] = value;
}

void* convert_union_to_ptr(void* p)
{
	return p;
}

struct cdts* cast_to_cdts(void* p)
{
	return (cdts*)p;
}

struct cdt* get_cdt_index(struct cdt* p, int index)
{
	return &p[index];
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

void* get_index(void* p, int index)
{
	return p + index;
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

int8_t get_int8_item(metaffi_int8* array, int index)
{
	return array[index];
}

metaffi_size get_metaffi_size_item(metaffi_size* array, int index)
{
	return array[index];
}

void set_metaffi_size_item(metaffi_size* array, int index, metaffi_size value)
{
	array[index] = value;
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
	"os"
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

func GetNullHandle() IDL.MetaFFITypeInfo {
	return IDL.MetaFFITypeInfo{IDL.NULL, "", IDL.METAFFI_TYPE_NULL, 0}
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

func createMultiDimSliceOfHandles(pcdt_arr *C.struct_cdt_metaffi_handle_array, dims []int, index int) []interface{} {

	if index == len(dims)-1 {
		// Base case: create a 1D slice
		slice := make([]interface{}, dims[index])
		for i := range slice {
			metaffi_handle_instance := (*C.struct_cdt_metaffi_handle)(unsafe.Pointer(uintptr(unsafe.Pointer(pcdt_arr.vals)) + uintptr(i)*unsafe.Sizeof(*pcdt_arr.vals)))
			var mhandle interface{}
			if metaffi_handle_instance.val == C.get_null_handle().val {
				mhandle = nil
			} else if metaffi_handle_instance.runtime_id == GO_RUNTIME_ID {
				mhandle = GetObject(Handle(metaffi_handle_instance.val))
			} else {
				mhandle = MetaFFIHandle{
					Val:       Handle(metaffi_handle_instance.val),
					RuntimeID: uint64(metaffi_handle_instance.runtime_id),
				}
			}
			slice[i] = mhandle
		}
		return slice
	}

	// Recursive case: create a multi-dimensional slice
	slice := make([]interface{}, dims[index])
	for i := range slice {
		slice[i] = createMultiDimSliceOfHandles(pcdt_arr, dims, index+1)
	}
	return slice
}

func convertToMultiDimSliceOfHandles(cdt_arr *C.struct_cdt_metaffi_handle_array) []interface{} {
	dims := make([]int, cdt_arr.dimensions)
	for i := range dims {
		dims[i] = int(C.get_int_item(cdt_arr.dimensions_lengths, C.int(i)))
	}
	return createMultiDimSliceOfHandles(cdt_arr, dims, 0)
}

type GoNumber interface {
	int | uint | int8 | uint8 | int16 | uint16 | int32 | uint32 | int64 | uint64 | float32 | float64
}

type CNumber interface {
	C.int8_t | C.uint8_t | C.int16_t | C.uint16_t | C.int32_t | C.uint32_t | C.int64_t | C.uint64_t | C.float | C.double
}

func createMultiDimSlice[ctype_t CNumber, gotype_t GoNumber](data_arr *ctype_t, dims []int, index int, sizeofPointer uintptr, sizeofElement uintptr) []interface{} {
	if index == len(dims)-1 {
		// Base case: create a 1D slice
		slice := make([]interface{}, dims[index])
		fmt.Printf("size of slice: %v\n", len(slice))
		for i := range slice {
			fmt.Printf("settings element: %v\n", i)
			elem := (*ctype_t)(unsafe.Pointer(uintptr(unsafe.Pointer(data_arr)) + uintptr(i)*sizeofElement))
			fmt.Printf("pointer to elem: %v\n", elem)
			slice[i] = *(*gotype_t)(unsafe.Pointer(elem))
		}
		return slice
	}

	// Recursive case: create a multi-dimensional slice
	slice := make([]interface{}, dims[index])
	for i := range slice {
		// *(int8_t**)(data_arr + i*sizeof(pointer))
		dereferenced := (*ctype_t)(*((**ctype_t)(unsafe.Pointer(uintptr(unsafe.Pointer(data_arr)) + uintptr(i)*sizeofPointer))))
		slice[i] = createMultiDimSlice[ctype_t, gotype_t](dereferenced, dims, index+1, sizeofPointer, sizeofElement)
	}
	return slice
}

func convertToMultiDimSliceOfNumbers[ctype_t CNumber, gotype_t GoNumber](dimensions C.metaffi_size, dimensions_length *C.metaffi_size, vals *ctype_t, sizeofElement uintptr) []interface{} {
	dims := make([]int, dimensions)
	for i := range dims {
		dims[i] = int(C.get_int_item(dimensions_length, C.int(i)))
	}
	x := 0
	return createMultiDimSlice[ctype_t, gotype_t](vals, dims, 0, unsafe.Sizeof(&x), sizeofElement)
}

// For string
func createMultiDimSliceOfString(strings *C.metaffi_string8, strings_lengths *C.metaffi_size, sizeofPointer uintptr, dims []int, index int) []interface{} {
	if index == len(dims)-1 {
		// Base case: create a 1D slice
		slice := make([]interface{}, dims[index])
		for i := range slice {
			str := *(*C.metaffi_string8)(unsafe.Pointer(uintptr(unsafe.Pointer(strings)) + uintptr(i)*sizeofPointer))
			strLen := *(*C.metaffi_size)(unsafe.Pointer(uintptr(unsafe.Pointer(strings_lengths)) + uintptr(i)*sizeofPointer))
			slice[i] = C.GoStringN(str, C.int(strLen))
		}
		return slice
	}

	// Recursive case: create a multi-dimensional slice
	slice := make([]interface{}, dims[index])
	for i := range slice {
		dereferenced_strings := (*C.metaffi_string8)(*((**C.metaffi_string8)(unsafe.Pointer(uintptr(unsafe.Pointer(strings)) + uintptr(i)*sizeofPointer))))
		dereferenced_lengths := (*C.metaffi_size)(*((**C.metaffi_size)(unsafe.Pointer(uintptr(unsafe.Pointer(strings_lengths)) + uintptr(i)*sizeofPointer))))
		slice[i] = createMultiDimSliceOfString(dereferenced_strings, dereferenced_lengths, sizeofPointer, dims, index+1)
	}
	return slice
}

func convertToMultiDimSliceOfStrings(cdt_arr *C.struct_cdt_metaffi_string8_array) []interface{} {
	dims := make([]int, cdt_arr.dimensions)
	for i := range dims {
		dims[i] = int(C.get_int_item(cdt_arr.dimensions_lengths, C.int(i)))
	}

	x := 0
	return createMultiDimSliceOfString(cdt_arr.vals, cdt_arr.vals_sizes, unsafe.Sizeof(&x), dims, 0)
}

// For bool
func createMultiDimSliceOfBool(cdt_arr *C.struct_cdt_metaffi_bool_array, dims []int, index int) []interface{} {
	if index == len(dims)-1 {
		// Base case: create a 1D slice
		slice := make([]interface{}, dims[index])
		for i := range slice {
			bool_instance := (*C.struct_cdt_metaffi_bool)(unsafe.Pointer(uintptr(unsafe.Pointer(cdt_arr.vals)) + uintptr(i)*unsafe.Sizeof(*cdt_arr.vals)))
			slice[i] = bool(bool_instance.val != 0)
		}
		return slice
	}

	// Recursive case: create a multi-dimensional slice
	slice := make([]interface{}, dims[index])
	for i := range slice {
		slice[i] = createMultiDimSliceOfBool(cdt_arr, dims, index+1)
	}
	return slice
}

func convertToMultiDimSliceOfBools(cdt_arr *C.struct_cdt_metaffi_bool_array) []interface{} {
	dims := make([]int, cdt_arr.dimensions)
	for i := range dims {
		dims[i] = int(C.get_int_item(cdt_arr.dimensions_lengths, C.int(i)))
	}
	return createMultiDimSliceOfBool(cdt_arr, dims, 0)
}

func FromCDTToGo(pdata unsafe.Pointer, i int) interface{} {
	fmt.Fprintf(os.Stderr, "+++++++++ In FromCDTToGo +++++++++\n")
	data := C.cast_to_cdt(pdata)
	var res interface{}
	index := C.int(i)
	pcdt := C.get_cdt_index(data, index)
	res_type := C.get_cdt_type(pcdt)

	switch uint64(res_type) {

	case IDL.METAFFI_TYPE_HANDLE: // handle
		pcdt_handle_res := ((*C.struct_cdt_metaffi_handle)(C.convert_union_to_ptr(unsafe.Pointer(&pcdt.cdt_val))))
		var h C.metaffi_handle = pcdt_handle_res.val

		if h == C.get_null_handle().val {
			res = nil
		}

		if pcdt_handle_res.runtime_id == GO_RUNTIME_ID {
			res = GetObject(Handle(h))
		} else {
			res = MetaFFIHandle{
				Val:       Handle(h),
				RuntimeID: uint64(pcdt_handle_res.runtime_id),
			}
		}

	case IDL.METAFFI_TYPE_HANDLE_ARRAY:
		cdt_handles_arr := ((*C.struct_cdt_metaffi_handle_array)(C.convert_union_to_ptr(unsafe.Pointer(&pcdt.cdt_val))))
		res = convertToMultiDimSliceOfHandles(cdt_handles_arr)

	case IDL.METAFFI_TYPE_FLOAT64:
		pcdt_float64 := ((*C.struct_cdt_metaffi_float64)(C.convert_union_to_ptr(unsafe.Pointer(&pcdt.cdt_val))))
		res = float64(pcdt_float64.val)

	case IDL.METAFFI_TYPE_FLOAT64_ARRAY:
		pcdt_float64_arr := ((*C.struct_cdt_metaffi_float64_array)(C.convert_union_to_ptr(unsafe.Pointer(&pcdt.cdt_val))))
		res = convertToMultiDimSliceOfNumbers[C.metaffi_float64, float64](pcdt_float64_arr.dimensions, pcdt_float64_arr.dimensions_lengths, pcdt_float64_arr.vals, C.sizeof_metaffi_float64)

	case IDL.METAFFI_TYPE_FLOAT32:
		pcdt_float32 := ((*C.struct_cdt_metaffi_float32)(C.convert_union_to_ptr(unsafe.Pointer(&pcdt.cdt_val))))
		res = float32(pcdt_float32.val)

	case IDL.METAFFI_TYPE_FLOAT32_ARRAY:
		pcdt_float32_arr := ((*C.struct_cdt_metaffi_float32_array)(C.convert_union_to_ptr(unsafe.Pointer(&pcdt.cdt_val))))
		res = convertToMultiDimSliceOfNumbers[C.metaffi_float32, float32](pcdt_float32_arr.dimensions, pcdt_float32_arr.dimensions_lengths, pcdt_float32_arr.vals, C.sizeof_metaffi_float32)

	case IDL.METAFFI_TYPE_INT8:
		pcdt_int8 := ((*C.struct_cdt_metaffi_int8)(C.convert_union_to_ptr(unsafe.Pointer(&pcdt.cdt_val))))
		res = int8(pcdt_int8.val)

	case IDL.METAFFI_TYPE_INT8_ARRAY:
		pcdt_int8_arr := ((*C.struct_cdt_metaffi_int8_array)(C.convert_union_to_ptr(unsafe.Pointer(&pcdt.cdt_val))))
		res = convertToMultiDimSliceOfNumbers[C.metaffi_int8, int8](pcdt_int8_arr.dimensions, pcdt_int8_arr.dimensions_lengths, pcdt_int8_arr.vals, C.sizeof_int8_t)

		// For uint8
	case IDL.METAFFI_TYPE_UINT8:
		pcdt_uint8 := ((*C.struct_cdt_metaffi_uint8)(C.convert_union_to_ptr(unsafe.Pointer(&pcdt.cdt_val))))
		res = uint8(pcdt_uint8.val)

	case IDL.METAFFI_TYPE_UINT8_ARRAY:
		pcdt_uint8_arr := ((*C.struct_cdt_metaffi_uint8_array)(C.convert_union_to_ptr(unsafe.Pointer(&pcdt.cdt_val))))
		res = convertToMultiDimSliceOfNumbers[C.metaffi_uint8, uint8](pcdt_uint8_arr.dimensions, pcdt_uint8_arr.dimensions_lengths, pcdt_uint8_arr.vals, C.sizeof_uint8_t)

		// For int16
	case IDL.METAFFI_TYPE_INT16:
		pcdt_int16 := ((*C.struct_cdt_metaffi_int16)(C.convert_union_to_ptr(unsafe.Pointer(&pcdt.cdt_val))))
		res = int16(pcdt_int16.val)

	case IDL.METAFFI_TYPE_INT16_ARRAY:
		pcdt_int16_arr := ((*C.struct_cdt_metaffi_int16_array)(C.convert_union_to_ptr(unsafe.Pointer(&pcdt.cdt_val))))
		res = convertToMultiDimSliceOfNumbers[C.metaffi_int16, int16](pcdt_int16_arr.dimensions, pcdt_int16_arr.dimensions_lengths, pcdt_int16_arr.vals, C.sizeof_int16_t)

		// For uint16
	case IDL.METAFFI_TYPE_UINT16:
		pcdt_uint16 := ((*C.struct_cdt_metaffi_uint16)(C.convert_union_to_ptr(unsafe.Pointer(&pcdt.cdt_val))))
		res = uint16(pcdt_uint16.val)

	case IDL.METAFFI_TYPE_UINT16_ARRAY:
		pcdt_uint16_arr := ((*C.struct_cdt_metaffi_uint16_array)(C.convert_union_to_ptr(unsafe.Pointer(&pcdt.cdt_val))))
		res = convertToMultiDimSliceOfNumbers[C.metaffi_uint16, uint16](pcdt_uint16_arr.dimensions, pcdt_uint16_arr.dimensions_lengths, pcdt_uint16_arr.vals, C.sizeof_uint16_t)

		// For int32
	case IDL.METAFFI_TYPE_INT32:
		pcdt_int32 := ((*C.struct_cdt_metaffi_int32)(C.convert_union_to_ptr(unsafe.Pointer(&pcdt.cdt_val))))
		res = int32(pcdt_int32.val)

	case IDL.METAFFI_TYPE_INT32_ARRAY:
		pcdt_int32_arr := ((*C.struct_cdt_metaffi_int32_array)(C.convert_union_to_ptr(unsafe.Pointer(&pcdt.cdt_val))))
		res = convertToMultiDimSliceOfNumbers[C.metaffi_int32, int32](pcdt_int32_arr.dimensions, pcdt_int32_arr.dimensions_lengths, pcdt_int32_arr.vals, C.sizeof_int32_t)

		// For uint32
	case IDL.METAFFI_TYPE_UINT32:
		pcdt_uint32 := ((*C.struct_cdt_metaffi_uint32)(C.convert_union_to_ptr(unsafe.Pointer(&pcdt.cdt_val))))
		res = uint32(pcdt_uint32.val)

	case IDL.METAFFI_TYPE_UINT32_ARRAY:
		pcdt_uint32_arr := ((*C.struct_cdt_metaffi_uint32_array)(C.convert_union_to_ptr(unsafe.Pointer(&pcdt.cdt_val))))
		res = convertToMultiDimSliceOfNumbers[C.metaffi_uint32, uint32](pcdt_uint32_arr.dimensions, pcdt_uint32_arr.dimensions_lengths, pcdt_uint32_arr.vals, C.sizeof_uint32_t)

		// For int64
	case IDL.METAFFI_TYPE_INT64:
		pcdt_int64 := (*C.struct_cdt_metaffi_int64)(C.convert_union_to_ptr(unsafe.Pointer(&pcdt.cdt_val)))
		res = int64(pcdt_int64.val)

	case IDL.METAFFI_TYPE_INT64_ARRAY:
		pcdt_int64_arr := (*C.struct_cdt_metaffi_int64_array)(C.convert_union_to_ptr(unsafe.Pointer(&pcdt.cdt_val)))
		res = convertToMultiDimSliceOfNumbers[C.metaffi_int64, int64](pcdt_int64_arr.dimensions, pcdt_int64_arr.dimensions_lengths, pcdt_int64_arr.vals, C.sizeof_metaffi_int64)

		// For uint64
	case IDL.METAFFI_TYPE_UINT64:
		pcdt_uint64 := (*C.struct_cdt_metaffi_uint64)(C.convert_union_to_ptr(unsafe.Pointer(&pcdt.cdt_val)))
		res = uint64(pcdt_uint64.val)

	case IDL.METAFFI_TYPE_UINT64_ARRAY:
		pcdt_uint64_arr := (*C.struct_cdt_metaffi_uint64_array)(C.convert_union_to_ptr(unsafe.Pointer(&pcdt.cdt_val)))
		res = convertToMultiDimSliceOfNumbers[C.metaffi_uint64, uint64](pcdt_uint64_arr.dimensions, pcdt_uint64_arr.dimensions_lengths, pcdt_uint64_arr.vals, C.sizeof_metaffi_uint64)

	case IDL.METAFFI_TYPE_STRING8: // string8
		pcdt_string8_res := (*C.struct_cdt_metaffi_string8)(C.convert_union_to_ptr(unsafe.Pointer(&pcdt.cdt_val)))
		res = C.GoStringN(pcdt_string8_res.val, C.int(pcdt_string8_res.length))

	case IDL.METAFFI_TYPE_STRING8_ARRAY: // []string8
		pcdt_string8_res_arr := (*C.struct_cdt_metaffi_string8_array)(C.convert_union_to_ptr(unsafe.Pointer(&pcdt.cdt_val)))
		res = convertToMultiDimSliceOfStrings(pcdt_string8_res_arr)

	case IDL.METAFFI_TYPE_BOOL: // bool
		pcdt_bool_res := ((*C.struct_cdt_metaffi_bool)(C.convert_union_to_ptr(unsafe.Pointer(&pcdt.cdt_val))))
		res = pcdt_bool_res.val != C.metaffi_bool(0)

	case IDL.METAFFI_TYPE_BOOL_ARRAY: // []bool
		pcdt_bool_res_arr := (*C.struct_cdt_metaffi_bool_array)(C.convert_union_to_ptr(unsafe.Pointer(&pcdt.cdt_val)))
		res = convertToMultiDimSliceOfBools(pcdt_bool_res_arr)

	default:
		panic(fmt.Errorf("Converting from CDT to Go failed at index %v. supported type: %v", index, res_type))
	}

	return res
}

func internalCopySliceToArray[T any](input interface{}, dims int, lengths []int, sizeOfT uintptr, setElement func(*T, interface{})) *T {
	val := reflect.ValueOf(input)

	arr := (*T)(C.malloc(C.size_t(lengths[0]) * C.size_t(sizeOfT)))

	for i := 0; i < val.Len(); i++ {
		if dims == 1 {
			setElement((*T)(unsafe.Pointer(uintptr(unsafe.Pointer(arr))+uintptr(i)*sizeOfT)), val.Index(i).Interface())
		} else {
			// Recurse into the next dimension
			for i := 0; i < val.Len(); i++ {

				// TODO: pass size of pointer
				var x *C.int
				p := (**T)(unsafe.Pointer(uintptr(unsafe.Pointer(arr)) + uintptr(i)*unsafe.Sizeof(x)))

				// TODO: pass another parameter - current dim, instead of shortening "lengths"
				*p = internalCopySliceToArray(val.Index(i).Interface(), dims-1, lengths[1:], sizeOfT, setElement)
			}
		}
	}

	return arr
}

func copySliceToArray[T any](input interface{}, dimensions int, lengths **C.metaffi_size, setElement func(*T, interface{})) *T {

	var dummy T
	sliceVal := reflect.ValueOf(input)
	if sliceVal.Kind() != reflect.Slice {
		panic("Given input is not a slice")
	}

	dimsInInput := 0
	*lengths = (*C.metaffi_size)(C.malloc(C.size_t(dimensions) * C.sizeof_metaffi_size))
	goLength := make([]int, dimensions)
	for sliceVal.Kind() == reflect.Slice {
		if dimsInInput >= dimensions {
			panic(fmt.Sprintf("Given slice dimensions (%v) is larger than expected dimensions (%v)", dimsInInput, dimensions))
		}

		C.set_metaffi_size_item(*lengths, C.int(dimsInInput), C.metaffi_size(sliceVal.Len())) // store dimensions length
		goLength[dimsInInput] = sliceVal.Len()
		sliceVal = sliceVal.Index(0)

		dimsInInput++
	}

	if dimsInInput != dimensions {
		panic(fmt.Sprintf("Given slice dimensions (%v) is not equal to expected dimensions (%v)", dimsInInput, dimensions))
	}

	return internalCopySliceToArray(input, dimensions, goLength, unsafe.Sizeof(dummy), setElement)

}

func internalCopyStringSliceToStringArray(input interface{}, dims int, lengths []int, curDim int, sizeOfPointer uintptr, setElement func(p *C.metaffi_string8, l *C.metaffi_size, val string)) (*C.metaffi_string8, *C.metaffi_size) {
	val := reflect.ValueOf(input)

	arr_str := (*C.metaffi_string8)(C.malloc(C.size_t(lengths[0]) * C.size_t(sizeOfPointer)))
	arr_str_lens := (*C.metaffi_size)(C.malloc(C.size_t(lengths[0]) * C.size_t(sizeOfPointer)))

	for i := 0; i < val.Len(); i++ {
		if dims == 1 {
			setElement((*C.metaffi_string8)(unsafe.Pointer(uintptr(unsafe.Pointer(arr_str))+uintptr(i)*sizeOfPointer)), (*C.metaffi_size)(unsafe.Pointer(uintptr(unsafe.Pointer(arr_str_lens))+uintptr(i)*C.sizeof_metaffi_size)), val.Index(i).Interface().(string))
		} else {
			// Recurse into the next dimension
			for i := 0; i < val.Len(); i++ {

				// TODO: pass size of pointer
				p := (**C.metaffi_string8)(unsafe.Pointer(uintptr(unsafe.Pointer(arr_str)) + uintptr(i)*unsafe.Sizeof(sizeOfPointer)))
				l := (**C.metaffi_size)(unsafe.Pointer(uintptr(unsafe.Pointer(arr_str_lens)) + uintptr(i)*unsafe.Sizeof(sizeOfPointer)))

				// TODO: pass another parameter - current dim, instead of shortening "lengths"
				*p, *l = internalCopyStringSliceToStringArray(val.Index(i).Interface(), dims-1, lengths, curDim+1, sizeOfPointer, setElement)
			}
		}
	}

	return arr_str, arr_str_lens
}

func copyStringSliceToStringArray(input interface{}, dimensions int, lengths **C.metaffi_size, setElement func(p *C.metaffi_string8, l *C.metaffi_size, val string)) (*C.metaffi_string8, *C.metaffi_size) {
	sliceVal := reflect.ValueOf(input)
	if sliceVal.Kind() != reflect.Slice {
		panic("Given input is not a slice")
	}

	dimsInInput := 0
	*lengths = (*C.metaffi_size)(C.malloc(C.size_t(dimensions) * C.sizeof_metaffi_size))
	goLength := make([]int, dimensions)
	for sliceVal.Kind() == reflect.Slice {
		if dimsInInput >= dimensions {
			panic(fmt.Sprintf("Given slice dimensions (%v) is larger than expected dimensions (%v)", dimsInInput, dimensions))
		}

		C.set_metaffi_size_item(*lengths, C.int(dimsInInput), C.metaffi_size(sliceVal.Len())) // store dimensions length
		goLength[dimsInInput] = sliceVal.Len()
		sliceVal = sliceVal.Index(0)

		dimsInInput++
	}

	if dimsInInput != dimensions {
		panic(fmt.Sprintf("Given slice dimensions (%v) is not equal to expected dimensions (%v)", dimsInInput, dimensions))
	}

	x := C.int(0)
	return internalCopyStringSliceToStringArray(input, dimensions, goLength, 0, unsafe.Sizeof(&x), setElement)
}

func setHandleElement(p *C.struct_cdt_metaffi_handle, val interface{}) {
	if h, ok := val.(MetaFFIHandle); ok {
		(*p).val = C.metaffi_handle(h.Val)
		(*p).runtime_id = C.metaffi_size(h.RuntimeID)
	} else {

		if val == nil {
			(*p).val = C.metaffi_handle(uintptr(0))
			(*p).runtime_id = 0
			return
		} else {
			(*p).val = C.metaffi_handle(SetObject(val))
			(*p).runtime_id = GO_RUNTIME_ID
		}
	}
}

// For int8 and C.metaffi_int8
func setElementInt8ToMetaffiInt8(p *C.metaffi_int8, val interface{}) {
	v := val.(int8)
	*p = C.metaffi_int8(v)
}

// For uint8 and C.metaffi_uint8
func setElementUint8ToMetaffiUint8(p *C.metaffi_uint8, val interface{}) {
	v := val.(uint8)
	*p = C.metaffi_uint8(v)
}

// For int16 and C.metaffi_int16
func setElementInt16ToMetaffiInt16(p *C.metaffi_int16, val interface{}) {
	v := val.(int16)
	*p = C.metaffi_int16(v)
}

// For uint16 and C.metaffi_uint16
func setElementUint16ToMetaffiUint16(p *C.metaffi_uint16, val interface{}) {
	v := val.(uint16)
	*p = C.metaffi_uint16(v)
}

// For int32 and C.metaffi_int32
func setElementInt32ToMetaffiInt32(p *C.metaffi_int32, val interface{}) {
	v := val.(int32)
	*p = C.metaffi_int32(v)
}

// For uint32 and C.metaffi_uint32
func setElementUint32ToMetaffiUint32(p *C.metaffi_uint32, val interface{}) {
	v := val.(uint32)
	*p = C.metaffi_uint32(v)
}

// For int64 and C.metaffi_int64
func setElementInt64ToMetaffiInt64(p *C.metaffi_int64, val interface{}) {
	v := val.(int64)
	*p = C.metaffi_int64(v)
}

// For uint64 and C.metaffi_uint64
func setElementUint64ToMetaffiUint64(p *C.metaffi_uint64, val interface{}) {
	v := val.(uint64)
	*p = C.metaffi_uint64(v)
}

// For int and C.metaffi_int64
func setElementIntToMetaffiInt64(p *C.metaffi_int64, val interface{}) {
	v := val.(int)
	*p = C.metaffi_int64(v)
}

// For uint and C.metaffi_uint64
func setElementUintToMetaffiUint64(p *C.metaffi_uint64, val interface{}) {
	v := val.(uint)
	*p = C.metaffi_uint64(v)
}

// For float32 and C.metaffi_float32
func setElementFloat32ToMetaffiFloat32(p *C.metaffi_float32, val interface{}) {
	v := val.(float32)
	*p = C.metaffi_float32(v)
}

// For float64 and C.metaffi_float64
func setElementFloat64ToMetaffiFloat64(p *C.metaffi_float64, val interface{}) {
	v := val.(float64)
	*p = C.metaffi_float64(v)
}

func setBoolElement(p *C.metaffi_bool, val interface{}) {
	if val.(bool) {
		*p = C.metaffi_bool(1)
	} else {
		*p = C.metaffi_bool(0)
	}
}

func setCDTStringElement(p *C.struct_cdt_metaffi_string8, val interface{}) {
	p.val = C.CString(val.(string))
	p.length = C.metaffi_size(len(val.(string)))
}

func setStringElement(p *C.metaffi_string8, l *C.metaffi_size, val string) {
	*p = C.CString(val)
	*l = C.metaffi_size(len(val))
}

func FromGoToCDT(input interface{}, pdata unsafe.Pointer, t IDL.MetaFFITypeInfo, i int) {

	pcdt := C.cast_to_cdt(pdata)
	index := C.int(i)

	cdt_to_set := C.get_cdt_index(pcdt, index)

	switch t.Type {

	case IDL.METAFFI_TYPE_HANDLE:
		cdt_to_set.free_required = 1
		C.set_cdt_type(cdt_to_set, C.metaffi_handle_type)
		pcdt_out_Handle_input := ((*C.struct_cdt_metaffi_handle)(C.convert_union_to_ptr(unsafe.Pointer(&cdt_to_set.cdt_val))))
		setHandleElement(pcdt_out_Handle_input, input)

	case IDL.METAFFI_TYPE_HANDLE_ARRAY:
		cdt_to_set.free_required = 1
		C.set_cdt_type(cdt_to_set, C.metaffi_handle_array_type)
		pcdt_handle_array := ((*C.struct_cdt_metaffi_handle_array)(C.convert_union_to_ptr(unsafe.Pointer(&cdt_to_set.cdt_val))))

		pcdt_handle_array.dimensions_lengths = (*C.metaffi_size)(C.malloc(C.size_t(unsafe.Sizeof(C.metaffi_size(0))) * C.size_t(t.Dimensions)))
		pcdt_handle_array.vals = copySliceToArray[C.struct_cdt_metaffi_handle](input, t.Dimensions, &pcdt_handle_array.dimensions_lengths, setHandleElement)
		pcdt_handle_array.dimensions = C.metaffi_size(t.Dimensions)

	case IDL.METAFFI_TYPE_FLOAT64:
		cdt_to_set.free_required = 0
		C.set_cdt_type(cdt_to_set, C.metaffi_float64_type)

		pcdt_out_float64_input := ((*C.struct_cdt_metaffi_float64)(C.convert_union_to_ptr(unsafe.Pointer(&cdt_to_set.cdt_val))))
		setElementFloat64ToMetaffiFloat64(&pcdt_out_float64_input.val, input)

	case IDL.METAFFI_TYPE_FLOAT64_ARRAY:
		cdt_to_set.free_required = 1
		C.set_cdt_type(cdt_to_set, C.metaffi_float64_array_type)
		pcdt_float64_array := ((*C.struct_cdt_metaffi_float64_array)(C.convert_union_to_ptr(unsafe.Pointer(&cdt_to_set.cdt_val))))

		pcdt_float64_array.dimensions_lengths = (*C.metaffi_size)(C.malloc(C.size_t(unsafe.Sizeof(C.metaffi_size(0))) * C.size_t(t.Dimensions)))
		pcdt_float64_array.vals = copySliceToArray[C.metaffi_float64](input, t.Dimensions, &pcdt_float64_array.dimensions_lengths, setElementFloat64ToMetaffiFloat64)
		pcdt_float64_array.dimensions = C.metaffi_size(t.Dimensions)

	case IDL.METAFFI_TYPE_FLOAT32:
		cdt_to_set.free_required = 0
		C.set_cdt_type(cdt_to_set, C.metaffi_float32_type)

		pcdt_out_float32_input := ((*C.struct_cdt_metaffi_float32)(C.convert_union_to_ptr(unsafe.Pointer(&cdt_to_set.cdt_val))))
		setElementFloat32ToMetaffiFloat32(&pcdt_out_float32_input.val, input)

	case IDL.METAFFI_TYPE_FLOAT32_ARRAY:
		cdt_to_set.free_required = 1
		C.set_cdt_type(cdt_to_set, C.metaffi_float32_array_type)
		pcdt_float32_array := ((*C.struct_cdt_metaffi_float32_array)(C.convert_union_to_ptr(unsafe.Pointer(&cdt_to_set.cdt_val))))

		pcdt_float32_array.dimensions_lengths = (*C.metaffi_size)(C.malloc(C.size_t(unsafe.Sizeof(C.metaffi_size(0))) * C.size_t(t.Dimensions)))
		pcdt_float32_array.vals = copySliceToArray[C.metaffi_float32](input, t.Dimensions, &pcdt_float32_array.dimensions_lengths, setElementFloat32ToMetaffiFloat32)
		pcdt_float32_array.dimensions = C.metaffi_size(t.Dimensions)

	case IDL.METAFFI_TYPE_INT8:
		cdt_to_set.free_required = 0
		C.set_cdt_type(cdt_to_set, C.metaffi_int8_type)

		pcdt_out_int8_input := ((*C.struct_cdt_metaffi_int8)(C.convert_union_to_ptr(unsafe.Pointer(&cdt_to_set.cdt_val))))
		setElementInt8ToMetaffiInt8(&pcdt_out_int8_input.val, input)

	case IDL.METAFFI_TYPE_INT8_ARRAY:
		cdt_to_set.free_required = 1
		C.set_cdt_type(cdt_to_set, C.metaffi_int8_array_type)
		pcdt_int8_array := ((*C.struct_cdt_metaffi_int8_array)(C.convert_union_to_ptr(unsafe.Pointer(&cdt_to_set.cdt_val))))

		pcdt_int8_array.dimensions_lengths = (*C.metaffi_size)(C.malloc(C.size_t(unsafe.Sizeof(C.metaffi_size(0))) * C.size_t(t.Dimensions)))
		pcdt_int8_array.vals = copySliceToArray[C.metaffi_int8](input, t.Dimensions, &pcdt_int8_array.dimensions_lengths, setElementInt8ToMetaffiInt8)
		pcdt_int8_array.dimensions = C.metaffi_size(t.Dimensions)

	case IDL.METAFFI_TYPE_UINT8:
		cdt_to_set.free_required = 0
		C.set_cdt_type(cdt_to_set, C.metaffi_uint8_type)

		pcdt_out_uint8_input := ((*C.struct_cdt_metaffi_uint8)(C.convert_union_to_ptr(unsafe.Pointer(&cdt_to_set.cdt_val))))
		setElementUint8ToMetaffiUint8(&pcdt_out_uint8_input.val, input)

	case IDL.METAFFI_TYPE_UINT8_ARRAY:
		cdt_to_set.free_required = 1
		C.set_cdt_type(cdt_to_set, C.metaffi_uint8_array_type)
		pcdt_uint8_array := ((*C.struct_cdt_metaffi_uint8_array)(C.convert_union_to_ptr(unsafe.Pointer(&cdt_to_set.cdt_val))))

		pcdt_uint8_array.dimensions_lengths = (*C.metaffi_size)(C.malloc(C.size_t(unsafe.Sizeof(C.metaffi_size(0))) * C.size_t(t.Dimensions)))
		pcdt_uint8_array.vals = copySliceToArray[C.metaffi_uint8](input, t.Dimensions, &pcdt_uint8_array.dimensions_lengths, setElementUint8ToMetaffiUint8)
		pcdt_uint8_array.dimensions = C.metaffi_size(t.Dimensions)

	case IDL.METAFFI_TYPE_INT16:
		cdt_to_set.free_required = 0
		C.set_cdt_type(cdt_to_set, C.metaffi_int16_type)

		pcdt_out_int16_input := ((*C.struct_cdt_metaffi_int16)(C.convert_union_to_ptr(unsafe.Pointer(&cdt_to_set.cdt_val))))
		setElementInt16ToMetaffiInt16(&pcdt_out_int16_input.val, input)

	case IDL.METAFFI_TYPE_INT16_ARRAY:
		cdt_to_set.free_required = 1
		C.set_cdt_type(cdt_to_set, C.metaffi_int16_array_type)
		pcdt_int16_array := ((*C.struct_cdt_metaffi_int16_array)(C.convert_union_to_ptr(unsafe.Pointer(&cdt_to_set.cdt_val))))

		pcdt_int16_array.dimensions_lengths = (*C.metaffi_size)(C.malloc(C.size_t(unsafe.Sizeof(C.metaffi_size(0))) * C.size_t(t.Dimensions)))
		pcdt_int16_array.vals = copySliceToArray[C.metaffi_int16](input, t.Dimensions, &pcdt_int16_array.dimensions_lengths, setElementInt16ToMetaffiInt16)
		pcdt_int16_array.dimensions = C.metaffi_size(t.Dimensions)

	case IDL.METAFFI_TYPE_UINT16:
		cdt_to_set.free_required = 0
		C.set_cdt_type(cdt_to_set, C.metaffi_uint16_type)

		pcdt_out_uint16_input := ((*C.struct_cdt_metaffi_uint16)(C.convert_union_to_ptr(unsafe.Pointer(&cdt_to_set.cdt_val))))
		setElementUint16ToMetaffiUint16(&pcdt_out_uint16_input.val, input)

	case IDL.METAFFI_TYPE_UINT16_ARRAY:
		cdt_to_set.free_required = 1
		C.set_cdt_type(cdt_to_set, C.metaffi_uint16_array_type)
		pcdt_uint16_array := ((*C.struct_cdt_metaffi_uint16_array)(C.convert_union_to_ptr(unsafe.Pointer(&cdt_to_set.cdt_val))))

		pcdt_uint16_array.dimensions_lengths = (*C.metaffi_size)(C.malloc(C.size_t(unsafe.Sizeof(C.metaffi_size(0))) * C.size_t(t.Dimensions)))
		pcdt_uint16_array.vals = copySliceToArray[C.metaffi_uint16](input, t.Dimensions, &pcdt_uint16_array.dimensions_lengths, setElementUint16ToMetaffiUint16)
		pcdt_uint16_array.dimensions = C.metaffi_size(t.Dimensions)

	case IDL.METAFFI_TYPE_INT32:
		cdt_to_set.free_required = 0
		C.set_cdt_type(cdt_to_set, C.metaffi_int32_type)

		pcdt_out_int32_input := ((*C.struct_cdt_metaffi_int32)(C.convert_union_to_ptr(unsafe.Pointer(&cdt_to_set.cdt_val))))
		setElementInt32ToMetaffiInt32(&pcdt_out_int32_input.val, input)

	case IDL.METAFFI_TYPE_INT32_ARRAY:
		cdt_to_set.free_required = 1
		C.set_cdt_type(cdt_to_set, C.metaffi_int32_array_type)
		pcdt_int32_array := ((*C.struct_cdt_metaffi_int32_array)(C.convert_union_to_ptr(unsafe.Pointer(&cdt_to_set.cdt_val))))

		pcdt_int32_array.dimensions_lengths = (*C.metaffi_size)(C.malloc(C.size_t(unsafe.Sizeof(C.metaffi_size(0))) * C.size_t(t.Dimensions)))
		pcdt_int32_array.vals = copySliceToArray[C.metaffi_int32](input, t.Dimensions, &pcdt_int32_array.dimensions_lengths, setElementInt32ToMetaffiInt32)
		pcdt_int32_array.dimensions = C.metaffi_size(t.Dimensions)

		// For uint32
	case IDL.METAFFI_TYPE_UINT32:
		cdt_to_set.free_required = 0
		C.set_cdt_type(cdt_to_set, C.metaffi_uint32_type)

		pcdt_out_uint32_input := ((*C.struct_cdt_metaffi_uint32)(C.convert_union_to_ptr(unsafe.Pointer(&cdt_to_set.cdt_val))))
		setElementUint32ToMetaffiUint32(&pcdt_out_uint32_input.val, input)

	case IDL.METAFFI_TYPE_UINT32_ARRAY:
		cdt_to_set.free_required = 1
		C.set_cdt_type(cdt_to_set, C.metaffi_uint32_array_type)
		pcdt_uint32_array := ((*C.struct_cdt_metaffi_uint32_array)(C.convert_union_to_ptr(unsafe.Pointer(&cdt_to_set.cdt_val))))

		pcdt_uint32_array.dimensions_lengths = (*C.metaffi_size)(C.malloc(C.size_t(unsafe.Sizeof(C.metaffi_size(0))) * C.size_t(t.Dimensions)))
		pcdt_uint32_array.vals = copySliceToArray[C.metaffi_uint32](input, t.Dimensions, &pcdt_uint32_array.dimensions_lengths, setElementUint32ToMetaffiUint32)
		pcdt_uint32_array.dimensions = C.metaffi_size(t.Dimensions)

		// For int64
	case IDL.METAFFI_TYPE_INT64:
		cdt_to_set.free_required = 0
		C.set_cdt_type(cdt_to_set, C.metaffi_int64_type)

		pcdt_out_int64_input := ((*C.struct_cdt_metaffi_int64)(C.convert_union_to_ptr(unsafe.Pointer(&cdt_to_set.cdt_val))))

		if _, ok := input.(int64); ok {
			setElementInt64ToMetaffiInt64(&pcdt_out_int64_input.val, input)
		} else {
			setElementIntToMetaffiInt64(&pcdt_out_int64_input.val, input)
		}

	case IDL.METAFFI_TYPE_INT64_ARRAY:
		cdt_to_set.free_required = 1
		C.set_cdt_type(cdt_to_set, C.metaffi_int64_array_type)
		pcdt_int64_array := ((*C.struct_cdt_metaffi_int64_array)(C.convert_union_to_ptr(unsafe.Pointer(&cdt_to_set.cdt_val))))

		pcdt_int64_array.dimensions_lengths = (*C.metaffi_size)(C.malloc(C.size_t(unsafe.Sizeof(C.metaffi_size(0))) * C.size_t(t.Dimensions)))

		// check if it is "int" or "int64" and set the appropriate setElementFunc
		var setElementFunc func(*C.metaffi_int64, interface{})
		curVal := reflect.ValueOf(input)
		for curVal.Kind() == reflect.Slice {
			curVal = curVal.Index(0)
		}
		if curVal.Kind() == reflect.Int64 {
			setElementFunc = setElementInt64ToMetaffiInt64
		} else {
			setElementFunc = setElementIntToMetaffiInt64
		}

		pcdt_int64_array.vals = copySliceToArray[C.metaffi_int64](input, t.Dimensions, &pcdt_int64_array.dimensions_lengths, setElementFunc)
		pcdt_int64_array.dimensions = C.metaffi_size(t.Dimensions)

		// For uint64
	case IDL.METAFFI_TYPE_UINT64:
		cdt_to_set.free_required = 0
		C.set_cdt_type(cdt_to_set, C.metaffi_uint64_type)

		pcdt_out_uint64_input := ((*C.struct_cdt_metaffi_uint64)(C.convert_union_to_ptr(unsafe.Pointer(&cdt_to_set.cdt_val))))

		if _, ok := input.(uint64); ok {
			setElementUint64ToMetaffiUint64(&pcdt_out_uint64_input.val, input)
		} else {
			setElementUintToMetaffiUint64(&pcdt_out_uint64_input.val, input)
		}

	case IDL.METAFFI_TYPE_UINT64_ARRAY:
		cdt_to_set.free_required = 1
		C.set_cdt_type(cdt_to_set, C.metaffi_uint64_array_type)
		pcdt_uint64_array := ((*C.struct_cdt_metaffi_uint64_array)(C.convert_union_to_ptr(unsafe.Pointer(&cdt_to_set.cdt_val))))

		// check if it is "uint" or "uint64" and set the appropriate setElementFunc
		var setElementFunc func(*C.metaffi_uint64, interface{})
		curVal := reflect.ValueOf(input)
		for curVal.Kind() == reflect.Slice {
			curVal = curVal.Index(0)
		}
		if curVal.Kind() == reflect.Uint64 {
			setElementFunc = setElementUint64ToMetaffiUint64
		} else {
			setElementFunc = setElementUintToMetaffiUint64
		}

		pcdt_uint64_array.dimensions_lengths = (*C.metaffi_size)(C.malloc(C.size_t(unsafe.Sizeof(C.metaffi_size(0))) * C.size_t(t.Dimensions)))
		pcdt_uint64_array.vals = copySliceToArray[C.metaffi_uint64](input, t.Dimensions, &pcdt_uint64_array.dimensions_lengths, setElementFunc)
		pcdt_uint64_array.dimensions = C.metaffi_size(t.Dimensions)

	case IDL.METAFFI_TYPE_BOOL:
		cdt_to_set.free_required = 0
		C.set_cdt_type(cdt_to_set, C.metaffi_bool_type)

		pcdt_out_bool_input := ((*C.struct_cdt_metaffi_bool)(C.convert_union_to_ptr(unsafe.Pointer(&cdt_to_set.cdt_val))))
		setBoolElement(&pcdt_out_bool_input.val, input)

	case IDL.METAFFI_TYPE_BOOL_ARRAY:
		cdt_to_set.free_required = 1
		C.set_cdt_type(cdt_to_set, C.metaffi_bool_array_type)
		pcdt_bool_array := ((*C.struct_cdt_metaffi_bool_array)(C.convert_union_to_ptr(unsafe.Pointer(&cdt_to_set.cdt_val))))

		pcdt_bool_array.dimensions_lengths = (*C.metaffi_size)(C.malloc(C.size_t(unsafe.Sizeof(C.metaffi_size(0))) * C.size_t(t.Dimensions)))
		pcdt_bool_array.vals = copySliceToArray[C.metaffi_bool](input, t.Dimensions, &pcdt_bool_array.dimensions_lengths, setBoolElement)
		pcdt_bool_array.dimensions = C.metaffi_size(t.Dimensions)

	case IDL.METAFFI_TYPE_STRING8:
		cdt_to_set.free_required = 1
		C.set_cdt_type(cdt_to_set, C.metaffi_string8_type)

		pcdt_out_string8_input := ((*C.struct_cdt_metaffi_string8)(C.convert_union_to_ptr(unsafe.Pointer(&cdt_to_set.cdt_val))))
		setCDTStringElement(pcdt_out_string8_input, input)

	case IDL.METAFFI_TYPE_STRING8_ARRAY:
		cdt_to_set.free_required = 1
		C.set_cdt_type(cdt_to_set, C.metaffi_string8_array_type)
		pcdt_string8_array := ((*C.struct_cdt_metaffi_string8_array)(C.convert_union_to_ptr(unsafe.Pointer(&cdt_to_set.cdt_val))))

		pcdt_string8_array.dimensions_lengths = (*C.metaffi_size)(C.malloc(C.size_t(unsafe.Sizeof(C.metaffi_size(0))) * C.size_t(t.Dimensions)))
		pcdt_string8_array.vals, pcdt_string8_array.vals_sizes = copyStringSliceToStringArray(input, t.Dimensions, &pcdt_string8_array.dimensions_lengths, setStringElement)
		pcdt_string8_array.dimensions = C.metaffi_size(t.Dimensions)

	default:
		panic(fmt.Errorf("Input value %v is not of a supported type, but of type: %v", "input", t.Type))
	}
}
