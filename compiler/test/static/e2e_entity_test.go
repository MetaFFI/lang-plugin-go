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
	_, _ = Return_int64()
	_, _ = Get_G_Name_Getter()
	_ = Set_G_Name_Setter("")
	_, _ = NewTestHandle()
	t.Log("expected symbols present")
}

// --- Primitives return ---

func TestReturnInt8(t *testing.T) {
	v, err := Return_int8()
	if err != nil {
		t.Fatal(err)
	}
	if v != 42 {
		t.Errorf("Return_int8: got %v, want 42", v)
	}
}

func TestReturnInt16(t *testing.T) {
	v, err := Return_int16()
	if err != nil {
		t.Fatal(err)
	}
	if v != 1000 {
		t.Errorf("Return_int16: got %v, want 1000", v)
	}
}

func TestReturnInt32(t *testing.T) {
	v, err := Return_int32()
	if err != nil {
		t.Fatal(err)
	}
	if v != 100000 {
		t.Errorf("Return_int32: got %v, want 100000", v)
	}
}

func TestReturnInt64(t *testing.T) {
	v, err := Return_int64()
	if err != nil {
		t.Fatal(err)
	}
	if v != 9223372036854775807 {
		t.Errorf("Return_int64: got %v, want 9223372036854775807", v)
	}
}

func TestReturnUint8(t *testing.T) {
	v, err := Return_uint8()
	if err != nil {
		t.Fatal(err)
	}
	if v != 255 {
		t.Errorf("Return_uint8: got %v, want 255", v)
	}
}

func TestReturnUint16(t *testing.T) {
	v, err := Return_uint16()
	if err != nil {
		t.Fatal(err)
	}
	if v != 65535 {
		t.Errorf("Return_uint16: got %v, want 65535", v)
	}
}

func TestReturnUint32(t *testing.T) {
	v, err := Return_uint32()
	if err != nil {
		t.Fatal(err)
	}
	if v != 4294967295 {
		t.Errorf("Return_uint32: got %v, want 4294967295", v)
	}
}

func TestReturnUint64(t *testing.T) {
	v, err := Return_uint64()
	if err != nil {
		t.Fatal(err)
	}
	if v != 18446744073709551615 {
		t.Errorf("Return_uint64: got %v", v)
	}
}

func TestReturnFloat32(t *testing.T) {
	v, err := Return_float32()
	if err != nil {
		t.Fatal(err)
	}
	if math.Abs(float64(v)-3.14159) > 1e-5 {
		t.Errorf("Return_float32: got %v", v)
	}
}

func TestReturnFloat64(t *testing.T) {
	v, err := Return_float64()
	if err != nil {
		t.Fatal(err)
	}
	if math.Abs(v-3.141592653589793) > 1e-10 {
		t.Errorf("Return_float64: got %v", v)
	}
}

func TestReturnBoolTrue(t *testing.T) {
	v, err := Return_bool_true()
	if err != nil {
		t.Fatal(err)
	}
	if !v {
		t.Error("Return_bool_true: got false, want true")
	}
}

func TestReturnBoolFalse(t *testing.T) {
	v, err := Return_bool_false()
	if err != nil {
		t.Fatal(err)
	}
	if v {
		t.Error("Return_bool_false: got true, want false")
	}
}

func TestReturnString8(t *testing.T) {
	v, err := Return_string8()
	if err != nil {
		t.Fatal(err)
	}
	if v != "Hello from test plugin" {
		t.Errorf("Return_string8: got %q", v)
	}
}

func TestReturnNull(t *testing.T) {
	_, err := Return_null()
	if err != nil {
		t.Fatal(err)
	}
}

// --- Primitives accept ---

func TestAcceptInt8(t *testing.T) {
	if err := Accept_int8(42); err != nil {
		t.Fatal(err)
	}
}

func TestAcceptInt64(t *testing.T) {
	if err := Accept_int64(12345678901234); err != nil {
		t.Fatal(err)
	}
}

