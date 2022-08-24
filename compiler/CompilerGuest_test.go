package main

import (
	compiler "github.com/MetaFFI/plugin-sdk/compiler/go/IDL"
	"io/ioutil"
	"os"
	"testing"
)

const idl_guest = `{"idl_filename": "test","idl_extension": ".json","idl_filename_with_extension": "test.json", "target_language": "test", "idl_full_path": "","modules": [{"name": "TestModule","comment": "Comments for TestModule\n","tags": null,"functions": [{"name": "F1","comment": "F1 comment\nparam1 comment\n","function_path": {"module": "$PWD/temp","package": "GoFuncs","function": "F1"},"parameter_type": "Params1","return_values_type": "Return1","parameters": [{"name": "p1","type": "float64","comment": "= 3.141592","tags": null,"dimensions": 0,"pass_method": ""},{"name": "p2","type": "float32","comment": "= 2.71","tags": null,"dimensions": 0,"pass_method": ""},{"name": "p3","type": "int8","comment": "= -10","tags": null,"dimensions": 0,"pass_method": ""},{"name": "p4","type": "int16","comment": "= -20","tags": null,"dimensions": 0,"pass_method": ""},{"name": "p5","type": "int32","comment": "= -30","tags": null,"dimensions": 0,"pass_method": ""},{"name": "p6","type": "int64","comment": "= -40","tags": null,"dimensions": 0,"pass_method": ""},{"name": "p7","type": "uint8","comment": "= 50","tags": null,"dimensions": 0,"pass_method": ""},{"name": "p8","type": "uint16","comment": "= 60","tags": null,"dimensions": 0,"pass_method": ""},{"name": "p9","type": "uint32","comment": "= 70","tags": null,"dimensions": 0,"pass_method": ""},{"name": "p10","type": "uint64","comment": "= 80","tags": null,"dimensions": 0,"pass_method": ""},{"name": "p11","type": "bool","comment": "= true","tags": null,"dimensions": 0,"pass_method": ""},{"name": "p12","type": "string8","comment": "= This is an input","tags": null,"dimensions": 0,"pass_method": ""},{"name": "p13","type": "string8","comment": "= {element one, element two}","tags": null,"dimensions": 1,"pass_method": ""},{"name": "p14","type": "uint8","comment": "= {2, 4, 6, 8, 10}","tags": null,"dimensions": 1,"pass_method": ""}],"return_values": [{"name": "r1","type": "float64","comment": "= 0.57721","tags": null,"dimensions": 0,"pass_method": ""},{"name": "r2","type": "float32","comment": "= 3.359","tags": null,"dimensions": 0,"pass_method": ""},{"name": "r3","type": "int8","comment": "= -11","tags": null,"dimensions": 0,"pass_method": ""},{"name": "r4","type": "int16","comment": "= -21","tags": null,"dimensions": 0,"pass_method": ""},{"name": "r5","type": "int32","comment": "= -31","tags": null,"dimensions": 0,"pass_method": ""},{"name": "r6","type": "int64","comment": "= -41","tags": null,"dimensions": 0,"pass_method": ""},{"name": "r7","type": "uint8","comment": "= 51","tags": null,"dimensions": 0,"pass_method": ""},{"name": "r8","type": "uint16","comment": "= 61","tags": null,"dimensions": 0,"pass_method": ""},{"name": "r9","type": "uint32","comment": "= 71","tags": null,"dimensions": 0,"pass_method": ""},{"name": "r10","type": "uint64","comment": "= 81","tags": null,"dimensions": 0,"pass_method": ""},{"name": "r11","type": "bool","comment": "= true","tags": null,"dimensions": 0,"pass_method": ""},{"name": "r12","type": "string8","comment": "= This is an output","tags": null,"dimensions": 0,"pass_method": ""},{"name": "r13","type": "string8","comment": "= {return one, return two}","tags": null,"dimensions": 1,"pass_method": ""},{"name": "r14","type": "uint8","comment": "= {20, 40, 60, 80, 100}","tags": null,"dimensions": 1,"pass_method": ""}]}]}]}`

