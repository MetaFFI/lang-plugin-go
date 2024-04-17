import os.path
import subprocess

from colorama import Fore

compile_command = ['javac', '*.java']
print(f'{Fore.BLUE}Running - {" ".join(compile_command)}{Fore.RESET}')
subprocess.run(compile_command, check=True, cwd=os.path.dirname(__file__))
