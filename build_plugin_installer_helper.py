from genericpath import isdir
import platform
import shutil
from typing import List, Tuple, Dict
import glob
import os
import shutil

gomods_to_restore = []

def get_files(win_metaffi_home: str, ubuntu_metaffi_home: str) -> Tuple[Dict[str, str], Dict[str, str]]:
	global gomods_to_restore

	pluginname = 'go'
	
	win_metaffi_home = win_metaffi_home.replace('\\', '/')+f'/{pluginname}/'
	ubuntu_metaffi_home = ubuntu_metaffi_home.replace('\\', '/')+f'/{pluginname}/'

	# make a backup of go.mod files, as we strip the "replace" directives from them
	for gomod in glob.glob(f'{win_metaffi_home}/**/go.mod', recursive=True):
		# make a copy
		shutil.copy(gomod, gomod+'.bak')

		gomods_to_restore.append(gomod)
		
		with open(gomod, 'r') as f:
			lines = f.readlines()
		
		with open(gomod, 'w') as f:
			for line in lines:
				if not line.startswith('replace'):
					f.write(line)

	win_files = {
		'xllr.go.dll': win_metaffi_home + 'xllr.go.dll',
		'metaffi.idl.go.dll': win_metaffi_home + 'metaffi.idl.go.dll',
		'metaffi.compiler.go.dll': win_metaffi_home + 'metaffi.compiler.go.dll'
	}

	# for each absolute path in the value of win_files, check if the file exists
	for key, value in win_files.items():
		if not os.path.isfile(value):
			raise FileNotFoundError(f'{value} not found - cannot build the installer')
		

	ubuntu_files = {
		'xllr.go.so': ubuntu_metaffi_home + 'xllr.go.so',
		'metaffi.idl.go.so': ubuntu_metaffi_home + 'metaffi.idl.go.so',
		'metaffi.compiler.go.so': ubuntu_metaffi_home + 'metaffi.compiler.go.so',
		'libboost_filesystem.so.1.87.0': ubuntu_metaffi_home + 'libboost_filesystem.so.1.87.0'
	}

	# * copy the api tests
	current_script_dir = os.path.dirname(os.path.abspath(__file__))
	api_tests_files = glob.glob(f'{current_script_dir}/api/tests/**', recursive=True)
	for file in api_tests_files:
		if '__pycache__' in file:
			continue

		if os.path.isfile(file):
			target = file.replace('\\', '/').removeprefix(current_script_dir.replace('\\', '/')+'/api/')
			win_files[target] = file
			ubuntu_files[target] = file

	# * uninstaller
	win_files['uninstall_plugin.py'] = os.path.dirname(os.path.abspath(__file__))+'/uninstall_plugin.py'
	ubuntu_files['uninstall_plugin.py'] = os.path.dirname(os.path.abspath(__file__))+'/uninstall_plugin.py'

	return win_files, ubuntu_files


def post_copy_files():
	global gomods_to_restore

	# restore the backup of go.mod files
	for gomodbak in gomods_to_restore:
		# force move
		shutil.move(gomodbak+'.bak', gomodbak)


def setup_environment():
	# make sure pycrosskit is installed
	try:
		# Attempt to import pycrosskit
		import pycrosskit
	except ImportError:
		print("pycrosskit for writing environment variables is missing")
		print("make sure to install requirements.txt and try again")
		exit(1)
			
	
	from pycrosskit.envariables import SysEnv

	# Set CGO_ENABLED=1
	if platform.system() == 'Windows':
		res = os.system('go env -w CGO_ENABLED=1')
		if res != 0:
			raise ValueError('Failed to set CGO_ENABLED=1')
	else:
		SysEnv().set('CGO_ENABLED', '1')

	# Get CGO_CFLAGS from the STDOUT of "go env CGO_CFLAGS" and append -I{metaffi_home}
	metaffi_home = os.environ.get('METAFFI_HOME')
	if metaffi_home is None:
		raise ValueError('METAFFI_HOME is not set')
	
	cgo_cflags = os.popen('go env CGO_CFLAGS').read().strip()
	cgo_cflags += f' -I"{metaffi_home}"'

	# Set CGO_CFLAGS
	if platform.system() == 'Windows':
		res = os.system(f'go env -w CGO_CFLAGS="{cgo_cflags}"')
		if res != 0:
			raise ValueError('Failed to set CGO_CFLAGS')
	else:
		SysEnv().set('CGO_CFLAGS', cgo_cflags)


def check_prerequisites() -> bool:
	# run "go version" and make sure it returns 0
	if os.system('go version') != 0:
		print('Go is not installed')
		return False
	
	return True


def print_prerequisites():
	pass


def get_version():
	return '0.3.0'

