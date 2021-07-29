package openffi

/*
#include <stdint.h>

typedef void* openffi_handle;

openffi_handle int_to_handle(intptr_t i)
{
	return (openffi_handle)i;
}
*/
import "C"
import "sync"

type handle C.openffi_handle

var(
	handlesToObjects map[C.openffi_handle]interface{}
	objectsToHandles map[interface{}]C.openffi_handle
	lock sync.RWMutex
)


func init(){
	handlesToObjects = make(map[C.openffi_handle]interface{})
	objectsToHandles = make(map[interface{}]C.openffi_handle)
}

// sets the object and returns a handle
// if object already set, it returns the existing handle
func SetObject(obj interface{}) handle{
	
	lock.Lock()
	defer lock.Unlock()

	if h, found := objectsToHandles[obj]; found{
		return h
	}

	handleID := C.int_to_handle(len(handlesToObjects)+1)

	handlesToObjects[handleID] = obj
	objectsToHandles[obj] = handleID

	return handle(handleID)
}


func GetObject(h handle) interface{}{

	lock.RLock()
	defer lock.RUnlock()

	if o, found := handlesToObjects[C.openffi_handle(h)]; found{
		return handle(o)
	} else {
		return nil
	}

}

func ContainsObject(obj interface{}) bool{

	lock.RLock()
	defer lock.RUnlock()

	_, found := objectsToHandles[obj]
	return found

}