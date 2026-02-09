#include <doctest/doctest.h>

#include "go_test_env.h"

// uuid, x/text, goquery (PuerkitoBio/goquery) - include when guest module has third-party entities built.

TEST_CASE("third_party - placeholder")
{
	auto& env = go_test_env();
	REQUIRE(env.guest_module.module_path().size() > 0u);
}
