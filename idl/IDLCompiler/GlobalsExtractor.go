package IDLCompiler

import (
	"fmt"
	"github.com/GreenFuze/go-parser"
	"github.com/MetaFFI/plugin-sdk/compiler/go/IDL"
	"strings"
)

// --------------------------------------------------------------------
func ExtractGlobals(gofile *parser.GoFile, metaffiGuestLib string) []*IDL.GlobalDefinition {

	globalsDefs := make([]*IDL.GlobalDefinition, 0)

	for _, gs := range gofile.GlobalConstants {
		if !IsPublic(gs.Name) {
			continue
		}

		var alias string
		if gs.Underlying != gs.Type {
			alias = gs.Type
		}
		global := IDL.NewGlobalDefinitionWithAlias(gs.Name, goTypeToMFFI(gs.Underlying), alias, "Get"+gs.Name, "")
		global.Getter.SetFunctionPath("metaffi_guest_lib", metaffiGuestLib)
		global.Getter.SetFunctionPath("entrypoint_function", "EntryPoint_"+global.Getter.Name)

		globalsDefs = append(globalsDefs, global)
	}

	for _, gv := range gofile.GlobalVariables {
		if !IsPublic(gv.Name) {
			continue
		}

		var alias string
		if gv.Underlying != gv.Type {
			alias = gv.Type
		}
		global := IDL.NewGlobalDefinitionWithAlias(gv.Name, goTypeToMFFI(gv.Underlying), alias, "Get"+gv.Name, "Set"+gv.Name)
		globalsDefs = append(globalsDefs, IDL.NewGlobalDefinitionWithAlias(gv.Name, goTypeToMFFI(gv.Underlying), alias, "Get"+gv.Name, "Set"+gv.Name))
		global.Getter.SetFunctionPath("metaffi_guest_lib", metaffiGuestLib)
		global.Getter.SetFunctionPath("entrypoint_function", "EntryPoint_"+global.Getter.Name)
		global.Setter.SetFunctionPath("metaffi_guest_lib", metaffiGuestLib)
		global.Setter.SetFunctionPath("entrypoint_function", "EntryPoint_"+global.Setter.Name)
	}

	return globalsDefs
}

// --------------------------------------------------------------------
func GetRequiredImport(gofile *parser.GoFile, fullType string) string {

	if strings.Contains(fullType, "...") {
		fullType = strings.ReplaceAll(fullType, "...", "[]") // replace ellipsis with array
	}

	if !strings.Contains(fullType, ".") {
		return ""
	}

	// get package name
	splitType := strings.Split(fullType, ".")

	packageName := splitType[len(splitType)-2] // get one before the last element
	packageName = strings.ReplaceAll(strings.ReplaceAll(packageName, "*", ""), "[]", "")

	if packageName == gofile.Package { // no import required
		return ""
	}

	for _, imp := range gofile.Imports {
		var impname string
		if imp.Name != "" {
			impname = strings.ReplaceAll(imp.Name, `"`, "")
		} else {
			impname = strings.ReplaceAll(imp.Path, `"`, "")
			impname = impname[strings.LastIndex(impname, ".")+1:]
		}

		if impname == packageName {
			return strings.ReplaceAll(imp.Path, `"`, "")
		}
	}

	panic(fmt.Errorf("package name \"%v\" (of type: %v) is used, but cannot find its import", packageName, fullType))
}

//--------------------------------------------------------------------
