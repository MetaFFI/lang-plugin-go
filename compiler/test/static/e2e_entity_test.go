//go:build ignore

// Package test: entity E2E tests for generated Go host (xllr.test).
// This file is copied into output/test/ and run with the generated host_MetaFFIHost.go.
// Plugin path is taken from METAFFI_TEST_PLUGIN_PATH (set by orchestration test).
package test

import (
	"math"
	"os"
	"strings"
	"testing"
)

var pluginPath string

func TestMain(m *testing.M) {
	pluginPath = os.Getenv("METAFFI_TEST_PLUGIN_PATH")
	if pluginPath == "" {
		panic("METAFFI_TEST_PLUGIN_PATH not set")
	}
	BindModuleToCode(pluginPath)
	os.Exit(m.Run())
}

// --- Compiler generation checks ---

func TestCompilerGeneratesFile(t *testing.T) {
	t.Log("generated host built and linked successfully")
}

func TestGeneratedHasExpectedContent(t *testing.T) {
	_, _ = Return_Int64()
	_, _ = Get_G_Name_Getter()
	_ = Set_G_Name_Setter("")
	_, _ = NewTestHandle()
	t.Log("expected symbols present")
}

// --- Primitives return ---

func TestReturnInt8(t *testing.T) {
	v, err := Return_Int8()
	if err != nil {
		t.Fatal(err)
	}
	if v != 42 {
		t.Errorf("Return_Int8: got %v, want 42", v)
	}
}

func TestReturnInt16(t *testing.T) {
	v, err := Return_Int16()
	if err != nil {
		t.Fatal(err)
	}
	if v != 1000 {
		t.Errorf("Return_Int16: got %v, want 1000", v)
	}
}

func TestReturnInt32(t *testing.T) {
	v, err := Return_Int32()
	if err != nil {
		t.Fatal(err)
	}
	if v != 100000 {
		t.Errorf("Return_Int32: got %v, want 100000", v)
	}
}

func TestReturnInt64(t *testing.T) {
	v, err := Return_Int64()
	if err != nil {
		t.Fatal(err)
	}
	if v != 9223372036854775807 {
		t.Errorf("Return_Int64: got %v, want 9223372036854775807", v)
	}
}

func TestReturnUint8(t *testing.T) {
	v, err := Return_Uint8()
	if err != nil {
		t.Fatal(err)
	}
	if v != 255 {
		t.Errorf("Return_Uint8: got %v, want 255", v)
	}
}

func TestReturnUint16(t *testing.T) {
	v, err := Return_Uint16()
	if err != nil {
		t.Fatal(err)
	}
	if v != 65535 {
		t.Errorf("Return_Uint16: got %v, want 65535", v)
	}
}

func TestReturnUint32(t *testing.T) {
	v, err := Return_Uint32()
	if err != nil {
		t.Fatal(err)
	}
	if v != 4294967295 {
		t.Errorf("Return_Uint32: got %v, want 4294967295", v)
	}
}

func TestReturnUint64(t *testing.T) {
	v, err := Return_Uint64()
	if err != nil {
		t.Fatal(err)
	}
	if v != 18446744073709551615 {
		t.Errorf("Return_Uint64: got %v", v)
	}
}

func TestReturnFloat32(t *testing.T) {
	v, err := Return_Float32()
	if err != nil {
		t.Fatal(err)
	}
	if math.Abs(float64(v)-3.14159) > 1e-5 {
		t.Errorf("Return_Float32: got %v", v)
	}
}

func TestReturnFloat64(t *testing.T) {
	v, err := Return_Float64()
	if err != nil {
		t.Fatal(err)
	}
	if math.Abs(v-3.141592653589793) > 1e-10 {
		t.Errorf("Return_Float64: got %v", v)
	}
}

func TestReturnBoolTrue(t *testing.T) {
	v, err := Return_Bool_True()
	if err != nil {
		t.Fatal(err)
	}
	if !v {
		t.Error("Return_Bool_True: got false, want true")
	}
}

func TestReturnBoolFalse(t *testing.T) {
	v, err := Return_Bool_False()
	if err != nil {
		t.Fatal(err)
	}
	if v {
		t.Error("Return_Bool_False: got true, want false")
	}
}

