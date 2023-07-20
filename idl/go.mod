module github.com/MetaFFI/lang-plugin-go/idl

go 1.16

require (
	github.com/GreenFuze/go-parser v0.0.0-20230319143659-8b32f99912ea
	github.com/MetaFFI/plugin-sdk v0.0.0-20220624151549-e43fd804ca7d
)

replace github.com/MetaFFI/plugin-sdk => ../plugin-sdk
replace github.com/GreenFuze/go-parser => ../../../GreenFuze/go-parser