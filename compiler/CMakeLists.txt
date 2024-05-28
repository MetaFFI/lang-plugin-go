
#------------------------------- Go Language Support - Compiler Plugin ------------------------------------------
set(target_name "metaffi.compiler.go")

# build
add_go_target(${target_name})

go_get(${target_name} WORKING_DIRECTORY ${CMAKE_CURRENT_LIST_DIR})
go_build(${target_name} WORKING_DIRECTORY ${CMAKE_CURRENT_LIST_DIR})
add_go_test("Go compiler test" WORKING_DIRECTORY ${CMAKE_CURRENT_LIST_DIR})

# install
install(FILES ${CMAKE_BINARY_DIR}/${target_name}${CMAKE_SHARED_LIBRARY_SUFFIX}
		DESTINATION .)

copy_file_from_project_binary_to_metaffi_home(${target_name} ${target_name}${CMAKE_SHARED_LIBRARY_SUFFIX})