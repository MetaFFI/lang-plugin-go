package metaffi

/*
#cgo !windows LDFLAGS: -L. -ldl
#cgo LDFLAGS: -Wl,--allow-multiple-definition

#include <stdlib.h>
#include <stdint.h>
#include <include/cdt_structs.h>
#include <include/cdt_capi_loader.h>
#include <include/cdt_capi_loader.c>

metaffi_char8* get_vals_from_char8_array(struct cdt_metaffi_char8_array* arr)
{
    return arr->vals;
}

metaffi_char8 get_char8_item(metaffi_char8* array, int index)
{
    return array[index];
}

struct cdt_metaffi_char8_array* get_arr_from_char8_array(struct cdt_metaffi_char8_array* arr, int index)
{
    return &(arr->arr[index]);
}

int get_dimension_from_char8_array(struct cdt_metaffi_char8_array* arr)
{
 return arr->dimension;
}

int get_length_from_char8_array(struct cdt_metaffi_char8_array* arr)
{
    return arr->length;
}

void set_vals_from_char8_array(struct cdt_metaffi_char8_array* arr, metaffi_char8* vals)
{
 arr->vals = vals;
}

void alloc_char8_arr(struct cdt_metaffi_char8_array* arr, int length)
{
 arr->arr = malloc(length * sizeof(struct cdt_metaffi_char8_array));
}

void set_dimension_from_char8_array(struct cdt_metaffi_char8_array* arr, int dimension)
{
 arr->dimension = dimension;
}

void set_length_from_char8_array(struct cdt_metaffi_char8_array* arr, int length)
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

type CDTMetaFFIChar8Array struct{}

func (this *CDTMetaFFIChar8Array) getVals(pointer unsafe.Pointer) unsafe.Pointer {
	return unsafe.Pointer(C.get_vals_from_char8_array((*C.struct_cdt_metaffi_char8_array)(pointer)))
}

func (this *CDTMetaFFIChar8Array) setVals(pointer unsafe.Pointer, vals unsafe.Pointer, length int) {
	C.set_vals_from_char8_array((*C.struct_cdt_metaffi_char8_array)(pointer), (*C.metaffi_char8)(vals))
	C.set_length_from_char8_array((*C.struct_cdt_metaffi_char8_array)(pointer), C.int(length))
}

func (this *CDTMetaFFIChar8Array) getArr(pointer unsafe.Pointer, index int) unsafe.Pointer {
	return unsafe.Pointer(C.get_arr_from_char8_array((*C.struct_cdt_metaffi_char8_array)(pointer), C.int(index)))
}

func (this *CDTMetaFFIChar8Array) allocArr(pointer unsafe.Pointer, length int) {
	C.alloc_char8_arr((*C.struct_cdt_metaffi_char8_array)(pointer), C.int(length))
}

func (this *CDTMetaFFIChar8Array) getDimension(pointer unsafe.Pointer) int {
	return int(C.get_dimension_from_char8_array((*C.struct_cdt_metaffi_char8_array)(pointer)))
}

func (this *CDTMetaFFIChar8Array) setDimension(pointer unsafe.Pointer, dimension int) {
	C.set_dimension_from_char8_array((*C.struct_cdt_metaffi_char8_array)(pointer), C.int(dimension))
}

func (this *CDTMetaFFIChar8Array) getLength(pointer unsafe.Pointer) int {
	return int(C.get_length_from_char8_array((*C.struct_cdt_metaffi_char8_array)(pointer)))
}

func (this *CDTMetaFFIChar8Array) setLength(pointer unsafe.Pointer, length int) {
	C.set_length_from_char8_array((*C.struct_cdt_metaffi_char8_array)(pointer), C.int(length))
}

func (this *CDTMetaFFIChar8Array) getElement(pointer unsafe.Pointer, index int) rune {
	return rune(C.get_char8_item((*C.metaffi_char8)(pointer), C.int(index)))
}

func GoRuneToMetaffiChar8(val interface{}) C.metaffi_char8 {
	r := val.(rune)
	if r < 0 || r > 127 {
		panic("rune value out of range for conversion to C.metaffi_char8")
	}
	return C.metaffi_char8(r)
}

func Get1DGoChar8Array[T any](index []int, otherArray interface{}, elemSize int) (out1DArray unsafe.Pointer, out1DArrayLength int) {
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
		val := v.Index(i).Interface()
		convertedVal := GoRuneToMetaffiChar8(val)
		reflect.NewAt(reflect.TypeOf(convertedVal), unsafe.Pointer(uintptr(cArray)+uintptr(i)*uintptr(elemSize))).Elem().Set(reflect.ValueOf(convertedVal))
	}

	return cArray, v.Len()
}
