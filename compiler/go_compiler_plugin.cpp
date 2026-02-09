#include "go_compiler_plugin.h"

#include <compiler/go/cpp_wrapper/go_compiler_wrapper.h>
#include <runtime/xllr_capi_loader.h>
#include <utils/env_utils.h>  // get_env_var

#include <cstring>
#include <stdexcept>

#ifdef _WIN32
	#define GO_COMPILER_EXE_SUFFIX ".exe"
#else
	#define GO_COMPILER_EXE_SUFFIX ""
#endif

static std::string resolve_go_compiler_path()
{
	std::string home = get_env_var("METAFFI_HOME");
	if (home.empty())
		throw std::runtime_error("METAFFI_HOME is not set");
	// lang-plugin-go copies go_compiler to METAFFI_HOME/go
	return home + "/go/go_compiler" + GO_COMPILER_EXE_SUFFIX;
}

GoCompilerPlugin::GoCompilerPlugin()
	: m_wrapper(std::make_unique<metaffi::compiler::go::GoCompilerWrapper>(resolve_go_compiler_path()))
{
}

GoCompilerPlugin::~GoCompilerPlugin() = default;

void GoCompilerPlugin::ensure_init()
{
	if (m_initialized)
		return;
	m_initialized = true;
}

void GoCompilerPlugin::init()
{
	ensure_init();
}

static void set_err(char** out_err, uint32_t* out_err_len, const std::string& msg)
{
	size_t len = msg.size();
	*out_err = xllr_alloc_string(msg.c_str(), static_cast<uint64_t>(len));
	*out_err_len = static_cast<uint32_t>(len);
}

void GoCompilerPlugin::compile_to_guest(const char* idl_def_json, uint32_t idl_def_json_length,
	const char* output_path, uint32_t output_path_length,
	const char* guest_options, uint32_t guest_options_length,
	char** out_err, uint32_t* out_err_len)
{
	ensure_init();
	try
	{
		std::string idl_json(idl_def_json, idl_def_json_length);
		std::string out_path(output_path, output_path_length);
		std::string opts(guest_options, guest_options_length);
		m_wrapper->compile_to_guest(idl_json, out_path, opts);
		*out_err = nullptr;
		*out_err_len = 0;
	}
	catch (const std::exception& e)
	{
		set_err(out_err, out_err_len, e.what());
	}
	catch (...)
	{
		set_err(out_err, out_err_len, "Unknown error in compile_to_guest");
	}
}

void GoCompilerPlugin::compile_from_host(const char* idl_def_json, uint32_t idl_def_json_length,
	const char* output_path, uint32_t output_path_length,
	const char* host_options, uint32_t host_options_length,
	char** out_err, uint32_t* out_err_len)
{
	ensure_init();
	try
	{
		std::string idl_json(idl_def_json, idl_def_json_length);
		std::string out_path(output_path, output_path_length);
		std::string opts(host_options, host_options_length);
		m_wrapper->compile_from_host(idl_json, out_path, opts);
		*out_err = nullptr;
		*out_err_len = 0;
	}
	catch (const std::exception& e)
	{
		set_err(out_err, out_err_len, e.what());
	}
	catch (...)
	{
		set_err(out_err, out_err_len, "Unknown error in compile_from_host");
	}
}
