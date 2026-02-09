#pragma once

#include <compiler/compiler_plugin_interface.h>

#include <cstdint>
#include <memory>
#include <string>

namespace metaffi { namespace compiler { namespace go {
	class GoCompilerWrapper;
}}}

/**
 * Go Compiler Plugin: implements compiler_plugin_interface by running
 * the sdk/compiler/go executable (go_compiler) via GoCompilerWrapper,
 * same pattern as go_idl_compiler.
 */
class GoCompilerPlugin : public compiler_plugin_interface
{
public:
	GoCompilerPlugin();
	~GoCompilerPlugin();

	void init() override;

	void compile_to_guest(const char* idl_def_json, uint32_t idl_def_json_length,
		const char* output_path, uint32_t output_path_length,
		const char* guest_options, uint32_t guest_options_length,
		char** out_err, uint32_t* out_err_len) override;

	void compile_from_host(const char* idl_def_json, uint32_t idl_def_json_length,
		const char* output_path, uint32_t output_path_length,
		const char* host_options, uint32_t host_options_length,
		char** out_err, uint32_t* out_err_len) override;

private:
	void ensure_init();

	std::unique_ptr<metaffi::compiler::go::GoCompilerWrapper> m_wrapper;
	bool m_initialized = false;
};
