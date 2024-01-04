module github.com/MetaFFI/lang-plugin-go/idl

go 1.16

require (
	github.com/GreenFuze/go-parser v0.0.0-20231114070054-7b0168006eb8
	github.com/MetaFFI/plugin-sdk v0.0.0-20240104091413-269ab95f95ad
)

replace github.com/MetaFFI/plugin-sdk => ../plugin-sdk

replace github.com/GreenFuze/go-parser => ../../../GreenFuze/go-parser
