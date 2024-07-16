module python3

require (
	github.com/MetaFFI/lang-plugin-go/api v0.0.0-20240704181417-f8f950a95efb
	github.com/MetaFFI/plugin-sdk v0.1.3-0.20240714120957-b4eeb5ecb980
)

require (
	github.com/GreenFuze/go-parser v0.0.0-20240414081507-f3897e769e4c // indirect
	github.com/MetaFFI/lang-plugin-go/compiler v0.0.0-20240704181417-f8f950a95efb // indirect
	github.com/MetaFFI/lang-plugin-go/go-runtime v0.0.0-20240704181417-f8f950a95efb // indirect
	github.com/MetaFFI/lang-plugin-go/idl v0.0.0-20240630050947-d2f24c54ac0b // indirect
	github.com/timandy/routine v1.1.3 // indirect
	golang.org/x/mod v0.17.0 // indirect
	golang.org/x/sync v0.7.0 // indirect
	golang.org/x/text v0.16.0 // indirect
	golang.org/x/tools v0.21.1-0.20240508182429-e35e4ccd0d2d // indirect
)

replace github.com/MetaFFI/lang-plugin-go/api => ../../../../api

replace github.com/MetaFFI/lang-plugin-go/go-runtime => ../../../../go-runtime

go 1.21.4
