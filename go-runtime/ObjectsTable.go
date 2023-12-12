package metaffi

/*
typedef void* metaffi_handle;

metaffi_handle int_to_handle(unsigned long long i)
{
	return (metaffi_handle)i;
}
*/
import "C"
import (
	"fmt"
	"sync"
)

const GO_RUNTIME_ID = 3958232544

type Handle C.metaffi_handle
type MetaFFIHandle struct {
	Val       Handle
	RuntimeID uint64
}

var (
	handlesToObjects map[C.metaffi_handle]interface{}
	objectsToHandles map[interface{}]C.metaffi_handle
	lock             sync.RWMutex
)

func init() {
	handlesToObjects = make(map[C.metaffi_handle]interface{})
	objectsToHandles = make(map[interface{}]C.metaffi_handle)
}

// sets the object and returns a handle
// if object already set, it returns the existing handle
func SetObject(obj interface{}) Handle {

	lock.Lock()
	defer lock.Unlock()

	if h, found := objectsToHandles[obj]; found {
		return Handle(h)
	}

	handleID := C.int_to_handle(C.ulonglong(len(handlesToObjects) + 1))

	handlesToObjects[handleID] = obj
	objectsToHandles[obj] = handleID

	return Handle(handleID)
}

func GetObject(h Handle) interface{} {

	lock.RLock()
	defer lock.RUnlock()

	if o, found := handlesToObjects[C.metaffi_handle(h)]; found {

		return o
	} else {
		return nil
	}

}

func ContainsObject(obj interface{}) bool {

	lock.RLock()
	defer lock.RUnlock()

	_, found := objectsToHandles[obj]
	return found

}

func ReleaseObject(h Handle) error {
	lock.Lock()
	defer lock.Unlock()

	obj, found := handlesToObjects[C.metaffi_handle(h)]
	if !found {
		return fmt.Errorf("Given handle (%v) is not found in MetaFFI Go's object table", h)
	}

	objectsToHandles[obj] = nil
	handlesToObjects[C.metaffi_handle(h)] = nil

	return nil
}
