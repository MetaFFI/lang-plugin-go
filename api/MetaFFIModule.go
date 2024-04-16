package api

import "C"
import (
	"fmt"
	"unsafe"

	goruntime "github.com/MetaFFI/lang-plugin-go/go-runtime"
	"github.com/MetaFFI/plugin-sdk/compiler/go/IDL"
)

type MetaFFIModule struct {
	runtime    *MetaFFIRuntime
	modulePath string
}

func (this *MetaFFIModule) Load(functionPath string, paramsMetaFFITypes []IDL.MetaFFIType, retvalMetaFFITypes []IDL.MetaFFIType) (ff func(...interface{}) ([]interface{}, error), err error) {
	var params []IDL.MetaFFITypeInfo
	if paramsMetaFFITypes != nil {
		params = make([]IDL.MetaFFITypeInfo, 0)
		for _, p := range paramsMetaFFITypes {
			params = append(params, IDL.MetaFFITypeInfo{StringType: p})
		}
	}

	var retvals []IDL.MetaFFITypeInfo
	if retvalMetaFFITypes != nil {
		retvals = make([]IDL.MetaFFITypeInfo, 0)
		for _, r := range retvalMetaFFITypes {
			retvals = append(retvals, IDL.MetaFFITypeInfo{StringType: r})
		}
	}

	return this.LoadWithAlias(functionPath, params, retvals)
}

func (this *MetaFFIModule) LoadWithAlias(functionPath string, paramsMetaFFITypes []IDL.MetaFFITypeInfo, retvalMetaFFITypes []IDL.MetaFFITypeInfo) (ff func(...interface{}) ([]interface{}, error), err error) {

	// convert Go's String metaffi types to INT metaffi types
	if paramsMetaFFITypes != nil {
		for i := 0; i < len(paramsMetaFFITypes); i++ {
			item := paramsMetaFFITypes[i]
			item.FillMetaFFITypeFromStringMetaFFIType()
			paramsMetaFFITypes[i] = item
		}
	}

	if retvalMetaFFITypes != nil {
		for i := 0; i < len(retvalMetaFFITypes); i++ {
			item := retvalMetaFFITypes[i]
			item.FillMetaFFITypeFromStringMetaFFIType()
			retvalMetaFFITypes[i] = item
		}
	}

	var pff *unsafe.Pointer
	pff, err = goruntime.XLLRLoadFunctionWithAliases(this.runtime.runtimePlugin, this.modulePath, functionPath, paramsMetaFFITypes, retvalMetaFFITypes)

	if err != nil { // failed
		return
	}

	ff = func(params ...interface{}) (retvals []interface{}, err error) {

		if len(params) != len(paramsMetaFFITypes) {
			return nil, fmt.Errorf("Expecting %v parameters, received %v parameters", len(paramsMetaFFITypes), len(params))
		}

		xcall_params, parametersCDTS, return_valuesCDTS := goruntime.XLLRAllocCDTSBuffer(uint64(len(params)), uint64(len(retvalMetaFFITypes)))

		paramsCount := len(params)
		retvalCount := len(retvalMetaFFITypes)

		if paramsCount > 0 {
			for i, p := range params {
				goruntime.FromGoToCDT(p, parametersCDTS, paramsMetaFFITypes[i], i)
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
			retvals[i] = goruntime.FromCDTToGo(return_valuesCDTS, i, nil)
		}

		if goruntime.GetCacheSize() < paramsCount+retvalCount {
			goruntime.CFree(xcall_params)
		}

		return
	}

	return
}
