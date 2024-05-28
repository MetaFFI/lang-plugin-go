package metaffi

/*
#cgo !windows LDFLAGS: -L. -ldl
#cgo LDFLAGS: -Wl,--allow-multiple-definition

#include <string.h>
#include <include/cdts_traverse_construct.h>
#include <include/xllr_capi_loader.h>
#include <stdlib.h>

void set_metaffi_type_info_type(struct metaffi_type_info* info, uint64_t type) {
    info->type = type;
}

metaffi_type get_metaffi_type(struct metaffi_type_info* info) {
    return info->type;
}

uint32_t* cast_char32_t_to_uint32_t(char32_t* input) {
    return (uint32_t*)input;
}

uint16_t* cast_char16_t_to_uint16_t(char16_t* input) {
    return (uint16_t*)input;
}

uint8_t* cast_char8_t_to_uint8_t(char8_t* input) {
    return (uint8_t*)input;
}

struct metaffi_type_info* cast_to_metaffi_type_info(void* input) {
    return (struct metaffi_type_info*)input;
}

metaffi_string8 cast_to_metaffi_string8(char* input) {
    return (metaffi_string8)input;
}

struct cdt_metaffi_handle get_null_handle(){
	struct cdt_metaffi_handle res;
	res.handle = NULL;
	res.runtime_id = 0;
	res.release = NULL;
	return res;
}

char* copy_string(char* s, int n) {
	char* cstr = (char*)malloc(n*sizeof(char) + 1);
	memcpy(cstr, s, n);
	cstr[n] = 0;
	return cstr;
}
*/
import "C"
import (
	"fmt"
	"github.com/MetaFFI/plugin-sdk/compiler/go/IDL"
	"golang.org/x/text/unicode/norm"
	"reflect"
	"unicode/utf16"
	"unicode/utf8"
	"unsafe"
)

//export onFloat64
func onFloat64(index *C.metaffi_size, indexSize C.metaffi_size, val C.metaffi_float64, context unsafe.Pointer) {
	tctxt := traverseContextTLS.Get()

	if indexSize == 0 {
		tctxt.Result = float64(val)
	} else { // within an array
		elem := getElement(index, indexSize, tctxt.Result)
		elem.SetFloat(float64(val))
	}
}

//export onFloat32
func onFloat32(index *C.metaffi_size, indexSize C.metaffi_size, val C.metaffi_float32, context unsafe.Pointer) {
	tctxt := traverseContextTLS.Get()

	if indexSize == 0 {
		tctxt.Result = float32(val)
	} else { // within an array
		elem := getElement(index, indexSize, tctxt.Result)
		elem.Set(reflect.ValueOf(float32(val)))
	}
}

//export onInt8
func onInt8(index *C.metaffi_size, indexSize C.metaffi_size, val C.metaffi_int8, context unsafe.Pointer) {
	tctxt := traverseContextTLS.Get()

	if indexSize == 0 {
		tctxt.Result = int8(val)
	} else { // within an array
		elem := getElement(index, indexSize, tctxt.Result)
		elem.Set(reflect.ValueOf(int8(val)))
	}
}

//export onUInt8
func onUInt8(index *C.metaffi_size, indexSize C.metaffi_size, val C.metaffi_uint8, context unsafe.Pointer) {
	tctxt := traverseContextTLS.Get()

	if indexSize == 0 {
		tctxt.Result = uint8(val)
	} else { // within an array
		elem := getElement(index, indexSize, tctxt.Result)
		elem.Set(reflect.ValueOf(uint8(val)))
	}
}

//export onInt16
func onInt16(index *C.metaffi_size, indexSize C.metaffi_size, val C.metaffi_int16, context unsafe.Pointer) {
	tctxt := traverseContextTLS.Get()

	if indexSize == 0 {
		tctxt.Result = int16(val)
	} else { // within an array
		elem := getElement(index, indexSize, tctxt.Result)
		elem.Set(reflect.ValueOf(int16(val)))
	}
}

