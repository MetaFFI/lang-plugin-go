package metaffi

/*
#cgo !windows LDFLAGS: -L. -ldl
#cgo LDFLAGS: -Wl,--allow-multiple-definition

#include <stdlib.h>
#include <stdint.h>
#include <include/cdt_structs.h>
#include <include/cdt_capi_loader.h>
#include <include/cdt_capi_loader.c>

metaffi_int64* get_vals_from_int64_array(struct cdt_metaffi_int64_array* arr)
{
    return arr->vals;
}

metaffi_int64 get_int64_item(metaffi_int64* array, int index)
{
    return array[index];
}

struct cdt_metaffi_int64_array* get_arr_from_int64_array(struct cdt_metaffi_int64_array* arr, int index)
{
    return &(arr->arr[index]);
}

int get_dimension_from_int64_array(struct cdt_metaffi_int64_array* arr)
{
 return arr->dimension;
}

int get_length_from_int64_array(struct cdt_metaffi_int64_array* arr)
{
    return arr->length;
}

void set_vals_from_int64_array(struct cdt_metaffi_int64_array* arr, metaffi_int64* vals)
{
 arr->vals = vals;
}

void alloc_int64_arr(struct cdt_metaffi_int64_array* arr, int length)
{
 arr->arr = malloc(length * sizeof(struct cdt_metaffi_int64_array));
}

void set_dimension_from_int64_array(struct cdt_metaffi_int64_array* arr, int dimension)
{
 arr->dimension = dimension;
}

void set_length_from_int64_array(struct cdt_metaffi_int64_array* arr, int length)
{
 arr->length = length;
}

*/
import "C"
import (
	"reflect"
	"unsafe"
)

type CDTMetaFFIInt64Array struct{}

func (this *CDTMetaFFIInt64Array) getVals(pointer unsafe.Pointer) unsafe.Pointer {
	return unsafe.Pointer(C.get_vals_from_int64_array((*C.struct_cdt_metaffi_int64_array)(pointer)))
}

func (this *CDTMetaFFIInt64Array) setVals(pointer unsafe.Pointer, vals unsafe.Pointer, length int) {
	C.set_vals_from_int64_array((*C.struct_cdt_metaffi_int64_array)(pointer), (*C.metaffi_int64)(vals))
	C.set_length_from_int64_array((*C.struct_cdt_metaffi_int64_array)(pointer), C.int(length))
}

func (this *CDTMetaFFIInt64Array) getArr(pointer unsafe.Pointer, index int) unsafe.Pointer {
	return unsafe.Pointer(C.get_arr_from_int64_array((*C.struct_cdt_metaffi_int64_array)(pointer), C.int(index)))
}

func (this *CDTMetaFFIInt64Array) allocArr(pointer unsafe.Pointer, length int) {
	C.alloc_int64_arr((*C.struct_cdt_metaffi_int64_array)(pointer), C.int(length))
}

func (this *CDTMetaFFIInt64Array) getDimension(pointer unsafe.Pointer) int {
	return int(C.get_dimension_from_int64_array((*C.struct_cdt_metaffi_int64_array)(pointer)))
}

func (this *CDTMetaFFIInt64Array) setDimension(pointer unsafe.Pointer, dimension int) {
	C.set_dimension_from_int64_array((*C.struct_cdt_metaffi_int64_array)(pointer), C.int(dimension))
}

func (this *CDTMetaFFIInt64Array) getLength(pointer unsafe.Pointer) int {
	return int(C.get_length_from_int64_array((*C.struct_cdt_metaffi_int64_array)(pointer)))
}

func (this *CDTMetaFFIInt64Array) setLength(pointer unsafe.Pointer, length int) {
	C.set_length_from_int64_array((*C.struct_cdt_metaffi_int64_array)(pointer), C.int(length))
}

func (this *CDTMetaFFIInt64Array) getElement(pointer unsafe.Pointer, index int) int64 {
	return int64(C.get_int64_item((*C.metaffi_int64)(pointer), C.int(index)))
}

func GoInt64ToMetaffiInt64(val interface{}) C.metaffi_int64 {
	return C.metaffi_int64(reflect.ValueOf(val).Convert(reflect.TypeOf(int64(0))).Interface().(int64))
}
