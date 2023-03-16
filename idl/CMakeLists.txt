
#------------------------------- "go" IDL Plugin ------------------------------------------

set(target_name "metaffi.idl.go")

# build
add_custom_target(${target_name} ALL)
go_get(${target_name} WORKING_DIRECTORY ${CMAKE_CURRENT_LIST_DIR})
go_build(${target_name} WORKING_DIRECTORY ${CMAKE_CURRENT_LIST_DIR})

# Unitest
add_go_test("go IDL plugin test")

# install
install(FILES ${CMAKE_BINARY_DIR}/${target_name}${CMAKE_SHARED_LIBRARY_SUFFIX}
		DESTINATION .)