module github.com/MetaFFI/lang-plugin-go/compiler

go 1.21

toolchain go1.22.0

require (
	github.com/MetaFFI/lang-plugin-go/idl v0.0.0-20240528133705-72be7ae7dd0f
	github.com/MetaFFI/plugin-sdk v0.1.2
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
