package IDLCompiler

import (
	"fmt"
	"github.com/GreenFuze/go-parser"
	"github.com/MetaFFI/plugin-sdk/compiler/go/IDL"
	"go/build"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"unicode"
)

var Imports map[string]bool

// --------------------------------------------------------------------
func IsPublic(name string) bool {
	if name == "" {
		return false
	}

	return unicode.IsUpper(rune(name[0]))
}

// --------------------------------------------------------------------

type GoIDLCompiler struct {
}

// --------------------------------------------------------------------

func NewGoIDLCompiler() *GoIDLCompiler {

	Imports = make(map[string]bool)

	// get Go file AST
	return &GoIDLCompiler{}
}

// --------------------------------------------------------------------

func (this *GoIDLCompiler) parseSource(goSourceCode string, gofilepath string, metaFFIGuestLib string, mod *IDL.ModuleDefinition) (bool, error) {

	gofile, err := parser.ParseSource(goSourceCode, gofilepath, true)
	if err != nil {
		return true, err
	}

	globals := ExtractGlobals(gofile, metaFFIGuestLib)
	LoadClasses(gofile, metaFFIGuestLib)
	LoadMethods(gofile, metaFFIGuestLib)
	classes := ExtractClasses()
	functions := ExtractFunctions(gofile, metaFFIGuestLib)

	for imp, _ := range Imports {
		mod.AddExternalResourceIfNotExist(imp)
	}

	importPath, _, err := gofile.ImportPath()
	if err != nil {
		return true, err
	}

	importPath = strings.Replace(importPath, "\\", "/", -1)

	mod.AddGlobals(globals)
	mod.AddFunctions(functions)

	for _, c := range classes {
		mod.AddClass(c)
	}

	mod.SetFunctionPath("package", gofile.Package)
	mod.SetFunctionPath("module", importPath)

	return true, nil
}

// --------------------------------------------------------------------

func (this *GoIDLCompiler) parseDir(dir string, metaFFIGuestLib string, mod *IDL.ModuleDefinition, isInGoROOT bool) (bool, error) {

	gofiles, err := parser.ParseDir(dir, true, func(file fs.FileInfo) bool {

		if file.IsDir() {
			return false
		}

		if strings.ToLower(filepath.Ext(file.Name())) != ".go" { // skip non-go files
			return false
		}

		if strings.HasSuffix(strings.ToLower(filepath.Base(file.Name())), "_test.go") { // skip test files
			return false
		}

		return true
	})
	if err != nil {
		return true, err
	}

	// Methods can be implemented in different files than their types
	// therefore, first load all classes, then the methods
	for _, gofile := range gofiles {
		LoadClasses(gofile, metaFFIGuestLib)
	}

	for _, gofile := range gofiles {
		LoadMethods(gofile, metaFFIGuestLib)
	}

	for _, c := range classes {
		mod.AddClass(c)
	}

	for _, gofile := range gofiles {
		globals := ExtractGlobals(gofile, metaFFIGuestLib)
		functions := ExtractFunctions(gofile, metaFFIGuestLib)

		for imp, _ := range Imports {
			mod.AddExternalResourceIfNotExist(imp)
		}

		importPath, _, err := gofile.ImportPath()
		if err != nil {
			return true, err
		}

		importPath = strings.Replace(importPath, "\\", "/", -1)

		mod.AddGlobals(globals)
		mod.AddFunctions(functions)

		if !isInGoROOT {
			mod.SetFunctionPath("package", gofile.Package)
			mod.SetFunctionPath("module", importPath)
		} else {
			mod.SetFunctionPath("package", removeGoPathModulePathFromPath(removeGoRootSrcPathFromPath(dir)))
			mod.SetFunctionPath("module", removeGoPathModulePathFromPath(removeGoRootSrcPathFromPath(dir)))
		}

	}

	return true, nil
}

//--------------------------------------------------------------------

