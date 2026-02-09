# Go Compiler Plugin: C++ shared library that implements compiler_plugin_interface
# by running sdk/compiler/go executable (go_compiler) via GoCompilerWrapper (same pattern as go_idl_compiler).

find_or_install_package(Boost COMPONENTS filesystem)

collect_c_cpp_files("${metaffi_sdk_root}/compiler/go/cpp_wrapper" sdk_go_compiler_wrapper)

set(go_compiler_plugin_src
	${CMAKE_CURRENT_LIST_DIR}/go_compiler_plugin.cpp
	${CMAKE_CURRENT_LIST_DIR}/compiler_api.cpp
	${sdk_go_compiler_wrapper_src}
)

set(go_compiler_plugin_includes
	${sdk_include_dir}
	${metaffi_sdk_root}
	${metaffi_sdk_root}/compiler
	${CMAKE_CURRENT_LIST_DIR}
	${Boost_INCLUDE_DIRS}
)

c_cpp_shared_lib(metaffi.compiler.go
	"${go_compiler_plugin_src};${sdk_src}"
	"${go_compiler_plugin_includes}"
	"Boost::filesystem"
	"./go"
)

# Ensure go_compiler executable is built and copy from SDK output to METAFFI_HOME/go (plugin loads from there)
# Source: METAFFI_HOME/sdk/compiler/go/go_compiler.exe (same as SDK build output)
add_dependencies(metaffi.compiler.go go_compiler_exe)
set(GO_COMPILER_EXE_FROM_SDK "$ENV{METAFFI_HOME}/sdk/compiler/go/go_compiler${CMAKE_EXECUTABLE_SUFFIX}")
add_custom_command(TARGET metaffi.compiler.go POST_BUILD
	COMMAND ${CMAKE_COMMAND} -E make_directory $ENV{METAFFI_HOME}/go
	COMMAND ${CMAKE_COMMAND} -E copy_if_different
		${GO_COMPILER_EXE_FROM_SDK}
		$ENV{METAFFI_HOME}/go/
	COMMENT "Copy go_compiler from SDK output to METAFFI_HOME/go"
)

set(metaffi.compiler.go metaffi.compiler.go PARENT_SCOPE)

# E2E tests for Go host compiler (generate host from xllr.test, then run go test)
add_subdirectory(test)