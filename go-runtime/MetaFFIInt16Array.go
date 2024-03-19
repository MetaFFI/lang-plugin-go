package metaffi

/*
#cgo !windows LDFLAGS: -L. -ldl
#cgo LDFLAGS: -Wl,--allow-multiple-definition

#include <stdlib.h>
#include <stdint.h>
#include <include/cdt_structs.h>
#include <include/cdt_capi_loader.h>
#include <include/cdt_capi_loader.c>

metaffi_int16* get_vals_from_int16_array(struct cdt_metaffi_int16_array* arr)
{
    return arr->vals;
}

metaffi_int16 get_int16_item(metaffi_int16* array, int index)
{
    return array[index];
}

struct cdt_metaffi_int16_array* get_arr_from_int16_array(struct cdt_metaffi_int16_array* arr, int index)
{
    return &(arr->arr[index]);
}

int get_dimension_from_int16_array(struct cdt_metaffi_int16_array* arr)
{
 return arr->dimension;
}

int get_length_from_int16_array(struct cdt_metaffi_int16_array* arr)
{
    return arr->length;
}

void set_vals_from_int16_array(struct cdt_metaffi_int16_array* arr, metaffi_int16* vals)
{
 arr->vals = vals;
}

void alloc_int16_arr(struct cdt_metaffi_int16_array* arr, int length)
{
 arr->arr = malloc(length * sizeof(struct cdt_metaffi_int16_array));
}

void set_dimension_from_int16_array(struct cdt_metaffi_int16_array* arr, int dimension)
{
 arr->dimension = dimension;
}

void set_length_from_int16_array(struct cdt_metaffi_int16_array* arr, int length)
{
 arr->length = length;
}

*/
import "C"
import (
	"reflect"
	"unsafe"
)

type CDTMetaFFIInt16Array struct{}

func (this *CDTMetaFFIInt16Array) getVals(pointer unsafe.Pointer) unsafe.Pointer {
	return unsafe.Pointer(C.get_vals_from_int16_array((*C.struct_cdt_metaffi_int16_array)(pointer)))
}

func (this *CDTMetaFFIInt16Array) setVals(pointer unsafe.Pointer, vals unsafe.Pointer, length int) {
	C.set_vals_from_int16_array((*C.struct_cdt_metaffi_int16_array)(pointer), (*C.metaffi_int16)(vals))
	C.set_length_from_int16_array((*C.struct_cdt_metaffi_int16_array)(pointer), C.int(length))
}

func (this *CDTMetaFFIInt16Array) getArr(pointer unsafe.Pointer, index int) unsafe.Pointer {
	return unsafe.Pointer(C.get_arr_from_int16_array((*C.struct_cdt_metaffi_int16_array)(pointer), C.int(index)))
}

func (this *CDTMetaFFIInt16Array) allocArr(pointer unsafe.Pointer, length int) {
	C.alloc_int16_arr((*C.struct_cdt_metaffi_int16_array)(pointer), C.int(length))
}

func (this *CDTMetaFFIInt16Array) getDimension(pointer unsafe.Pointer) int {
	return int(C.get_dimension_from_int16_array((*C.struct_cdt_metaffi_int16_array)(pointer)))
}

func (this *CDTMetaFFIInt16Array) setDimension(pointer unsafe.Pointer, dimension int) {
	C.set_dimension_from_int16_array((*C.struct_cdt_metaffi_int16_array)(pointer), C.int(dimension))
}

func (this *CDTMetaFFIInt16Array) getLength(pointer unsafe.Pointer) int {
	return int(C.get_length_from_int16_array((*C.struct_cdt_metaffi_int16_array)(pointer)))
}

func (this *CDTMetaFFIInt16Array) setLength(pointer unsafe.Pointer, length int) {
	C.set_length_from_int16_array((*C.struct_cdt_metaffi_int16_array)(pointer), C.int(length))
}

func (this *CDTMetaFFIInt16Array) getElement(pointer unsafe.Pointer, index int) int16 {
	return int16(C.get_int16_item((*C.metaffi_int16)(pointer), C.int(index)))
}

func GoInt16ToMetaffiInt16(val interface{}) C.metaffi_int16 {
	return C.metaffi_int16(reflect.ValueOf(val).Convert(reflect.TypeOf(int16(0))).Interface().(int16))
}
