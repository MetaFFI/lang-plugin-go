module github.com/MetaFFI/lang-plugin-go/idl

go 1.16

require (
	github.com/GreenFuze/go-parser v0.0.0-20240113080500-83cff8210dd0
	github.com/MetaFFI/plugin-sdk v0.0.0-20240111115655-a2f1bf60dbf5
)

replace github.com/MetaFFI/plugin-sdk => ../plugin-sdk

replace github.com/GreenFuze/go-parser => ../../../GreenFuze/go-parser
