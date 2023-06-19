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
	"GetTypeOrAlias":                  getTypeOrAlias,
	"HandleNoneGoObject":              handleNoneGoObject,
	"IsGoRuntimePackNeeded":           isGoRuntimePackNeeded,
}

// --------------------------------------------------------------------
func isGoRuntimePackNeeded(idl *IDL.IDLDefinition) bool {

	for _, m := range idl.Modules {
		for _, g := range m.Globals {
			if g.IsHandle() {
				return true
			}
		}

		for _, f := range m.Functions {
			for _, p := range f.Parameters {
				if p.IsHandle() {
					return true
				}
			}

			for _, r := range f.ReturnValues {
				if r.IsHandle() {
					return true
				}
			}
		}

		if len(m.Classes) > 0 {
			return true
		}
	}

	return false
}

// --------------------------------------------------------------------
func handleNoneGoObject(arg *IDL.ArgDefinition, module *IDL.ModuleDefinition) string {
	if arg.IsTypeAlias() && module.IsContainsClass(arg.TypeAlias) {
		return asPublic(arg.TypeAlias) + "{ h: obj }" // construct a "wrapper class" for handle
	} else {
		return "obj" // just return handle
	}
}

// --------------------------------------------------------------------
func getTypeOrAlias(arg *IDL.ArgDefinition, module *IDL.ModuleDefinition) string {
	if arg.IsTypeAlias() && module.IsContainsClass(arg.TypeAlias) {
		return asPublic(arg.TypeAlias)
	} else {
		return asPublic(string(IDL.HANDLE))
	}
}

// --------------------------------------------------------------------
func generateMethodParams(meth *IDL.MethodDefinition, mod *IDL.ModuleDefinition) string {
	//{{range $index, $elem := $f.Parameters}}{{if gt $index 0}}{{if gt $index 1}},{{end}} {{$elem.Name}} {{ConvertToGoType $elem}}{{end}}{{end}}

	res := make([]string, 0)

	for i, p := range meth.Parameters {
		if i == 0 {
			if !meth.InstanceRequired {
				res = append(res, fmt.Sprintf("%v %v", p.Name, convertToGoType(p, mod)))
			}
			continue
		}

		res = append(res, fmt.Sprintf("%v %v", p.Name, convertToGoType(p, mod)))
	}

	return strings.Join(res, ",")
}

// --------------------------------------------------------------------
func generateMethodName(meth *IDL.MethodDefinition) string {
	if meth.InstanceRequired {
		return toGoNameConv(meth.GetNameWithOverloadIndex())
	} else {
		return fmt.Sprintf("%v_%v", toGoNameConv(meth.GetParent().Name), toGoNameConv(meth.GetNameWithOverloadIndex()))
	}
}

// --------------------------------------------------------------------
func generateMethodReceiverCode(meth *IDL.MethodDefinition) string {
	if meth.InstanceRequired {
		return fmt.Sprintf("(this *%v)", asPublic(meth.GetParent().Name))
	} else {
		return "" // No receiver
	}
}

// --------------------------------------------------------------------
func getCDTReturnValueIndex(params []*IDL.ArgDefinition, retvals []*IDL.ArgDefinition) int {
	return 1 // return values are always at index 1
}

// --------------------------------------------------------------------
func getCDTParametersIndex(params []*IDL.ArgDefinition) int {
	if len(params) > 0 {
		return 0
	} else {
		panic("Both parameters and return values are 0 - parameters should not be used")
	}
}

// --------------------------------------------------------------------
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

// --------------------------------------------------------------------
func generateCodeXCall(className string, funcName string, params []*IDL.ArgDefinition, retvals []*IDL.ArgDefinition) string {
	/*
		err = XLLRXCallParamsRet(metaffi_positional_args_Gett_id, xcall_params)  // call function pointer metaffi_positional_args_Gett_id via XLLR

		// check errors
		if err != nil{
			err = fmt.Errorf("Function failed metaffi_positional_args.Gett. Error: %v", err)
			return
		}
	*/
	var name string
	if className != "" {
		name += className + "_"
	}

	name += funcName

	code := fmt.Sprintf("\terr = XLLR%v(%v_id", xcall(params, retvals), name)

	if len(params) > 0 || len(retvals) > 0 {
		code += ", xcall_params"
	}

	code += fmt.Sprintf(")  // call function pointer %v_id via XLLR\n", name)

	code += "\t\n"
	code += "\t// check errors\n"
	code += "\tif err != nil{\n"
	code += "\t\terr = fmt.Errorf(\"Failed calling function" + className + "." + funcName + ". Error: %v\", err)\n"
	code += "\t\treturn\n"
	code += "\t}\n"

	return code
}

// --------------------------------------------------------------------
func xcall(params []*IDL.ArgDefinition, retvals []*IDL.ArgDefinition) string {

	// name of xcall
	if len(params) > 0 && len(retvals) > 0 {
		return "XCallParamsRet"
	} else if len(params) > 0 {
		return "XCallParamsNoRet"
	} else if len(retvals) > 0 {
		return "XCallNoParamsRet"
	} else {
		return "XCallNoParamsNoRet"
	}
}

