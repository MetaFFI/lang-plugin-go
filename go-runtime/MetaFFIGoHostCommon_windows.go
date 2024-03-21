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

void set_cdt_string8(struct cdt* p, metaffi_string8 val)
{
	p->cdt_val.metaffi_string8_val = val;
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

//--------------------------------------------------------------------

func setValAtIndices(arr interface{}, indices []int, val interface{}) error {
	v := reflect.ValueOf(arr)
	if v.Kind() != reflect.Slice {
		return fmt.Errorf("arr is not a slice. Kind: %v", v.Kind())
	}

	valRef := reflect.ValueOf(val)
	if valRef.Kind() != reflect.Slice {
		return fmt.Errorf("val is not a slice. Kind: %v", valRef.Kind())
	}

	for i, index := range indices {
		if v.Kind() != reflect.Slice {
			return fmt.Errorf("invalid index: not a slice")
		}
		if index < 0 || index >= v.Len() {
			return fmt.Errorf("invalid index: out of range")
		}

		if i == len(indices)-1 {
			v.Index(index).Set(valRef)
		} else {
			v = v.Index(index)
		}
	}
	return nil
}

func createMultidimArray(dimensions int, elementType reflect.Type, length int) interface{} {
	if dimensions <= 0 {
		return nil
	}

	// Start with a single-dimensional slice
	sliceType := reflect.SliceOf(elementType)

	// Add more dimensions
	for i := 1; i < dimensions; i++ {
		sliceType = reflect.SliceOf(sliceType)
	}

	// Create a slice of the final type
	slice := reflect.MakeSlice(sliceType, length, length)

	return slice.Interface()
}

type QueueItem struct {
	Array unsafe.Pointer
	Index []int
}

type cdtArray[T any] interface {
	getVals(pointer unsafe.Pointer) unsafe.Pointer
	getArr(pointer unsafe.Pointer, index int) unsafe.Pointer
	allocArr(pointer unsafe.Pointer, length int)
	getDimension(pointer unsafe.Pointer) int
	getLength(pointer unsafe.Pointer) int
	getElement(pointer unsafe.Pointer, index int) T
	setVals(pointer unsafe.Pointer, vals unsafe.Pointer, length int)
	setDimension(pointer unsafe.Pointer, dimension int)
	setLength(pointer unsafe.Pointer, length int)
}

func traverseMultiDimArray[T any](arr unsafe.Pointer, t reflect.Type, cdtFields cdtArray[T]) interface{} {
	queue := make([]QueueItem, 0)
	index := make([]int, cdtFields.getDimension(arr))
	var otherArray interface{}
	queue = append(queue, QueueItem{Array: arr, Index: index})

	for len(queue) > 0 {
		currentItem := queue[0]
		current := currentItem.Array
		index = currentItem.Index
		queue = queue[1:]

		if cdtFields.getDimension(current) == 1 {
			// Get the vals from the C struct
			vals := cdtFields.getVals(current)
			// Convert the C array to a Go slice
			curLen := cdtFields.getLength(current)
			goVals := reflect.MakeSlice(reflect.SliceOf(t), curLen, curLen)
			for i := 0; i < curLen; i++ {
				goVals.Index(i).Set(reflect.ValueOf(cdtFields.getElement(vals, i)))
			}

			if otherArray == nil {
				otherArray = goVals.Interface() // In the case we're traversing a 1D array
			} else {
				// Navigate to the correct position in the otherArray and append the goVals
				err := setValAtIndices(otherArray, index[:len(index)-1], goVals.Interface())
				if err != nil {
					panic(err)
				}
			}

		} else {
			// Create a new slice with the specified length
			newArray := createMultidimArray(cdtFields.getDimension(current), t, cdtFields.getLength(current))

			// Navigate to the correct position in the otherArray
			if otherArray == nil {
				otherArray = newArray
			} else {
				err := setValAtIndices(otherArray, index[:len(index)-cdtFields.getDimension(current)], newArray)
				if err != nil {
					panic(err)
				}
			}

			for i := 0; i < cdtFields.getLength(current); i++ {
				newIndex := make([]int, cdtFields.getDimension(arr))
				copy(newIndex, index)
				newIndex[len(index)-cdtFields.getDimension(current)] = i
				queue = append(queue, QueueItem{Array: cdtFields.getArr(current, i), Index: newIndex})
			}
		}
	}
	return otherArray
}

func FromCDTToGo(pdata unsafe.Pointer, i int, objectType reflect.Type) interface{} {

	data := C.cast_to_cdt(pdata)
	var res interface{}
	index := C.int(i)
	pcdt := C.get_cdt_index(data, index)
	res_type := C.get_cdt_type(pcdt)

	switch uint64(res_type) {

	case IDL.METAFFI_TYPE_FLOAT32:
		pfloat32 := (*C.metaffi_float32)(C.convert_union_to_ptr(unsafe.Pointer(&pcdt.cdt_val)))
		res = float32(*pfloat32)

	case IDL.METAFFI_TYPE_FLOAT32_ARRAY:
		pcdt_float32_arr := unsafe.Pointer((*C.struct_cdt_metaffi_float32_array)(C.convert_union_to_ptr(unsafe.Pointer(&pcdt.cdt_val))))
		res = traverseMultiDimArray[float32](pcdt_float32_arr, reflect.TypeOf(float32(0)), &CDTMetaFFIFloat32Array{})

	case IDL.METAFFI_TYPE_FLOAT64:
		pfloat64 := (*C.metaffi_float64)(C.convert_union_to_ptr(unsafe.Pointer(&pcdt.cdt_val)))
		res = float64(*pfloat64)

	case IDL.METAFFI_TYPE_FLOAT64_ARRAY:
		pcdt_float64_arr := unsafe.Pointer((*C.struct_cdt_metaffi_float64_array)(C.convert_union_to_ptr(unsafe.Pointer(&pcdt.cdt_val))))
		res = traverseMultiDimArray[float64](pcdt_float64_arr, reflect.TypeOf(float64(0)), &CDTMetaFFIFloat64Array{})

	case IDL.METAFFI_TYPE_INT8:
		pint8 := (*C.metaffi_int8)(C.convert_union_to_ptr(unsafe.Pointer(&pcdt.cdt_val)))
		res = int8(*pint8)

	case IDL.METAFFI_TYPE_INT8_ARRAY:
		pcdt_int8_arr := unsafe.Pointer((*C.struct_cdt_metaffi_int8_array)(C.convert_union_to_ptr(unsafe.Pointer(&pcdt.cdt_val))))
		res = traverseMultiDimArray[int8](pcdt_int8_arr, reflect.TypeOf(int8(0)), &CDTMetaFFIInt8Array{})

	case IDL.METAFFI_TYPE_UINT8:
		puint8 := (*C.metaffi_uint8)(C.convert_union_to_ptr(unsafe.Pointer(&pcdt.cdt_val)))
		res = uint8(*puint8)

	case IDL.METAFFI_TYPE_UINT8_ARRAY:
		pcdt_uint8_arr := unsafe.Pointer((*C.struct_cdt_metaffi_uint8_array)(C.convert_union_to_ptr(unsafe.Pointer(&pcdt.cdt_val))))
		res = traverseMultiDimArray[uint8](pcdt_uint8_arr, reflect.TypeOf(uint8(0)), &CDTMetaFFIUint8Array{})

		// For int16
	case IDL.METAFFI_TYPE_INT16:
		pint16 := (*C.metaffi_int16)(C.convert_union_to_ptr(unsafe.Pointer(&pcdt.cdt_val)))
		res = int16(*pint16)

	case IDL.METAFFI_TYPE_INT16_ARRAY:
		pcdt_int16_arr := (*C.struct_cdt_metaffi_int16_array)(C.convert_union_to_ptr(unsafe.Pointer(&pcdt.cdt_val)))
		res = traverseMultiDimArray[int16](unsafe.Pointer(pcdt_int16_arr), reflect.TypeOf(int16(0)), &CDTMetaFFIInt16Array{})

		// For uint16
	case IDL.METAFFI_TYPE_UINT16:
		puint16 := ((*C.metaffi_uint16)(C.convert_union_to_ptr(unsafe.Pointer(&pcdt.cdt_val))))
		res = uint16(*puint16)

	case IDL.METAFFI_TYPE_UINT16_ARRAY:
		pcdt_uint16_arr := (*C.struct_cdt_metaffi_uint16_array)(C.convert_union_to_ptr(unsafe.Pointer(&pcdt.cdt_val)))
		res = traverseMultiDimArray[uint16](unsafe.Pointer(pcdt_uint16_arr), reflect.TypeOf(uint16(0)), &CDTMetaFFIUint16Array{})

		// For int32
	case IDL.METAFFI_TYPE_INT32:
		pint32 := ((*C.metaffi_int32)(C.convert_union_to_ptr(unsafe.Pointer(&pcdt.cdt_val))))
		res = int32(*pint32)

	case IDL.METAFFI_TYPE_INT32_ARRAY:
		pcdt_int32_arr := (*C.struct_cdt_metaffi_int32_array)(C.convert_union_to_ptr(unsafe.Pointer(&pcdt.cdt_val)))
		res = traverseMultiDimArray[int32](unsafe.Pointer(pcdt_int32_arr), reflect.TypeOf(int32(0)), &CDTMetaFFIInt32Array{})

		// For uint32
	case IDL.METAFFI_TYPE_UINT32:
		puint32 := ((*C.metaffi_uint32)(C.convert_union_to_ptr(unsafe.Pointer(&pcdt.cdt_val))))
		res = uint32(*puint32)

	case IDL.METAFFI_TYPE_UINT32_ARRAY:
		pcdt_uint32_arr := (*C.struct_cdt_metaffi_uint32_array)(C.convert_union_to_ptr(unsafe.Pointer(&pcdt.cdt_val)))
		res = traverseMultiDimArray[uint32](unsafe.Pointer(pcdt_uint32_arr), reflect.TypeOf(uint32(0)), &CDTMetaFFIUint32Array{})

		// For int64
	case IDL.METAFFI_TYPE_INT64:
		pint64 := (*C.metaffi_int64)(C.convert_union_to_ptr(unsafe.Pointer(&pcdt.cdt_val)))
		res = int64(*pint64)

	case IDL.METAFFI_TYPE_INT64_ARRAY:
		pcdt_int64_arr := (*C.struct_cdt_metaffi_int64_array)(C.convert_union_to_ptr(unsafe.Pointer(&pcdt.cdt_val)))
		res = traverseMultiDimArray[int64](unsafe.Pointer(pcdt_int64_arr), reflect.TypeOf(int64(0)), &CDTMetaFFIInt64Array{})

		// For uint64
	case IDL.METAFFI_TYPE_UINT64:
		puint64 := (*C.metaffi_uint64)(C.convert_union_to_ptr(unsafe.Pointer(&pcdt.cdt_val)))
		res = uint64(*puint64)

	case IDL.METAFFI_TYPE_UINT64_ARRAY:
		pcdt_uint64_arr := (*C.struct_cdt_metaffi_uint64_array)(C.convert_union_to_ptr(unsafe.Pointer(&pcdt.cdt_val)))
		res = traverseMultiDimArray[uint64](unsafe.Pointer(pcdt_uint64_arr), reflect.TypeOf(uint64(0)), &CDTMetaFFIUint64Array{})

	case IDL.METAFFI_TYPE_BOOL: // bool
		pbool := ((*C.metaffi_bool)(C.convert_union_to_ptr(unsafe.Pointer(&pcdt.cdt_val))))
		res = *pbool != C.metaffi_bool(0)

	case IDL.METAFFI_TYPE_BOOL_ARRAY: // []bool
		pcdt_bool_arr := (*C.struct_cdt_metaffi_bool_array)(C.convert_union_to_ptr(unsafe.Pointer(&pcdt.cdt_val)))
		res = traverseMultiDimArray[bool](unsafe.Pointer(pcdt_bool_arr), reflect.TypeOf(false), &CDTMetaFFIBoolArray{})

	case IDL.METAFFI_TYPE_CHAR8: // char8
		pchar8 := (*C.metaffi_char8)(C.convert_union_to_ptr(unsafe.Pointer(&pcdt.cdt_val)))
		res = rune(*pchar8)

	case IDL.METAFFI_TYPE_CHAR8_ARRAY: // []char8
		pcdt_char8_arr := (*C.struct_cdt_metaffi_char8_array)(C.convert_union_to_ptr(unsafe.Pointer(&pcdt.cdt_val)))
		res = traverseMultiDimArray[rune](unsafe.Pointer(pcdt_char8_arr), reflect.TypeOf(' '), &CDTMetaFFIChar8Array{})

	case IDL.METAFFI_TYPE_STRING8: // string8
		pstring8 := (*C.metaffi_string8)(C.convert_union_to_ptr(unsafe.Pointer(&pcdt.cdt_val)))
		res = C.GoString((*C.char)(unsafe.Pointer(*pstring8)))

	case IDL.METAFFI_TYPE_STRING8_ARRAY: // []string8
		pcdt_string8_arr := (*C.struct_cdt_metaffi_string8_array)(C.convert_union_to_ptr(unsafe.Pointer(&pcdt.cdt_val)))
		res = traverseMultiDimArray[string](unsafe.Pointer(pcdt_string8_arr), reflect.TypeOf(""), &CDTMetaFFIString8Array{})

	case IDL.METAFFI_TYPE_CHAR16: // char16
		pchar16 := (*C.metaffi_char16)(C.convert_union_to_ptr(unsafe.Pointer(&pcdt.cdt_val)))
		res = GetRuneFromUTF16(unsafe.Pointer(pchar16))

	case IDL.METAFFI_TYPE_STRING16: // string16
		pstring16 := (*C.metaffi_string16)(C.convert_union_to_ptr(unsafe.Pointer(&pcdt.cdt_val)))
		res = ConvertUTF16ToGoString(unsafe.Pointer(pstring16))

	case IDL.METAFFI_TYPE_STRING16_ARRAY: // string16
		pcdt_string16_arr := (*C.struct_cdt_metaffi_string16_array)(C.convert_union_to_ptr(unsafe.Pointer(&pcdt.cdt_val)))
		res = traverseMultiDimArray[string](unsafe.Pointer(pcdt_string16_arr), reflect.TypeOf(""), &CDTMetaFFIString16Array{})

	case IDL.METAFFI_TYPE_HANDLE: // handle
		pcdt_handle_res := ((*C.struct_cdt_metaffi_handle)(C.convert_union_to_ptr(unsafe.Pointer(&pcdt.cdt_val))))
		res = GetGoObject(pcdt_handle_res)

	case IDL.METAFFI_TYPE_HANDLE_ARRAY:
		pcdt_handle_arr := (*C.struct_cdt_metaffi_handle_array)(C.convert_union_to_ptr(unsafe.Pointer(&pcdt.cdt_val)))
		if objectType == nil {
			objectType = reflect.TypeOf((*interface{})(nil)).Elem()
		}
		res = traverseMultiDimArray[interface{}](unsafe.Pointer(pcdt_handle_arr), objectType, &CDTMetaFFIHandleArray{})

	case IDL.METAFFI_TYPE_CALLABLE:
		panic("Callable type not implemented yet")

	default:
		panic(fmt.Errorf("Converting from CDT to Go failed at index %v. supported type: %v", index, res_type))
	}

	return res
}

func GetMetaFFITypeInfo(input interface{}) (IDL.MetaFFITypeInfo, reflect.Type) {

	if input == nil {
		return IDL.MetaFFITypeInfo{IDL.NULL, "", IDL.METAFFI_TYPE_NULL, 0}, nil
	}

	t := reflect.TypeOf(input)
	var metaFFIType IDL.MetaFFIType
	var alias string
	var dimensions int

	// Check if it's a slice
	if t.Kind() == reflect.Slice {
		dimensions = 1
		t = t.Elem()
		for t.Kind() == reflect.Slice {
			dimensions++
			t = t.Elem()
		}
	}

	// Check if it's a primitive type
	switch t.Kind() {
	case reflect.Float64:
		metaFFIType = IDL.FLOAT64
	case reflect.Float32:
		metaFFIType = IDL.FLOAT32
	case reflect.Int8:
		metaFFIType = IDL.INT8
	case reflect.Int16:
		metaFFIType = IDL.INT16
	case reflect.Int32:
		metaFFIType = IDL.INT32
	case reflect.Int64:
		metaFFIType = IDL.INT64
	case reflect.Uint8:
		metaFFIType = IDL.UINT8
	case reflect.Uint16:
		metaFFIType = IDL.UINT16
	case reflect.Uint32:
		metaFFIType = IDL.UINT32
	case reflect.Uint64:
		metaFFIType = IDL.UINT64
	case reflect.Bool:
		metaFFIType = IDL.BOOL
	case reflect.String:
		metaFFIType = IDL.STRING8
	default:
		metaFFIType = IDL.HANDLE
	}

	// If it's a slice, append "_array" to the metaFFIType
	if dimensions > 0 {
		metaFFIType = IDL.MetaFFIType(string(metaFFIType) + "_array")
	}

	// Check if it's a named type
	if t.Name() != "" && t.Name() != string(metaFFIType) {
		alias = t.Name()
	}

	return IDL.MetaFFITypeInfo{
		StringType: metaFFIType,
		Alias:      alias,
		Type:       IDL.TypeStringToTypeEnum[metaFFIType],
		Dimensions: dimensions,
	}, t
}

type getArrayFunc func(index []int, otherArray interface{}) int
type get1DArrayFunc[T any] func(index []int, otherArray interface{}, elemSize int) (out1DArray unsafe.Pointer, out1DArrayLength int)

func constructMultiDimArray[T any](arr unsafe.Pointer, cdtArrayFields cdtArray[T], elemSize int, dimensions int, otherArray interface{}, getArray getArrayFunc, get1DArray get1DArrayFunc[T]) {
	queue := make([]struct {
		Index      []int
		CurrentArr unsafe.Pointer
		CurrentDim int
	}, 0)

	index := make([]int, dimensions)
	queue = append(queue, struct {
		Index      []int
		CurrentArr unsafe.Pointer
		CurrentDim int
	}{Index: index, CurrentArr: arr, CurrentDim: dimensions})

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		if current.CurrentDim == 1 {
			newArr, length := get1DArray(current.Index[:dimensions-1], otherArray, elemSize)
			cdtArrayFields.setDimension(current.CurrentArr, 1)
			cdtArrayFields.setVals(current.CurrentArr, newArr, length)
		} else {
			length := getArray(current.Index[:dimensions-current.CurrentDim], otherArray)
			cdtArrayFields.setDimension(current.CurrentArr, current.CurrentDim)
			cdtArrayFields.setLength(current.CurrentArr, length)
			cdtArrayFields.allocArr(current.CurrentArr, length)
			for i := 0; i < length; i++ {
				newIndex := make([]int, dimensions)
				copy(newIndex, current.Index)
				newIndex[dimensions-current.CurrentDim] = i
				queue = append(queue, struct {
					Index      []int
					CurrentArr unsafe.Pointer
					CurrentDim int
				}{Index: newIndex, CurrentArr: cdtArrayFields.getArr(current.CurrentArr, i), CurrentDim: current.CurrentDim - 1})
			}
		}
	}
}

