#include <doctest/doctest.h>

#include "go_test_env.h"

// CallTransformer, ReturnTransformer, CallFunction, CallTransformerWithError, CallTransformerRecover
// Full tests require make_callable + call_xcall; placeholders verify entities exist.

TEST_CASE("callbacks - load CallTransformer")
{
	auto& env = go_test_env();
	CHECK_NOTHROW(env.guest_module.load_entity("callable=CallTransformer", { metaffi_callable_type, metaffi_string8_type }, { metaffi_string8_type }));
}

TEST_CASE("callbacks - load ReturnTransformer")
{
	auto& env = go_test_env();
	CHECK_NOTHROW(env.guest_module.load_entity("callable=ReturnTransformer", {}, { metaffi_callable_type }));
}
