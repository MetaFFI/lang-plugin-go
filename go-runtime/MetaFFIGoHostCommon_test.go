package metaffi

import (
	"github.com/MetaFFI/plugin-sdk/compiler/go/IDL"
	"reflect"
	"testing"
	"unsafe"
)

func TestFloat32Array(t *testing.T) {

	input := Get3DFloat32ArrayCDTS()

	res := FromCDTToGo(unsafe.Pointer(input.pcdt), 0, nil)

	if GetCDTSType(input, 0) != IDL.METAFFI_TYPE_FLOAT32_ARRAY {
		t.Fatalf("pcdt.type is not of type METAFFI_TYPE_FLOAT32_ARRAY")
	}

	output_data := res.([][][]float32)

	// Check the outer array length
	if len(output_data) != 3 {
		t.Fatalf("Outer array lengths do not match. Expected: 3, Got: %v", len(output_data))
	}

	// Define the expected 3D array
	expected := [][][]float32{
		{{1.0, 2.0, 3.0}, {4.0, 5.0, 6.0}, {7.0, 8.0, 9.0}},
		{{10.0, 11.0, 12.0}, {13.0, 14.0, 15.0}, {16.0, 17.0, 18.0}},
		{{19.0, 20.0, 21.0}, {22.0, 23.0, 24.0}, {25.0, 26.0, 27.0}},
	}

	// Compare output_data with expected data
	for i, outer := range output_data {
		for j, middle := range outer {
			for k, val := range middle {
				if val != expected[i][j][k] {
					t.Errorf("Values at index [%v][%v][%v] do not match. Expected: %v, Got: %v", i, j, k, expected[i][j][k], val)
				}
			}
		}
	}
}

func TestGoCDTInt8(t *testing.T) {
	pcdts := GetCDTS()

	var input int8 = 123

	typeInfo := IDL.MetaFFITypeInfo{
		StringType: IDL.INT8,
		Alias:      "",
		Type:       IDL.METAFFI_TYPE_INT8,
		Dimensions: 0,
	}

	FromGoToCDT(input, unsafe.Pointer(pcdts.pcdt), typeInfo, 0)

	if GetCDTSType(pcdts, 0) != IDL.METAFFI_TYPE_INT8 {
		t.Errorf("pcdt.type is not of type METAFFI_TYPE_INT8")
	}

	output := FromCDTToGo(unsafe.Pointer(pcdts.pcdt), 0, nil)

	output_data, ok := output.(int8)
	if !ok {
		t.Errorf("output is not of type int8")
	}

	if input != output_data {
		t.Errorf("input and output are not equal. input: %v, output: %v", input, output_data)
	}
}

func TestGoCDTInt8Array(t *testing.T) {
	pcdts := GetCDTS()

	var input []int8 = []int8{123, 42, 54}

	typeInfo := IDL.MetaFFITypeInfo{
		StringType: IDL.INT8_ARRAY,
		Alias:      "",
		Type:       IDL.METAFFI_TYPE_INT8_ARRAY,
		Dimensions: 1,
	}

	FromGoToCDT(input, unsafe.Pointer(pcdts.pcdt), typeInfo, 0)

	if GetCDTSType(pcdts, 0) != IDL.METAFFI_TYPE_INT8_ARRAY {
		t.Fatalf("pcdt.type is not of type METAFFI_TYPE_INT8")
	}

	output := FromCDTToGo(unsafe.Pointer(pcdts.pcdt), 0, nil).([]int8)

	if len(output) != len(input) {
		t.Fatalf("length of input and output are not equal. input: %v, output: %v", len(input), len(output))
	}

	for i, v := range input {
		if v != output[i] {
			t.Errorf("input and output are not equal. input: %v, output: %v", v, output[i])
		}
	}
}

func TestGoCDTUInt8CArray(t *testing.T) {

	input := Get2DUInt8ArrayCDTS()
	output := FromCDTToGo(unsafe.Pointer(input.pcdt), 0, nil).([][]uint8)

	if len(output) != 3 {
		t.Errorf("length of input and output are not equal. input: %v, output: %v", 3, len(output))
	}

	if output[0][0] != 0 {
		t.Errorf("input and output are not equal. input: %v, output: %v", 0, output[0][0])
	}

	if output[0][1] != 1 {
		t.Errorf("input and output are not equal. input: %v, output: %v", 1, output[0][1])
	}

	if output[0][2] != 2 {
		t.Errorf("input and output are not equal. input: %v, output: %v", 2, output[0][2])
	}

	if output[1][0] != 3 {
		t.Errorf("input and output are not equal. input: %v, output: %v", 3, output[1][0])
	}

	if output[1][1] != 4 {
		t.Errorf("input and output are not equal. input: %v, output: %v", 4, output[1][1])
	}

	if output[1][2] != 5 {
		t.Errorf("input and output are not equal. input: %v, output: %v", 5, output[1][2])
	}

	if output[2][0] != 6 {
		t.Errorf("input and output are not equal. input: %v, output: %v", 6, output[2][0])
	}

	if output[2][1] != 7 {
		t.Errorf("input and output are not equal. input: %v, output: %v", 7, output[2][1])
	}

	if output[2][2] != 8 {
		t.Errorf("input and output are not equal. input: %v, output: %v", 8, output[2][2])
	}
}

