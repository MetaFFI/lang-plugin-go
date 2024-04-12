module github.com/MetaFFI/lang-plugin-go/compiler

go 1.21

toolchain go1.22.0

require (
	github.com/MetaFFI/lang-plugin-go/go-runtime v0.0.0-20240320181653-de27c21f7ebc
	github.com/MetaFFI/lang-plugin-go/idl v0.0.0-20240412094403-424634bf312c
	github.com/MetaFFI/plugin-sdk v0.0.0-20240412085519-d3a8fca4727b
	github.com/pkg/profile v1.7.0
	golang.org/x/text v0.14.0
)

require (
	github.com/GreenFuze/go-parser v0.0.0-20240113080500-83cff8210dd0 // indirect
	golang.org/x/mod v0.17.0 // indirect
	golang.org/x/sync v0.7.0 // indirect
	golang.org/x/tools v0.20.0 // indirect
)

replace github.com/MetaFFI/plugin-sdk => ../plugin-sdk

replace github.com/MetaFFI/lang-plugin-go/go-runtime => ../go-runtime

replace github.com/MetaFFI/lang-plugin-go/idl => ../idl

replace github.com/GreenFuze/go-parser => ../../../GreenFuze/go-parser
