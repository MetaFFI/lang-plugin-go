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
