#include <doctest/doctest.h>

#include "go_test_env.h"

TEST_CASE("errors - ReturnsAnError")
{
	auto& env = go_test_env();
	auto ret = env.guest_module.load_entity("callable=ReturnsAnError", {}, {});
	CHECK_THROWS(ret.call<>());
}

TEST_CASE("errors - ReturnErrorTuple")
{
	auto& env = go_test_env();
	auto ret = env.guest_module.load_entity(
		"callable=ReturnErrorTuple",
		{ metaffi_bool_type },
		{ metaffi_bool_type });
	auto [ok] = ret.call<bool>(true);
	CHECK(ok == true);
}

TEST_CASE("errors - Panics")
{
	auto& env = go_test_env();
	auto panics = env.guest_module.load_entity("callable=Panics", {}, {});
	CHECK_THROWS(panics.call<>());
}
