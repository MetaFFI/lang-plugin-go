package metaffi

/*
#cgo !windows LDFLAGS: -L. -ldl
#cgo LDFLAGS: -Wl,--allow-multiple-definition

#include <stdlib.h>
#include <stdbool.h>
#include <include/cdt_structs.h>
#include <include/cdt_capi_loader.h>
#include <include/cdt_capi_loader.c>

metaffi_bool* get_vals_from_bool_array(struct cdt_metaffi_bool_array* arr)
{
    return arr->vals;
}

metaffi_bool get_bool_item(metaffi_bool* array, int index)
{
    return array[index];
}

struct cdt_metaffi_bool_array* get_arr_from_bool_array(struct cdt_metaffi_bool_array* arr, int index)
{
    return &(arr->arr[index]);
}

int get_dimension_from_bool_array(struct cdt_metaffi_bool_array* arr)
{
 return arr->dimension;
}

int get_length_from_bool_array(struct cdt_metaffi_bool_array* arr)
{
    return arr->length;
}

void set_vals_from_bool_array(struct cdt_metaffi_bool_array* arr, metaffi_bool* vals)
{
 arr->vals = vals;
}

void alloc_bool_arr(struct cdt_metaffi_bool_array* arr, int length)
{
 arr->arr = malloc(length * sizeof(struct cdt_metaffi_bool_array));
}

void set_dimension_from_bool_array(struct cdt_metaffi_bool_array* arr, int dimension)
{
 arr->dimension = dimension;
}

void set_length_from_bool_array(struct cdt_metaffi_bool_array* arr, int length)
{
 arr->length = length;
}

*/
import "C"
import (
	"unsafe"
)

type CDTMetaFFIBoolArray struct{}

func (this *CDTMetaFFIBoolArray) getVals(pointer unsafe.Pointer) unsafe.Pointer {
	return unsafe.Pointer(C.get_vals_from_bool_array((*C.struct_cdt_metaffi_bool_array)(pointer)))
}

func (this *CDTMetaFFIBoolArray) setVals(pointer unsafe.Pointer, vals unsafe.Pointer, length int) {
	C.set_vals_from_bool_array((*C.struct_cdt_metaffi_bool_array)(pointer), (*C.metaffi_bool)(vals))
	C.set_length_from_bool_array((*C.struct_cdt_metaffi_bool_array)(pointer), C.int(length))
}

func (this *CDTMetaFFIBoolArray) getArr(pointer unsafe.Pointer, index int) unsafe.Pointer {
	return unsafe.Pointer(C.get_arr_from_bool_array((*C.struct_cdt_metaffi_bool_array)(pointer), C.int(index)))
}

func (this *CDTMetaFFIBoolArray) allocArr(pointer unsafe.Pointer, length int) {
	C.alloc_bool_arr((*C.struct_cdt_metaffi_bool_array)(pointer), C.int(length))
}

func (this *CDTMetaFFIBoolArray) getDimension(pointer unsafe.Pointer) int {
	return int(C.get_dimension_from_bool_array((*C.struct_cdt_metaffi_bool_array)(pointer)))
}

func (this *CDTMetaFFIBoolArray) setDimension(pointer unsafe.Pointer, dimension int) {
	C.set_dimension_from_bool_array((*C.struct_cdt_metaffi_bool_array)(pointer), C.int(dimension))
}

func (this *CDTMetaFFIBoolArray) getLength(pointer unsafe.Pointer) int {
	return int(C.get_length_from_bool_array((*C.struct_cdt_metaffi_bool_array)(pointer)))
}

func (this *CDTMetaFFIBoolArray) setLength(pointer unsafe.Pointer, length int) {
	C.set_length_from_bool_array((*C.struct_cdt_metaffi_bool_array)(pointer), C.int(length))
}

func (this *CDTMetaFFIBoolArray) getElement(pointer unsafe.Pointer, index int) bool {
	return C.get_bool_item((*C.metaffi_bool)(pointer), C.int(index)) != 0
}

func GoBoolToMetaffiBool(val interface{}) C.metaffi_bool {
	if val.(bool) {
		return C.metaffi_bool(1)
	} else {
		return C.metaffi_bool(0)
	}
}
