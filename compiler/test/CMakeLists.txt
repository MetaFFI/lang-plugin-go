# E2E Tests for Go Host Compiler
#
# Validates the Go host compiler by:
# 1. Running the metaffi compiler to generate host code from xllr.test IDL
# 2. Building the generated code and running go test against xllr.test
#
# Requires: metaffi.compiler.go, metaffi, xllr.test (build before running tests).
# Environment: METAFFI_SOURCE_ROOT, METAFFI_HOME must be set (test fails if not).

add_go_test(go_host_compiler_e2e_test
	WORKING_DIRECTORY ${CMAKE_CURRENT_LIST_DIR}
)

set(go_host_compiler_e2e_test go_host_compiler_e2e_test PARENT_SCOPE)
