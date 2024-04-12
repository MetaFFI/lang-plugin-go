module github.com/MetaFFI/lang-plugin-go/idl

go 1.16

require (
	github.com/GreenFuze/go-parser v0.0.0-20240113080500-83cff8210dd0
	github.com/MetaFFI/plugin-sdk v0.0.0-20240412085519-d3a8fca4727b
	golang.org/x/tools v0.20.0 // indirect
)

replace github.com/MetaFFI/plugin-sdk => ../plugin-sdk

replace github.com/GreenFuze/go-parser => ../../../GreenFuze/go-parser
