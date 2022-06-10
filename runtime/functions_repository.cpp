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
int64_t functions_repository::load_function(const std::string& function_path)
{
	metaffi::utils::function_path_parser fp(function_path);
	
	std::string metaffi_guest_lib_name = fp[function_path_entry_metaffi_guest_lib];
	
	if(metaffi_guest_lib_name.empty()){
		throw std::runtime_error("Guest library is not defined");
	}
	
	
	auto it = this->modules.find(metaffi_guest_lib_name);
	
	std::shared_ptr<boost::dll::shared_library> metaffi_guest_lib;
	if(it == this->modules.end())
	{
		// if module not found - load it
		std::shared_ptr<boost::dll::shared_library> mod = metaffi::utils::load_library(metaffi_guest_lib_name);
		this->modules[metaffi_guest_lib_name] = mod;
		metaffi_guest_lib = mod;
	}
	else
	{
		metaffi_guest_lib = it->second;
	}
	
	// load function (from guest module)
	auto foreign_function = metaffi::utils::load_func<foreign_function_entrypoint_signature>(*metaffi_guest_lib, fp[function_path_entry_entrypoint_function]);
	
	int64_t function_id = (int64_t)this->functions.size();
	this->functions.push_back(foreign_function);
	
	return function_id;
}
//--------------------------------------------------------------------
std::shared_ptr<foreign_function_entrypoint> functions_repository::get_function(int64_t function_id)
{
	if(function_id < 0 || function_id > this->functions.size()-1)
	{
		std::stringstream ss;
		ss << "invalid function id " << function_id;
		throw std::runtime_error(ss.str().c_str());
	}
	
	return this->functions[function_id];
}
//--------------------------------------------------------------------
