package metaffi

/*
#cgo !windows LDFLAGS: -L. -ldl
#cgo LDFLAGS: -Wl,--allow-multiple-definition

#include <include/xllr_capi_loader.h>
#include <include/xllr_capi_loader.c>

void call_plugin_xcall_no_params_no_ret(void** ppv, char** err, uint64_t* out_err)
{
	void* pvoidxcall = ppv[0];
	void* pctxt = ppv[1];

	(((void(*)(void*,char**,uint64_t*))pvoidxcall)(pctxt, err, out_err));
}

void call_plugin_xcall_no_params_ret(void** ppv, struct cdts* cdts, char** err, uint64_t* out_err)
{
	void* pvoidxcall = ppv[0];
	void* pctxt = ppv[1];

	(((void(*)(void*,void*,char**,uint64_t*))pvoidxcall)(pctxt, cdts, err, out_err));
}

void call_plugin_xcall_params_no_ret(void** ppv, struct cdts* cdts, char** err, uint64_t* out_err)
{
	void* pvoidxcall = ppv[0];
	void* pctxt = ppv[1];

	(((void(*)(void*,void*,char**,uint64_t*))pvoidxcall)(pctxt, cdts, err, out_err));
}

void call_plugin_xcall_params_ret(void** ppv, struct cdts* cdts, char** err, uint64_t* out_err)
{
	void* pvoidxcall = ppv[0];
	void* pctxt = ppv[1];

	(((void(*)(void*,void*,char**,uint64_t*))pvoidxcall)(pctxt, cdts, err, out_err));
}


struct cdts* cast_to_cdts(void* p)
{
	return (struct cdts*)p;
}

*/
import "C"
import (
	"fmt"
	"github.com/MetaFFI/plugin-sdk/compiler/go/IDL"
	"os"
	"unsafe"
)

func init() {
	err := C.load_xllr()
	if err != nil {
		panic("Failed to load MetaFFI XLLR functions: " + C.GoString(err))
	}
}

func XLLRLoadFunction(runtimePlugin string, modulePath string, functionPath string, paramsTypes []uint64, retvalsTypes []uint64) (*unsafe.Pointer, error) {

	var params []IDL.MetaFFITypeInfo
	if paramsTypes != nil {
		params = make([]IDL.MetaFFITypeInfo, 0)
		for _, p := range paramsTypes {
			params = append(params, IDL.MetaFFITypeInfo{Type: p})
		}
	}

	var retvals []IDL.MetaFFITypeInfo
	if retvalsTypes != nil {
		retvals = make([]IDL.MetaFFITypeInfo, 0)
		for _, r := range retvalsTypes {
			retvals = append(retvals, IDL.MetaFFITypeInfo{Type: r})
		}
	}

	return XLLRLoadFunctionWithAliases(runtimePlugin, modulePath, functionPath, params, retvals)
}

func XLLRLoadFunctionWithAliases(runtimePlugin string, modulePath string, functionPath string, paramsTypes []IDL.MetaFFITypeInfo, retvalsTypes []IDL.MetaFFITypeInfo) (*unsafe.Pointer, error) {

	pruntimePlugin := C.CString(runtimePlugin)
	defer CFree(unsafe.Pointer(pruntimePlugin))

	pmodulePath := C.CString(modulePath)
	defer CFree(unsafe.Pointer(pmodulePath))

	ppath := C.CString(functionPath)
	defer CFree(unsafe.Pointer(ppath))

	var out_err *C.char
	var out_err_len C.uint32_t
	out_err_len = C.uint32_t(0)

	var pparamTypes *C.struct_metaffi_type_info
	if paramsTypes != nil {
		pparamTypes = createMetaFFITypeInfoArray(paramsTypes)
		defer freeMetaFFITypeInfoArray(pparamTypes, len(paramsTypes))
	}

	pparamTypesLen := (C.uint8_t)(len(paramsTypes))

	var ppretvalsTypes *C.struct_metaffi_type_info
	if retvalsTypes != nil {
		ppretvalsTypes = createMetaFFITypeInfoArray(retvalsTypes)
		defer freeMetaFFITypeInfoArray(ppretvalsTypes, len(retvalsTypes))
	}
	pretvalsTypesLen := (C.uint8_t)(len(retvalsTypes))

	id := C.xllr_load_function(pruntimePlugin, C.uint(len(runtimePlugin)),
		pmodulePath, C.uint(len(modulePath)),
		ppath, C.uint(len(functionPath)),
		pparamTypes, ppretvalsTypes,
		pparamTypesLen, pretvalsTypesLen,
		&out_err, &out_err_len)

	if id == nil {
		return nil, fmt.Errorf("Failed to load foreign entity entrypoint \"%v\": %v", functionPath, string(C.GoBytes(unsafe.Pointer(out_err), C.int(out_err_len))))
	}

	return id, nil
}

