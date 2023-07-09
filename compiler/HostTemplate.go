package main

import "C"
import "runtime"

const HostHeaderTemplate = `
// Code generated by MetaFFI. DO NOT EDIT.
// Host code for {{.IDLFilenameWithExtension}}
`

const HostPackageTemplate = `package {{.Package}}
`

const HostImportsTemplate = `
import "fmt"
import "unsafe"
{{if IsGoRuntimePackNeeded .}}import . "github.com/MetaFFI/lang-plugin-go/go-runtime"{{end}}
`

func GetHostHelperFunctions() string {
	if runtime.GOOS == "windows" {
		return HostHelperFunctionsWindows
	} else {
		return HostHelperFunctionsNonWindows
	}
}

func GetHostHelperFunctionsName() string {
	if runtime.GOOS == "windows" {
		return "HostHelperFunctionsWindows"
	} else {
		return "HostHelperFunctionsNonWindows"
	}
}

const HostHelperFunctionsWindows = `
{{ $idl := . }}

// function IDs
{{range $mindex, $m := .Modules}}

{{range $findex, $f := $m.Globals}}
{{if $f.Getter}}var {{$f.Getter.GetNameWithOverloadIndex}}_id unsafe.Pointer{{end}}
{{if $f.Setter}}var {{$f.Setter.GetNameWithOverloadIndex}}_id unsafe.Pointer{{end}}
{{end}}{{/* End globals */}}

{{range $findex, $f := $m.Functions}}
var {{$f.GetNameWithOverloadIndex}}_id unsafe.Pointer
{{end}}{{/* End Functions */}}

{{range $cindex, $c := $m.Classes}}
{{range $findex, $f := $c.Fields}}
{{if $f.Getter}}var {{$c.Name}}_{{$f.Getter.GetNameWithOverloadIndex}}_id unsafe.Pointer{{end}}
{{if $f.Setter}}var {{$c.Name}}_{{$f.Setter.GetNameWithOverloadIndex}}_id unsafe.Pointer{{end}}
{{end}}{{/* End Fields */}}
{{range $findex, $f := $c.Methods}}
var {{$c.Name}}_{{$f.GetNameWithOverloadIndex}}_id unsafe.Pointer
{{end}}{{/* End Methods */}}
{{range $findex, $f := $c.Constructors}}
var {{$c.Name}}_{{$f.GetNameWithOverloadIndex}}_id unsafe.Pointer
{{end}}{{/* End Constructor */}}
{{if $c.Releaser}}
var {{$c.Name}}_{{$c.Releaser.GetNameWithOverloadIndex}}_id unsafe.Pointer
{{end}}{{/* End Releaser */}}
{{end}}{{/* End Classes */}}
{{end}}{{/* End modules */}}

func MetaFFILoad(modulePath string){
	LoadCDTCAPI()

	runtimePlugin := "xllr.{{.TargetLanguage}}"
	err := XLLRLoadRuntimePlugin(runtimePlugin)
	if err != nil{
		panic(err)
	}

	// load functions
	loadFF := func(modulePath string, fpath string, paramsCount int8, retvalCount int8) unsafe.Pointer{
		id, err := XLLRLoadFunction(runtimePlugin, modulePath, fpath, nil, paramsCount, retvalCount)
		if err != nil{ // failed
			panic(err)
		}

		return id
	}

	{{range $mindex, $m := .Modules}}
	{{range $findex, $f := $m.Globals}}
	{{if $f.Getter}}{{$f.Getter.GetNameWithOverloadIndex}}_id = loadFF(modulePath, ` + "`" + `{{$f.Getter.FunctionPathAsString $idl}}` + "`" + `, {{len $f.Getter.Parameters}}, {{len $f.Getter.ReturnValues}}){{end}}
	{{if $f.Setter}}{{$f.Setter.GetNameWithOverloadIndex}}_id = loadFF(modulePath, ` + "`" + `{{$f.Setter.FunctionPathAsString $idl}}` + "`" + `, {{len $f.Setter.Parameters}}, {{len $f.Setter.ReturnValues}}){{end}}
	{{end}}{{/* End globals */}}

	{{range $findex, $f := $m.Functions}}
	{{$f.GetNameWithOverloadIndex}}_id = loadFF(modulePath, ` + "`" + `{{$f.FunctionPathAsString $idl}}` + "`" + `, {{len $f.Parameters}}, {{len $f.ReturnValues}})
	{{end}}{{/* End Functions */}}

	{{range $cindex, $c := $m.Classes}}
	{{range $findex, $f := $c.Fields}}
	{{if $f.Getter}}{{$c.Name}}_{{$f.Getter.GetNameWithOverloadIndex}}_id = loadFF(modulePath, ` + "`" + `{{$f.Getter.FunctionPathAsString $idl}}` + "`" + `, {{len $f.Getter.Parameters}}, {{len $f.Getter.ReturnValues}}){{end}}
	{{if $f.Setter}}{{$c.Name}}_{{$f.Setter.GetNameWithOverloadIndex}}_id = loadFF(modulePath, ` + "`" + `{{$f.Setter.FunctionPathAsString $idl}}` + "`" + `, {{len $f.Setter.Parameters}}, {{len $f.Setter.ReturnValues}}){{end}}
	{{end}}{{/* End Fields */}}
	{{range $findex, $f := $c.Methods}}
	{{$c.Name}}_{{$f.GetNameWithOverloadIndex}}_id = loadFF(modulePath, ` + "`" + `{{$f.FunctionPathAsString $idl}}` + "`" + `, {{len $f.Parameters}}, {{len $f.ReturnValues}})
	{{end}}{{/* End Methods */}}
	{{range $findex, $f := $c.Constructors}}
	{{$c.Name}}_{{$f.GetNameWithOverloadIndex}}_id = loadFF(modulePath, ` + "`" + `{{$f.FunctionPathAsString $idl}}` + "`" + `, {{len $f.Parameters}}, {{len $f.ReturnValues}})
	{{end}}{{/* End Constructor */}}
	{{if $c.Releaser}}
	{{$c.Name}}_{{$c.Releaser.GetNameWithOverloadIndex}}_id = loadFF(modulePath, ` + "`" + `{{$c.Releaser.FunctionPathAsString $idl}}` + "`" + `, {{len $c.Releaser.Parameters}}, {{len $c.Releaser.ReturnValues}})
	{{end}}{{/* End Releaser */}}
	{{end}}{{/* End Classes */}}
	{{end}}{{/* End modules */}}
}

func Free(){
	err := XLLRFreeRuntimePlugin("xllr.{{.TargetLanguage}}")
	if err != nil{ panic(err) }
}

`

