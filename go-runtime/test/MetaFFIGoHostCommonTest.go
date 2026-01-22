package main

import "C"
import (
	"fmt"
	. "github.com/MetaFFI/lang-plugin-go/go-runtime"
	"github.com/MetaFFI/sdk/idl_entities/go/IDL"
	_ "net/http/pprof" // Import to register pprof handlers
	"os"
	"reflect"
	"runtime"
	"unsafe"
)

func main() {

	//TestGoCDTInt8()
	//TestGoCDTInt8Array()
	//TestGoCDTInt82DArray()
	//TestGoCDTString()
	//TestGoCDTStringArray()
	//TestGoCDT3DStringArray()
	//TestGoCDTHandleGoObject()
	//TestGoCDTHandleArray()
	//TestGoCDTHandleArraySameType()
	//TestGoCDTHandle3DArray()
	TestReturnErrWithNil()
	//TestGoToCDTMetaFFIHandle()

	runtime.GC()

}

func TestGoCDTInt8() {
	pcdts := GetCDTS()
	defer func() {
		FreeCDTS(pcdts)
	}()

	var input int8 = 123

	typeInfo := IDL.MetaFFITypeInfo{
		StringType: IDL.INT8,
		Alias:      "",
		Type:       IDL.METAFFI_TYPE_INT8,
		Dimensions: 0,
	}

	FromGoToCDT(input, unsafe.Pointer(pcdts.arr), typeInfo, 0)

	if GetCDTSType(pcdts, 0) != IDL.METAFFI_TYPE_INT8 {
		_, _ = fmt.Fprintf(os.Stderr, "pcdt.type is not of type METAFFI_TYPE_INT8")
	}

	output := FromCDTToGo(unsafe.Pointer(pcdts.arr), 0, nil)

	output_data, ok := output.(int8)
	if !ok {
		_, _ = fmt.Fprintf(os.Stderr, "output is not of type int8")
	}

	if input != output_data {
		_, _ = fmt.Fprintf(os.Stderr, "input and output are not equal. input: %v, output: %v", input, output_data)
	}
}

func TestGoCDTInt8Array() {
	pcdts := GetCDTS()
	defer func() {
		FreeCDTS(pcdts)
	}()

	var input []int8 = []int8{123, 42, 54}

	typeInfo := IDL.MetaFFITypeInfo{
		StringType: IDL.INT8_ARRAY,
		Alias:      "",
		Type:       IDL.METAFFI_TYPE_INT8_ARRAY,
		Dimensions: 1,
	}

	FromGoToCDT(input, unsafe.Pointer(pcdts.arr), typeInfo, 0)

	if GetCDTSType(pcdts, 0) != IDL.METAFFI_TYPE_INT8_ARRAY {
		panic(fmt.Sprintf("pcdt.type is not of type METAFFI_TYPE_INT8"))
	}

	output := FromCDTToGo(unsafe.Pointer(pcdts.arr), 0, nil).([]int8)

	if len(output) != len(input) {
		panic(fmt.Sprintf("length of input and output are not equal. input: %v, output: %v", len(input), len(output)))
	}

	for i, v := range input {
		if v != output[i] {
			_, _ = fmt.Fprintf(os.Stderr, "input and output are not equal. input: %v, output: %v", v, output[i])
		}
	}
}

func TestGoCDTInt82DArray() {
	pcdts := GetCDTS()
	defer func() {
		FreeCDTS(pcdts)
	}()

	var input [][]int8 = [][]int8{
		{1, 2, 3},
		{4, 5, 6},
		{7, 8, 9},
	}

	typeInfo := IDL.MetaFFITypeInfo{
		StringType: IDL.INT8_ARRAY,
		Alias:      "",
		Type:       IDL.METAFFI_TYPE_INT8_ARRAY,
		Dimensions: 2,
	}

	FromGoToCDT(input, unsafe.Pointer(pcdts.arr), typeInfo, 0)

	if GetCDTSType(pcdts, 0) != IDL.METAFFI_TYPE_INT8_ARRAY {
		panic(fmt.Sprintf("pcdt.type is not of type METAFFI_TYPE_INT8_ARRAY (%d), but of type %v", IDL.METAFFI_TYPE_INT8_ARRAY, GetCDTSType(pcdts, 0)))
	}

	output := FromCDTToGo(unsafe.Pointer(pcdts.arr), 0, nil).([][]int8)

	if len(output) != len(input) {
		panic(fmt.Sprintf("length of input and output are not equal. input: %v, output: %v", len(input), len(output)))
	}

	// deep compare output to input
	for i, v := range input {
		if len(v) != len(output[i]) {
			panic(fmt.Sprintf("length of input and output are not equal. input: %v, output: %v", len(v), len(output[i])))
		}
		for j, v2 := range v {
			if v2 != output[i][j] {
				_, _ = fmt.Fprintf(os.Stderr, "input and output are not equal. input: %v, output: %v", v2, output[i][j])
			}
		}
	}
}

