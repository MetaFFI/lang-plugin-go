#include <runtime/runtime_plugin_api.h>
#include <utils/scope_guard.hpp>
#include <utils/entity_path_parser.h>
#include <runtime_manager/go/runtime_manager.h>
#include <runtime_manager/go/module.h>
#include <runtime_manager/go/entity.h>

#include <memory>
#include <unordered_map>
#include <mutex>
#include <sstream>
#include <cstring>
#include <iostream>
#include <cstdlib>

#include <boost/algorithm/string.hpp>

namespace
{
	bool go_plugin_log_enabled()
	{
		static const bool enabled = []() -> bool
		{
			const char* raw = std::getenv("METAFFI_GO_PLUGIN_DEBUG_LOG");
			if(!raw)
			{
				return false;
			}

			std::string val(raw);
			boost::algorithm::to_lower(val);
			return val == "1" || val == "true" || val == "yes" || val == "on";
		}();
		return enabled;
	}
}

#define GO_PLUGIN_LOG(msg) do { if(go_plugin_log_enabled()) { std::cerr << "[go_plugin] " << msg << std::endl; } } while(0)

using namespace metaffi::utils;

namespace
{
	void set_err(char** err, const char* msg)
	{
		if (err == nullptr) return;
		const std::size_t len = std::strlen(msg);
		char* buf = static_cast<char*>(malloc(len + 1));
		if (buf) {
			std::memcpy(buf, msg, len + 1);
			*err = buf;
		}
	}

	// Keeps Module and Entity alive for the lifetime of the xcall
	struct GoXCallHolder
	{
		std::shared_ptr<Module> module;
		std::shared_ptr<Entity> entity;
	};

	// Single global runtime manager (Go has no external runtime to load)
	go_runtime_manager& get_go_manager()
	{
		static go_runtime_manager manager;
		return manager;
	}

	// Cache loaded modules by path so we don't load the same library multiple times
	std::unordered_map<std::string, std::shared_ptr<Module>>& get_module_cache()
	{
		static std::unordered_map<std::string, std::shared_ptr<Module>> cache;
		return cache;
	}

	static std::mutex g_module_cache_mutex;

	// Build Go guest entrypoint symbol name from entity_path (e.g. "EntryPoint_HelloWorld")
	std::string build_entrypoint_name(const char* entity_path)
	{
		entity_path_parser fpp(entity_path);
		std::stringstream fp;
		fp << "EntryPoint_";

		if (fpp.contains("callable"))
		{
			std::string callable_name = fpp["callable"];
			boost::replace_all(callable_name, ".", "_");
			fp << callable_name;
			if (callable_name.size() >= 12 && callable_name.substr(callable_name.size() - 12) == "_EmptyStruct")
				fp << "_MetaFFI";
		}
		else if (fpp.contains("global"))
		{
			if (fpp.contains("getter")) fp << "Get";
			else if (fpp.contains("setter")) fp << "Set";
			else throw std::runtime_error("global action is not specified (getter/setter)");
			fp << fpp["global"];
		}
		else if (fpp.contains("field"))
		{
			std::string action = fpp.contains("getter") ? "_Get" : fpp.contains("setter") ? "_Set" : "";
			if (action.empty()) throw std::runtime_error("field action is not specified (getter/setter)");
			std::string fieldName = fpp["field"];
			boost::replace_all(fieldName, ".", action);
			fp << fieldName;
		}
		else
			throw std::runtime_error("entity_path must contain callable, global, or field");

		return fp.str();
	}
}

//--------------------------------------------------------------------
void load_runtime(char** err)
{
	// No-op: Go compiles to standalone binaries; there is no external runtime to load.
	(void)err;
}

//--------------------------------------------------------------------
void free_runtime(char** err)
{
	// No-op: Go shared libraries cannot be unloaded (dlclose limitation).
	(void)err;
}

//--------------------------------------------------------------------
xcall* load_entity(const char* module_path, const char* entity_path, metaffi_type_info* params_types, int8_t params_count, metaffi_type_info* retvals_types, int8_t retval_count, char** err)
{
	(void)params_types;
	(void)params_count;
	(void)retvals_types;
	(void)retval_count;

	try
	{
		if (module_path == nullptr || module_path[0] == '\0')
		{
			set_err(err, "module_path cannot be null or empty");
			return nullptr;
		}
		if (entity_path == nullptr || entity_path[0] == '\0')
		{
			set_err(err, "entity_path cannot be null or empty");
			return nullptr;
		}

		GO_PLUGIN_LOG("load_entity: module_path=" << module_path << " entity_path=" << entity_path);

		go_runtime_manager& manager = get_go_manager();
		std::shared_ptr<Module> module;

		{
			std::lock_guard<std::mutex> lock(g_module_cache_mutex);
			auto& cache = get_module_cache();
			auto it = cache.find(module_path);
			if (it != cache.end())
			{
				module = it->second;
				GO_PLUGIN_LOG("load_entity: module from cache");
			}
			else
			{
				GO_PLUGIN_LOG("load_entity: loading module...");
				manager.load_runtime();
				module = manager.load_module(module_path);
				cache[module_path] = module;
				GO_PLUGIN_LOG("load_entity: module loaded");
			}
		}

		// Module::load_entity expects entity_path (e.g. "callable=CallTransformer") and maps it to symbol EntryPoint_* internally
		GO_PLUGIN_LOG("load_entity: calling module->load_entity(entity_path)");
		std::shared_ptr<Entity> entity = module->load_entity(entity_path);
		GO_PLUGIN_LOG("load_entity: got entity");
		void* func_ptr = entity->get_function_pointer();
		GO_PLUGIN_LOG("load_entity: func_ptr=" << func_ptr);
		if (func_ptr == nullptr)
		{
			set_err(err, "load_entity: function pointer is null");
			return nullptr;
		}

		GoXCallHolder* holder = new GoXCallHolder{std::move(module), std::move(entity)};
		xcall* pxcall = new xcall(func_ptr, holder);
		GO_PLUGIN_LOG("load_entity: created xcall, returning");
		return pxcall;
	}
	catch (const std::exception& e)
	{
		set_err(err, e.what());
		return nullptr;
	}
}

