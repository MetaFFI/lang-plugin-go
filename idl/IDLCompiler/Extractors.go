package IDLCompiler

import (
	"github.com/MetaFFI/plugin-sdk/compiler/go/IDL"
	"strings"
)

// Declare a global variable map
var primitives = map[string]string{
	"bool":       "true",
	"byte":       "true",
	"complex64":  "true",
	"complex128": "true",
	"float32":    "true",
	"float64":    "true",
	"int":        "true",
	"int8":       "true",
	"int16":      "true",
	"int32":      "true",
	"int64":      "true",
	"rune":       "true",
	"string":     "true",
	"uint":       "true",
	"uint8":      "true",
	"uint16":     "true",
	"uint32":     "true",
	"uint64":     "true",
	"uintptr":    "true",
}

// --------------------------------------------------------------------
func isPrimitiveType(typename string) bool {
	typename = strings.ReplaceAll(typename, "[]", "")
	typename = strings.ReplaceAll(typename, "*", "")

	_, isPrimitive := primitives[typename]
	return isPrimitive
}

// --------------------------------------------------------------------
func goTypeToMFFI(typename string) IDL.MetaFFIType {

	typename = strings.ReplaceAll(typename, "[]", "")
	typename = strings.ReplaceAll(typename, "*", "")

	switch typename {
	case "string":
		return IDL.STRING8
	case "int":
		return IDL.INT64
	case "int8":
		return IDL.INT8
	case "int16":
		return IDL.INT16
	case "int32":
		return IDL.INT32
	case "int64":
		return IDL.INT64
	case "untyped int":
		return IDL.INT64
	case "uint":
		return IDL.UINT64
	case "uint8":
		return IDL.UINT8
	case "uint16":
		return IDL.UINT16
	case "uint32":
		return IDL.UINT32
	case "uint64":
		return IDL.UINT64
	case "float32":
		return IDL.FLOAT32
	case "float64":
		return IDL.FLOAT64
	case "bool":
		return IDL.BOOL
	case "interface{}":
		return IDL.ANY

	default:
		return IDL.HANDLE
	}

}

//--------------------------------------------------------------------
