package metaffi

/*
#cgo !windows LDFLAGS: -L. -ldl
#cgo LDFLAGS: -Wl,--allow-multiple-definition

#include <stdlib.h>
#include <stdint.h>
#include <include/cdt_structs.h>
#include <include/cdt_capi_loader.h>
#include <include/cdt_capi_loader.c>

metaffi_uint32* get_vals_from_uint32_array(struct cdt_metaffi_uint32_array* arr)
{
    return arr->vals;
}

metaffi_uint32 get_uint32_item(metaffi_uint32* array, int index)
{
    return array[index];
}

struct cdt_metaffi_uint32_array* get_arr_from_uint32_array(struct cdt_metaffi_uint32_array* arr, int index)
{
    return &(arr->arr[index]);
}

int get_dimension_from_uint32_array(struct cdt_metaffi_uint32_array* arr)
{
 return arr->dimension;
}

int get_length_from_uint32_array(struct cdt_metaffi_uint32_array* arr)
{
    return arr->length;
}

void set_vals_from_uint32_array(struct cdt_metaffi_uint32_array* arr, metaffi_uint32* vals)
{
 arr->vals = vals;
}

void alloc_uint32_arr(struct cdt_metaffi_uint32_array* arr, int length)
{
 arr->arr = malloc(length * sizeof(struct cdt_metaffi_uint32_array));
}

void set_dimension_from_uint32_array(struct cdt_metaffi_uint32_array* arr, int dimension)
{
 arr->dimension = dimension;
}

void set_length_from_uint32_array(struct cdt_metaffi_uint32_array* arr, int length)
{
 arr->length = length;
}

*/
import "C"
import (
	"reflect"
	"unsafe"
)

type CDTMetaFFIUint32Array struct{}

func (this *CDTMetaFFIUint32Array) getVals(pointer unsafe.Pointer) unsafe.Pointer {
	return unsafe.Pointer(C.get_vals_from_uint32_array((*C.struct_cdt_metaffi_uint32_array)(pointer)))
}

func (this *CDTMetaFFIUint32Array) setVals(pointer unsafe.Pointer, vals unsafe.Pointer, length int) {
	C.set_vals_from_uint32_array((*C.struct_cdt_metaffi_uint32_array)(pointer), (*C.metaffi_uint32)(vals))
	C.set_length_from_uint32_array((*C.struct_cdt_metaffi_uint32_array)(pointer), C.int(length))
}

func (this *CDTMetaFFIUint32Array) getArr(pointer unsafe.Pointer, index int) unsafe.Pointer {
	return unsafe.Pointer(C.get_arr_from_uint32_array((*C.struct_cdt_metaffi_uint32_array)(pointer), C.int(index)))
}

func (this *CDTMetaFFIUint32Array) allocArr(pointer unsafe.Pointer, length int) {
	C.alloc_uint32_arr((*C.struct_cdt_metaffi_uint32_array)(pointer), C.int(length))
}

func (this *CDTMetaFFIUint32Array) getDimension(pointer unsafe.Pointer) int {
	return int(C.get_dimension_from_uint32_array((*C.struct_cdt_metaffi_uint32_array)(pointer)))
}

func (this *CDTMetaFFIUint32Array) setDimension(pointer unsafe.Pointer, dimension int) {
	C.set_dimension_from_uint32_array((*C.struct_cdt_metaffi_uint32_array)(pointer), C.int(dimension))
}

func (this *CDTMetaFFIUint32Array) getLength(pointer unsafe.Pointer) int {
	return int(C.get_length_from_uint32_array((*C.struct_cdt_metaffi_uint32_array)(pointer)))
}

func (this *CDTMetaFFIUint32Array) setLength(pointer unsafe.Pointer, length int) {
	C.set_length_from_uint32_array((*C.struct_cdt_metaffi_uint32_array)(pointer), C.int(length))
}

func (this *CDTMetaFFIUint32Array) getElement(pointer unsafe.Pointer, index int) uint32 {
	return uint32(C.get_uint32_item((*C.metaffi_uint32)(pointer), C.int(index)))
}

func GoUint32ToMetaffiUint32(val interface{}) C.metaffi_uint32 {
	return C.metaffi_uint32(reflect.ValueOf(val).Convert(reflect.TypeOf(uint32(0))).Interface().(uint32))
}
