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
	cdts_param_ret[0].pcdt->cdt_val.metaffi_uint8_array_val.dimensions = 2;
	cdts_param_ret[0].pcdt->cdt_val.metaffi_uint8_array_val.dimensions_lengths = malloc(2 * sizeof(metaffi_size));
	cdts_param_ret[0].pcdt->cdt_val.metaffi_uint8_array_val.dimensions_lengths[0] = 3;
	cdts_param_ret[0].pcdt->cdt_val.metaffi_uint8_array_val.dimensions_lengths[1] = 3;
	uint8_t data[3][3] = { {0,1,2}, {3,4,5}, {6,7,8} };
	cdts_param_ret[0].pcdt->cdt_val.metaffi_uint8_array_val.vals = (uint8_t*)malloc(3 * sizeof(uint8_t*));
	for(int i = 0; i < 3; i++)
	{
		((uint8_t**)cdts_param_ret[0].pcdt->cdt_val.metaffi_uint8_array_val.vals)[i] = (uint8_t*)malloc(3 * sizeof(uint8_t));
		for(int j = 0; j < 3; j++)
		{
			((uint8_t**)cdts_param_ret[0].pcdt->cdt_val.metaffi_uint8_array_val.vals)[i][j] = data[i][j];
		}
	}
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
