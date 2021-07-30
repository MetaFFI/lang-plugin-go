package main

import (
	compiler "github.com/MetaFFI/plugin-sdk/compiler/go"
	"os"
	"os/exec"
	"testing"
)

const idl_host = `{"idl_filename": "test","idl_extension": ".proto","idl_filename_with_extension": "test.proto","idl_full_path": "","modules": [{"name": "Service1","target_language": "test","comment": "Comments for Service1\n","tags": {"metaffi_function_path": "package=main","metaffi_target_language": "test"},"functions": [{"name": "f1","comment": "F1 comment\nparam1 comment\n","tags": {"metaffi_function_path": "function=F1,metaffi_guest_lib=$PWD/temp/test_MetaFFIGuest.so"},"path_to_foreign_function": {"module": "$PWD/temp","package": "GoFuncs","function": "F1"},"parameter_type": "Params1","return_values_type": "Return1","parameters": [{"name": "p1","type": "float64","comment": "= 3.141592","tags": null,"dimensions": 0,"pass_method": ""},{"name": "p2","type": "float32","comment": "= 2.71","tags": null,"dimensions": 0,"pass_method": ""},{"name": "p3","type": "int8","comment": "= -10","tags": null,"dimensions": 0,"pass_method": ""},{"name": "p4","type": "int16","comment": "= -20","tags": null,"dimensions": 0,"pass_method": ""},{"name": "p5","type": "int32","comment": "= -30","tags": null,"dimensions": 0,"pass_method": ""},{"name": "p6","type": "int64","comment": "= -40","tags": null,"dimensions": 0,"pass_method": ""},{"name": "p7","type": "uint8","comment": "= 50","tags": null,"dimensions": 0,"pass_method": ""},{"name": "p8","type": "uint16","comment": "= 60","tags": null,"dimensions": 0,"pass_method": ""},{"name": "p9","type": "uint32","comment": "= 70","tags": null,"dimensions": 0,"pass_method": ""},{"name": "p10","type": "uint64","comment": "= 80","tags": null,"dimensions": 0,"pass_method": ""},{"name": "p11","type": "bool","comment": "= true","tags": null,"dimensions": 0,"pass_method": ""},{"name": "p12","type": "string8","comment": "= This is an input","tags": null,"dimensions": 0,"pass_method": ""},{"name": "p13","type": "string8","comment": "= {element one, element two}","tags": null,"dimensions": 1,"pass_method": ""},{"name": "p14","type": "uint8","comment": "= {2, 4, 6, 8, 10}","tags": null,"dimensions": 1,"pass_method": ""}],"return_values": [{"name": "r1","type": "float64","comment": "= 0.57721","tags": null,"dimensions": 0,"pass_method": ""},{"name": "r2","type": "float32","comment": "= 3.359","tags": null,"dimensions": 0,"pass_method": ""},{"name": "r3","type": "int8","comment": "= -11","tags": null,"dimensions": 0,"pass_method": ""},{"name": "r4","type": "int16","comment": "= -21","tags": null,"dimensions": 0,"pass_method": ""},{"name": "r5","type": "int32","comment": "= -31","tags": null,"dimensions": 0,"pass_method": ""},{"name": "r6","type": "int64","comment": "= -41","tags": null,"dimensions": 0,"pass_method": ""},{"name": "r7","type": "uint8","comment": "= 51","tags": null,"dimensions": 0,"pass_method": ""},{"name": "r8","type": "uint16","comment": "= 61","tags": null,"dimensions": 0,"pass_method": ""},{"name": "r9","type": "uint32","comment": "= 71","tags": null,"dimensions": 0,"pass_method": ""},{"name": "r10","type": "uint64","comment": "= 81","tags": null,"dimensions": 0,"pass_method": ""},{"name": "r11","type": "bool","comment": "= true","tags": null,"dimensions": 0,"pass_method": ""},{"name": "r12","type": "string8","comment": "= This is an output","tags": null,"dimensions": 0,"pass_method": ""},{"name": "r13","type": "string8","comment": "= {return one, return two}","tags": null,"dimensions": 1,"pass_method": ""},{"name": "r14","type": "uint8","comment": "= {20, 40, 60, 80, 100}","tags": null,"dimensions": 1,"pass_method": ""}]}]}]}`
const compilerTestCode = `package main

import (
	"testing"
)

func TestHostCompiler(t *testing.T){

	r1,r2,r3,r4,r5,r6,r7,r8,r9,r10,r11,r12,r13,r14, err := F1(3.141592,
															2.71,
															-10,
															-20,
															-30,
															-40,
															50,
															60,
															70,
															80,
															true,
															"This is an input",
															[]string{"element one", "element two"},
															[]byte{2, 4, 6, 8, 10})

	if err != nil{
		t.Fatalf("Failed with error: %v\n", err)
	}

	if r1 != 0.57721{ t.Fatalf("r1 != 0.57721. r1=%v", r1) }
	if r2 != 3.359{ t.Fatalf("r2 != 3.359") }
	
	if r3 != -11{ t.Fatalf("r3 != -11") }
	if r4 != -21{ t.Fatalf("r4 != -21") }
	if r5 != -31{ t.Fatalf("r5 != -31") }
	if r6 != -41{ t.Fatalf("r6 != -41") }
	
	if r7 != 51{ t.Fatalf("r7 != 51") }
	if r8 != 61{ t.Fatalf("r8 != 61") }
	if r9 != 71{ t.Fatalf("r9 != 71") }
	if r10 != 81{ t.Fatalf("r10 != 81") }

	if !r11 { t.Fatalf("r11 == false") }

	if r12 != "This is an output" { t.Fatalf("r12 != \"This is an output\"") }
	
	if(len(r13) != 2){ t.Fatalf("len(r13) != 2") }
	if r13[0] != "return one" { t.Fatalf("r13[0] != \"return one\"") }
	if r13[1] != "return two" { t.Fatalf("r13[1] != \"return two\"") }
	
	if(len(r14) != 5){ t.Fatalf("len(r14) != 5") }
	if r14[0] != 20 { t.Fatalf("r14[0] != 20") }
	if r14[1] != 40 { t.Fatalf("r14[1] != 40") }
	if r14[2] != 60 { t.Fatalf("r14[2] != 60") }
	if r14[3] != 80 { t.Fatalf("r14[3] != 80") }
	if r14[4] != 100 { t.Fatalf("r14[4] != 100") }
}`

//--------------------------------------------------------------------
func TestHost(t *testing.T){

	def, err := compiler.NewIDLDefinition(idl_host)
	if err != nil{
		t.Fatal(err)
		return
	}

	_ = os.RemoveAll("temp")

	err = os.Mkdir("temp", 0700)
	if err != nil{
		t.Fatal(err)
		return
	}

	defer func(){
		err = os.RemoveAll("temp")
		if err != nil{
			t.Fatal(err)
			return
		}
	}()

	cmp := NewCompiler(def, "./temp")
	_, err = cmp.CompileHost(nil)
	if err != nil{
		t.Fatal(err)
		return
	}

	err = os.WriteFile("./temp/CompilerTestCode_test.go", []byte(compilerTestCode), 0700)
	if err != nil {
		t.Fatal(err)
		return
	}

	buildCmd := exec.Command("go", "test", "-v")
	buildCmd.Dir = "./temp"
	output, err := buildCmd.CombinedOutput()
	if err != nil{
		println(string(output))
		t.Fatalf("Failed building Go Host test code with error: %v.\nOutput:\n%v", err, string(output))
	}

}
//--------------------------------------------------------------------