func getArray(index []int, otherArray interface{}) int {
	v := reflect.ValueOf(otherArray)

	if len(index) > 0 {
		for _, idx := range index {
			if v.Kind() != reflect.Slice {
				panic("Error: Invalid index, not a slice. Kind: " + v.Kind().String())
			}
			if idx < 0 || idx >= v.Len() {
				panic(fmt.Sprintf("Error: Invalid index, out of range. Requested index: %v, Length: %v", idx, v.Len()))
			}
			v = v.Index(idx)
		}
	}

	if v.Kind() != reflect.Slice {
		panic("Error: Final value is not a slice. Kind: " + v.Kind().String())
	}

	return v.Len()
}

func get1DArray[T any](index []int, otherArray interface{}, elemSize int) (out1DArray unsafe.Pointer, out1DArrayLength int) {

	var t T
	v := reflect.ValueOf(otherArray)

	// Traverse the multidimensional slice according to the index
	for _, idx := range index {
		if v.Kind() != reflect.Slice {
			panic("Error: Invalid index, not a slice. Kind: " + v.Kind().String())
		}
		if idx < 0 || idx >= v.Len() {
			panic(fmt.Sprintf("Error: Invalid index, out of range. Requested index: %v, Length: %v", idx, v.Len()))
		}
		v = v.Index(idx)
	}

	// Ensure the final value is a slice
	if v.Kind() != reflect.Slice {
		panic("Error: Final value is not a slice. Kind: " + v.Kind().String())
	}

	// Create a C-array of the appropriate size
	s := C.size_t(v.Len() * elemSize)
	cArray := C.malloc(s)

	// Copy the Go array to the C-array
	for i := 0; i < v.Len(); i++ {
		reflect.NewAt(reflect.TypeOf(t), unsafe.Pointer(uintptr(cArray)+uintptr(i)*uintptr(elemSize))).Elem().Set(v.Index(i))
	}

	return cArray, v.Len()
}

