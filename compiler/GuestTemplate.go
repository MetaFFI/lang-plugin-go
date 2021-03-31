package main

const GuestHeaderTemplate = `
// Code generated by OpenFFI. Modify only in marked places.
// Guest code for {{.IDLFilenameWithExtension}}
{{ $pfn := .IDLFilename}}

package main
`

const GuestImportsTemplate = `
import "github.com/golang/protobuf/proto"
import "fmt"
{{range $mindex, $i := .Imports}}
import "{{$i}}"{{end}}

`

const GuestCImport = `
/*
#include <stdint.h>
*/
import "C"
`

const GuestMainFunction = `
func main(){} // main function must be declared to create dynamic library
`

const GuestHelperFunctions = `
func errToOutError(out_err **C.char, out_err_len *C.uint64_t, is_error *C.uint8_t, customText string, err error){
	*is_error = 1
	txt := customText+err.Error()
	*out_err = C.CString(txt)
	*out_err_len = C.uint64_t(len(txt))
}

func panicHandler(out_err **C.char, out_err_len *C.uint64_t, is_error *C.uint8_t){
	
	if rec := recover(); rec != nil{
		fmt.Println("Caught Panic")

		msg := "Panic in Go function. Panic Data: "
		switch recType := rec.(type){
			case error: msg += (rec.(error)).Error()
			case string: msg += rec.(string)
			default: msg += fmt.Sprintf("Panic with type: %v - %v", recType, rec)
		}

		*is_error = 1
		*out_err = C.CString(msg)
		*out_err_len = C.uint64_t(len(msg))
	}
}
`

const GuestFunctionXLLRTemplate = `
// add functions
{{range $mindex, $m := .Modules}}

{{range $findex, $f := $m.Functions}}

// Call to foreign {{$f.PathToForeignFunction.function}}
//export EntryPoint_{{$f.PathToForeignFunction.function}}
func EntryPoint_{{$f.PathToForeignFunction.function}}(in_params *C.char, in_params_len C.uint64_t, out_params **C.char, out_params_len *C.uint64_t, out_ret **C.char, out_ret_len *C.uint64_t, is_error *C.uint8_t){

	// catch panics and return them as errors
	defer panicHandler(out_ret, out_ret_len, is_error)
	
	*is_error = 0

	// deserialize parameters
	inParams := C.GoStringN(in_params, C.int(in_params_len))
	req := {{$f.ParametersType}}{}
	err := proto.Unmarshal([]byte(inParams), &req)
	if err != nil{
		errToOutError(out_ret, out_ret_len, is_error, "Failed to unmarshal parameters", err)
		return
	}
	
	// call original function
	{{range $index, $elem := $f.ReturnValues}}{{if $index}},{{end}}{{$elem.Name}}{{end}}{{if $f.ReturnValues}} := {{end}}OpenFFI{{$f.PathToForeignFunction.function}}({{range $index, $elem := $f.Parameters}}{{if $index}},{{end}}{{if eq $elem.PassMethod "by_pointer"}}&{{end}}req.{{AsPublic $elem.Name}}{{end}})
	
	ret := {{$f.ParametersType}}{}

	// === fill out_ret
	// if one of the returned parameters is of interface type Error, check if error, and if so, return error
	{{range $index, $elem := $f.ReturnValues}}
	if err, isError := interface{}({{$elem.Name}}).(error); isError{
		errToOutError(out_ret, out_ret_len, is_error, "Error returned", err)
		return
	} else {
		ret.{{AsPublic $elem.Name}} = {{if eq $elem.PassMethod "by_pointer"}}&{{end}}{{$elem.Name}}
	}	
	{{end}}

	// serialize results
	serializedRet, err := proto.Marshal(&ret)
	if err != nil{
		errToOutError(out_ret, out_ret_len, is_error, "Failed to marshal return values into protobuf", err)
		return
	}

	// write serialized results to out_ret
	serializedRetStr := string(serializedRet)
	*out_ret = C.CString(serializedRetStr)
	*out_ret_len = C.uint64_t(len(serializedRetStr))

	// === fill out_params
	serializedParams, err := proto.Marshal(&req)
	if err != nil{
		errToOutError(out_ret, out_ret_len, is_error, "Failed to marshal parameter values into protobuf", err)
		return
	}
	
	if out_params != nil && out_params_len != nil{
		// write serialized parameters to out_params
		serializedParamsStr := string(serializedParams)
		*out_params = C.CString(serializedParamsStr)
		*out_params_len = C.uint64_t(len(serializedParamsStr))
	}
	
}

{{end}}

{{end}}

`

const GuestFunctionTemplate = `
{{range $mindex, $m := .Modules}}
{{range $findex, $f := $m.Functions}}
func OpenFFI_{{$f.PathToForeignFunction.function}}({{range $index, $elem := $f.Parameters}}{{if $index}}, {{end}}{{$elem.Name}} {{if eq $elem.PassMethod "by_pointer"}}*{{end}}{{if $elem.InnerTypes}}*{{end}}{{$elem.Type}}{{if eq $elem.Type "map"}}[{{$elem.MapKeyType}}]{{$elem.MapValueType}}{{end}}{{end}}){

	// Call original function as it is defined in the IDL. Modify to suit your needs.
	return {{$f.PathToForeignFunction.function}}({{range $index, $elem := $f.Parameters}}{{if $index}},{{end}}{{$elem.Name}}{{end}})
}
{{end}}{{end}}
`