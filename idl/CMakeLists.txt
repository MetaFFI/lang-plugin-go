# IDL plugin for Go: C++ plugin that invokes SDK go_idl_compiler executable

find_or_install_package(Boost COMPONENTS filesystem)

collect_c_cpp_files(${CMAKE_CURRENT_LIST_DIR} metaffi_idl_go)
collect_c_cpp_files("${metaffi_sdk_root}/idl_compiler/go/cpp_wrapper" sdk_go_idl_wrapper)

set(metaffi_idl_go_all_src
	${sdk_src}
	${metaffi_idl_go_src}
	${sdk_go_idl_wrapper_src}
)

set(metaffi_idl_go_include ${sdk_include_dir};${metaffi_sdk_root})

c_cpp_shared_lib(metaffi.idl.go
	"${metaffi_idl_go_all_src}"
	"${metaffi_idl_go_include};${Boost_INCLUDE_DIRS}"
	"Boost::filesystem"
	"./go")

# Ensure go_idl_compiler is built and copy from SDK output to METAFFI_HOME/go (plugin loads from there)
# Source: METAFFI_HOME/sdk/idl_compiler/go/go_idl_compiler.exe (same as SDK build output)
add_dependencies(metaffi.idl.go go_idl_compiler_exe)
set(GO_IDL_COMPILER_EXE_FROM_SDK "$ENV{METAFFI_HOME}/sdk/idl_compiler/go/go_idl_compiler${CMAKE_EXECUTABLE_SUFFIX}")
add_custom_command(TARGET metaffi.idl.go POST_BUILD
	COMMAND ${CMAKE_COMMAND} -E make_directory $ENV{METAFFI_HOME}/go
	COMMAND ${CMAKE_COMMAND} -E copy_if_different
		${GO_IDL_COMPILER_EXE_FROM_SDK}
		$ENV{METAFFI_HOME}/go/
	COMMENT "Copy go_idl_compiler from SDK output to METAFFI_HOME/go"
)

set(metaffi.idl.go metaffi.idl.go PARENT_SCOPE)
