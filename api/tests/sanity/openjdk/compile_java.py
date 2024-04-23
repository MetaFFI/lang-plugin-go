import glob
import os.path
import subprocess

from colorama import Fore

compile_command = ['javac']
compile_command.extend(glob.glob(os.path.join(os.path.dirname(__file__), '*.java')))
print(f'{Fore.BLUE}Running - {" ".join(compile_command)}{Fore.RESET}')
subprocess.run(compile_command, check=True, cwd=os.path.dirname(__file__))
