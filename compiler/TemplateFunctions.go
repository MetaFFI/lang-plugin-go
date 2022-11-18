package main

import "C"
import (
	"fmt"
	"os"
	"strings"
	
	"github.com/MetaFFI/plugin-sdk/compiler/go/IDL"
)

var templatesFuncMap = map[string]interface{}{
	"AsPublic":                        asPublic,
	"ToGoNameConv":                    toGoNameConv,
	"CastIfNeeded":                    castIfNeeded,
	"ParamActual":                     paramActual,
	"GetEnvVar":                       getEnvVar,
	"Sizeof":                          Sizeof,
	"CalculateArgsLength":             calculateArgsLength,
	"CalculateArgLength":              calculateArgLength,
	"Add":                             add,
	"IsInteger":                       isInteger,
	"IsParametersOrReturnValues":      isParametersOrReturnValues,
	"ConvertToCType":                  convertToCType,
	"ConvertToGoType":                 convertToGoType,
	"GetNumericTypes":                 getNumericTypes,
	"GetMetaFFIType":                  getMetaFFIType,
	"GetMetaFFIArrayType":             getMetaFFIArrayType,
	"GetMetaFFIStringTypes":           getMetaFFIStringTypes,
	"MakeMetaFFIType":                 makeMetaFFIType,
	"MethodNameNotExists":             methodNameNotExists,
	"GenerateCodeAllocateCDTS":        generateCodeAllocateCDTS,
	"GenerateCodeXCall":               generateCodeXCall,
	"GenerateCodeEntryPointSignature": generateCodeEntrypointSignature,
	"GetCDTReturnValueIndex":          getCDTReturnValueIndex,
	"GetCDTParametersIndex":           getCDTParametersIndex,
	"GenerateMethodReceiverCode":      generateMethodReceiverCode,
	"GenerateMethodName":              generateMethodName,
	"GenerateMethodParams":            generateMethodParams,
}

//--------------------------------------------------------------------
func generateMethodParams(meth *IDL.MethodDefinition) string {
	//{{range $index, $elem := $f.Parameters}}{{if gt $index 0}}{{if gt $index 1}},{{end}} {{$elem.Name}} {{ConvertToGoType $elem}}{{end}}{{end}}
	
	res := make([]string, 0)
	
	for i, p := range meth.Parameters {
		if i == 0 {
			if !meth.InstanceRequired {
				res = append(res, fmt.Sprintf("%v %v", p.Name, convertToGoType(p)))
			}
			continue
		}
		
		res = append(res, fmt.Sprintf("%v %v", p.Name, convertToGoType(p)))
	}
	
	return strings.Join(res, ",")
}

//--------------------------------------------------------------------
func generateMethodName(meth *IDL.MethodDefinition) string {
	if meth.InstanceRequired {
		return toGoNameConv(meth.Name)
	} else {
		return fmt.Sprintf("%v_%v", toGoNameConv(meth.GetParent().Name), toGoNameConv(meth.Name))
	}
}

//--------------------------------------------------------------------
func generateMethodReceiverCode(meth *IDL.MethodDefinition) string {
	if meth.InstanceRequired {
		return fmt.Sprintf("(this *%v)", meth.GetParent().Name)
	} else {
		return "" // No receiver
	}
}

//--------------------------------------------------------------------
func getCDTReturnValueIndex(params []*IDL.ArgDefinition, retvals []*IDL.ArgDefinition) int {
	if len(params) > 0 {
		return 1
	} else if len(retvals) > 0 {
		return 0
	} else {
		panic("Both parameters and return values are 0 - return values should not be used")
	}
}

//--------------------------------------------------------------------
func getCDTParametersIndex(params []*IDL.ArgDefinition) int {
	if len(params) > 0 {
		return 0
	} else {
		panic("Both parameters and return values are 0 - parameters should not be used")
	}
}

//--------------------------------------------------------------------
func generateCodeEntrypointSignature(className string, funcName string, params []*IDL.ArgDefinition, retvals []*IDL.ArgDefinition) string {
	// {{$c.Name}}_{{$f.Name}}(parameters *C.struct_cdt, parameters_length C.uint64_t, return_values *C.struct_cdt, return_values_length C.uint64_t, out_err **C.char, out_err_len *C.uint64_t)
	
	name := ""
	if className != "" {
		name += className + "_"
	}
	
	name += funcName
	
	if len(params) > 0 || len(retvals) > 0 {
		return fmt.Sprintf("%v(xcall_params *C.struct_cdts, out_err **C.char, out_err_len *C.uint64_t)", name)
	} else {
		return fmt.Sprintf("%v(out_err **C.char, out_err_len *C.uint64_t)", name)
	}
}

