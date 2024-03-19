package metaffi

/*
#cgo !windows LDFLAGS: -L. -ldl
#cgo LDFLAGS: -Wl,--allow-multiple-definition

#include <stdlib.h>
#include <stdint.h>
#include <include/cdt_structs.h>
#include <include/cdt_capi_loader.h>
#include <include/cdt_capi_loader.c>

struct cdts* get_cdts()
{
	cdts* res = (cdts*)malloc(sizeof(struct cdts));
	res->pcdt = (cdt*)malloc(sizeof(struct cdt));
	res->len = 1;
}

uint64_t get_cdts_type(struct cdts* pcdts, int index)
{
	return pcdts->pcdt[index].type;
}

void get_2d_int8(struct cdts* pcdts)
{
	printf("++++ HERE ++++\n");

	metaffi_int8** int8_array = (metaffi_int8**)pcdts->pcdt[0].cdt_val.metaffi_int8_array_val.vals;

	for (int i = 0; i < 3; i++)
	{
		for (int j = 0; j < 3; j++)
		{
			printf("||| %d\n", int8_array[i][j]);
		}
	}
}

struct cdts* get_2d_uint8_array_cdts()
{
	cdts* cdts_param_ret = get_cdts();
	cdts_param_ret[0].pcdt->type = metaffi_uint8_array_type;
	cdts_param_ret[0].pcdt->cdt_val.metaffi_uint8_array_val.dimension = 2;
	cdts_param_ret[0].pcdt->cdt_val.metaffi_uint8_array_val.length = 3;
	cdts_param_ret[0].pcdt->cdt_val.metaffi_uint8_array_val.arr = (struct cdt_metaffi_uint8_array*)malloc(3 * sizeof(struct cdt_metaffi_uint8_array));
	cdts_param_ret[0].pcdt->cdt_val.metaffi_uint8_array_val.arr[0].dimension = 1;
	cdts_param_ret[0].pcdt->cdt_val.metaffi_uint8_array_val.arr[0].length = 3;
	cdts_param_ret[0].pcdt->cdt_val.metaffi_uint8_array_val.arr[0].vals = (uint8_t*)malloc(3 * sizeof(uint8_t));
	cdts_param_ret[0].pcdt->cdt_val.metaffi_uint8_array_val.arr[0].vals[0] = 0;
	cdts_param_ret[0].pcdt->cdt_val.metaffi_uint8_array_val.arr[0].vals[1] = 1;
	cdts_param_ret[0].pcdt->cdt_val.metaffi_uint8_array_val.arr[0].vals[2] = 2;

	cdts_param_ret[0].pcdt->cdt_val.metaffi_uint8_array_val.arr[1].dimension = 1;
	cdts_param_ret[0].pcdt->cdt_val.metaffi_uint8_array_val.arr[1].length = 3;
	cdts_param_ret[0].pcdt->cdt_val.metaffi_uint8_array_val.arr[1].vals = (uint8_t*)malloc(3 * sizeof(uint8_t));
	cdts_param_ret[0].pcdt->cdt_val.metaffi_uint8_array_val.arr[1].vals[0] = 3;
	cdts_param_ret[0].pcdt->cdt_val.metaffi_uint8_array_val.arr[1].vals[1] = 4;
	cdts_param_ret[0].pcdt->cdt_val.metaffi_uint8_array_val.arr[1].vals[2] = 5;

	cdts_param_ret[0].pcdt->cdt_val.metaffi_uint8_array_val.arr[2].dimension = 1;
	cdts_param_ret[0].pcdt->cdt_val.metaffi_uint8_array_val.arr[2].length = 3;
	cdts_param_ret[0].pcdt->cdt_val.metaffi_uint8_array_val.arr[2].vals = (uint8_t*)malloc(3 * sizeof(uint8_t));
	cdts_param_ret[0].pcdt->cdt_val.metaffi_uint8_array_val.arr[2].vals[0] = 6;
	cdts_param_ret[0].pcdt->cdt_val.metaffi_uint8_array_val.arr[2].vals[1] = 7;
	cdts_param_ret[0].pcdt->cdt_val.metaffi_uint8_array_val.arr[2].vals[2] = 8;

	return cdts_param_ret;
}

struct cdts* get_3d_float32_array_cdts()
{
	cdts* cdts_param_ret = get_cdts();
	cdts_param_ret[0].pcdt->type = metaffi_float32_array_type;
	cdts_param_ret[0].pcdt->cdt_val.metaffi_float32_array_val.dimension = 3;
	cdts_param_ret[0].pcdt->cdt_val.metaffi_float32_array_val.length = 3;
	cdts_param_ret[0].pcdt->cdt_val.metaffi_float32_array_val.arr = (struct cdt_metaffi_float32_array*)malloc(3 * sizeof(struct cdt_metaffi_float32_array));

	cdts_param_ret[0].pcdt->cdt_val.metaffi_float32_array_val.arr[0].dimension = 2;
	cdts_param_ret[0].pcdt->cdt_val.metaffi_float32_array_val.arr[0].length = 3;
	cdts_param_ret[0].pcdt->cdt_val.metaffi_float32_array_val.arr[0].arr = (struct cdt_metaffi_float32_array*)malloc(3 * sizeof(struct cdt_metaffi_float32_array));

	cdts_param_ret[0].pcdt->cdt_val.metaffi_float32_array_val.arr[0].arr[0].dimension = 1;
	cdts_param_ret[0].pcdt->cdt_val.metaffi_float32_array_val.arr[0].arr[0].length = 3;
	cdts_param_ret[0].pcdt->cdt_val.metaffi_float32_array_val.arr[0].arr[0].vals = (float*)malloc(3 * sizeof(float));
	cdts_param_ret[0].pcdt->cdt_val.metaffi_float32_array_val.arr[0].arr[0].vals[0] = 1.0f;
	cdts_param_ret[0].pcdt->cdt_val.metaffi_float32_array_val.arr[0].arr[0].vals[1] = 2.0f;
	cdts_param_ret[0].pcdt->cdt_val.metaffi_float32_array_val.arr[0].arr[0].vals[2] = 3.0f;

	cdts_param_ret[0].pcdt->cdt_val.metaffi_float32_array_val.arr[0].arr[1].dimension = 1;
	cdts_param_ret[0].pcdt->cdt_val.metaffi_float32_array_val.arr[0].arr[1].length = 3;
	cdts_param_ret[0].pcdt->cdt_val.metaffi_float32_array_val.arr[0].arr[1].vals = (float*)malloc(3 * sizeof(float));
	cdts_param_ret[0].pcdt->cdt_val.metaffi_float32_array_val.arr[0].arr[1].vals[0] = 4.0f;
	cdts_param_ret[0].pcdt->cdt_val.metaffi_float32_array_val.arr[0].arr[1].vals[1] = 5.0f;
	cdts_param_ret[0].pcdt->cdt_val.metaffi_float32_array_val.arr[0].arr[1].vals[2] = 6.0f;

	cdts_param_ret[0].pcdt->cdt_val.metaffi_float32_array_val.arr[0].arr[2].dimension = 1;
	cdts_param_ret[0].pcdt->cdt_val.metaffi_float32_array_val.arr[0].arr[2].length = 3;
	cdts_param_ret[0].pcdt->cdt_val.metaffi_float32_array_val.arr[0].arr[2].vals = (float*)malloc(3 * sizeof(float));
	cdts_param_ret[0].pcdt->cdt_val.metaffi_float32_array_val.arr[0].arr[2].vals[0] = 7.0f;
	cdts_param_ret[0].pcdt->cdt_val.metaffi_float32_array_val.arr[0].arr[2].vals[1] = 8.0f;
	cdts_param_ret[0].pcdt->cdt_val.metaffi_float32_array_val.arr[0].arr[2].vals[2] = 9.0f;

	cdts_param_ret[0].pcdt->cdt_val.metaffi_float32_array_val.arr[1].dimension = 2;
	cdts_param_ret[0].pcdt->cdt_val.metaffi_float32_array_val.arr[1].length = 3;
	cdts_param_ret[0].pcdt->cdt_val.metaffi_float32_array_val.arr[1].arr = (struct cdt_metaffi_float32_array*)malloc(3 * sizeof(struct cdt_metaffi_float32_array));

	cdts_param_ret[0].pcdt->cdt_val.metaffi_float32_array_val.arr[1].arr[0].dimension = 1;
	cdts_param_ret[0].pcdt->cdt_val.metaffi_float32_array_val.arr[1].arr[0].length = 3;
	cdts_param_ret[0].pcdt->cdt_val.metaffi_float32_array_val.arr[1].arr[0].vals = (float*)malloc(3 * sizeof(float));
	cdts_param_ret[0].pcdt->cdt_val.metaffi_float32_array_val.arr[1].arr[0].vals[0] = 10.0f;
	cdts_param_ret[0].pcdt->cdt_val.metaffi_float32_array_val.arr[1].arr[0].vals[1] = 11.0f;
	cdts_param_ret[0].pcdt->cdt_val.metaffi_float32_array_val.arr[1].arr[0].vals[2] = 12.0f;

	cdts_param_ret[0].pcdt->cdt_val.metaffi_float32_array_val.arr[1].arr[1].dimension = 1;
	cdts_param_ret[0].pcdt->cdt_val.metaffi_float32_array_val.arr[1].arr[1].length = 3;
	cdts_param_ret[0].pcdt->cdt_val.metaffi_float32_array_val.arr[1].arr[1].vals = (float*)malloc(3 * sizeof(float));
	cdts_param_ret[0].pcdt->cdt_val.metaffi_float32_array_val.arr[1].arr[1].vals[0] = 13.0f;
	cdts_param_ret[0].pcdt->cdt_val.metaffi_float32_array_val.arr[1].arr[1].vals[1] = 14.0f;
	cdts_param_ret[0].pcdt->cdt_val.metaffi_float32_array_val.arr[1].arr[1].vals[2] = 15.0f;

	cdts_param_ret[0].pcdt->cdt_val.metaffi_float32_array_val.arr[1].arr[2].dimension = 1;
	cdts_param_ret[0].pcdt->cdt_val.metaffi_float32_array_val.arr[1].arr[2].length = 3;
	cdts_param_ret[0].pcdt->cdt_val.metaffi_float32_array_val.arr[1].arr[2].vals = (float*)malloc(3 * sizeof(float));
	cdts_param_ret[0].pcdt->cdt_val.metaffi_float32_array_val.arr[1].arr[2].vals[0] = 16.0f;
	cdts_param_ret[0].pcdt->cdt_val.metaffi_float32_array_val.arr[1].arr[2].vals[1] = 17.0f;
	cdts_param_ret[0].pcdt->cdt_val.metaffi_float32_array_val.arr[1].arr[2].vals[2] = 18.0f;

	cdts_param_ret[0].pcdt->cdt_val.metaffi_float32_array_val.arr[2].dimension = 2;
	cdts_param_ret[0].pcdt->cdt_val.metaffi_float32_array_val.arr[2].length = 3;
	cdts_param_ret[0].pcdt->cdt_val.metaffi_float32_array_val.arr[2].arr = (struct cdt_metaffi_float32_array*)malloc(3 * sizeof(struct cdt_metaffi_float32_array));

	cdts_param_ret[0].pcdt->cdt_val.metaffi_float32_array_val.arr[2].arr[0].dimension = 1;
	cdts_param_ret[0].pcdt->cdt_val.metaffi_float32_array_val.arr[2].arr[0].length = 3;
	cdts_param_ret[0].pcdt->cdt_val.metaffi_float32_array_val.arr[2].arr[0].vals = (float*)malloc(3 * sizeof(float));
	cdts_param_ret[0].pcdt->cdt_val.metaffi_float32_array_val.arr[2].arr[0].vals[0] = 19.0f;
	cdts_param_ret[0].pcdt->cdt_val.metaffi_float32_array_val.arr[2].arr[0].vals[1] = 20.0f;
	cdts_param_ret[0].pcdt->cdt_val.metaffi_float32_array_val.arr[2].arr[0].vals[2] = 21.0f;

	cdts_param_ret[0].pcdt->cdt_val.metaffi_float32_array_val.arr[2].arr[1].dimension = 1;
	cdts_param_ret[0].pcdt->cdt_val.metaffi_float32_array_val.arr[2].arr[1].length = 3;
	cdts_param_ret[0].pcdt->cdt_val.metaffi_float32_array_val.arr[2].arr[1].vals = (float*)malloc(3 * sizeof(float));
	cdts_param_ret[0].pcdt->cdt_val.metaffi_float32_array_val.arr[2].arr[1].vals[0] = 22.0f;
	cdts_param_ret[0].pcdt->cdt_val.metaffi_float32_array_val.arr[2].arr[1].vals[1] = 23.0f;
	cdts_param_ret[0].pcdt->cdt_val.metaffi_float32_array_val.arr[2].arr[1].vals[2] = 24.0f;

	cdts_param_ret[0].pcdt->cdt_val.metaffi_float32_array_val.arr[2].arr[2].dimension = 1;
	cdts_param_ret[0].pcdt->cdt_val.metaffi_float32_array_val.arr[2].arr[2].length = 3;
	cdts_param_ret[0].pcdt->cdt_val.metaffi_float32_array_val.arr[2].arr[2].vals = (float*)malloc(3 * sizeof(float));
	cdts_param_ret[0].pcdt->cdt_val.metaffi_float32_array_val.arr[2].arr[2].vals[0] = 25.0f;
	cdts_param_ret[0].pcdt->cdt_val.metaffi_float32_array_val.arr[2].arr[2].vals[1] = 26.0f;
	cdts_param_ret[0].pcdt->cdt_val.metaffi_float32_array_val.arr[2].arr[2].vals[2] = 27.0f;

	return cdts_param_ret;
}


*/
import "C"

func GetCDTS() *C.cdts {
	return C.get_cdts()
}

func GetCDTSType(pcdts *C.cdts, index int) uint64 {
	return uint64(C.get_cdts_type(pcdts, C.int(index)))
}

func Print2DInt8(data *C.struct_cdts) {

	C.get_2d_int8(data)
}

func Get2DUInt8ArrayCDTS() *C.cdts {
	return C.get_2d_uint8_array_cdts()
}

func Get3DFloat32ArrayCDTS() *C.cdts {
	return C.get_3d_float32_array_cdts()
}
