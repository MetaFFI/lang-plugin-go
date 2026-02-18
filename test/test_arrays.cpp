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

TEST_CASE("packed array: sum 1d int64 array")
{
	auto& env = go_test_env();

	auto sum1d = env.guest_module.load_entity(
		"callable=Sum1dInt64Array",
		{metaffi_int64_packed_array_type},
		{metaffi_int64_type});
	auto [sum_val] = sum1d.call<int64_t>(std::vector<int64_t>({1, 2, 3, 4, 5}));
	CHECK(sum_val == 15);
}

TEST_CASE("packed array: echo 1d int64 array")
{
	auto& env = go_test_env();

	auto echo = env.guest_module.load_entity(
		"callable=Echo1dInt64Array",
		{metaffi_int64_packed_array_type},
		{metaffi_int64_packed_array_type});
	auto [result] = echo.call<std::vector<int64_t>>(std::vector<int64_t>({100, 200, 300}));
	CHECK(result == std::vector<int64_t>({100, 200, 300}));
}

TEST_CASE("packed array: echo 1d float64 array")
{
	auto& env = go_test_env();

	auto echo = env.guest_module.load_entity(
		"callable=Echo1dFloat64Array",
		{metaffi_float64_packed_array_type},
		{metaffi_float64_packed_array_type});
	auto [result] = echo.call<std::vector<double>>(std::vector<double>({1.5, 2.5, 3.5}));
	REQUIRE(result.size() == 3);
	CHECK(result[0] == doctest::Approx(1.5));
	CHECK(result[1] == doctest::Approx(2.5));
	CHECK(result[2] == doctest::Approx(3.5));
}

TEST_CASE("packed array: make 1d int64 array")
{
	auto& env = go_test_env();

	auto make1d = env.guest_module.load_entity(
		"callable=Make1dInt64Array",
		{},
		{metaffi_int64_packed_array_type});
	auto [arr] = make1d.call<std::vector<int64_t>>();
	CHECK(arr == std::vector<int64_t>({10, 20, 30, 40, 50}));
}
