package main

import (
	idl "github.com/MetaFFI/lang-plugin-go/idl/IDLCompiler"
	"os"
	"testing"
)

func TestLibrary(t *testing.T) {
	idlCompiler := idl.NewGoIDLCompiler()
	idlDef, _, err := idlCompiler.ParseIDL("", "text/template")
	if err != nil {
		t.Fatal(err)
		return
	}

	_ = os.RemoveAll("temp_guest")
	err = os.Mkdir("temp_guest", 0700)
	if err != nil {
		t.Fatal(err)
		return
	}

	defer func() {
		err = os.RemoveAll("temp_guest")
		if err != nil {
			t.Fatal(err)
			return
		}
	}()

	cmp := NewGuestCompiler()
	err = cmp.Compile(idlDef, "temp_guest", "", nil)
	if err != nil {
		t.Fatal(err)
		return
	}
}
