module github.com/MetaFFI/lang-plugin-go/api

go 1.21.4

require (
	github.com/MetaFFI/lang-plugin-go/go-runtime v0.0.0-20240528185259-d34b330dbc2d
	github.com/MetaFFI/plugin-sdk v0.1.2
)

require (
	github.com/GreenFuze/go-parser v0.0.0-20240113080500-83cff8210dd0 // indirect
	github.com/MetaFFI/lang-plugin-go/idl v0.0.0-20240303182820-2df581898d4d // indirect
	github.com/timandy/routine v1.1.3 // indirect
	golang.org/x/mod v0.15.0 // indirect
	golang.org/x/text v0.15.0 // indirect
	golang.org/x/tools v0.18.0 // indirect
)

replace github.com/MetaFFI/lang-plugin-go/go-runtime => ../go-runtime
replace github.com/MetaFFI/plugin-sdk => ../plugin-sdk