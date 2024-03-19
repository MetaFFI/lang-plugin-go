package metaffi

/*
#cgo !windows LDFLAGS: -L. -ldl
#cgo LDFLAGS: -Wl,--allow-multiple-definition

#include <stdlib.h>
#include <stdint.h>
#include <include/cdt_structs.h>
#include <include/cdt_capi_loader.h>
#include <include/cdt_capi_loader.c>

metaffi_int8* get_vals_from_int8_array(struct cdt_metaffi_int8_array* arr)
{
    return arr->vals;
}

metaffi_int8 get_int8_item(metaffi_int8* array, int index)
{
    return array[index];
}

struct cdt_metaffi_int8_array* get_arr_from_int8_array(struct cdt_metaffi_int8_array* arr, int index)
{
    return &(arr->arr[index]);
}

int get_dimension_from_int8_array(struct cdt_metaffi_int8_array* arr)
{
	return arr->dimension;
}

int get_length_from_int8_array(struct cdt_metaffi_int8_array* arr)
{
    return arr->length;
}

void set_vals_from_int8_array(struct cdt_metaffi_int8_array* arr, metaffi_int8* vals)
{
	arr->vals = vals;
}

void alloc_int8_arr(struct cdt_metaffi_int8_array* arr, int length)
{
	arr->arr = malloc(length * sizeof(struct cdt_metaffi_int8_array));
}

void set_dimension_from_int8_array(struct cdt_metaffi_int8_array* arr, int dimension)
{
	arr->dimension = dimension;
}

void set_length_from_int8_array(struct cdt_metaffi_int8_array* arr, int length)
{
	arr->length = length;
}

*/
import "C"
import (
	"reflect"
	"unsafe"
)

type CDTMetaFFIInt8Array struct{}

func (this *CDTMetaFFIInt8Array) getVals(pointer unsafe.Pointer) unsafe.Pointer {
	return unsafe.Pointer(C.get_vals_from_int8_array((*C.struct_cdt_metaffi_int8_array)(pointer)))
}

func (this *CDTMetaFFIInt8Array) setVals(pointer unsafe.Pointer, vals unsafe.Pointer, length int) {
	C.set_vals_from_int8_array((*C.struct_cdt_metaffi_int8_array)(pointer), (*C.metaffi_int8)(vals))
	C.set_length_from_int8_array((*C.struct_cdt_metaffi_int8_array)(pointer), C.int(length))
}

func (this *CDTMetaFFIInt8Array) getArr(pointer unsafe.Pointer, index int) unsafe.Pointer {
	return unsafe.Pointer(C.get_arr_from_int8_array((*C.struct_cdt_metaffi_int8_array)(pointer), C.int(index)))
}

func (this *CDTMetaFFIInt8Array) allocArr(pointer unsafe.Pointer, length int) {
	C.alloc_int8_arr((*C.struct_cdt_metaffi_int8_array)(pointer), C.int(length))
}

func (this *CDTMetaFFIInt8Array) getDimension(pointer unsafe.Pointer) int {
	return int(C.get_dimension_from_int8_array((*C.struct_cdt_metaffi_int8_array)(pointer)))
}

func (this *CDTMetaFFIInt8Array) setDimension(pointer unsafe.Pointer, dimension int) {
	C.set_dimension_from_int8_array((*C.struct_cdt_metaffi_int8_array)(pointer), C.int(dimension))
}

func (this *CDTMetaFFIInt8Array) getLength(pointer unsafe.Pointer) int {
	return int(C.get_length_from_int8_array((*C.struct_cdt_metaffi_int8_array)(pointer)))
}

func (this *CDTMetaFFIInt8Array) setLength(pointer unsafe.Pointer, length int) {
	C.set_length_from_int8_array((*C.struct_cdt_metaffi_int8_array)(pointer), C.int(length))
}

func (this *CDTMetaFFIInt8Array) getElement(pointer unsafe.Pointer, index int) int8 {
	return int8(C.get_int8_item((*C.metaffi_int8)(pointer), C.int(index)))
}

func GoInt8ToMetaffiInt8(val interface{}) C.metaffi_int8 {
	return C.metaffi_int8(reflect.ValueOf(val).Convert(reflect.TypeOf(int8(0))).Interface().(int8))
}
