#include <runtime/runtime_plugin_api.h>
#include <utils/scope_guard.hpp>
#include <boost/filesystem.hpp>
#include "functions_repository.h"

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


//--------------------------------------------------------------------
void load_runtime(char** /*err*/, uint32_t* /*err_len*/){ /* No runtime to load */ }
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
	char** out_err, uint64_t *out_err_len,
	uint64_t args_len,
	va_list params
)
{
	try
	{
		// get module
		std::shared_ptr<foreign_function_entrypoint> func = functions_repository::get_instance().get_function(function_id);
		
		// call function
		(*func)(out_err, out_err_len,
		        args_len,
		        params);
	}
	catch_err((char**)out_err, out_err_len, exc.what());
}
//--------------------------------------------------------------------
