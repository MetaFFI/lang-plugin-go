package metaffi

/*
#cgo !windows LDFLAGS: -L. -ldl
#cgo LDFLAGS: -Wl,--allow-multiple-definition

#include <stdlib.h>
#include <stdint.h>
#include <include/cdt_structs.h>
#include <include/cdt_capi_loader.h>
#include <include/cdt_capi_loader.c>

metaffi_uint16* get_vals_from_uint16_array(struct cdt_metaffi_uint16_array* arr)
{
    return arr->vals;
}

metaffi_uint16 get_uint16_item(metaffi_uint16* array, int index)
{
    return array[index];
}

struct cdt_metaffi_uint16_array* get_arr_from_uint16_array(struct cdt_metaffi_uint16_array* arr, int index)
{
    return &(arr->arr[index]);
}

int get_dimension_from_uint16_array(struct cdt_metaffi_uint16_array* arr)
{
 return arr->dimension;
}

int get_length_from_uint16_array(struct cdt_metaffi_uint16_array* arr)
{
    return arr->length;
}

void set_vals_from_uint16_array(struct cdt_metaffi_uint16_array* arr, metaffi_uint16* vals)
{
 arr->vals = vals;
}

void alloc_uint16_arr(struct cdt_metaffi_uint16_array* arr, int length)
{
 arr->arr = malloc(length * sizeof(struct cdt_metaffi_uint16_array));
}

void set_dimension_from_uint16_array(struct cdt_metaffi_uint16_array* arr, int dimension)
{
 arr->dimension = dimension;
}

void set_length_from_uint16_array(struct cdt_metaffi_uint16_array* arr, int length)
{
 arr->length = length;
}

*/
import "C"
import (
	"reflect"
	"unsafe"
)

type CDTMetaFFIUint16Array struct{}

func (this *CDTMetaFFIUint16Array) getVals(pointer unsafe.Pointer) unsafe.Pointer {
	return unsafe.Pointer(C.get_vals_from_uint16_array((*C.struct_cdt_metaffi_uint16_array)(pointer)))
}

func (this *CDTMetaFFIUint16Array) setVals(pointer unsafe.Pointer, vals unsafe.Pointer, length int) {
	C.set_vals_from_uint16_array((*C.struct_cdt_metaffi_uint16_array)(pointer), (*C.metaffi_uint16)(vals))
	C.set_length_from_uint16_array((*C.struct_cdt_metaffi_uint16_array)(pointer), C.int(length))
}

func (this *CDTMetaFFIUint16Array) getArr(pointer unsafe.Pointer, index int) unsafe.Pointer {
	return unsafe.Pointer(C.get_arr_from_uint16_array((*C.struct_cdt_metaffi_uint16_array)(pointer), C.int(index)))
}

func (this *CDTMetaFFIUint16Array) allocArr(pointer unsafe.Pointer, length int) {
	C.alloc_uint16_arr((*C.struct_cdt_metaffi_uint16_array)(pointer), C.int(length))
}

func (this *CDTMetaFFIUint16Array) getDimension(pointer unsafe.Pointer) int {
	return int(C.get_dimension_from_uint16_array((*C.struct_cdt_metaffi_uint16_array)(pointer)))
}

func (this *CDTMetaFFIUint16Array) setDimension(pointer unsafe.Pointer, dimension int) {
	C.set_dimension_from_uint16_array((*C.struct_cdt_metaffi_uint16_array)(pointer), C.int(dimension))
}

func (this *CDTMetaFFIUint16Array) getLength(pointer unsafe.Pointer) int {
	return int(C.get_length_from_uint16_array((*C.struct_cdt_metaffi_uint16_array)(pointer)))
}

func (this *CDTMetaFFIUint16Array) setLength(pointer unsafe.Pointer, length int) {
	C.set_length_from_uint16_array((*C.struct_cdt_metaffi_uint16_array)(pointer), C.int(length))
}

func (this *CDTMetaFFIUint16Array) getElement(pointer unsafe.Pointer, index int) uint16 {
	return uint16(C.get_uint16_item((*C.metaffi_uint16)(pointer), C.int(index)))
}

func GoUint16ToMetaffiUint16(val interface{}) C.metaffi_uint16 {
	return C.metaffi_uint16(reflect.ValueOf(val).Convert(reflect.TypeOf(uint16(0))).Interface().(uint16))
}
