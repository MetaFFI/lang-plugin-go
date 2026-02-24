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
set(GO_RUNTIME_TEST_DIR "${CMAKE_CURRENT_SOURCE_DIR}/test")
add_custom_command(
    OUTPUT ${CMAKE_CURRENT_SOURCE_DIR}/test/TestRuntime_MetaFFIGuest.*
    COMMAND $ENV{METAFFI_HOME}/metaffi${CMAKE_EXECUTABLE_SUFFIX} -c --idl TestRuntime.go -g go
    WORKING_DIRECTORY ${GO_RUNTIME_TEST_DIR}
    DEPENDS ${CMAKE_CURRENT_SOURCE_DIR}/test/TestRuntime.go
)
add_custom_target(build_go_guest DEPENDS ${CMAKE_CURRENT_SOURCE_DIR}/test/TestRuntime.go)
add_dependencies(build_go_guest metaffi)

# run go_api_test.cpp doctest unit test (links plugin and uses same SDK sources for in-process test)
c_cpp_exe(go_api_test
		"${xllr.go_src};${sdk_go_plugin_src};go_api_test.cpp"
		"${sdk_go_plugin_include};${Boost_INCLUDE_DIRS};${doctest_INCLUDE_DIRS};"
		"doctest::doctest;Boost::filesystem;xllr.go"
		"./go")
add_dependencies(go_api_test xllr.go build_go_guest)

add_test(NAME go_api_test_build_guest
		COMMAND $ENV{METAFFI_HOME}/metaffi${CMAKE_EXECUTABLE_SUFFIX} -c --idl TestRuntime.go -g go
		WORKING_DIRECTORY ${CMAKE_CURRENT_LIST_DIR}/test)
set_tests_properties(go_api_test_build_guest PROPERTIES FIXTURES_SETUP go_guest_built)

add_test(NAME go_api_test
		COMMAND $ENV{METAFFI_HOME}/go/go_api_test)
set_tests_properties(go_api_test PROPERTIES FIXTURES_REQUIRED go_guest_built)

set(go_runtime_test_cgo_flags "-I$ENV{METAFFI_HOME} -I$ENV{METAFFI_HOME}/include -I$ENV{METAFFI_SOURCE_ROOT}/sdk")
if(WIN32)
	set(go_runtime_test_path "$ENV{METAFFI_HOME}/sdk/api/cpp;$ENV{METAFFI_HOME};$ENV{METAFFI_HOME}/go;$ENV{PATH}")
	string(REPLACE ";" "\\;" go_runtime_test_path "${go_runtime_test_path}")
	set_tests_properties(go_api_test_build_guest PROPERTIES
		ENVIRONMENT "PATH=${go_runtime_test_path};CGO_CFLAGS=${go_runtime_test_cgo_flags};CGO_CPPFLAGS=${go_runtime_test_cgo_flags}")
	set_tests_properties(go_api_test PROPERTIES
		ENVIRONMENT "PATH=${go_runtime_test_path};CGO_CFLAGS=${go_runtime_test_cgo_flags};CGO_CPPFLAGS=${go_runtime_test_cgo_flags}")
elseif(APPLE)
	set_tests_properties(go_api_test_build_guest PROPERTIES
		ENVIRONMENT "DYLD_LIBRARY_PATH=$ENV{METAFFI_HOME}/sdk/api/cpp:$ENV{DYLD_LIBRARY_PATH};CGO_CFLAGS=${go_runtime_test_cgo_flags};CGO_CPPFLAGS=${go_runtime_test_cgo_flags}")
	set_tests_properties(go_api_test PROPERTIES
		ENVIRONMENT "DYLD_LIBRARY_PATH=$ENV{METAFFI_HOME}/sdk/api/cpp:$ENV{DYLD_LIBRARY_PATH};CGO_CFLAGS=${go_runtime_test_cgo_flags};CGO_CPPFLAGS=${go_runtime_test_cgo_flags}")
else()
	set_tests_properties(go_api_test_build_guest PROPERTIES
		ENVIRONMENT "LD_LIBRARY_PATH=$ENV{METAFFI_HOME}/sdk/api/cpp:$ENV{LD_LIBRARY_PATH};CGO_CFLAGS=${go_runtime_test_cgo_flags};CGO_CPPFLAGS=${go_runtime_test_cgo_flags}")
	set_tests_properties(go_api_test PROPERTIES
		ENVIRONMENT "LD_LIBRARY_PATH=$ENV{METAFFI_HOME}/sdk/api/cpp:$ENV{LD_LIBRARY_PATH};CGO_CFLAGS=${go_runtime_test_cgo_flags};CGO_CPPFLAGS=${go_runtime_test_cgo_flags}")
endif()

set(go_api_test go_api_test PARENT_SCOPE)