const HostHelperFunctionsNonWindows = `
// function IDs
{{range $mindex, $m := .Modules}}

{{range $findex, $f := $m.Globals}}
{{if $f.Getter}}var {{$f.Getter.GetNameWithOverloadIndex}}_id unsafe.Pointer{{end}}
{{if $f.Setter}}var {{$f.Setter.GetNameWithOverloadIndex}}_id unsafe.Pointer{{end}}
{{end}}{{/* End globals */}}

{{range $findex, $f := $m.Functions}}
var {{$f.GetNameWithOverloadIndex}}_id unsafe.Pointer
{{end}}{{/* End Functions */}}

{{range $cindex, $c := $m.Classes}}
{{range $findex, $f := $c.Fields}}
{{if $f.Getter}}var {{$c.Name}}_{{$f.Getter.GetNameWithOverloadIndex}}_id unsafe.Pointer{{end}}
{{if $f.Setter}}var {{$c.Name}}_{{$f.Setter.GetNameWithOverloadIndex}}_id unsafe.Pointer{{end}}
{{end}}{{/* End Fields */}}
{{range $findex, $f := $c.Methods}}
var {{$c.Name}}_{{$f.GetNameWithOverloadIndex}}_id unsafe.Pointer
{{end}}{{/* End Methods */}}
{{range $findex, $f := $c.Constructors}}
var {{$c.Name}}_{{$f.GetNameWithOverloadIndex}}_id unsafe.Pointer
{{end}}{{/* End Constructor */}}
{{if $c.Releaser}}
var {{$c.Name}}_{{$c.Releaser.GetNameWithOverloadIndex}}_id unsafe.Pointer
{{end}}{{/* End Releaser */}}
{{end}}{{/* End Classes */}}
{{end}}{{/* End modules */}}

func Load(modulePath string){
	LoadCDTCAPI()

	runtimePlugin := "xllr.{{.TargetLanguage}}"
	err := XLLRLoadRuntimePlugin(runtimePlugin)
	if err != nil{
		panic(err)
	}

	// load functions
	loadFF := func(modulePath string, fpath string, paramsCount int8, retvalCount int8) unsafe.Pointer{
		id, err := XLLRLoadFunction(runtimePlugin, modulePath, fpath, nil, paramsCount, retvalCount)
		if err != nil{ // failed
			panic(err)
		}
		return id
	}

	{{ $idl := . }}
	{{range $mindex, $m := .Modules}}
	{{range $findex, $f := $m.Globals}}
	{{if $f.Getter}}{{$f.Getter.GetNameWithOverloadIndex}}_id = loadFF(modulePath, "{{$f.Getter.FunctionPathAsString $idl}}", {{len $f.Getter.Parameters}}, {{len $f.Getter.ReturnValues}} ){{end}}
	{{if $f.Setter}}{{$f.Setter.GetNameWithOverloadIndex}}_id = loadFF(modulePath, "{{$f.Setter.FunctionPathAsString $idl}}", {{len $f.Setter.Parameters}}, {{len $f.Setter.ReturnValues}} ){{end}}
	{{end}}{{/* End globals */}}
	
	{{range $findex, $f := $m.Functions}}
	{{$f.GetNameWithOverloadIndex}}_id = loadFF(modulePath, "{{$f.FunctionPathAsString $idl}}", {{len $f.Parameters}}, {{len $f.ReturnValues}})
	{{end}}{{/* End Functions */}}

	{{range $cindex, $c := $m.Classes}}
	{{range $findex, $f := $c.Fields}}
	{{if $f.Getter}}{{$c.Name}}_{{$f.Getter.GetNameWithOverloadIndex}}_id = loadFF(modulePath, "{{$f.Getter.FunctionPathAsString $idl}}", {{len $f.Getter.Parameters}}, {{len $f.Getter.ReturnValues}}){{end}}
	{{if $f.Setter}}{{$c.Name}}_{{$f.Setter.GetNameWithOverloadIndex}}_id = loadFF(modulePath, "{{$f.Setter.FunctionPathAsString $idl}}", {{len $f.Setter.Parameters}}, {{len $f.Setter.ReturnValues}}){{end}}
	{{end}}{{/* End Fields */}}
	{{range $findex, $f := $c.Methods}}
	{{$c.Name}}_{{$f.GetNameWithOverloadIndex}}_id = loadFF(modulePath, "{{$f.FunctionPathAsString $idl}}", {{len $f.Parameters}}, {{len $f.ReturnValues}})
	{{end}}{{/* End Methods */}}
	{{range $findex, $f := $c.Constructors}}
	{{$c.Name}}_{{$f.GetNameWithOverloadIndex}}_id = loadFF(modulePath, "{{$f.FunctionPathAsString $idl}}", {{len $f.Parameters}}, {{len $f.ReturnValues}})
	{{end}}{{/* End Constructor */}}
	{{if $c.Releaser}}
	{{$c.Name}}_{{$c.Releaser.GetNameWithOverloadIndex}}_id = loadFF(modulePath, "{{$c.Releaser.FunctionPathAsString $idl}}", {{len $c.Releaser.Parameters}}, {{len $c.Releaser.ReturnValues}})
	{{end}}{{/* End Releaser */}}
	{{end}}{{/* End Classes */}}
	{{end}}{{/* End modules */}}
}

func Free(){
	err := XLLRFreeRuntimePlugin("xllr.{{.TargetLanguage}}")
	if err != nil{ panic(err) }
}

`

