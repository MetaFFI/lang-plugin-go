#pragma once
#include <memory>
#include <unordered_map>
#include <utils/foreign_function.h>
#include <sstream>

//--------------------------------------------------------------------
class functions_repository
{
private: // variable
	static std::unique_ptr<functions_repository> instance;

private: // methods
	std::unordered_map<std::string, std::shared_ptr<boost::dll::shared_library>> modules;
	std::vector<std::shared_ptr<foreign_function_params_ret_entrypoint>> functions_params_ret;
	std::vector<std::shared_ptr<foreign_function_params_no_ret_entrypoint>> functions_params_no_ret;
	std::vector<std::shared_ptr<foreign_function_no_params_ret_entrypoint>> functions_no_params_ret;
	std::vector<std::shared_ptr<foreign_function_no_params_no_ret_entrypoint>> functions_no_params_no_ret;
	
public: // static functions
	static functions_repository& get_instance();
	static void free_instance();
	
public: // methods
	functions_repository() = default;
	~functions_repository() = default;
	
	void* load_function(const std::string& module_path, const std::string& entrypoint_name, int params_count, int retval_count);
	
};
//--------------------------------------------------------------------