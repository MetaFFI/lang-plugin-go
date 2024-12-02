# load conan packages
find_or_install_package(Boost COMPONENTS filesystem)
find_or_install_package(doctest)


# unit test go-runtime module. run ./test/run_test.py
add_test(NAME go_runtime_test
		COMMAND ${PYTHON_EXECUTABLE} ${PYTHON_EXECUTABLE_ARG} run_test.py
		WORKING_DIRECTORY ${CMAKE_CURRENT_SOURCE_DIR}/test
)