func TestReturnString8(t *testing.T) {
	v, err := Return_String8()
	if err != nil {
		t.Fatal(err)
	}
	if v != "Hello from test plugin" {
		t.Errorf("Return_String8: got %q", v)
	}
}

func TestReturnNull(t *testing.T) {
	_, err := Return_Null()
	if err != nil {
		t.Fatal(err)
	}
}

// --- Primitives accept ---

func TestAcceptInt8(t *testing.T) {
	if err := Accept_Int8(42); err != nil {
		t.Fatal(err)
	}
}

func TestAcceptInt64(t *testing.T) {
	if err := Accept_Int64(12345678901234); err != nil {
		t.Fatal(err)
	}
}

func TestAcceptFloat64(t *testing.T) {
	if err := Accept_Float64(3.14159); err != nil {
		t.Fatal(err)
	}
}

func TestAcceptBool(t *testing.T) {
	if err := Accept_Bool(true); err != nil {
		t.Fatal(err)
	}
	if err := Accept_Bool(false); err != nil {
		t.Fatal(err)
	}
}

func TestAcceptString8(t *testing.T) {
	if err := Accept_String8("test string"); err != nil {
		t.Fatal(err)
	}
}

// --- Echo ---

func TestEchoInt64(t *testing.T) {
	for _, in := range []int64{123, -456, 0} {
		out, err := Echo_Int64(in)
		if err != nil {
			t.Fatal(err)
		}
		if out != in {
			t.Errorf("Echo_Int64(%d): got %d", in, out)
		}
	}
}

func TestEchoFloat64(t *testing.T) {
	for _, in := range []float64{3.14, -2.718} {
		out, err := Echo_Float64(in)
		if err != nil {
			t.Fatal(err)
		}
		if math.Abs(out-in) > 1e-10 {
			t.Errorf("Echo_Float64(%v): got %v", in, out)
		}
	}
}

func TestEchoString8(t *testing.T) {
	for _, in := range []string{"test", "", "hello world"} {
		out, err := Echo_String8(in)
		if err != nil {
			t.Fatal(err)
		}
		if out != in {
			t.Errorf("Echo_String8(%q): got %q", in, out)
		}
	}
}

func TestEchoBool(t *testing.T) {
	for _, in := range []bool{true, false} {
		out, err := Echo_Bool(in)
		if err != nil {
			t.Fatal(err)
		}
		if out != in {
			t.Errorf("Echo_Bool(%v): got %v", in, out)
		}
	}
}

// --- Arithmetic ---

func TestAddInt64(t *testing.T) {
	out, err := Add_Int64(3, 4)
	if err != nil {
		t.Fatal(err)
	}
	if out != 7 {
		t.Errorf("Add_Int64(3,4): got %v, want 7", out)
	}
	out, _ = Add_Int64(-5, 10)
	if out != 5 {
		t.Errorf("Add_Int64(-5,10): got %v, want 5", out)
	}
}

func TestAddFloat64(t *testing.T) {
	out, err := Add_Float64(1.5, 2.5)
	if err != nil {
		t.Fatal(err)
	}
	if math.Abs(out-4.0) > 1e-10 {
		t.Errorf("Add_Float64: got %v", out)
	}
}

func TestConcatStrings(t *testing.T) {
	out, err := Concat_Strings("Hello", " World")
	if err != nil {
		t.Fatal(err)
	}
	if out != "Hello World" {
		t.Errorf("Concat_Strings: got %q", out)
	}
}

// --- Arrays ---

func TestReturnInt64Array1D(t *testing.T) {
	out, err := Return_Int64_Array_1d()
	if err != nil {
		t.Fatal(err)
	}
	want := []int64{1, 2, 3}
	if len(out) != len(want) {
		t.Fatalf("len: got %d, want %d", len(out), len(want))
	}
	for i := range want {
		if out[i] != want[i] {
			t.Errorf("at %d: got %v, want %v", i, out[i], want[i])
		}
	}
}

