# python script to run unittests for api using subprocess
import shutil
import subprocess
import os
import sys
from colorama import init, Fore
import platform
import glob

# Initialize colorama
init()

# Get the current path of this Python script
current_path = os.path.dirname(os.path.abspath(__file__))

def get_extension_by_platform() -> str:
	if platform.system() == 'Windows':
		return '.dll'
	elif platform.system() == 'Darwin':
		return '.dylib'
	else:
		return '.so'


def run_script(script_path):
	print(f'{Fore.CYAN}Running script: {script_path}{Fore.RESET}')
	
	if script_path.endswith('.py'):
		# Python script
		python_command = 'py' if platform.system() == 'Windows' else 'python3.11'
		command = [python_command, script_path]
	else:
		raise ValueError(f'Unsupported script file type: {script_path}')
	
	script_dir = os.path.dirname(os.path.abspath(script_path))
	
	process = subprocess.Popen(command, cwd=script_dir, stdout=subprocess.PIPE, stderr=subprocess.PIPE, text=True)
	
	stdout, stderr = process.communicate()
	print(stdout, end='')
	print(stderr, file=sys.stderr, end='')
	
	if process.returncode != 0:
		raise subprocess.CalledProcessError(process.returncode, command)

def get_terminal_after_path_symbol() -> str:
	if platform.system() == 'Windows':
		return '>'
	else:
		return '$'

def run_unittest(script_path):
	print(f'{Fore.CYAN}Running unittest: {script_path}{Fore.RESET}')
	
	terminal_path = os.path.dirname(script_path)+get_terminal_after_path_symbol()

	if script_path.endswith('.py'):
		# Python unittest
		python_command = 'py' if platform.system() == 'Windows' else 'python3.11'
		command = [python_command, '-m', 'unittest', script_path]
	elif script_path.endswith('.java'):
		# Java JUnit test
		junit_jar = os.path.join(current_path, 'junit-platform-console-standalone-1.10.2.jar')
		hamcrest_jar = os.path.join(current_path, 'hamcrest-core-1.3.jar')
		bridge_jar = os.path.join(os.environ['METAFFI_HOME']+'/openjdk/', 'xllr.openjdk.bridge.jar')
		api_jar = os.path.join(os.environ['METAFFI_HOME']+'/openjdk/', 'metaffi.api.jar')
		class_name = os.path.splitext(os.path.basename(script_path))[0]
		class_path = f'.{os.pathsep}{junit_jar}{os.pathsep}{hamcrest_jar}{os.pathsep}{bridge_jar}{os.pathsep}{api_jar}'
		
		# Compile the Java source file
		compile_command = ['javac', '-cp', class_path, script_path]
		print(f'{Fore.BLUE}{os.getcwd()+get_terminal_after_path_symbol()} - {" ".join(compile_command)}{Fore.RESET}')
		subprocess.run(compile_command, check=True)
		
		# Run the JUnit test
		command = ['java', '-jar', junit_jar, '-cp', class_path, '-c', class_name]
	elif script_path.endswith('.go'):
		
		goget_command = ['go', 'get', '-v']
		print(f'{Fore.BLUE}{terminal_path} - {" ".join(goget_command)}{Fore.RESET}')
		subprocess.run(goget_command, check=True, cwd=os.path.dirname(script_path))

		goget_command = ['go', 'get', '-v', 'github.com/MetaFFI/plugin-sdk@main']
		print(f'{Fore.BLUE}{terminal_path} - {" ".join(goget_command)}{Fore.RESET}')
		subprocess.run(goget_command, check=True, cwd=os.path.dirname(script_path))

		goget_command = ['go', 'get', '-v', 'github.com/MetaFFI/lang-plugin-go/compiler@main']
		print(f'{Fore.BLUE}{terminal_path} - {" ".join(goget_command)}{Fore.RESET}')
		subprocess.run(goget_command, check=True, cwd=os.path.dirname(script_path))

		goget_command = ['go', 'get', '-v', 'github.com/MetaFFI/lang-plugin-go/go-runtime@main']
		print(f'{Fore.BLUE}{terminal_path} - {" ".join(goget_command)}{Fore.RESET}')
		subprocess.run(goget_command, check=True, cwd=os.path.dirname(script_path))

		goget_command = ['go', 'get', '-v', 'github.com/MetaFFI/lang-plugin-go/api@main']
		print(f'{Fore.BLUE}{terminal_path} - {" ".join(goget_command)}{Fore.RESET}')
		subprocess.run(goget_command, check=True, cwd=os.path.dirname(script_path))
		
		command = ['go', 'run', script_path]
	else:
		raise ValueError(f'Unsupported unittest file type: {script_path}')
	
	script_dir = os.path.dirname(os.path.abspath(script_path))
	
	print(f'{Fore.BLUE}{terminal_path} - {" ".join(command)}{Fore.RESET}')
	process = subprocess.Popen(command, cwd=script_dir, stdout=subprocess.PIPE, stderr=subprocess.PIPE, text=True)
	
	stdout, stderr = process.communicate()
	
	print(stdout, end='')
	print(stderr, file=sys.stderr, end='')
	
	if process.returncode != 0:
		raise subprocess.CalledProcessError(process.returncode, command)
	
	# If it's a Java unittest, delete the compiled .class file
	if script_path.endswith('.java'):
		class_files = glob.glob(os.path.join(script_dir, "*.class"))
		for file in class_files:
			os.remove(file)


