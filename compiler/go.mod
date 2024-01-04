module github.com/MetaFFI/lang-plugin-go/compiler

go 1.16

require (
	github.com/GreenFuze/go-parser v0.0.0-20231114070054-7b0168006eb8 // indirect
	github.com/MetaFFI/lang-plugin-go/go-runtime v0.0.0-20231230163933-89db4845c863
	github.com/MetaFFI/lang-plugin-go/idl v0.0.0-20231230163933-89db4845c863
	github.com/MetaFFI/plugin-sdk v0.0.0-20231223172124-e3162c9aa768
	github.com/google/pprof v0.0.0-20231229205709-960ae82b1e42 // indirect
	github.com/pkg/profile v1.7.0
	golang.org/x/text v0.14.0
)

replace github.com/MetaFFI/plugin-sdk => ../plugin-sdk

replace github.com/MetaFFI/lang-plugin-go/go-runtime => ../go-runtime

replace github.com/MetaFFI/lang-plugin-go/idl => ../idl

replace github.com/GreenFuze/go-parser => ../../../GreenFuze/go-parser
