package IDLCompiler

import (
	"fmt"
	"github.com/GreenFuze/go-parser"
	"github.com/MetaFFI/plugin-sdk/compiler/go/IDL"
	"strings"
)

var classes map[string]*IDL.ClassDefinition

// --------------------------------------------------------------------
func init() {
	classes = make(map[string]*IDL.ClassDefinition)
}

// --------------------------------------------------------------------

func LoadClasses(gofile *parser.GoFile, metaffiGuestLib string) {

	for _, i := range gofile.Interfaces {
		if !IsPublic(i.Name) {
			continue
		}

		interfaceDef := IDL.NewClassDefinition(i.Name)
		interfaceDef.Comment = i.Comments
		classes[interfaceDef.Name] = interfaceDef
	}

	for _, s := range gofile.Structs {

		if !IsPublic(s.Name) {
			continue
		}

		structDef := IDL.NewClassDefinition(s.Name)
		structDef.Comment = s.Comments

		for _, f := range s.Fields {
			if !IsPublic(f.Name) {
				continue
			}

			fdecl := IDL.NewFieldDefinition(structDef, f.Name, goTypeToMFFI(f.Type), "Get"+f.Name, "Set"+f.Name, true)
			fdecl.TypeAlias = f.Type
			fdecl.Getter.SetTag("receiver_pointer", "true")
			fdecl.Getter.SetFunctionPath("metaffi_guest_lib", metaffiGuestLib)
			fdecl.Getter.SetFunctionPath("entrypoint_function", "EntryPoint_"+structDef.Name+"_"+fdecl.Getter.Name)

			fdecl.Setter.SetTag("receiver_pointer", "true")
			fdecl.Setter.SetFunctionPath("metaffi_guest_lib", metaffiGuestLib)
			fdecl.Setter.SetFunctionPath("entrypoint_function", "EntryPoint_"+structDef.Name+"_"+fdecl.Setter.Name)

			structDef.AddField(fdecl)
		}

		classes[structDef.Name] = structDef
	}
}

// --------------------------------------------------------------------

func LoadMethods(gofile *parser.GoFile, metaffiGuestLib string) {
	for _, meth := range gofile.StructMethods {
		if !IsPublic(meth.Name) {
			continue
		}

		if meth.Receivers == nil || len(meth.Receivers) == 0 {
			continue // function, not a method
		}

		cls := getReceiverClass(meth.Receivers[0])

		funcDecl := IDL.NewFunctionDefinition(meth.Name)
		funcDecl.Comment = meth.Comments

		for i, p := range meth.Params {

			if strings.HasPrefix(p.Type, "...") { // if ellipsis
				p.Type = strings.ReplaceAll(p.Type, "...", "[]")
				p.Underlying = strings.ReplaceAll(p.Underlying, "...", "[]")
				funcDecl.SetTag("variadic_parameter", p.Name)
			}

			var alias string
			if p.Underlying != p.Type || !isPrimitiveType(p.Type) {
				alias = p.Type
			}

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

		for i, p := range meth.Results {
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

		methDecl := IDL.NewMethodDefinitionWithFunction(cls, funcDecl, true)
		if strings.Contains(meth.Receivers[0], "*") {
			methDecl.SetTag("receiver_pointer", "true")
		}

		methDecl.SetFunctionPath("metaffi_guest_lib", metaffiGuestLib)
		methDecl.SetFunctionPath("entrypoint_function", "EntryPoint_"+cls.Name+"_"+methDecl.Name)

		cls.AddMethod(methDecl)
	}

	res := make([]*IDL.ClassDefinition, 0)
	for _, c := range classes {
		c.Releaser.SetFunctionPath("metaffi_guest_lib", metaffiGuestLib)
		c.Releaser.SetFunctionPath("entrypoint_function", "EntryPoint_"+c.Name+"_"+c.Releaser.Name)
		res = append(res, c)
	}

}

// --------------------------------------------------------------------

func ExtractClasses() []*IDL.ClassDefinition {
	res := make([]*IDL.ClassDefinition, 0)

	for _, c := range classes {
		res = append(res, c)
	}

	return res
}

// --------------------------------------------------------------------

func getReceiverClass(s string) *IDL.ClassDefinition {
	s = strings.ReplaceAll(s, "[", "")
	s = strings.ReplaceAll(s, "]", "")
	s = strings.ReplaceAll(s, "*", "")

	cls, found := classes[s]
	if !found {
		if strings.Contains(s, ".") { // remove "package name" if there is one
			s = s[strings.Index(s, ".")+1:]
		}

		cls, found = classes[s]

		if !found {
			panic(fmt.Errorf("Cannot find receiver class %v", s))
		}
	}

	return cls
}

//--------------------------------------------------------------------