func XLLRXCallParamsRet(pff *unsafe.Pointer, parameters unsafe.Pointer) error {

	// TODO: Free error message, in case of returned error
	// 		 The problem is that some plugins return strings that cannot be freed - FIX THIS!

	var out_err *C.char
	var out_err_len C.uint64_t
	out_err_len = C.uint64_t(0)

	C.call_plugin_xcall_params_ret(pff, C.cast_to_cdts(parameters), &out_err, &out_err_len)

	if out_err_len != C.uint64_t(0) {
		return fmt.Errorf("%v", string(C.GoBytes(unsafe.Pointer(out_err), C.int(out_err_len))))
	}

	return nil
}

func XLLRXCallNoParamsRet(pff *unsafe.Pointer, return_values unsafe.Pointer) error {

	var out_err *C.char
	var out_err_len C.uint64_t
	out_err_len = C.uint64_t(0)

	C.call_plugin_xcall_no_params_ret(pff, C.cast_to_cdts(return_values), &out_err, &out_err_len)

	if out_err_len != C.uint64_t(0) {
		return fmt.Errorf("%v", string(C.GoBytes(unsafe.Pointer(out_err), C.int(out_err_len))))
	}

	return nil
}

func XLLRXCallParamsNoRet(pff *unsafe.Pointer, parameters unsafe.Pointer) error {

	var out_err *C.char
	var out_err_len C.uint64_t
	out_err_len = C.uint64_t(0)

	C.call_plugin_xcall_params_no_ret(pff, C.cast_to_cdts(parameters), &out_err, &out_err_len)

	if out_err_len != C.uint64_t(0) {
		return fmt.Errorf("%v", string(C.GoBytes(unsafe.Pointer(out_err), C.int(out_err_len))))
	}

	return nil
}

func XLLRXCallNoParamsNoRet(pff *unsafe.Pointer) error {

	var out_err *C.char
	var out_err_len C.uint64_t
	out_err_len = C.uint64_t(0)

	C.call_plugin_xcall_no_params_no_ret(pff, &out_err, &out_err_len)

	if out_err_len != C.uint64_t(0) {
		return fmt.Errorf("%v", string(C.GoBytes(unsafe.Pointer(out_err), C.int(out_err_len))))
	}

	return nil
}

func XLLRLoadRuntimePlugin(runtimePlugin string) error {

	pruntime_plugin := C.CString(runtimePlugin)
	defer CFree(unsafe.Pointer(pruntime_plugin))

	// load foreign runtime
	var out_err *C.char
	var out_err_len C.uint32_t
	out_err_len = C.uint32_t(0)

	C.xllr_load_runtime_plugin(pruntime_plugin, C.uint(len(runtimePlugin)), &out_err, &out_err_len)

	if out_err_len != C.uint32_t(0) {
		return fmt.Errorf("Failed to load runtime %v: %v", runtimePlugin, string(C.GoBytes(unsafe.Pointer(out_err), C.int(out_err_len))))
	}

	return nil
}

func XLLRFreeRuntimePlugin(runtimePlugin string) error {

	pruntime_plugin := C.CString(runtimePlugin)
	defer CFree(unsafe.Pointer(pruntime_plugin))

	var out_err *C.char
	var out_err_len C.uint32_t
	out_err_len = C.uint32_t(0)

	C.xllr_free_runtime_plugin(pruntime_plugin, C.uint(len(runtimePlugin)), &out_err, &out_err_len)

	if out_err_len != C.uint32_t(0) {
		return fmt.Errorf("Failed to free runtime %v: %v", runtimePlugin, string(C.GoBytes(unsafe.Pointer(out_err), C.int(out_err_len))))
	}

	return nil
}

func ConstructCDTS(cdts *C.struct_cdts, callbacks *C.struct_construct_cdts_callbacks) {

	var err *C.char = nil
	C.xllr_construct_cdts(cdts, callbacks, &err)

	if err != nil {
		panic(C.GoString(err))
	}
}

func ConstructCDT(cdt *C.struct_cdt, callbacks *C.struct_construct_cdts_callbacks) {
	var err *C.char = nil
	fmt.Fprintf(os.Stderr, "ConstructCDT 1 - +++\n")
	C.xllr_construct_cdt(nil, nil, &err)
	fmt.Fprintf(os.Stderr, "ConstructCDT 2 - %v\n", err)
	if err != nil {
		panic(C.GoString(err))
	}
	fmt.Fprintf(os.Stderr, "ConstructCDT 3\n")
}

func TraverseCDTS(cdts *C.struct_cdts, callbacks *C.struct_traverse_cdts_callbacks) {
	var err *C.char = nil
	C.xllr_traverse_cdts(cdts, callbacks, &err)

	if err != nil {
		panic(C.GoString(err))
	}
}

func TraverseCDT(cdt *C.struct_cdt, callbacks *C.struct_traverse_cdts_callbacks) {
	var err *C.char = nil
	C.xllr_traverse_cdt(cdt, callbacks, &err)

	if err != nil {
		panic(C.GoString(err))
	}
}

func CFree(p unsafe.Pointer) {
	C.free(p)
}
