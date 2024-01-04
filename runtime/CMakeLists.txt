
#------------------------------- Go Language Support - Runtime Plugin ------------------------------------------

add_boost(filesystem)

set(target_name "xllr.go")

# build xllr.go
get_dir_sources(${CMAKE_CURRENT_LIST_DIR} ${target_name}_sources)
get_sdk_utils_sources(utils_srcs)
get_sdk_runtime_sources(runtime_srcs)

add_library(${target_name} SHARED ${${target_name}_sources} ${utils_srcs} ${runtime_srcs})
target_include_directories(${target_name} PRIVATE ${CMAKE_CURRENT_LIST_DIR} ${Boost_INCLUDE_DIRS} ${METAFFI_SDK})
set_target_properties(${target_name} PROPERTIES PREFIX "" POSITION_INDEPENDENT_CODE ON)
target_link_libraries(${target_name} PRIVATE Boost::filesystem)
link_to_dl_for_non_windows(${target_name})

set_target_properties( ${target_name}
		PROPERTIES
		ARCHIVE_OUTPUT_DIRECTORY "${PROJECT_BINARY_DIR}"
		LIBRARY_OUTPUT_DIRECTORY "${PROJECT_BINARY_DIR}"
		RUNTIME_OUTPUT_DIRECTORY "${PROJECT_BINARY_DIR}"
		)


# unitests
add_executable(go_api_test "go_api_test.cpp" ${utils_srcs} ${runtime_srcs})
target_link_libraries(go_api_test PRIVATE Catch2::Catch2 ${target_name} Boost::filesystem)
target_include_directories(go_api_test PRIVATE ${METAFFI_SDK})
add_test(NAME "Go Runtime API" COMMAND go_api_test)
if(NOT WIN32)
	add_custom_command(TARGET go_api_test POST_BUILD COMMAND patchelf --set-rpath "$ORIGIN:$ORIGIN/lib:$ENV{METAFFI_HOME}:${Boost_LIBRARY_DIRS}" $<TARGET_FILE:go_api_test>)
endif()


# install
install_target(${target_name} .)