func TestAcceptFloat64(t *testing.T) {
	if err := Accept_float64(3.14159); err != nil {
		t.Fatal(err)
	}
}

func TestAcceptBool(t *testing.T) {
	if err := Accept_bool(true); err != nil {
		t.Fatal(err)
	}
	if err := Accept_bool(false); err != nil {
		t.Fatal(err)
	}
}

func TestAcceptString8(t *testing.T) {
	if err := Accept_string8("test string"); err != nil {
		t.Fatal(err)
	}
}

// --- Echo ---

func TestEchoInt64(t *testing.T) {
	for _, in := range []int64{123, -456, 0} {
		out, err := Echo_int64(in)
		if err != nil {
			t.Fatal(err)
		}
		if out != in {
			t.Errorf("Echo_int64(%d): got %d", in, out)
		}
	}
}

func TestEchoFloat64(t *testing.T) {
	for _, in := range []float64{3.14, -2.718} {
		out, err := Echo_float64(in)
		if err != nil {
			t.Fatal(err)
		}
		if math.Abs(out-in) > 1e-10 {
			t.Errorf("Echo_float64(%v): got %v", in, out)
		}
	}
}

func TestEchoString8(t *testing.T) {
	for _, in := range []string{"test", "", "hello world"} {
		out, err := Echo_string8(in)
		if err != nil {
			t.Fatal(err)
		}
		if out != in {
			t.Errorf("Echo_string8(%q): got %q", in, out)
		}
	}
}

func TestEchoBool(t *testing.T) {
	for _, in := range []bool{true, false} {
		out, err := Echo_bool(in)
		if err != nil {
			t.Fatal(err)
		}
		if out != in {
			t.Errorf("Echo_bool(%v): got %v", in, out)
		}
	}
}

// --- Arithmetic ---

func TestAddInt64(t *testing.T) {
	out, err := Add_int64(3, 4)
	if err != nil {
		t.Fatal(err)
	}
	if out != 7 {
		t.Errorf("Add_int64(3,4): got %v, want 7", out)
	}
	out, _ = Add_int64(-5, 10)
	if out != 5 {
		t.Errorf("Add_int64(-5,10): got %v, want 5", out)
	}
}

func TestAddFloat64(t *testing.T) {
	out, err := Add_float64(1.5, 2.5)
	if err != nil {
		t.Fatal(err)
	}
	if math.Abs(out-4.0) > 1e-10 {
		t.Errorf("Add_float64: got %v", out)
	}
}

func TestConcatStrings(t *testing.T) {
	out, err := Concat_strings("Hello", " World")
	if err != nil {
		t.Fatal(err)
	}
	if out != "Hello World" {
		t.Errorf("Concat_strings: got %q", out)
	}
}

// --- Arrays ---

func TestReturnInt64Array1D(t *testing.T) {
	out, err := Return_int64_array_1d()
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
	out, err := Return_int64_array_2d()
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
	out, err := Return_string_array()
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
	out, err := Sum_int64_array([]int64{1, 2, 3, 4, 5})
	if err != nil {
		t.Fatal(err)
	}
	if out != 15 {
		t.Errorf("Sum_int64_array: got %v, want 15", out)
	}
}

