package main

import (
	"github.com/MetaFFI/lang-plugin-go/api"
	metaffi "github.com/MetaFFI/lang-plugin-go/go-runtime"
	"github.com/MetaFFI/plugin-sdk/compiler/go/IDL"
)

var runtime *api.MetaFFIRuntime

type Response struct {
	instance metaffi.MetaFFIHandle
	getText  func(...interface{}) ([]interface{}, error)
}

func NewResponse(h metaffi.MetaFFIHandle) (*Response, error) {
	this := &Response{instance: h}

	mod, err := runtime.LoadModule("requests.Response")
	if err != nil {
		return nil, err
	}

	this.getText, err = mod.Load("callable=Response.text.fget,instance_required", []IDL.MetaFFIType{IDL.HANDLE}, []IDL.MetaFFIType{IDL.STRING8})
	if err != nil {
		return nil, err
	}

	return this, nil
}

func (this *Response) GetText() (string, error) {
	r, err := this.getText(this.instance)
	if err != nil {
		return "", err
	}

	return r[0].(string), nil
}

//--------------------------------------------------------------------

type Requests struct {
	get func(...interface{}) ([]interface{}, error)
}

func NewRequests() (*Requests, error) {
	module, err := runtime.LoadModule("requests")
	if err != nil {
		return nil, err
	}

	this := &Requests{}

	this.get, err = module.Load("callable=get", []IDL.MetaFFIType{IDL.STRING8}, []IDL.MetaFFIType{IDL.HANDLE})
	if err != nil {
		return nil, err
	}

	return this, nil
}

func (this *Requests) Get(url string) (*Response, error) {
	r, err := this.get(url)
	if err != nil {
		return nil, err
	}

	return NewResponse(r[0].(metaffi.MetaFFIHandle))
}

//--------------------------------------------------------------------

type Tag struct {
	instance metaffi.MetaFFIHandle
	get      func(...interface{}) ([]interface{}, error)
}

func NewTag(h metaffi.MetaFFIHandle) (*Tag, error) {
	tag := &Tag{instance: h}

	tagModule, err := runtime.LoadModule("bs4.element.Tag")
	if err != nil {
		return nil, err
	}

	tag.get, err = tagModule.Load("callable=Tag.get", []IDL.MetaFFIType{IDL.HANDLE, IDL.STRING8}, []IDL.MetaFFIType{IDL.STRING8})
	if err != nil {
		return nil, err
	}

	return tag, nil
}

func (this *Tag) Get(attribute string) (string, error) {
	res, err := this.get(this.instance, attribute)
	if err != nil {
		return "", err
	}

	return res[0].(string), nil
}

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

func (this *PyList) Len() (interface{}, error) {
	res, err := this.len(this.instance)
	if err != nil {
		return nil, err
	}

	return res[0].(int64), nil
}

//--------------------------------------------------------------------

type BeautifulSoup struct {
	pfindAll func(...interface{}) ([]interface{}, error)
	instance metaffi.MetaFFIHandle
}

func NewBeautifulSoup(source string, parser string) (*BeautifulSoup, error) {
	mod, err := runtime.LoadModule("bs4")
	if err != nil {
		return nil, err
	}

	res := &BeautifulSoup{}

	constructor, err := mod.Load("callable=BeautifulSoup", []IDL.MetaFFIType{IDL.STRING8, IDL.STRING8}, []IDL.MetaFFIType{IDL.HANDLE})
	if err != nil {
		return nil, err
	}

	res.pfindAll, err = mod.Load("callable=BeautifulSoup.find_all", []IDL.MetaFFIType{IDL.HANDLE, IDL.STRING8}, []IDL.MetaFFIType{IDL.HANDLE})
	if err != nil {
		return nil, err
	}

	instance, err := constructor(source, parser)
	if err != nil {
		return nil, err
	}
	res.instance = instance[0].(metaffi.MetaFFIHandle)
	return res, nil
}

func (this *BeautifulSoup) FindAll(tag string) ([]*Tag, error) {
	listHandle, err := this.pfindAll(this.instance, tag)
	if err != nil {
		return nil, err
	}

	pylist, err := NewPyListFromHandle(listHandle[0].(metaffi.MetaFFIHandle))
	if err != nil {
		return nil, err
	}

	length, err := pylist.Len()
	if err != nil {
		return nil, err
	}

	tags := make([]*Tag, 0)
	var i int64 = 0
	for ; i < length.(int64); i++ {
		tagHandle, err := pylist.Get(i)
		if err != nil {
			return nil, err
		}

		tag, err := NewTag(tagHandle.(metaffi.MetaFFIHandle))
		if err != nil {
			return nil, err
		}

		tags = append(tags, tag)
	}

	return tags, nil
}

//--------------------------------------------------------------------

func main() {
	runtime = api.NewMetaFFIRuntime("python311")
	err := runtime.LoadRuntimePlugin()
	if err != nil {
		panic(err)
	}

	req, err := NewRequests()
	if err != nil {
		panic(err)
	}

	res, err := req.Get("https://microsoft.com/")
	if err != nil {
		panic(err)
	}

	html, err := res.GetText()
	if err != nil {
		panic(err)
	}

	bs, err := NewBeautifulSoup(html, "html.parser")
	if err != nil {
		panic(err)
	}

	links, err := bs.FindAll("a")
	if err != nil {
		panic(err)
	}

	for _, tag := range links {
		l, err := tag.Get("href")
		if err != nil {
			panic(err)
		}

		println(l)
	}

}
