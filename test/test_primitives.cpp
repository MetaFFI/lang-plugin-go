#include <doctest/doctest.h>

#include "go_test_env.h"

// AcceptsPrimitives, EchoBytes, ToUpperRune (from primitives.go)

TEST_CASE("primitives - placeholder")
{
	auto& env = go_test_env();
	REQUIRE(env.guest_module.module_path().size() > 0u);
}
