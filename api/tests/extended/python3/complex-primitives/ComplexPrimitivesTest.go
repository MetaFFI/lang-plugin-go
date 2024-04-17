package main

import (
	"fmt"
	"github.com/MetaFFI/lang-plugin-go/api"
	metaffi "github.com/MetaFFI/lang-plugin-go/go-runtime"
	"github.com/MetaFFI/plugin-sdk/compiler/go/IDL"
)

var runtime *api.MetaFFIRuntime

//--------------------------------------------------------------------

type PyList struct {
	instance metaffi.MetaFFIHandle
	append   func(...interface{}) ([]interface{}, error)
	get      func(...interface{}) ([]interface{}, error)
	len      func(...interface{}) ([]interface{}, error)
}

func NewPyListFromHandle(h metaffi.MetaFFIHandle) (*PyList, error) {
	this := &PyList{}
	this.instance = h

	mod, err := runtime.LoadModule("builtins")
	if err != nil {
		return nil, err
	}

	this.append, err = mod.Load("callable=list.append,instance_required", []IDL.MetaFFIType{IDL.HANDLE, IDL.ANY}, nil)
	if err != nil {
		return nil, err
	}

	this.get, err = mod.Load("callable=list.__getitem__,instance_required", []IDL.MetaFFIType{IDL.HANDLE, IDL.INT64}, []IDL.MetaFFIType{IDL.ANY})
	if err != nil {
		return nil, err
	}

	this.len, err = mod.Load("callable=list.__len__,instance_required", []IDL.MetaFFIType{IDL.HANDLE}, []IDL.MetaFFIType{IDL.INT64})
	if err != nil {
		return nil, err
	}

	return this, nil
}

func NewPyList() (*PyList, error) {
	this := &PyList{}

	mod, err := runtime.LoadModule("builtins")
	if err != nil {
		return nil, err
	}

	constructor, err := mod.Load("callable=list", nil, []IDL.MetaFFIType{IDL.HANDLE})
	if err != nil {
		return nil, err
	}

	this.append, err = mod.Load("callable=list.append,instance_required", []IDL.MetaFFIType{IDL.HANDLE, IDL.ANY}, nil)
	if err != nil {
		return nil, err
	}

	this.get, err = mod.Load("callable=list.__getitem__,instance_required", []IDL.MetaFFIType{IDL.HANDLE, IDL.INT64}, []IDL.MetaFFIType{IDL.ANY})
	if err != nil {
		return nil, err
	}

	this.len, err = mod.Load("callable=list.__len__,instance_required", []IDL.MetaFFIType{IDL.HANDLE}, []IDL.MetaFFIType{IDL.INT64})
	if err != nil {
		return nil, err
	}

	instance, err := constructor()
	if err != nil {
		return nil, err
	}

	this.instance = instance[0].(metaffi.MetaFFIHandle)

	return this, nil
}

func (this *PyList) Append(obj interface{}) error {
	_, err := this.append(this.instance, obj)
	if err != nil {
		return err
	}

	return err
}

func (this *PyList) Get(i int64) (interface{}, error) {
	res, err := this.get(this.instance, i)
	if err != nil {
		return nil, err
	}

	return res[0], nil
}

func (this *PyList) Len() (int64, error) {
	res, err := this.len(this.instance)
	if err != nil {
		return -1, err
	}

	return res[0].(int64), nil
}

//--------------------------------------------------------------------

type PyDict struct {
	Instance metaffi.MetaFFIHandle
	set      func(...interface{}) ([]interface{}, error)
	get      func(...interface{}) ([]interface{}, error)
	len      func(...interface{}) ([]interface{}, error)
}

func NewPyDictFromHandle(h metaffi.MetaFFIHandle) (*PyDict, error) {
	this := &PyDict{}
	this.Instance = h

	mod, err := runtime.LoadModule("builtins")
	if err != nil {
		return nil, err
	}

	this.set, err = mod.Load("callable=dict.__setitem__,instance_required", []IDL.MetaFFIType{IDL.HANDLE, IDL.STRING8, IDL.ANY}, nil)
	if err != nil {
		return nil, err
	}

	this.get, err = mod.Load("callable=dict.__getitem__,instance_required", []IDL.MetaFFIType{IDL.HANDLE, IDL.STRING8}, []IDL.MetaFFIType{IDL.ANY})
	if err != nil {
		return nil, err
	}

	this.len, err = mod.Load("callable=dict.__len__,instance_required", []IDL.MetaFFIType{IDL.HANDLE}, []IDL.MetaFFIType{IDL.INT64})
	if err != nil {
		return nil, err
	}

	return this, nil
}

