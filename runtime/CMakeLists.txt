
#------------------------------- Go Language Support - Runtime Plugin ------------------------------------------
project(lang-plugin-go-runtime VERSION 0.0.1)

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

# Copy to metaffi directory
copy_file_from_project_binary_to_metaffi_home(${target_name} "${target_name}${CMAKE_SHARED_LIBRARY_SUFFIX}")

# install
install_target(${target_name} .)