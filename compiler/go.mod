module github.com/MetaFFI/lang-plugin-go/compiler

go 1.16

require (
	github.com/Masterminds/sprig/v3 v3.2.3 // indirect
	github.com/MetaFFI/plugin-sdk v0.0.0-00010101000000-000000000000
	golang.org/x/text v0.4.0
)

replace github.com/MetaFFI/plugin-sdk => ../plugin-sdk
replace github.com/MetaFFI/lang-plugin-go/go-runtime => ../go-runtime
