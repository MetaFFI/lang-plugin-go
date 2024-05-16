module github.com/MetaFFI/lang-plugin-go/go-runtime

go 1.21

toolchain go1.22.0

require (
	github.com/MetaFFI/plugin-sdk v0.0.0-20240418113454-40cb0644f6c7
	github.com/timandy/routine v1.1.3
	golang.org/x/text v0.14.0
)

replace github.com/MetaFFI/plugin-sdk => ../plugin-sdk