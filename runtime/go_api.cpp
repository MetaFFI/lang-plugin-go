#include <runtime/runtime_plugin_api.h>
#include <utils/scope_guard.hpp>
#include <boost/filesystem.hpp>
#include <boost/thread.hpp>
#include "functions_repository.h"
#include <utils/xllr_api_wrapper.h>
#include <utils/library_loader.h>

using namespace openffi::utils;

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
			openffi::utils::load_library("xllr.go.loader"); // the dynamic library loads GO runtime
			xllr.set_runtime_flag(GO_RUNTIME, GO_RUNTIME_LENGTH);
		}
	}
	catch_err(err, err_len, exc.what());
}
//--------------------------------------------------------------------
void free_runtime(char** /*err*/, uint32_t* /*err_len*/){ /* No runtime free */ }
//--------------------------------------------------------------------
int64_t load_function(const char* function_path, uint32_t function_path_len, char** err, uint32_t* err_len)
{
	/*
	 * Load modules into modules repository - make sure every module is loaded once
	 */
	try
	{
		return functions_repository::get_instance().load_function(std::string(function_path, function_path_len));
	}
	catch(std::exception& exc)
	{
		handle_err(err, err_len, exc.what());
	}
	
	return -1;
}
//--------------------------------------------------------------------
void free_function(int64_t function_id, char** /*err*/, uint32_t* /*err_len*/)
{
	/*
	 * Go doesn't support freeing libraries
	 */
}
//--------------------------------------------------------------------
void call(
	int64_t function_id,
	cdt* parameters, uint64_t parameters_len,
	cdt* return_values, uint64_t return_values_len,
	char** out_err, uint64_t* out_err_len
)
{
	try
	{
		// get module
		std::shared_ptr<foreign_function_entrypoint> func = functions_repository::get_instance().get_function(function_id);
		
		// call function
		(*func)(parameters, parameters_len,
		        return_values, return_values_len,
		        out_err, out_err_len);
	}
	catch_err((char**)out_err, out_err_len, exc.what());
}
//--------------------------------------------------------------------
