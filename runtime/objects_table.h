#pragma once
#include <runtime/openffi_primitives.h>
#include <map>

extern "C"
{
	openffi_handle set_object(void* obj);
	void* get_object(openffi_handle handle);
	void remove_object(openffi_handle handle);
	int contains_object(void* obj);
}

static std::map<openffi_handle, void*> objects;
static std::map<void*, openffi_handle> objects_to_handles;

