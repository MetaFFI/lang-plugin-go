#include <runtime/xllr_capi_loader.h>
#include <idl_compiler/idl_plugin_interface.h>
#include <idl_compiler/go/cpp_wrapper/go_idl_compiler_wrapper.h>
#include <utils/logger.hpp>
#include <string>
#include <fstream>
#include <cstring>

static auto LOG = metaffi::get_logger("go.idl");

class GoIDLPlugin : public idl_plugin_interface
{
	metaffi::idl_compiler::GoIDLCompilerWrapper compiler_;

public:
	void init() override
	{
		// Resolve go_idl_compiler path: METAFFI_HOME/go/go_idl_compiler[.exe] (lang-plugin-go copies it there)
		const char* home = std::getenv("METAFFI_HOME");
		if (home && home[0])
		{
			std::string path = home;
			path += "/go/go_idl_compiler";
#ifdef _WIN32
			path += ".exe";
#endif
			compiler_.set_executable_path(path);
		}
		METAFFI_INFO(LOG, "Go IDL plugin initialized");
	}

	void parse_idl(const char* source_code, uint32_t source_code_length,
		const char* file_or_path, uint32_t file_or_path_length,
		char** out_idl_def_json, uint32_t* out_idl_def_json_length,
		char** out_err, uint32_t* out_err_len) override
	{
		try
		{
			std::string path(file_or_path, file_or_path_length);
			if (source_code && source_code_length > 0)
			{
				// Write source to the given path (or temp) then compile
				std::ofstream f(path);
				if (!f)
				{
					*out_err = xllr_alloc_string("Failed to open path for writing", 31);
					*out_err_len = 31;
					*out_idl_def_json = nullptr;
					*out_idl_def_json_length = 0;
					return;
				}
				f.write(source_code, source_code_length);
				f.close();
			}
			if (path.empty())
			{
				*out_err = xllr_alloc_string("file_or_path is required", 23);
				*out_err_len = 23;
				*out_idl_def_json = nullptr;
				*out_idl_def_json_length = 0;
				return;
			}
			std::string idl_json = compiler_.compile(path);
			*out_idl_def_json = xllr_alloc_string(idl_json.c_str(), static_cast<uint64_t>(idl_json.size()));
			*out_idl_def_json_length = static_cast<uint32_t>(idl_json.size());
			*out_err = nullptr;
			*out_err_len = 0;
		}
		catch (const std::exception& e)
		{
			const char* msg = e.what();
			size_t len = std::strlen(msg);
			*out_err = xllr_alloc_string(msg, static_cast<uint64_t>(len));
			*out_err_len = static_cast<uint32_t>(len);
			*out_idl_def_json = nullptr;
			*out_idl_def_json_length = 0;
		}
		catch (...)
		{
			const char* msg = "Unknown error in parse_idl";
			size_t len = std::strlen(msg);
			*out_err = xllr_alloc_string(msg, static_cast<uint64_t>(len));
			*out_err_len = static_cast<uint32_t>(len);
			*out_idl_def_json = nullptr;
			*out_idl_def_json_length = 0;
		}
	}
};

static GoIDLPlugin s_go_idl_plugin;
static bool s_go_idl_initialized = false;

extern "C"
{
	void init_plugin(void)
	{
		if (!s_go_idl_initialized)
		{
			s_go_idl_plugin.init();
			s_go_idl_initialized = true;
		}
	}

	void parse_idl(const char* source_code, uint32_t source_code_length,
		const char* file_path, uint32_t file_path_length,
		char** out_idl_def_json, uint32_t* out_idl_def_json_length,
		char** out_err, uint32_t* out_err_len)
	{
		if (!s_go_idl_initialized)
		{
			s_go_idl_plugin.init();
			s_go_idl_initialized = true;
		}
		s_go_idl_plugin.parse_idl(source_code, source_code_length,
			file_path, file_path_length,
			out_idl_def_json, out_idl_def_json_length,
			out_err, out_err_len);
	}
}