//export onUInt16
func onUInt16(index *C.metaffi_size, indexSize C.metaffi_size, val C.metaffi_uint16, context unsafe.Pointer) {
	tctxt := traverseContextTLS.Get()

	if indexSize == 0 {
		tctxt.Result = uint16(val)
	} else { // within an array
		elem := getElement(index, indexSize, tctxt.Result)
		elem.Set(reflect.ValueOf(uint16(val)))
	}
}

//export onInt32
func onInt32(index *C.metaffi_size, indexSize C.metaffi_size, val C.metaffi_int32, context unsafe.Pointer) {
	tctxt := traverseContextTLS.Get()

	if indexSize == 0 {
		tctxt.Result = int32(val)
	} else { // within an array
		elem := getElement(index, indexSize, tctxt.Result)
		elem.Set(reflect.ValueOf(int32(val)))
	}
}

//export onUInt32
func onUInt32(index *C.metaffi_size, indexSize C.metaffi_size, val C.metaffi_uint32, context unsafe.Pointer) {
	tctxt := traverseContextTLS.Get()

	if indexSize == 0 {
		tctxt.Result = uint32(val)
	} else { // within an array
		elem := getElement(index, indexSize, tctxt.Result)
		elem.Set(reflect.ValueOf(uint32(val)))
	}
}

//export onInt64
func onInt64(index *C.metaffi_size, indexSize C.metaffi_size, val C.metaffi_int64, context unsafe.Pointer) {
	tctxt := traverseContextTLS.Get()

	if indexSize == 0 {
		tctxt.Result = int64(val)
	} else { // within an array
		elem := getElement(index, indexSize, tctxt.Result)
		elem.Set(reflect.ValueOf(int64(val)))
	}
}

//export onBool
func onBool(index *C.metaffi_size, indexSize C.metaffi_size, val C.metaffi_bool, context unsafe.Pointer) {
	tctxt := traverseContextTLS.Get()

	if indexSize == 0 {
		var goval bool
		if val == C.metaffi_bool(0) {
			goval = false
		} else {
			goval = true
		}

		tctxt.Result = goval
	} else { // within an array
		elem := getElement(index, indexSize, tctxt.Result)

		if val == C.metaffi_bool(0) {
			elem.Set(reflect.ValueOf(false))
		} else {
			elem.Set(reflect.ValueOf(true))
		}
	}
}

//export onUInt64
func onUInt64(index *C.metaffi_size, indexSize C.metaffi_size, val C.metaffi_uint64, context unsafe.Pointer) {
	tctxt := traverseContextTLS.Get()

	if indexSize == 0 {
		tctxt.Result = uint64(val)
	} else { // within an array
		elem := getElement(index, indexSize, tctxt.Result)
		elem.Set(reflect.ValueOf(uint64(val)))
	}
}

//export onHandle
func onHandle(index *C.metaffi_size, indexSize C.metaffi_size, val *C.struct_cdt_metaffi_handle, context unsafe.Pointer) {
	tctxt := traverseContextTLS.Get()

	if indexSize == 0 {
		// If not Go, return CDTMetaFFIHandle
		tctxt.Result = GetGoObject(val)

	} else { // within an array
		elem := getElement(index, indexSize, tctxt.Result)
		elem.Set(reflect.ValueOf(GetGoObject(val)))
	}
}

//export onCallable
func onCallable(index *C.metaffi_size, indexSize C.metaffi_size, val *C.struct_cdt_metaffi_callable, context unsafe.Pointer) {
	panic("Not supported yet")
}

//export onNull
func onNull(index *C.metaffi_size, indexSize C.metaffi_size, context unsafe.Pointer) {
	tctxt := traverseContextTLS.Get()

	if indexSize == 0 {
		// If not Go, return CDTMetaFFIHandle
		tctxt.Result = nil

	} else { // within an array
		elem := getElement(index, indexSize, tctxt.Result)
		elem.Set(reflect.ValueOf(nil))
	}
}

