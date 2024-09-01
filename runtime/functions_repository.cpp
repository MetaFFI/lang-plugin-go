#include "functions_repository.h"
#include <boost/filesystem.hpp>
#include <utils/function_loader.hpp>
#include <utils/entity_path_parser.h>
#include <utils/library_loader.h>
#include <sstream>

std::unique_ptr<functions_repository> functions_repository::instance;


//--------------------------------------------------------------------
functions_repository& functions_repository::get_instance()
{
	// IMPORTANT! This singleton is not thread-safe! If this needs to be thread-safe, use std::once.
	
	if(!functions_repository::instance)
	{
		functions_repository::instance = std::make_unique<functions_repository>();
	}
	
	return *functions_repository::instance;
}
//--------------------------------------------------------------------
void functions_repository::free_instance()
{
	functions_repository::instance = nullptr;
}
//--------------------------------------------------------------------
void* functions_repository::load_function(const std::string& module_path, const std::string& entrypoint_name, int params_count, int retval_count)
{
	if(module_path.empty()){
		throw std::runtime_error("Guest library is not defined");
	}
	
	auto it = this->modules.find(module_path);
	
	std::shared_ptr<boost::dll::shared_library> metaffi_guest_lib;
	if(it == this->modules.end())
	{
		if(!boost::filesystem::exists(module_path))
		{
			throw std::invalid_argument("given module path is not found");
		}
		
		// if module not found - load it
		
		std::shared_ptr<boost::dll::shared_library> mod = std::make_shared<boost::dll::shared_library>();
		mod->load(module_path);
		this->modules[module_path] = mod;
		metaffi_guest_lib = mod;
	}
	else
	{
		metaffi_guest_lib = it->second;
	}
	
	if(params_count > 0 && retval_count > 0)
	{
		return (void*)metaffi_guest_lib->get<foreign_function_entrypoint_signature_params_ret>(entrypoint_name);
	}
	else if(params_count > 0)
	{
		return (void*)metaffi_guest_lib->get<foreign_function_entrypoint_signature_params_no_ret>(entrypoint_name);
	}
	else if(retval_count > 0)
	{
		return (void*)metaffi_guest_lib->get<foreign_function_entrypoint_signature_no_params_ret>(entrypoint_name);
	}
	else
	{
		return (void*)metaffi_guest_lib->get<foreign_function_entrypoint_signature_no_params_no_ret>(entrypoint_name);
	}
}
//--------------------------------------------------------------------

