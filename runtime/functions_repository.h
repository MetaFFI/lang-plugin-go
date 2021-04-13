#pragma once
#include <memory>
#include <unordered_map>
#include <utils/foreign_function.h>

//--------------------------------------------------------------------
class functions_repository
{
private: // variable
	static std::unique_ptr<functions_repository> instance;

private: // methods
	std::unordered_map<std::string, std::shared_ptr<boost::dll::shared_library>> modules;
	std::vector<std::shared_ptr<foreign_function_entrypoint>> functions;
	
public: // static functions
	static functions_repository& get_instance();
	static void free_instance();
	
public: // methods
	functions_repository() = default;
	~functions_repository() = default;
	
	int64_t load_function(const std::string& function_path);
	std::shared_ptr<foreign_function_entrypoint> get_function(int64_t function_id);
};
//--------------------------------------------------------------------