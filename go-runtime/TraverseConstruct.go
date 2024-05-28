package metaffi

/*
#cgo !windows LDFLAGS: -L. -ldl
#cgo LDFLAGS: -Wl,--allow-multiple-definition

#include <include/xllr_capi_loader.h>
#include <include/xllr_capi_loader.c>
#include <include/cdts_traverse_construct.h>
#include <stdint.h>


#include "traverse_construct_go_callbacks.h"


struct traverse_cdts_callbacks* initialize_traverse_cdts_callbacks() {
    struct traverse_cdts_callbacks* tcc = malloc(sizeof(struct traverse_cdts_callbacks));
    tcc->context = 0;
    tcc->on_float64 = onFloat64;
    tcc->on_float32 = onFloat32;
    tcc->on_int8 = onInt8;
    tcc->on_uint8 = onUInt8;
    tcc->on_int16 = onInt16;
    tcc->on_uint16 = onUInt16;
    tcc->on_int32 = onInt32;
    tcc->on_uint32 = onUInt32;
    tcc->on_int64 = onInt64;
    tcc->on_uint64 = onUInt64;
    tcc->on_bool = onBool;
    tcc->on_char8 = onChar8;
    tcc->on_string8 = onString8;
    tcc->on_char16 = onChar16;
    tcc->on_string16 = onString16;
    tcc->on_char32 = onChar32;
    tcc->on_string32 = onString32;
    tcc->on_handle = onHandle;
    tcc->on_callable = onCallable;
    tcc->on_null = onNull;
    tcc->on_array = onArray;
    return tcc;
}

struct construct_cdts_callbacks* initialize_construct_cdts_callbacks() {
    struct construct_cdts_callbacks* ccc = malloc(sizeof(struct construct_cdts_callbacks));
    ccc->context = 0;
    ccc->get_float64 = getFloat64;
    ccc->get_float32 = getFloat32;
    ccc->get_int8 = getInt8;
    ccc->get_uint8 = getUInt8;
    ccc->get_int16 = getInt16;
    ccc->get_uint16 = getUInt16;
    ccc->get_int32 = getInt32;
    ccc->get_uint32 = getUInt32;
    ccc->get_int64 = getInt64;
    ccc->get_uint64 = getUInt64;
    ccc->get_bool = getBool;
    ccc->get_char8 = getChar8;
    ccc->get_string8 = getString8;
    ccc->get_char16 = getChar16;
    ccc->get_string16 = getString16;
    ccc->get_char32 = getChar32;
    ccc->get_string32 = getString32;
    ccc->get_handle = getHandle;
    ccc->get_callable = getCallable;
    ccc->get_array_metadata = getArrayMetadata;
    ccc->construct_cdt_array = constructCDTArray;
	ccc->get_type_info = getTypeInfo;
    return ccc;
}

void GoMetaFFIHandleTocdt_metaffi_handle(struct cdt_metaffi_handle* p , void* handle, uint64_t runtime_id, void* release) {
	p->handle = (metaffi_handle)handle;
	p->runtime_id = runtime_id;
	p->release = (void (*)(struct cdt_metaffi_handle*))release;
}

*/
import "C"
import (
	"fmt"
	"reflect"
	"unsafe"
)

func getElement(index *C.metaffi_size, indexSize C.metaffi_size, root interface{}) reflect.Value {

	if index == nil {
		return reflect.ValueOf(root)
	}

	// traverse the root object to get the element
	// Convert C array to Go slice
	indexSlice := (*[1 << 30]uint64)(unsafe.Pointer(index))[:indexSize:indexSize]

	// Traverse the root object
	v := reflect.ValueOf(root)
	for _, idx := range indexSlice {
		if v.Kind() == reflect.Slice || v.Kind() == reflect.Array {
			if idx < uint64(v.Len()) {
				v = v.Index(int(idx))
			} else {
				panic(fmt.Sprintf("Index out of range: %v. Length: %v", idx, v.Len()))
			}
		} else {
			panic(fmt.Sprintf("Unsupported type: %T", v.Interface()))
		}
	}

	return v
}

func createMultiDimSlice(length int, dimensions int, elemType reflect.Type) interface{} {
	if dimensions <= 0 {
		return nil
	}

	// Create a slice type for each dimension
	for i := 0; i < dimensions; i++ {
		elemType = reflect.SliceOf(elemType)
	}

	// Create a slice of the final type
	slice := reflect.MakeSlice(elemType, length, length)

	return slice.Interface()
}

func getGoTypeFromMetaFFIType(metaffiType C.metaffi_type, commonGoType reflect.Type) reflect.Type {

	switch metaffiType {
	case C.metaffi_float64_type:
		return reflect.TypeOf(float64(0))
	case C.metaffi_float32_type:
		return reflect.TypeOf(float32(0))
	case C.metaffi_int8_type:
		return reflect.TypeOf(int8(0))
	case C.metaffi_uint8_type:
		return reflect.TypeOf(uint8(0))
	case C.metaffi_int16_type:
		return reflect.TypeOf(int16(0))
	case C.metaffi_uint16_type:
		return reflect.TypeOf(uint16(0))
	case C.metaffi_int32_type:
		return reflect.TypeOf(int32(0))
	case C.metaffi_uint32_type:
		return reflect.TypeOf(uint32(0))
	case C.metaffi_int64_type:
		return reflect.TypeOf(int64(0))
	case C.metaffi_uint64_type:
		return reflect.TypeOf(uint64(0))
	case C.metaffi_bool_type:
		return reflect.TypeOf(false)
	case C.metaffi_char8_type:
		return reflect.TypeOf(rune(0))
	case C.metaffi_string8_type:
		return reflect.TypeOf(string(""))
	case C.metaffi_char16_type:
		return reflect.TypeOf(rune(0))
	case C.metaffi_string16_type:
		return reflect.TypeOf(string(""))
	case C.metaffi_char32_type:
		return reflect.TypeOf(rune(0))
	case C.metaffi_string32_type:
		return reflect.TypeOf("")
	case C.metaffi_any_type:
		fallthrough
	case C.metaffi_handle_type:
		if commonGoType == nil {
			panic("metaffi_handle_type requires a common Go type")
		}
		return commonGoType
	case C.metaffi_callable_type:
		return reflect.TypeOf(func() {})
	default:
		panic(fmt.Sprintf("Cannot find requested MetaFFI Type: %v", metaffiType))
	}
}