# --------------------------------------------

# sanity tests

# --------------------------------------------

# run Go -> python3.11 tests
print(f'{Fore.MAGENTA}Testing Sanity Go -> Python3.11{Fore.RESET} - {Fore.YELLOW}RUNNING{Fore.RESET}')

# Define the paths to the scripts to be run
test_sanity_python311_path = os.path.join(current_path, 'sanity', 'python3', 'MetaFFIAPI_python3.go')

# Run the scripts
run_unittest(test_sanity_python311_path)

print(f'{Fore.MAGENTA}Testing Sanity Go -> Python3.11{Fore.RESET} - {Fore.GREEN}PASSED{Fore.RESET}')

# --------------------------------------------

# run Go -> OpenJDK tests
if platform.system() != 'Windows':
	print(f'{Fore.MAGENTA}Testing Sanity Go -> OpenJDK{Fore.RESET} - {Fore.YELLOW}RUNNING{Fore.RESET}')
	
	# Define the paths to the scripts to be run
	build_sanity_go_script_path = os.path.join(current_path, 'sanity', 'openjdk', 'compile_java.py')
	
	run_script(build_sanity_go_script_path)
	
	if os.path.exists(os.path.join(current_path, 'sanity', 'openjdk', 'TestRuntime.class')):
		dest_dir = os.path.join(current_path, 'sanity', 'openjdk', 'sanity')
		if not os.path.exists(dest_dir):
			os.makedirs(dest_dir)
		
		src_file = os.path.join(current_path, 'sanity', 'openjdk', 'TestRuntime.class')
		dest_file = os.path.join(dest_dir, 'TestRuntime.class')
		shutil.move(src_file, dest_file)
		
		src_file = os.path.join(current_path, 'sanity', 'openjdk', 'TestMap.class')
		dest_file = os.path.join(dest_dir, 'TestMap.class')
		shutil.move(src_file, dest_file)
		
	test_sanity_go_path = os.path.join(current_path, 'sanity', 'openjdk', 'MetaFFIAPI_openjdk.go')
	
	
	run_unittest(test_sanity_go_path)
	
	print(f'{Fore.MAGENTA}Testing Sanity Go -> OpenJDK{Fore.RESET} - {Fore.GREEN}PASSED{Fore.RESET}')
else:
	print(f'{Fore.MAGENTA}Testing Sanity Go -> OpenJDK{Fore.RESET} - {Fore.RED}SKIPPING due to Go bug in loading JVM in windows (https://github.com/golang/go/issues/58542){Fore.RESET}')

# --------------------------------------------

# extended tests

# --------------------------------------------

# run Go -> python3.11 tests

print(f'{Fore.MAGENTA}Testing Extended Go -> Python3.11{Fore.RESET} - {Fore.YELLOW}RUNNING{Fore.RESET}')

# Define the path to the unittest script
test_extended_bs4_path = os.path.join(current_path, 'extended', 'python3', 'beautifulsoup', 'BeautifulSoupTest.go')
run_unittest(test_extended_bs4_path)

test_extended_py_complex_primitives_path = os.path.join(current_path, 'extended', 'python3', 'complex-primitives', 'ComplexPrimitivesTest.go')
run_unittest(test_extended_py_complex_primitives_path)

print(f'{Fore.MAGENTA}Testing Extended Go -> Python3.11{Fore.RESET} - {Fore.GREEN}PASSED{Fore.RESET}')

# --------------------------------------------

# run Go -> OpenJDK tests
if platform.system() != 'Windows':
	print(f'{Fore.MAGENTA}Testing Extended Go -> OpenJDK{Fore.RESET} - {Fore.YELLOW}RUNNING{Fore.RESET}')
	
	# Define the paths to the scripts to be run
	test_extended_openjdk_bytes_arrays_path = os.path.join(current_path, 'extended', 'openjdk', 'log4j', 'Log4j.go')
	run_unittest(test_extended_openjdk_bytes_arrays_path)
	
	print(f'{Fore.MAGENTA}Testing Extended Go -> OpenJDK{Fore.RESET} - {Fore.GREEN}PASSED{Fore.RESET}')
else:
	print(f'{Fore.MAGENTA}Testing Sanity Go -> OpenJDK{Fore.RESET} - {Fore.RED}SKIPPING due to Go bug in loading JVM in windows (https://github.com/golang/go/issues/58542){Fore.RESET}')