func TestReturnInt64Array2D(t *testing.T) {
	out, err := Return_Int64_Array_2d()
	if err != nil {
		t.Fatal(err)
	}
	want := [][]int64{{1, 2}, {3, 4}}
	if len(out) != len(want) {
		t.Fatalf("len: got %d, want %d", len(out), len(want))
	}
	for i := range want {
		for j := range want[i] {
			if out[i][j] != want[i][j] {
				t.Errorf("at %d,%d: got %v, want %v", i, j, out[i][j], want[i][j])
			}
		}
	}
}

func TestReturnStringArray(t *testing.T) {
	out, err := Return_String_Array()
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"one", "two", "three"}
	for i := range want {
		if out[i] != want[i] {
			t.Errorf("at %d: got %q, want %q", i, out[i], want[i])
		}
	}
}

func TestSumInt64Array(t *testing.T) {
	out, err := Sum_Int64_Array([]int64{1, 2, 3, 4, 5})
	if err != nil {
		t.Fatal(err)
	}
	if out != 15 {
		t.Errorf("Sum_Int64_Array: got %v, want 15", out)
	}
}

func TestEchoInt64Array(t *testing.T) {
	in := []int64{10, 20, 30}
	out, err := Echo_Int64_Array(in)
	if err != nil {
		t.Fatal(err)
	}
	for i := range in {
		if out[i] != in[i] {
			t.Errorf("at %d: got %v, want %v", i, out[i], in[i])
		}
	}
}

func TestJoinStrings(t *testing.T) {
	out, err := Join_Strings([]string{"a", "b", "c"})
	if err != nil {
		t.Fatal(err)
	}
	if out != "a, b, c" {
		t.Errorf("Join_Strings: got %q, want \"a, b, c\"", out)
	}
}

// --- TestHandle ---

func TestTestHandleConstructor(t *testing.T) {
	h, err := NewTestHandle()
	if err != nil {
		t.Fatal(err)
	}
	if h == nil {
		t.Fatal("NewTestHandle returned nil")
	}
	id, err := h.Get_Id_Getter()
	if err != nil {
		t.Fatal(err)
	}
	if id < 1 {
		t.Errorf("id: got %v, want >= 1", id)
	}
	data, err := h.Get_Data_Getter()
	if err != nil {
		t.Fatal(err)
	}
	if data != "test_data" {
		t.Errorf("data: got %q, want \"test_data\"", data)
	}
}

func TestTestHandleDataSetter(t *testing.T) {
	h, err := NewTestHandle()
	if err != nil {
		t.Fatal(err)
	}
	if err := h.Set_Data_Setter("new_value"); err != nil {
		t.Fatal(err)
	}
	data, err := h.Get_Data_Getter()
	if err != nil {
		t.Fatal(err)
	}
	if data != "new_value" {
		t.Errorf("after set: got %q, want \"new_value\"", data)
	}
}

func TestTestHandleAppendToData(t *testing.T) {
	h, err := NewTestHandle()
	if err != nil {
		t.Fatal(err)
	}
	if err := h.Set_Data_Setter("base"); err != nil {
		t.Fatal(err)
	}
	if err := h.Append_To_Data("_suffix"); err != nil {
		t.Fatal(err)
	}
	data, err := h.Get_Data_Getter()
	if err != nil {
		t.Fatal(err)
	}
	if data != "base_suffix" {
		t.Errorf("after append: got %q, want \"base_suffix\"", data)
	}
}

// --- Global g_name ---

func TestGlobalGName(t *testing.T) {
	if err := Set_G_Name_Setter("test_value"); err != nil {
		t.Fatal(err)
	}
	v, err := Get_G_Name_Getter()
	if err != nil {
		t.Fatal(err)
	}
	if v != "test_value" {
		t.Errorf("get_g_name: got %q, want \"test_value\"", v)
	}
}

// --- Errors ---

func TestThrowError(t *testing.T) {
	err := Throw_Error()
	if err == nil {
		t.Fatal("Throw_Error: expected error")
	}
	if !strings.Contains(err.Error(), "Test error thrown intentionally") {
		t.Errorf("Throw_Error: message missing expected text: %v", err)
	}
}

func TestThrowWithMessage(t *testing.T) {
	err := Throw_With_Message("Custom error message")
	if err == nil {
		t.Fatal("Throw_With_Message: expected error")
	}
	if !strings.Contains(err.Error(), "Custom error message") {
		t.Errorf("Throw_With_Message: %v", err)
	}
}

