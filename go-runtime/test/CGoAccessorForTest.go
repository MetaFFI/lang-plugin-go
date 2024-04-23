package main

/*
#cgo !windows LDFLAGS: -L. -ldl
#cgo LDFLAGS: -Wl,--allow-multiple-definition

#include <stdlib.h>
#include <stdint.h>
#include <include/cdt.h>

struct cdts* get_cdts()
{
	struct cdts* res = (struct cdts*)calloc(1, sizeof(struct cdts));
	res->arr = (struct cdt*)calloc(1, sizeof(struct cdt));
	res->length = 1;
	res->fixed_dimensions = 0;
}

uint64_t get_cdts_type(struct cdts* pcdts, int index)
{
	return pcdts->arr[index].type;
}

struct cdts* get_2d_uint8_array_cdts()
{
	struct cdts* cdts_param_ret = get_cdts();
	cdts_param_ret[0].arr->type = metaffi_uint8_array_type;
	cdts_param_ret[0].arr->cdt_val.array_val.fixed_dimensions = 2;
	cdts_param_ret[0].arr->cdt_val.array_val.length = 3;
	cdts_param_ret[0].arr->cdt_val.array_val.arr = (struct cdt*)malloc(3 * sizeof(struct cdt));
	cdts_param_ret[0].arr->cdt_val.array_val.arr[0].type = metaffi_uint8_array_type;
	cdts_param_ret[0].arr->cdt_val.array_val.arr[0].cdt_val.array_val.fixed_dimensions = 1;
	cdts_param_ret[0].arr->cdt_val.array_val.arr[0].cdt_val.array_val.length = 3;
	cdts_param_ret[0].arr->cdt_val.array_val.arr[0].cdt_val.array_val.arr = (struct cdt*)malloc(3 * sizeof(struct cdt));
	cdts_param_ret[0].arr->cdt_val.array_val.arr[0].cdt_val.array_val.arr[0].type = metaffi_uint8_type;
	cdts_param_ret[0].arr->cdt_val.array_val.arr[0].cdt_val.array_val.arr[0].cdt_val.uint8_val = 0;
	cdts_param_ret[0].arr->cdt_val.array_val.arr[0].cdt_val.array_val.arr[1].type = metaffi_uint8_type;
	cdts_param_ret[0].arr->cdt_val.array_val.arr[0].cdt_val.array_val.arr[1].cdt_val.uint8_val = 1;
	cdts_param_ret[0].arr->cdt_val.array_val.arr[0].cdt_val.array_val.arr[2].type = metaffi_uint8_type;
	cdts_param_ret[0].arr->cdt_val.array_val.arr[0].cdt_val.array_val.arr[2].cdt_val.uint8_val = 2;

	cdts_param_ret[0].arr->cdt_val.array_val.arr[1].type = metaffi_uint8_array_type;
	cdts_param_ret[0].arr->cdt_val.array_val.arr[1].cdt_val.array_val.fixed_dimensions = 1;
	cdts_param_ret[0].arr->cdt_val.array_val.arr[1].cdt_val.array_val.length = 3;
	cdts_param_ret[0].arr->cdt_val.array_val.arr[1].cdt_val.array_val.arr = (struct cdt*)malloc(3 * sizeof(struct cdt));
	cdts_param_ret[0].arr->cdt_val.array_val.arr[1].cdt_val.array_val.arr[0].type = metaffi_uint8_type;
	cdts_param_ret[0].arr->cdt_val.array_val.arr[1].cdt_val.array_val.arr[0].cdt_val.uint8_val = 3;
	cdts_param_ret[0].arr->cdt_val.array_val.arr[1].cdt_val.array_val.arr[1].type = metaffi_uint8_type;
	cdts_param_ret[0].arr->cdt_val.array_val.arr[1].cdt_val.array_val.arr[1].cdt_val.uint8_val = 4;
	cdts_param_ret[0].arr->cdt_val.array_val.arr[1].cdt_val.array_val.arr[2].type = metaffi_uint8_type;
	cdts_param_ret[0].arr->cdt_val.array_val.arr[1].cdt_val.array_val.arr[2].cdt_val.uint8_val = 5;

	cdts_param_ret[0].arr->cdt_val.array_val.arr[2].type = metaffi_uint8_array_type;
	cdts_param_ret[0].arr->cdt_val.array_val.arr[2].cdt_val.array_val.fixed_dimensions = 1;
	cdts_param_ret[0].arr->cdt_val.array_val.arr[2].cdt_val.array_val.length = 3;
	cdts_param_ret[0].arr->cdt_val.array_val.arr[2].cdt_val.array_val.arr = (struct cdt*)malloc(3 * sizeof(struct cdt));
	cdts_param_ret[0].arr->cdt_val.array_val.arr[2].cdt_val.array_val.arr[0].type = metaffi_uint8_type;
	cdts_param_ret[0].arr->cdt_val.array_val.arr[2].cdt_val.array_val.arr[0].cdt_val.uint8_val = 6;
	cdts_param_ret[0].arr->cdt_val.array_val.arr[2].cdt_val.array_val.arr[1].type = metaffi_uint8_type;
	cdts_param_ret[0].arr->cdt_val.array_val.arr[2].cdt_val.array_val.arr[1].cdt_val.uint8_val = 7;
	cdts_param_ret[0].arr->cdt_val.array_val.arr[2].cdt_val.array_val.arr[2].type = metaffi_uint8_type;
	cdts_param_ret[0].arr->cdt_val.array_val.arr[2].cdt_val.array_val.arr[2].cdt_val.uint8_val = 8;

	return cdts_param_ret;
}

struct cdts* get_3d_float32_array_cdts()
{
	struct cdts* cdts_param_ret = get_cdts();
	cdts_param_ret[0].arr->type = metaffi_float32_array_type;
	cdts_param_ret[0].arr->cdt_val.array_val.fixed_dimensions = 3;
	cdts_param_ret[0].arr->cdt_val.array_val.length = 3;
	cdts_param_ret[0].arr->cdt_val.array_val.arr = (struct cdt*)malloc(3 * sizeof(struct cdt));

	for (int i = 0; i < 3; i++) {
		cdts_param_ret[0].arr->cdt_val.array_val.arr[i].type = metaffi_float32_array_type;
		cdts_param_ret[0].arr->cdt_val.array_val.arr[i].cdt_val.array_val.fixed_dimensions = 2;
		cdts_param_ret[0].arr->cdt_val.array_val.arr[i].cdt_val.array_val.length = 3;
		cdts_param_ret[0].arr->cdt_val.array_val.arr[i].cdt_val.array_val.arr = (struct cdt*)malloc(3 * sizeof(struct cdt));

		for (int j = 0; j < 3; j++) {
			cdts_param_ret[0].arr->cdt_val.array_val.arr[i].cdt_val.array_val.arr[j].type = metaffi_float32_array_type;
			cdts_param_ret[0].arr->cdt_val.array_val.arr[i].cdt_val.array_val.arr[j].cdt_val.array_val.fixed_dimensions = 1;
			cdts_param_ret[0].arr->cdt_val.array_val.arr[i].cdt_val.array_val.arr[j].cdt_val.array_val.length = 3;
			cdts_param_ret[0].arr->cdt_val.array_val.arr[i].cdt_val.array_val.arr[j].cdt_val.array_val.arr = (struct cdt*)malloc(3 * sizeof(struct cdt));

			for (int k = 0; k < 3; k++) {
				cdts_param_ret[0].arr->cdt_val.array_val.arr[i].cdt_val.array_val.arr[j].cdt_val.array_val.arr[k].type = metaffi_float32_type;
				cdts_param_ret[0].arr->cdt_val.array_val.arr[i].cdt_val.array_val.arr[j].cdt_val.array_val.arr[k].cdt_val.float32_val = (i * 9.0f) + (j * 3.0f) + k + 1.0f;
			}
		}
	}

	return cdts_param_ret;
}


*/
import "C"

func GetCDTS() *C.struct_cdts {
	return C.get_cdts()
}

func GetCDTSType(pcdts *C.struct_cdts, index int) uint64 {
	return uint64(C.get_cdts_type(pcdts, C.int(index)))
}

func Get2DUInt8ArrayCDTS() *C.struct_cdts {
	return C.get_2d_uint8_array_cdts()
}

func Get3DFloat32ArrayCDTS() *C.struct_cdts {
	return C.get_3d_float32_array_cdts()
}
