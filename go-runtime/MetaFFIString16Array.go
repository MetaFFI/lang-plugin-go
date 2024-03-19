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

metaffi_string16* get_vals_from_string16_array(struct cdt_metaffi_string16_array* arr)
{
    return arr->vals;
}

metaffi_string16 get_string16_item(metaffi_string16* array, int index)
{
    return array[index];
}

struct cdt_metaffi_string16_array* get_arr_from_string16_array(struct cdt_metaffi_string16_array* arr, int index)
{
    return &(arr->arr[index]);
}

int get_dimension_from_string16_array(struct cdt_metaffi_string16_array* arr)
{
 return arr->dimension;
}

int get_length_from_string16_array(struct cdt_metaffi_string16_array* arr)
{
    return arr->length;
}

void set_vals_from_string16_array(struct cdt_metaffi_string16_array* arr, metaffi_string16* vals)
{
 arr->vals = vals;
}

void alloc_string16_arr(struct cdt_metaffi_string16_array* arr, int length)
{
 arr->arr = malloc(length * sizeof(struct cdt_metaffi_string16_array));
}

void set_dimension_from_string16_array(struct cdt_metaffi_string16_array* arr, int dimension)
{
 arr->dimension = dimension;
}

void set_length_from_string16_array(struct cdt_metaffi_string16_array* arr, int length)
{
 arr->length = length;
}

*/
import "C"
import (
	"fmt"
	"reflect"
	"unicode/utf16"
	"unsafe"
)

type CDTMetaFFIString16Array struct{}

func (this *CDTMetaFFIString16Array) getVals(pointer unsafe.Pointer) unsafe.Pointer {
	return unsafe.Pointer(C.get_vals_from_string16_array((*C.struct_cdt_metaffi_string16_array)(pointer)))
}

func (this *CDTMetaFFIString16Array) setVals(pointer unsafe.Pointer, vals unsafe.Pointer, length int) {
	C.set_vals_from_string16_array((*C.struct_cdt_metaffi_string16_array)(pointer), (*C.metaffi_string16)(vals))
	C.set_length_from_string16_array((*C.struct_cdt_metaffi_string16_array)(pointer), C.int(length))
}

func (this *CDTMetaFFIString16Array) getArr(pointer unsafe.Pointer, index int) unsafe.Pointer {
	return unsafe.Pointer(C.get_arr_from_string16_array((*C.struct_cdt_metaffi_string16_array)(pointer), C.int(index)))
}

func (this *CDTMetaFFIString16Array) allocArr(pointer unsafe.Pointer, length int) {
	C.alloc_string16_arr((*C.struct_cdt_metaffi_string16_array)(pointer), C.int(length))
}

func (this *CDTMetaFFIString16Array) getDimension(pointer unsafe.Pointer) int {
	return int(C.get_dimension_from_string16_array((*C.struct_cdt_metaffi_string16_array)(pointer)))
}

func (this *CDTMetaFFIString16Array) setDimension(pointer unsafe.Pointer, dimension int) {
	C.set_dimension_from_string16_array((*C.struct_cdt_metaffi_string16_array)(pointer), C.int(dimension))
}

func (this *CDTMetaFFIString16Array) getLength(pointer unsafe.Pointer) int {
	return int(C.get_length_from_string16_array((*C.struct_cdt_metaffi_string16_array)(pointer)))
}

func (this *CDTMetaFFIString16Array) setLength(pointer unsafe.Pointer, length int) {
	C.set_length_from_string16_array((*C.struct_cdt_metaffi_string16_array)(pointer), C.int(length))
}

func ConvertUTF16ToGoString(pointer unsafe.Pointer) string {
	length := 0
	for {
		if *(*uint16)(unsafe.Pointer(uintptr(pointer) + uintptr(length)*2)) == 0 {
			break
		}
		length++
	}
	goSlice := (*[1 << 30]uint16)(pointer)[:length:length]

	// Decode the UTF-16 string to a Go string
	runes := utf16.Decode(goSlice)
	return string(runes)
}

func (this *CDTMetaFFIString16Array) getElement(pointer unsafe.Pointer, index int) string {
	item := C.get_string16_item((*C.metaffi_string16)(pointer), C.int(index))
	return ConvertUTF16ToGoString(unsafe.Pointer(item))
}

func GoStringToMetaffiString16(val interface{}) C.metaffi_string16 {
	str := val.(string)
	runes := []rune(str)
	utf16Str := utf16.Encode(runes)
	return C.metaffi_string16(unsafe.Pointer(&utf16Str[0]))
}

func Get1DGoString16Array[T any](index []int, otherArray interface{}, elemSize int) (out1DArray unsafe.Pointer, out1DArrayLength int) {
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
		pstring16 := (*C.metaffi_string16)(unsafe.Pointer(uintptr(cArray) + uintptr(i)*uintptr(elemSize)))
		*pstring16 = GoStringToMetaffiString16(v.Index(i).Interface().(string))
	}

	return cArray, v.Len()
}