func TestErrorIfNegativePositive(t *testing.T) {
	if err := Error_If_Negative(42); err != nil {
		t.Fatal(err)
	}
	if err := Error_If_Negative(0); err != nil {
		t.Fatal(err)
	}
}

func TestErrorIfNegativeNegative(t *testing.T) {
	err := Error_If_Negative(-1)
	if err == nil {
		t.Fatal("Error_If_Negative(-1): expected error")
	}
}

// --- Multiple returns ---

func TestReturnTwoValues(t *testing.T) {
	n, s, err := Return_Two_Values()
	if err != nil {
		t.Fatal(err)
	}
	if n != 42 || s != "answer" {
		t.Errorf("Return_Two_Values: got (%v, %q), want (42, \"answer\")", n, s)
	}
}

func TestReturnThreeValues(t *testing.T) {
	a, b, c, err := Return_Three_Values()
	if err != nil {
		t.Fatal(err)
	}
	if a != 1 || math.Abs(b-2.5) > 1e-10 || !c {
		t.Errorf("Return_Three_Values: got (%v, %v, %v)", a, b, c)
	}
}

func TestSwapValues(t *testing.T) {
	s, n, err := Swap_Values(123, "hello")
	if err != nil {
		t.Fatal(err)
	}
	if s != "hello" || n != 123 {
		t.Errorf("Swap_Values: got (%q, %v), want (\"hello\", 123)", s, n)
	}
}

// --- Callables ---

func TestCallCallbackAdd(t *testing.T) {
	adder := func(a, b int64) int64 { return a + b }
	out, err := Call_Callback_Add(adder)
	if err != nil {
		t.Fatal(err)
	}
	if out != 7 {
		t.Errorf("Call_Callback_Add: got %v, want 7", out)
	}
}

func TestCallCallbackString(t *testing.T) {
	echo := func(s string) string { return s }
	out, err := Call_Callback_String(echo)
	if err != nil {
		t.Fatal(err)
	}
	if out != "test" {
		t.Errorf("Call_Callback_String: got %q, want \"test\"", out)
	}
}

func TestReturnAdderCallback(t *testing.T) {
	ret, err := Return_Adder_Callback()
	if err != nil {
		t.Fatal(err)
	}
	if ret == nil {
		t.Fatal("Return_Adder_Callback: got nil")
	}

	var result []interface{}
	switch callable := ret.(type) {
	case func(...interface{}) ([]interface{}, error):
		result, err = callable(int64(10), int64(20))
	case interface {
		Call(...interface{}) ([]interface{}, error)
	}:
		result, err = callable.Call(int64(10), int64(20))
	default:
		t.Fatalf("Return_Adder_Callback: unsupported callable type %T", ret)
	}
	if err != nil {
		t.Fatalf("callable invocation failed: %v", err)
	}
	if len(result) != 1 {
		t.Fatalf("callable.Call: expected 1 return, got %d", len(result))
	}
	if v, ok := result[0].(int64); !ok || v != 30 {
		t.Fatalf("callable.Call(10, 20): got %v (%T), want 30", result[0], result[0])
	}
}

// --- Any type ---

func TestAcceptAnyInt64(t *testing.T) {
	out, err := Accept_Any(int64(42))
	if err != nil {
		t.Fatal(err)
	}
	if v, ok := out.(int64); !ok || v != 142 {
		t.Errorf("Accept_Any(42): got %v (%T), want 142", out, out)
	}
}

func TestAcceptAnyFloat64(t *testing.T) {
	out, err := Accept_Any(3.14)
	if err != nil {
		t.Fatal(err)
	}
	if v, ok := out.(float64); !ok || math.Abs(v-6.28) > 1e-10 {
		t.Errorf("Accept_Any(3.14): got %v (%T), want 6.28", out, out)
	}
}

func TestAcceptAnyString(t *testing.T) {
	out, err := Accept_Any("hello")
	if err != nil {
		t.Fatal(err)
	}
	if v, ok := out.(string); !ok || v != "echoed: hello" {
		t.Errorf("Accept_Any(\"hello\"): got %v (%T), want \"echoed: hello\"", out, out)
	}
}