const GuestCode = `
package GoFuncs

import "fmt"

func F1(p1 float64, p2 float32, p3 int8, p4 int16, p5 int32, p6 int64, p7 uint8, p8 uint16, p9 uint32, p10 uint64, p11 bool, p12 string, p13 []string, p14 []byte) (float64, float32, int8, int16, int32, int64, uint8, uint16, uint32, uint64, bool, string, []string, []uint8){

	/* This function expects the parameters (in that order):
		double = 3.141592
	    float = 2.71f

	    int8 = -10
	    int16 = -20
	    int32 = -30
	    int64 = -40

	    uint8 = 50
	    uint16 = 60
	    uint32 = 70
	    uint64 = 80

	    bool = 1

	    string = "This is an input"
	    string[] = {"element one", "element two"}

	    bytes = {2, 4, 6, 8, 10}
	*/

	println("Hello from Go F1")

	if p1 != 3.141592{
		panic("p1 != 3.141592")
	}

	if p2 != 2.71{
		panic("p2 != 2.71")
	}

	if p3 != -10{
		panic("p3 != -10")
	}

	if p4 != -20{
		panic("p4 != -20")
	}

	if p5 != -30{
		panic("p5 != -30")
	}

	if p6 != -40{
		panic("p6 != -40")
	}

	if p7 != 50{
		panic("p7 != 50")
	}

	if p8 != 60{
		panic("p8 != 60")
	}

	if p9 != 70{
		panic("p9 != 70")
	}

	if p10 != 80{
		panic("p10 != 80")
	}

	if !p11 {
		panic("p11 == false")
	}

	if p12 != "This is an input"{
		panic("p12 != \"This is an input\"")
	}

	if len(p13) != 2{
		panic(fmt.Sprintf("len(p13) != 2. len(p13)=%v", len(p13)))
	}

	if p13[0] != "element one"{
		panic("p13[0] != \"element one\"")
	}

	if p13[1] != "element two"{
		panic("p13[1] != \"element two\"")
	}

	if len(p14) != 5{
		panic("len(p14) != 5")
	}

	if p14[0] != 2 || p14[1] != 4 || p14[2] != 6 || p14[3] != 8 || p14[4] != 10{
		panic("p14[0] != 2 || p14[1] != 4 || p14[2] != 6 || p14[3] != 8 || p14[4] != 10")
	}

	var r1 float64 = 0.57721
	var r2 float32 = 3.359
	
	var r3 int8 = -11
	var r4 int16 = -21
	var r5 int32 = -31
	var r6 int64 = -41
	
	var r7 uint8 = 51
	var r8 uint16 = 61
	var r9 uint32 = 71
	var r10 uint64 = 81

	var r11 bool = true

	var r12 string = "This is an output"
	var r13 []string = []string{ "return one", "return two" }

	var r14 []uint8 = []uint8{ 20, 40, 60, 80, 100 }

	return r1, r2, r3, r4, r5, r6, r7, r8, r9, r10, r11, r12, r13, r14
}`

//--------------------------------------------------------------------
func TestGuest(t *testing.T) {
	
	skipMessage := "Cannot execute guest test from Go, as Go->Go is currently unsupported due to Go inability to load Go shared modules.\n"
	skipMessage += "To run the test, execute the C code in CompilerGuest_testHelper.go from C/C++."
	t.Skip(skipMessage)
	
	def, err := compiler.NewIDLDefinitionFromJSON(idl_guest)
	if err != nil {
		t.Fatal(err)
		return
	}
	
	_ = os.RemoveAll("temp")
	
	err = os.Mkdir("temp", 0700)
	if err != nil {
		t.Fatal(err)
		return
	}
	
	err = ioutil.WriteFile("./temp/GuestCode.go", []byte(GuestCode), 0600)
	if err != nil {
		t.Fatal(err)
		return
	}
	
	err = ioutil.WriteFile("./temp/go.mod", []byte("module GoFuncs"), 0600)
	
	defer func() {
		err = os.RemoveAll("temp")
		if err != nil {
			t.Fatal(err)
			return
		}
	}()
	
	cmp := NewGuestCompiler()
	err = cmp.Compile(def, "temp", "", "", "")
	if err != nil {
		t.Fatal(err)
		return
	}
	
	if CallHostMock() != 0 {
		t.Fatal("Failed calling guest")
	}
}

//--------------------------------------------------------------------
