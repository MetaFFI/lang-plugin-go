/**
 * C exports for the Go compiler plugin (init_plugin, compile_to_guest, compile_from_host).
 * Required by the MetaFFI compiler plugin loader.
 */

#include "go_compiler_plugin.h"
#include <runtime/xllr_capi_loader.h>
#include <cstring>

namespace
{
	GoCompilerPlugin g_plugin;
	bool g_initialized = false;

	void set_error(char** out_err, uint32_t* out_err_len, const char* msg)
	{
		size_t len = std::strlen(msg);
		*out_err = xllr_alloc_string(msg, static_cast<uint64_t>(len));
		*out_err_len = static_cast<uint32_t>(len);
	}

	void clear_error(char** out_err, uint32_t* out_err_len)
	{
		*out_err = nullptr;
		*out_err_len = 0;
	}
}

extern "C"
{
	void init_plugin()
	{
		try
		{
			if (!g_initialized)
			{
				g_plugin.init();
				g_initialized = true;
			}
		}
		catch (const std::exception& e)
		{
			// init has no out_err; caller may check by calling compile_* and getting an error
			(void)e;
		}
	}

	void compile_to_guest(
		const char* idl_def_json, uint32_t idl_def_json_length,
		const char* output_path, uint32_t output_path_length,
		const char* guest_options, uint32_t guest_options_length,
		char** out_err, uint32_t* out_err_len)
	{
		try
		{
			if (!g_initialized)
			{
				g_plugin.init();
				g_initialized = true;
			}
			g_plugin.compile_to_guest(idl_def_json, idl_def_json_length,
				output_path, output_path_length,
				guest_options, guest_options_length,
				out_err, out_err_len);
		}
		catch (const std::exception& e)
		{
			set_error(out_err, out_err_len, e.what());
		}
		catch (...)
		{
			set_error(out_err, out_err_len, "Unknown error in compile_to_guest");
		}
	}

	void compile_from_host(
		const char* idl_def_json, uint32_t idl_def_json_length,
		const char* output_path, uint32_t output_path_length,
		const char* host_options, uint32_t host_options_length,
		char** out_err, uint32_t* out_err_len)
	{
		try
		{
			if (!g_initialized)
			{
				g_plugin.init();
				g_initialized = true;
			}
			g_plugin.compile_from_host(idl_def_json, idl_def_json_length,
				output_path, output_path_length,
				host_options, host_options_length,
				out_err, out_err_len);
		}
		catch (const std::exception& e)
		{
			set_error(out_err, out_err_len, e.what());
		}
		catch (...)
		{
			set_error(out_err, out_err_len, "Unknown error in compile_from_host");
		}
	}
}
