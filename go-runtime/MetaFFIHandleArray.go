package metaffi

/*
#cgo !windows LDFLAGS: -L. -ldl
#cgo LDFLAGS: -Wl,--allow-multiple-definition

#include <stdlib.h>
#include <stdint.h>
#include <include/cdt_structs.h>
#include <include/cdt_capi_loader.h>
#include <include/cdt_capi_loader.c>

struct cdt_metaffi_handle* get_vals_from_handle_array(struct cdt_metaffi_handle_array* arr)
{
    return arr->vals;
}

struct cdt_metaffi_handle* get_handle_item(struct cdt_metaffi_handle* array, int index)
{
    return &array[index];
}

struct cdt_metaffi_handle_array* get_arr_from_handle_array(struct cdt_metaffi_handle_array* arr, int index)
{
    return &(arr->arr[index]);
}

int get_dimension_from_handle_array(struct cdt_metaffi_handle_array* arr)
{
 return arr->dimension;
}

int get_length_from_handle_array(struct cdt_metaffi_handle_array* arr)
{
    return arr->length;
}

void set_vals_from_handle_array(struct cdt_metaffi_handle_array* arr, struct cdt_metaffi_handle* vals)
{
 arr->vals = vals;
}

void alloc_handle_arr(struct cdt_metaffi_handle_array* arr, int length)
{
 arr->arr = malloc(length * sizeof(struct cdt_metaffi_handle_array));
}

void set_dimension_from_handle_array(struct cdt_metaffi_handle_array* arr, int dimension)
{
 arr->dimension = dimension;
}

void set_length_from_handle_array(struct cdt_metaffi_handle_array* arr, int length)
{
 arr->length = length;
}

*/
import "C"
import (
	"fmt"
	"reflect"
	"unsafe"
)

type CDTMetaFFIHandleArray struct{}

func (this *CDTMetaFFIHandleArray) getVals(pointer unsafe.Pointer) unsafe.Pointer {
	return unsafe.Pointer(C.get_vals_from_handle_array((*C.struct_cdt_metaffi_handle_array)(pointer)))
}

func (this *CDTMetaFFIHandleArray) setVals(pointer unsafe.Pointer, vals unsafe.Pointer, length int) {
	C.set_vals_from_handle_array((*C.struct_cdt_metaffi_handle_array)(pointer), (*C.struct_cdt_metaffi_handle)(vals))
	C.set_length_from_handle_array((*C.struct_cdt_metaffi_handle_array)(pointer), C.int(length))
}

func (this *CDTMetaFFIHandleArray) getArr(pointer unsafe.Pointer, index int) unsafe.Pointer {
	return unsafe.Pointer(C.get_arr_from_handle_array((*C.struct_cdt_metaffi_handle_array)(pointer), C.int(index)))
}

func (this *CDTMetaFFIHandleArray) allocArr(pointer unsafe.Pointer, length int) {
	C.alloc_handle_arr((*C.struct_cdt_metaffi_handle_array)(pointer), C.int(length))
}

func (this *CDTMetaFFIHandleArray) getDimension(pointer unsafe.Pointer) int {
	return int(C.get_dimension_from_handle_array((*C.struct_cdt_metaffi_handle_array)(pointer)))
}

func (this *CDTMetaFFIHandleArray) setDimension(pointer unsafe.Pointer, dimension int) {
	C.set_dimension_from_handle_array((*C.struct_cdt_metaffi_handle_array)(pointer), C.int(dimension))
}

func (this *CDTMetaFFIHandleArray) getLength(pointer unsafe.Pointer) int {
	return int(C.get_length_from_handle_array((*C.struct_cdt_metaffi_handle_array)(pointer)))
}

func (this *CDTMetaFFIHandleArray) setLength(pointer unsafe.Pointer, length int) {
	C.set_length_from_handle_array((*C.struct_cdt_metaffi_handle_array)(pointer), C.int(length))
}

func GetGoObject(h *C.struct_cdt_metaffi_handle) interface{} {

	if uintptr(h.val) == uintptr(0) {
		return nil
	}

	if h.runtime_id == GO_RUNTIME_ID {
		return GetObject(Handle(h.val))
	} else {
		return MetaFFIHandle{
			Val:       Handle(h.val),
			RuntimeID: uint64(h.runtime_id),
		}
	}
}

func GoObjectToMetaffiHandle(p *C.struct_cdt_metaffi_handle, val interface{}) {
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

func (this *CDTMetaFFIHandleArray) getElement(pointer unsafe.Pointer, index int) interface{} {
	return GetGoObject((*C.struct_cdt_metaffi_handle)(C.get_handle_item((*C.struct_cdt_metaffi_handle)(pointer), C.int(index))))
}

func Get1DGoObjectArray[T any](index []int, otherArray interface{}, elemSize int) (out1DArray unsafe.Pointer, out1DArrayLength int) {
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
		pcdt_handle := (*C.struct_cdt_metaffi_handle)(unsafe.Pointer(uintptr(cArray) + uintptr(i)*uintptr(elemSize)))
		GoObjectToMetaffiHandle(pcdt_handle, v.Index(i).Interface())
	}

	return cArray, v.Len()
}
