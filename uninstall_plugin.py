import importlib
import sys
import subprocess

import os

if __name__ == "__main__":
	# remove -I"{metaffi_home}" from CGO_CFLAGS
	
	# make sure pycrosskit is installed
	try:
		# Attempt to import pycrosskit
		import pycrosskit
	except ImportError:
		print("pycrosskit for writing environment variables is missing")
		print("make sure to install requirements.txt and try again")
		exit(1)
	
	from pycrosskit.envariables import SysEnv

	# get CGO_CFLAGS
	cgo_cflags = os.popen('go env CGO_CFLAGS').read().strip()
	
	# remove -I"{metaffi_home}" from CGO_CFLAGS
	metaffi_home = os.environ.get('METAFFI_HOME')
	if metaffi_home is None:
		raise ValueError('METAFFI_HOME is not set')

	if f'-I"{metaffi_home}"' in cgo_cflags:
		cgo_cflags = cgo_cflags.replace(f' -I"{metaffi_home}"', '')
		SysEnv().set('CGO_CFLAGS', cgo_cflags)

	# remove the plugins directory
	import shutil
	import os
	
	# get the path to the plugins directory
	metaffi_home = os.getenv('METAFFI_HOME')
	assert metaffi_home is not None, 'METAFFI_HOME is not set'
	
	plugins_dir = os.path.join(metaffi_home, 'go')
	
	# remove the plugins directory
	shutil.rmtree(plugins_dir)



		
