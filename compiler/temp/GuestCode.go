
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
}