#include <doctest/doctest.h>

#include "go_test_env.h"

// MakeStringList, MakeStringIntMap, MakeIntSet, MakeNestedMap, MakeSomeClassList, MakeMapAny

TEST_CASE("collections - placeholder")
{
	auto& env = go_test_env();
	REQUIRE(env.guest_module.module_path().size() > 0u);
}
