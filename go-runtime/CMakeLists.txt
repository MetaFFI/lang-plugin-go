
#------------------------------- Go Language Support - Go-Runtime package ------------------------------------------

# Add METAFFI_HOME environment variable to "CGO_CFLAGS", so MetaFFIGoHostCommon.go can compile its CGo code

if("$ENV{CGO_CFLAGS}" STREQUAL "")
	message(WARNING "CGO_CFLAGS must be set to METAFFI_HOME to use the Go plugin")
	SET(ENV{CGO_CFLAGS} $ENV{METAFFI_HOME})
endif()


if(WIN32)
	set(PYEXEC "py")
else()
	set(PYEXEC "python3")
endif()


add_test(NAME GoRuntime
		COMMAND ${PYEXEC} "run_test.py"
		WORKING_DIRECTORY "${CMAKE_CURRENT_LIST_DIR}/test")