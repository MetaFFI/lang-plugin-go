package main

import "C"
import (
	"fmt"
	compiler "github.com/OpenFFI/plugin-sdk/compiler/go"
	"os"
	"strings"
)

var templatesFuncMap = map[string]interface{}{
	"AsPublic": asPublic,
	"ToGoNameConv": toGoNameConv,
	"CastIfNeeded": castIfNeeded,
	"PassParameter": passParameter,
	"ConvertToC": convertToC,
	"ConvertToGo": convertToGo,
	"ParamActual": paramActual,
	"FormalCParameter": formalCParameter,
	"GetEnvVar": getEnvVar,
	"Sizeof": Sizeof,
}
//--------------------------------------------------------------------
func formalCParameter(field *compiler.FieldDefinition, direction string) string{
	ctype := "openffi_"+field.Type

	if direction == "out"{
		ctype += "*"
	}

	if field.IsArray{
		ctype += "*"
	}

	cname := direction+"_"+field.Name

	result := fmt.Sprintf("%v %v", ctype, cname)

	if strings.Index(field.Type, compiler.STRING) == 0{

		if field.IsArray{
			if direction == "out"{
				result += ", openffi_size** "+direction+"_"+field.Name+"_sizes, openffi_size* "+direction+"_"+field.Name+"_len"
			} else {
				result += ", openffi_size* "+direction+"_"+field.Name+"_sizes, openffi_size "+direction+"_"+field.Name+"_len"
			}
		} else {
			 // add length
			if direction == "out"{
				result += ", openffi_size* "+direction+"_"+field.Name+"_len"
			} else {
				result += ", openffi_size "+direction+"_"+field.Name+"_len"
			}
		}

	} else {

		if field.IsArray { // add another parameter for length
			if direction == "out"{
				result += ", openffi_size* "+direction+"_"+field.Name+"_len"
			} else {
				result += ", openffi_size "+direction+"_"+field.Name+"_len"
			}
		}
	}

	return result
}
//--------------------------------------------------------------------
func Sizeof(field *compiler.FieldDefinition) string{
	return fmt.Sprintf("C.sizeof_openffi_%v", field.Type)
}
//--------------------------------------------------------------------
func getEnvVar(env string) string{
	return os.Getenv(env)
}
//--------------------------------------------------------------------
func paramActual(field *compiler.FieldDefinition, direction string, namePrefix string) string{

	var prefix string
	if namePrefix != ""{
		prefix = namePrefix +"_"
	} else {
		prefix = direction +"_"
	}


	switch field.Type {
		case compiler.STRING: fallthrough
		case compiler.STRING8: fallthrough
		case compiler.STRING16: fallthrough
		case compiler.STRING32:
			if field.IsArray{
				if direction == "out"{
					return fmt.Sprintf("&"+prefix+field.Name+",&"+prefix+field.Name+"_sizes"+",&"+prefix+field.Name+"_len")
				} else {
					return fmt.Sprintf(prefix+field.Name+","+prefix+field.Name+"_sizes"+","+prefix+field.Name+"_len")
				}

			} else {

				if direction == "out"{
					return fmt.Sprintf("&"+prefix+field.Name+",&"+prefix+field.Name+"_len")
				} else {
					return fmt.Sprintf(prefix+field.Name+","+prefix+field.Name+"_len")
				}
			}

		default:
			if field.IsArray{
				if direction == "out"{
					return fmt.Sprintf("&"+prefix+field.Name+",&"+prefix+field.Name+"_len")
				} else {
					return fmt.Sprintf(prefix+field.Name+","+prefix+field.Name+"_len")
				}

			} else {
				if direction == "out"{
					return fmt.Sprintf("&"+prefix+field.Name)
				} else {
					return fmt.Sprintf(prefix+field.Name)
				}
			}
	}
}
//--------------------------------------------------------------------
func convertToGo(field *compiler.FieldDefinition, fieldPrefix, varPrefix string) string{

	varName := varPrefix+"_"+field.Name
	fieldName := fieldPrefix+"_"+field.Name

	switch field.Type{
		case compiler.FLOAT64: fallthrough
		case compiler.FLOAT32: fallthrough
		case compiler.INT8: fallthrough
		case compiler.INT16: fallthrough
		case compiler.INT32: fallthrough
		case compiler.INT64: fallthrough
		case compiler.UINT8: fallthrough
		case compiler.UINT16: fallthrough
		case compiler.UINT32: fallthrough
		case compiler.UINT64:
			if field.IsArray{
				return varName+" := make([]"+field.Type+", 0); for _, val := range "+fieldName+"{ "+varName+" = append("+varName+", "+field.Type+"("+fieldName+")) }"
			} else {
				return varName+" := "+field.Type+"("+fieldName+")"
			}

		case compiler.BOOL:
			if field.IsArray{
				return varName+" := make([]"+field.Type+", 0); for _, val := range "+fieldName+"{ "+varName+" = append("+varName+", "+field.Type+"("+fieldName+")) }"
			} else {
				return fmt.Sprintf("%v := %v != C.openffi_bool(0)", varName, fieldName)
			}

		case compiler.STRING: fallthrough
		case compiler.STRING8: fallthrough
		case compiler.STRING16: fallthrough
		case compiler.STRING32:
			if field.IsArray{
				return varName+" := make([]string, 0); for i:=C.openffi_size(0) ; i<"+fieldName+"_len ; i++ { val_size := C.openffi_size(0); val := C.get_string_element(C.int(i), "+fieldName+", "+fieldName+"_sizes, &val_size); "+varName+" = append("+varName+", C.GoStringN(val, C.int(val_size))) }"
			} else {
				return fmt.Sprintf("%v := C.GoStringN(%v, C.int(%v))", varName, fieldName, fieldName+"_len")
			}

		default:
			panic("Unsupported OpenFFI Type "+field.Type)
	}
}
//--------------------------------------------------------------------
func convertToC(field *compiler.FieldDefinition, direction string) string{

	varName := direction +"_"+field.Name

	res := ""

	switch field.Type{
		case compiler.FLOAT64: fallthrough
		case compiler.FLOAT32: fallthrough
		case compiler.INT8:  fallthrough
		case compiler.INT16: fallthrough
		case compiler.INT32: fallthrough
		case compiler.INT64: fallthrough
		case compiler.UINT8: fallthrough
		case compiler.UINT16: fallthrough
		case compiler.UINT32: fallthrough
		case compiler.UINT64:
			if direction != "out" {
				if field.IsArray {
					res = varName + `_arr := make([]C.openffi_` + field.Type + `, 0); `
					res += varName + `_len := C.openffi_size(len(` + field.Name + `)); `
					res += `for _, val := range ` + field.Name + ` { ` + varName + `_arr = append(` + varName + `_arr, C.openffi_` + field.Type + `(val)) }; `
					res += varName + ` := &(` + varName + `_arr[0]); `
				} else {
					return fmt.Sprintf(`%v := C.openffi_`+field.Type+`(%v)`, varName, field.Name)
				}
			} else {
				if field.IsArray {
					return varName+` C.openffi_`+field.Type+`(0); `+varName+`_len := C.openffi_size(0);`
				} else {
					return varName+` := C.openffi_`+field.Type+`(0)`
				}
			}

		case compiler.BOOL:
			if direction != "out" {
				if field.IsArray {
					res = varName + `_arr := make([]C.openffi_` + field.Type + `, 0); `
					res += varName + `_len := C.openffi_size(len(` + field.Name + `)); `
					res += `for _, val := range ` + field.Name + ` { var cval C.openffi_bool; if val{ cval=C.openffi_bool(1) } else { cval=C.openffi_bool(0) }; ` + varName + `_arr = append(` + varName + `_arr, cval) }; `
					res += varName + ` := &(` + varName + `_arr[0]); `

				} else {
					res = `var ` + varName + ` C.openffi_` + field.Type + `; if ` + field.Name + ` { ` + varName + ` = C.openffi_bool(1) } else { ` + varName + ` = C.openffi_bool(0) }`
				}
			} else {
				if field.IsArray {
					return varName+` := C.openffi_`+field.Type+`(0); `+varName+`_len := C.openffi_size(0);`
				} else {
					return varName+` := C.openffi_`+field.Type+`(0)`
				}
			}

		case compiler.STRING: fallthrough
		case compiler.STRING8: fallthrough
		case compiler.STRING16: fallthrough
		case compiler.STRING32:
			if direction != "out" {
				if field.IsArray {
					res = varName + `_arr := make([]C.openffi_` + field.Type + `, 0); `
					res += varName + `_go_sizes := make([]C.openffi_size, 0); `
					res += varName + `_len := C.openffi_size(len(` + field.Name + `)); `
					res += `for _, val := range ` + field.Name + ` { curCval := C.CString(val); ` + varName + `_arr = append(` + varName + `_arr, curCval); ` + varName + `_go_sizes = append(` + varName + `_go_sizes, C.openffi_size(len(val)));  defer C.free(unsafe.Pointer(curCval)) }; `
					res += varName + ` := &(` + varName + `_arr[0]); `
					res += varName + `_sizes := &(` + varName + `_go_sizes[0]); `
				} else {
					res = fmt.Sprintf("%v := C.CString(%v); %v_len := C.ulong(len(%v)); defer C.free(unsafe.Pointer(%v))", varName, field.Name, varName, field.Name, varName)
				}
			} else {
				if field.IsArray {
					return `var `+varName+` *C.openffi_string; var `+varName+`_sizes *C.openffi_size; `+varName+`_len := C.openffi_size(0);`
				} else {
					return varName+` := C.openffi_`+field.Type+`(0)`
				}
			}

		default:
			panic("Unsupported OpenFFI Type "+field.Type)
	}

	return res
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