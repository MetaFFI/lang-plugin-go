#include <doctest/doctest.h>

#include "go_test_env.h"

// GetColor, ColorName (enum), sub.Echo

TEST_CASE("enum_sub - placeholder")
{
	auto& env = go_test_env();
	REQUIRE(env.guest_module.module_path().size() > 0u);
}
