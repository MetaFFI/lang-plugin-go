#pragma once

#include <metaffi/api/metaffi_api.h>

#include <filesystem>
#include <string>

struct GoTestEnv
{
	metaffi::api::MetaFFIRuntime runtime;
	metaffi::api::MetaFFIModule guest_module;

	GoTestEnv();
	~GoTestEnv();
};

GoTestEnv& go_test_env();
