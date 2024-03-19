module github.com/MetaFFI/lang-plugin-go/compiler

go 1.16

require (
	github.com/MetaFFI/lang-plugin-go/go-runtime v0.0.0-20240104094200-93bfb07792f3
	github.com/MetaFFI/lang-plugin-go/idl v0.0.0-20240308210859-fe74e9c86f50
	github.com/MetaFFI/plugin-sdk v0.0.0-20240318182623-6d89a99d0061
	github.com/pkg/profile v1.7.0
	golang.org/x/text v0.14.0
)

replace github.com/MetaFFI/plugin-sdk => ../plugin-sdk

replace github.com/MetaFFI/lang-plugin-go/go-runtime => ../go-runtime

replace github.com/MetaFFI/lang-plugin-go/idl => ../idl

replace github.com/GreenFuze/go-parser => ../../../GreenFuze/go-parser
