#include <doctest/doctest.h>

#include "go_test_env.h"
#include "go_wrappers.h"

#include <memory>
#include <string>
#include <vector>

TEST_CASE("core - HelloWorld")
{
	auto& env = go_test_env();
	auto hello = env.guest_module.load_entity("callable=HelloWorld", {}, { metaffi_string8_type });
	auto [msg] = hello.call<std::string>();
	CHECK(msg == "Hello World, from Go");
}

TEST_CASE("core - DivIntegers")
{
	auto& env = go_test_env();
	auto div = env.guest_module.load_entity(
		"callable=DivIntegers",
		{ metaffi_int64_type, metaffi_int64_type },
		{ metaffi_float64_type });
	auto [val] = div.call<double>(int64_t(10), int64_t(2));
	CHECK(val == doctest::Approx(5.0));
}

TEST_CASE("core - JoinStrings")
{
	auto& env = go_test_env();
	auto join = env.guest_module.load_entity(
		"callable=JoinStrings",
		{ metaffi_string8_array_type },
		{ metaffi_string8_type });
	std::vector<std::string> parts = { "a", "b", "c" };
	auto [joined] = join.call<std::string>(parts);
	CHECK(joined == "a,b,c");
}

TEST_CASE("core - WaitABit")
{
	auto& env = go_test_env();
	auto wait = env.guest_module.load_entity(
		"callable=WaitABit",
		{ metaffi_int64_type },
		{});
	CHECK_NOTHROW(wait.call<>(int64_t(1)));
}

TEST_CASE("core - ReturnNull")
{
	auto& env = go_test_env();
	auto ret_null = env.guest_module.load_entity("callable=ReturnNull", {}, { metaffi_null_type });
	auto [null_val] = ret_null.call<metaffi_variant>();
	CHECK(std::holds_alternative<cdt_metaffi_handle>(null_val));
	CHECK(std::get<cdt_metaffi_handle>(null_val).handle == nullptr);
}

TEST_CASE("core - ReturnMultipleReturnValues")
{
	auto& env = go_test_env();
	auto ret = env.guest_module.load_entity(
		"callable=ReturnMultipleReturnValues",
		{},
		{ metaffi_int64_type, metaffi_string8_type, metaffi_float64_type, metaffi_any_type, metaffi_uint8_array_type, metaffi_handle_type });
	auto [i, s, f, a, b, h] = ret.call<int64_t, std::string, double, metaffi_variant, std::vector<uint8_t>, cdt_metaffi_handle*>();
	CHECK(i == 1);
	CHECK(s == "string");
	CHECK(f == doctest::Approx(3.0));
	CHECK(b.size() == 3u);
	CHECK(h != nullptr);
}

TEST_CASE("core - ReturnAny")
{
	auto& env = go_test_env();
	auto ret = env.guest_module.load_entity(
		"callable=ReturnAny",
		{ metaffi_int64_type },
		{ metaffi_any_type });
	auto [v0] = ret.call<metaffi_variant>(int64_t(0));
	CHECK(std::holds_alternative<metaffi_int64>(v0));
	CHECK(std::get<metaffi_int64>(v0) == 1);
}

TEST_CASE("core - AcceptsAny")
{
	auto& env = go_test_env();
	auto acc = env.guest_module.load_entity(
		"callable=AcceptsAny",
		{ metaffi_int64_type, metaffi_any_type },
		{});
	CHECK_NOTHROW(acc.call<>(int64_t(0), int64_t(1)));
}

TEST_CASE("core - CallCallbackAdd")
{
	auto& env = go_test_env();
	// CallCallbackAdd(add func(int64,int64) int64) -> (int64, error).
	// Full test requires make_callable + call_xcall; for now verify entity loads.
	CHECK_NOTHROW(env.guest_module.load_entity(
		"callable=CallCallbackAdd",
		{ metaffi_callable_type },
		{ metaffi_int64_type }));
}

TEST_CASE("core - ReturnCallbackAdd")
{
	auto& env = go_test_env();
	auto ret_cb = env.guest_module.load_entity(
		"callable=ReturnCallbackAdd",
		{},
		{ metaffi_callable_type });
	auto [cb] = ret_cb.call<cdt_metaffi_callable*>();
	// Callable returned from Go – verify non-null and has a valid xcall
	CHECK(cb != nullptr);
	CHECK(cb->val != nullptr);
}