// --------------------------------------------------------------------
func generateCodeAllocateCDTS(params []*IDL.ArgDefinition, retval []*IDL.ArgDefinition) string {
	/*
		xcall_params, parametersCDTS, return_valuesCDTS := XLLRAllocCDTSBuffer(1, 1)
	*/

	if len(params) == 0 && len(retval) == 0 {
		return ""
	}

	code := "xcall_params, "
	if len(params) > 0 {
		code += "parametersCDTS, "
	} else {
		code += "_, "
	}

	if len(retval) > 0 {
		code += "return_valuesCDTS "
	} else {
		code += "_ "
	}

	code += fmt.Sprintf(":= XLLRAllocCDTSBuffer(%v, %v)\n", len(params), len(retval))

	return code
}

// --------------------------------------------------------------------
func methodNameNotExists(c *IDL.ClassDefinition, fieldName string, prefix string) bool {
	for _, m := range c.Methods {
		if m.Name == prefix+fieldName {
			return false
		}
	}

	return true
}

// --------------------------------------------------------------------
func convertToGoType(def *IDL.ArgDefinition, mod *IDL.ModuleDefinition) string {

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
		if def.IsTypeAlias() && mod.IsContainsClass(def.TypeAlias) {
			res = asPublic(def.TypeAlias)
		} else {
			res = "interface{}"
		}
	default:
		res = string(def.Type)
	}

	if def.IsArray() && t != IDL.ANY {
		res = "[]" + res
	}

	return res
}

// --------------------------------------------------------------------
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

// --------------------------------------------------------------------
func isParametersOrReturnValues(f *IDL.FunctionDefinition) bool {
	return len(f.Parameters) > 0 || len(f.ReturnValues) > 0
}

// --------------------------------------------------------------------
func isInteger(t string) bool {
	return strings.Index(t, "int") == 0
}

// --------------------------------------------------------------------
func add(x int, y int) int {
	return x + y
}

// --------------------------------------------------------------------
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

// --------------------------------------------------------------------
func calculateArgsLength(fields []*IDL.ArgDefinition) int {

	length := 0

	for _, f := range fields {
		length += calculateArgLength(f)
	}

	return length
}

// --------------------------------------------------------------------
func Sizeof(field *IDL.ArgDefinition) string {
	return fmt.Sprintf("C.sizeof_metaffi_%v", field.Type)
}

// --------------------------------------------------------------------
func getEnvVar(env string, is_path bool) string {
	res := os.Getenv(env)
	if is_path {
		res = strings.ReplaceAll(res, "\\", "/")
	}
	return res
}

// --------------------------------------------------------------------
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

// --------------------------------------------------------------------
func asPublic(elem string) string {

	return toGoNameConv(elem)
	//if len(elem) == 0 {
	//	return ""
	//} else if len(elem) == 1 {
	//	return strings.ToUpper(elem)
	//} else {
	//	return strings.ToUpper(elem[0:1]) + elem[1:]
	//}
}

// --------------------------------------------------------------------
func countUnderscores(s string) (int, int) {
	startCount := 0
	endCount := 0
	for i := 0; i < len(s); i++ {
		if s[i] == '_' {
			startCount++
		} else {
			break
		}
	}
	for i := len(s) - 1; i >= 0; i-- {
		if s[i] == '_' {
			endCount++
		} else {
			break
		}
	}
	return startCount, endCount
}

// --------------------------------------------------------------------
func toGoNameConv(elem string) string {

	underscoreAtStart, underscoreAtEnd := countUnderscores(elem)

	elem = strings.Replace(elem, "_", " ", -1)
	elem = strings.Title(elem)
	elem = strings.Replace(elem, " ", "", -1)

	if underscoreAtEnd > 0 {
		elem += strings.Repeat("_", underscoreAtEnd)
	}

	if underscoreAtStart > 0 { // This is because Go doesn't support _ at the beginning of the element.
		elem = "U_" + elem
	}

	return elem
}

// --------------------------------------------------------------------
func castIfNeeded(elem string) string {
	if strings.Contains(elem, "int") {
		return "int(" + elem + ")"
	}
	return elem
}

// --------------------------------------------------------------------
func getNumericTypes() (numericTypes []string) {
	return []string{"Handle", "float64", "float32", "int8", "int16", "int32", "int64", "uint8", "uint16", "uint32", "uint64"}
}

// --------------------------------------------------------------------
func makeMetaFFIType(t string) string {
	return "metaffi_" + strings.ToLower(t)
}

// --------------------------------------------------------------------
func getMetaFFIStringTypes() (numericTypes []string) {
	return []string{"string8"}
}

// --------------------------------------------------------------------
func getMetaFFIType(numericType string) (numericTypes uint64) {
	return IDL.TypeStringToTypeEnum[IDL.MetaFFIType(numericType)]
}

// --------------------------------------------------------------------
func getMetaFFIArrayType(numericType string) (numericTypes uint64) {
	return IDL.TypeStringToTypeEnum[IDL.MetaFFIType(numericType+"_array")]
}

//--------------------------------------------------------------------
