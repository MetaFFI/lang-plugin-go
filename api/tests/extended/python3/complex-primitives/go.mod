module complex-primitives

go 1.23.0

toolchain go1.23.1

require (
	github.com/MetaFFI/lang-plugin-go/api v0.3.1-0.20250411075025-78733fa107b1
	github.com/MetaFFI/lang-plugin-go/go-runtime v0.3.1-0.20250411075025-78733fa107b1
	github.com/MetaFFI/sdk/idl_entities/go v0.0.0
)

replace github.com/MetaFFI/sdk/idl_entities/go => ../../../../../../sdk/idl_entities/go

// for debug - uncomment
//replace github.com/MetaFFI/lang-plugin-go/api => ../../../../api
//replace github.com/MetaFFI/lang-plugin-go/go-runtime => ../../../../go-runtime
