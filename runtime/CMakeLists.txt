find_package(Boost REQUIRED COMPONENTS filesystem) # Boost library

include(CTest)
enable_testing()

set(target_name "xllr.go")

file(GLOB SRC_FILES ${CMAKE_CURRENT_LIST_DIR}/*.cpp ${CMAKE_CURRENT_LIST_DIR}/*.h ${METAFFI_SDK}/utils/*.cpp)
# build xllr.go
add_library(${target_name}
			SHARED
			${SRC_FILES}
)
target_include_directories(${target_name} PUBLIC ${CMAKE_CURRENT_LIST_DIR} ${Boost_INCLUDE_DIRS})

set_property(TARGET ${target_name} PROPERTY POSITION_INDEPENDENT_CODE ON)

# remove the "lib" prefix in linux to get the expected filename "xllr.go.[ext]"
SET_TARGET_PROPERTIES(${target_name} PROPERTIES PREFIX "")

target_link_libraries(${target_name} Boost::filesystem)

if(WIN32)
	target_link_directories(${target_name} PUBLIC ${BOOST_LIBRARYDIR})
else()
	target_link_libraries(${target_name} ${Boost_LIBRARIES})
endif()

# Copy to metaffi directory
add_custom_command(TARGET ${target_name} POST_BUILD
		COMMAND "${CMAKE_COMMAND}" -E copy
		"${CMAKE_CURRENT_BINARY_DIR}/${target_name}${CMAKE_SHARED_LIBRARY_SUFFIX}"
		"$ENV{METAFFI_HOME}/${target_name}${CMAKE_SHARED_LIBRARY_SUFFIX}"
		COMMENT "Copy ${CMAKE_CURRENT_BINARY_DIR}/${target_name}${CMAKE_SHARED_LIBRARY_SUFFIX} -> $ENV{METAFFI_HOME}/${target_name}${CMAKE_SHARED_LIBRARY_SUFFIX}")


# install
install(FILES ${CMAKE_CURRENT_BINARY_DIR}/${target_name}${CMAKE_SHARED_LIBRARY_SUFFIX}
		DESTINATION "lib")