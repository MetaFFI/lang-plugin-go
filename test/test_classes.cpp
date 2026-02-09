#include <doctest/doctest.h>

#include "go_test_env.h"

// SomeClass (Print), TestMap (NewTestMap, Set, Get, Contains, GetName/SetName), Base/Derived, Outer/Inner
// TODO: implement when guest exports are confirmed.

TEST_CASE("classes - placeholder")
{
	auto& env = go_test_env();
	REQUIRE(env.guest_module.module_path().size() > 0u);
}