//export onArray
func onArray(index *C.metaffi_size, indexSize C.metaffi_size, val *C.struct_cdts, fixedDimensions C.metaffi_int64, commonType C.metaffi_type, context unsafe.Pointer) C.metaffi_bool {
	tctxt := traverseContextTLS.Get()

	if commonType&C.metaffi_type(C.metaffi_array_type) != 0 {
		commonType = commonType & ^C.metaffi_type(C.metaffi_array_type)
	}

	// if metaffi_any_type, get the common metaffi type
	var tempForAnyTempDynamicChecking C.metaffi_type = 0
	if commonType&C.metaffi_type(C.metaffi_any_type) != 0 {
		cdts := CDTS{c: val}
		for i := C.metaffi_size(0); i < cdts.GetLength(); i++ {
			elem := cdts.GetCDT(int(i))
			if tempForAnyTempDynamicChecking == 0 {
				tempForAnyTempDynamicChecking = elem.GetTypeVal()
			} else if tempForAnyTempDynamicChecking != elem.GetTypeVal() { // no common type - use metaffi_any_type
				commonType = C.metaffi_any_type
				break
			}
		}
	}

	var commonGoType reflect.Type = nil
	if commonType&C.metaffi_type(C.metaffi_handle_type) != 0 { // if metaffi_handle, get the Go common type

		cdts := CDTS{c: val}
		for i := C.metaffi_size(0); i < cdts.GetLength(); i++ {

			if cdts.GetCDT(int(i)).GetTypeVal()&C.metaffi_array_type == 0 {
				elem := GetGoObject(cdts.GetCDT(int(i)).GetHandleVal().Val)

				if elem == nil {
					panic(fmt.Sprintf("Go Object returned nil - Handle: %v %v", cdts.GetCDT(int(i)).GetHandleVal().Val.handle, cdts.GetCDT(int(i)).GetHandleVal().Val.runtime_id))
				}

				curType := reflect.ValueOf(elem).Type()
				if commonGoType == nil {
					commonGoType = curType
				} else if commonGoType != curType { // no common type - use interface{}
					commonGoType = reflect.TypeFor[interface{}]()
					break
				}
			} else {
				commonGoType = reflect.TypeFor[interface{}]()
			}
		}
	} else if commonType == C.metaffi_any_type {
		commonGoType = reflect.TypeFor[interface{}]()
	}

	if indexSize == 0 { // check roots

		// if handle, get the common type
		tctxt.Result = createMultiDimSlice(int(val.length), int(fixedDimensions), getGoTypeFromMetaFFIType(commonType, commonGoType))
	} else { // within an array
		elem := getElement(index, indexSize, tctxt.Result)
		elem.Set(reflect.ValueOf(createMultiDimSlice(int(val.length), int(fixedDimensions), getGoTypeFromMetaFFIType(commonType, commonGoType))))
	}

	return C.metaffi_bool(1)
}

//export onChar8
func onChar8(index *C.metaffi_size, indexSize C.metaffi_size, val C.struct_metaffi_char8, context unsafe.Pointer) {
	tctxt := traverseContextTLS.Get()

	// Convert C array to Go slice
	goSlice := []uint8{
		uint8(*C.cast_char8_t_to_uint8_t(&val.c[0])),
		uint8(*C.cast_char8_t_to_uint8_t(&val.c[1])),
		uint8(*C.cast_char8_t_to_uint8_t(&val.c[2])),
		uint8(*C.cast_char8_t_to_uint8_t(&val.c[3])),
	}
	decoded := string(goSlice)

	if indexSize == 0 {
		// If not Go, return CDTMetaFFIHandle
		tctxt.Result = decoded[0]
	} else { // within an array
		elem := getElement(index, indexSize, tctxt.Result)
		elem.Set(reflect.ValueOf(decoded[0]))
	}
}

