package metaffi

/*
#cgo !windows LDFLAGS: -L. -ldl
#cgo LDFLAGS: -Wl,--allow-multiple-definition

#include <stdlib.h>
#include <stdint.h>
#include <uchar.h>
#include <include/cdt_structs.h>
#include <include/cdt_capi_loader.h>
#include <include/cdt_capi_loader.c>

metaffi_char32* get_vals_from_char32_array(struct cdt_metaffi_char32_array* arr)
{
    return arr->vals;
}

metaffi_char32 get_char32_item(metaffi_char32* array, int index)
{
    return array[index];
}

struct cdt_metaffi_char32_array* get_arr_from_char32_array(struct cdt_metaffi_char32_array* arr, int index)
{
    return &(arr->arr[index]);
}

int get_dimension_from_char32_array(struct cdt_metaffi_char32_array* arr)
{
	return arr->dimension;
}

int get_length_from_char32_array(struct cdt_metaffi_char32_array* arr)
{
	return arr->length;
}

void set_vals_from_char32_array(struct cdt_metaffi_char32_array* arr, metaffi_char32* vals)
{
	arr->vals = vals;
}

void alloc_char32_arr(struct cdt_metaffi_char32_array* arr, int length)
{
	arr->arr = malloc(length * sizeof(struct cdt_metaffi_char32_array));
}

void set_dimension_from_char32_array(struct cdt_metaffi_char32_array* arr, int dimension)
{
	arr->dimension = dimension;
}

void set_length_from_char32_array(struct cdt_metaffi_char32_array* arr, int length)
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

type CDTMetaFFIChar32Array struct{}

func (this *CDTMetaFFIChar32Array) getVals(pointer unsafe.Pointer) unsafe.Pointer {
	return unsafe.Pointer(C.get_vals_from_char32_array((*C.struct_cdt_metaffi_char32_array)(pointer)))
}

func (this *CDTMetaFFIChar32Array) setVals(pointer unsafe.Pointer, vals unsafe.Pointer, length int) {
	C.set_vals_from_char32_array((*C.struct_cdt_metaffi_char32_array)(pointer), (*C.metaffi_char32)(vals))
	C.set_length_from_char32_array((*C.struct_cdt_metaffi_char32_array)(pointer), C.int(length))
}

func (this *CDTMetaFFIChar32Array) getArr(pointer unsafe.Pointer, index int) unsafe.Pointer {
	return unsafe.Pointer(C.get_arr_from_char32_array((*C.struct_cdt_metaffi_char32_array)(pointer), C.int(index)))
}

func (this *CDTMetaFFIChar32Array) allocArr(pointer unsafe.Pointer, length int) {
	C.alloc_char32_arr((*C.struct_cdt_metaffi_char32_array)(pointer), C.int(length))
}

func (this *CDTMetaFFIChar32Array) getDimension(pointer unsafe.Pointer) int {
	return int(C.get_dimension_from_char32_array((*C.struct_cdt_metaffi_char32_array)(pointer)))
}

func (this *CDTMetaFFIChar32Array) setDimension(pointer unsafe.Pointer, dimension int) {
	C.set_dimension_from_char32_array((*C.struct_cdt_metaffi_char32_array)(pointer), C.int(dimension))
}

func (this *CDTMetaFFIChar32Array) getLength(pointer unsafe.Pointer) int {
	return int(C.get_length_from_char32_array((*C.struct_cdt_metaffi_char32_array)(pointer)))
}

func (this *CDTMetaFFIChar32Array) setLength(pointer unsafe.Pointer, length int) {
	C.set_length_from_char32_array((*C.struct_cdt_metaffi_char32_array)(pointer), C.int(length))
}

func (this *CDTMetaFFIChar32Array) getElement(pointer unsafe.Pointer, index int) rune {
	return rune(C.get_char32_item((*C.metaffi_char32)(pointer), C.int(index)))
}

func GoRuneToMetaffiChar32(val interface{}) C.metaffi_char32 {
	r := val.(rune)
	if r < 0 {
		panic("invalid rune. Outside of unicode point range. value out of range for conversion to C.metaffi_char32")
	}
	return C.metaffi_char32(r)
}

func Get1DGoChar32Array[T any](index []int, otherArray interface{}, elemSize int) (out1DArray unsafe.Pointer, out1DArrayLength int) {
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
		reflect.NewAt(reflect.TypeOf(rune(0)), unsafe.Pointer(uintptr(cArray)+uintptr(i)*uintptr(elemSize))).Elem().Set(v.Index(i))
	}

	return cArray, v.Len()
}