//--------------------------------------------------------------------
// Context struct matching Go's go_callable_context (GoCallable.go)
#pragma pack(push, 1)
struct go_callable_context
{
	unsigned long long func_handle;
	int8_t params_count;
	int8_t retval_count;
};
#pragma pack(pop)

xcall* make_callable(void* make_callable_context, metaffi_type_info* params_types, int8_t params_count, metaffi_type_info* retvals_types, int8_t retval_count, char** err)
{
	try
	{
		// make_callable_context is a Go handle (metaffi_handle, stored in the Go handle table).
		// We need a loaded Go module to look up the dispatcher symbols.
		std::shared_ptr<Module> module;
		{
			std::lock_guard<std::mutex> lock(g_module_cache_mutex);
			auto& cache = get_module_cache();
			if (cache.empty())
			{
				set_err(err, "make_callable: no Go modules loaded");
				return nullptr;
			}
			// All Go guest modules link the same SDK, so any module has the dispatchers.
			module = cache.begin()->second;
		}

		bool has_params = params_count > 0;
		bool has_retvals = retval_count > 0;

		const char* dispatcher_name =
			has_params && has_retvals   ? "GoCallable_ParamsRet"
			: !has_params && has_retvals ? "GoCallable_NoParamsRet"
			: has_params && !has_retvals ? "GoCallable_ParamsNoRet"
			                             : "GoCallable_NoParamsNoRet";

		void* dispatcher = module->get_symbol(dispatcher_name);
		if (!dispatcher)
		{
			std::string msg = std::string("make_callable: dispatcher '") + dispatcher_name + "' not found in Go module";
			set_err(err, msg.c_str());
			return nullptr;
		}

		// Allocate a context matching Go's go_callable_context
		auto* ctx = static_cast<go_callable_context*>(malloc(sizeof(go_callable_context)));
		ctx->func_handle = reinterpret_cast<unsigned long long>(make_callable_context);
		ctx->params_count = params_count;
		ctx->retval_count = retval_count;

		xcall* pxcall = new xcall(dispatcher, ctx);
		GO_PLUGIN_LOG("make_callable: created xcall dispatcher=" << dispatcher_name << " handle=" << ctx->func_handle);
		return pxcall;
	}
	catch (const std::exception& e)
	{
		set_err(err, e.what());
		return nullptr;
	}
}

//--------------------------------------------------------------------
void free_xcall(xcall* pxcall, char** err)
{
	if (pxcall == nullptr) return;
	try
	{
		void* context = pxcall->pxcall_and_context[1];
		if (context != nullptr)
		{
			GoXCallHolder* holder = static_cast<GoXCallHolder*>(context);
			delete holder;
		}
		delete pxcall;
	}
	catch (const std::exception& e)
	{
		set_err(err, e.what());
	}
}

//--------------------------------------------------------------------
// XLLR calls these with the xcall* as context; we forward to the function pointer stored in the xcall.
void xcall_params_ret(void* context, cdts params_ret[2], char** out_err)
{
	GO_PLUGIN_LOG("xcall_params_ret: entry context=" << context);
	if (context == nullptr) { set_err(out_err, "xcall_params_ret: context is null"); return; }
	xcall* pxcall = static_cast<xcall*>(context);
	GO_PLUGIN_LOG("xcall_params_ret: pxcall=" << static_cast<void*>(pxcall) << " func=" << (pxcall ? pxcall->pxcall_and_context[0] : nullptr));
	(*pxcall)(params_ret, out_err);
	GO_PLUGIN_LOG("xcall_params_ret: exit");
}

void xcall_params_no_ret(void* context, cdts params_ret[2], char** out_err)
{
	GO_PLUGIN_LOG("xcall_params_no_ret: entry context=" << context);
	if (context == nullptr) { set_err(out_err, "xcall_params_no_ret: context is null"); return; }
	xcall* pxcall = static_cast<xcall*>(context);
	// Convention: params_ret[0]=params, params_ret[1]=unused; pass params slot to Go handler
	(*pxcall)(&params_ret[0], out_err);
	GO_PLUGIN_LOG("xcall_params_no_ret: exit");
}

void xcall_no_params_ret(void* context, cdts params_ret[2], char** out_err)
{
	GO_PLUGIN_LOG("xcall_no_params_ret: entry context=" << context);
	if (context == nullptr) { set_err(out_err, "xcall_no_params_ret: context is null"); return; }
	xcall* pxcall = static_cast<xcall*>(context);
	GO_PLUGIN_LOG("xcall_no_params_ret: pxcall=" << static_cast<void*>(pxcall) << " func=" << (pxcall ? pxcall->pxcall_and_context[0] : nullptr));
	// Convention: params_ret[0]=params (unused), params_ret[1]=retvals; pass retvals slot to Go handler
	(*pxcall)(&params_ret[1], out_err);
	GO_PLUGIN_LOG("xcall_no_params_ret: exit");
}

void xcall_no_params_no_ret(void* context, char** out_err)
{
	GO_PLUGIN_LOG("xcall_no_params_no_ret: entry context=" << context);
	if (context == nullptr) { set_err(out_err, "xcall_no_params_no_ret: context is null"); return; }
	xcall* pxcall = static_cast<xcall*>(context);
	(*pxcall)(out_err);
	GO_PLUGIN_LOG("xcall_no_params_no_ret: exit");
}
