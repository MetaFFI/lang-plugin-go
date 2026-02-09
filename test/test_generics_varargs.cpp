#include <doctest/doctest.h>

#include "go_test_env.h"

// NewIntBox, NewStringBox, Box Get/Set, Sum (varargs), Join (varargs)

TEST_CASE("generics_varargs - placeholder")
{
	auto& env = go_test_env();
	REQUIRE(env.guest_module.module_path().size() > 0u);
}
