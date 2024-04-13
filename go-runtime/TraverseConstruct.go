package metaffi

/*
#cgo !windows LDFLAGS: -L. -ldl
#cgo LDFLAGS: -Wl,--allow-multiple-definition

#include <include/xllr_capi_loader.h>
#include <include/xllr_capi_loader.c>
#include <include/cdts_traverse_construct.h>
#include <stdint.h>


#include "traverse_construct_go_callbacks.h"


struct traverse_cdts_callbacks initialize_traverse_cdts_callbacks(void* context) {
    struct traverse_cdts_callbacks tcc;
    tcc.context = (void*)context;
    tcc.on_float64 = onFloat64;
    tcc.on_float32 = onFloat32;
    tcc.on_int8 = onInt8;
    tcc.on_uint8 = onUInt8;
    tcc.on_int16 = onInt16;
    tcc.on_uint16 = onUInt16;
    tcc.on_int32 = onInt32;
    tcc.on_uint32 = onUInt32;
    tcc.on_int64 = onInt64;
    tcc.on_uint64 = onUInt64;
    tcc.on_bool = onBool;
    tcc.on_char8 = onChar8;
    tcc.on_string8 = onString8;
    tcc.on_char16 = onChar16;
    tcc.on_string16 = onString16;
    tcc.on_char32 = onChar32;
    tcc.on_string32 = onString32;
    tcc.on_handle = onHandle;
    tcc.on_callable = onCallable;
    tcc.on_null = onNull;
    tcc.on_array = onArray;
    return tcc;
}

struct construct_cdts_callbacks* initialize_construct_cdts_callbacks() {
    struct construct_cdts_callbacks* ccc = malloc(sizeof(struct construct_cdts_callbacks));
    ccc->context = 0;
    ccc->get_float64 = 0;//getFloat64;
    ccc->get_float32 = 0;//getFloat32;
    ccc->get_int8 = 0;//getInt8;
    ccc->get_uint8 = 0;//getUInt8;
    ccc->get_int16 = 0;//getInt16;
    ccc->get_uint16 = 0;//getUInt16;
    ccc->get_int32 = 0;//getInt32;
    ccc->get_uint32 = 0;//getUInt32;
    ccc->get_int64 = 0;//getInt64;
    ccc->get_uint64 = 0;//getUInt64;
    ccc->get_bool = 0;//getBool;
    ccc->get_char8 = 0;//getChar8;
    ccc->get_string8 = 0;//getString8;
    ccc->get_char16 = 0;//getChar16;
    ccc->get_string16 = 0;//getString16;
    ccc->get_char32 = 0;//getChar32;
    ccc->get_string32 = 0;//getString32;
    ccc->get_handle = 0;//getHandle;
    ccc->get_callable = 0;//getCallable;
    ccc->get_array_metadata = 0;//getArrayMetadata;
    ccc->construct_cdt_array = 0;//constructCDTArray;
	ccc->get_type_info = 0;//getTypeInfo;
    return ccc;
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

func NewTraverseCDTSCallbacks() C.struct_traverse_cdts_callbacks {
	return C.initialize_traverse_cdts_callbacks(nil)
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

	if uintptr(h.val) == uintptr(0) {
		return nil
	}

	if h.runtime_id == GO_RUNTIME_ID {
		return GetObject(Handle(h.val))
	} else {
		return MetaFFIHandle{
			Val:       Handle(h.val),
			RuntimeID: uint64(h.runtime_id),
		}
	}
}

func GoObjectToMetaffiHandle(p *C.struct_cdt_metaffi_handle, val interface{}) {
	if h, ok := val.(MetaFFIHandle); ok {
		(*p).val = C.metaffi_handle(h.Val)
		(*p).runtime_id = C.metaffi_size(h.RuntimeID)
	} else {

		if val == nil {
			(*p).val = C.metaffi_handle(uintptr(0))
			(*p).runtime_id = 0
			return
		} else {
			(*p).val = C.metaffi_handle(SetObject(val))
			(*p).runtime_id = GO_RUNTIME_ID
		}
	}
}

//func TraverseCDT(item C.struct_cdt, currentIndex []C.metaffi_size) {
//	if item._type == C.metaffi_any_type {
//		panic("traversed CDT must have a concrete type, not dynamic type like metaffi_any_type")
//	}
//
//	commonType := C.metaffi_any_type
//	typeToUse := item._type
//	if typeToUse&C.metaffi_array_type != 0 && typeToUse != C.metaffi_array_type {
//		commonType = typeToUse &^ C.metaffi_array_type
//		typeToUse = C.metaffi_array_type
//	}
//
//	switch typeToUse {
//	case C.metaffi_float64_type:
//		onFloat64(&currentIndex[0], C.metaffi_size(len(currentIndex)), item.cdt_val.float64_val, nil)
//	case C.metaffi_float32_type:
//		onFloat32(&currentIndex[0], C.metaffi_size(len(currentIndex)), item.cdt_val.float32_val, nil)
//	case C.metaffi_int8_type:
//		onInt8(&currentIndex[0], C.metaffi_size(len(currentIndex)), item.cdt_val.int8_val, nil)
//	case C.metaffi_uint8_type:
//		onUInt8(&currentIndex[0], C.metaffi_size(len(currentIndex)), item.cdt_val.uint8_val, nil)
//	case C.metaffi_int16_type:
//		onInt16(&currentIndex[0], C.metaffi_size(len(currentIndex)), item.cdt_val.int16_val, nil)
//	case C.metaffi_uint16_type:
//		onUInt16(&currentIndex[0], C.metaffi_size(len(currentIndex)), item.cdt_val.uint16_val, nil)
//	case C.metaffi_int32_type:
//		onInt32(&currentIndex[0], C.metaffi_size(len(currentIndex)), item.cdt_val.int32_val, nil)
//	case C.metaffi_uint32_type:
//		onUInt32(&currentIndex[0], C.metaffi_size(len(currentIndex)), item.cdt_val.uint32_val, nil)
//	case C.metaffi_int64_type:
//		onInt64(&currentIndex[0], C.metaffi_size(len(currentIndex)), item.cdt_val.int64_val, nil)
//	case C.metaffi_uint64_type:
//		onUInt64(&currentIndex[0], C.metaffi_size(len(currentIndex)), item.cdt_val.uint64_val, nil)
//	case C.metaffi_bool_type:
//		onBool(&currentIndex[0], C.metaffi_size(len(currentIndex)), item.cdt_val.bool_val, nil)
//	case C.metaffi_char8_type:
//		onChar8(&currentIndex[0], C.metaffi_size(len(currentIndex)), item.cdt_val.char8_val, nil)
//	case C.metaffi_string8_type:
//		onString8(&currentIndex[0], C.metaffi_size(len(currentIndex)), item.cdt_val.string8_val, nil)
//	case C.metaffi_char16_type:
//		onChar16(&currentIndex[0], C.metaffi_size(len(currentIndex)), item.cdt_val.char16_val, nil)
//	case C.metaffi_string16_type:
//		onString16(&currentIndex[0], C.metaffi_size(len(currentIndex)), item.cdt_val.string16_val, 0, nil) // TODO length
//	case C.metaffi_char32_type:
//		onChar32(&currentIndex[0], C.metaffi_size(len(currentIndex)), item.cdt_val.char32_val, nil)
//	case C.metaffi_string32_type:
//		onString32(&currentIndex[0], C.metaffi_size(len(currentIndex)), item.cdt_val.string32_val, 0, nil) // TODO length
//	case C.metaffi_handle_type:
//		onHandle(&currentIndex[0], C.metaffi_size(len(currentIndex)), item.cdt_val.handle_val, nil)
//	case C.metaffi_callable_type:
//		onCallable(&currentIndex[0], C.metaffi_size(len(currentIndex)), item.cdt_val.callable_val, nil)
//	case C.metaffi_null_type:
//		onNull(&currentIndex[0], C.metaffi_size(len(currentIndex)), nil)
//	case C.metaffi_array_type:
//		continueTraverse := onArray(&currentIndex[0], C.metaffi_size(len(currentIndex)), item.cdt_val.array_val, item.cdt_val.array_val.fixed_dimensions, commonType, nil)
//		if continueTraverse != 0 {
//			traverseCDTS(*item.cdt_val.array_val, currentIndex)
//		}
//	default:
//		panic(fmt.Sprintf("Unknown type while traversing CDTS: %v", item._type))
//	}
//}
//
//func traverseCDTS(arr C.struct_cdts, startingIndex []C.metaffi_size) {
//	if arr.length == 0 { // empty CDTS
//		return
//	}
//
//	queue := make([]struct {
//		index []C.metaffi_size
//		pcdt  *C.struct_cdt
//	}, 0)
//
//	for i := C.metaffi_size(0); i < arr.length; i++ {
//		index := append(startingIndex, i)
//		queue = append(queue, struct {
//			index []C.metaffi_size
//			pcdt  *C.struct_cdt
//		}{index: index, pcdt: &arr.arr[i]})
//	}
//
//	for len(queue) > 0 {
//		current := queue[0]
//		queue = queue[1:]
//
//		TraverseCDT(*current.pcdt, current.index)
//	}
//}

//func ConstructCDT(item *C.struct_cdt, callbacks *C.struct_construct_cdts_callbacks, currentIndex []C.metaffi_size, knownType C.struct_metaffi_type_info) {
// var ti C.struct_metaffi_type_info
// if knownType._type == C.metaffi_any_type {
//  ti = callbacks.get_type_info(callbacks.context, &currentIndex[0], C.metaffi_size(len(currentIndex)))
// } else {
//  ti = knownType
// }
//
// if ti._type == C.metaffi_any_type {
//  panic("get_type_info must return a concrete type, not dynamic type like metaffi_any_type")
// }
//
// if ti._type != C.metaffi_any_type && ti.fixed_dimensions > 0 {
//  ti._type |= C.metaffi_array_type
// }
//
// item._type = ti._type
//
// var commonType C.metaffi_type = 0
// if ti._type&C.metaffi_array_type != 0 && ti._type != C.metaffi_array_type {
//  commonType = ti._type &^ C.metaffi_array_type
//  item._type = C.metaffi_array_type
// }
//
// switch item._type {
// case C.metaffi_float64_type:
//  item.cdt_val.float64_val = callbacks.get_float64(callbacks.context, &currentIndex[0], C.metaffi_size(len(currentIndex)))
//  item.free_required = false
// case C.metaffi_float32_type:
//  item.cdt_val.float32_val = callbacks.get_float32(callbacks.context, &currentIndex[0], C.metaffi_size(len(currentIndex)))
//  item.free_required = false
// case C.metaffi_int8_type:
//  item.cdt_val.int8_val = callbacks.get_int8(callbacks.context, &currentIndex[0], C.metaffi_size(len(currentIndex)))
//  item.free_required = false
// case C.metaffi_uint8_type:
//  item.cdt_val.uint8_val = callbacks.get_uint8(callbacks.context, &currentIndex[0], C.metaffi_size(len(currentIndex)))
//  item.free_required = false
// case C.metaffi_int16_type:
//  item.cdt_val.int16_val = callbacks.get_int16(callbacks.context, &currentIndex[0], C.metaffi_size(len(currentIndex)))
//  item.free_required = false
// case C.metaffi_uint16_type:
//  item.cdt_val.uint16_val = callbacks.get_uint16(callbacks.context, &currentIndex[0], C.metaffi_size(len(currentIndex)))
//  item.free_required = false
// case C.metaffi_int32_type:
//  item.cdt_val.int32_val = callbacks.get_int32(callbacks.context, &currentIndex[0], C.metaffi_size(len(currentIndex)))
//  item.free_required = false
// case C.metaffi_uint32_type:
//  item.cdt_val.uint32_val = callbacks.get_uint32(callbacks.context, &currentIndex[0], C.metaffi_size(len(currentIndex)))
//  item.free_required = false
// case C.metaffi_int64_type:
//  item.cdt_val.int64_val = callbacks.get_int64(callbacks.context, &currentIndex[0], C.metaffi_size(len(currentIndex)))
//  item.free_required = false
// case C.metaffi_uint64_type:
//  item.cdt_val.uint64_val = callbacks.get_uint64(callbacks.context, &currentIndex[0], C.metaffi_size(len(currentIndex)))
//  item.free_required = false
// case C.metaffi_bool_type:
//  item.cdt_val.bool_val = callbacks.get_bool(callbacks.context, &currentIndex[0], C.metaffi_size(len(currentIndex)))
//  item.free_required = false
// case C.metaffi_char8_type:
//  item.cdt_val.char8_val = callbacks.get_char8(callbacks.context, &currentIndex[0], C.metaffi_size(len(currentIndex)))
//  item.free_required = false
// case C.metaffi_string8_type:
//  item.cdt_val.string8_val = callbacks.get_string8(callbacks.context, &currentIndex[0], C.metaffi_size(len(currentIndex)))
//  item.free_required = true
// case C.metaffi_char16_type:
//  item.cdt_val.char16_val = callbacks.get_char16(callbacks.context, &currentIndex[0], C.metaffi_size(len(currentIndex)))
//  item.free_required = false
// case C.metaffi_string16_type:
//  item.cdt_val.string16_val = callbacks.get_string16(callbacks.context, &currentIndex[0], C.metaffi_size(len(currentIndex)))
//  item.free_required = true
// case C.metaffi_char32_type:
//  item.cdt_val.char32_val = callbacks.get_char32(callbacks.context, &currentIndex[0], C.metaffi_size(len(currentIndex)))
//  item.free_required = false
// case C.metaffi_string32_type:
//  item.cdt_val.string32_val = callbacks.get_string32(callbacks.context, &currentIndex[0], C.metaffi_size(len(currentIndex)))
//  item.free_required = true
// case C.metaffi_handle_type:
//  item.cdt_val.handle_val = callbacks.get_handle(callbacks.context, &currentIndex[0], C.metaffi_size(len(currentIndex)))
//  item.free_required = true
// case C.metaffi_null_type:
//  item.free_required = false
// case C.metaffi_array_type:
//  item.cdt_val.array_val.fixed_dimensions = C.INT_MIN
//  var isManuallyConstructArray C.metaffi_bool = 0
//  var isFixedDimension C.metaffi_bool = 0
//  var is1DArray C.metaffi_bool = 0
//  arrayLength := callbacks.get_array_metadata(callbacks.context, &currentIndex[0], C.metaffi_size(len(currentIndex)), &isFixedDimension, &is1DArray, &commonType, &isManuallyConstructArray)
//  item._type = commonType == C.metaffi_any_type ? C.metaffi_array_type : (commonType | C.metaffi_array_type)
//  item.free_required = true
//  item.cdt_val.array_val.arr = (*C.struct_cdt)(C.malloc(C.size_t(arrayLength) * C.size_t(unsafe.Sizeof(C.struct_cdt{}))))
//  item.cdt_val.array_val.length = arrayLength
//  if isManuallyConstructArray != 0 {
//   callbacks.construct_cdt_array(callbacks.context, &currentIndex[0], C.metaffi_size(len(currentIndex)), item.cdt_val.array_val)
//  } else {
//   for i := 0; i < int(item.cdt_val.array_val.length); i++ {
//    newIndex := append(currentIndex, C.metaffi_size(i))
//    newItem := &item.cdt_val.array_val.arr[i]
//    ConstructCDT(newItem, callbacks, newIndex, knownType)
//   }
//  }
// case C.metaffi_callable_type:
//  item.cdt_val.callable_val = callbacks.get_callable(callbacks.context, &currentIndex[0], C.metaffi_size(len(currentIndex)))
//  item.free_required = true
// default:
//  panic(fmt.Sprintf("Unknown type while constructing CDTS: %v", item._type))
// }
//}
