package metaffi

/*
	typedef void* metaffi_handle;
    void Releaser(h metaffi_handle);
*/
import "C"


func GetReleaserCFunction() unsafe.Pointer{
	return unsafe.Pointer(C.Releaser)
}