func NewPyDict() (*PyDict, error) {
	this := &PyDict{}

	mod, err := runtime.LoadModule("builtins")
	if err != nil {
		return nil, err
	}

	constructor, err := mod.Load("callable=dict", nil, []IDL.MetaFFIType{IDL.HANDLE})
	if err != nil {
		return nil, err
	}

	this.set, err = mod.Load("callable=dict.__setitem__,instance_required", []IDL.MetaFFIType{IDL.HANDLE, IDL.STRING8, IDL.ANY}, nil)
	if err != nil {
		return nil, err
	}

	this.get, err = mod.Load("callable=dict.__getitem__,instance_required", []IDL.MetaFFIType{IDL.HANDLE, IDL.STRING8}, []IDL.MetaFFIType{IDL.ANY})
	if err != nil {
		return nil, err
	}

	this.len, err = mod.Load("callable=dict.__len__,instance_required", []IDL.MetaFFIType{IDL.HANDLE}, []IDL.MetaFFIType{IDL.INT64})
	if err != nil {
		return nil, err
	}

	instance, err := constructor()
	if err != nil {
		return nil, err
	}

	this.Instance = instance[0].(metaffi.MetaFFIHandle)

	return this, nil
}

func (this *PyDict) Set(k string, obj interface{}) error {
	_, err := this.set(this.Instance, k, obj)
	if err != nil {
		return err
	}

	return err
}

func (this *PyDict) Get(k string) (interface{}, error) {
	res, err := this.get(this.Instance, k)
	if err != nil {
		return nil, err
	}

	return res[0], nil
}

func (this *PyDict) Len() (int64, error) {
	res, err := this.len(this.Instance)
	if err != nil {
		return -1, err
	}

	return res[0].(int64), nil
}

//--------------------------------------------------------------------

