package metaffi

/*
#cgo !windows LDFLAGS: -L. -ldl
#cgo LDFLAGS: -Wl,--allow-multiple-definition

#include <include/cdt.h>
#include <include/metaffi_primitives.h>
#include <stdlib.h>

metaffi_type get_cdt_type(struct cdt* c) {
    return c->type;
}

void set_cdt_type(struct cdt* c, metaffi_type t) {
    c->type = t;
}

metaffi_float32 get_cdt_float32_val(struct cdt* c) {
    return c->cdt_val.float32_val;
}

void set_cdt_float32_val(struct cdt* c, metaffi_float32 val) {
    c->cdt_val.float32_val = val;
}

metaffi_float64 get_cdt_float64_val(struct cdt* c) {
    return c->cdt_val.float64_val;
}

void set_cdt_float64_val(struct cdt* c, metaffi_float64 val) {
    c->cdt_val.float64_val = val;
}

metaffi_bool get_cdt_bool_val(struct cdt* c) {
    return c->cdt_val.bool_val;
}

void set_cdt_bool_val(struct cdt* c, metaffi_bool val) {
    c->cdt_val.bool_val = val;
}

struct metaffi_char8* get_cdt_char8_val(struct cdt* c) {
    return &c->cdt_val.char8_val;
}

void set_cdt_char8_val(struct cdt* c, struct metaffi_char8* val) {
    c->cdt_val.char8_val = *val;
}

char* get_cdt_string8_val(struct cdt* c) {
	return (char*)c->cdt_val.string8_val;
}

void set_cdt_string8_val(struct cdt* c, char* val) {
	c->cdt_val.string8_val = val;
}

struct metaffi_char16* get_cdt_char16_val(struct cdt* c) {
    return &c->cdt_val.char16_val;
}

void set_cdt_char16_val(struct cdt* c, struct metaffi_char16* val) {
    c->cdt_val.char16_val = *val;
}

metaffi_string16 get_cdt_string16_val(struct cdt* c) {
	return c->cdt_val.string16_val;
}

void set_cdt_string16_val(struct cdt* c, metaffi_string16 val) {
	c->cdt_val.string16_val = val;
}

struct metaffi_char32* get_cdt_char32_val(struct cdt* c) {
    return &c->cdt_val.char32_val;
}

void set_cdt_char32_val(struct cdt* c, struct metaffi_char32* val) {
    c->cdt_val.char32_val = *val;
}

struct cdt_metaffi_handle* get_cdt_handle_val(struct cdt* c) {
    return &c->cdt_val.handle_val;
}

void set_cdt_handle_val(struct cdt* c, struct cdt_metaffi_handle* val) {
    c->cdt_val.handle_val = *val;
}

struct cdt_metaffi_callable* get_cdt_callable_val(struct cdt* c) {
    return c->cdt_val.callable_val;
}

void set_cdt_callable_val(struct cdt* c, struct cdt_metaffi_callable* val) {
    c->cdt_val.callable_val = val;
}

metaffi_int8 get_cdt_int8_val(struct cdt* c) {
    return c->cdt_val.int8_val;
}

void set_cdt_int8_val(struct cdt* c, metaffi_int8 val) {
    c->cdt_val.int8_val = val;
}

metaffi_uint8 get_cdt_uint8_val(struct cdt* c) {
    return c->cdt_val.uint8_val;
}

void set_cdt_uint8_val(struct cdt* c, metaffi_uint8 val) {
    c->cdt_val.uint8_val = val;
}

metaffi_int16 get_cdt_int16_val(struct cdt* c) {
    return c->cdt_val.int16_val;
}

void set_cdt_int16_val(struct cdt* c, metaffi_int16 val) {
    c->cdt_val.int16_val = val;
}

metaffi_uint16 get_cdt_uint16_val(struct cdt* c) {
    return c->cdt_val.uint16_val;
}

void set_cdt_uint16_val(struct cdt* c, metaffi_uint16 val) {
    c->cdt_val.uint16_val = val;
}

metaffi_int32 get_cdt_int32_val(struct cdt* c) {
    return c->cdt_val.int32_val;
}

void set_cdt_int32_val(struct cdt* c, metaffi_int32 val) {
    c->cdt_val.int32_val = val;
}

metaffi_uint32 get_cdt_uint32_val(struct cdt* c) {
    return c->cdt_val.uint32_val;
}

void set_cdt_uint32_val(struct cdt* c, metaffi_uint32 val) {
    c->cdt_val.uint32_val = val;
}

metaffi_int64 get_cdt_int64_val(struct cdt* c) {
    return c->cdt_val.int64_val;
}

void set_cdt_int64_val(struct cdt* c, metaffi_int64 val) {
    c->cdt_val.int64_val = val;
}

metaffi_uint64 get_cdt_uint64_val(struct cdt* c) {
    return c->cdt_val.uint64_val;
}

void set_cdt_uint64_val(struct cdt* c, metaffi_uint64 val) {
    c->cdt_val.uint64_val = val;
}

struct cdt* get_cdt_at_index(struct cdts* pcdts, int val) {
	return &pcdts->arr[val];
}
*/
import "C"
import "unsafe"