func TestGoCDTString() {
	pcdts := GetCDTS()
	defer func() {
		FreeCDTS(pcdts)
	}()

	var input string = "Hello, World!"

	typeInfo := IDL.MetaFFITypeInfo{
		StringType: IDL.STRING8,
		Alias:      "",
		Type:       IDL.METAFFI_TYPE_STRING8,
		Dimensions: 0,
	}

	FromGoToCDT(input, unsafe.Pointer(pcdts.arr), typeInfo, 0)

	if GetCDTSType(pcdts, 0) != IDL.METAFFI_TYPE_STRING8 {
		panic(fmt.Sprintf("pcdt.type is not of type METAFFI_TYPE_INT8"))
	}

	output := FromCDTToGo(unsafe.Pointer(pcdts.arr), 0, nil)

	output_data, ok := output.(string)
	if !ok {
		panic(fmt.Sprintf("output is not of type int8"))
	}

	if input != output_data {
		panic(fmt.Sprintf("input and output are not equal. input: %v, output: %v", input, output_data))
	}
}

func TestGoCDTStringArray() {
	pcdts := GetCDTS()
	defer func() {
		FreeCDTS(pcdts)
	}()

	input := []string{"one", "two", "three"}

	typeInfo := IDL.MetaFFITypeInfo{
		StringType: IDL.STRING8_ARRAY,
		Alias:      "",
		Type:       IDL.METAFFI_TYPE_STRING8_ARRAY,
		Dimensions: 1,
	}

	FromGoToCDT(input, unsafe.Pointer(pcdts.arr), typeInfo, 0)

	if GetCDTSType(pcdts, 0) != IDL.METAFFI_TYPE_STRING8_ARRAY {
		panic(fmt.Sprintf("pcdt.type is not of type METAFFI_TYPE_STRING8_ARRAY"))
	}

	output := FromCDTToGo(unsafe.Pointer(pcdts.arr), 0, nil).([]string)

	if len(output) != len(input) {
		panic(fmt.Sprintf("length of input and output are not equal. input: %v, output: %v", len(input), len(output)))
	}

	for i, v := range output {
		if v != input[i] {
			_, _ = fmt.Fprintf(os.Stderr, "input and output are not equal. output: %v, input: %v", v, input[i])
		}
	}
}

func TestGoCDT3DStringArray() {
	pcdts := GetCDTS()
	defer func() {
		FreeCDTS(pcdts)
	}()

	input := [][][]string{
		{{"one"}, {"two1", "two2"}, {"three1", "three2", "three3"}},
		{{"four1", "four2", "four3", "four4"}, {"five1", "five2", "five3", "five4", "five5"}, {"six1", "six2", "six3", "six4", "six5", "six6"}},
	}

	typeInfo := IDL.MetaFFITypeInfo{
		StringType: IDL.STRING8_ARRAY,
		Alias:      "",
		Type:       IDL.METAFFI_TYPE_STRING8_ARRAY,
		Dimensions: 3,
	}

	FromGoToCDT(input, unsafe.Pointer(pcdts.arr), typeInfo, 0)

	if GetCDTSType(pcdts, 0) != IDL.METAFFI_TYPE_STRING8_ARRAY {
		panic(fmt.Sprintf("pcdt.type is not of type METAFFI_TYPE_STRING8_ARRAY"))
	}

	output := FromCDTToGo(unsafe.Pointer(pcdts.arr), 0, nil).([][][]string)

	if len(output) != len(input) {
		panic(fmt.Sprintf("length of input and output are not equal. input: %v, output: %v", len(input), len(output)))
	}

	for i, v := range output {
		for j, v2 := range v {
			for k, v3 := range v2 {
				if v3 != input[i][j][k] {
					_, _ = fmt.Fprintf(os.Stderr, "input and output are not equal. output: %v, input: %v", v3, input[i][j][k])
				}
			}

		}
	}
}

