"""
Plugin hooks for the MetaFFI Go plugin.

Invoked by the CLI installer:
  python plugin_hooks.py --check-prerequisites
  python plugin_hooks.py --setup-environment
  python plugin_hooks.py --pre-uninstall
"""

import os
import platform
import sys


def check_prerequisites() -> bool:
	"""Return True if prerequisites met. Print message and return False if not."""

	if os.system('go version') != 0:
		print('Go is not installed')
		return False

	return True


def setup_environment():
	"""Called after files are installed. Set env vars, etc."""

	try:
		import pycrosskit
	except ImportError:
		print("pycrosskit for writing environment variables is missing")
		print("make sure to install requirements.txt and try again")
		sys.exit(1)

	from pycrosskit.envariables import SysEnv

	# Set CGO_ENABLED=1
	if platform.system() == 'Windows':
		res = os.system('go env -w CGO_ENABLED=1')
		if res != 0:
			raise ValueError('Failed to set CGO_ENABLED=1')
	else:
		SysEnv().set('CGO_ENABLED', '1')

	# Append -I"{METAFFI_HOME}" to CGO_CFLAGS
	metaffi_home = os.environ.get('METAFFI_HOME')
	if metaffi_home is None:
		raise ValueError('METAFFI_HOME is not set')

	cgo_cflags = os.popen('go env CGO_CFLAGS').read().strip()
	if f'-I"{metaffi_home}"' not in cgo_cflags:
		cgo_cflags += f' -I"{metaffi_home}"'

	if platform.system() == 'Windows':
		res = os.system(f'go env -w CGO_CFLAGS="{cgo_cflags}"')
		if res != 0:
			raise ValueError('Failed to set CGO_CFLAGS')
	else:
		SysEnv().set('CGO_CFLAGS', cgo_cflags)


def pre_uninstall():
	"""Called before plugin directory is removed. Clean up env vars, etc."""

	try:
		import pycrosskit
	except ImportError:
		# pycrosskit not available — nothing to clean up
		return

	from pycrosskit.envariables import SysEnv

	metaffi_home = os.environ.get('METAFFI_HOME')
	if metaffi_home is None:
		return

	# Remove -I"{METAFFI_HOME}" from CGO_CFLAGS
	cgo_cflags = os.popen('go env CGO_CFLAGS').read().strip()
	needle = f' -I"{metaffi_home}"'
	if needle in cgo_cflags:
		cgo_cflags = cgo_cflags.replace(needle, '')
		if platform.system() == 'Windows':
			os.system(f'go env -w CGO_CFLAGS="{cgo_cflags}"')
		else:
			SysEnv().set('CGO_CFLAGS', cgo_cflags)


if __name__ == "__main__":
	if len(sys.argv) < 2:
		print("Usage: python plugin_hooks.py --check-prerequisites|--setup-environment|--pre-uninstall")
		sys.exit(1)

	action = sys.argv[1]

	if action == '--check-prerequisites':
		ok = check_prerequisites()
		sys.exit(0 if ok else 1)

	elif action == '--setup-environment':
		setup_environment()
		sys.exit(0)

	elif action == '--pre-uninstall':
		pre_uninstall()
		sys.exit(0)

	else:
		print(f"Unknown action: {action}")
		sys.exit(1)