type CDTS struct {
	c *C.struct_cdts
}

func NewCDTSFromCDTS(c unsafe.Pointer) *CDTS {
	return &CDTS{c: (*C.struct_cdts)(c)}
}

func (cdts *CDTS) GetCDT(index int) *CDT {
	return &CDT{c: C.get_cdt_at_index(cdts.c, C.int(index))}
}

func (cdts *CDTS) GetLength() C.metaffi_size {
	return cdts.c.length
}

func (cdts *CDTS) GetFixedDimensions() C.metaffi_int64 {
	return cdts.c.fixed_dimensions
}

//------------------------------------------------------------

type CDT struct {
	c *C.struct_cdt
}

func (cdt *CDT) GetTypeVal() C.metaffi_type {
	return C.get_cdt_type(cdt.c)
}

func (cdt *CDT) SetTypeVal(t C.metaffi_type) {
	C.set_cdt_type(cdt.c, t)
}

func (cdt *CDT) GetFreeRequired() C.metaffi_bool {
	return cdt.c.free_required
}

func (cdt *CDT) GetFloat32Val() C.metaffi_float32 {
	return C.get_cdt_float32_val(cdt.c)
}

func (cdt *CDT) SetFloat32Val(val C.metaffi_float32) {
	C.set_cdt_float32_val(cdt.c, val)
}

func (cdt *CDT) GetFloat64Val() C.metaffi_float64 {
	return C.get_cdt_float64_val(cdt.c)
}

func (cdt *CDT) SetFloat64Val(val C.metaffi_float64) {
	C.set_cdt_float64_val(cdt.c, val)
}

func (cdt *CDT) GetBoolVal() C.metaffi_bool {
	return C.get_cdt_bool_val(cdt.c)
}

func (cdt *CDT) SetBoolVal(val C.metaffi_bool) {
	C.set_cdt_bool_val(cdt.c, val)
}

func (cdt *CDT) GetInt8Val() C.metaffi_int8 {
	return C.get_cdt_int8_val(cdt.c)
}

func (cdt *CDT) SetInt8Val(val C.metaffi_int8) {
	C.set_cdt_int8_val(cdt.c, val)
}

func (cdt *CDT) GetUInt8Val() C.metaffi_uint8 {
	return C.get_cdt_uint8_val(cdt.c)
}

func (cdt *CDT) SetUInt8Val(val C.metaffi_uint8) {
	C.set_cdt_uint8_val(cdt.c, val)
}

func (cdt *CDT) GetInt16Val() C.metaffi_int16 {
	return C.get_cdt_int16_val(cdt.c)
}

func (cdt *CDT) SetInt16Val(val C.metaffi_int16) {
	C.set_cdt_int16_val(cdt.c, val)
}

func (cdt *CDT) GetUInt16Val() C.metaffi_uint16 {
	return C.get_cdt_uint16_val(cdt.c)
}

func (cdt *CDT) SetUInt16Val(val C.metaffi_uint16) {
	C.set_cdt_uint16_val(cdt.c, val)
}

func (cdt *CDT) GetInt32Val() C.metaffi_int32 {
	return C.get_cdt_int32_val(cdt.c)
}

func (cdt *CDT) SetInt32Val(val C.metaffi_int32) {
	C.set_cdt_int32_val(cdt.c, val)
}

func (cdt *CDT) GetUInt32Val() C.metaffi_uint32 {
	return C.get_cdt_uint32_val(cdt.c)
}

func (cdt *CDT) SetUInt32Val(val C.metaffi_uint32) {
	C.set_cdt_uint32_val(cdt.c, val)
}

func (cdt *CDT) GetInt64Val() C.metaffi_int64 {
	return C.get_cdt_int64_val(cdt.c)
}

func (cdt *CDT) SetInt64Val(val C.metaffi_int64) {
	C.set_cdt_int64_val(cdt.c, val)
}

func (cdt *CDT) GetUInt64Val() C.metaffi_uint64 {
	return C.get_cdt_uint64_val(cdt.c)
}

func (cdt *CDT) SetUInt64Val(val C.metaffi_uint64) {
	C.set_cdt_uint64_val(cdt.c, val)
}

func (cdt *CDT) GetFloat32() float32 {
	return float32(C.get_cdt_float32_val(cdt.c))
}

func (cdt *CDT) SetFloat32(val float32) {
	C.set_cdt_float32_val(cdt.c, C.metaffi_float32(val))
}

func (cdt *CDT) GetFloat64() float64 {
	return float64(C.get_cdt_float64_val(cdt.c))
}

