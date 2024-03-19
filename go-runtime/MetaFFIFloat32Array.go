package metaffi

/*
#cgo !windows LDFLAGS: -L. -ldl
#cgo LDFLAGS: -Wl,--allow-multiple-definition

#include <stdlib.h>
#include <stdint.h>
#include <include/cdt_structs.h>
#include <include/cdt_capi_loader.h>
#include <include/cdt_capi_loader.c>

metaffi_float32* get_vals_from_float32_array(struct cdt_metaffi_float32_array* arr)
{
    return arr->vals;
}

metaffi_float32 get_float32_item(metaffi_float32* array, int index)
{
    return array[index];
}

struct cdt_metaffi_float32_array* get_arr_from_float32_array(struct cdt_metaffi_float32_array* arr, int index)
{
    return &(arr->arr[index]);
}

int get_dimension_from_float32_array(struct cdt_metaffi_float32_array* arr)
{
	return arr->dimension;
}

int get_length_from_float32_array(struct cdt_metaffi_float32_array* arr)
{
	return arr->length;
}

void set_vals_from_float32_array(struct cdt_metaffi_float32_array* arr, metaffi_float32* vals, int length)
{
	arr->vals = vals;
	arr->length = length;
}

void alloc_float32_arr(struct cdt_metaffi_float32_array* arr, int length)
{
	arr->arr = malloc(length * sizeof(struct cdt_metaffi_float32_array));
}

void set_dimension_from_float32_array(struct cdt_metaffi_float32_array* arr, int dimension)
{
	arr->dimension = dimension;
}

void set_length_from_float32_array(struct cdt_metaffi_float32_array* arr, int length)
{
	arr->length = length;
}

*/
import "C"
import (
	"reflect"
	"unsafe"
)

type CDTMetaFFIFloat32Array struct{}

func (this *CDTMetaFFIFloat32Array) getVals(pointer unsafe.Pointer) unsafe.Pointer {
	return unsafe.Pointer(C.get_vals_from_float32_array((*C.struct_cdt_metaffi_float32_array)(pointer)))
}

func (this *CDTMetaFFIFloat32Array) setVals(pointer unsafe.Pointer, vals unsafe.Pointer, length int) {
	C.set_vals_from_float32_array((*C.struct_cdt_metaffi_float32_array)(pointer), (*C.metaffi_float32)(vals), C.int(length))
}

func (this *CDTMetaFFIFloat32Array) getArr(pointer unsafe.Pointer, index int) unsafe.Pointer {
	return unsafe.Pointer(C.get_arr_from_float32_array((*C.struct_cdt_metaffi_float32_array)(pointer), C.int(index)))
}

func (this *CDTMetaFFIFloat32Array) allocArr(pointer unsafe.Pointer, length int) {
	C.alloc_float32_arr((*C.struct_cdt_metaffi_float32_array)(pointer), C.int(length))
}

func (this *CDTMetaFFIFloat32Array) getDimension(pointer unsafe.Pointer) int {
	return int(C.get_dimension_from_float32_array((*C.struct_cdt_metaffi_float32_array)(pointer)))
}

func (this *CDTMetaFFIFloat32Array) setDimension(pointer unsafe.Pointer, dimension int) {
	C.set_dimension_from_float32_array((*C.struct_cdt_metaffi_float32_array)(pointer), C.int(dimension))
}

func (this *CDTMetaFFIFloat32Array) getLength(pointer unsafe.Pointer) int {
	return int(C.get_length_from_float32_array((*C.struct_cdt_metaffi_float32_array)(pointer)))
}

func (this *CDTMetaFFIFloat32Array) setLength(pointer unsafe.Pointer, length int) {
	C.set_length_from_float32_array((*C.struct_cdt_metaffi_float32_array)(pointer), C.int(length))
}

func (this *CDTMetaFFIFloat32Array) getElement(pointer unsafe.Pointer, index int) float32 {
	return float32(C.get_float32_item((*C.metaffi_float32)(pointer), C.int(index)))
}

func GoFloat32ToMetaffiFloat32(val interface{}) C.metaffi_float32 {
	return C.metaffi_float32(reflect.ValueOf(val).Convert(reflect.TypeOf(float32(0))).Interface().(float32))
}
