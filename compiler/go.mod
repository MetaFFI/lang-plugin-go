module github.com/MetaFFI/lang-plugin-go/compiler

go 1.16

require (
	github.com/MetaFFI/lang-plugin-go/go-runtime v0.0.0-20240104094200-93bfb07792f3
	github.com/MetaFFI/lang-plugin-go/idl v0.0.0-20240304151931-a18399726aa3
	github.com/MetaFFI/plugin-sdk v0.0.0-20240304151550-3a68760a89d1
	github.com/pkg/profile v1.7.0
	golang.org/x/text v0.14.0
	golang.org/x/tools v0.18.0 // indirect
)

replace github.com/MetaFFI/plugin-sdk => ../plugin-sdk

replace github.com/MetaFFI/lang-plugin-go/go-runtime => ../go-runtime

replace github.com/MetaFFI/lang-plugin-go/idl => ../idl

replace github.com/GreenFuze/go-parser => ../../../GreenFuze/go-parser
