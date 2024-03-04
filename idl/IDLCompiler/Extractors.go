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

	isArray := false
	if strings.Count(typename, "[]") > 0 {
		isArray = true
	}

	typename = strings.ReplaceAll(typename, "[]", "")
	typename = strings.ReplaceAll(typename, "*", "")

	switch typename {
	case "string":
		if !isArray {
			return IDL.STRING8
		} else {
			return IDL.STRING8_ARRAY
		}
	case "int":
		if !isArray {
			return IDL.INT64
		} else {
			return IDL.INT64_ARRAY
		}
	case "int8":
		if !isArray {
			return IDL.INT8
		} else {
			return IDL.INT8_ARRAY
		}
	case "int16":
		if !isArray {
			return IDL.INT16
		} else {
			return IDL.INT16_ARRAY
		}
	case "int32":
		if !isArray {
			return IDL.INT32
		} else {
			return IDL.INT32_ARRAY
		}
	case "int64":
		if !isArray {
			return IDL.INT64
		} else {
			return IDL.INT64_ARRAY
		}
	case "untyped int":
		if !isArray {
			return IDL.INT64
		} else {
			return IDL.INT64_ARRAY
		}
	case "uint":
		if !isArray {
			return IDL.UINT64
		} else {
			return IDL.UINT64_ARRAY
		}
	case "byte":
		fallthrough
	case "uint8":
		if !isArray {
			return IDL.UINT8
		} else {
			return IDL.UINT8_ARRAY
		}
	case "uint16":
		if !isArray {
			return IDL.UINT16
		} else {
			return IDL.UINT16_ARRAY
		}
	case "uint32":
		if !isArray {
			return IDL.UINT32
		} else {
			return IDL.UINT32_ARRAY
		}
	case "uint64":
		if !isArray {
			return IDL.UINT64
		} else {
			return IDL.UINT64_ARRAY
		}
	case "float32":
		if !isArray {
			return IDL.FLOAT32
		} else {
			return IDL.FLOAT32_ARRAY
		}
	case "float64":
		if !isArray {
			return IDL.FLOAT64
		} else {
			return IDL.FLOAT64_ARRAY
		}
	case "bool":
		if !isArray {
			return IDL.BOOL
		} else {
			return IDL.BOOL_ARRAY
		}
	case "any":
		fallthrough
	case "interface{}":
		if !isArray {
			return IDL.ANY
		} else {
			return IDL.ANY_ARRAY
		}

	default:
		if !isArray {
			return IDL.HANDLE
		} else {
			return IDL.HANDLE_ARRAY
		}
	}

}

//--------------------------------------------------------------------
