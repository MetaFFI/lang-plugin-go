package metaffi

import (
	"github.com/MetaFFI/plugin-sdk/compiler/go/IDL"
	"testing"
	"unsafe"
)

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

	output := FromCDTToGo(unsafe.Pointer(pcdts.pcdt), 0)

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
		t.Errorf("pcdt.type is not of type METAFFI_TYPE_INT8")
	}

	output := FromCDTToGo(unsafe.Pointer(pcdts.pcdt), 0)

	output_data, ok := output.([]interface{})
	if !ok {
		t.Errorf("output is not of type int8")
	}

	if len(output_data) != len(input) {
		t.Errorf("length of input and output are not equal. input: %v, output: %v", len(input), len(output_data))
	}

	for i, v := range input {
		if v != output_data[i] {
			t.Errorf("input and output are not equal. input: %v, output: %v", v, output_data[i])
		}
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

	output := FromCDTToGo(unsafe.Pointer(pcdts.pcdt), 0)

	if len(output.([]interface{})) != len(input) {
		t.Errorf("length of input and output are not equal. input: %v, output: %v", len(input), len(output.([]interface{})))
	}

	for i, v := range input {
		if len(v) != len(output.([]interface{})[i].([]interface{})) {
			t.Errorf("length of input and output are not equal. input: %v, output: %v", len(v), len(output.([]interface{})[i].([]interface{})))
		}
		for j, v2 := range v {
			if v2 != output.([]interface{})[i].([]interface{})[j].(int8) {
				t.Errorf("input and output are not equal. input: %v, output: %v", v2, output.([]interface{})[i].([]interface{})[j].(int8))
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

	output := FromCDTToGo(unsafe.Pointer(pcdts.pcdt), 0)

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

	output := FromCDTToGo(unsafe.Pointer(pcdts.pcdt), 0)

	output_data, ok := output.([]interface{})
	if !ok {
		t.Fatalf("output is not of type int8")
	}

	if len(output_data) != len(input) {
		t.Fatalf("length of input and output are not equal. input: %v, output: %v", len(input), len(output_data))

	}

	for i, v := range output.([]interface{}) {
		if v.(string) != input[i] {
			t.Errorf("input and output are not equal. output: %v, input: %v", v, input[i])
		}
	}
}

func TestGoCDT3DStringArray(t *testing.T) {
	pcdts := GetCDTS()

	input := [][][]string{
		{{"one", "two", "three"}, {"four", "five", "six"}, {"seven", "eight", "nine"}},
		{{"ten", "eleven", "twelve"}, {"thirteen", "fourteen", "fifteen"}, {"sixteen", "seventeen", "eighteen"}},
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

	output := FromCDTToGo(unsafe.Pointer(pcdts.pcdt), 0)

	output_data, ok := output.([]interface{})
	if !ok {
		t.Fatalf("output is not of type int8")
	}

	if len(output_data) != len(input) {
		t.Fatalf("length of input and output are not equal. input: %v, output: %v", len(input), len(output_data))
	}

	for i, v := range output.([]interface{}) {
		for j, v2 := range v.([]interface{}) {
			for k, v3 := range v2.([]interface{}) {
				if v3.(string) != input[i][j][k] {
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

	output := FromCDTToGo(unsafe.Pointer(pcdts.pcdt), 0)

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

	output := FromCDTToGo(unsafe.Pointer(pcdts.pcdt), 0)

	output_data, ok := output.(MetaFFIHandle)
	if !ok {
		t.Errorf("output is not of type int8")
	}

	if input != output_data {
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

	if GetCDTSType(pcdts, 0) != IDL.METAFFI_TYPE_HANDLE {
		t.Errorf("pcdt.type is not of type METAFFI_TYPE_HANDLE")
	}

	output := FromCDTToGo(unsafe.Pointer(pcdts.pcdt), 0)

	output_data, ok := output.([]interface{})
	if !ok {
		t.Errorf("output is not of type int8")
	}

	if input[0].(test) != output_data[0].(test) {
		t.Errorf("input and output are not equal. input: %v, output: %v", input[0].(test), output_data[0].(test))
	}

	if input[1].(MetaFFIHandle) != output_data[1].(MetaFFIHandle) {
		t.Errorf("input and output are not equal. input: %v, output: %v", input[1].(MetaFFIHandle), output_data[1].(MetaFFIHandle))
	}

	if input[2].(test) != output_data[2].(test) {
		t.Errorf("input and output are not equal. input: %v, output: %v", input[2].(test), output_data[2].(test))
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

	if GetCDTSType(pcdts, 0) != IDL.METAFFI_TYPE_HANDLE {
		t.Errorf("pcdt.type is not of type METAFFI_TYPE_HANDLE")
	}

	output := FromCDTToGo(unsafe.Pointer(pcdts.pcdt), 0)

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
					if d != v3.(MetaFFIHandle) {
						t.Errorf("input and output are not equal. input: %v, output: %v", d, v3)
					}
				}
			}
		}
	}
}
