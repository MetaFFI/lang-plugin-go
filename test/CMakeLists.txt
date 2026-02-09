# Go host tests (C++ calling Go guest module via xllr.go)

find_or_install_package(Boost COMPONENTS filesystem)
find_or_install_package(doctest)

set(go_host_test_src
	${CMAKE_CURRENT_LIST_DIR}/test_main.cpp
	${CMAKE_CURRENT_LIST_DIR}/go_test_env.cpp
	${CMAKE_CURRENT_LIST_DIR}/go_wrappers.cpp
	${CMAKE_CURRENT_LIST_DIR}/test_core.cpp
	${CMAKE_CURRENT_LIST_DIR}/test_errors.cpp
	${CMAKE_CURRENT_LIST_DIR}/test_callbacks.cpp
	${CMAKE_CURRENT_LIST_DIR}/test_arrays.cpp
	${CMAKE_CURRENT_LIST_DIR}/test_classes.cpp
	${CMAKE_CURRENT_LIST_DIR}/test_module_state.cpp
	${CMAKE_CURRENT_LIST_DIR}/test_collections.cpp
	${CMAKE_CURRENT_LIST_DIR}/test_primitives.cpp
	${CMAKE_CURRENT_LIST_DIR}/test_generics_varargs.cpp
	${CMAKE_CURRENT_LIST_DIR}/test_enum_sub.cpp
	${CMAKE_CURRENT_LIST_DIR}/test_third_party.cpp
)

set(go_host_test_includes
	${sdk_include_dir}
	${metaffi_sdk_root}/api/cpp/include
	${CMAKE_CURRENT_LIST_DIR}
	${doctest_INCLUDE_DIRS}
	${Boost_INCLUDE_DIRS}
)

set(go_host_test_libs
	doctest::doctest
	Boost::filesystem
	metaffi.api.cpp
)

c_cpp_exe(go_host_test
	"${go_host_test_src}"
	"${go_host_test_includes}"
	"${go_host_test_libs}"
	"$ENV{METAFFI_HOME}"
)

# Build guest module with metaffi first (IDL + guest compiler). Only if this succeeds do we run go_host_test.
set(go_guest_idl_dir "${metaffi_sdk_root}/test_modules/guest_modules/go")
add_test(NAME go_host_test_build_guest COMMAND $ENV{METAFFI_HOME}/metaffi${CMAKE_EXECUTABLE_SUFFIX}
	-c --idl "${go_guest_idl_dir}" -g go -o "${go_guest_idl_dir}/test_bin")
set_tests_properties(go_host_test_build_guest PROPERTIES FIXTURES_SETUP go_host_fixture)
if(WIN32)
	set(go_host_test_path "$ENV{METAFFI_HOME}/sdk/api/cpp;$ENV{PATH}")
	string(REPLACE ";" "\\;" go_host_test_path "${go_host_test_path}")
	set_tests_properties(go_host_test_build_guest PROPERTIES ENVIRONMENT "PATH=${go_host_test_path}")
elseif(APPLE)
	set_tests_properties(go_host_test_build_guest PROPERTIES
		ENVIRONMENT "DYLD_LIBRARY_PATH=$ENV{METAFFI_HOME}/sdk/api/cpp:$ENV{DYLD_LIBRARY_PATH}")
else()
	set_tests_properties(go_host_test_build_guest PROPERTIES
		ENVIRONMENT "LD_LIBRARY_PATH=$ENV{METAFFI_HOME}/sdk/api/cpp:$ENV{LD_LIBRARY_PATH}")
endif()

add_test(NAME go_host_test COMMAND $ENV{METAFFI_HOME}/go_host_test${CMAKE_EXECUTABLE_SUFFIX})
set_tests_properties(go_host_test PROPERTIES FIXTURES_REQUIRED go_host_fixture)
if(WIN32)
	set(go_host_test_path "$ENV{METAFFI_HOME}/sdk/api/cpp;$ENV{PATH}")
	string(REPLACE ";" "\\;" go_host_test_path "${go_host_test_path}")
	set_tests_properties(go_host_test PROPERTIES
		ENVIRONMENT "PATH=${go_host_test_path}"
	)
elseif(APPLE)
	set_tests_properties(go_host_test PROPERTIES
		ENVIRONMENT "DYLD_LIBRARY_PATH=$ENV{METAFFI_HOME}/sdk/api/cpp:$ENV{DYLD_LIBRARY_PATH}"
	)
else()
	set_tests_properties(go_host_test PROPERTIES
		ENVIRONMENT "LD_LIBRARY_PATH=$ENV{METAFFI_HOME}/sdk/api/cpp:$ENV{LD_LIBRARY_PATH}"
	)
endif()
set(go_host_test go_host_test PARENT_SCOPE)
