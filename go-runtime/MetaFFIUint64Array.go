package metaffi

/*
#cgo !windows LDFLAGS: -L. -ldl
#cgo LDFLAGS: -Wl,--allow-multiple-definition

#include <stdlib.h>
#include <stdint.h>
#include <include/cdt_structs.h>
#include <include/cdt_capi_loader.h>
#include <include/cdt_capi_loader.c>

metaffi_uint64* get_vals_from_uint64_array(struct cdt_metaffi_uint64_array* arr)
{
    return arr->vals;
}

metaffi_uint64 get_uint64_item(metaffi_uint64* array, int index)
{
    return array[index];
}

struct cdt_metaffi_uint64_array* get_arr_from_uint64_array(struct cdt_metaffi_uint64_array* arr, int index)
{
    return &(arr->arr[index]);
}

int get_dimension_from_uint64_array(struct cdt_metaffi_uint64_array* arr)
{
 return arr->dimension;
}

int get_length_from_uint64_array(struct cdt_metaffi_uint64_array* arr)
{
    return arr->length;
}

void set_vals_from_uint64_array(struct cdt_metaffi_uint64_array* arr, metaffi_uint64* vals)
{
 arr->vals = vals;
}

void alloc_uint64_arr(struct cdt_metaffi_uint64_array* arr, int length)
{
 arr->arr = malloc(length * sizeof(struct cdt_metaffi_uint64_array));
}

void set_dimension_from_uint64_array(struct cdt_metaffi_uint64_array* arr, int dimension)
{
 arr->dimension = dimension;
}

void set_length_from_uint64_array(struct cdt_metaffi_uint64_array* arr, int length)
{
 arr->length = length;
}

*/
import "C"
import (
	"reflect"
	"unsafe"
)

type CDTMetaFFIUint64Array struct{}

func (this *CDTMetaFFIUint64Array) getVals(pointer unsafe.Pointer) unsafe.Pointer {
	return unsafe.Pointer(C.get_vals_from_uint64_array((*C.struct_cdt_metaffi_uint64_array)(pointer)))
}

func (this *CDTMetaFFIUint64Array) setVals(pointer unsafe.Pointer, vals unsafe.Pointer, length int) {
	C.set_vals_from_uint64_array((*C.struct_cdt_metaffi_uint64_array)(pointer), (*C.metaffi_uint64)(vals))
	C.set_length_from_uint64_array((*C.struct_cdt_metaffi_uint64_array)(pointer), C.int(length))
}

func (this *CDTMetaFFIUint64Array) getArr(pointer unsafe.Pointer, index int) unsafe.Pointer {
	return unsafe.Pointer(C.get_arr_from_uint64_array((*C.struct_cdt_metaffi_uint64_array)(pointer), C.int(index)))
}

func (this *CDTMetaFFIUint64Array) allocArr(pointer unsafe.Pointer, length int) {
	C.alloc_uint64_arr((*C.struct_cdt_metaffi_uint64_array)(pointer), C.int(length))
}

func (this *CDTMetaFFIUint64Array) getDimension(pointer unsafe.Pointer) int {
	return int(C.get_dimension_from_uint64_array((*C.struct_cdt_metaffi_uint64_array)(pointer)))
}

func (this *CDTMetaFFIUint64Array) setDimension(pointer unsafe.Pointer, dimension int) {
	C.set_dimension_from_uint64_array((*C.struct_cdt_metaffi_uint64_array)(pointer), C.int(dimension))
}

func (this *CDTMetaFFIUint64Array) getLength(pointer unsafe.Pointer) int {
	return int(C.get_length_from_uint64_array((*C.struct_cdt_metaffi_uint64_array)(pointer)))
}

func (this *CDTMetaFFIUint64Array) setLength(pointer unsafe.Pointer, length int) {
	C.set_length_from_uint64_array((*C.struct_cdt_metaffi_uint64_array)(pointer), C.int(length))
}

func (this *CDTMetaFFIUint64Array) getElement(pointer unsafe.Pointer, index int) uint64 {
	return uint64(C.get_uint64_item((*C.metaffi_uint64)(pointer), C.int(index)))
}

func GoUint64ToMetaffiUint64(val interface{}) C.metaffi_uint64 {
	return C.metaffi_uint64(reflect.ValueOf(val).Convert(reflect.TypeOf(uint64(0))).Interface().(uint64))
}
