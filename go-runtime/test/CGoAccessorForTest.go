package main

/*
#cgo !windows LDFLAGS: -L. -ldl
#cgo LDFLAGS: -Wl,--allow-multiple-definition
#cgo CFLAGS: -I"C:/src/github.com/MetaFFI/output/windows/x64/debug"

#include <stdlib.h>
#include <stdint.h>
#include <include/cdt.h>
#include <stdio.h>
#include <include/xllr_capi_loader.h>

uint64_t get_cdts_type(struct cdts* pcdts, int64_t index)
{
	return pcdts->arr[index].type;
}

struct cdt* get_arr_cdt_index(struct cdts* pcdts, int64_t index)
{
	return &pcdts->arr[index];
}

void* get_handle_value(struct cdt* cdt_handle)
{
	return cdt_handle->cdt_val.handle_val->handle;
}

void set_array_val(struct cdt *cdt, struct cdts *val) {
    cdt->cdt_val.array_val = val;
}

void set_array_val_details(struct cdts *array_val, int fixed_dimensions, int length) {
    array_val->fixed_dimensions = fixed_dimensions;
    array_val->length = length;
}

void set_cdt_val_array_val(struct cdts *cdts, struct cdts *array_val) {
    cdts->arr->cdt_val.array_val = array_val;
}

void set_float32_val(struct cdt *cdt, float val) {
    cdt->cdt_val.float32_val = val;
}

void set_uint8_val(struct cdt *cdt, uint8_t val) {
	cdt->cdt_val.uint8_val = val;
}

*/
import "C"
import "unsafe"

func GetCDTS() *C.struct_cdts {
	res := (*C.struct_cdts)(C.calloc(1, C.sizeof_struct_cdts))
	res.arr = (*C.struct_cdt)(C.calloc(1, C.sizeof_struct_cdt))
	res.length = 1
	res.fixed_dimensions = 0
	return res
}

func GetCDT(pcdts *C.struct_cdts, index int) *C.struct_cdt {
	return C.get_arr_cdt_index(pcdts, C.int64_t(index))
}

func FreeCDTS(pcdts *C.struct_cdts) {
	C.free(unsafe.Pointer(pcdts.arr))
	pcdts.arr = nil
	C.free(unsafe.Pointer(pcdts))
	pcdts = nil
}

func GetCDTHandleValue(pcdt *C.struct_cdt) uintptr {
	return uintptr(C.get_handle_value(pcdt))
}

func GetCDTSType(pcdts *C.struct_cdts, index int) uint64 {
	return uint64(C.get_cdts_type(pcdts, C.int64_t(index)))
}

func init() {
	err := C.load_xllr()
	if err != nil {
		panic("Failed to load MetaFFI XLLR functions: " + C.GoString(err))
	}
}
