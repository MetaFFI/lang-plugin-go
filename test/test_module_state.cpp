#include <doctest/doctest.h>

#include "go_test_env.h"

// FiveSeconds, GetCounter, SetCounter, IncCounter (from state.go)

TEST_CASE("module_state - placeholder")
{
	auto& env = go_test_env();
	REQUIRE(env.guest_module.module_path().size() > 0u);
}
