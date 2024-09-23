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
	
	win_metaffi_home = win_metaffi_home.replace('\\', '/')
	ubuntu_metaffi_home = ubuntu_metaffi_home.replace('\\', '/')

	# make a backup of go.mod files, as we strip the "replace" directives from them
	for gomod in glob.glob(f'{win_metaffi_home}/go/**/go.mod', recursive=True):
		# make a copy
		shutil.copy(gomod, gomod+'.bak')

		gomods_to_restore.append(gomod)
		
		with open(gomod, 'r') as f:
			lines = f.readlines()
		
		with open(gomod, 'w') as f:
			for line in lines:
				if not line.startswith('replace'):
					f.write(line)

	win_files = {}
	for file in glob.glob(win_metaffi_home + f'/{pluginname}/**', recursive=True):		
		if os.path.isfile(file) and '__' not in file:
			file = file.replace('\\', '/')
			win_files[file.removeprefix(win_metaffi_home+f'/{pluginname}/')] = file

	assert len(win_files) > 0, f'No files found in {win_metaffi_home}/{pluginname}'

	ubuntu_files = {}
	for file in glob.glob(ubuntu_metaffi_home + f'/{pluginname}/**', recursive=True):
		if os.path.isfile(file) and '__' not in file:
			file = file.replace('\\', '/')
			ubuntu_files[file.removeprefix(ubuntu_metaffi_home+f'/{pluginname}/')] = file

	assert len(ubuntu_files) > 0, f'No files found in {ubuntu_metaffi_home}/{pluginname}'

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
	print("""Prerequisites:\n\tGo""")

