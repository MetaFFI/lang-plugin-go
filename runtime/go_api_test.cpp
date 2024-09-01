#define DOCTEST_CONFIG_IMPLEMENT_WITH_MAIN
#include "runtime_id.h"
#include <doctest/doctest.h>
#include <filesystem>
#include <runtime/runtime_plugin_api.h>
#include <utils/scope_guard.hpp>

std::string original;
std::filesystem::path module_path;

#define go_xcall_scope_guard(name) \
	metaffi::utils::scope_guard sg_##name([&name](){char* err = nullptr;\
	free_xcall(name, &err);\
	if(err){                          \
            std::string strerr(err);                       \
            xllr_alloc_string(err, strlen(err));                       \
			FAIL(strerr); }})
 
struct GlobalSetup {
	GlobalSetup()
	{
		module_path = std::filesystem::path(__FILE__);
		module_path = module_path.parent_path();
		module_path.append("test");
#ifdef _WIN32
		module_path.append("TestRuntime_MetaFFIGuest.dll");
#else
		module_path.append("TestRuntime_MetaFFIGuest.so");
#endif

		if(std::getenv("METAFFI_HOME") == nullptr)
		{
			std::cerr << "METAFFI_HOME" << " is not set" << std::endl;
			exit(1);
		}
		
		const char* err = load_xllr();
		if(err)
		{
			std::cerr << "failed to load XLLR functions: " << err << std::endl;
			exit(1);
		}
		
	}

	~GlobalSetup() = default;
};

static GlobalSetup setup;

char* err = nullptr;

xcall* cppload_function(const std::string& mod_path,
                        const std::string& entity_path,
                        std::vector<metaffi_type_info> params_types,
                        std::vector<metaffi_type_info> retvals_types)
{
	err = nullptr;
	uint32_t err_len_load = 0;

	metaffi_type_info* params_types_arr = params_types.empty() ? nullptr : params_types.data();
	metaffi_type_info* retvals_types_arr = retvals_types.empty() ? nullptr : retvals_types.data();

	xcall* pxcall = load_entity(mod_path.c_str(),
	                            entity_path.c_str(),
	                            params_types_arr, params_types.size(),
	                            retvals_types_arr, retvals_types.size(),
	                            &err);

	if(err)
	{
		FAIL(std::string(err));
	}
	REQUIRE((err_len_load == 0));
	REQUIRE((pxcall->pxcall_and_context[0] != nullptr));
	REQUIRE((pxcall->pxcall_and_context[1] == nullptr));// no context in Go

	return pxcall;
};


