package metaffi

/*
#cgo !windows LDFLAGS: -L. -ldl
#cgo LDFLAGS: -Wl,--allow-multiple-definition

#include <stdlib.h>
#include <stdint.h>
#include <include/cdt_structs.h>
#include <include/cdt_capi_loader.h>
#include <include/cdt_capi_loader.c>

metaffi_int32* get_vals_from_int32_array(struct cdt_metaffi_int32_array* arr)
{
    return arr->vals;
}

metaffi_int32 get_int32_item(metaffi_int32* array, int index)
{
    return array[index];
}

struct cdt_metaffi_int32_array* get_arr_from_int32_array(struct cdt_metaffi_int32_array* arr, int index)
{
    return &(arr->arr[index]);
}

int get_dimension_from_int32_array(struct cdt_metaffi_int32_array* arr)
{
 return arr->dimension;
}

int get_length_from_int32_array(struct cdt_metaffi_int32_array* arr)
{
    return arr->length;
}

void set_vals_from_int32_array(struct cdt_metaffi_int32_array* arr, metaffi_int32* vals)
{
 arr->vals = vals;
}

void alloc_int32_arr(struct cdt_metaffi_int32_array* arr, int length)
{
 arr->arr = malloc(length * sizeof(struct cdt_metaffi_int32_array));
}

void set_dimension_from_int32_array(struct cdt_metaffi_int32_array* arr, int dimension)
{
 arr->dimension = dimension;
}

void set_length_from_int32_array(struct cdt_metaffi_int32_array* arr, int length)
{
 arr->length = length;
}

*/
import "C"
import (
	"reflect"
	"unsafe"
)

type CDTMetaFFIInt32Array struct{}

func (this *CDTMetaFFIInt32Array) getVals(pointer unsafe.Pointer) unsafe.Pointer {
	return unsafe.Pointer(C.get_vals_from_int32_array((*C.struct_cdt_metaffi_int32_array)(pointer)))
}

func (this *CDTMetaFFIInt32Array) setVals(pointer unsafe.Pointer, vals unsafe.Pointer, length int) {
	C.set_vals_from_int32_array((*C.struct_cdt_metaffi_int32_array)(pointer), (*C.metaffi_int32)(vals))
	C.set_length_from_int32_array((*C.struct_cdt_metaffi_int32_array)(pointer), C.int(length))
}

func (this *CDTMetaFFIInt32Array) getArr(pointer unsafe.Pointer, index int) unsafe.Pointer {
	return unsafe.Pointer(C.get_arr_from_int32_array((*C.struct_cdt_metaffi_int32_array)(pointer), C.int(index)))
}

func (this *CDTMetaFFIInt32Array) allocArr(pointer unsafe.Pointer, length int) {
	C.alloc_int32_arr((*C.struct_cdt_metaffi_int32_array)(pointer), C.int(length))
}

func (this *CDTMetaFFIInt32Array) getDimension(pointer unsafe.Pointer) int {
	return int(C.get_dimension_from_int32_array((*C.struct_cdt_metaffi_int32_array)(pointer)))
}

func (this *CDTMetaFFIInt32Array) setDimension(pointer unsafe.Pointer, dimension int) {
	C.set_dimension_from_int32_array((*C.struct_cdt_metaffi_int32_array)(pointer), C.int(dimension))
}

func (this *CDTMetaFFIInt32Array) getLength(pointer unsafe.Pointer) int {
	return int(C.get_length_from_int32_array((*C.struct_cdt_metaffi_int32_array)(pointer)))
}

func (this *CDTMetaFFIInt32Array) setLength(pointer unsafe.Pointer, length int) {
	C.set_length_from_int32_array((*C.struct_cdt_metaffi_int32_array)(pointer), C.int(length))
}

func (this *CDTMetaFFIInt32Array) getElement(pointer unsafe.Pointer, index int) int32 {
	return int32(C.get_int32_item((*C.metaffi_int32)(pointer), C.int(index)))
}

func GoInt32ToMetaffiInt32(val interface{}) C.metaffi_int32 {
	return C.metaffi_int32(reflect.ValueOf(val).Convert(reflect.TypeOf(int32(0))).Interface().(int32))
}
