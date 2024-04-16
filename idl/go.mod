module github.com/MetaFFI/lang-plugin-go/idl

go 1.16

require (
	github.com/GreenFuze/go-parser v0.0.0-20240414081507-f3897e769e4c
	github.com/MetaFFI/plugin-sdk v0.0.0-20240416150902-5f975a29af46
	golang.org/x/tools v0.20.0 // indirect
)

replace github.com/MetaFFI/plugin-sdk => ../plugin-sdk

replace github.com/GreenFuze/go-parser => ../../../GreenFuze/go-parser
