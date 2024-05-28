package TestRuntime

import (
	"strings"
	"time"
	"fmt"
)

func HelloWorld() {
	println("Hello World, From Go!")
}

func ReturnsAnError(){
	panic("An error from ReturnsAnError")
}

func DivIntegers(x int, y int) float32{

	if y == 0{
		panic("Divisor is 0")
	}

	return float32(x) / float32(y)
}

func JoinStrings(arrs []string) string{
	return strings.Join(arrs, ",")
}

const FiveSeconds = time.Second*5
func WaitABit(d time.Duration) error{
	time.Sleep(d)
	return nil
}

//--------------------------------------------

type SomeClass struct{}

func (s SomeClass) Print() {
	fmt.Println("Hello from inner class")
}

func GetSomeClasses() []SomeClass {
	return []SomeClass{{}, {}, {}}
}

func ExpectThreeSomeClasses(arr []SomeClass) {
	if len(arr) != 3 {
		panic("Array length is not 3")
	}
}

func ExpectThreeBuffers(buffers [][]byte) {
	if len(buffers) != 3 {
		panic("Buffers length is not 3")
	}
}

func GetThreeBuffers() [][]byte {
	buffers := make([][]byte, 3)
	for i := 0; i < 3; i++ {
		buffers[i] = []byte{1, 2, 3}
	}
	return buffers
}


//--------------------------------------------

type TestMap struct{
	m map[string]interface{}
	Name string
}

func NewTestMap() *TestMap{
	return &TestMap{ 
		m: make(map[string]interface{}),
		Name: "TestMap Name",
	}
}

func (this *TestMap) Set(k string, v interface{}){
	fmt.Printf("Setting: %v %T\n", k, v)
	this.m[k] = v
}

func (this *TestMap) Get(k string) interface{}{
	v := this.m[k]
	fmt.Printf("Getting: %v %T\n", k, v)
	return v
}

func (this *TestMap) Contains(k string) bool{
	_, found := this.m[k]
	fmt.Printf("Contains: %v %v\n", k, found)
	return found
}
