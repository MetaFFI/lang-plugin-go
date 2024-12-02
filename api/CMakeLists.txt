# load conan packages
find_or_install_package(Boost COMPONENTS filesystem)
find_or_install_package(doctest)


# * ---- Cross Language tests ----
add_test(NAME go_api_cross_pl_tests
		COMMAND ${PYTHON_EXECUTABLE} ${PYTHON_EXECUTABLE_ARG} run_api_tests.py
		WORKING_DIRECTORY ${CMAKE_CURRENT_SOURCE_DIR}/tests
)