type OnFloat64Func func(index *C.metaffi_size, indexSize C.metaffi_size, val C.metaffi_float64, context unsafe.Pointer)

func NewTraverseCDTSCallbacks() *C.struct_traverse_cdts_callbacks {
	return C.initialize_traverse_cdts_callbacks()
}

//--------------------------------------------------------------------

func getMetaFFITypeFromGoType(v reflect.Value) (detectedType C.metaffi_type, is1DArray bool) {
	t := v.Type()
	arrayMask := C.metaffi_type(0)
	is1DArray = false

	if t.Kind() == reflect.Slice {
		is1DArray = true
		arrayMask = C.metaffi_array_type
		t = t.Elem()
	}

	for t.Kind() == reflect.Slice { // for multi-dimensional arrays
		is1DArray = false
		t = t.Elem()
	}

	switch t.Kind() {
	case reflect.Float64:
		return arrayMask | C.metaffi_float64_type, is1DArray
	case reflect.Float32:
		return arrayMask | C.metaffi_float32_type, is1DArray
	case reflect.Int8:
		return arrayMask | C.metaffi_int8_type, is1DArray
	case reflect.Uint8:
		return arrayMask | C.metaffi_uint8_type, is1DArray
	case reflect.Int16:
		return arrayMask | C.metaffi_int16_type, is1DArray
	case reflect.Uint16:
		return arrayMask | C.metaffi_uint16_type, is1DArray
	case reflect.Int32:
		return arrayMask | C.metaffi_int32_type, is1DArray
	case reflect.Uint32:
		return arrayMask | C.metaffi_uint32_type, is1DArray
	case reflect.Int64:
		return arrayMask | C.metaffi_int64_type, is1DArray
	case reflect.Uint64:
		return arrayMask | C.metaffi_uint64_type, is1DArray
	case reflect.Bool:
		return arrayMask | C.metaffi_bool_type, is1DArray
	case reflect.String:
		return arrayMask | C.metaffi_string8_type, is1DArray
	case reflect.Func:
		return arrayMask | C.metaffi_callable_type, is1DArray
	case reflect.Interface:
		if t.NumMethod() == 0 {

			// []interface{}
			if arrayMask == C.metaffi_array_type {

				// if one of the elements is a slice - 1D array is false
				for i := 0; i < v.Len(); i++ {
					curv := v.Index(i)
					if curv.Kind() == reflect.Slice {
						is1DArray = false
						detectedType, _ = getMetaFFITypeFromGoType(curv)
						break
					} else if curv.Elem().Kind() == reflect.Slice {
						is1DArray = false
						detectedType, _ = getMetaFFITypeFromGoType(curv.Elem())
						break
					} else {
						var isInner1DArray bool
						detectedType, isInner1DArray = getMetaFFITypeFromGoType(curv)
						if isInner1DArray && is1DArray {
							is1DArray = false
						}
					}
				}

			} else { // interface{}

				var isInner1DArray bool
				detectedType, isInner1DArray = getMetaFFITypeFromGoType(v.Elem())
				if isInner1DArray && is1DArray {
					is1DArray = false
				}
			}

			return arrayMask | detectedType, is1DArray

		} else {
			return arrayMask | C.metaffi_handle_type, is1DArray // interface of a struct - handle
		}
	default:
		return arrayMask | C.metaffi_handle_type, is1DArray
	}
}

func NewConstructCDTSCallbacks() *C.struct_construct_cdts_callbacks {
	return C.initialize_construct_cdts_callbacks()
}

func GetGoObject(h *C.struct_cdt_metaffi_handle) interface{} {

	if uintptr(h.handle) == uintptr(0) {
		return nil
	}

	if h.runtime_id == GO_RUNTIME_ID {
		return GetObject(Handle(h.handle))
	} else {
		return MetaFFIHandle{
			Val:       Handle(h.handle),
			RuntimeID: uint64(h.runtime_id),
			CReleaser: unsafe.Pointer(h.release),
		}
	}
}

func GoObjectToMetaffiHandle(p *C.struct_cdt_metaffi_handle, val interface{}) {
	if h, ok := val.(MetaFFIHandle); ok {
		C.GoMetaFFIHandleTocdt_metaffi_handle(p, unsafe.Pointer(h.Val), C.uint64_t(h.RuntimeID), h.CReleaser)
	} else {

		// set Go object into cdt_metaffi_handle
		if val == nil {
			(*p).handle = C.metaffi_handle(uintptr(0))
			(*p).runtime_id = 0
			(*p).release = nil
		} else {
			C.GoMetaFFIHandleTocdt_metaffi_handle(p, unsafe.Pointer(SetObject(val)), GO_RUNTIME_ID, GetReleaserCFunction()
		}
	}
}