//--------------------------------------------------------------------
func generateCodeXCall(className string, funcName string, params []*IDL.ArgDefinition, retvals []*IDL.ArgDefinition) string {
	/*
		var out_err *C.char
		var out_err_len C.uint64_t
		out_err_len = C.uint64_t(0)
	
		C.xllr_xcall(pruntime_plugin, runtime_plugin_length,
				C.int64_t({{$f.Getter.Name}}_id),
				parametersCDTS, {{$paramsLength}},
				return_valuesCDTS, {{$returnLength}},
				&out_err, &out_err_len)
	
		// check errors
		if out_err_len != 0{
			err = fmt.Errorf("Function failed. Error: %v", string(C.GoBytes(unsafe.Pointer(out_err), C.int(out_err_len))))
			return
		}
	*/
	var name string
	if className != "" {
		name += className + "_"
	}
	
	name += funcName
	
	code := "\tvar out_err *C.char\n"
	code += "\tvar out_err_len C.uint64_t\n"
	code += "\tout_err_len = C.uint64_t(0)\n"
	code += "\t\n"
	
	if len(params) > 0 || len(retvals) > 0 {
		code += fmt.Sprintf("\tC.xllr_%v(%v_id, xcall_params, &out_err, &out_err_len)  // call function pointer %v_id via XLLR\n", xcall(params, retvals), name, name)
	} else {
		code += fmt.Sprintf("\tC.xllr_%v(%v_id, &out_err, &out_err_len)  // call function pointer %v_id via XLLR\n", xcall(params, retvals), name, name)
	}
	
	code += "\t\n"
	code += "\t// check errors\n"
	code += "\tif out_err_len != 0{\n"
	code += "\t\terr = fmt.Errorf(\"Function failed. Error: %v\", string(C.GoBytes(unsafe.Pointer(out_err), C.int(out_err_len))))\n"
	code += "\t\treturn\n"
	code += "\t}\n"
	
	return code
}

//--------------------------------------------------------------------
func xcall(params []*IDL.ArgDefinition, retvals []*IDL.ArgDefinition) string {
	
	// name of xcall
	if len(params) > 0 && len(retvals) > 0 {
		return "xcall_params_ret"
	} else if len(params) > 0 {
		return "xcall_params_no_ret"
	} else if len(retvals) > 0 {
		return "xcall_no_params_ret"
	} else {
		return "xcall_no_params_no_ret"
	}
}

//--------------------------------------------------------------------
func generateCodeAllocateCDTS(params []*IDL.ArgDefinition, retval []*IDL.ArgDefinition) string {
	/*
		parametersCDTS := C.xllr_alloc_cdts_buffer( {{$paramsLength}} )
		return_valuesCDTS := C.xllr_alloc_cdts_buffer( {{$returnLength}} )
	*/
	
	if len(params) > 0 || len(retval) > 0 { // use convert_host_params_to_cdts to allocate CDTS
		code := fmt.Sprintf("xcall_params := C.xllr_alloc_cdts_buffer(%v, %v)\n", len(params), len(retval))
		
		code += "\txcall_params_slice := (*[1 << 30]C.cdts)(unsafe.Pointer(xcall_params))[:2:2]\n"
		
		if len(params) > 0 {
			code += "\tparametersCDTS := xcall_params_slice[0].pcdt\n"
			
			if len(retval) > 0 {
				code += "\treturn_valuesCDTS := xcall_params_slice[1].pcdt\n"
			}
		} else {
			code += "\treturn_valuesCDTS := xcall_params_slice[0].pcdt\n"
		}
		
		return code
		
	} else {
		return ""
	}
}

//--------------------------------------------------------------------
func methodNameNotExists(c *IDL.ClassDefinition, fieldName string, prefix string) bool {
	for _, m := range c.Methods {
		if m.Name == prefix+fieldName {
			return false
		}
	}
	
	return true
}

//--------------------------------------------------------------------
func convertToGoType(def *IDL.ArgDefinition) string {
	
	var res string
	
	t := IDL.MetaFFIType(strings.ReplaceAll(string(def.Type), "_array", ""))
	
	switch t {
	case IDL.STRING8:
		fallthrough
	case IDL.STRING16:
		fallthrough
	case IDL.STRING32:
		res = "string"
	case IDL.ANY:
		res = "interface{}"
	case IDL.HANDLE:
		if def.IsTypeAlias() {
			res = def.TypeAlias
		} else {
			res = "interface{}"
		}
	default:
		res = string(def.Type)
	}
	
	if def.IsArray() {
		res = "[]" + res
	}
	
	return res
}

