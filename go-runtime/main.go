package main

import (
	"sync"
	"unsafe"
)

/*
void* int_to_vptr(long long i)
{
	return (void*)i;
}
 */
import "C"

var objects map[unsafe.Pointer]interface{}
var lock sync.RWMutex

func init(){
	objects = make(map[unsafe.Pointer]interface{})
}

func Set(v interface{}) unsafe.Pointer{
	lock.Lock()
	defer lock.Unlock()
	id := C.int_to_vptr(C.longlong(len(objects)))
	objects[id] = v
	return id
}

func Get(p unsafe.Pointer) (interface{}, bool){
	lock.RLock()
	defer lock.RUnlock()

	v, exists := objects[p]
	return v, exists
}