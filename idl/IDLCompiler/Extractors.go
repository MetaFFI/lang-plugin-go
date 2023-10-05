package IDLCompiler

import (
	"github.com/MetaFFI/plugin-sdk/compiler/go/IDL"
	"strings"
)

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