TEST_SUITE("go runtime api")
{
	TEST_CASE("HelloWorld")
	{
		std::string entity_path = "callable=HelloWorld";
		xcall* phello_world = cppload_function(module_path.string(), entity_path, {}, {});
		go_xcall_scope_guard(phello_world);
		(*phello_world)(&err);
		if(err) { FAIL(std::string(err)); }
	}

	TEST_CASE("runtime_test_target.returns_an_error")
	{
		std::string entity_path = "callable=ReturnsAnError";
		xcall* preturns_an_error = cppload_function(module_path.string(), entity_path, {}, {});
		go_xcall_scope_guard(preturns_an_error);
		(*preturns_an_error)(&err);
		REQUIRE((err != nullptr));
		xllr_free_string(err);
		
	}

	TEST_CASE("runtime_test_target.div_integers")
	{
		std::string entity_path = "callable=DivIntegers";
		std::vector<metaffi_type_info> params_types = {metaffi_type_info(metaffi_int64_type),
		                                               metaffi_type_info(metaffi_int64_type)};
		std::vector<metaffi_type_info> retvals_types = {metaffi_type_info(metaffi_float32_type)};

		xcall* pdiv_integers = cppload_function(module_path.string(), entity_path, params_types, retvals_types);
		go_xcall_scope_guard(pdiv_integers);
		
		cdts* pcdts = xllr_alloc_cdts_buffer(params_types.size(), retvals_types.size());
		cdts_scope_guard(pcdts);
		cdts& params = pcdts[0];
		cdts& retval = pcdts[1];

		params[0] = ((metaffi_int64) 10);
		params[1] = ((metaffi_int64) 2);

		(*pdiv_integers)(pcdts, &err);
		if(err) { FAIL(std::string(err)); }


		REQUIRE((retval[0].type == metaffi_float32_type));
		REQUIRE((retval[0].cdt_val.float32_val == 5.0f));
	}

	TEST_CASE("runtime_test_target.join_strings")
	{
		std::string entity_path = "callable=JoinStrings";
		std::vector<metaffi_type_info> params_types = {metaffi_type_info(metaffi_string8_array_type)};
		std::vector<metaffi_type_info> retvals_types = {metaffi_type_info(metaffi_string8_type)};

		xcall* join_strings = cppload_function(module_path.string(), entity_path, params_types, retvals_types);
		go_xcall_scope_guard(join_strings);
		
		cdts* pcdts = xllr_alloc_cdts_buffer(params_types.size(), retvals_types.size());
		cdts_scope_guard(pcdts);
		cdts& pcdts_params = pcdts[0];
		cdts& pcdts_retvals = pcdts[1];
		pcdts_params[0].set_new_array(3, 1, metaffi_string8_type);

		std::u8string one = u8"one";
		std::u8string two = u8"two";
		std::u8string three = u8"three";

		pcdts_params[0].cdt_val.array_val->arr[0].set_string(one.c_str(), false);
		pcdts_params[0].cdt_val.array_val->arr[1].set_string(two.c_str(), false);
		pcdts_params[0].cdt_val.array_val->arr[2].set_string(three.c_str(), false);

		(*join_strings)(pcdts, &err);
		if(err) { FAIL(std::string(err)); }


		REQUIRE((pcdts_retvals[0].type == metaffi_string8_type));
		REQUIRE((std::u8string_view(pcdts_retvals[0].cdt_val.string8_val) == u8"one,two,three"));
	}

	TEST_CASE("runtime_test_target.SomeClass")
	{
		std::string entity_path = "callable=GetSomeClasses";
		std::vector<metaffi_type_info> retvals_getSomeClasses_types = {{metaffi_handle_array_type, (char*) "SomeClass[]", true, 1}};
		xcall* pgetSomeClasses = cppload_function(module_path.string(), entity_path, {}, retvals_getSomeClasses_types);
		go_xcall_scope_guard(pgetSomeClasses);
		
		entity_path = "callable=ExpectThreeSomeClasses";
		std::vector<metaffi_type_info> params_expectThreeSomeClasses_types = {{metaffi_handle_array_type, (char*) "SomeClass[]", true, 1}};
		xcall* pexpectThreeSomeClasses = cppload_function(module_path.string(), entity_path, params_expectThreeSomeClasses_types, {});
		go_xcall_scope_guard(pexpectThreeSomeClasses);

		entity_path = "callable=SomeClass.Print";
		std::vector<metaffi_type_info> params_SomeClassPrint_types = {metaffi_type_info{metaffi_handle_type}};

		xcall* pSomeClassPrint = cppload_function(module_path.string(), entity_path, params_SomeClassPrint_types, {});
		go_xcall_scope_guard(pSomeClassPrint);
		
		cdts* pcdts = xllr_alloc_cdts_buffer(0, 1);
		cdts_scope_guard(pcdts);
		cdts& pcdts_params = pcdts[0];
		cdts& pcdts_retvals = pcdts[1];

		(*pgetSomeClasses)(pcdts, &err);
		if(err) { FAIL(std::string(err)); }
		
		REQUIRE((pcdts_retvals[0].type == metaffi_handle_array_type));
		REQUIRE((pcdts_retvals[0].cdt_val.array_val->fixed_dimensions == 1));
		REQUIRE((pcdts_retvals[0].cdt_val.array_val->length == 3));

		REQUIRE((pcdts_retvals[0].cdt_val.array_val->arr[0].type == metaffi_handle_type));
		REQUIRE((pcdts_retvals[0].cdt_val.array_val->arr[1].type == metaffi_handle_type));
		REQUIRE((pcdts_retvals[0].cdt_val.array_val->arr[2].type == metaffi_handle_type));

		REQUIRE((pcdts_retvals[0].cdt_val.array_val->arr[0].cdt_val.handle_val->runtime_id == GO_RUNTIME_ID));
		REQUIRE((pcdts_retvals[0].cdt_val.array_val->arr[1].cdt_val.handle_val->runtime_id == GO_RUNTIME_ID));
		REQUIRE((pcdts_retvals[0].cdt_val.array_val->arr[2].cdt_val.handle_val->runtime_id == GO_RUNTIME_ID));

		REQUIRE((pcdts_retvals[0].cdt_val.array_val->arr[0].cdt_val.handle_val->handle != nullptr));
		REQUIRE((pcdts_retvals[0].cdt_val.array_val->arr[1].cdt_val.handle_val->handle != nullptr));
		REQUIRE((pcdts_retvals[0].cdt_val.array_val->arr[2].cdt_val.handle_val->handle != nullptr));

		std::vector<cdt_metaffi_handle*> arr = {pcdts_retvals[0].cdt_val.array_val->arr[0].cdt_val.handle_val,
		                                       pcdts_retvals[0].cdt_val.array_val->arr[1].cdt_val.handle_val,
		                                       pcdts_retvals[0].cdt_val.array_val->arr[2].cdt_val.handle_val};
		//--------------------------------------------------------------------

		cdts* pcdts2 = (cdts*)xllr_alloc_cdts_buffer(1, 0);
		cdts_scope_guard(pcdts2);
		cdts& pcdts_params2 = pcdts2[0];
		cdts& pcdts_retvals2 = pcdts2[1];

		pcdts_params2[0].set_new_array(3, 1, metaffi_handle_type);
		pcdts_params2[0].cdt_val.array_val->arr[0].set_handle(arr[0]);
		pcdts_params2[0].cdt_val.array_val->arr[1].set_handle(arr[1]);
		pcdts_params2[0].cdt_val.array_val->arr[2].set_handle(arr[2]);

		(*pexpectThreeSomeClasses)(pcdts2, &err);
		if(err) { FAIL(std::string(err)); }

		//--------------------------------------------------------------------

		cdts* pcdts3 = (cdts*) xllr_alloc_cdts_buffer(1, 0);
		cdts_scope_guard(pcdts3);
		cdts& pcdts_params3 = pcdts3[0];
		cdts& pcdts_retvals3 = pcdts3[1];

		pcdts_params3[0].set_handle(arr[1]);// use the 2nd instance

		(*pSomeClassPrint)(pcdts3, &err);
		if(err) { FAIL(std::string(err)); }
	}


	TEST_CASE("runtime_test_target.ThreeBuffers")
	{
		std::string entity_path = "callable=ExpectThreeBuffers";
		std::vector<metaffi_type_info> params_expectThreeBuffers_types = {{metaffi_uint8_array_type, nullptr, false, 2}};

		xcall* pexpectThreeBuffers = cppload_function(module_path.string(), entity_path, params_expectThreeBuffers_types, {});
		go_xcall_scope_guard(pexpectThreeBuffers);
		
		entity_path = "callable=GetThreeBuffers";
		std::vector<metaffi_type_info> retval_getThreeBuffers_types = {{metaffi_uint8_array_type, nullptr, false, 2}};

		xcall* pgetThreeBuffers = cppload_function(module_path.string(), entity_path, {}, retval_getThreeBuffers_types);
		go_xcall_scope_guard(pgetThreeBuffers);
		
		// pass 3 buffers
		cdts* pcdts = xllr_alloc_cdts_buffer(1, 0);
		cdts_scope_guard(pcdts);
		cdts& pcdts_params = pcdts[0];
		cdts& pcdts_retvals = pcdts[1];

		pcdts_params[0].set_new_array(3, 2, metaffi_uint8_array_type);
		metaffi_uint8 data[3][3] = {{0, 1, 2}, {3, 4, 5}, {6, 7, 8}};
		pcdts_params[0].cdt_val.array_val->arr[0].set_new_array(3, 1, metaffi_uint8_type);
		pcdts_params[0].cdt_val.array_val->arr[1].set_new_array(3, 1, metaffi_uint8_type);
		pcdts_params[0].cdt_val.array_val->arr[2].set_new_array(3, 1, metaffi_uint8_type);

		pcdts_params[0].cdt_val.array_val->arr[0].cdt_val.array_val->arr[0] = (data[0][0]);
		pcdts_params[0].cdt_val.array_val->arr[0].cdt_val.array_val->arr[1] = (data[0][1]);
		pcdts_params[0].cdt_val.array_val->arr[0].cdt_val.array_val->arr[2] = (data[0][2]);

		pcdts_params[0].cdt_val.array_val->arr[1].cdt_val.array_val->arr[0] = (data[1][0]);
		pcdts_params[0].cdt_val.array_val->arr[1].cdt_val.array_val->arr[1] = (data[1][1]);
		pcdts_params[0].cdt_val.array_val->arr[1].cdt_val.array_val->arr[2] = (data[1][2]);

		pcdts_params[0].cdt_val.array_val->arr[2].cdt_val.array_val->arr[0] = (data[2][0]);
		pcdts_params[0].cdt_val.array_val->arr[2].cdt_val.array_val->arr[1] = (data[2][1]);
		pcdts_params[0].cdt_val.array_val->arr[2].cdt_val.array_val->arr[2] = (data[2][2]);

		(*pexpectThreeBuffers)(pcdts, &err);
		if(err) { FAIL(std::string(err)); }


		// get 3 buffers
		cdts* pcdts2 = (cdts*)xllr_alloc_cdts_buffer(0, 1);
		cdts_scope_guard(pcdts2);
		
		cdts& pcdts_params2 = pcdts[0];
		cdts& pcdts_retvals2 = pcdts[1];

		(*pgetThreeBuffers)(pcdts2, &err);
		if(err) { FAIL(std::string(err)); }


		REQUIRE((pcdts_retvals2[0].type == metaffi_uint8_array_type));
		REQUIRE((pcdts_retvals2[0].cdt_val.array_val->fixed_dimensions == 2));
		REQUIRE((pcdts_retvals2[0].cdt_val.array_val->length == 3));
		for(int i = 0; i < 3; i++)
		{
			REQUIRE((pcdts_retvals2[0].cdt_val.array_val->arr[i].cdt_val.array_val->length == 3));
			for(int j = 0; j < 3; j++)
			{
				REQUIRE((pcdts_retvals2[0].cdt_val.array_val->arr[i].cdt_val.array_val->arr[j].cdt_val.uint8_val == j + 1));
			}
		}
	}

	TEST_CASE("runtime_test_target.testmap.set_get_contains")
	{
		// create new testmap
		std::string entity_path = "callable=NewTestMap";
		std::vector<metaffi_type_info> retvals_types = {metaffi_type_info(metaffi_handle_type)};

		xcall* pnew_testmap = cppload_function(module_path.string(), entity_path, {}, retvals_types);
		go_xcall_scope_guard(pnew_testmap);
		cdts* pcdts = xllr_alloc_cdts_buffer(0, 1);
		cdts_scope_guard(pcdts);
		cdts& params_cdts = pcdts[0];
		cdts& retvals_cdts = pcdts[1];

		(*pnew_testmap)((cdts*) pcdts, &err);
		if(err) { FAIL(std::string(err)); }


		REQUIRE((retvals_cdts[0].type == metaffi_handle_type));
		REQUIRE((retvals_cdts[0].cdt_val.handle_val->handle != nullptr));
		REQUIRE((retvals_cdts[0].cdt_val.handle_val->runtime_id == GO_RUNTIME_ID));

		cdt_metaffi_handle* testmap_instance = retvals_cdts[0].cdt_val.handle_val;

		// set
		entity_path = "callable=TestMap.Set,instance_required";
		std::vector<metaffi_type_info> params_types = {metaffi_type_info(metaffi_handle_type),
		                                               metaffi_type_info(metaffi_string8_type),
		                                               metaffi_type_info(metaffi_any_type)};

		xcall* p_testmap_set = cppload_function(module_path.string(), entity_path, {}, retvals_types);
		go_xcall_scope_guard(p_testmap_set);
		cdts* pcdts2 = (cdts*) xllr_alloc_cdts_buffer(3, 0);
		cdts_scope_guard(pcdts2);
		cdts& params_cdts2 = pcdts2[0];
		cdts& retvals_cdts2 = pcdts2[1];

		params_cdts2[0].set_handle(testmap_instance);
		params_cdts2[1].set_string((metaffi_string8)u8"key", false);
		params_cdts2[2] = ((int32_t) 42);

		(*p_testmap_set)(pcdts2, &err);
		if(err) { FAIL(std::string(err)); }


		// contains
		entity_path = "callable=TestMap.Contains,instance_required";
		params_types[0].type = metaffi_handle_type;
		params_types[1].type = metaffi_string8_type;
		retvals_types[0].type = metaffi_bool_type;

		xcall* p_testmap_contains = cppload_function(module_path.string(), entity_path, {}, retvals_types);
		go_xcall_scope_guard(p_testmap_contains);
		cdts* pcdts3 = (cdts*) xllr_alloc_cdts_buffer(2, 1);
		cdts_scope_guard(pcdts3);
		cdts& params_cdts3 = pcdts3[0];
		cdts& retvals_cdts3 = pcdts3[1];

		params_cdts3[0].set_handle(testmap_instance);
		params_cdts3[1].set_string((metaffi_string8) u8"key", true);

		(*p_testmap_contains)(pcdts3, &err);
		if(err) { FAIL(std::string(err)); }


		REQUIRE((retvals_cdts3[0].type == metaffi_bool_type));
		REQUIRE((retvals_cdts3[0].cdt_val.bool_val != 0));

		// get
		entity_path = "callable=TestMap.Get,instance_required";
		params_types[0].type = metaffi_handle_type;
		params_types[1].type = metaffi_string8_type;
		retvals_types[0].type = metaffi_any_type;

		xcall* p_testmap_get = cppload_function(module_path.string(), entity_path, {}, retvals_types);
		go_xcall_scope_guard(p_testmap_get);
		
		cdts* pcdts4 = (cdts*) xllr_alloc_cdts_buffer(2, 1);
		cdts_scope_guard(pcdts4);
		cdts& params_cdts4 = pcdts4[0];
		cdts& retvals_cdts4 = pcdts4[1];

		params_cdts4[0].set_handle(testmap_instance);
		params_cdts4[1].set_string((char8_t*) u8"key", true);

		(*p_testmap_get)((cdts*) pcdts4, &err);
		if(err) { FAIL(std::string(err)); }


		REQUIRE((retvals_cdts4[0].type == metaffi_int32_type));
		REQUIRE((retvals_cdts4[0].cdt_val.int32_val == 42));
	}

	TEST_CASE("runtime_test_target.testmap.set_get_contains_cpp_object")
	{
		// create new testmap
		std::string entity_path = "callable=NewTestMap";
		std::vector<metaffi_type_info> retvals_types = {metaffi_type_info(metaffi_handle_type)};

		xcall* pnew_testmap = cppload_function(module_path.string(), entity_path, {}, retvals_types);
		go_xcall_scope_guard(pnew_testmap);
		
		cdts* pcdts = xllr_alloc_cdts_buffer(0, 1);
		cdts_scope_guard(pcdts);
		cdts& pcdts_params = pcdts[0];
		cdts& pcdts_retvals = pcdts[1];

		(*pnew_testmap)((cdts*) pcdts, &err);
		if(err) { FAIL(std::string(err)); }


		REQUIRE((pcdts_retvals[0].type == metaffi_handle_type));
		REQUIRE((pcdts_retvals[0].cdt_val.handle_val->handle != nullptr));
		REQUIRE((pcdts_retvals[0].cdt_val.handle_val->runtime_id == GO_RUNTIME_ID));

		cdt_metaffi_handle* testmap_instance = pcdts_retvals[0].cdt_val.handle_val;

		// set
		entity_path = "callable=TestMap.Set,instance_required";
		std::vector<metaffi_type_info> params_types = {metaffi_type_info(metaffi_handle_type),
		                                               metaffi_type_info(metaffi_string8_type),
		                                               metaffi_type_info(metaffi_any_type)};

		xcall* p_testmap_set = cppload_function(module_path.string(), entity_path, {}, retvals_types);
		go_xcall_scope_guard(p_testmap_set);
		
		cdts* pcdts2 = (cdts*) xllr_alloc_cdts_buffer(3, 0);
		cdts_scope_guard(pcdts2);
		cdts& pcdts_params2 = pcdts2[0];
		cdts& pcdts_retvals2 = pcdts2[1];

		std::vector<int> vec_to_insert = {1, 2, 3};

		pcdts_params2[0].set_handle(testmap_instance);
		pcdts_params2[1].set_string((metaffi_string8) u8"key", true);
		pcdts_params2[2].set_handle(new cdt_metaffi_handle{&vec_to_insert, 733, nullptr});

		(*p_testmap_set)((cdts*) pcdts2, &err);
		if(err) { FAIL(std::string(err)); }


		// contains
		entity_path = "callable=TestMap.Contains,instance_required";
		params_types[0].type = metaffi_handle_type;
		params_types[1].type = metaffi_string8_type;
		retvals_types[0].type = metaffi_bool_type;

		xcall* p_testmap_contains = cppload_function(module_path.string(), entity_path, {}, retvals_types);
		go_xcall_scope_guard(p_testmap_contains);
		
		cdts* pcdts3 = (cdts*) xllr_alloc_cdts_buffer(2, 1);
		cdts_scope_guard(pcdts3);
		cdts& pcdts_params3 = pcdts3[0];
		cdts& pcdts_retvals3 = pcdts3[1];

		pcdts_params3[0].set_handle(testmap_instance);
		pcdts_params3[1].set_string((metaffi_string8) u8"key", true);

		(*p_testmap_contains)(pcdts3, &err);
		if(err) { FAIL(std::string(err)); }


		REQUIRE((pcdts_retvals3[0].type == metaffi_bool_type));
		REQUIRE((pcdts_retvals3[0].cdt_val.bool_val != 0));

		// get
		entity_path = "callable=TestMap.Get,instance_required";
		params_types[0].type = metaffi_handle_type;
		params_types[1].type = metaffi_string8_type;
		retvals_types[0].type = metaffi_any_type;

		xcall* p_testmap_get = cppload_function(module_path.string(), entity_path, {}, retvals_types);
		go_xcall_scope_guard(p_testmap_get);
		
		cdts* pcdts4 = (cdts*) xllr_alloc_cdts_buffer(2, 1);
		cdts_scope_guard(pcdts4);
		cdts& pcdts_params4 = pcdts4[0];
		cdts& pcdts_retvals4 = pcdts4[1];

		pcdts_params4[0].set_handle(testmap_instance);
		pcdts_params4[1].set_string((char8_t*) u8"key", true);

		(*p_testmap_get)(pcdts4, &err);
		if(err) { FAIL(std::string(err)); }


		REQUIRE((pcdts_retvals4[0].type == metaffi_handle_type));
		auto& vector_pulled = *(std::vector<int>*) pcdts_retvals4[0].cdt_val.handle_val->handle;

		REQUIRE((vector_pulled[0] == 1));
		REQUIRE((vector_pulled[1] == 2));
		REQUIRE((vector_pulled[2] == 3));
	}

	TEST_CASE("runtime_test_target.testmap.get_set_name")
	{
		// load constructor
		std::string entity_path = "callable=NewTestMap";
		std::vector<metaffi_type_info> retvals_types = {metaffi_type_info(metaffi_handle_type)};

		xcall* pnew_testmap = cppload_function(module_path.string(), entity_path, {}, retvals_types);
		go_xcall_scope_guard(pnew_testmap);
		// load getter
		entity_path = "field=TestMap.Name,instance_required,getter";
		std::vector<metaffi_type_info> params_types = {metaffi_type_info(metaffi_handle_type)};
		retvals_types = {metaffi_type_info(metaffi_string8_type)};

		xcall* pget_name = cppload_function(module_path.string(), entity_path, {}, retvals_types);
		go_xcall_scope_guard(pget_name);
		
		// load setter
		entity_path = "field=TestMap.Name,instance_required,setter";
		params_types[0].type = metaffi_handle_type;
		retvals_types[0].type = metaffi_string8_type;

		xcall* pset_name = cppload_function(module_path.string(), entity_path, {}, retvals_types);
		go_xcall_scope_guard(pset_name);
		
		// create new testmap
		cdts* pcdts = xllr_alloc_cdts_buffer(0, 1);
		cdts_scope_guard(pcdts);
		cdts& pcdts_params = pcdts[0];
		cdts& pcdts_retvals = pcdts[1];

		(*pnew_testmap)(pcdts, &err);
		if(err) { FAIL(std::string(err)); }


		REQUIRE((pcdts_retvals[0].type == metaffi_handle_type));
		REQUIRE((pcdts_retvals[0].cdt_val.handle_val->handle != nullptr));
		REQUIRE((pcdts_retvals[0].cdt_val.handle_val->runtime_id == GO_RUNTIME_ID));

		cdt_metaffi_handle* testmap_instance = pcdts_retvals[0].cdt_val.handle_val;


		// get name
		cdts* pcdts2 = (cdts*) xllr_alloc_cdts_buffer(1, 1);
		cdts_scope_guard(pcdts2);
		cdts& pcdts_params2 = pcdts2[0];
		cdts& pcdts_retvals2 = pcdts2[1];

		pcdts_params2[0].set_handle(testmap_instance);

		(*pget_name)(pcdts2, &err);
		if(err) { FAIL(std::string(err)); }


		REQUIRE((pcdts_retvals2[0].type == metaffi_string8_type));
		REQUIRE((std::u8string(pcdts_retvals2[0].cdt_val.string8_val) == u8"TestMap Name"));

		// set name to "name is my name"

		cdts* pcdts3 = (cdts*) xllr_alloc_cdts_buffer(2, 0);
		cdts_scope_guard(pcdts3);
		cdts& pcdts_params3 = pcdts3[0];
		cdts& pcdts_retvals3 = pcdts3[1];

		pcdts_params3[0].set_handle(testmap_instance);
		pcdts_params3[1].set_string((metaffi_string8) u8"name is my name", true);

		(*pset_name)(pcdts3, &err);
		if(err) { FAIL(std::string(err)); }


		// get name again and make sure it is "name is my name"
		cdts* pcdts4 = (cdts*) xllr_alloc_cdts_buffer(1, 1);
		cdts_scope_guard(pcdts4);
		cdts& pcdts_params4 = pcdts4[0];
		cdts& pcdts_retvals4 = pcdts4[1];

		pcdts_params4[0] .set_handle(testmap_instance);

		(*pget_name)(pcdts4, &err);
		if(err) { FAIL(std::string(err)); }
		
		REQUIRE((pcdts_retvals4[0].type == metaffi_string8_type));
		REQUIRE((std::u8string(pcdts_retvals4[0].cdt_val.string8_val) == u8"name is my name"));
	}

	TEST_CASE("runtime_test_target.testmap.get_set_name_from_empty_struct")
	{
		// load constructor
		std::string entity_path = "callable=TestMap.EmptyStruct";
		std::vector<metaffi_type_info> retvals_types = {metaffi_type_info(metaffi_handle_type)};

		xcall* pnew_testmap = cppload_function(module_path.string(), entity_path, {}, retvals_types);
		go_xcall_scope_guard(pnew_testmap);
		
		// load getter
		entity_path = "field=TestMap.Name,instance_required,getter";
		std::vector<metaffi_type_info> params_types = {metaffi_type_info{metaffi_handle_type}};
		retvals_types = {metaffi_type_info{metaffi_string8_type}};

		xcall* pget_name = cppload_function(module_path.string(), entity_path, {}, retvals_types);
		go_xcall_scope_guard(pget_name);
		
		// load setter
		entity_path = "field=TestMap.Name,instance_required,setter";
		params_types = {metaffi_type_info{metaffi_handle_type}, metaffi_type_info{metaffi_string8_type}};

		xcall* pset_name = cppload_function(module_path.string(), entity_path, {}, retvals_types);
		go_xcall_scope_guard(pset_name);
		
		// create new testmap
		cdts* pcdts = xllr_alloc_cdts_buffer(0, 1);
		cdts_scope_guard(pcdts);
		cdts& pcdts_params = pcdts[0];
		cdts& pcdts_retvals = pcdts[1];

		(*pnew_testmap)(pcdts, &err);
		if(err) { FAIL(std::string(err)); }
		
		REQUIRE((pcdts_retvals[0].type == metaffi_handle_type));
		REQUIRE((pcdts_retvals[0].cdt_val.handle_val->handle != nullptr));
		REQUIRE((pcdts_retvals[0].cdt_val.handle_val->runtime_id == GO_RUNTIME_ID));

		cdt_metaffi_handle* testmap_instance = pcdts_retvals[0].cdt_val.handle_val;

		// get name
		cdts* pcdts2 = (cdts*) xllr_alloc_cdts_buffer(1, 1);
		cdts_scope_guard(pcdts2);
		cdts& pcdts_params2 = pcdts2[0];
		cdts& pcdts_retvals2 = pcdts2[1];

		pcdts_params2[0].set_handle(testmap_instance);

		(*pget_name)(pcdts2, &err);
		if(err) { FAIL(std::string(err)); }


		REQUIRE((pcdts_retvals2[0].type == metaffi_string8_type));
		REQUIRE((std::u8string(pcdts_retvals2[0].cdt_val.string8_val).empty()));

		// set name to "name is my name"

		cdts* pcdts3 = (cdts*) xllr_alloc_cdts_buffer(2, 0);
		cdts_scope_guard(pcdts3);
		cdts& pcdts_params3 = pcdts3[0];
		cdts& pcdts_retvals3 = pcdts3[1];

		pcdts_params3[0] .set_handle(testmap_instance);
		pcdts_params3[1] .set_string((metaffi_string8) u8"name is my name", true);

		(*pset_name)(pcdts3, &err);
		if(err) { FAIL(std::string(err)); }


		// get name again and make sure it is "name is my name"
		cdts* pcdts4 = (cdts*)xllr_alloc_cdts_buffer(1, 1);
		cdts_scope_guard(pcdts4);
		cdts& pcdts_params4 = pcdts4[0];
		cdts& pcdts_retvals4 = pcdts4[1];

		pcdts_params4[0] .set_handle(testmap_instance);

		(*pget_name)(pcdts4, &err);
		if(err) { FAIL(std::string(err)); }


		REQUIRE((pcdts_retvals4[0].type == metaffi_string8_type));
		REQUIRE((std::u8string(pcdts_retvals4[0].cdt_val.string8_val) == u8"name is my name"));
	}

	TEST_CASE("runtime_test_target.wait_a_bit")
	{
		// get five_seconds global
		std::vector<metaffi_type_info> var_type = {metaffi_type_info(metaffi_int64_type, "time.Duration", true)};
		std::string variable_path = "global=FiveSeconds,getter";
		xcall* pfive_seconds_getter = cppload_function(module_path.string(), variable_path, {}, var_type);
		go_xcall_scope_guard(pfive_seconds_getter);
		
		cdts* pcdts = xllr_alloc_cdts_buffer(0, 1);
		cdts_scope_guard(pcdts);
		cdts& pcdts_params = pcdts[0];
		cdts& pcdts_retvals = pcdts[1];

		(*pfive_seconds_getter)(pcdts, &err);
		if(err) { FAIL(std::string(err)); }


		REQUIRE((pcdts_retvals[0].type == metaffi_int64_type));
		REQUIRE((pcdts_retvals[0].cdt_val.int64_val == 5000000000));

		metaffi_int64 five = pcdts_retvals[0].cdt_val.int64_val;

		// call wait_a_bit
		std::string entity_path = "callable=WaitABit";
		std::vector<metaffi_type_info> params_types = {metaffi_type_info(metaffi_int64_type)};

		xcall* pwait_a_bit = cppload_function(module_path.string(), entity_path, params_types, {});
		go_xcall_scope_guard(pwait_a_bit);
		
		cdts* pcdts2 = (cdts*) xllr_alloc_cdts_buffer(1, 0);
		cdts_scope_guard(pcdts2);
		cdts& pcdts_params2 = pcdts2[0];
		cdts& pcdts_retvals2 = pcdts2[1];

		pcdts_params2[0] = (five);

		(*pwait_a_bit)(pcdts2, &err);
		if(err) { FAIL(std::string(err)); }
	}
}