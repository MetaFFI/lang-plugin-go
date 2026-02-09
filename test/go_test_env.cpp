#include "go_test_env.h"

#include <utils/env_utils.h>

#include <stdexcept>

#if defined(_WIN32) || defined(_WIN64)
#define GO_GUEST_LIB_SUFFIX ".dll"
#else
#define GO_GUEST_LIB_SUFFIX ".so"
#endif

namespace
{
std::string require_env(const char* name)
{
	std::string value = get_env_var(name);
	if (value.empty())
	{
		std::string msg = "Environment variable not set: ";
		msg += name;
		throw std::runtime_error(msg);
	}
	return value;
}

std::string resolve_guest_module_path()
{
	std::string root = require_env("METAFFI_SOURCE_ROOT");
	std::filesystem::path path(root);
	path /= "sdk";
	path /= "test_modules";
	path /= "guest_modules";
	path /= "go";
	path /= "test_bin";
	path /= std::string("guest_MetaFFIGuest") + GO_GUEST_LIB_SUFFIX;

	std::string path_str = path.string();
	if (!std::filesystem::exists(path_str))
	{
		throw std::runtime_error("Go guest module not found: " + path_str +
			". Build the guest into test_bin/ first (e.g. run MetaFFI Go compiler on sdk/test_modules/guest_modules/go; output is guest_MetaFFIGuest.dll/.so).");
	}
	return path_str;
}
} // namespace

GoTestEnv::GoTestEnv()
	: runtime("go")
	, guest_module(runtime.runtime_plugin(), resolve_guest_module_path())
{
	runtime.load_runtime_plugin();
}

GoTestEnv::~GoTestEnv()
{
	runtime.release_runtime_plugin();
}

GoTestEnv& go_test_env()
{
	static GoTestEnv env;
	return env;
}