func TestGoCDTHandleGoObject() {
	pcdts := GetCDTS()
	defer func() {
		FreeCDTS(pcdts)
	}()

	type test struct {
		A int
	}

	input := test{A: 26}

	typeInfo := IDL.MetaFFITypeInfo{
		StringType: IDL.HANDLE,
		Alias:      "",
		Type:       IDL.METAFFI_TYPE_HANDLE,
		Dimensions: 0,
	}

	FromGoToCDT(input, unsafe.Pointer(pcdts.arr), typeInfo, 0)

	if GetCDTSType(pcdts, 0) != IDL.METAFFI_TYPE_HANDLE {
		panic(fmt.Sprintf("pcdt.type is not of type METAFFI_TYPE_HANDLE"))
	}

	output := FromCDTToGo(unsafe.Pointer(pcdts.arr), 0, nil)

	output_data, ok := output.(test)
	if !ok {
		panic(fmt.Sprintf("output is not of type test"))
	}

	if input != output_data {
		panic(fmt.Sprintf("input and output are not equal. input: %v, output: %v", input, output_data))
	}
}

func TestGoToCDTMetaFFIHandle() {
	pcdts := GetCDTS()
	defer func() {
		FreeCDTS(pcdts)
	}()

	type test struct {
		A int
	}

	input := MetaFFIHandle{Val: Handle(unsafe.Pointer(uintptr(123))), RuntimeID: 1010}

	typeInfo := IDL.MetaFFITypeInfo{
		StringType: IDL.HANDLE,
		Alias:      "",
		Type:       IDL.METAFFI_TYPE_HANDLE,
		Dimensions: 0,
	}

	FromGoToCDT(input, unsafe.Pointer(pcdts.arr), typeInfo, 0)

	if GetCDTSType(pcdts, 0) != IDL.METAFFI_TYPE_HANDLE {
		panic("pcdt.type is not of type METAFFI_TYPE_HANDLE")
	}

	output := FromCDTToGo(unsafe.Pointer(pcdts.arr), 0, nil)

	output_data, ok := output.(MetaFFIHandle)
	if !ok {
		panic("output is not of type int8")
	}

	if input.Val != output_data.Val || input.RuntimeID != output_data.RuntimeID {
		panic("input and output are not equal. input: %v, output: %v")
	}
}

func TestGoCDTHandleArray() {
	pcdts := GetCDTS()
	defer func() {
		FreeCDTS(pcdts)
	}()

	type test struct {
		A int
	}

	input := []interface{}{test{A: 26}, MetaFFIHandle{
		Val:       Handle(unsafe.Pointer(uintptr(123))),
		RuntimeID: 2020,
	}, test{A: 27}}

	typeInfo := IDL.MetaFFITypeInfo{
		StringType: IDL.HANDLE,
		Alias:      "",
		Type:       IDL.METAFFI_TYPE_HANDLE_ARRAY,
		Dimensions: 1,
	}

	FromGoToCDT(input, unsafe.Pointer(pcdts.arr), typeInfo, 0)

	if GetCDTSType(pcdts, 0)&IDL.METAFFI_TYPE_HANDLE_ARRAY == 0 {
		panic(fmt.Sprintf("pcdt.type is not of type METAFFI_TYPE_HANDLE_ARRAY. Returned: %v", GetCDTSType(pcdts, 0)))
	}

	output := FromCDTToGo(unsafe.Pointer(pcdts.arr), 0, nil)

	output_data, ok := output.([]interface{})
	if !ok {
		panic(fmt.Sprintf("output is not of type []interface{}"))
	}

	if input[0].(test) != output_data[0].(test) {
		panic(fmt.Sprintf("input and output are not equal. input: %v, output: %v", input[0].(test), output_data[0].(test)))
	}

	if input[1].(MetaFFIHandle).Val != output_data[1].(MetaFFIHandle).Val || input[1].(MetaFFIHandle).RuntimeID != output_data[1].(MetaFFIHandle).RuntimeID {
		panic(fmt.Sprintf("input and output are not equal. input: %v, output: %v", input[1].(MetaFFIHandle), output_data[1].(MetaFFIHandle)))
	}

	if input[2].(test) != output_data[2].(test) {
		panic(fmt.Sprintf("input and output are not equal. input: %v, output: %v", input[2].(test), output_data[2].(test)))
	}
}