//export onString8
func onString8(index *C.metaffi_size, indexSize C.metaffi_size, val *C.char, context unsafe.Pointer) {
	tctxt := traverseContextTLS.Get()

	if indexSize == 0 {
		// Convert C string to Go string
		goBytes := C.GoBytes(unsafe.Pointer(val), C.int(C.strlen(val)))
		goString := string(goBytes)

		// If not Go, return CDTMetaFFIHandle
		tctxt.Result = goString
	} else { // within an array
		elem := getElement(index, indexSize, tctxt.Result)

		goBytes := C.GoBytes(unsafe.Pointer(val), C.int(C.strlen(val)))
		goString := string(goBytes)

		elem.Set(reflect.ValueOf(goString))
	}
}

//export onChar16
func onChar16(index *C.metaffi_size, indexSize C.metaffi_size, val C.struct_metaffi_char16, context unsafe.Pointer) {
	tctxt := traverseContextTLS.Get()

	// Convert C array to Go slice
	goSlice := []uint16{
		uint16(*C.cast_char16_t_to_uint16_t(&val.c[0])),
		uint16(*C.cast_char16_t_to_uint16_t(&val.c[1])),
	}
	decoded := utf16.Decode(goSlice)

	if indexSize == 0 {
		// If not Go, return CDTMetaFFIHandle
		tctxt.Result = decoded[0]
	} else { // within an array
		elem := getElement(index, indexSize, tctxt.Result)
		elem.Set(reflect.ValueOf(decoded[0]))
	}
}

//export onString16
func onString16(index *C.metaffi_size, indexSize C.metaffi_size, val *C.char16_t, length C.int, context unsafe.Pointer) {
	tctxt := traverseContextTLS.Get()

	// Convert C array to Go slice
	cSlice := (*[1 << 30]C.char16_t)(unsafe.Pointer(val))[:length:length]

	// Convert C array to Go slice
	goSlice := make([]uint16, length)
	for i, c := range cSlice {
		goSlice[i] = uint16(c)
	}

	// Decode UTF-16 to Go string
	goString := string(utf16.Decode(goSlice))

	if indexSize == 0 {
		// If not Go, return CDTMetaFFIHandle
		tctxt.Result = goString
	} else { // within an array
		elem := getElement(index, indexSize, tctxt.Result)
		elem.Set(reflect.ValueOf(goString))
	}
}

//export onChar32
func onChar32(index *C.metaffi_size, indexSize C.metaffi_size, val C.struct_metaffi_char32, context unsafe.Pointer) {
	tctxt := traverseContextTLS.Get()

	// Convert C array to Go rune
	goRune := rune(*C.cast_char32_t_to_uint32_t(&val.c))

	// Convert rune to string
	decoded := string(goRune)

	if indexSize == 0 {
		// If not Go, return CDTMetaFFIHandle
		tctxt.Result = decoded
	} else { // within an array
		elem := getElement(index, indexSize, tctxt.Result)
		elem.Set(reflect.ValueOf(decoded))
	}
}

//export onString32
func onString32(index *C.metaffi_size, indexSize C.metaffi_size, val *C.char32_t, length C.int, context unsafe.Pointer) {
	tctxt := traverseContextTLS.Get()

	// Convert C array to Go slice
	cSlice := (*[1 << 30]C.char32_t)(unsafe.Pointer(val))[:length:length]

	// Convert C array to Go slice
	goSlice := make([]rune, length)
	for i, c := range cSlice {
		goSlice[i] = rune(c)
	}

	// Convert rune slice to string
	goString := string(goSlice)

	if indexSize == 0 {
		// If not Go, return CDTMetaFFIHandle
		tctxt.Result = goString
	} else { // within an array
		elem := getElement(index, indexSize, tctxt.Result)
		elem.Set(reflect.ValueOf(goString))
	}
}

//export constructCDTArray
func constructCDTArray(index *C.metaffi_size, indexSize C.metaffi_size, manuallyFillArray *C.struct_cdts, context unsafe.Pointer) {
	panic("Shouldn't be called")
}