func (this *GoIDLCompiler) parseFile(gofilepath string, metaFFIGuestLib string, mod *IDL.ModuleDefinition) (bool, error) {

	data, err := os.ReadFile(gofilepath)
	if err != nil {
		return true, err
	}

	return this.parseSource(string(data), gofilepath, metaFFIGuestLib, mod)
}

//--------------------------------------------------------------------

func getGOPATHModulePath() string {
	return build.Default.GOPATH + "/pkg/mod/"
}

func getGOROOTSrcPath() string {
	return build.Default.GOROOT + "/src/"
}

func removeGoPathModulePathFromPath(path string) string {
	return strings.ReplaceAll(path, getGOPATHModulePath(), "")
}
func removeGoRootSrcPathFromPath(path string) string {
	return strings.ReplaceAll(path, getGOROOTSrcPath(), "")
}

// --------------------------------------------------------------------
func replaceUpper(s string) string {
	result := ""
	for _, c := range s {
		if c >= 'A' && c <= 'Z' {
			result += "!" + strings.ToLower(string(c))
		} else {
			result += string(c)
		}
	}
	return result
}

func findCodePath(path string) (fullPath string, foundPath os.FileInfo, isInGoROOT bool, isInGoPATH bool, err error) {
	foundPath, errDirect := os.Stat(path)
	if errDirect == nil {
		return path, foundPath, false, false, nil
	}

	foundPath, errGoRoot := os.Stat(getGOROOTSrcPath() + path)
	if errGoRoot == nil {
		return getGOROOTSrcPath() + path, foundPath, true, false, nil
	}

	foundPath, errGoPath := os.Stat(getGOPATHModulePath() + replaceUpper(path))
	if errGoPath == nil {
		return getGOPATHModulePath() + replaceUpper(path), foundPath, false, true, nil
	}

	return "", nil, false, false, fmt.Errorf("Couldn't find given path. Errors:\n%v\n%v\n%v", errDirect, errGoRoot, errGoPath)
}

func (this *GoIDLCompiler) ParseIDL(goSourceCode string, gofilepath string) (*IDL.IDLDefinition, bool, error) {

	idl := IDL.NewIDLDefinition(gofilepath, "go")

	// parse AST and build IDLDefinition

	if goSourceCode != "" { // source code is available, generate IDL from source code
		module := IDL.NewModuleDefinition("go")
		_, err := this.parseSource(goSourceCode, gofilepath, idl.MetaFFIGuestLib, module)
		if err != nil {
			return nil, true, err
		}
		idl.AddModule(module)

	} else { // no given source code - only a path

		if gofilepath == "" {
			return nil, true, fmt.Errorf("No given source code or path")
		}

		// if the gofilepath is a file, read source code and generate IDL, if path, generate IDL from all source code files
		gofilepath, fi, isInGoROOT, _, err := findCodePath(gofilepath)
		if err != nil {
			return nil, true, err
		}

		if !fi.IsDir() {

			// a file
			module := IDL.NewModuleDefinition(strings.ReplaceAll(filepath.Base(gofilepath), filepath.Ext(gofilepath), ""))
			_, err = this.parseFile(gofilepath, idl.MetaFFIGuestLib, module)
			if err != nil {
				return nil, true, err
			}
			idl.AddModule(module)

		} else {

			// directory
			var module *IDL.ModuleDefinition
			if isInGoROOT {
				module = IDL.NewModuleDefinition(removeGoRootSrcPathFromPath(gofilepath))
			} else {
				module = IDL.NewModuleDefinition(filepath.Base(gofilepath))
			}

			_, err = this.parseDir(gofilepath, idl.MetaFFIGuestLib, module, isInGoROOT)
			if err != nil {
				return nil, true, err
			}

			idl.AddModule(module)
		}
	}

	idl.FinalizeConstruction()

	return idl, true, nil
}

//--------------------------------------------------------------------