const HostFunctionStubsTemplate = `
{{ $pfn := .IDLFilename}}
{{ $idl := . }}
{{range $mindex, $m := .Modules}}

{{range $findex, $f := $m.Globals}}
{{if $f.Getter}}
{{if $f.Comment}}/*
{{$f.Comment}}
*/{{end}}
func {{ToGoNameConv $f.Getter.GetNameWithOverloadIndex}}_MetaFFIGetter() (instance {{ConvertToGoType $f.ArgDefinition $m}}, err error){
	{{ $paramsLength := len $f.Getter.Parameters }}{{ $returnLength := len $f.Getter.ReturnValues }}
	{{GenerateCodeAllocateCDTS $f.Getter.Parameters $f.Getter.ReturnValues}}

	// parameters
	{{range $index, $elem := $f.Getter.Parameters}}
	FromGoToCDT({{$elem.Name}}, xcall_params, {{$index}})
	{{end}}{{/* End Parameters */}}

	{{GenerateCodeXCall "" $f.Getter.GetNameWithOverloadIndex $f.Getter.Parameters $f.Getter.ReturnValues}}
	
	{{range $index, $elem := $f.Getter.ReturnValues}}
	{{$elem.Name}}AsInterface := FromCDTToGo(return_valuesCDTS, {{$index}})
	var {{$elem.Name}} {{ConvertToGoType $elem $m}}
	
	if {{$elem.Name}}AsInterface != nil{
		{{if $elem.IsAny}}
		{{$elem.Name}} = {{$elem.Name}}AsInterface
		
		{{else if not $elem.IsHandle}}
		{{$elem.Name}} = {{if $elem.IsTypeAlias}}{{if $elem.IsArray}}[]{{end}}{{GetTypeOrAlias $elem $m}}{{else}}{{ConvertToGoType $elem $m}}{{end}}({{$elem.Name}}AsInterface.({{ConvertToGoType $elem $m}}))
		{{else}} {{/* Handle */}}		
		{{/* Go object */}}
		
		{{if not $elem.IsArray}}
		if obj, ok := {{$elem.Name}}AsInterface.(Handle); ok{ // None Go object			
			{{$elem.Name}} = {{HandleNoneGoObject $elem $m}}
		} else {
			{{$elem.Name}} = {{if and $elem.IsTypeAlias (not $elem.IsHandleTypeAlias)}}{{if $elem.IsArray}}[]{{end}}{{GetTypeOrAlias $elem $m}}{{else}}{{ConvertToGoType $elem $m}}{{end}}({{$elem.Name}}AsInterface.({{ConvertToGoType $elem $m}}))
		}
		{{else}}
		{{if $elem.IsTypeAlias}} {{/* a type is specified */}}
		if len({{$elem.Name}}AsInterface.([]interface{})) > 0{
			{{$elem.Name}} = make([]{{GetTypeOrAlias $elem $m}}, len({{$elem.Name}}AsInterface.([]interface{})))
			if _, ok := {{$elem.Name}}AsInterface.([]interface{})[0].(Handle); ok{
				for i, h := range {{$elem.Name}}AsInterface.([]interface{}){
					{{$elem.Name}}[i] = {{GetTypeOrAlias $elem $m}}{ h: h.(Handle) }
				}
			} else {
				for i, obj := range {{$elem.Name}}AsInterface.([]interface{}){
					{{$elem.Name}}[i] = obj.({{GetTypeOrAlias $elem $m}})
				}
			}
		}
		{{else}}
		{{$elem.Name}} = {{$elem.Name}}AsInterface
		{{end}}
		{{end}}

		{{end}}{{/* end handling types */}}
	}

	{{end}}{{/* End return values */}}

	return {{range $index, $elem := $f.Getter.ReturnValues}}{{if $index}},{{end}}{{$elem.Name}}{{end}}{{if gt $returnLength 0}},{{end}} nil
}
{{end}}{{/* End Getter */}}
{{if $f.Setter}}
func {{ToGoNameConv $f.Setter.GetNameWithOverloadIndex}}_MetaFFISetter({{$f.Name}} {{ConvertToGoType $f.ArgDefinition $m}}) (err error){
	{{ $paramsLength := len $f.Setter.Parameters }}{{ $returnLength := len $f.Setter.ReturnValues }}

	{{GenerateCodeAllocateCDTS $f.Setter.Parameters $f.Setter.ReturnValues}}
	
	// parameters
	{{range $index, $elem := $f.Setter.Parameters}}
	FromGoToCDT({{$elem.Name}}, parametersCDTS, {{$index}})
	{{end}}{{/* End Parameters */}}

	{{GenerateCodeXCall "" $f.Getter.GetNameWithOverloadIndex $f.Setter.Parameters $f.Setter.ReturnValues}}
	
	{{range $index, $elem := $f.Setter.ReturnValues}}
	{{$elem.Name}}AsInterface := FromCDTToGo(return_valuesCDTS, {{$index}})
	if {{$elem.Name}}AsInterface != nil{
		{{if $elem.IsAny}}
		{{$elem.Name}} = {{$elem.Name}}AsInterface
		
		{{else if not $elem.IsHandle}}
		{{$elem.Name}} = {{if $elem.IsTypeAlias}}{{if $elem.IsArray}}[]{{end}}{{GetTypeOrAlias $elem $m}}{{else}}{{ConvertToGoType $elem $m}}{{end}}({{$elem.Name}}AsInterface.({{ConvertToGoType $elem $m}}))
		{{else}} {{/* Handle */}}		
		{{/* Go object */}}
		
		{{if not $elem.IsArray}}
		if obj, ok := {{$elem.Name}}AsInterface.(Handle); ok{ // None Go object			
			{{$elem.Name}} = {{HandleNoneGoObject $elem $m}}
		} else {
			{{$elem.Name}} = {{if and $elem.IsTypeAlias (not $elem.IsHandleTypeAlias)}}{{if $elem.IsArray}}[]{{end}}{{GetTypeOrAlias $elem $m}}{{else}}{{ConvertToGoType $elem $m}}{{end}}({{$elem.Name}}AsInterface.({{ConvertToGoType $elem $m}}))
		}
		{{else}}
		{{if $elem.IsTypeAlias}} {{/* a type is specified */}}
		if len({{$elem.Name}}AsInterface.([]interface{})) > 0{
			{{$elem.Name}} = make([]{{GetTypeOrAlias $elem $m}}, len({{$elem.Name}}AsInterface.([]interface{})))
			if _, ok := {{$elem.Name}}AsInterface.([]interface{})[0].(Handle); ok{
				for i, h := range {{$elem.Name}}AsInterface.([]interface{}){
					{{$elem.Name}}[i] = {{GetTypeOrAlias $elem $m}}{ h: h.(Handle) }
				}
			} else {
				for i, obj := range {{$elem.Name}}AsInterface.([]interface{}){
					{{$elem.Name}}[i] = obj.({{GetTypeOrAlias $elem $m}})
				}
			}
		}
		{{else}}
		{{$elem.Name}} = {{$elem.Name}}AsInterface
		{{end}}
		{{end}}

		{{end}}{{/* end handling types */}}
	}
	
	{{end}}{{/* End return values */}}

	return {{range $index, $elem := $f.Setter.ReturnValues}}{{if $index}},{{end}}{{$elem.Name}}{{end}}{{if gt $returnLength 0}},{{end}} nil
}
{{end}}{{/* End Setter */}}
{{end}}{{/* End Global */}}


// Code to call foreign functions in module {{$m.Name}} via XLLR
{{range $findex, $f := $m.Functions}}
// Call to foreign {{$f.Name}}
{{if $f.Comment}}/*
{{$f.Comment}}
*/{{end}}
{{range $index, $elem := $f.Parameters}}
{{if $elem.Comment}}// {{$elem.Name}} - {{$elem.Comment}}{{end}}{{end}}{{/* End Parameters comments */}}
func {{ToGoNameConv $f.GetNameWithOverloadIndex}}({{range $index, $elem := $f.Parameters}}{{if $index}},{{end}} {{$elem.Name}} {{ConvertToGoType $elem $m}}{{end}}) ({{range $index, $elem := $f.ReturnValues}}{{if $index}},{{end}}{{$elem.Name}} {{ConvertToGoType $elem $m}}{{end}}{{if $f.ReturnValues}},{{end}} err error){

	{{ $paramsLength := len $f.Parameters }}{{ $returnLength := len $f.ReturnValues }}

	{{GenerateCodeAllocateCDTS $f.Parameters $f.ReturnValues}}
	
	// parameters
	{{range $index, $elem := $f.Parameters}}
	FromGoToCDT({{$elem.Name}}, parametersCDTS, {{$index}})
	{{end}}{{/* End Parameters */}}

	{{GenerateCodeXCall "" $f.GetNameWithOverloadIndex $f.Parameters $f.ReturnValues}}
	
	{{range $index, $elem := $f.ReturnValues}}
	{{$elem.Name}}AsInterface := FromCDTToGo(return_valuesCDTS, {{$index}})
	if {{$elem.Name}}AsInterface != nil{
		{{if $elem.IsAny}}
		{{$elem.Name}} = {{$elem.Name}}AsInterface
		
		{{else if not $elem.IsHandle}}
		{{$elem.Name}} = {{if $elem.IsTypeAlias}}{{if $elem.IsArray}}[]{{end}}{{GetTypeOrAlias $elem $m}}{{else}}{{ConvertToGoType $elem $m}}{{end}}({{$elem.Name}}AsInterface.({{ConvertToGoType $elem $m}}))
		{{else}} {{/* Handle */}}		
		{{/* Go object */}}
		
		{{if not $elem.IsArray}}
		if obj, ok := {{$elem.Name}}AsInterface.(Handle); ok{ // None Go object			
			{{$elem.Name}} = {{HandleNoneGoObject $elem $m}}
		} else {
			{{$elem.Name}} = {{if and $elem.IsTypeAlias (not $elem.IsHandleTypeAlias)}}{{if $elem.IsArray}}[]{{end}}{{GetTypeOrAlias $elem $m}}{{else}}{{ConvertToGoType $elem $m}}{{end}}({{$elem.Name}}AsInterface.({{ConvertToGoType $elem $m}}))
		}
		{{else}}
		{{if $elem.IsTypeAlias}} {{/* a type is specified */}}
		if len({{$elem.Name}}AsInterface.([]interface{})) > 0{
			{{$elem.Name}} = make([]{{GetTypeOrAlias $elem $m}}, len({{$elem.Name}}AsInterface.([]interface{})))
			if _, ok := {{$elem.Name}}AsInterface.([]interface{})[0].(Handle); ok{
				for i, h := range {{$elem.Name}}AsInterface.([]interface{}){
					{{$elem.Name}}[i] = {{GetTypeOrAlias $elem $m}}{ h: h.(Handle) }
				}
			} else {
				for i, obj := range {{$elem.Name}}AsInterface.([]interface{}){
					{{$elem.Name}}[i] = obj.({{GetTypeOrAlias $elem $m}})
				}
			}
		}
		{{else}}
		{{$elem.Name}} = {{$elem.Name}}AsInterface
		{{end}}
		{{end}}

		{{end}}{{/* end handling types */}}
	}
	
	{{end}}{{/* End return values */}}

	return {{range $index, $elem := $f.ReturnValues}}{{if $index}},{{end}}{{$elem.Name}}{{end}}{{if gt $returnLength 0}},{{end}} nil
}
{{end}}{{/* End Function */}}

{{range $cindex, $c := $m.Classes}}
type {{AsPublic $c.Name}} struct{
	h Handle
}

{{range $findex, $f := $c.Constructors}}
func New{{ToGoNameConv $f.GetNameWithOverloadIndex}}({{range $index, $elem := $f.Parameters}}{{if $index}},{{end}} {{$elem.Name}} {{ConvertToGoType $elem $m}}{{end}}) (instance *{{AsPublic $c.Name}}, err error){
	{{ $paramsLength := len $f.Parameters }}{{ $returnLength := len $f.ReturnValues }}

	{{GenerateCodeAllocateCDTS $f.Parameters $f.ReturnValues}}
	
	// parameters
	{{range $index, $elem := $f.Parameters}}
	FromGoToCDT({{$elem.Name}}, parametersCDTS, {{$index}})
	{{end}}{{/* End Parameters */}}

	{{GenerateCodeXCall $c.Name $f.GetNameWithOverloadIndex $f.Parameters $f.ReturnValues}}
	
	inst := &{{AsPublic $c.Name}}{}

	{{range $index, $elem := $f.ReturnValues}}
	{{$elem.Name}}AsInterface := FromCDTToGo(return_valuesCDTS, {{$index}})
	if {{$elem.Name}}AsInterface != nil{
		inst.h = {{$elem.Name}}AsInterface.(Handle)
	} else {
		return nil, fmt.Errorf("Object creation returned nil")
	}
		
	{{end}}{{/* End return values */}}

	return inst, nil	
}
{{end}}{{/* End Constructor */}}

func (this *{{AsPublic $c.Name}}) GetHandle() Handle{
	return this.h
}

func (this *{{AsPublic $c.Name}}) SetHandle(h Handle){
	this.h = h
}

{{range $findex, $f := $c.Fields}}
{{if $f.Getter}}
func {{GenerateMethodReceiverCode $f.Getter}} {{GenerateMethodName $f.Getter}}_MetaFFIGetter({{GenerateMethodParams $f.Getter $m}}) ({{range $index, $elem := $f.Getter.ReturnValues}}{{if $index}},{{end}}{{$elem.Name}} {{ConvertToGoType $elem $m}}{{end}}{{if $f.Getter.ReturnValues}},{{end}} err error){
	
	{{ $paramsLength := len $f.Getter.Parameters }}{{ $returnLength := len $f.Getter.ReturnValues }}

	{{GenerateCodeAllocateCDTS $f.Getter.Parameters $f.Getter.ReturnValues}}
	
	// get parameters
	{{if $f.Getter.InstanceRequired}}
	FromGoToCDT(this.h, parametersCDTS, 0)
	{{range $index, $elem := $f.Getter.Parameters}}{{if gt $index 0}}
	FromGoToCDT(this.h, parametersCDTS, {{$index}})
	{{end}}{{end}}{{/* End Parameters */}}
	{{else}}
	{{range $index, $elem := $f.Getter.Parameters}}
	FromGoToCDT(this.h, parametersCDTS, {{$index}})
	{{end}} {{/* End Parameters */}}
	{{end}} {{/* End InstanceRequired */}}


	{{GenerateCodeXCall $c.Name $f.Getter.GetNameWithOverloadIndex $f.Getter.Parameters $f.Getter.ReturnValues}}
	
	{{range $index, $elem := $f.Getter.ReturnValues}}
	{{$elem.Name}}AsInterface := FromCDTToGo(return_valuesCDTS, {{$index}})
	if {{$elem.Name}}AsInterface != nil{
		{{if $elem.IsAny}}
		// any
		{{$elem.Name}} = {{$elem.Name}}AsInterface
		
		{{else if not $elem.IsHandle}}
		// not handle
		{{$elem.Name}} = {{if $elem.IsTypeAlias}}{{if $elem.IsArray}}[]{{end}}{{GetTypeOrAlias $elem $m}}{{else}}{{ConvertToGoType $elem $m}}{{end}}({{$elem.Name}}AsInterface.({{ConvertToGoType $elem $m}}))
		
		{{else}} {{/* Handle */}}
		// handle
		{{if not $elem.IsArray}}
		if obj, ok := {{$elem.Name}}AsInterface.(Handle); ok{ // None Go object			
			{{$elem.Name}} = {{HandleNoneGoObject $elem $m}}
		} else {
			{{$elem.Name}} = {{if and $elem.IsTypeAlias (not $elem.IsHandleTypeAlias)}}{{if $elem.IsArray}}[]{{end}}{{GetTypeOrAlias $elem $m}}{{else}}{{ConvertToGoType $elem $m}}{{end}}({{$elem.Name}}AsInterface.({{ConvertToGoType $elem $m}}))
		}
		{{else}}
		{{if $elem.IsTypeAlias}} {{/* a type is specified */}}
		if len({{$elem.Name}}AsInterface.([]interface{})) > 0{
			{{$elem.Name}} = make([]{{GetTypeOrAlias $elem $m}}, len({{$elem.Name}}AsInterface.([]interface{})))
			if _, ok := {{$elem.Name}}AsInterface.([]interface{})[0].(Handle); ok{
				for i, h := range {{$elem.Name}}AsInterface.([]interface{}){
					{{$elem.Name}}[i] = {{GetTypeOrAlias $elem $m}}{ h: h.(Handle) }
				}
			} else {
				for i, obj := range {{$elem.Name}}AsInterface.([]interface{}){
					{{$elem.Name}}[i] = obj.({{GetTypeOrAlias $elem $m}})
				}
			}
		}
		{{else}}
		{{$elem.Name}} = {{$elem.Name}}AsInterface
		{{end}}
		{{end}}

		{{end}}{{/* end handling types */}}
	}

	{{end}}{{/* End return values */}}

	return {{range $index, $elem := $f.Getter.ReturnValues}}{{if $index}},{{end}}{{$elem.Name}}{{end}}{{if gt $returnLength 0}},{{end}} nil	
}
{{end}}{{/* End Getter */}}
{{if $f.Setter}}
func {{GenerateMethodReceiverCode $f.Setter}} {{GenerateMethodName $f.Setter}}_MetaFFISetter({{GenerateMethodParams $f.Setter $m}}) ({{range $index, $elem := $f.Setter.ReturnValues}}{{if $index}},{{end}}{{$elem.Name}} {{ConvertToGoType $elem $m}}{{end}}{{if $f.Setter.ReturnValues}},{{end}} err error){
	
	{{ $paramsLength := len $f.Setter.Parameters }}{{ $returnLength := len $f.Setter.ReturnValues }}

	{{GenerateCodeAllocateCDTS $f.Setter.Parameters $f.Setter.ReturnValues}}
	
	// parameters
	FromGoToCDT(this.h, parametersCDTS, 0) // object
	{{range $index, $elem := $f.Setter.Parameters}}{{if gt $index 0}}
	FromGoToCDT({{$elem.Name}}, parametersCDTS, {{$index}})
	{{end}}{{end}}{{/* End Parameters */}}

	{{GenerateCodeXCall $c.Name $f.Setter.GetNameWithOverloadIndex $f.Setter.Parameters $f.Setter.ReturnValues}}
	
	{{range $index, $elem := $f.Setter.ReturnValues}}
	{{$elem.Name}}AsInterface := FromCDTToGo(return_valuesCDTS, {{$index}})
	if {{$elem.Name}}AsInterface != nil{
		{{if $elem.IsAny}}
		{{$elem.Name}} = {{$elem.Name}}AsInterface
		
		{{else if not $elem.IsHandle}}
		{{$elem.Name}} = {{if $elem.IsTypeAlias}}{{if $elem.IsArray}}[]{{end}}{{GetTypeOrAlias $elem $m}}{{else}}{{ConvertToGoType $elem $m}}{{end}}({{$elem.Name}}AsInterface.({{ConvertToGoType $elem $m}}))
		{{else}} {{/* Handle */}}		
		{{/* Go object */}}
		
		{{if not $elem.IsArray}}
		if obj, ok := {{$elem.Name}}AsInterface.(Handle); ok{ // None Go object			
			{{$elem.Name}} = {{HandleNoneGoObject $elem $m}}
		} else {
			{{$elem.Name}} = {{if and $elem.IsTypeAlias (not $elem.IsHandleTypeAlias)}}{{if $elem.IsArray}}[]{{end}}{{GetTypeOrAlias $elem $m}}{{else}}{{ConvertToGoType $elem $m}}{{end}}({{$elem.Name}}AsInterface.({{ConvertToGoType $elem $m}}))
		}
		{{else}}
		{{if $elem.IsTypeAlias}} {{/* a type is specified */}}
		if len({{$elem.Name}}AsInterface.([]interface{})) > 0{
			{{$elem.Name}} = make([]{{GetTypeOrAlias $elem $m}}, len({{$elem.Name}}AsInterface.([]interface{})))
			if _, ok := {{$elem.Name}}AsInterface.([]interface{})[0].(Handle); ok{
				for i, h := range {{$elem.Name}}AsInterface.([]interface{}){
					{{$elem.Name}}[i] = {{GetTypeOrAlias $elem $m}}{ h: h.(Handle) }
				}
			} else {
				for i, obj := range {{$elem.Name}}AsInterface.([]interface{}){
					{{$elem.Name}}[i] = obj.({{GetTypeOrAlias $elem $m}})
				}
			}
		}
		{{else}}
		{{$elem.Name}} = {{$elem.Name}}AsInterface
		{{end}}
		{{end}}
		
		{{end}}{{/* end handling types */}}
	}
	
	{{end}}{{/* End return values */}}

	return {{range $index, $elem := $f.Setter.ReturnValues}}{{if $index}},{{end}}{{$elem.Name}}{{end}}{{if gt $returnLength 0}},{{end}} nil	
}
{{end}}{{/* End Setter */}}
{{end}}{{/* End Fields */}}
{{range $findex, $f := $c.Methods}}
func {{GenerateMethodReceiverCode $f}} {{GenerateMethodName $f}}({{GenerateMethodParams $f $m}}) ({{range $index, $elem := $f.ReturnValues}}{{if $index}},{{end}}{{$elem.Name}} {{ConvertToGoType $elem $m}}{{end}}{{if $f.ReturnValues}},{{end}} err error){
	
	{{ $paramsLength := len $f.Parameters }}{{ $returnLength := len $f.ReturnValues }}

	{{GenerateCodeAllocateCDTS $f.Parameters $f.ReturnValues}}
	
	// parameters
	{{if $f.InstanceRequired}}
	FromGoToCDT(this.h, parametersCDTS, 0) // object
	{{range $index, $elem := $f.Parameters}}{{if gt $index 0}}
	FromGoToCDT({{$elem.Name}}, parametersCDTS, {{$index}})
	{{end}}{{end}}
	{{else}}
	{{range $index, $elem := $f.Parameters}}
	FromGoToCDT({{$elem.Name}}, parametersCDTS, {{$index}})
	{{end}}
	{{end}}

	{{GenerateCodeXCall $c.Name $f.GetNameWithOverloadIndex $f.Parameters $f.ReturnValues}}
	
	{{range $index, $elem := $f.ReturnValues}}
	{{$elem.Name}}AsInterface := FromCDTToGo(return_valuesCDTS, {{$index}})
	if {{$elem.Name}}AsInterface != nil{
		{{if $elem.IsAny}}
		{{$elem.Name}} = {{$elem.Name}}AsInterface
		
		{{else if not $elem.IsHandle}}
		{{$elem.Name}} = {{if $elem.IsTypeAlias}}{{if $elem.IsArray}}[]{{end}}{{GetTypeOrAlias $elem $m}}{{else}}{{ConvertToGoType $elem $m}}{{end}}({{$elem.Name}}AsInterface.({{ConvertToGoType $elem $m}}))
		{{else}} {{/* Handle */}}		
		{{/* Go object */}}
		
		{{if not $elem.IsArray}}
		if obj, ok := {{$elem.Name}}AsInterface.(Handle); ok{ // None Go object			
			{{$elem.Name}} = {{HandleNoneGoObject $elem $m}}
		} else {
			{{$elem.Name}} = {{if and $elem.IsTypeAlias (not $elem.IsHandleTypeAlias)}}{{if $elem.IsArray}}[]{{end}}{{GetTypeOrAlias $elem $m}}{{else}}{{ConvertToGoType $elem $m}}{{end}}({{$elem.Name}}AsInterface.({{ConvertToGoType $elem $m}}))
		}
		{{else}}
		{{if $elem.IsTypeAlias}} {{/* a type is specified */}}
		if len({{$elem.Name}}AsInterface.([]interface{})) > 0{
			{{$elem.Name}} = make([]{{GetTypeOrAlias $elem $m}}, len({{$elem.Name}}AsInterface.([]interface{})))
			if _, ok := {{$elem.Name}}AsInterface.([]interface{})[0].(Handle); ok{
				for i, h := range {{$elem.Name}}AsInterface.([]interface{}){
					{{$elem.Name}}[i] = {{GetTypeOrAlias $elem $m}}{ h: h.(Handle) }
				}
			} else {
				for i, obj := range {{$elem.Name}}AsInterface.([]interface{}){
					{{$elem.Name}}[i] = obj.({{GetTypeOrAlias $elem $m}})
				}
			}
		}
		{{else}}
		{{$elem.Name}} = {{$elem.Name}}AsInterface
		{{end}}
		{{end}}

		{{end}}{{/* end handling types */}}
	}
	
	{{end}}{{/* End return values */}}

	return {{range $index, $elem := $f.ReturnValues}}{{if $index}},{{end}}{{$elem.Name}}{{end}}{{if gt $returnLength 0}},{{end}} nil
}
{{end}}{{/* End Methods */}}
{{if $c.Releaser}}{{ $f := $c.Releaser}}
func (this *{{AsPublic $c.Name}}) {{ToGoNameConv $f.GetNameWithOverloadIndex}}({{range $index, $elem := $f.Parameters}}{{if $index}},{{end}} {{$elem.Name}} {{ConvertToGoType $elem $m}}{{end}}) ({{range $index, $elem := $f.ReturnValues}}{{if $index}},{{end}}{{$elem.Name}} {{ConvertToGoType $elem $m}}{{end}}{{if $f.ReturnValues}},{{end}} err error){
	{{ $paramsLength := len $f.Parameters }}{{ $returnLength := len $f.ReturnValues }}
	{{GenerateCodeAllocateCDTS $f.Parameters $f.ReturnValues}}
	
	// parameters
	FromGoToCDT(this.h, parametersCDTS, 0) // object
	{{range $index, $elem := $f.Parameters}}{{if gt $index 0}}
	FromGoToCDT({{$elem.Name}}, parametersCDTS, {{$index}})
	{{end}}{{end}}{{/* End Parameters */}}

	{{GenerateCodeXCall $c.Name $f.GetNameWithOverloadIndex $f.Parameters $f.ReturnValues}}
	
	{{range $index, $elem := $f.ReturnValues}}
	{{$elem.Name}}AsInterface := FromCDTToGo(return_valuesCDTS, {{$index}})
	if {{$elem.Name}}AsInterface != nil{
		{{if $elem.IsAny}}
		{{$elem.Name}} = {{$elem.Name}}AsInterface
		
		{{else if not $elem.IsHandle}}
		{{$elem.Name}} = {{if $elem.IsTypeAlias}}{{if $elem.IsArray}}[]{{end}}{{GetTypeOrAlias $elem $m}}{{else}}{{ConvertToGoType $elem $m}}{{end}}({{$elem.Name}}AsInterface.({{ConvertToGoType $elem $m}}))
		{{else}} {{/* Handle */}}		
		{{/* Go object */}}
		
		{{if not $elem.IsArray}}
		if obj, ok := {{$elem.Name}}AsInterface.(Handle); ok{ // None Go object			
			{{$elem.Name}} = {{HandleNoneGoObject $elem $m}}
		} else {
			{{$elem.Name}} = {{if and $elem.IsTypeAlias (not $elem.IsHandleTypeAlias)}}{{if $elem.IsArray}}[]{{end}}{{GetTypeOrAlias $elem $m}}{{else}}{{ConvertToGoType $elem $m}}{{end}}({{$elem.Name}}AsInterface.({{ConvertToGoType $elem $m}}))
		}
		{{else}}
		{{if $elem.IsTypeAlias}} {{/* a type is specified */}}
		if len({{$elem.Name}}AsInterface.([]interface{})) > 0{
			{{$elem.Name}} = make([]{{GetTypeOrAlias $elem $m}}, len({{$elem.Name}}AsInterface.([]interface{})))
			if _, ok := {{$elem.Name}}AsInterface.([]interface{})[0].(Handle); ok{
				for i, h := range {{$elem.Name}}AsInterface.([]interface{}){
					{{$elem.Name}}[i] = {{GetTypeOrAlias $elem $m}}{ h: h.(Handle) }
				}
			} else {
				for i, obj := range {{$elem.Name}}AsInterface.([]interface{}){
					{{$elem.Name}}[i] = obj.({{GetTypeOrAlias $elem $m}})
				}
			}
		}
		{{else}}
		{{$elem.Name}} = {{$elem.Name}}AsInterface
		{{end}}
		{{end}}

		{{end}}{{/* end handling types */}}
	}
	
	{{end}}{{/* End return values */}}

	return {{range $index, $elem := $f.ReturnValues}}{{if $index}},{{end}}{{$elem.Name}}{{end}}{{if gt $returnLength 0}},{{end}} nil
}
{{end}}{{/* End Releaser */}}
{{end}}{{/* End Classes */}}
{{end}}{{/* End modules */}}

`
