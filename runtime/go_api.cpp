#include <runtime/runtime_plugin_api.h>
#include <utils/scope_guard.hpp>
#include <boost/thread.hpp>
#include "functions_repository.h"
#include <utils/entity_path_parser.h>

using namespace metaffi::utils;

#define handle_err(err, desc) \
	{auto err_len = strlen( desc ); \
	*err = (char*)malloc(err_len + 1); \
	std::copy(desc, desc + err_len, *err);   \
    (*err)[err_len] = '\0';}

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
void load_runtime(char** err)
{
	// go runtime loads when loading the module
}
//--------------------------------------------------------------------
void free_runtime(char** /*err*/){ /* No runtime free */ }
//--------------------------------------------------------------------
xcall* load_entity(const char* module_path, const char* entity_path, metaffi_type_info* params_types, int8_t params_count, metaffi_type_info* retvals_types, int8_t retval_count, char** err)
{
	/*
	 * Load modules into modules repository - make sure every module is loaded once
	 */
	try
	{
		// build from function path the correct entrypoint
		metaffi::utils::entity_path_parser fpp(entity_path);
		
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
		
		void* pfunc = functions_repository::get_instance().load_function(module_path, fp.str(), params_count, retval_count);
		xcall* pxcall = new xcall(pfunc, nullptr);
		
		return pxcall;
	}
	catch(std::exception& exc)
	{
		handle_err(err, exc.what());
	}
	
	return nullptr;
}
//--------------------------------------------------------------------
xcall* make_callable(void* make_callable_context, metaffi_type_info* params_types, int8_t params_count, metaffi_type_info* retvals_types, int8_t retval_count, char** err)
{
	return nullptr;
}
//--------------------------------------------------------------------
void free_xcall(xcall* pxcall, char** /*err*/)
{
	delete pxcall;
	pxcall = nullptr;
}
//--------------------------------------------------------------------