//export getRootElementsCount
func getRootElementsCount(context unsafe.Pointer) C.metaffi_size {
	panic("As we parse CDT and not the root CDTS, this function should not be called")
}

//export getTypeInfo
func getTypeInfo(index *C.metaffi_size, indexSize C.metaffi_size, _ unsafe.Pointer) C.struct_metaffi_type_info {

	cctxt := constructContextTLS.Get()

	if index == nil { // root
		var mt C.struct_metaffi_type_info
		mt.is_free_alias = C.metaffi_bool(0)
		C.set_metaffi_type_info_type(&mt, C.uint64_t(cctxt.TypeInfo.Type))

		if C.get_metaffi_type(&mt) == C.metaffi_any_type {
			mffitype, _ := getMetaFFITypeFromGoType(reflect.ValueOf(cctxt.Input))
			C.set_metaffi_type_info_type(&mt, C.uint64_t(mffitype))
		}
		return mt
	} else {
		val := getElement(index, indexSize, cctxt.Input)

		detectedType, _ := getMetaFFITypeFromGoType(val)
		ti := IDL.MetaFFITypeInfo{Type: uint64(detectedType)}

		idlTypeInfo := ti.AsCMetaFFITypeInfo()
		cTypeInfo := C.cast_to_metaffi_type_info(unsafe.Pointer(&idlTypeInfo))

		res := *cTypeInfo
		return res
	}
}

//export getFloat64
func getFloat64(index *C.metaffi_size, indexSize C.metaffi_size, _ unsafe.Pointer) C.metaffi_float64 {
	cctxt := constructContextTLS.Get()
	val := getElement(index, indexSize, cctxt.Input)
	return C.metaffi_float64(val.Interface().(float64))
}

//export getFloat32
func getFloat32(index *C.metaffi_size, indexSize C.metaffi_size, _ unsafe.Pointer) C.metaffi_float32 {
	cctxt := constructContextTLS.Get()
	val := getElement(index, indexSize, cctxt.Input)
	return C.metaffi_float32(val.Interface().(float32))
}

//export getInt8
func getInt8(index *C.metaffi_size, indexSize C.metaffi_size, _ unsafe.Pointer) C.metaffi_int8 {
	cctxt := constructContextTLS.Get()
	val := getElement(index, indexSize, cctxt.Input)
	return C.metaffi_int8(val.Interface().(int8))
}

//export getUInt8
func getUInt8(index *C.metaffi_size, indexSize C.metaffi_size, _ unsafe.Pointer) C.metaffi_uint8 {
	cctxt := constructContextTLS.Get()
	val := getElement(index, indexSize, cctxt.Input)
	return C.metaffi_uint8(val.Interface().(uint8))
}

//export getInt16
func getInt16(index *C.metaffi_size, indexSize C.metaffi_size, _ unsafe.Pointer) C.metaffi_int16 {
	cctxt := constructContextTLS.Get()
	val := getElement(index, indexSize, cctxt.Input)
	return C.metaffi_int16(val.Interface().(int16))
}

//export getUInt16
func getUInt16(index *C.metaffi_size, indexSize C.metaffi_size, _ unsafe.Pointer) C.metaffi_uint16 {
	cctxt := constructContextTLS.Get()
	val := getElement(index, indexSize, cctxt.Input)
	return C.metaffi_uint16(val.Interface().(uint16))
}

//export getInt32
func getInt32(index *C.metaffi_size, indexSize C.metaffi_size, _ unsafe.Pointer) C.metaffi_int32 {
	cctxt := constructContextTLS.Get()
	val := getElement(index, indexSize, cctxt.Input)
	return C.metaffi_int32(val.Interface().(int32))
}

//export getUInt32
func getUInt32(index *C.metaffi_size, indexSize C.metaffi_size, _ unsafe.Pointer) C.metaffi_uint32 {
	cctxt := constructContextTLS.Get()
	val := getElement(index, indexSize, cctxt.Input)
	return C.metaffi_uint32(val.Interface().(uint32))
}

