#include <doctest/doctest.h>

#include "go_test_env.h"

#include <vector>

// GetThreeBuffers, ExpectThreeBuffers, GetSomeClasses, ExpectThreeSomeClasses,
// Make2DArray, Make3DArray, MakeRaggedArray, Sum3DArray, SumRaggedArray
// TODO: align with guest_modules/go/arrays.go entity names and implement.

TEST_CASE("arrays - placeholder")
{
	auto& env = go_test_env();
	REQUIRE(env.guest_module.module_path().size() > 0u);
}