func TestEchoInt64Array(t *testing.T) {
	in := []int64{10, 20, 30}
	out, err := Echo_int64_array(in)
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
	out, err := Join_strings([]string{"a", "b", "c"})
	if err != nil {
		t.Fatal(err)
	}
	if out != "a, b, c" {
		t.Errorf("Join_strings: got %q, want \"a, b, c\"", out)
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
	if err := h.Append_to_data("_suffix"); err != nil {
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
	_, err := Throw_error()
	if err == nil {
		t.Fatal("Throw_error: expected error")
	}
	if !strings.Contains(err.Error(), "Test error thrown intentionally") {
		t.Errorf("Throw_error: message missing expected text: %v", err)
	}
}

func TestThrowWithMessage(t *testing.T) {
	_, err := Throw_with_message("Custom error message")
	if err == nil {
		t.Fatal("Throw_with_message: expected error")
	}
	if !strings.Contains(err.Error(), "Custom error message") {
		t.Errorf("Throw_with_message: %v", err)
	}
}

func TestErrorIfNegativePositive(t *testing.T) {
	if err := Error_if_negative(42); err != nil {
		t.Fatal(err)
	}
	if err := Error_if_negative(0); err != nil {
		t.Fatal(err)
	}
}

func TestErrorIfNegativeNegative(t *testing.T) {
	err := Error_if_negative(-1)
	if err == nil {
		t.Fatal("Error_if_negative(-1): expected error")
	}
}

// --- Multiple returns ---

func TestReturnTwoValues(t *testing.T) {
	n, s, err := Return_two_values()
	if err != nil {
		t.Fatal(err)
	}
	if n != 42 || s != "answer" {
		t.Errorf("Return_two_values: got (%v, %q), want (42, \"answer\")", n, s)
	}
}

func TestReturnThreeValues(t *testing.T) {
	a, b, c, err := Return_three_values()
	if err != nil {
		t.Fatal(err)
	}
	if a != 1 || math.Abs(b-2.5) > 1e-10 || !c {
		t.Errorf("Return_three_values: got (%v, %v, %v)", a, b, c)
	}
}

func TestSwapValues(t *testing.T) {
	s, n, err := Swap_values(123, "hello")
	if err != nil {
		t.Fatal(err)
	}
	if s != "hello" || n != 123 {
		t.Errorf("Swap_values: got (%q, %v), want (\"hello\", 123)", s, n)
	}
}

// --- Callables ---

func TestCallCallbackAdd(t *testing.T) {
	adder := func(a, b int64) int64 { return a + b }
	out, err := Call_callback_add(adder)
	if err != nil {
		t.Fatal(err)
	}
	if out != 7 {
		t.Errorf("Call_callback_add: got %v, want 7", out)
	}
}

func TestCallCallbackString(t *testing.T) {
	echo := func(s string) string { return s }
	out, err := Call_callback_string(echo)
	if err != nil {
		t.Fatal(err)
	}
	if out != "test" {
		t.Errorf("Call_callback_string: got %q, want \"test\"", out)
	}
}

func TestReturnAdderCallback(t *testing.T) {
	ret, err := Return_adder_callback()
	if err != nil {
		t.Fatal(err)
	}
	if ret == nil {
		t.Fatal("Return_adder_callback: got nil")
	}
	callable, ok := ret.(*MetaFFICallable)
	if !ok {
		t.Fatalf("Return_adder_callback: expected *MetaFFICallable, got %T", ret)
	}
	result, err := callable.Call(int64(10), int64(20))
	if err != nil {
		t.Fatalf("callable.Call(10, 20): %v", err)
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
	out, err := Accept_any(int64(42))
	if err != nil {
		t.Fatal(err)
	}
	if v, ok := out.(int64); !ok || v != 142 {
		t.Errorf("Accept_any(42): got %v (%T), want 142", out, out)
	}
}

func TestAcceptAnyFloat64(t *testing.T) {
	out, err := Accept_any(3.14)
	if err != nil {
		t.Fatal(err)
	}
	if v, ok := out.(float64); !ok || math.Abs(v-6.28) > 1e-10 {
		t.Errorf("Accept_any(3.14): got %v (%T), want 6.28", out, out)
	}
}

func TestAcceptAnyString(t *testing.T) {
	out, err := Accept_any("hello")
	if err != nil {
		t.Fatal(err)
	}
	if v, ok := out.(string); !ok || v != "echoed: hello" {
		t.Errorf("Accept_any(\"hello\"): got %v (%T), want \"echoed: hello\"", out, out)
	}
}