//export getInt64
func getInt64(index *C.metaffi_size, indexSize C.metaffi_size, _ unsafe.Pointer) C.metaffi_int64 {
	cctxt := constructContextTLS.Get()
	val := getElement(index, indexSize, cctxt.Input)

	// val might be an alias to int64
	return C.metaffi_int64(val.Int())
}

//export getUInt64
func getUInt64(index *C.metaffi_size, indexSize C.metaffi_size, _ unsafe.Pointer) C.metaffi_uint64 {
	cctxt := constructContextTLS.Get()
	val := getElement(index, indexSize, cctxt.Input)
	return C.metaffi_uint64(val.Interface().(uint64))
}

//export getBool
func getBool(index *C.metaffi_size, indexSize C.metaffi_size, _ unsafe.Pointer) C.metaffi_bool {
	cctxt := constructContextTLS.Get()
	val := getElement(index, indexSize, cctxt.Input)
	if val.Interface().(bool) {
		return C.metaffi_bool(1)
	} else {
		return C.metaffi_bool(0)
	}
}

//export getChar8
func getChar8(index *C.metaffi_size, indexSize C.metaffi_size, _ unsafe.Pointer) C.struct_metaffi_char8 {
	cctxt := constructContextTLS.Get()
	val := getElement(index, indexSize, cctxt.Input)
	r := val.Interface().(rune)

	// Convert the rune to a UTF-8 encoded byte array
	utf8Bytes := make([]byte, 4)
	utf8.EncodeRune(utf8Bytes, r)

	// Create a C.struct_metaffi_char8 and copy the bytes into it
	var char8 C.struct_metaffi_char8
	for i, b := range utf8Bytes {
		char8.c[i] = C.uchar(b) // change here
	}

	return char8
}

//export getString8
func getString8(index *C.metaffi_size, indexSize C.metaffi_size, freeRequired *C.metaffi_bool, _ unsafe.Pointer) C.metaffi_string8 {
	cctxt := constructContextTLS.Get()
	*freeRequired = C.metaffi_bool(1)
	val := getElement(index, indexSize, cctxt.Input)
	s := val.Interface().(string)
	cstr := C.CString(s)
	defer C.free(unsafe.Pointer(cstr))
	return C.cast_to_metaffi_string8(C.xllr_alloc_string(cstr, C.uint64_t(len(s))))
}

//export getChar16
func getChar16(index *C.metaffi_size, indexSize C.metaffi_size, _ unsafe.Pointer) C.struct_metaffi_char16 {
	// Convert the context to a Go slice
	cctxt := constructContextTLS.Get()
	val := getElement(index, indexSize, cctxt.Input)

	// Retrieve the rune at the specified index
	r := val.Interface().(rune)

	// Convert the rune to a UTF-16 byte array
	utf16Bytes := utf16.Encode([]rune{r})

	// Create a C.struct_metaffi_char16 and copy the bytes into it
	var char16 C.struct_metaffi_char16
	char16.c[0] = C.char16_t(utf16Bytes[0])
	if len(utf16Bytes) > 1 {
		char16.c[1] = C.char16_t(utf16Bytes[1])
	}

	return char16
}

//export getString16
func getString16(index *C.metaffi_size, indexSize C.metaffi_size, freeRequired *C.metaffi_bool, _ unsafe.Pointer) *C.char16_t {
	cctxt := constructContextTLS.Get()
	*freeRequired = C.metaffi_bool(1)
	val := getElement(index, indexSize, cctxt.Input)
	s := val.Interface().(string)

	// Convert the string to UTF-16
	utf16Str := utf16.Encode([]rune(s))

	// Allocate memory for the UTF-16 string
	mem := C.malloc(C.size_t(len(utf16Str)) * C.size_t(unsafe.Sizeof(C.char16_t(0))))

	// Cast the allocated memory to a pointer to a C.char16_t
	p := (*[1 << 30]C.char16_t)(mem)

	// Copy the UTF-16 string to the allocated memory
	for i, v := range utf16Str {
		p[i] = C.char16_t(v)
	}

	// Return the pointer to the allocated memory
	return (*C.char16_t)(mem)
}

