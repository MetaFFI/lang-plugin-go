module github.com/MetaFFI/lang-plugin-go/compiler

go 1.16

require (
	github.com/GreenFuze/go-parser v0.0.0-20230720034525-39975f58f232 // indirect
	github.com/Masterminds/sprig/v3 v3.2.3 // indirect
	github.com/MetaFFI/lang-plugin-go/idl v0.0.0-20230720034620-0023f0398e0e
	github.com/MetaFFI/plugin-sdk v0.0.0-20230620130440-aa38ebb77ba2
	golang.org/x/text v0.11.0
)

replace github.com/MetaFFI/plugin-sdk => ../plugin-sdk

replace github.com/MetaFFI/lang-plugin-go/go-runtime => ../go-runtime

replace github.com/MetaFFI/lang-plugin-go/idl => ../idl
