#include "objects_table.h"
#include <shared_mutex>

std::shared_mutex m;

//--------------------------------------------------------------------
metaffi_handle set_object(void* obj)
{
	std::unique_lock<std::shared_mutex> l(m);
	
	auto it = objects_to_handles.find(obj);
	if(it != objects_to_handles.end())
	{
		return it->second;
	}

	metaffi_handle id = (metaffi_handle)(objects.size()+1);
	objects[id] = obj;
	objects_to_handles[obj] = id;
	
	return id;
}
//--------------------------------------------------------------------
void* get_object(metaffi_handle handle)
{
	std::shared_lock<std::shared_mutex> l(m);
	
	auto it = objects.find(handle);
	if(it == objects.end())
	{
		return nullptr;
	}
	
	return it->second;
}
//--------------------------------------------------------------------
void remove_object(metaffi_handle handle)
{
	std::unique_lock<std::shared_mutex> l(m);
	
	auto it = objects.find(handle);
	if(it == objects.end())
	{
		return;
	}
	
	objects_to_handles.erase(it->second);
	objects.erase(it);
}
//--------------------------------------------------------------------
int contains_object(void* obj)
{
	return objects_to_handles.find(obj) != objects_to_handles.end();
}
//--------------------------------------------------------------------
