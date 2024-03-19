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

metaffi_string32* get_vals_from_string32_array(struct cdt_metaffi_string32_array* arr)
{
    return arr->vals;
}

metaffi_string32 get_string32_item(metaffi_string32* array, int index)
{
    return array[index];
}

struct cdt_metaffi_string32_array* get_arr_from_string32_array(struct cdt_metaffi_string32_array* arr, int index)
{
    return &(arr->arr[index]);
}

int get_dimension_from_string32_array(struct cdt_metaffi_string32_array* arr)
{
 return arr->dimension;
}

int get_length_from_string32_array(struct cdt_metaffi_string32_array* arr)
{
    return arr->length;
}

void set_vals_from_string32_array(struct cdt_metaffi_string32_array* arr, metaffi_string32* vals)
{
 arr->vals = vals;
}

void alloc_string32_arr(struct cdt_metaffi_string32_array* arr, int length)
{
 arr->arr = malloc(length * sizeof(struct cdt_metaffi_string32_array));
}

void set_dimension_from_string32_array(struct cdt_metaffi_string32_array* arr, int dimension)
{
 arr->dimension = dimension;
}

void set_length_from_string32_array(struct cdt_metaffi_string32_array* arr, int length)
{
 arr->length = length;
}

*/
import "C"
import (
	"fmt"
	"golang.org/x/text/encoding/unicode/utf32"
	"golang.org/x/text/transform"
	"reflect"
	"unsafe"
)

type CDTMetaFFIString32Array struct{}

func (this *CDTMetaFFIString32Array) getVals(pointer unsafe.Pointer) unsafe.Pointer {
	return unsafe.Pointer(C.get_vals_from_string32_array((*C.struct_cdt_metaffi_string32_array)(pointer)))
}

func (this *CDTMetaFFIString32Array) setVals(pointer unsafe.Pointer, vals unsafe.Pointer, length int) {
	C.set_vals_from_string32_array((*C.struct_cdt_metaffi_string32_array)(pointer), (*C.metaffi_string32)(vals))
	C.set_length_from_string32_array((*C.struct_cdt_metaffi_string32_array)(pointer), C.int(length))
}

func (this *CDTMetaFFIString32Array) getArr(pointer unsafe.Pointer, index int) unsafe.Pointer {
	return unsafe.Pointer(C.get_arr_from_string32_array((*C.struct_cdt_metaffi_string32_array)(pointer), C.int(index)))
}

func (this *CDTMetaFFIString32Array) allocArr(pointer unsafe.Pointer, length int) {
	C.alloc_string32_arr((*C.struct_cdt_metaffi_string32_array)(pointer), C.int(length))
}

func (this *CDTMetaFFIString32Array) getDimension(pointer unsafe.Pointer) int {
	return int(C.get_dimension_from_string32_array((*C.struct_cdt_metaffi_string32_array)(pointer)))
}

func (this *CDTMetaFFIString32Array) setDimension(pointer unsafe.Pointer, dimension int) {
	C.set_dimension_from_string32_array((*C.struct_cdt_metaffi_string32_array)(pointer), C.int(dimension))
}

func (this *CDTMetaFFIString32Array) getLength(pointer unsafe.Pointer) int {
	return int(C.get_length_from_string32_array((*C.struct_cdt_metaffi_string32_array)(pointer)))
}

func (this *CDTMetaFFIString32Array) setLength(pointer unsafe.Pointer, length int) {
	C.set_length_from_string32_array((*C.struct_cdt_metaffi_string32_array)(pointer), C.int(length))
}

func (this *CDTMetaFFIString32Array) getElement(pointer unsafe.Pointer, index int) string {
	item := C.get_string32_item((*C.metaffi_string32)(pointer), C.int(index))

	// Convert the C array to a Go array
	length := 0
	for {
		if *(*rune)(unsafe.Pointer(uintptr(unsafe.Pointer(item)) + uintptr(length)*4)) == 0 {
			break
		}
		length++
	}
	goSlice := (*[1 << 30]rune)(unsafe.Pointer(item))[:length:length]

	// Decode the UTF-32 string to a Go string
	utf32Decoder := utf32.UTF32(utf32.BigEndian, utf32.IgnoreBOM).NewDecoder()
	str, _, err := transform.String(utf32Decoder, string(goSlice))

	if err != nil {
		panic(err)
	}

	return str
}

func GoStringToMetaffiString32(val interface{}) C.metaffi_string32 {
	str := val.(string)
	runes := []rune(str)
	utf32Encoder := utf32.UTF32(utf32.BigEndian, utf32.IgnoreBOM).NewEncoder()
	encoded := make([]byte, len(runes)*4)
	pointer := C.malloc(C.size_t(len(runes) * 4))
	for i := 0; i < len(runes); i++ {
		nDst, _, err := utf32Encoder.Transform(encoded[i*4:], []byte(string(runes[i])), true)
		if err != nil {
			panic(err)
		}
		C.memcpy(unsafe.Pointer(uintptr(pointer)+uintptr(i*4)), unsafe.Pointer(&encoded[i*4]), C.size_t(nDst))
	}
	return C.metaffi_string32(pointer)
}

func Get1DGoString32Array[T any](index []int, otherArray interface{}, elemSize int) (out1DArray unsafe.Pointer, out1DArrayLength int) {
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
		pstring32 := (*C.metaffi_string32)(unsafe.Pointer(uintptr(cArray) + uintptr(i)*uintptr(elemSize)))
		*pstring32 = GoStringToMetaffiString32(v.Index(i).Interface())
	}

	return cArray, v.Len()
}
