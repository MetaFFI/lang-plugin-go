package main

import (
	"testing"
)

var src string = `
package TestFuncs

import (
	"strings"
	"time"
)

type MyString string

const (
	NumberConst = 100000
	FiveSeconds = time.Second*5
)

var(
	NumberVar = 3653
	StringVar MyString = "Test"
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

func WaitABit(d time.Duration) error{
	time.Sleep(d)
	return nil
}

type TestMap struct{
	m map[string]interface{}
	Name string
}

type KV interface{
	Set(string, interface{})
	Get(string) interface{}
}

func NewTestMap() *TestMap{
	return &TestMap{ 
		m: make(map[string]interface{}),
		Name: "TestMap Name",
	}
}

func (this *TestMap) Set(k string, v interface{}){
	this.m[k] = v
}

func (this *TestMap) Get(k string) interface{}{
	return this.m[k]
}

func (this *TestMap) Contains(k string) bool{
	_, found := this.m[k]
	return found
}

`

func TestGoIDLCompiler_Compile(t *testing.T) {
	comp := NewGoIDLCompiler()

	_, _, err := comp.ParseIDL(src, "")
	if err != nil {
		t.Fatal(err)
	}
}

func TestGoIDLCompiler_CompilePackage(t *testing.T) {
	comp := NewGoIDLCompiler()

	idl, _, err := comp.ParseIDL("", "text/template")
	if err != nil {
		t.Fatal(err)
	}

	_, err = idl.ToJSON()
	if err != nil {
		t.Fatal(err)
	}
}
