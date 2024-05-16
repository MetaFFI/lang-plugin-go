package main

/*
#cgo !windows LDFLAGS: -L. -ldl
#cgo LDFLAGS: -Wl,--allow-multiple-definition

#include <stdlib.h>
#include <stdint.h>
#include <include/cdt.h>
#include <stdio.h>

uint64_t get_cdts_type(struct cdts* pcdts, int index)
{
	return pcdts->arr[index].type;
}

struct cdt* get_arr_cdt_index(struct cdts* pcdts, int index)
{
	return &pcdts->arr[index];
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

func GetCDTS() *C.struct_cdts {
	res := (*C.struct_cdts)(C.calloc(1, C.sizeof_struct_cdts))
	res.arr = (*C.struct_cdt)(C.calloc(1, C.sizeof_struct_cdt))
	res.length = 1
	res.fixed_dimensions = 0
	return res
}

func GetCDTSType(pcdts *C.struct_cdts, index int) uint64 {
	return uint64(C.get_cdts_type(pcdts, C.int(index)))
}
