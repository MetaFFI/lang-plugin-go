package main

import (
	"fmt"
	"github.com/OpenFFI/plugin-sdk/compiler/go"
	"html/template"
	"strings"
)

//--------------------------------------------------------------------
type Compiler struct{
	def *compiler.IDLDefinition
	serializationCode map[string]string
}
//--------------------------------------------------------------------
func NewCompiler(def *compiler.IDLDefinition, serializationCode map[string]string) *Compiler {
	return &Compiler{def: def, serializationCode: serializationCode}
}
//--------------------------------------------------------------------
func (this *Compiler) CompileGuest() (string, error){

	temp, err := template.New("guest").Parse(GuestTemplate)
	if err != nil{
		return "", fmt.Errorf("Failed to parse guest template. error: %v", err)
	}

	strbuf := strings.Builder{}

	err = temp.Execute(&strbuf, this)
	if err != nil{
		return "", fmt.Errorf("Failed to build guest template, err: %v", err)
	}

	return strbuf.String(), nil
}
//--------------------------------------------------------------------
// @protobufFileName - The name of the protobuf python generated
func (this *Compiler) CompileHost() (string, error){

	temp, err := template.New("host").Parse(HostTemplate)
	if err != nil{
		return "", fmt.Errorf("Failed to parse host template. error: %v", err)
	}

	strbuf := strings.Builder{}

	err = temp.Execute(&strbuf, this)
	if err != nil{
		return "", fmt.Errorf("Failed to build host template, err: %v", err)
	}

	return strbuf.String(), nil
}
//--------------------------------------------------------------------