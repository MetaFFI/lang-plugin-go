#define CATCH_CONFIG_MAIN
#include <catch2/catch.hpp>
#include <runtime/runtime_plugin_api.h>
#include <filesystem>
#include <runtime/cdts_wrapper.h>

TEST_CASE( "go runtime api", "[goruntime]" )
{
	std::filesystem::path module_path(__FILE__);
	module_path = module_path.parent_path();
	module_path.append("test");
	module_path.append("TestRuntime_MetaFFIGuest.dll");
	
	char* err = nullptr;
	uint32_t err_len = 0;
	
	SECTION("Load Go Runtime")
	{
		REQUIRE(std::getenv("METAFFI_HOME") != nullptr);
		
		load_runtime(&err, &err_len);
		
		if(err){
			FAIL(err);
		}
		
		REQUIRE(err_len == 0);
	}

	SECTION("HelloWorld")
	{
		std::string function_path = "callable=HelloWorld";
		void** phello_world = load_function(module_path.string().c_str(), module_path.string().length(),
											function_path.c_str(), function_path.length(),
											nullptr, nullptr,
											  0, 0,
											  &err, &err_len);

		if(err){ FAIL(err); }
		REQUIRE(err_len == 0);
		REQUIRE(phello_world[0] != nullptr);
		REQUIRE(phello_world[1] == nullptr); // no context

		uint64_t long_err_len = 0;
		((void(*)(void*,char**, uint64_t*))phello_world[0])(nullptr, &err, &long_err_len);
		if(err){ FAIL(err); }
		REQUIRE(long_err_len == 0);
	}

	SECTION("runtime_test_target.returns_an_error")
	{
		std::string function_path = "callable=ReturnsAnError";
		void** preturns_an_error = load_function(module_path.string().c_str(), module_path.string().length(),
		                                    function_path.c_str(), function_path.length(),
											nullptr, nullptr,
		                                    0, 0,
		                                    &err, &err_len);

		if(err){ FAIL(err); }
		REQUIRE(err_len == 0);
		REQUIRE(preturns_an_error[0] != nullptr);
		REQUIRE(preturns_an_error[1] == nullptr);

		uint64_t long_err_len = 0;
		((void(*)(void*,char**, uint64_t*))preturns_an_error[0])(nullptr, &err, &long_err_len);
		REQUIRE(err != nullptr);
		REQUIRE(long_err_len > 0);

		SUCCEED(std::string(err, long_err_len).c_str());
	}

	SECTION("runtime_test_target.div_integers")
	{
		std::string function_path = "callable=DivIntegers";
		metaffi_type_with_alias params_types[] = {metaffi::runtime::make_type_with_alias(metaffi_int64_type), metaffi::runtime::make_type_with_alias(metaffi_int64_type)};
		metaffi_type_with_alias retvals_types[] = {metaffi::runtime::make_type_with_alias(metaffi_float32_type)};
		
		uint8_t params_count = 2;
		uint8_t retvals_count = 1;
		
		void** pdiv_integers = load_function(module_path.string().c_str(), module_path.string().length(),
		                                         function_path.c_str(), function_path.length(),
		                                        params_types, retvals_types,
		                                     params_count, retvals_count,
		                                         &err, &err_len);

		if(err){ FAIL(err); }
		REQUIRE(err_len == 0);
		REQUIRE(pdiv_integers[0] != nullptr);
		REQUIRE(pdiv_integers[1] == nullptr);
		
		
		cdts* cdts_param_ret = (cdts*)xllr_alloc_cdts_buffer(params_count, retvals_count);
		metaffi::runtime::cdts_wrapper wrapper(cdts_param_ret[0].pcdt, cdts_param_ret[0].len, false);
		wrapper[0]->type = metaffi_int64_type;
		wrapper[0]->cdt_val.metaffi_int64_val.val = 10;
		wrapper[1]->type = metaffi_int64_type;
		wrapper[1]->cdt_val.metaffi_int64_val.val = 2;
		
		uint64_t long_err_len = 0;
		((void(*)(void*,cdts*,char**,uint64_t*))pdiv_integers[0])(pdiv_integers[1], (cdts*)cdts_param_ret, &err, &long_err_len);
		if(err){ FAIL(err); }
		REQUIRE(long_err_len == 0);
		
		metaffi::runtime::cdts_wrapper wrapper_ret(cdts_param_ret[1].pcdt, cdts_param_ret[1].len, false);
		REQUIRE(wrapper_ret[0]->type == metaffi_float32_type);
		REQUIRE(wrapper_ret[0]->cdt_val.metaffi_float32_val.val == 5.0);
		
		if(cdts_param_ret[0].len + cdts_param_ret[1].len > cdts_cache_size){
			free(cdts_param_ret);
		}
	}
	
	SECTION("runtime_test_target.join_strings")
	{
		std::string function_path = "callable=JoinStrings";
		metaffi_type_with_alias params_types[] = {metaffi::runtime::make_type_with_alias(metaffi_string8_array_type)};
		metaffi_type_with_alias retvals_types[] = {metaffi::runtime::make_type_with_alias(metaffi_string8_type)};
		
		void** join_strings = load_function(module_path.string().c_str(), module_path.string().length(),
		                                     function_path.c_str(), function_path.length(),
		                                     params_types, retvals_types,
		                                     1, 1,
		                                     &err, &err_len);
		
		if(err){ FAIL(err); }
		REQUIRE(err_len == 0);
		REQUIRE(join_strings[0] != nullptr);
		REQUIRE(join_strings[1] == nullptr);
		
		
		cdts* cdts_param_ret = (cdts*)xllr_alloc_cdts_buffer(1, 1);
		metaffi::runtime::cdts_wrapper wrapper(cdts_param_ret[0].pcdt, cdts_param_ret[0].len, false);
		metaffi_size array_dimensions;
		metaffi_size array_length[] = {3};
		metaffi_string8 values[] = {(metaffi_string8)"one", (metaffi_string8)"two", (metaffi_string8)"three"};
		metaffi_size vals_length[] = {strlen("one"), strlen("two"), strlen("three")};
		wrapper[0]->type = metaffi_string8_array_type;
		wrapper[0]->cdt_val.metaffi_string8_array_val.dimensions = 1;
		wrapper[0]->cdt_val.metaffi_string8_array_val.dimensions_lengths = (metaffi_size*)array_length;
		wrapper[0]->cdt_val.metaffi_string8_array_val.vals = values;
		wrapper[0]->cdt_val.metaffi_string8_array_val.vals_sizes = vals_length;
		
		uint64_t long_err_len = 0;
		((void(*)(void*,cdts*,char**,uint64_t*))join_strings[0])(join_strings[1], (cdts*)cdts_param_ret, &err, &long_err_len);
		if(err){ FAIL(err); }
		REQUIRE(long_err_len == 0);
		
		metaffi::runtime::cdts_wrapper wrapper_ret(cdts_param_ret[1].pcdt, cdts_param_ret[1].len, false);
		REQUIRE(wrapper_ret[0]->type == metaffi_string8_type);
		
		std::string returned(wrapper_ret[0]->cdt_val.metaffi_string8_val.val, wrapper_ret[0]->cdt_val.metaffi_string8_val.length);
		REQUIRE(returned == "one,two,three");
		
		if(cdts_param_ret[0].len + cdts_param_ret[1].len > cdts_cache_size){
			free(cdts_param_ret);
		}
	}
	
	SECTION("runtime_test_target.testmap.set_get_contains")
	{
		// create new testmap
		std::string function_path = "callable=NewTestMap";
		metaffi_type_with_alias retvals_types[] = {metaffi::runtime::make_type_with_alias(metaffi_handle_type)};
		
		void** pnew_testmap = load_function(module_path.string().c_str(), module_path.string().length(),
		                                    function_path.c_str(), function_path.length(),
		                                    nullptr, retvals_types,
		                                    0, 1,
		                                    &err, &err_len);
		
		if(err){ FAIL(err); }
		REQUIRE(err_len == 0);
		REQUIRE(pnew_testmap[0] != nullptr);
		REQUIRE(pnew_testmap[1] == nullptr);
		
		cdts* cdts_param_ret = (cdts*)xllr_alloc_cdts_buffer(0, 1);
		
		uint64_t long_err_len = 0;
		((void(*)(void*,cdts*,char**,uint64_t*))pnew_testmap[0])(pnew_testmap[1], (cdts*)cdts_param_ret, &err, &long_err_len);
		if(err){ FAIL(err); }
		REQUIRE(long_err_len == 0);
		
		metaffi::runtime::cdts_wrapper wrapper_ret(cdts_param_ret[1].pcdt, cdts_param_ret[1].len, false);
		REQUIRE(wrapper_ret[0]->type == metaffi_handle_type);
		REQUIRE(wrapper_ret[0]->cdt_val.metaffi_handle_val.val != nullptr);
		
		if(cdts_param_ret[0].len + cdts_param_ret[1].len > cdts_cache_size){
			free(cdts_param_ret);
		}
		
		metaffi_handle testmap_instance = wrapper_ret[0]->cdt_val.metaffi_handle_val.val;
		
		// set
		function_path = "callable=TestMap.Set,instance_required";
		metaffi_type_with_alias params_types[] = {metaffi::runtime::make_type_with_alias(metaffi_handle_type),
												  metaffi::runtime::make_type_with_alias(metaffi_string8_type),
												  metaffi::runtime::make_type_with_alias(metaffi_any_type)};
		
		void** p_testmap_set = load_function(module_path.string().c_str(), module_path.string().length(),
		                                    function_path.c_str(), function_path.length(),
		                                    nullptr, retvals_types,
		                                    3, 0,
		                                    &err, &err_len);
		
		if(err){ FAIL(err); }
		REQUIRE(err_len == 0);
		REQUIRE(p_testmap_set[0] != nullptr);
		REQUIRE(p_testmap_set[1] == nullptr);
		
		cdts_param_ret = (cdts*)xllr_alloc_cdts_buffer(3, 0);
		metaffi::runtime::cdts_wrapper wrapper(cdts_param_ret[0].pcdt, cdts_param_ret[0].len, false);
		wrapper[0]->type = metaffi_handle_type;
		wrapper[0]->cdt_val.metaffi_handle_val.val = testmap_instance;
		wrapper[1]->type = metaffi_string8_type;
		wrapper[1]->cdt_val.metaffi_string8_val.val = (char*)"key";
		wrapper[1]->cdt_val.metaffi_string8_val.length = strlen("key");
		wrapper[2]->type = metaffi_int64_type;
		wrapper[2]->cdt_val.metaffi_int64_val.val = 42;
		
		long_err_len = 0;
		((void(*)(void*,cdts*,char**,uint64_t*))p_testmap_set[0])(p_testmap_set[1], (cdts*)cdts_param_ret, &err, &long_err_len);
		if(err){ FAIL(err); }
		REQUIRE(long_err_len == 0);
		
		if(cdts_param_ret[0].len + cdts_param_ret[1].len > cdts_cache_size){
			free(cdts_param_ret);
		}
		
		// contains
		function_path = "callable=TestMap.Contains,instance_required";
		params_types[0].type = metaffi_handle_type;
		params_types[1].type = metaffi_string8_type;
		retvals_types[0].type = metaffi_bool_type;
		
		void** p_testmap_contains = load_function(module_path.string().c_str(), module_path.string().length(),
		                                     function_path.c_str(), function_path.length(),
		                                     nullptr, retvals_types,
		                                     2, 1,
		                                     &err, &err_len);
		
		if(err){ FAIL(err); }
		REQUIRE(err_len == 0);
		REQUIRE(p_testmap_contains[0] != nullptr);
		REQUIRE(p_testmap_contains[1] == nullptr);
		
		cdts_param_ret = (cdts*)xllr_alloc_cdts_buffer(2, 1);
		metaffi::runtime::cdts_wrapper wrapper_contains_params(cdts_param_ret[0].pcdt, cdts_param_ret[0].len, false);
		wrapper_contains_params[0]->type = metaffi_handle_type;
		wrapper_contains_params[0]->cdt_val.metaffi_handle_val.val = testmap_instance;
		wrapper_contains_params[1]->type = metaffi_string8_type;
		wrapper_contains_params[1]->cdt_val.metaffi_string8_val.val = (char*)"key";
		wrapper_contains_params[1]->cdt_val.metaffi_string8_val.length = strlen("key");
		
		long_err_len = 0;
		((void(*)(void*,cdts*,char**,uint64_t*))p_testmap_contains[0])(p_testmap_contains[1], (cdts*)cdts_param_ret, &err, &long_err_len);
		if(err){ FAIL(err); }
		REQUIRE(long_err_len == 0);
		
		metaffi::runtime::cdts_wrapper wrapper_contains_ret(cdts_param_ret[1].pcdt, cdts_param_ret[1].len, false);
		REQUIRE(wrapper_contains_ret[0]->type == metaffi_bool_type);
		REQUIRE(wrapper_contains_ret[0]->cdt_val.metaffi_bool_val.val != 0);
		
		if(cdts_param_ret[0].len + cdts_param_ret[1].len > cdts_cache_size){
			free(cdts_param_ret);
		}
		
		// get
		function_path = "callable=TestMap.Get,instance_required";
		params_types[0].type = metaffi_handle_type;
		params_types[1].type = metaffi_string8_type;
		retvals_types[0].type = metaffi_any_type;
		
		void** p_testmap_get = load_function(module_path.string().c_str(), module_path.string().length(),
		                                          function_path.c_str(), function_path.length(),
		                                          nullptr, retvals_types,
		                                          1, 1,
		                                          &err, &err_len);
		
		if(err){ FAIL(err); }
		REQUIRE(err_len == 0);
		REQUIRE(p_testmap_get[0] != nullptr);
		REQUIRE(p_testmap_get[1] == nullptr);
		
		cdts_param_ret = (cdts*)xllr_alloc_cdts_buffer(2, 1);
		metaffi::runtime::cdts_wrapper wrapper_get_params(cdts_param_ret[0].pcdt, cdts_param_ret[0].len, false);
		wrapper_get_params[0]->type = metaffi_handle_type;
		wrapper_get_params[0]->cdt_val.metaffi_handle_val.val = testmap_instance;
		wrapper_get_params[1]->type = metaffi_string8_type;
		wrapper_get_params[1]->cdt_val.metaffi_string8_val.val = (char*)"key";
		wrapper_get_params[1]->cdt_val.metaffi_string8_val.length = strlen("key");
		
		long_err_len = 0;
		((void(*)(void*,cdts*,char**,uint64_t*))p_testmap_get[0])(p_testmap_get[1], (cdts*)cdts_param_ret, &err, &long_err_len);
		if(err){ FAIL(err); }
		REQUIRE(long_err_len == 0);
		
		metaffi::runtime::cdts_wrapper wrapper_get_ret(cdts_param_ret[1].pcdt, cdts_param_ret[1].len, false);
		REQUIRE(wrapper_get_ret[0]->type == metaffi_int64_type);
		REQUIRE(wrapper_get_ret[0]->cdt_val.metaffi_int64_val.val == 42);
		
		if(cdts_param_ret[0].len + cdts_param_ret[1].len > cdts_cache_size){
			free(cdts_param_ret);
		}
	}
	
	SECTION("runtime_test_target.testmap.get_set_name")
	{
		// load constructor
		std::string function_path = "callable=NewTestMap";
		metaffi_type_with_alias retvals_types[] = {metaffi::runtime::make_type_with_alias(metaffi_handle_type)};
		
		void** pnew_testmap = load_function(module_path.string().c_str(), module_path.string().length(),
		                                    function_path.c_str(), function_path.length(),
		                                    nullptr, retvals_types,
		                                    0, 1,
		                                    &err, &err_len);
		
		if(err){ FAIL(err); }
		REQUIRE(err_len == 0);
		REQUIRE(pnew_testmap[0] != nullptr);
		REQUIRE(pnew_testmap[1] == nullptr);
		
		
		// load getter
		function_path = "field=TestMap.Name,instance_required,getter";
		metaffi_type_with_alias params_types[2] = {0};
		params_types[0].type = metaffi_handle_type;
		retvals_types[0].type = metaffi_string8_type;
		
		void** pget_name = load_function(module_path.string().c_str(), module_path.string().length(),
		                                    function_path.c_str(), function_path.length(),
		                                    nullptr, retvals_types,
		                                    1, 1,
		                                    &err, &err_len);
		
		if(err){ FAIL(err); }
		REQUIRE(err_len == 0);
		REQUIRE(pget_name[0] != nullptr);
		REQUIRE(pget_name[1] == nullptr);
		
		// load setter
		function_path = "field=TestMap.Name,instance_required,setter";
		params_types[0].type = metaffi_handle_type;
		retvals_types[0].type = metaffi_string8_type;
		
		void** pset_name = load_function(module_path.string().c_str(), module_path.string().length(),
		                                 function_path.c_str(), function_path.length(),
		                                 nullptr, retvals_types,
		                                 2, 0,
		                                 &err, &err_len);
		
		if(err){ FAIL(err); }
		REQUIRE(err_len == 0);
		REQUIRE(pset_name[0] != nullptr);
		REQUIRE(pset_name[1] == nullptr);
		
		// create new testmap
		cdts* cdts_param_ret = (cdts*)xllr_alloc_cdts_buffer(0, 1);
		uint64_t long_err_len = 0;
		((void(*)(void*,cdts*,char**,uint64_t*))pnew_testmap[0])(pnew_testmap[1], (cdts*)cdts_param_ret, &err, &long_err_len);

		if(err){ FAIL(err); }
		REQUIRE(long_err_len == 0);
		
		metaffi::runtime::cdts_wrapper wrapper_ret(cdts_param_ret[1].pcdt, cdts_param_ret[1].len, false);
		REQUIRE(wrapper_ret[0]->type == metaffi_handle_type);
		REQUIRE(wrapper_ret[0]->cdt_val.metaffi_handle_val.val != nullptr);
		
		if(cdts_param_ret[0].len + cdts_param_ret[1].len > cdts_cache_size){
			free(cdts_param_ret);
		}
		
		metaffi_handle testmap_instance = wrapper_ret[0]->cdt_val.metaffi_handle_val.val;
		
		
		// get name
		
		cdts_param_ret = (cdts*)xllr_alloc_cdts_buffer(1, 1);
		metaffi::runtime::cdts_wrapper wrapper_get_params(cdts_param_ret[0].pcdt, cdts_param_ret[0].len, false);
		wrapper_get_params[0]->type = metaffi_handle_type;
		wrapper_get_params[0]->cdt_val.metaffi_handle_val.val = testmap_instance;
		
		long_err_len = 0;
		((void(*)(void*,cdts*,char**,uint64_t*))pget_name[0])(pget_name[1], (cdts*)cdts_param_ret, &err, &long_err_len);
		if(err){ FAIL(err); }
		REQUIRE(long_err_len == 0);
		
		metaffi::runtime::cdts_wrapper wrapper_get_ret(cdts_param_ret[1].pcdt, cdts_param_ret[1].len, false);
		REQUIRE(wrapper_get_ret[0]->type == metaffi_string8_type);
		REQUIRE(std::string(wrapper_get_ret[0]->cdt_val.metaffi_string8_val.val, wrapper_get_ret[0]->cdt_val.metaffi_string8_val.length) == "TestMap Name");
		
		// set name to "name is my name"
		
		cdts_param_ret = (cdts*)xllr_alloc_cdts_buffer(2, 0);
		metaffi::runtime::cdts_wrapper wrapper_set_params(cdts_param_ret[0].pcdt, cdts_param_ret[0].len, false);
		wrapper_set_params[0]->type = metaffi_handle_type;
		wrapper_set_params[0]->cdt_val.metaffi_handle_val.val = testmap_instance;
		wrapper_set_params[1]->type = metaffi_string8_type;
		wrapper_set_params[1]->cdt_val.metaffi_string8_val.val = (char*)"name is my name";
		wrapper_set_params[1]->cdt_val.metaffi_string8_val.length = strlen("name is my name");
		
		long_err_len = 0;
		((void(*)(void*,cdts*,char**,uint64_t*))pset_name[0])(pset_name[1], (cdts*)cdts_param_ret, &err, &long_err_len);
		if(err){ FAIL(err); }
		REQUIRE(long_err_len == 0);
		
		// get name again and make sure it is "name is my name"
		cdts* last_get_params = (cdts*)xllr_alloc_cdts_buffer(1, 1);
		wrapper_get_params = metaffi::runtime::cdts_wrapper(last_get_params[0].pcdt, last_get_params[0].len, false);
		wrapper_get_params[0]->type = metaffi_handle_type;
		wrapper_get_params[0]->cdt_val.metaffi_handle_val.val = testmap_instance;
		
		long_err_len = 0;
		((void(*)(void*,cdts*,char**,uint64_t*))pget_name[0])(pget_name[1], (cdts*)last_get_params, &err, &long_err_len);
		if(err){ FAIL(err); }
		REQUIRE(long_err_len == 0);
		
		metaffi::runtime::cdts_wrapper last_get_wrapper(last_get_params[1].pcdt, last_get_params[1].len, false);
		REQUIRE(wrapper_get_ret[0]->type == metaffi_string8_type);
		REQUIRE(std::string(wrapper_get_ret[0]->cdt_val.metaffi_string8_val.val, last_get_wrapper[0]->cdt_val.metaffi_string8_val.length) == "name is my name");
	}
	
	SECTION("runtime_test_target.testmap.get_set_name_from_empty_struct")
	{
		// load constructor
		std::string function_path = "callable=TestMap.EmptyStruct";
		metaffi_type_with_alias retvals_types[] = {metaffi::runtime::make_type_with_alias(metaffi_handle_type)};
		
		void** pnew_testmap = load_function(module_path.string().c_str(), module_path.string().length(),
		                                    function_path.c_str(), function_path.length(),
		                                    nullptr, retvals_types,
		                                    0, 1,
		                                    &err, &err_len);
		
		if(err){ FAIL(err); }
		REQUIRE(err_len == 0);
		REQUIRE(pnew_testmap[0] != nullptr);
		REQUIRE(pnew_testmap[1] == nullptr);
		
		
		// load getter
		function_path = "field=TestMap.Name,instance_required,getter";
		metaffi_type_with_alias params_types[2] = {0};
		params_types[0].type = metaffi_handle_type;
		retvals_types[0].type = metaffi_string8_type;
		
		void** pget_name = load_function(module_path.string().c_str(), module_path.string().length(),
		                                 function_path.c_str(), function_path.length(),
		                                 nullptr, retvals_types,
		                                 1, 1,
		                                 &err, &err_len);
		
		if(err){ FAIL(err); }
		REQUIRE(err_len == 0);
		REQUIRE(pget_name[0] != nullptr);
		REQUIRE(pget_name[1] == nullptr);
		
		// load setter
		function_path = "field=TestMap.Name,instance_required,setter";
		params_types[0].type = metaffi_handle_type;
		retvals_types[0].type = metaffi_string8_type;
		
		void** pset_name = load_function(module_path.string().c_str(), module_path.string().length(),
		                                 function_path.c_str(), function_path.length(),
		                                 nullptr, retvals_types,
		                                 2, 0,
		                                 &err, &err_len);
		
		if(err){ FAIL(err); }
		REQUIRE(err_len == 0);
		REQUIRE(pset_name[0] != nullptr);
		REQUIRE(pset_name[1] == nullptr);
		
		// create new testmap
		cdts* cdts_param_ret = (cdts*)xllr_alloc_cdts_buffer(0, 1);
		uint64_t long_err_len = 0;
		((void(*)(void*,cdts*,char**,uint64_t*))pnew_testmap[0])(pnew_testmap[1], (cdts*)cdts_param_ret, &err, &long_err_len);
		
		if(err){ FAIL(err); }
		REQUIRE(long_err_len == 0);
		
		metaffi::runtime::cdts_wrapper wrapper_ret(cdts_param_ret[1].pcdt, cdts_param_ret[1].len, false);
		REQUIRE(wrapper_ret[0]->type == metaffi_handle_type);
		REQUIRE(wrapper_ret[0]->cdt_val.metaffi_handle_val.val != nullptr);
		
		if(cdts_param_ret[0].len + cdts_param_ret[1].len > cdts_cache_size){
			free(cdts_param_ret);
		}
		
		metaffi_handle testmap_instance = wrapper_ret[0]->cdt_val.metaffi_handle_val.val;
		
		
		// get name
		
		cdts_param_ret = (cdts*)xllr_alloc_cdts_buffer(1, 1);
		metaffi::runtime::cdts_wrapper wrapper_get_params(cdts_param_ret[0].pcdt, cdts_param_ret[0].len, false);
		wrapper_get_params[0]->type = metaffi_handle_type;
		wrapper_get_params[0]->cdt_val.metaffi_handle_val.val = testmap_instance;
		
		long_err_len = 0;
		((void(*)(void*,cdts*,char**,uint64_t*))pget_name[0])(pget_name[1], (cdts*)cdts_param_ret, &err, &long_err_len);
		if(err){ FAIL(err); }
		REQUIRE(long_err_len == 0);
		
		metaffi::runtime::cdts_wrapper wrapper_get_ret(cdts_param_ret[1].pcdt, cdts_param_ret[1].len, false);
		REQUIRE(wrapper_get_ret[0]->type == metaffi_string8_type);
		REQUIRE(std::string(wrapper_get_ret[0]->cdt_val.metaffi_string8_val.val, wrapper_get_ret[0]->cdt_val.metaffi_string8_val.length) == "");
		
		// set name to "name is my name"
		
		cdts_param_ret = (cdts*)xllr_alloc_cdts_buffer(2, 0);
		metaffi::runtime::cdts_wrapper wrapper_set_params(cdts_param_ret[0].pcdt, cdts_param_ret[0].len, false);
		wrapper_set_params[0]->type = metaffi_handle_type;
		wrapper_set_params[0]->cdt_val.metaffi_handle_val.val = testmap_instance;
		wrapper_set_params[1]->type = metaffi_string8_type;
		wrapper_set_params[1]->cdt_val.metaffi_string8_val.val = (char*)"name is my name";
		wrapper_set_params[1]->cdt_val.metaffi_string8_val.length = strlen("name is my name");
		
		long_err_len = 0;
		((void(*)(void*,cdts*,char**,uint64_t*))pset_name[0])(pset_name[1], (cdts*)cdts_param_ret, &err, &long_err_len);
		if(err){ FAIL(err); }
		REQUIRE(long_err_len == 0);
		
		// get name again and make sure it is "name is my name"
		cdts* last_get_params = (cdts*)xllr_alloc_cdts_buffer(1, 1);
		wrapper_get_params = metaffi::runtime::cdts_wrapper(last_get_params[0].pcdt, last_get_params[0].len, false);
		wrapper_get_params[0]->type = metaffi_handle_type;
		wrapper_get_params[0]->cdt_val.metaffi_handle_val.val = testmap_instance;
		
		long_err_len = 0;
		((void(*)(void*,cdts*,char**,uint64_t*))pget_name[0])(pget_name[1], (cdts*)last_get_params, &err, &long_err_len);
		if(err){ FAIL(err); }
		REQUIRE(long_err_len == 0);
		
		metaffi::runtime::cdts_wrapper last_get_wrapper(last_get_params[1].pcdt, last_get_params[1].len, false);
		REQUIRE(wrapper_get_ret[0]->type == metaffi_string8_type);
		REQUIRE(std::string(wrapper_get_ret[0]->cdt_val.metaffi_string8_val.val, last_get_wrapper[0]->cdt_val.metaffi_string8_val.length) == "name is my name");
	}
	
	SECTION("runtime_test_target.wait_a_bit")
	{
		// get five_seconds global
		metaffi_type_with_alias var_type[] = {metaffi::runtime::make_type_with_alias(metaffi_int64_type)};
		std::string variable_path = "global=FiveSeconds,getter";
		void** pfive_seconds_getter = load_function(module_path.string().c_str(), module_path.string().length(),
		                                            variable_path.c_str(), variable_path.length(),
		                                            nullptr, var_type,
		                                            0, 1,
		                                            &err, &err_len);
		if(err){ FAIL(err); }
		REQUIRE(err_len == 0);
		REQUIRE(pfive_seconds_getter[0] != nullptr);
		REQUIRE(pfive_seconds_getter[1] == nullptr);
		
		cdts* getter_ret = (cdts*)xllr_alloc_cdts_buffer(0, 1);
		
		uint64_t long_err_len = 0;
		((void(*)(void*,cdts*,char**,uint64_t*))pfive_seconds_getter[0])(pfive_seconds_getter[1], (cdts*)getter_ret, &err, &long_err_len);
		if(err){ FAIL(err); }
		
		REQUIRE(getter_ret->pcdt->type == metaffi_int64_type);
		REQUIRE(getter_ret->pcdt->cdt_val.metaffi_int64_val.val == 5000000000);
		
		int64_t five = getter_ret->pcdt->cdt_val.metaffi_int64_val.val;
		
		// call wait_a_bit
		std::string function_path = "callable=WaitABit";
		metaffi_type_with_alias params_types[] = {metaffi::runtime::make_type_with_alias(metaffi_int64_type)};
		
		void** pwait_a_bit = load_function(module_path.string().c_str(), module_path.string().length(),
		                                   function_path.c_str(), function_path.length(),
		                                   params_types, nullptr,
		                                   1, 0,
		                                   &err, &err_len);
		
		if(err){ FAIL(err); }
		REQUIRE(err_len == 0);
		REQUIRE(pwait_a_bit[0] != nullptr);
		REQUIRE(pwait_a_bit[1] == nullptr);
		
		
		cdts* cdts_param_ret = (cdts*)xllr_alloc_cdts_buffer(1, 0);
		metaffi::runtime::cdts_wrapper wrapper(cdts_param_ret[0].pcdt, cdts_param_ret[0].len, false);
		wrapper[0]->type = metaffi_int64_type;
		wrapper[0]->cdt_val.metaffi_int64_val.val = five;
		
		long_err_len = 0;
		((void(*)(void*,cdts*,char**,uint64_t*))pwait_a_bit[0])(pwait_a_bit[1], (cdts*)cdts_param_ret, &err, &long_err_len);
		if(err){ FAIL(err); }
		REQUIRE(long_err_len == 0);
		
		if(cdts_param_ret[0].len + cdts_param_ret[1].len > cdts_cache_size){
			free(cdts_param_ret);
		}
	}
	
}