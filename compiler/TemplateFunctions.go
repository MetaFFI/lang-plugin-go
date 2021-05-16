package main

import (
	"fmt"
	compiler "github.com/OpenFFI/plugin-sdk/compiler/go"
	"strings"
)

var templatesFuncMap = map[string]interface{}{
	"AsPublic": asPublic,
	"ToGoNameConv": toGoNameConv,
	"CastIfNeeded": castIfNeeded,
	"PassParameter": passParameter,
	"ConvertToC": convertToC,
	"ConvertToGo": convertToGo,
}
//--------------------------------------------------------------------
func convertToGo(varName string, elem interface{}) string{
	field := elem.(*compiler.FieldDefinition)

	switch field.Type{
		case compiler.FLOAT64:	return fmt.Sprintf("%v := flaot64(%v)", varName, field.Name)
		case compiler.FLOAT32:	return fmt.Sprintf("%v := float32(%v)", varName, field.Name)

		case compiler.INT8: return fmt.Sprintf("%v := int8(%v)", varName, field.Name)
		case compiler.INT16: return fmt.Sprintf("%v := int16(%v)", varName, field.Name)
		case compiler.INT32: return fmt.Sprintf("%v := int32(%v)", varName, field.Name)
		case compiler.INT64: return fmt.Sprintf("%v := int64(%v)", varName, field.Name)

		case compiler.UINT8: return fmt.Sprintf("%v := uint8(%v)", varName, field.Name)
		case compiler.UINT16: return fmt.Sprintf("%v := uint16(%v)", varName, field.Name)
		case compiler.UINT32: return fmt.Sprintf("%v := uint32(%v)", varName, field.Name)
		case compiler.UINT64: return fmt.Sprintf("%v := uint64(%v)", varName, field.Name)

		case compiler.BOOL: return fmt.Sprintf("%v := %v != C.int8(0)", varName, field.Name)

		case compiler.STRING: return fmt.Sprintf("%v := C.GoStringN(%v, %v)", varName, field.Name, field.Name+"_len")
		case compiler.STRING8: return fmt.Sprintf("%v := C.GoStringN(%v, %v)", varName, field.Name, field.Name+"_len")
		case compiler.STRING16: return fmt.Sprintf("%v := C.GoStringN(%v, %v)", varName, field.Name, field.Name+"_len")
		case compiler.STRING32: return fmt.Sprintf("%v := C.GoStringN(%v, %v)", varName, field.Name, field.Name+"_len")

		case compiler.BYTES: return fmt.Sprintf("%v := C.CString(%v); %v_len := len(%v); defer C.free(%v)", varName, field.Name, varName, field.Name, varName)

		default:
			panic("Unsupported OpenFFI Type "+field.Type)
	}
}
//--------------------------------------------------------------------
func convertToC(varName string, elem interface{}) string{
	field := elem.(*compiler.FieldDefinition)

	switch field.Type{
		case compiler.FLOAT64:	return fmt.Sprintf("%v := C.double(%v)", varName, field.Name)
		case compiler.FLOAT32:	return fmt.Sprintf("%v := C.float(%v)", varName, field.Name)

		case compiler.INT8: return fmt.Sprintf("%v := C.int8_t(%v)", varName, field.Name)
		case compiler.INT16: return fmt.Sprintf("%v := C.int16_t(%v)", varName, field.Name)
		case compiler.INT32: return fmt.Sprintf("%v := C.int32_t(%v)", varName, field.Name)
		case compiler.INT64: return fmt.Sprintf("%v := C.int64_t(%v)", varName, field.Name)

		case compiler.UINT8: return fmt.Sprintf("%v := C.uint8_t(%v)", varName, field.Name)
		case compiler.UINT16: return fmt.Sprintf("%v := C.uint16_t(%v)", varName, field.Name)
		case compiler.UINT32: return fmt.Sprintf("%v := C.uint32_t(%v)", varName, field.Name)
		case compiler.UINT64: return fmt.Sprintf("%v := C.uint64_t(%v)", varName, field.Name)

		case compiler.BOOL: return fmt.Sprintf("var %v C.int8_t; if %v { %v = C.int8_t(1) } else { %v = C.int8_t(0) }", varName, field.Name, varName, varName)

		case compiler.STRING: return fmt.Sprintf("%v := C.CString(%v); %v_len := len(%v); defer C.free(%v)", varName, field.Name, varName, field.Name, varName)
		case compiler.STRING8: return fmt.Sprintf("%v := C.CString(%v); %v_len := len(%v); defer C.free(%v)", varName, field.Name, varName, field.Name, varName)
		case compiler.STRING16: return fmt.Sprintf("%v := C.CString(%v); %v_len := len(%v); defer C.free(%v)", varName, field.Name, varName, field.Name, varName)
		case compiler.STRING32: return fmt.Sprintf("%v := C.CString(%v); %v_len := len(%v); defer C.free(%v)", varName, field.Name, varName, field.Name, varName)

		case compiler.BYTES: return fmt.Sprintf("%v := C.CString(%v); %v_len := len(%v); defer C.free(%v)", varName, field.Name, varName, field.Name, varName)

		default:
			panic("Unsupported OpenFFI Type "+field.Type)
	}
}
//--------------------------------------------------------------------
func asPublic(elem string) string {
	if len(elem) == 0 {
		return ""
	} else if len(elem) == 1 {
		return strings.ToUpper(elem)
	} else {
		return strings.ToUpper(elem[0:1]) + elem[1:]
	}
}
//--------------------------------------------------------------------
func toGoNameConv(elem string) string{
	elem = strings.Replace(elem, "_", " ", -1)
	elem = strings.Title(elem)
	return strings.Replace(elem, " ", "", -1)
}
//--------------------------------------------------------------------
func passParameter(p interface{}) string{
	param := p.(*compiler.FieldDefinition)
	res := "req."+asPublic(param.Name)

	if param.PassMethod == "by_pointer"{
		res = "&"+res
	}

	if strings.Contains(param.Type, "int"){
		res = "int("+res+")"
	}

	return res
}
//--------------------------------------------------------------------
func castIfNeeded(elem string) string{
	if strings.Contains(elem, "int"){
		return "int("+elem+")"
	}
	return elem
}
//--------------------------------------------------------------------