func TestGoCDTHandleArraySameType() {
	pcdts := GetCDTS()
	defer func() {
		FreeCDTS(pcdts)
	}()

	type test struct {
		A int
	}

	input := []test{{A: 26}, {A: 27}}

	typeInfo := IDL.MetaFFITypeInfo{
		StringType: IDL.HANDLE,
		Alias:      "",
		Type:       IDL.METAFFI_TYPE_HANDLE_ARRAY,
		Dimensions: 1,
	}

	FromGoToCDT(input, unsafe.Pointer(pcdts.arr), typeInfo, 0)

	if GetCDTSType(pcdts, 0) != IDL.METAFFI_TYPE_HANDLE_ARRAY {
		panic(fmt.Sprintf("pcdt.type is not of type METAFFI_TYPE_HANDLE_ARRAY"))
	}

	output := FromCDTToGo(unsafe.Pointer(pcdts.arr), 0, reflect.TypeFor[test]()).([]test)

	if input[0] != output[0] {
		panic(fmt.Sprintf("input and output are not equal. input: %v, output: %v", input[0], output[0]))
	}

	if input[1] != output[1] {
		panic(fmt.Sprintf("input and output are not equal. input: %v, output: %v", input[1], output[1]))
	}
}

func TestGoCDTHandle3DArray() {
	pcdts := GetCDTS()
	defer func() {
		FreeCDTS(pcdts)
	}()

	type test struct {
		A int
	}

	input := [][][]interface{}{
		{{test{A: 26}, MetaFFIHandle{
			Val:       Handle(unsafe.Pointer(uintptr(123))),
			RuntimeID: 2020,
		}, test{A: 27}}},
		{{test{A: 28}, MetaFFIHandle{
			Val:       Handle(unsafe.Pointer(uintptr(124))),
			RuntimeID: 2021,
		}, test{A: 29}}},
		{{test{A: 30}, MetaFFIHandle{
			Val:       Handle(unsafe.Pointer(uintptr(125))),
			RuntimeID: 2022,
		}, test{A: 31}}},
	}

	typeInfo := IDL.MetaFFITypeInfo{
		StringType: IDL.HANDLE,
		Alias:      "",
		Type:       IDL.METAFFI_TYPE_HANDLE_ARRAY,
		Dimensions: 3,
	}

	FromGoToCDT(input, unsafe.Pointer(pcdts.arr), typeInfo, 0)

	if GetCDTSType(pcdts, 0) != IDL.METAFFI_TYPE_HANDLE_ARRAY {
		panic(fmt.Sprintf("pcdt.type is not of type METAFFI_TYPE_HANDLE_ARRAY"))
	}

	output := FromCDTToGo(unsafe.Pointer(pcdts.arr), 0, nil)

	_, ok := output.([][][]interface{})
	if !ok {
		panic(fmt.Sprintf("output is not of type []interface, but of type %v", output))
	}

	for i, v := range output.([][][]interface{}) {
		for j, v2 := range v {
			for k, v3 := range v2 {
				if d, ok := input[i][j][k].(test); ok {
					if d != v3.(test) {
						panic(fmt.Sprintf("input and output are not equal. input: %v, output: %v", d, v3))
					}
				} else if d, ok := input[i][j][k].(MetaFFIHandle); ok {
					if d.Val != v3.(MetaFFIHandle).Val || d.RuntimeID != v3.(MetaFFIHandle).RuntimeID {
						panic(fmt.Sprintf("input and output are not equal. input: %v, output: %v", d, v3))
					}
				} else {
					_, _ = fmt.Fprintf(os.Stderr, "input and output are not equal. input: %v, output: %v", input[i][j][k], v3)
				}
			}
		}
	}
}

func TestReturnErrWithNil() {
	var err error

	pcdts := GetCDTS()
	defer func() {
		FreeCDTS(pcdts)
	}()

	typeInfo := IDL.MetaFFITypeInfo{
		StringType: IDL.HANDLE,
		Alias:      "error",
		Type:       IDL.METAFFI_TYPE_HANDLE,
		Dimensions: 0,
	}
	FromGoToCDT(err, unsafe.Pointer(pcdts.arr), typeInfo, 0)

	cdttype := GetCDTSType(pcdts, 0)

	if cdttype != IDL.METAFFI_TYPE_HANDLE {
		panic(fmt.Sprintf("pcdt.type is not of type METAFFI_TYPE_HANDLE"))
	}

	pcdt := GetCDT(pcdts, 0)
	if GetCDTHandleValue(pcdt) != uintptr(0) {
		panic(fmt.Sprintf("pcdt.handle is not nil"))
	}
}
