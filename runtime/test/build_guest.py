import os
import sys

rc = os.system('metaffi -c --idl TestRuntime.go -g go')
if rc != 0:
    print(f'build_guest.py: metaffi command failed with exit code {rc}', file=sys.stderr)
    sys.exit(1)
