package main

import "C"
import (
	"fmt"
	"os"
	"strings"

	"github.com/MetaFFI/plugin-sdk/compiler/go/IDL"
)

var templatesFuncMap = map[string]interface{}{
	"AsPublic":         asPublic,
	"ToGoNameConv":     toGoNameConv,
	"CastIfNeeded":     castIfNeeded,
	"ParamActual":                paramActual,
	"GetEnvVar":                  getEnvVar,
	"Sizeof":                     Sizeof,
	"CalculateArgsLength":        calculateArgsLength,
	"CalculateArgLength":         calculateArgLength,
	"Add":                        add,
	"IsInteger":                  isInteger,
	"IsParametersOrReturnValues": isParametersOrReturnValues,
	"ConvertToCType":             convertToCType,
	"ConvertToGoType":            convertToGoType,
	"GetNumericTypes":            getNumericTypes,
	"GetMetaFFIType":      getMetaFFIType,
	"GetMetaFFIArrayType":        getMetaFFIArrayType,
	"GetMetaFFIStringTypes":      getMetaFFIStringTypes,
	"MakeMetaFFIType": makeMetaFFIType,
	"IsMetaFFIGoRuntimeNeeded": isMetaFFIGoRuntimeNeeded,
	"MethodNameNotExists": methodNameNotExists,
}
//--------------------------------------------------------------------
func methodNameNotExists(c *IDL.ClassDefinition, fieldName string, prefix string) bool{
	for _, m := range c.Methods{
		if m.Name == prefix+fieldName{
			return false
		}
	}

	return true
}
//--------------------------------------------------------------------
func isMetaFFIGoRuntimeNeeded(defs []*IDL.ModuleDefinition) bool{

	isHanleOrAny := func(f *IDL.ArgDefinition) bool{
		return f.IsHandle() || f.IsAny()
	}

	for _, def := range defs {
		for _, f := range def.Functions {
			for _, p := range f.Parameters {
				if isHanleOrAny(p) {
					return true
				}
			}

			for _, r := range f.ReturnValues {
				if isHanleOrAny(r) {
					return true
				}
			}
		}
	}

	return false
}
//--------------------------------------------------------------------
func convertToGoType(def *IDL.ArgDefinition) string{

	var res string

	switch def.Type {
		case IDL.STRING8: fallthrough
		case IDL.STRING16: fallthrough
		case IDL.STRING32:
			res = "string"
		case IDL.ANY: return "interface{}"
		case IDL.HANDLE: return "interface{}"
		default:
			res = string(def.Type)
	}

	if def.IsArray(){
		res = "[]"+res
	}

	return res
}
//--------------------------------------------------------------------
func convertToCType(metaffiType IDL.MetaFFIType) string{
	switch metaffiType {
		case "float32": return "float"
		case "float64": return "double"
		default:
			return string("C."+metaffiType)
	}
}
//--------------------------------------------------------------------
func isParametersOrReturnValues(f *IDL.FunctionDefinition) bool{
	return len(f.Parameters) > 0 || len(f.ReturnValues) > 0
}
//--------------------------------------------------------------------
func isInteger(t string) bool{
	return strings.Index(t, "int") == 0
}
//--------------------------------------------------------------------
func add(x int, y int) int{
	return x + y
}
//--------------------------------------------------------------------
func calculateArgLength(f *IDL.ArgDefinition) int{

	if f.IsString(){
		if f.IsArray(){
			return 3 // pointer to string array, pointer to sizes array, length of array
		} else {
			return 2 // pointer to string, size of string
		}
	} else {
		if f.IsArray(){
			return 2 // pointer to type array, length of array
		} else {
			return 1 // value
		}
	}
}
//--------------------------------------------------------------------
func calculateArgsLength(fields []*IDL.ArgDefinition) int{

	length := 0

	for _, f := range fields{
		length += calculateArgLength(f)
	}

	return length
}
//--------------------------------------------------------------------
func Sizeof(field *IDL.ArgDefinition) string{
	return fmt.Sprintf("C.sizeof_metaffi_%v", field.Type)
}
//--------------------------------------------------------------------
func getEnvVar(env string) string{
	return os.Getenv(env)
}
//--------------------------------------------------------------------
func paramActual(field *IDL.ArgDefinition, direction string, namePrefix string) string{

	var prefix string
	if namePrefix != ""{
		prefix = namePrefix +"_"
	} else {
		prefix = direction +"_"
	}


	switch field.Type {
		case IDL.STRING8: fallthrough
		case IDL.STRING16: fallthrough
		case IDL.STRING32:
			if field.IsArray(){
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
			if field.IsArray(){
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
func castIfNeeded(elem string) string{
	if strings.Contains(elem, "int"){
		return "int("+elem+")"
	}
	return elem
}
//--------------------------------------------------------------------
func getNumericTypes() (numericTypes []string ){
	return []string{ "Handle", "float64", "float32", "int8", "int16", "int32", "int64", "uint8", "uint16", "uint32", "uint64" }
}
//--------------------------------------------------------------------
func makeMetaFFIType(t string) string{
	return "metaffi_"+strings.ToLower(t)
}
//--------------------------------------------------------------------
func getMetaFFIStringTypes() (numericTypes []string ){
	return []string{ "string8" }
}
//--------------------------------------------------------------------
func getMetaFFIType(numericType string) (numericTypes uint64){
	return IDL.TypeStringToTypeEnum[IDL.MetaFFIType(numericType)]
}
//--------------------------------------------------------------------
func getMetaFFIArrayType(numericType string) (numericTypes uint64){
	return IDL.TypeStringToTypeEnum[IDL.MetaFFIType(numericType+"_array")]
}
//--------------------------------------------------------------------