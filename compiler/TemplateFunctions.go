package main

import "C"
import (
	"fmt"
	compiler "github.com/OpenFFI/plugin-sdk/compiler/go"
	"os"
	"strings"
)

var templatesFuncMap = map[string]interface{}{
	"AsPublic":         asPublic,
	"ToGoNameConv":     toGoNameConv,
	"CastIfNeeded":     castIfNeeded,
	"ParamActual":      paramActual,
	"GetEnvVar":        getEnvVar,
	"Sizeof":           Sizeof,
	"CalculateArgsLength": calculateArgsLength,
	"CalculateArgLength": calculateArgLength,
	"Add": add,
	"IsInteger": isInteger,
	"IsParametersOrReturnValues": isParametersOrReturnValues,
	"ConvertToCType": convertToCType,
	"ConvertToGoType": convertToGoType,
}
//--------------------------------------------------------------------
func convertToGoType(def *compiler.FieldDefinition) string{

	var res string

	switch def.Type {
		case "string8": fallthrough
		case "string16": fallthrough
		case "string32":
			res = "string"
		default:
			res = string(def.Type)
	}

	if def.IsArray(){
		res = "[]"+res
	}

	return res
}
//--------------------------------------------------------------------
func convertToCType(openffiType string) string{
	switch openffiType {
		case "float32": return "float"
		case "float64": return "double"
		default:
			return "C."+openffiType
	}
}
//--------------------------------------------------------------------
func isParametersOrReturnValues(f *compiler.FunctionDefinition) bool{
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
func calculateArgLength(f *compiler.FieldDefinition) int{

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
func calculateArgsLength(fields []*compiler.FieldDefinition) int{

	length := 0

	for _, f := range fields{
		length += calculateArgLength(f)
	}

	return length
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
		case compiler.STRING8: fallthrough
		case compiler.STRING16: fallthrough
		case compiler.STRING32:
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