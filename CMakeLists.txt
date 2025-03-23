
# add SDK
add_subdirectory("${CMAKE_CURRENT_LIST_DIR}/plugin-sdk")

# runtime
add_subdirectory("${CMAKE_CURRENT_LIST_DIR}/runtime")

# go-runtime
add_subdirectory("${CMAKE_CURRENT_LIST_DIR}/go-runtime")

# idl
add_subdirectory("${CMAKE_CURRENT_LIST_DIR}/idl")

# compiler
add_subdirectory("${CMAKE_CURRENT_LIST_DIR}/compiler")

# api
add_subdirectory("${CMAKE_CURRENT_LIST_DIR}/api")

add_custom_target(go ALL
	DEPENDS xllr.go metaffi.compiler.go metaffi.idl.go go_api_test
)
set_target_properties(go PROPERTIES EXCLUDE_FROM_ALL TRUE)

# set go to parent scope
set(go ${go} PARENT_SCOPE)