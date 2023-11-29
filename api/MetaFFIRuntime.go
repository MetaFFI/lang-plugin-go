package api

import (
	goruntime "github.com/MetaFFI/lang-plugin-go/go-runtime"
	"os"
)

func init() {
	goruntime.LoadCDTCAPI()
}

type MetaFFIRuntime struct {
	runtimePlugin string
}

func NewMetaFFIRuntime(runtimePlugin string) *MetaFFIRuntime {
	return &MetaFFIRuntime{runtimePlugin: "xllr." + runtimePlugin}
}

func (this *MetaFFIRuntime) LoadRuntimePlugin() error {
	return goruntime.XLLRLoadRuntimePlugin(this.runtimePlugin)
}

func (this *MetaFFIRuntime) ReleaseRuntimePlugin() error {
	return goruntime.XLLRFreeRuntimePlugin(this.runtimePlugin)
}

func (this *MetaFFIRuntime) LoadModule(modulePath string) (*MetaFFIModule, error) {
	_, err := os.Stat(modulePath)
	if err != nil {
		return nil, err
	}

	return &MetaFFIModule{
		runtime:    this,
		modulePath: modulePath,
	}, nil
}
