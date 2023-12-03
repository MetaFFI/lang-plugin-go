#include <runtime/runtime_plugin_api.h>
#include <utils/scope_guard.hpp>
#include <boost/thread.hpp>
#include "functions_repository.h"
#include <utils/xllr_api_wrapper.h>
#include <utils/function_path_parser.h>

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
void** load_function(const char* module_path, uint32_t module_path_len, const char* function_path, uint32_t function_path_len, metaffi_types_with_alias_ptr params_types, metaffi_types_with_alias_ptr retvals_types, uint8_t params_count, uint8_t retval_count, char** err, uint32_t* err_len)
{
	/*
	 * Load modules into modules repository - make sure every module is loaded once
	 */
	try
	{
		// build from function path the correct entrypoint
		metaffi::utils::function_path_parser fpp(std::string(function_path, function_path_len));
		
		std::stringstream fp;
		fp << "EntryPoint_";
		
		
		if(fpp.contains("callable"))
		{
			std::string callable_name = fpp["callable"];
			boost::replace_all(callable_name, ".", "_");
			
			fp << callable_name;
			
			if(callable_name.ends_with("_EmptyStruct")){
				fp << "_MetaFFI";
			}
		}
		else if(fpp.contains("global"))
		{
			fp << (fpp.contains("getter") ? "Get" :
				   fpp.contains("setter") ? "Set" :
				   throw std::runtime_error("global action is not specified"));
			
			fp << fpp["global"];
		}
		else if(fpp.contains("field"))
		{
			std::string action = (fpp.contains("getter") ? "_Get" :
							       fpp.contains("setter") ? "_Set" :
							       throw std::runtime_error("global action is not specified"));
			
			std::string fieldName = fpp["field"];
			boost::replace_all(fieldName, ".", action);
			
			fp << fieldName;
		}
		
		void* pfunc = functions_repository::get_instance().load_function(std::string(module_path, module_path_len), fp.str(), params_count, retval_count);
		void** res = (void**)malloc(sizeof(void*)*2);
		res[0] = pfunc;
		res[1] = nullptr; // no context required
		
		return res;
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