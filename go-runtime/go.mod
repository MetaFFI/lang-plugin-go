module github.com/MetaFFI/lang-plugin-go/go-runtime

go 1.21

toolchain go1.22.0

require (
	github.com/MetaFFI/plugin-sdk v0.1.2
	github.com/timandy/routine v1.1.3
	golang.org/x/text v0.15.0
)

// for tests
//replace github.com/MetaFFI/plugin-sdk => ../plugin-sdk
