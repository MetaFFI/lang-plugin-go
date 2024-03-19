package metaffi

/*
#cgo !windows LDFLAGS: -L. -ldl
#cgo LDFLAGS: -Wl,--allow-multiple-definition

#include <stdlib.h>
#include <stdint.h>
#include <include/cdt_structs.h>
#include <include/cdt_capi_loader.h>
#include <include/cdt_capi_loader.c>

metaffi_float64* get_vals_from_float64_array(struct cdt_metaffi_float64_array* arr)
{
    return arr->vals;
}

metaffi_float64 get_float64_item(metaffi_float64* array, int index)
{
    return array[index];
}

struct cdt_metaffi_float64_array* get_arr_from_float64_array(struct cdt_metaffi_float64_array* arr, int index)
{
    return &(arr->arr[index]);
}

int get_dimension_from_float64_array(struct cdt_metaffi_float64_array* arr)
{
	return arr->dimension;
}

int get_length_from_float64_array(struct cdt_metaffi_float64_array* arr)
{
	return arr->length;
}

void set_vals_from_float64_array(struct cdt_metaffi_float64_array* arr, metaffi_float64* vals, int length)
{
	arr->vals = vals;
	arr->length = length;
}

void alloc_float64_arr(struct cdt_metaffi_float64_array* arr, int length)
{
	arr->arr = malloc(length * sizeof(struct cdt_metaffi_float64_array));
}

void set_dimension_from_float64_array(struct cdt_metaffi_float64_array* arr, int dimension)
{
	arr->dimension = dimension;
}

void set_length_from_float64_array(struct cdt_metaffi_float64_array* arr, int length)
{
	arr->length = length;
}

*/
import "C"
import (
	"reflect"
	"unsafe"
)

type CDTMetaFFIFloat64Array struct{}

func (this *CDTMetaFFIFloat64Array) getVals(pointer unsafe.Pointer) unsafe.Pointer {
	return unsafe.Pointer(C.get_vals_from_float64_array((*C.struct_cdt_metaffi_float64_array)(pointer)))
}

func (this *CDTMetaFFIFloat64Array) setVals(pointer unsafe.Pointer, vals unsafe.Pointer, length int) {
	C.set_vals_from_float64_array((*C.struct_cdt_metaffi_float64_array)(pointer), (*C.metaffi_float64)(vals), C.int(length))
}

func (this *CDTMetaFFIFloat64Array) getArr(pointer unsafe.Pointer, index int) unsafe.Pointer {
	return unsafe.Pointer(C.get_arr_from_float64_array((*C.struct_cdt_metaffi_float64_array)(pointer), C.int(index)))
}

func (this *CDTMetaFFIFloat64Array) allocArr(pointer unsafe.Pointer, length int) {
	C.alloc_float64_arr((*C.struct_cdt_metaffi_float64_array)(pointer), C.int(length))
}

func (this *CDTMetaFFIFloat64Array) getDimension(pointer unsafe.Pointer) int {
	return int(C.get_dimension_from_float64_array((*C.struct_cdt_metaffi_float64_array)(pointer)))
}

func (this *CDTMetaFFIFloat64Array) setDimension(pointer unsafe.Pointer, dimension int) {
	C.set_dimension_from_float64_array((*C.struct_cdt_metaffi_float64_array)(pointer), C.int(dimension))
}

func (this *CDTMetaFFIFloat64Array) getLength(pointer unsafe.Pointer) int {
	return int(C.get_length_from_float64_array((*C.struct_cdt_metaffi_float64_array)(pointer)))
}

func (this *CDTMetaFFIFloat64Array) setLength(pointer unsafe.Pointer, length int) {
	C.set_length_from_float64_array((*C.struct_cdt_metaffi_float64_array)(pointer), C.int(length))
}

func (this *CDTMetaFFIFloat64Array) getElement(pointer unsafe.Pointer, index int) float64 {
	return float64(C.get_float64_item((*C.metaffi_float64)(pointer), C.int(index)))
}

func GoFloat64ToMetaffiFloat64(val interface{}) C.metaffi_float64 {
	return C.metaffi_float64(reflect.ValueOf(val).Convert(reflect.TypeOf(float64(0))).Interface().(float64))
}
