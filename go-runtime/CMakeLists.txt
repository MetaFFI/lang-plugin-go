
#------------------------------- Go Language Support - Go-Runtime package ------------------------------------------

# Add METAFFI_HOME environment variable to "CGO_CFLAGS", so MetaFFIGoHostCommon.go can compile its CGo code

if("$ENV{CGO_CFLAGS}" STREQUAL "")
	message(WARNING "CGO_CFLAGS must be set to METAFFI_HOME to use the Go plugin")
	SET(ENV{CGO_CFLAGS} $ENV{METAFFI_HOME})
endif()