func (cdt *CDT) SetFloat64(val float64) {
	C.set_cdt_float64_val(cdt.c, C.metaffi_float64(val))
}

func (cdt *CDT) GetInt8() int8 {
	return int8(C.get_cdt_int8_val(cdt.c))
}

func (cdt *CDT) SetInt8(val int8) {
	C.set_cdt_int8_val(cdt.c, C.metaffi_int8(val))
}

func (cdt *CDT) GetUInt16() uint16 {
	return uint16(C.get_cdt_uint16_val(cdt.c))
}

func (cdt *CDT) SetUInt16(val uint16) {
	C.set_cdt_uint16_val(cdt.c, C.metaffi_uint16(val))
}

func (cdt *CDT) GetInt32() int32 {
	return int32(C.get_cdt_int32_val(cdt.c))
}

func (cdt *CDT) SetInt32(val int32) {
	C.set_cdt_int32_val(cdt.c, C.metaffi_int32(val))
}

func (cdt *CDT) GetUInt32() uint32 {
	return uint32(C.get_cdt_uint32_val(cdt.c))
}

func (cdt *CDT) SetUInt32(val uint32) {
	C.set_cdt_uint32_val(cdt.c, C.metaffi_uint32(val))
}

func (cdt *CDT) GetInt64() int64 {
	return int64(C.get_cdt_int64_val(cdt.c))
}

func (cdt *CDT) SetInt64(val int64) {
	C.set_cdt_int64_val(cdt.c, C.metaffi_int64(val))
}

func (cdt *CDT) GetUInt64() uint64 {
	return uint64(C.get_cdt_uint64_val(cdt.c))
}

func (cdt *CDT) SetUInt64(val uint64) {
	C.set_cdt_uint64_val(cdt.c, C.metaffi_uint64(val))
}

func (cdt *CDT) GetBool() bool {
	return C.get_cdt_bool_val(cdt.c) != 0
}

func (cdt *CDT) SetBool(val bool) {
	var cVal C.metaffi_bool
	if val {
		cVal = 1
	} else {
		cVal = 0
	}
	C.set_cdt_bool_val(cdt.c, cVal)
}

func (cdt *CDT) GetChar8() *MetaFFIChar8 {
	val := C.get_cdt_char8_val(cdt.c)
	return &MetaFFIChar8{Val: val}
}

func (cdt *CDT) SetChar8(val *MetaFFIChar8) {
	C.set_cdt_char8_val(cdt.c, val.Val)
}

func (cdt *CDT) GetChar16() *MetaFFIChar16 {
	val := C.get_cdt_char16_val(cdt.c)
	return &MetaFFIChar16{Val: val}
}

func (cdt *CDT) SetChar16(val *MetaFFIChar16) {
	C.set_cdt_char16_val(cdt.c, val.Val)
}

func (cdt *CDT) GetChar32() *MetaFFIChar32 {
	val := C.get_cdt_char32_val(cdt.c)
	return &MetaFFIChar32{Val: val}
}

func (cdt *CDT) SetChar32(val *MetaFFIChar32) {
	C.set_cdt_char32_val(cdt.c, val.Val)
}

func (cdt *CDT) GetString8() string {
	return C.GoString(C.get_cdt_string8_val(cdt.c))
}

func (cdt *CDT) SetString8(val string) {
	cVal := C.CString(val)
	C.set_cdt_string8_val(cdt.c, cVal)
}

func (cdt *CDT) GetHandleVal() *CDTMetaFFIHandle {
	val := C.get_cdt_handle_val(cdt.c)
	return &CDTMetaFFIHandle{Val: val}
}

func (cdt *CDT) SetHandleVal(val *CDTMetaFFIHandle) {
	C.set_cdt_handle_val(cdt.c, val.Val)
}

func (cdt *CDT) GetCallableVal() *MetaFFICallable {
	val := C.get_cdt_callable_val(cdt.c)
	return &MetaFFICallable{Val: val}
}

func (cdt *CDT) SetCallableVal(val *MetaFFICallable) {
	C.set_cdt_callable_val(cdt.c, val.Val)
}

//------------------------------------------------------------

type MetaFFIChar8 struct {
	Val *C.struct_metaffi_char8
}

//------------------------------------------------------------

type MetaFFIChar16 struct {
	Val *C.struct_metaffi_char16
}

//------------------------------------------------------------

type MetaFFIChar32 struct {
	Val *C.struct_metaffi_char32
}

//------------------------------------------------------------

type CDTMetaFFIHandle struct {
	Val *C.struct_cdt_metaffi_handle
}

func (handle *CDTMetaFFIHandle) GetHandle() C.metaffi_handle {
	return handle.Val.val
}

//------------------------------------------------------------

type MetaFFICallable struct {
	Val *C.struct_cdt_metaffi_callable
}

//--------------------------------------------------------------------