func TestGoCDTInt82DArray(t *testing.T) {
	pcdts := GetCDTS()

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

	FromGoToCDT(input, unsafe.Pointer(pcdts.pcdt), typeInfo, 0)

	if GetCDTSType(pcdts, 0) != IDL.METAFFI_TYPE_INT8_ARRAY {
		t.Errorf("pcdt.type is not of type METAFFI_TYPE_INT8")
	}

	output := FromCDTToGo(unsafe.Pointer(pcdts.pcdt), 0, nil).([][]int8)

	if len(output) != len(input) {
		t.Fatalf("length of input and output are not equal. input: %v, output: %v", len(input), len(output))
	}

	// deep compare output to input
	for i, v := range input {
		if len(v) != len(output[i]) {
			t.Fatalf("length of input and output are not equal. input: %v, output: %v", len(v), len(output[i]))
		}
		for j, v2 := range v {
			if v2 != output[i][j] {
				t.Errorf("input and output are not equal. input: %v, output: %v", v2, output[i][j])
			}
		}
	}
}

func TestGoCDTString(t *testing.T) {
	pcdts := GetCDTS()

	var input string = "Hello, World!"

	typeInfo := IDL.MetaFFITypeInfo{
		StringType: IDL.STRING8,
		Alias:      "",
		Type:       IDL.METAFFI_TYPE_STRING8,
		Dimensions: 0,
	}

	FromGoToCDT(input, unsafe.Pointer(pcdts.pcdt), typeInfo, 0)

	if GetCDTSType(pcdts, 0) != IDL.METAFFI_TYPE_STRING8 {
		t.Errorf("pcdt.type is not of type METAFFI_TYPE_INT8")
	}

	output := FromCDTToGo(unsafe.Pointer(pcdts.pcdt), 0, nil)

	output_data, ok := output.(string)
	if !ok {
		t.Errorf("output is not of type int8")
	}

	if input != output_data {
		t.Errorf("input and output are not equal. input: %v, output: %v", input, output_data)
	}
}

func TestGoCDTStringArray(t *testing.T) {
	pcdts := GetCDTS()

	input := []string{"one", "two", "three"}

	typeInfo := IDL.MetaFFITypeInfo{
		StringType: IDL.STRING8_ARRAY,
		Alias:      "",
		Type:       IDL.METAFFI_TYPE_STRING8_ARRAY,
		Dimensions: 1,
	}

	FromGoToCDT(input, unsafe.Pointer(pcdts.pcdt), typeInfo, 0)

	if GetCDTSType(pcdts, 0) != IDL.METAFFI_TYPE_STRING8_ARRAY {
		t.Fatalf("pcdt.type is not of type METAFFI_TYPE_STRING8_ARRAY")
	}

	output := FromCDTToGo(unsafe.Pointer(pcdts.pcdt), 0, nil).([]string)

	if len(output) != len(input) {
		t.Fatalf("length of input and output are not equal. input: %v, output: %v", len(input), len(output))
	}

	for i, v := range output {
		if v != input[i] {
			t.Errorf("input and output are not equal. output: %v, input: %v", v, input[i])
		}
	}
}

func TestGoCDT3DStringArray(t *testing.T) {
	pcdts := GetCDTS()

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

	FromGoToCDT(input, unsafe.Pointer(pcdts.pcdt), typeInfo, 0)

	if GetCDTSType(pcdts, 0) != IDL.METAFFI_TYPE_STRING8_ARRAY {
		t.Fatalf("pcdt.type is not of type METAFFI_TYPE_STRING8_ARRAY")
	}

	output := FromCDTToGo(unsafe.Pointer(pcdts.pcdt), 0, nil).([][][]string)

	if len(output) != len(input) {
		t.Fatalf("length of input and output are not equal. input: %v, output: %v", len(input), len(output))
	}

	for i, v := range output {
		for j, v2 := range v {
			for k, v3 := range v2 {
				if v3 != input[i][j][k] {
					t.Errorf("input and output are not equal. output: %v, input: %v", v3, input[i][j][k])
				}
			}

		}
	}
}