func main() {

	// load runtime
	runtime = api.NewMetaFFIRuntime("python311")
	err := runtime.LoadRuntimePlugin()
	if err != nil {
		panic(err)
	}
	module, err := runtime.LoadModule("extended_test.py")
	if err != nil {
		panic(err)
	}

	// create object
	constructor, err := module.Load("callable=extended_test", nil, []IDL.MetaFFIType{IDL.HANDLE})
	if err != nil {
		panic(err)
	}

	tmp, err := constructor()
	if err != nil {
		panic(err)
	}

	instance := tmp[0].(metaffi.MetaFFIHandle)

	// load functions
	psetX, err := module.Load("callable=extended_test.x.fset,instance_required",
		[]IDL.MetaFFIType{IDL.HANDLE, IDL.INT64},
		nil)
	if err != nil {
		panic(err)
	}

	pgetX, err := module.Load("callable=extended_test.x.fget,instance_required",
		[]IDL.MetaFFIType{IDL.HANDLE},
		[]IDL.MetaFFIType{IDL.INT64})
	if err != nil {
		panic(err)
	}

	ppositional_or_named, err := module.Load("callable=extended_test.positional_or_named,instance_required",
		[]IDL.MetaFFIType{IDL.HANDLE, IDL.STRING8},
		[]IDL.MetaFFIType{IDL.STRING8})
	if err != nil {
		panic(err)
	}

	ppositional_or_named_as_named, err := module.Load("callable=extended_test.positional_or_named,instance_required,named_args",
		[]IDL.MetaFFIType{IDL.HANDLE, IDL.HANDLE},
		[]IDL.MetaFFIType{IDL.STRING8})
	if err != nil {
		panic(err)
	}

	plist_args, err := module.Load("callable=extended_test.list_args,instance_required",
		[]IDL.MetaFFIType{IDL.HANDLE},
		[]IDL.MetaFFIType{IDL.HANDLE})
	if err != nil {
		panic(err)
	}

	plist_args_without_default, err := module.Load("callable=extended_test.list_args,instance_required",
		[]IDL.MetaFFIType{IDL.HANDLE, IDL.STRING8},
		[]IDL.MetaFFIType{IDL.HANDLE})
	if err != nil {
		panic(err)
	}

	plist_args_without_default_and_varargs, err := module.Load("callable=extended_test.list_args,instance_required,varargs",
		[]IDL.MetaFFIType{IDL.HANDLE, IDL.STRING8, IDL.STRING8_ARRAY},
		[]IDL.MetaFFIType{IDL.HANDLE})
	if err != nil {
		panic(err)
	}

	pdict_args, err := module.Load("callable=extended_test.dict_args,instance_required",
		[]IDL.MetaFFIType{IDL.HANDLE},
		[]IDL.MetaFFIType{IDL.HANDLE})
	if err != nil {
		panic(err)
	}

	pdict_args_without_default, err := module.Load("callable=extended_test.dict_args,instance_required",
		[]IDL.MetaFFIType{IDL.HANDLE, IDL.STRING8},
		[]IDL.MetaFFIType{IDL.HANDLE})
	if err != nil {
		panic(err)
	}

	pdict_args_without_default_and_kwargs, err := module.Load("callable=extended_test.dict_args,instance_required,named_args",
		[]IDL.MetaFFIType{IDL.HANDLE, IDL.STRING8, IDL.HANDLE},
		[]IDL.MetaFFIType{IDL.HANDLE})
	if err != nil {
		panic(err)
	}

	pnamed_only, err := module.Load("callable=extended_test.named_only,instance_required,named_args",
		[]IDL.MetaFFIType{IDL.HANDLE, IDL.HANDLE},
		[]IDL.MetaFFIType{IDL.ANY})
	if err != nil {
		panic(err)
	}

	ppositional_only, err := module.Load("callable=extended_test.positional_only,instance_required",
		[]IDL.MetaFFIType{IDL.HANDLE, IDL.STRING8, IDL.STRING8},
		[]IDL.MetaFFIType{IDL.ANY})
	if err != nil {
		panic(err)
	}

	parg_positional_arg_named, err := module.Load("callable=extended_test.arg_positional_arg_named,instance_required",
		[]IDL.MetaFFIType{IDL.HANDLE},
		[]IDL.MetaFFIType{IDL.HANDLE})
	if err != nil {
		panic(err)
	}

	parg_positional_arg_named_with_kwargs, err := module.Load("callable=extended_test.arg_positional_arg_named,instance_required,named_args",
		[]IDL.MetaFFIType{IDL.HANDLE, IDL.HANDLE},
		[]IDL.MetaFFIType{IDL.HANDLE})
	if err != nil {
		panic(err)
	}

	parg_positional_arg_named_without_default, err := module.Load("callable=extended_test.arg_positional_arg_named,instance_required",
		[]IDL.MetaFFIType{IDL.HANDLE, IDL.STRING8},
		[]IDL.MetaFFIType{IDL.HANDLE})
	if err != nil {
		panic(err)
	}

	parg_positional_arg_named_without_default_with_kwargs, err := module.Load("callable=extended_test.arg_positional_arg_named,instance_required,named_args",
		[]IDL.MetaFFIType{IDL.HANDLE, IDL.STRING8, IDL.HANDLE},
		[]IDL.MetaFFIType{IDL.HANDLE})
	if err != nil {
		panic(err)
	}

	parg_positional_arg_named_without_default_with_varargs, err := module.Load("callable=extended_test.arg_positional_arg_named,instance_required,varargs",
		[]IDL.MetaFFIType{IDL.HANDLE, IDL.STRING8, IDL.STRING8_ARRAY},
		[]IDL.MetaFFIType{IDL.STRING8_ARRAY})
	if err != nil {
		panic(err)
	}

	parg_positional_arg_named_without_default_with_varargs_and_kwargs, err := module.Load("callable=extended_test.arg_positional_arg_named,instance_required,varargs,named_args",
		[]IDL.MetaFFIType{IDL.HANDLE, IDL.STRING8, IDL.STRING8_ARRAY, IDL.HANDLE},
		[]IDL.MetaFFIType{IDL.STRING8_ARRAY})
	if err != nil {
		panic(err)
	}

	// call tests
	_, err = psetX(instance, int64(4))
	if err != nil {
		panic(err)
	}
	val, err := pgetX(instance)
	if err != nil {
		panic(err)
	}
	if val[0].(int64) != 4 {
		panic(fmt.Sprintf("expected 4, received %v", val[0]))
	}

	//--------

	val, err = ppositional_or_named(instance, "PositionalOrNamed")
	if err != nil {
		panic(err)
	}
	if val[0].(string) != "PositionalOrNamed" {
		panic(fmt.Sprintf("expected 'PositionalOrNamed', received '%v'", val[0]))
	}

	pydict, err := NewPyDict()
	if err != nil {
		panic(err)
	}
	err = pydict.Set("value", "PositionalOrNamed")
	if err != nil {
		panic(err)
	}
	val, err = ppositional_or_named_as_named(instance, pydict.Instance)
	if err != nil {
		panic(err)
	}
	if val[0].(string) != "PositionalOrNamed" {
		panic(fmt.Sprintf("expected 'PositionalOrNamed', received '%v'", val[0]))
	}

	//-------

	val, err = plist_args(instance)
	if err != nil {
		panic(err)
	}

	strings := val[0].([]string)

	if strings[0] != "default" {
		panic(fmt.Sprintf("expected 'default', got '%v'", val[0].(string)))
	}

	//---------------

	val, err = plist_args_without_default(instance, "None Default")
	if err != nil {
		panic(err)
	}

	strings = val[0].([]string)

	if strings[0] != "None Default" {
		panic(fmt.Sprintf("expected 'default', got '%v'", val[0].(string)))
	}

	//---------------

	val, err = plist_args_without_default_and_varargs(instance, "None-Default 2", []string{"arg1", "arg2", "arg3"})
	if err != nil {
		panic(err)
	}

	strings = val[0].([]string)

	if strings[0] != "None-Default 2" {
		panic(fmt.Sprintf("expected 'default', got '%v'", strings[0]))
	}

	if strings[1] != "arg1" {
		panic(fmt.Sprintf("expected 'default', got '%v'", strings[1]))
	}

	if strings[2] != "arg2" {
		panic(fmt.Sprintf("expected 'default', got '%v'", strings[2]))
	}

	if strings[3] != "arg3" {
		panic(fmt.Sprintf("expected 'default', got '%v'", strings[3]))
	}

	//---------------

	val, err = pdict_args(instance) // dict_args()
	if err != nil {
		panic(err)
	}

	strings = val[0].([]string)

	if strings[0] != "default" {
		panic(fmt.Sprintf("expected 'default', got '%v'", strings[0]))
	}

	//-------

	val, err = pdict_args_without_default(instance, "none-default") // dict_args()
	if err != nil {
		panic(err)
	}

	strings = val[0].([]string)

	if strings[0] != "none-default" {
		panic(fmt.Sprintf("expected 'default', got '%v'", strings[0]))
	}

	//-------

	pydict, err = NewPyDict()
	if err != nil {
		panic(err)
	}
	err = pydict.Set("key1", "val1")
	if err != nil {
		panic(err)
	}
	err = pydict.Set("key2", "val2")
	if err != nil {
		panic(err)
	}

	val, err = pdict_args_without_default_and_kwargs(instance, "none-default", pydict.Instance) // dict_args()
	if err != nil {
		panic(err)
	}

	strings = val[0].([]string)

	if strings[0] != "none-default" {
		panic(fmt.Sprintf("expected 'none-default', got '%v'", strings[0]))
	}

	if strings[1] != "key1" {
		panic(fmt.Sprintf("expected 'key1', got '%v'", strings[1]))
	}

	if strings[2] != "val1" {
		panic(fmt.Sprintf("expected 'val1', got '%v'", strings[2]))
	}

	if strings[3] != "key2" {
		panic(fmt.Sprintf("expected 'key2', got '%v'", strings[3]))
	}

	if strings[4] != "val2" {
		panic(fmt.Sprintf("expected 'val2', got '%v'", strings[4]))
	}

	//-------

	pydict, err = NewPyDict()
	if err != nil {
		panic(err)
	}
	err = pydict.Set("named", "test")
	if err != nil {
		panic(err)
	}

	val, err = pnamed_only(instance, pydict.Instance) // named_only()
	if err != nil {
		panic(err)
	}

	if val[0].(string) != "test" {
		panic(fmt.Sprintf("expected 'test', got '%v'", val[0].(string)))
	}

	//----------

	val, err = ppositional_only(instance, "word1", "word2") // positional_only()
	if err != nil {
		panic(err)
	}

	if val[0].(string) != "word1 word2" {
		panic(fmt.Sprintf("expected 'word1 word2', got '%v'", val[0].(string)))
	}

	//----------

	val, err = parg_positional_arg_named(instance) // arg_positional_arg_named()
	if err != nil {
		panic(err)
	}

	strings = val[0].([]string)

	if strings[0] != "default" {
		panic(fmt.Sprintf("expected 'default', got '%v'", strings[0]))
	}

	//-----

	val, err = parg_positional_arg_named_without_default(instance, "positional arg") // arg_positional_arg_named()
	if err != nil {
		panic(err)
	}

	strings = val[0].([]string)

	if strings[0] != "positional arg" {
		panic(fmt.Sprintf("expected 'positional arg', got '%v'", strings[0]))
	}

	//-----

	val, err = parg_positional_arg_named_without_default_with_varargs(instance, "positional arg", []string{"var positional arg"}) // arg_positional_arg_named()
	if err != nil {
		panic(err)
	}

	strings = val[0].([]string)

	if strings[0] != "positional arg" {
		panic(fmt.Sprintf("expected 'positional arg', got '%v'", strings[0]))
	}

	if strings[1] != "var positional arg" {
		panic(fmt.Sprintf("expected 'var positional arg', got '%v'", strings[1]))
	}

	//-----

	pydict, err = NewPyDict()
	if err != nil {
		panic(err)
	}
	err = pydict.Set("key1", "val1")
	if err != nil {
		panic(err)
	}
	val, err = parg_positional_arg_named_without_default_with_kwargs(instance, "positional arg", pydict.Instance) // arg_positional_arg_named()
	if err != nil {
		panic(err)
	}

	strings = val[0].([]string)

	if strings[0] != "positional arg" {
		panic(fmt.Sprintf("expected 'positional arg', got '%v'", strings[0]))
	}

	if strings[1] != "key1" {
		panic(fmt.Sprintf("expected 'key1', got '%v'", strings[1]))
	}

	if strings[2] != "val1" {
		panic(fmt.Sprintf("expected 'val1', got '%v'", strings[2]))
	}

	//-----

	pydict, err = NewPyDict()
	if err != nil {
		panic(err)
	}
	err = pydict.Set("key1", "val1")
	if err != nil {
		panic(err)
	}

	val, err = parg_positional_arg_named_with_kwargs(instance, pydict.Instance) // arg_positional_arg_named()
	if err != nil {
		panic(err)
	}

	strings = val[0].([]string)

	if strings[0] != "default" {
		panic(fmt.Sprintf("expected 'default', got '%v'", strings[0]))
	}

	if strings[1] != "key1" {
		panic(fmt.Sprintf("expected 'key1', got '%v'", strings[1]))
	}

	if strings[2] != "val1" {
		panic(fmt.Sprintf("expected 'val1', got '%v'", strings[2]))
	}

	//-----

	pydict, err = NewPyDict()
	if err != nil {
		panic(err)
	}
	err = pydict.Set("key1", "val1")
	if err != nil {
		panic(err)
	}

	val, err = parg_positional_arg_named_without_default_with_varargs_and_kwargs(instance, "positional arg", []string{"var positional arg"}, pydict.Instance) // arg_positional_arg_named()
	if err != nil {
		panic(err)
	}

	strings = val[0].([]string)

	if strings[0] != "positional arg" {
		panic(fmt.Sprintf("expected 'positional arg', got '%v'", strings[0]))
	}

	if strings[1] != "var positional arg" {
		panic(fmt.Sprintf("expected 'var positional arg', got '%v'", strings[1]))
	}

	if strings[2] != "key1" {
		panic(fmt.Sprintf("expected 'key1', got '%v'", strings[2]))
	}

	if strings[3] != "val1" {
		panic(fmt.Sprintf("expected 'val1', got '%v'", strings[3]))
	}
}
