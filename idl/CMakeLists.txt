# load conan packages
find_or_install_package(Boost COMPONENTS filesystem)
find_or_install_package(doctest)


# build Go metaffi.compiler.go dynamic library
go_build(metaffi.idl.go
		SOURCE_DIR
			"${CMAKE_CURRENT_LIST_DIR}"
		OUTPUT_DIR
			"./go"
		OUTPUT_NAME
			metaffi.idl.go
		DEPENDENT
			${CMAKE_CURRENT_LIST_DIR}
)
set(metaffi.idl.go metaffi.idl.go PARENT_SCOPE)

# go test for metaffi.compiler.go
add_test(NAME metaffi_idl_go_test
		COMMAND go test
		WORKING_DIRECTORY ${CMAKE_CURRENT_SOURCE_DIR}/IDLCompiler
)

