package IDLCompiler

import (
	"fmt"
	"github.com/GreenFuze/go-parser"
	"github.com/MetaFFI/plugin-sdk/compiler/go/IDL"
	"strings"
)

// --------------------------------------------------------------------
func ExtractFunctions(gofile *parser.GoFile, metaffiGuestLib string) []*IDL.FunctionDefinition {

	functions := make([]*IDL.FunctionDefinition, 0)

	for _, f := range gofile.StructMethods {

		if !IsPublic(f.Name) {
			continue
		}

		if f.Receivers != nil && len(f.Receivers) > 0 {
			continue // method, not a function
		}

		funcDecl := IDL.NewFunctionDefinition(f.Name)
		funcDecl.Comment = f.Comments

		for i, p := range f.Params {
			var alias string
			alias = p.Underlying

			var name string
			if p.Name != "" {
				name = p.Name
			} else {
				name = fmt.Sprintf("p%v", i)
			}
			funcDecl.AddParameter(IDL.NewArgArrayDefinitionWithAlias(name, goTypeToMFFI(p.Underlying), strings.Count(p.Underlying, "[]"), alias))

			if alias != "" {
				imp := GetRequiredImport(gofile, alias)
				if imp != "" {
					Imports[imp] = true
				}
			}
		}

		for i, p := range f.Results {
			var alias string
			if p.Underlying != p.Type {
				alias = p.Type
			}

			var name string
			if p.Name != "" {
				name = p.Name
			} else {
				name = fmt.Sprintf("r%v", i)
			}

			funcDecl.AddReturnValues(IDL.NewArgArrayDefinitionWithAlias(name, goTypeToMFFI(p.Underlying), strings.Count(p.Underlying, "[]"), alias))

			if alias != "" {
				imp := GetRequiredImport(gofile, alias)
				if imp != "" {
					Imports[imp] = true
				}
			}
		}

		funcDecl.SetFunctionPath("metaffi_guest_lib", metaffiGuestLib)
		funcDecl.SetFunctionPath("entrypoint_function", "EntryPoint_"+funcDecl.Name)

		// check if constructor
		if cls := isConstructorFunction(f); cls != nil {
			funcDecl.SetFunctionPath("entrypoint_function", "EntryPoint_"+cls.Name+"_"+funcDecl.Name)
			cls.AddConstructor(IDL.NewConstructorDefinitionFromFunctionDefinition(funcDecl))
		} else {
			functions = append(functions, funcDecl)
		}
	}

	return functions
}

// --------------------------------------------------------------------
func isConstructorFunction(f *parser.GoStructMethod) *IDL.ClassDefinition {

	// if there is only 1 struct (defined in file) returning from function - it is a constructor

	var foundClass *IDL.ClassDefinition

	for _, r := range f.Results {
		rType := r.Type
		rType = strings.ReplaceAll(rType, "[", "")
		rType = strings.ReplaceAll(rType, "]", "")
		rType = strings.ReplaceAll(rType, "*", "")

		if cls, found := classes[rType]; found {
			if foundClass != nil {
				return nil // 2 classes returning - not a constructor
			}

			foundClass = cls
		}
	}

	return foundClass
}

//--------------------------------------------------------------------