func TestGoCDTHandleGoObject(t *testing.T) {
	pcdts := GetCDTS()

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

	FromGoToCDT(input, unsafe.Pointer(pcdts.pcdt), typeInfo, 0)

	if GetCDTSType(pcdts, 0) != IDL.METAFFI_TYPE_HANDLE {
		t.Errorf("pcdt.type is not of type METAFFI_TYPE_HANDLE")
	}

	output := FromCDTToGo(unsafe.Pointer(pcdts.pcdt), 0, nil)

	output_data, ok := output.(test)
	if !ok {
		t.Errorf("output is not of type int8")
	}

	if input != output_data {
		t.Errorf("input and output are not equal. input: %v, output: %v", input, output_data)
	}
}

func TestGoCDTHandleNonGoObject(t *testing.T) {
	pcdts := GetCDTS()

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

	FromGoToCDT(input, unsafe.Pointer(pcdts.pcdt), typeInfo, 0)

	if GetCDTSType(pcdts, 0) != IDL.METAFFI_TYPE_HANDLE {
		t.Errorf("pcdt.type is not of type METAFFI_TYPE_HANDLE")
	}

	output := FromCDTToGo(unsafe.Pointer(pcdts.pcdt), 0, nil)

	output_data, ok := output.(MetaFFIHandle)
	if !ok {
		t.Errorf("output is not of type int8")
	}

	if input.Val != output_data.Val || input.RuntimeID != output_data.RuntimeID {
		t.Errorf("input and output are not equal. input: %v, output: %v", input, output_data)
	}
}

func TestGoCDTHandleArray(t *testing.T) {
	pcdts := GetCDTS()

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

	FromGoToCDT(input, unsafe.Pointer(pcdts.pcdt), typeInfo, 0)

	if GetCDTSType(pcdts, 0) != IDL.METAFFI_TYPE_HANDLE_ARRAY {
		t.Errorf("pcdt.type is not of type METAFFI_TYPE_HANDLE_ARRAY")
	}

	output := FromCDTToGo(unsafe.Pointer(pcdts.pcdt), 0, nil)

	output_data, ok := output.([]interface{})
	if !ok {
		t.Errorf("output is not of type []interface{}")
	}

	if input[0].(test) != output_data[0].(test) {
		t.Errorf("input and output are not equal. input: %v, output: %v", input[0].(test), output_data[0].(test))
	}

	if input[1].(MetaFFIHandle).Val != output_data[1].(MetaFFIHandle).Val || input[1].(MetaFFIHandle).RuntimeID != output_data[1].(MetaFFIHandle).RuntimeID {
		t.Errorf("input and output are not equal. input: %v, output: %v", input[1].(MetaFFIHandle), output_data[1].(MetaFFIHandle))
	}

	if input[2].(test) != output_data[2].(test) {
		t.Errorf("input and output are not equal. input: %v, output: %v", input[2].(test), output_data[2].(test))
	}
}

func TestGoCDTHandleArraySameType(t *testing.T) {
	pcdts := GetCDTS()

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

	FromGoToCDT(input, unsafe.Pointer(pcdts.pcdt), typeInfo, 0)

	if GetCDTSType(pcdts, 0) != IDL.METAFFI_TYPE_HANDLE_ARRAY {
		t.Errorf("pcdt.type is not of type METAFFI_TYPE_HANDLE_ARRAY")
	}

	output := FromCDTToGo(unsafe.Pointer(pcdts.pcdt), 0, reflect.TypeFor[test]()).([]test)

	if input[0] != output[0] {
		t.Errorf("input and output are not equal. input: %v, output: %v", input[0], output[0])
	}

	if input[1] != output[1] {
		t.Errorf("input and output are not equal. input: %v, output: %v", input[1], output[1])
	}
}

func TestGoCDTHandle3DArray(t *testing.T) {
	pcdts := GetCDTS()

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

	FromGoToCDT(input, unsafe.Pointer(pcdts.pcdt), typeInfo, 0)

	if GetCDTSType(pcdts, 0) != IDL.METAFFI_TYPE_HANDLE_ARRAY {
		t.Errorf("pcdt.type is not of type METAFFI_TYPE_HANDLE_ARRAY")
	}

	output := FromCDTToGo(unsafe.Pointer(pcdts.pcdt), 0, nil)

	_, ok := output.([][][]interface{})
	if !ok {
		t.Fatalf("output is not of type []interface, but of type %v", output)
	}

	for i, v := range output.([][][]interface{}) {
		for j, v2 := range v {
			for k, v3 := range v2 {
				if d, ok := input[i][j][k].(test); ok {
					if d != v3.(test) {
						t.Errorf("input and output are not equal. input: %v, output: %v", d, v3)
					}
				} else if d, ok := input[i][j][k].(MetaFFIHandle); ok {
					if d.Val != v3.(MetaFFIHandle).Val || d.RuntimeID != v3.(MetaFFIHandle).RuntimeID {
						t.Errorf("input and output are not equal. input: %v, output: %v", d, v3)
					}
				} else {
					t.Errorf("input and output are not equal. input: %v, output: %v", input[i][j][k], v3)
				}
			}
		}
	}
}
