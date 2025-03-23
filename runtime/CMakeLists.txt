# load conan packages
find_or_install_package(Boost COMPONENTS filesystem)
find_or_install_package(doctest)


# build dynamic library

collect_c_cpp_files(${CMAKE_CURRENT_LIST_DIR} xllr.go)

# get cpython from "deps" as vcpkg doesn't have python3.11 and python3.12
c_cpp_shared_lib(xllr.go
		"${xllr.go_src};${sdk_src}"
		"${sdk_include_dir};${Boost_INCLUDE_DIRS};"
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

# run go_api_test.cpp doctest unit test
c_cpp_exe(go_api_test
		"${xllr.openjdk_src};${sdk_src};go_api_test.cpp"
		"${sdk_include_dir};${Boost_INCLUDE_DIRS};${doctest_INCLUDE_DIRS};"
		"doctest::doctest;Boost::filesystem;xllr.go"
		".")

add_custom_command(TARGET go_api_test POST_BUILD
    COMMAND ${CMAKE_COMMAND} -E remove $ENV{METAFFI_HOME}/xllr.go.dll
)
add_dependencies(go_api_test xllr.go build_go_guest)

add_test(NAME go_api_test
		COMMAND ${PYTHON_EXECUTABLE} ${CMAKE_CURRENT_LIST_DIR}/test/build_guest.py && $ENV{METAFFI_HOME}/go_api_test)

set(go_api_test go_api_test PARENT_SCOPE)