func FromGoToCDT(input interface{}, pdata unsafe.Pointer, t IDL.MetaFFITypeInfo, i int) {

	pcdt := C.cast_to_cdt(pdata)
	index := C.int(i)

	if t.Type == IDL.METAFFI_TYPE_ANY {
		// detect the type of the input
		t, _ = GetMetaFFITypeInfo(input)
	}

	cdt_to_set := C.get_cdt_index(pcdt, index)

	switch t.Type {
	case IDL.METAFFI_TYPE_FLOAT32:
		cdt_to_set.free_required = 0
		C.set_cdt_type(cdt_to_set, C.metaffi_float32_type)
		pfloat32 := (*C.metaffi_float32)(C.convert_union_to_ptr(unsafe.Pointer(&cdt_to_set.cdt_val)))
		*pfloat32 = GoFloat32ToMetaffiFloat32(input)

	case IDL.METAFFI_TYPE_FLOAT32_ARRAY:
		cdt_to_set.free_required = 1
		C.set_cdt_type(cdt_to_set, C.metaffi_float32_array_type)
		pcdt_float32_array := (*C.struct_cdt_metaffi_float32_array)(C.convert_union_to_ptr(unsafe.Pointer(&cdt_to_set.cdt_val)))
		constructMultiDimArray[float32](unsafe.Pointer(pcdt_float32_array), &CDTMetaFFIFloat32Array{}, C.sizeof_metaffi_float32, t.Dimensions, input, getArray, get1DArray[float32])

	case IDL.METAFFI_TYPE_FLOAT64:
		cdt_to_set.free_required = 0
		C.set_cdt_type(cdt_to_set, C.metaffi_float64_type)
		pfloat64 := (*C.metaffi_float64)(C.convert_union_to_ptr(unsafe.Pointer(&cdt_to_set.cdt_val)))
		*pfloat64 = GoFloat64ToMetaffiFloat64(input)

	case IDL.METAFFI_TYPE_FLOAT64_ARRAY:
		cdt_to_set.free_required = 1
		C.set_cdt_type(cdt_to_set, C.metaffi_float64_array_type)
		pcdt_float64_array := (*C.struct_cdt_metaffi_float64_array)(C.convert_union_to_ptr(unsafe.Pointer(&cdt_to_set.cdt_val)))
		constructMultiDimArray[float64](unsafe.Pointer(pcdt_float64_array), &CDTMetaFFIFloat64Array{}, C.sizeof_metaffi_float64, t.Dimensions, input, getArray, get1DArray[float64])

	case IDL.METAFFI_TYPE_INT8:
		cdt_to_set.free_required = 0
		C.set_cdt_type(cdt_to_set, C.metaffi_int8_type)
		pint8 := (*C.metaffi_int8)(C.convert_union_to_ptr(unsafe.Pointer(&cdt_to_set.cdt_val)))
		*pint8 = GoInt8ToMetaffiInt8(input)

	case IDL.METAFFI_TYPE_INT8_ARRAY:
		cdt_to_set.free_required = 1
		C.set_cdt_type(cdt_to_set, C.metaffi_int8_array_type)
		pcdt_int8_array := unsafe.Pointer((*C.struct_cdt_metaffi_int8_array)(C.convert_union_to_ptr(unsafe.Pointer(&cdt_to_set.cdt_val))))
		constructMultiDimArray[int8](pcdt_int8_array, &CDTMetaFFIInt8Array{}, C.sizeof_metaffi_int8, t.Dimensions, input, getArray, get1DArray[int8])

	case IDL.METAFFI_TYPE_UINT8:
		cdt_to_set.free_required = 0
		C.set_cdt_type(cdt_to_set, C.metaffi_uint8_type)
		puint8 := (*C.metaffi_uint8)(C.convert_union_to_ptr(unsafe.Pointer(&cdt_to_set.cdt_val)))
		*puint8 = GoUint8ToMetaffiUint8(input)

	case IDL.METAFFI_TYPE_UINT8_ARRAY:
		cdt_to_set.free_required = 1
		C.set_cdt_type(cdt_to_set, C.metaffi_uint8_array_type)
		pcdt_uint8_array := unsafe.Pointer((*C.struct_cdt_metaffi_uint8_array)(C.convert_union_to_ptr(unsafe.Pointer(&cdt_to_set.cdt_val))))
		constructMultiDimArray[uint8](pcdt_uint8_array, &CDTMetaFFIUint8Array{}, C.sizeof_metaffi_uint8, t.Dimensions, input, getArray, get1DArray[uint8])

	case IDL.METAFFI_TYPE_INT16:
		cdt_to_set.free_required = 0
		C.set_cdt_type(cdt_to_set, C.metaffi_int16_type)
		pint16 := (*C.metaffi_int16)(C.convert_union_to_ptr(unsafe.Pointer(&cdt_to_set.cdt_val)))
		*pint16 = GoInt16ToMetaffiInt16(input)

	case IDL.METAFFI_TYPE_INT16_ARRAY:
		cdt_to_set.free_required = 1
		C.set_cdt_type(cdt_to_set, C.metaffi_int16_array_type)
		pcdt_int16_array := unsafe.Pointer((*C.struct_cdt_metaffi_int16_array)(C.convert_union_to_ptr(unsafe.Pointer(&cdt_to_set.cdt_val))))
		constructMultiDimArray[int16](pcdt_int16_array, &CDTMetaFFIInt16Array{}, C.sizeof_metaffi_int16, t.Dimensions, input, getArray, get1DArray[int16])

	case IDL.METAFFI_TYPE_UINT16:
		cdt_to_set.free_required = 0
		C.set_cdt_type(cdt_to_set, C.metaffi_uint16_type)
		puint16 := (*C.metaffi_uint16)(C.convert_union_to_ptr(unsafe.Pointer(&cdt_to_set.cdt_val)))
		*puint16 = GoUint16ToMetaffiUint16(input)

	case IDL.METAFFI_TYPE_UINT16_ARRAY:
		cdt_to_set.free_required = 1
		C.set_cdt_type(cdt_to_set, C.metaffi_uint16_array_type)
		pcdt_uint16_array := unsafe.Pointer((*C.struct_cdt_metaffi_uint16_array)(C.convert_union_to_ptr(unsafe.Pointer(&cdt_to_set.cdt_val))))
		constructMultiDimArray[uint16](pcdt_uint16_array, &CDTMetaFFIUint16Array{}, C.sizeof_metaffi_uint16, t.Dimensions, input, getArray, get1DArray[uint16])

	case IDL.METAFFI_TYPE_INT32:
		cdt_to_set.free_required = 0
		C.set_cdt_type(cdt_to_set, C.metaffi_int32_type)
		pint32 := (*C.metaffi_int32)(C.convert_union_to_ptr(unsafe.Pointer(&cdt_to_set.cdt_val)))
		*pint32 = GoInt32ToMetaffiInt32(input)

	case IDL.METAFFI_TYPE_INT32_ARRAY:
		cdt_to_set.free_required = 1
		C.set_cdt_type(cdt_to_set, C.metaffi_int32_array_type)
		pcdt_int32_array := unsafe.Pointer((*C.struct_cdt_metaffi_int32_array)(C.convert_union_to_ptr(unsafe.Pointer(&cdt_to_set.cdt_val))))
		constructMultiDimArray[int32](pcdt_int32_array, &CDTMetaFFIInt32Array{}, C.sizeof_metaffi_int32, t.Dimensions, input, getArray, get1DArray[int32])

		// For uint32
	case IDL.METAFFI_TYPE_UINT32:
		cdt_to_set.free_required = 0
		C.set_cdt_type(cdt_to_set, C.metaffi_uint32_type)
		puint32 := (*C.metaffi_uint32)(C.convert_union_to_ptr(unsafe.Pointer(&cdt_to_set.cdt_val)))
		*puint32 = GoUint32ToMetaffiUint32(input)

	case IDL.METAFFI_TYPE_UINT32_ARRAY:
		cdt_to_set.free_required = 1
		C.set_cdt_type(cdt_to_set, C.metaffi_uint32_array_type)
		pcdt_uint32_array := unsafe.Pointer((*C.struct_cdt_metaffi_uint32_array)(C.convert_union_to_ptr(unsafe.Pointer(&cdt_to_set.cdt_val))))
		constructMultiDimArray[uint32](pcdt_uint32_array, &CDTMetaFFIUint32Array{}, C.sizeof_metaffi_uint32, t.Dimensions, input, getArray, get1DArray[uint32])

		// For int64
	case IDL.METAFFI_TYPE_INT64:
		cdt_to_set.free_required = 0
		C.set_cdt_type(cdt_to_set, C.metaffi_int64_type)
		pint64 := (*C.metaffi_int64)(C.convert_union_to_ptr(unsafe.Pointer(&cdt_to_set.cdt_val)))
		*pint64 = GoInt64ToMetaffiInt64(input)

	case IDL.METAFFI_TYPE_INT64_ARRAY:
		cdt_to_set.free_required = 1
		C.set_cdt_type(cdt_to_set, C.metaffi_int64_array_type)
		pcdt_int64_array := unsafe.Pointer((*C.struct_cdt_metaffi_int64_array)(C.convert_union_to_ptr(unsafe.Pointer(&cdt_to_set.cdt_val))))
		constructMultiDimArray[int64](pcdt_int64_array, &CDTMetaFFIInt64Array{}, C.sizeof_metaffi_int64, t.Dimensions, input, getArray, get1DArray[int64])

		// For uint64
	case IDL.METAFFI_TYPE_UINT64:
		cdt_to_set.free_required = 0
		C.set_cdt_type(cdt_to_set, C.metaffi_uint64_type)
		puint64 := (*C.metaffi_uint64)(C.convert_union_to_ptr(unsafe.Pointer(&cdt_to_set.cdt_val)))
		*puint64 = GoUint64ToMetaffiUint64(input)

	case IDL.METAFFI_TYPE_UINT64_ARRAY:
		cdt_to_set.free_required = 1
		C.set_cdt_type(cdt_to_set, C.metaffi_uint64_array_type)
		pcdt_uint64_array := unsafe.Pointer((*C.struct_cdt_metaffi_uint64_array)(C.convert_union_to_ptr(unsafe.Pointer(&cdt_to_set.cdt_val))))
		constructMultiDimArray[uint64](pcdt_uint64_array, &CDTMetaFFIUint64Array{}, C.sizeof_metaffi_uint64, t.Dimensions, input, getArray, get1DArray[uint64])

	case IDL.METAFFI_TYPE_NULL:
		cdt_to_set.free_required = 0
		C.set_cdt_type(cdt_to_set, C.metaffi_null_type)

	case IDL.METAFFI_TYPE_HANDLE:
		cdt_to_set.free_required = 1
		C.set_cdt_type(cdt_to_set, C.metaffi_handle_type)
		pcdt_handle := (*C.struct_cdt_metaffi_handle)(C.convert_union_to_ptr(unsafe.Pointer(&cdt_to_set.cdt_val)))
		GoObjectToMetaffiHandle(pcdt_handle, input)

	case IDL.METAFFI_TYPE_HANDLE_ARRAY:
		cdt_to_set.free_required = 1
		C.set_cdt_type(cdt_to_set, C.metaffi_handle_array_type)
		pcdt_handle_array := unsafe.Pointer((*C.struct_cdt_metaffi_handle_array)(C.convert_union_to_ptr(unsafe.Pointer(&cdt_to_set.cdt_val))))
		constructMultiDimArray[interface{}](pcdt_handle_array, &CDTMetaFFIHandleArray{}, C.sizeof_struct_cdt_metaffi_handle, t.Dimensions, input, getArray, Get1DGoObjectArray[interface{}])

	case IDL.METAFFI_TYPE_BOOL:
		cdt_to_set.free_required = 0
		C.set_cdt_type(cdt_to_set, C.metaffi_bool_type)
		pbool := (*C.metaffi_bool)(C.convert_union_to_ptr(unsafe.Pointer(&cdt_to_set.cdt_val)))
		*pbool = GoBoolToMetaffiBool(input)

	case IDL.METAFFI_TYPE_BOOL_ARRAY:
		cdt_to_set.free_required = 1
		C.set_cdt_type(cdt_to_set, C.metaffi_bool_array_type)
		pcdt_bool_array := unsafe.Pointer((*C.struct_cdt_metaffi_bool_array)(C.convert_union_to_ptr(unsafe.Pointer(&cdt_to_set.cdt_val))))
		constructMultiDimArray[bool](pcdt_bool_array, &CDTMetaFFIBoolArray{}, C.sizeof_metaffi_bool, t.Dimensions, input, getArray, get1DArray[bool])

	case IDL.METAFFI_TYPE_CHAR8:
		cdt_to_set.free_required = 0
		C.set_cdt_type(cdt_to_set, C.metaffi_char8_type)
		pchar8 := (*C.metaffi_char8)(C.convert_union_to_ptr(unsafe.Pointer(&cdt_to_set.cdt_val)))
		*pchar8 = GoRuneToMetaffiChar8(input)

	case IDL.METAFFI_TYPE_CHAR8_ARRAY:
		cdt_to_set.free_required = 1
		C.set_cdt_type(cdt_to_set, C.metaffi_char8_array_type)
		pcdt_char8_array := unsafe.Pointer((*C.struct_cdt_metaffi_char8_array)(C.convert_union_to_ptr(unsafe.Pointer(&cdt_to_set.cdt_val))))
		constructMultiDimArray[rune](pcdt_char8_array, &CDTMetaFFIChar8Array{}, C.sizeof_metaffi_char8, t.Dimensions, input, getArray, Get1DGoChar8Array[rune])

	case IDL.METAFFI_TYPE_STRING8:
		cdt_to_set.free_required = 1
		C.set_cdt_type(cdt_to_set, C.metaffi_string8_type)
		C.set_cdt_string8(cdt_to_set, GoStringToMetaffiString8(input))

	case IDL.METAFFI_TYPE_STRING8_ARRAY:
		cdt_to_set.free_required = 1
		C.set_cdt_type(cdt_to_set, C.metaffi_string8_array_type)
		pcdt_string8_array := unsafe.Pointer((*C.struct_cdt_metaffi_string8_array)(C.convert_union_to_ptr(unsafe.Pointer(&cdt_to_set.cdt_val))))
		constructMultiDimArray[string](pcdt_string8_array, &CDTMetaFFIString8Array{}, C.sizeof_metaffi_string8, t.Dimensions, input, getArray, Get1DGoString8Array[string])

	case IDL.METAFFI_TYPE_CHAR16:
		cdt_to_set.free_required = 0
		C.set_cdt_type(cdt_to_set, C.metaffi_char16_type)
		pchar16 := (*C.metaffi_char16)(C.convert_union_to_ptr(unsafe.Pointer(&cdt_to_set.cdt_val)))
		*pchar16 = GoRuneToMetaffiChar16(input)

	case IDL.METAFFI_TYPE_CHAR16_ARRAY:
		cdt_to_set.free_required = 1
		C.set_cdt_type(cdt_to_set, C.metaffi_char16_array_type)
		pcdt_char16_array := unsafe.Pointer((*C.struct_cdt_metaffi_char16_array)(C.convert_union_to_ptr(unsafe.Pointer(&cdt_to_set.cdt_val))))
		constructMultiDimArray[rune](pcdt_char16_array, &CDTMetaFFIChar16Array{}, C.sizeof_metaffi_char16, t.Dimensions, input, getArray, Get1DGoChar16Array[rune])

	case IDL.METAFFI_TYPE_STRING16:
		cdt_to_set.free_required = 1
		C.set_cdt_type(cdt_to_set, C.metaffi_string16_type)
		pstring16 := (*C.metaffi_string16)(C.convert_union_to_ptr(unsafe.Pointer(&cdt_to_set.cdt_val)))
		*pstring16 = GoStringToMetaffiString16(input)

	case IDL.METAFFI_TYPE_STRING16_ARRAY:
		cdt_to_set.free_required = 1
		C.set_cdt_type(cdt_to_set, C.metaffi_string16_array_type)
		pcdt_string16_array := unsafe.Pointer((*C.struct_cdt_metaffi_string16_array)(C.convert_union_to_ptr(unsafe.Pointer(&cdt_to_set.cdt_val))))
		constructMultiDimArray[string](pcdt_string16_array, &CDTMetaFFIString16Array{}, C.sizeof_metaffi_string16, t.Dimensions, input, getArray, Get1DGoString16Array[string])

	case IDL.METAFFI_TYPE_CHAR32:
		cdt_to_set.free_required = 0
		C.set_cdt_type(cdt_to_set, C.metaffi_char32_type)
		pchar32 := (*C.metaffi_char32)(C.convert_union_to_ptr(unsafe.Pointer(&cdt_to_set.cdt_val)))
		*pchar32 = GoRuneToMetaffiChar32(input)

	case IDL.METAFFI_TYPE_CHAR32_ARRAY:
		cdt_to_set.free_required = 1
		C.set_cdt_type(cdt_to_set, C.metaffi_char32_array_type)
		pcdt_char32_array := unsafe.Pointer((*C.struct_cdt_metaffi_char32_array)(C.convert_union_to_ptr(unsafe.Pointer(&cdt_to_set.cdt_val))))
		constructMultiDimArray[rune](pcdt_char32_array, &CDTMetaFFIChar32Array{}, C.sizeof_metaffi_char32, t.Dimensions, input, getArray, Get1DGoChar32Array[rune])

	case IDL.METAFFI_TYPE_STRING32:
		cdt_to_set.free_required = 1
		C.set_cdt_type(cdt_to_set, C.metaffi_string32_type)
		pstring32 := (*C.metaffi_string32)(C.convert_union_to_ptr(unsafe.Pointer(&cdt_to_set.cdt_val)))
		*pstring32 = GoStringToMetaffiString32(input)

	case IDL.METAFFI_TYPE_STRING32_ARRAY:
		cdt_to_set.free_required = 1
		C.set_cdt_type(cdt_to_set, C.metaffi_string32_array_type)
		pcdt_string32_array := unsafe.Pointer((*C.struct_cdt_metaffi_string32_array)(C.convert_union_to_ptr(unsafe.Pointer(&cdt_to_set.cdt_val))))
		constructMultiDimArray[string](pcdt_string32_array, &CDTMetaFFIString32Array{}, C.sizeof_metaffi_string32, t.Dimensions, input, getArray, Get1DGoString32Array[string])

	case IDL.METAFFI_TYPE_CALLABLE:
		panic("Callable type not implemented yet")

	default:
		panic(fmt.Errorf("%v MetaFFIType is not supported yet from Go", t.Type))
	}
}
