module github.com/MetaFFI/lang-plugin-go/compiler

go 1.16

require (
	github.com/MetaFFI/lang-plugin-go/go-runtime v0.0.0-20240320181653-de27c21f7ebc
	github.com/MetaFFI/lang-plugin-go/idl v0.0.0-20240321074310-7a53880ea54b
	github.com/MetaFFI/plugin-sdk v0.0.0-20240319194700-7aa7e30c4fb3
	github.com/pkg/profile v1.7.0
	golang.org/x/text v0.14.0
)

replace github.com/MetaFFI/plugin-sdk => ../plugin-sdk

replace github.com/MetaFFI/lang-plugin-go/go-runtime => ../go-runtime

replace github.com/MetaFFI/lang-plugin-go/idl => ../idl

replace github.com/GreenFuze/go-parser => ../../../GreenFuze/go-parser
