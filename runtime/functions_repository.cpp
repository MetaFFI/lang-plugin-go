#include "functions_repository.h"
#include <boost/filesystem.hpp>
#include <utils/function_loader.hpp>
#include <utils/function_path_parser.h>
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
void* functions_repository::load_function(const std::string& module_path, const std::string& function_path, int params_count, int retval_count)
{
	metaffi::utils::function_path_parser fp(function_path);
	
	if(module_path.empty()){
		throw std::runtime_error("Guest library is not defined");
	}
	
	auto it = this->modules.find(module_path);
	
	std::shared_ptr<boost::dll::shared_library> metaffi_guest_lib;
	if(it == this->modules.end())
	{
		// if module not found - load it
		std::shared_ptr<boost::dll::shared_library> mod = metaffi::utils::load_library(module_path);
		this->modules[module_path] = mod;
		metaffi_guest_lib = mod;
	}
	else
	{
		metaffi_guest_lib = it->second;
	}
	
	if(params_count > 0 && retval_count > 0)
	{
		auto foreign_function = metaffi::utils::load_func<foreign_function_entrypoint_signature_params_ret>(*metaffi_guest_lib, fp[function_path_entry_entrypoint_function]);
		this->functions_params_ret.push_back(foreign_function);
		return (void*)(pforeign_function_entrypoint_signature_params_ret)(*((ppforeign_function_entrypoint_signature_params_ret)foreign_function.get()));
	}
	else if(params_count > 0)
	{
		auto foreign_function = metaffi::utils::load_func<foreign_function_entrypoint_signature_params_no_ret>(*metaffi_guest_lib, fp[function_path_entry_entrypoint_function]);
		this->functions_params_no_ret.push_back(foreign_function);
		return (void*)(pforeign_function_entrypoint_signature_params_no_ret)(*((ppforeign_function_entrypoint_signature_params_no_ret)foreign_function.get()));
	}
	else if(retval_count > 0)
	{
		auto foreign_function = metaffi::utils::load_func<foreign_function_entrypoint_signature_no_params_ret>(*metaffi_guest_lib, fp[function_path_entry_entrypoint_function]);
		this->functions_no_params_ret.push_back(foreign_function);
		return (void*)(pforeign_function_entrypoint_signature_no_params_ret)(*((ppforeign_function_entrypoint_signature_no_params_ret)foreign_function.get()));
	}
	else
	{
		auto foreign_function = metaffi::utils::load_func<foreign_function_entrypoint_signature_no_params_no_ret>(*metaffi_guest_lib, fp[function_path_entry_entrypoint_function]);
		this->functions_no_params_no_ret.push_back(foreign_function);
		return (void*)(pforeign_function_entrypoint_signature_no_params_no_ret)(*((ppforeign_function_entrypoint_signature_no_params_no_ret)foreign_function.get()));
	}
}
//--------------------------------------------------------------------