//export getChar32
func getChar32(index *C.metaffi_size, indexSize C.metaffi_size, _ unsafe.Pointer) C.struct_metaffi_char32 {

	cctxt := constructContextTLS.Get()
	val := getElement(index, indexSize, cctxt.Input)

	r := val.Interface().(rune)

	// Convert the rune to a char32_t type
	char32 := C.char32_t(r)

	// Return a new C.struct_metaffi_char32 with the converted rune
	return C.struct_metaffi_char32{c: char32}
}

//export getString32
func getString32(index *C.metaffi_size, indexSize C.metaffi_size, freeRequired *C.metaffi_bool, _ unsafe.Pointer) C.metaffi_string32 {
	cctxt := constructContextTLS.Get()
	*freeRequired = C.metaffi_bool(1)
	val := getElement(index, indexSize, cctxt.Input)
	str := val.Interface().(string)

	// Normalize the string to its canonical form
	normalized := norm.NFC.String(str)

	// Convert the normalized string to a slice of runes
	runes := []rune(normalized)

	// Allocate a C array to hold the runes
	cArray := C.malloc(C.size_t(len(runes)) * C.size_t(unsafe.Sizeof(C.char32_t(0))))

	// Create a Go array backed by the C array
	goArray := (*[1 << 30]C.char32_t)(cArray)

	// Copy the runes to the C array
	for i, r := range runes {
		goArray[i] = C.char32_t(r)
	}

	// Return the C array
	return (*C.char32_t)(cArray)
}

//export getHandle
func getHandle(index *C.metaffi_size, indexSize C.metaffi_size, freeRequired *C.metaffi_bool, _ unsafe.Pointer) *C.struct_cdt_metaffi_handle {
	*freeRequired = C.metaffi_bool(1)
	cctxt := constructContextTLS.Get()
	val := getElement(index, indexSize, cctxt.Input)

	if !val.IsValid() {
		return nil
	}

	cdt_handle := (*C.struct_cdt_metaffi_handle)(C.malloc(C.size_t(unsafe.Sizeof(C.struct_cdt_metaffi_handle{}))))
	GoObjectToMetaffiHandle(cdt_handle, val.Interface())

	return cdt_handle
}

//export getCallable
func getCallable(index *C.metaffi_size, indexSize C.metaffi_size, freeRequired *C.metaffi_bool, _ unsafe.Pointer) *C.struct_cdt_metaffi_callable {
	panic("Not implemented yet")
}

//export getArrayMetadata
func getArrayMetadata(index *C.metaffi_size, indexSize C.metaffi_size, isFixedDimension *C.metaffi_bool, is1DArray *C.metaffi_bool, commonType *C.metaffi_type, isManuallyConstructArray *C.metaffi_bool, _ unsafe.Pointer) C.metaffi_size {
	// Initialize isFixedDimension to true
	*isFixedDimension = C.metaffi_bool(1)
	// Initialize is1DArray to true
	*is1DArray = C.metaffi_bool(1)
	// Initialize isManuallyConstructArray to false
	*isManuallyConstructArray = C.metaffi_bool(0)

	// Get the Go slice from the context
	cctxt := constructContextTLS.Get()

	// get array
	v := getElement(index, indexSize, cctxt.Input)

	arrayOfType, is1darray := getMetaFFITypeFromGoType(v)
	if is1darray {
		*is1DArray = C.metaffi_bool(1)
	} else {
		*is1DArray = C.metaffi_bool(0)
	}

	*commonType = arrayOfType & ^C.metaffi_type(C.metaffi_array_type)

	// Return the length of the slice
	return C.metaffi_size(v.Len())
}
