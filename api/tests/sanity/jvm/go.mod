module jvm

require (
	github.com/MetaFFI/lang-plugin-go/api v0.3.1-0.20250411075025-78733fa107b1
	github.com/MetaFFI/sdk/idl_entities/go v0.0.0
)

require (
	github.com/GreenFuze/go-parser v0.0.0-20240414081507-f3897e769e4c // indirect
	github.com/MetaFFI/lang-plugin-go/compiler v0.3.1-0.20250411075025-78733fa107b1 // indirect
	github.com/MetaFFI/lang-plugin-go/go-runtime v0.3.1-0.20250411075025-78733fa107b1 // indirect
	github.com/MetaFFI/lang-plugin-go/idl v0.0.0-20240630050947-d2f24c54ac0b // indirect
	github.com/timandy/routine v1.1.4 // indirect
	golang.org/x/mod v0.17.0 // indirect
	golang.org/x/sync v0.8.0 // indirect
	golang.org/x/text v0.18.0 // indirect
	golang.org/x/tools v0.21.1-0.20240508182429-e35e4ccd0d2d // indirect
)


// for debug - uncomment
//replace github.com/MetaFFI/lang-plugin-go/api => ../../../../api
//replace github.com/MetaFFI/lang-plugin-go/go-runtime => ../../../../go-runtime

go 1.23.0

toolchain go1.23.1
