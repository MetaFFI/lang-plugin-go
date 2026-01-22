module github.com/MetaFFI/lang-plugin-go/api

go 1.23.0

toolchain go1.24.4

require (
	github.com/MetaFFI/lang-plugin-go/go-runtime v0.3.0
	github.com/MetaFFI/sdk/idl_entities/go v0.0.0
)

require (
	github.com/timandy/routine v1.1.5 // indirect
	golang.org/x/text v0.27.0 // indirect
)

replace github.com/MetaFFI/lang-plugin-go/go-runtime => ../go-runtime

replace github.com/MetaFFI/sdk/idl_entities/go => ../../sdk/idl_entities/go
