module github.com/MetaFFI/lang-plugin-go/compiler

go 1.21

toolchain go1.22.0

require (
	github.com/MetaFFI/lang-plugin-go/go-runtime v0.0.0-20240320181653-de27c21f7ebc
	github.com/MetaFFI/lang-plugin-go/idl v0.0.0-20240520163600-78763ce1ebd9
	github.com/MetaFFI/plugin-sdk v0.1.2
	github.com/pkg/profile v1.7.0
	golang.org/x/text v0.15.0
)

require (
	github.com/GreenFuze/go-parser v0.0.0-20240414081507-f3897e769e4c // indirect
	golang.org/x/mod v0.17.0 // indirect
	golang.org/x/sync v0.7.0 // indirect
	golang.org/x/tools v0.21.0 // indirect
)

replace github.com/MetaFFI/plugin-sdk => ../plugin-sdk

replace github.com/MetaFFI/lang-plugin-go/go-runtime => ../go-runtime

replace github.com/MetaFFI/lang-plugin-go/idl => ../idl

replace github.com/GreenFuze/go-parser => ../../../GreenFuze/go-parser
