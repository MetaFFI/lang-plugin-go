
# add SDK
if(NOT DEFINED sdk_src)
	message(FATAL_ERROR "sdk_src is not defined. Ensure sdk/CMakeLists.txt is added before lang-plugin-go.")
endif()

# runtime
add_subdirectory("${CMAKE_CURRENT_LIST_DIR}/runtime")


# idl
add_subdirectory("${CMAKE_CURRENT_LIST_DIR}/idl")

# compiler
add_subdirectory("${CMAKE_CURRENT_LIST_DIR}/compiler")

# test (C++ host tests calling Go guest module)
add_subdirectory("${CMAKE_CURRENT_LIST_DIR}/test")

add_custom_target(go)
add_dependencies(go xllr.go metaffi.idl.go metaffi.compiler.go go_host_test)

# set go to parent scope
set(go ${go} PARENT_SCOPE)
