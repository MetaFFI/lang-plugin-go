#pragma once
#include <runtime/metaffi_primitives.h>
#include <map>

extern "C"
{
	metaffi_handle set_object(void* obj);
	void* get_object(metaffi_handle handle);
	void remove_object(metaffi_handle handle);
	int contains_object(void* obj);
}

static std::map<metaffi_handle, void*> objects;
static std::map<void*, metaffi_handle> objects_to_handles;

