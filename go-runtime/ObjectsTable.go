package openffi

/*
typedef void* openffi_handle;

openffi_handle int_to_handle(int i)
{
	return (openffi_handle)i;
}
*/
import "C"

import "sync"

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
func SetObject(obj interface{}) C.openffi_handle{
	
	lock.Lock()
	defer lock.Unlock()

	if h, found := objectsToHandles[obj]; found{
		return h
	}

	handleID := C.int_to_handle(len(handlesToObjects)+1)

	handlesToObjects[handleID] = obj
	objectsToHandles[obj] = handleID

	return handleID
}


func GetObject(h C.openffi_handle) interface{}{

	lock.RLock()
	defer lock.RUnlock()

	if o, found := handlesToObjects[h]; found{
		return o
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