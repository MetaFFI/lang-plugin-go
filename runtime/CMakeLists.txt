# load conan packages
find_or_install_package(Boost COMPONENTS filesystem)
find_or_install_package(doctest)

collect_c_cpp_files(${CMAKE_CURRENT_LIST_DIR} xllr.go)

# Collect SDK Go runtime_manager sources (runtime_manager.cpp, module.cpp, entity.cpp; test excluded by collect)
collect_c_cpp_files("${metaffi_sdk_root}/runtime_manager/go" sdk_go_runtime_manager)

set(sdk_go_plugin_src
	${sdk_src}
	${sdk_go_runtime_manager_src}
)

# SDK include: sdk_include_dir + metaffi_sdk_root for <runtime_manager/go/...> and <utils/...>
set(sdk_go_plugin_include ${sdk_include_dir};${metaffi_sdk_root})

# Build dynamic library (xllr.go plugin using SDK go_runtime_manager)
c_cpp_shared_lib(xllr.go
		"${xllr.go_src};${sdk_go_plugin_src}"
		"${sdk_go_plugin_include};${Boost_INCLUDE_DIRS};"
		"Boost::filesystem;"
		"./go")

set(xllr.go xllr.go PARENT_SCOPE)

# build go guest
add_custom_command(
    OUTPUT ${CMAKE_CURRENT_SOURCE_DIR}/test/TestRuntime_MetaFFIGuest.*
    COMMAND ${PYTHON_EXECUTABLE} ${PYTHON_EXECUTABLE_ARG} build_guest.py
    WORKING_DIRECTORY ${CMAKE_CURRENT_SOURCE_DIR}/test
    DEPENDS ${CMAKE_CURRENT_SOURCE_DIR}/test/TestRuntime.go
)
add_custom_target(build_go_guest DEPENDS ${CMAKE_CURRENT_SOURCE_DIR}/test/TestRuntime.go)

# run go_api_test.cpp doctest unit test (links plugin and uses same SDK sources for in-process test)
c_cpp_exe(go_api_test
		"${xllr.go_src};${sdk_go_plugin_src};go_api_test.cpp"
		"${sdk_go_plugin_include};${Boost_INCLUDE_DIRS};${doctest_INCLUDE_DIRS};"
		"doctest::doctest;Boost::filesystem;xllr.go"
		".")

add_custom_command(TARGET go_api_test POST_BUILD
    COMMAND ${CMAKE_COMMAND} -E remove $ENV{METAFFI_HOME}/xllr.go.dll
)
add_dependencies(go_api_test xllr.go build_go_guest)

add_test(NAME go_api_test_build_guest
		COMMAND ${PYTHON_EXECUTABLE} ${CMAKE_CURRENT_LIST_DIR}/test/build_guest.py
		WORKING_DIRECTORY ${CMAKE_CURRENT_LIST_DIR}/test)
set_tests_properties(go_api_test_build_guest PROPERTIES FIXTURES_SETUP go_guest_built)

add_test(NAME go_api_test
		COMMAND $ENV{METAFFI_HOME}/go_api_test)
set_tests_properties(go_api_test PROPERTIES FIXTURES_REQUIRED go_guest_built)

set(go_api_test go_api_test PARENT_SCOPE)
