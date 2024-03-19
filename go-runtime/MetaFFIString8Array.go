package metaffi

/*
#cgo !windows LDFLAGS: -L. -ldl
#cgo LDFLAGS: -Wl,--allow-multiple-definition

#include <stdlib.h>
#include <stdint.h>
#include <include/cdt_structs.h>
#include <include/cdt_capi_loader.h>
#include <include/cdt_capi_loader.c>

metaffi_string8* get_vals_from_string8_array(struct cdt_metaffi_string8_array* arr)
{
    return arr->vals;
}

metaffi_string8 get_string8_item(metaffi_string8* array, int index)
{
    return array[index];
}

struct cdt_metaffi_string8_array* get_arr_from_string8_array(struct cdt_metaffi_string8_array* arr, int index)
{
    return &(arr->arr[index]);
}

int get_dimension_from_string8_array(struct cdt_metaffi_string8_array* arr)
{
 return arr->dimension;
}

int get_length_from_string8_array(struct cdt_metaffi_string8_array* arr)
{
    return arr->length;
}

void set_vals_from_string8_array(struct cdt_metaffi_string8_array* arr, metaffi_string8* vals)
{
 arr->vals = vals;
}

void alloc_string8_arr(struct cdt_metaffi_string8_array* arr, int length)
{
 arr->arr = malloc(length * sizeof(struct cdt_metaffi_string8_array));
}

void set_dimension_from_string8_array(struct cdt_metaffi_string8_array* arr, int dimension)
{
 arr->dimension = dimension;
}

void set_length_from_string8_array(struct cdt_metaffi_string8_array* arr, int length)
{
 arr->length = length;
}

void print_string8(metaffi_string8 str)
{
	printf("given str: %s\n", str);
}

*/
import "C"
import (
	"fmt"
	"reflect"
	"unsafe"
)

type CDTMetaFFIString8Array struct{}

func (this *CDTMetaFFIString8Array) getVals(pointer unsafe.Pointer) unsafe.Pointer {
	return unsafe.Pointer(C.get_vals_from_string8_array((*C.struct_cdt_metaffi_string8_array)(pointer)))
}

func (this *CDTMetaFFIString8Array) setVals(pointer unsafe.Pointer, vals unsafe.Pointer, length int) {
	C.set_vals_from_string8_array((*C.struct_cdt_metaffi_string8_array)(pointer), (*C.metaffi_string8)(vals))
	C.set_length_from_string8_array((*C.struct_cdt_metaffi_string8_array)(pointer), C.int(length))
}

func (this *CDTMetaFFIString8Array) getArr(pointer unsafe.Pointer, index int) unsafe.Pointer {
	return unsafe.Pointer(C.get_arr_from_string8_array((*C.struct_cdt_metaffi_string8_array)(pointer), C.int(index)))
}

func (this *CDTMetaFFIString8Array) allocArr(pointer unsafe.Pointer, length int) {
	C.alloc_string8_arr((*C.struct_cdt_metaffi_string8_array)(pointer), C.int(length))
}

func (this *CDTMetaFFIString8Array) getDimension(pointer unsafe.Pointer) int {
	return int(C.get_dimension_from_string8_array((*C.struct_cdt_metaffi_string8_array)(pointer)))
}

func (this *CDTMetaFFIString8Array) setDimension(pointer unsafe.Pointer, dimension int) {
	C.set_dimension_from_string8_array((*C.struct_cdt_metaffi_string8_array)(pointer), C.int(dimension))
}

func (this *CDTMetaFFIString8Array) getLength(pointer unsafe.Pointer) int {
	return int(C.get_length_from_string8_array((*C.struct_cdt_metaffi_string8_array)(pointer)))
}

func (this *CDTMetaFFIString8Array) setLength(pointer unsafe.Pointer, length int) {
	C.set_length_from_string8_array((*C.struct_cdt_metaffi_string8_array)(pointer), C.int(length))
}

func (this *CDTMetaFFIString8Array) getElement(pointer unsafe.Pointer, index int) string {
	return C.GoString((*C.char)(unsafe.Pointer(C.get_string8_item((*C.metaffi_string8)(pointer), C.int(index)))))
}

func GoStringToMetaffiString8(val interface{}) C.metaffi_string8 {
	return C.metaffi_string8(unsafe.Pointer(C.CString(val.(string))))
}

func Get1DGoString8Array[T any](index []int, otherArray interface{}, _ int) (out1DArray unsafe.Pointer, out1DArrayLength int) {
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
	s := C.size_t(v.Len() * int(unsafe.Sizeof(uintptr(0))))
	cArray := (*[1 << 30]*C.char)(C.malloc(s))

	// Copy the Go array to the C-array
	for i := 0; i < v.Len(); i++ {
		goStr := v.Index(i).String()
		cStr := C.CString(goStr)
		cArray[i] = cStr
	}

	return unsafe.Pointer(cArray), v.Len()
}
