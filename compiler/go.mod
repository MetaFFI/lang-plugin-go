module github.com/MetaFFI/lang-plugin-go/compiler

go 1.16

require (
	github.com/MetaFFI/plugin-sdk v0.0.0-00010101000000-000000000000
	golang.org/x/text v0.3.7
)

replace github.com/MetaFFI/plugin-sdk => ../../plugin-sdk

replace github.com/MetaFFI/lang-plugin-go/go-runtime => ../go-runtime