//--------------------------------------------------------------------
func convertToCType(metaffiType IDL.MetaFFIType) string {
	switch metaffiType {
	case "float32":
		return "float"
	case "float64":
		return "double"
	default:
		return string("C." + metaffiType)
	}
}

//--------------------------------------------------------------------
func isParametersOrReturnValues(f *IDL.FunctionDefinition) bool {
	return len(f.Parameters) > 0 || len(f.ReturnValues) > 0
}

//--------------------------------------------------------------------
func isInteger(t string) bool {
	return strings.Index(t, "int") == 0
}

//--------------------------------------------------------------------
func add(x int, y int) int {
	return x + y
}

//--------------------------------------------------------------------
func calculateArgLength(f *IDL.ArgDefinition) int {
	
	if f.IsString() {
		if f.IsArray() {
			return 3 // pointer to string array, pointer to sizes array, length of array
		} else {
			return 2 // pointer to string, size of string
		}
	} else {
		if f.IsArray() {
			return 2 // pointer to type array, length of array
		} else {
			return 1 // value
		}
	}
}

//--------------------------------------------------------------------
func calculateArgsLength(fields []*IDL.ArgDefinition) int {
	
	length := 0
	
	for _, f := range fields {
		length += calculateArgLength(f)
	}
	
	return length
}

//--------------------------------------------------------------------
func Sizeof(field *IDL.ArgDefinition) string {
	return fmt.Sprintf("C.sizeof_metaffi_%v", field.Type)
}

//--------------------------------------------------------------------
func getEnvVar(env string, is_path bool) string {
	res := os.Getenv(env)
	if is_path {
		res = strings.ReplaceAll(res, "\\", "/")
	}
	return res
}

//--------------------------------------------------------------------
func paramActual(field *IDL.ArgDefinition, direction string, namePrefix string) string {
	
	var prefix string
	if namePrefix != "" {
		prefix = namePrefix + "_"
	} else {
		prefix = direction + "_"
	}
	
	switch field.Type {
	case IDL.STRING8:
		fallthrough
	case IDL.STRING16:
		fallthrough
	case IDL.STRING32:
		if field.IsArray() {
			if direction == "out" {
				return fmt.Sprintf("&" + prefix + field.Name + ",&" + prefix + field.Name + "_sizes" + ",&" + prefix + field.Name + "_len")
			} else {
				return fmt.Sprintf(prefix + field.Name + "," + prefix + field.Name + "_sizes" + "," + prefix + field.Name + "_len")
			}
			
		} else {
			
			if direction == "out" {
				return fmt.Sprintf("&" + prefix + field.Name + ",&" + prefix + field.Name + "_len")
			} else {
				return fmt.Sprintf(prefix + field.Name + "," + prefix + field.Name + "_len")
			}
		}
	
	default:
		if field.IsArray() {
			if direction == "out" {
				return fmt.Sprintf("&" + prefix + field.Name + ",&" + prefix + field.Name + "_len")
			} else {
				return fmt.Sprintf(prefix + field.Name + "," + prefix + field.Name + "_len")
			}
			
		} else {
			if direction == "out" {
				return fmt.Sprintf("&" + prefix + field.Name)
			} else {
				return fmt.Sprintf(prefix + field.Name)
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
func toGoNameConv(elem string) string {
	elem = strings.Replace(elem, "_", " ", -1)
	elem = strings.Title(elem)
	return strings.Replace(elem, " ", "", -1)
}

//--------------------------------------------------------------------
func castIfNeeded(elem string) string {
	if strings.Contains(elem, "int") {
		return "int(" + elem + ")"
	}
	return elem
}

//--------------------------------------------------------------------
func getNumericTypes() (numericTypes []string) {
	return []string{"Handle", "float64", "float32", "int8", "int16", "int32", "int64", "uint8", "uint16", "uint32", "uint64"}
}

//--------------------------------------------------------------------
func makeMetaFFIType(t string) string {
	return "metaffi_" + strings.ToLower(t)
}

//--------------------------------------------------------------------
func getMetaFFIStringTypes() (numericTypes []string) {
	return []string{"string8"}
}

//--------------------------------------------------------------------
func getMetaFFIType(numericType string) (numericTypes uint64) {
	return IDL.TypeStringToTypeEnum[IDL.MetaFFIType(numericType)]
}

//--------------------------------------------------------------------
func getMetaFFIArrayType(numericType string) (numericTypes uint64) {
	return IDL.TypeStringToTypeEnum[IDL.MetaFFIType(numericType+"_array")]
}

//--------------------------------------------------------------------
