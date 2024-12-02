# load conan packages
find_or_install_package(Boost COMPONENTS filesystem)
find_or_install_package(doctest)


# build Go metaffi.compiler.go dynamic library
go_build(metaffi.compiler.go
		SOURCE_DIR
			"${CMAKE_CURRENT_LIST_DIR}"
		OUTPUT_DIR
			"./go"
		OUTPUT_NAME
			metaffi.compiler.go
		DEPENDENT
			${CMAKE_CURRENT_LIST_DIR}
)
set(metaffi.compiler.go metaffi.compiler.go PARENT_SCOPE)


# go test for metaffi.compiler.go
add_test(NAME metaffi_compiler_go_test
		COMMAND go test
		WORKING_DIRECTORY ${CMAKE_CURRENT_SOURCE_DIR}
)

