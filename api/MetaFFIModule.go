package api

import "C"
import (
	"fmt"
	goruntime "github.com/MetaFFI/lang-plugin-go/go-runtime"
	"github.com/MetaFFI/plugin-sdk/compiler/go/IDL"
	"unsafe"
)

type MetaFFIModule struct {
	runtime    *MetaFFIRuntime
	modulePath string
}

func (this *MetaFFIModule) Load(functionPath string, paramsMetaFFITypes []IDL.MetaFFIType, retvalMetaFFITypes []IDL.MetaFFIType) (ff func(...interface{}) ([]interface{}, error), err error) {

	var paramTypes []uint64
	if paramsMetaFFITypes != nil {
		paramTypes = make([]uint64, len(paramsMetaFFITypes))
	}

	var retvalTypes []uint64
	if retvalMetaFFITypes != nil {
		retvalTypes = make([]uint64, len(retvalMetaFFITypes))
	}

	for i, p := range paramsMetaFFITypes {
		paramTypes[i] = IDL.TypeStringToTypeEnum[p]
	}

	for i, r := range retvalMetaFFITypes {
		retvalTypes[i] = IDL.TypeStringToTypeEnum[r]
	}

	var pff *unsafe.Pointer
	pff, err = goruntime.XLLRLoadFunction(this.runtime.runtimePlugin, this.modulePath, functionPath, paramTypes, retvalTypes)
	if err != nil { // failed
		return
	}

	ff = func(params ...interface{}) (retvals []interface{}, err error) {

		if len(params) != len(paramsMetaFFITypes) {
			return nil, fmt.Errorf("Expecting %v parameters, received %v parameters", len(paramsMetaFFITypes), len(params))
		}

		xcall_params, parametersCDTS, return_valuesCDTS := goruntime.XLLRAllocCDTSBuffer(goruntime.IntToMetaFFISize(len(params)), goruntime.IntToMetaFFISize(len(retvalMetaFFITypes)))

		paramsCount := len(params)
		retvalCount := len(retvalMetaFFITypes)

		if paramsCount > 0 {
			for i, p := range params {
				goruntime.FromGoToCDT(p, parametersCDTS, i)
			}
		}

		if paramsCount > 0 && retvalCount > 0 {
			err = goruntime.XLLRXCallParamsRet(pff, xcall_params)
		} else if paramsCount > 0 && retvalCount == 0 {
			err = goruntime.XLLRXCallParamsNoRet(pff, xcall_params)
		} else if paramsCount == 0 && retvalCount > 0 {
			err = goruntime.XLLRXCallNoParamsRet(pff, xcall_params)
		} else {
			err = goruntime.XLLRXCallNoParamsNoRet(pff)
		}

		if err != nil {
			return nil, err
		}

		if retvalCount == 0 { // no return values
			return
		}

		retvals = make([]interface{}, retvalCount)
		for i := 0; i < int(retvalCount); i++ {
			retvals[i] = goruntime.FromCDTToGo(return_valuesCDTS, i)
		}

		if goruntime.GetCacheSize() < paramsCount+retvalCount {
			goruntime.CFree(xcall_params)
		}

		return
	}

	return
}
