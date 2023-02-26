#include <runtime/runtime_plugin_api.h>
#include <utils/scope_guard.hpp>
#include <boost/thread.hpp>
#include "functions_repository.h"
#include <utils/xllr_api_wrapper.h>

using namespace metaffi::utils;

#define handle_err(err, err_len, desc) \
	*err_len = strlen( desc ); \
	*err = (char*)malloc(*err_len + 1); \
	strcpy(*err, desc ); \
	memset((*err+*err_len), 0, 1);

#define catch_err(err, err_len, desc) \
catch(std::exception& exc) \
{\
	handle_err(err, err_len, desc);\
}

#define handle_err_str(err, err_len, descstr) \
	*err_len = descstr.length(); \
	*err = (char*)malloc(*err_len + 1); \
	descstr.copy(*err, *err_len, 0); \
	memset((*err+*err_len), 0, 1);


#define TRUE 1
#define FALSE 0

#define GO_RUNTIME "go_runtime"
#define GO_RUNTIME_LENGTH 10

boost::mutex runtime_flags_lock;


//--------------------------------------------------------------------
void load_runtime(char** err, uint32_t* err_len)
{
	// load xllr.go.goruntime OR xllr.go.runtime
	// xllr.go.loader - loads Go runtime
	
	// loads the right dynamic library according to GO_RUNTIME flag in XLLR
	
	try
	{
		xllr_api_wrapper xllr;
		boost::unique_lock<boost::mutex> l(runtime_flags_lock);
		if (!xllr.is_runtime_flag_set(GO_RUNTIME, GO_RUNTIME_LENGTH)) // go is NOT loaded
		{
			xllr.set_runtime_flag(GO_RUNTIME, GO_RUNTIME_LENGTH);
		}
	}
	catch_err(err, err_len, exc.what());
}
//--------------------------------------------------------------------
void free_runtime(char** /*err*/, uint32_t* /*err_len*/){ /* No runtime free */ }
//--------------------------------------------------------------------
void* load_function(const char* module_path, uint32_t module_path_len, const char* function_path, uint32_t function_path_len, int8_t params_count, int8_t retval_count, char** err, uint32_t* err_len)
{
	/*
	 * Load modules into modules repository - make sure every module is loaded once
	 */
	try
	{
		return functions_repository::get_instance().load_function(std::string(module_path, module_path_len), std::string(function_path, function_path_len), params_count, retval_count);
	}
	catch(std::exception& exc)
	{
		handle_err(err, err_len, exc.what());
	}
	
	return nullptr;
}
//--------------------------------------------------------------------
void free_function(void* pff, char** /*err*/, uint32_t* /*err_len*/)
{
	/*
	 * Go doesn't support freeing libraries
	 */
}
//--------------------------------------------------------------------