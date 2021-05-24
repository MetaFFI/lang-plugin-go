package main

import (
	"testing"
)

func TestHostCompiler(t *testing.T){

	res, err := F1(3.141592,
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

	if len(res) != 2{
		t.Fatalf("Expected result length is 2. Received: %v\n", len(res))
	}

	if res[0] != "return one"{
		t.Fatalf("Expected res[0] is \"return one\". Received: %v\n", res[0])
	}

	if res[1] != "return two"{
		t.Fatalf("Expected res[1] is \"return two\". Received: %v\n", res[1])
	}
}