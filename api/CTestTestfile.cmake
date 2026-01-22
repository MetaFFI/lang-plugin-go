# CMake generated Testfile for 
# Source directory: C:/src/github.com/MetaFFI/lang-plugin-go/api
# Build directory: C:/src/github.com/MetaFFI/lang-plugin-go/api
# 
# This file includes the relevant testing commands required for 
# testing this directory and lists subdirectories to be tested as well.
if(CTEST_CONFIGURATION_TYPE MATCHES "^([Dd][Ee][Bb][Uu][Gg])$")
  add_test([=[go_api_cross_pl_tests]=] "C:/Users/green/AppData/Local/Programs/Python/Python311/python3.exe" "run_api_tests.py")
  set_tests_properties([=[go_api_cross_pl_tests]=] PROPERTIES  WORKING_DIRECTORY "C:/src/github.com/MetaFFI/lang-plugin-go/api/tests" _BACKTRACE_TRIPLES "C:/src/github.com/MetaFFI/lang-plugin-go/api/CMakeLists.txt;7;add_test;C:/src/github.com/MetaFFI/lang-plugin-go/api/CMakeLists.txt;0;")
elseif(CTEST_CONFIGURATION_TYPE MATCHES "^([Rr][Ee][Ll][Ee][Aa][Ss][Ee])$")
  add_test([=[go_api_cross_pl_tests]=] "C:/Users/green/AppData/Local/Programs/Python/Python311/python3.exe" "run_api_tests.py")
  set_tests_properties([=[go_api_cross_pl_tests]=] PROPERTIES  WORKING_DIRECTORY "C:/src/github.com/MetaFFI/lang-plugin-go/api/tests" _BACKTRACE_TRIPLES "C:/src/github.com/MetaFFI/lang-plugin-go/api/CMakeLists.txt;7;add_test;C:/src/github.com/MetaFFI/lang-plugin-go/api/CMakeLists.txt;0;")
elseif(CTEST_CONFIGURATION_TYPE MATCHES "^([Mm][Ii][Nn][Ss][Ii][Zz][Ee][Rr][Ee][Ll])$")
  add_test([=[go_api_cross_pl_tests]=] "C:/Users/green/AppData/Local/Programs/Python/Python311/python3.exe" "run_api_tests.py")
  set_tests_properties([=[go_api_cross_pl_tests]=] PROPERTIES  WORKING_DIRECTORY "C:/src/github.com/MetaFFI/lang-plugin-go/api/tests" _BACKTRACE_TRIPLES "C:/src/github.com/MetaFFI/lang-plugin-go/api/CMakeLists.txt;7;add_test;C:/src/github.com/MetaFFI/lang-plugin-go/api/CMakeLists.txt;0;")
elseif(CTEST_CONFIGURATION_TYPE MATCHES "^([Rr][Ee][Ll][Ww][Ii][Tt][Hh][Dd][Ee][Bb][Ii][Nn][Ff][Oo])$")
  add_test([=[go_api_cross_pl_tests]=] "C:/Users/green/AppData/Local/Programs/Python/Python311/python3.exe" "run_api_tests.py")
  set_tests_properties([=[go_api_cross_pl_tests]=] PROPERTIES  WORKING_DIRECTORY "C:/src/github.com/MetaFFI/lang-plugin-go/api/tests" _BACKTRACE_TRIPLES "C:/src/github.com/MetaFFI/lang-plugin-go/api/CMakeLists.txt;7;add_test;C:/src/github.com/MetaFFI/lang-plugin-go/api/CMakeLists.txt;0;")
else()
  add_test([=[go_api_cross_pl_tests]=] NOT_AVAILABLE)
endif()
