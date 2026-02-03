module log4j

go 1.21.4

require (
	github.com/MetaFFI/lang-plugin-go/api v0.3.1-0.20250406132325-4454c0dd3c0a
	github.com/MetaFFI/lang-plugin-go/go-runtime v0.3.1-0.20250406132325-4454c0dd3c0a
	github.com/MetaFFI/sdk/idl_entities/go v0.0.0
)

require (
	github.com/MetaFFI/lang-plugin-go/compiler v0.3.1-0.20250406132325-4454c0dd3c0a // indirect
	github.com/timandy/routine v1.1.4 // indirect
	golang.org/x/text v0.18.0 // indirect
)


// for debug - uncomment
//replace github.com/MetaFFI/lang-plugin-go/api => ../../../../api
//replace github.com/MetaFFI/lang-plugin-go/go-runtime => ../../../../go-runtime
