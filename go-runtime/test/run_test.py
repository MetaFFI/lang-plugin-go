import subprocess
import os

build_process = subprocess.run(["go", "mod", "tidy"])
if build_process.returncode != 0:
	print("Failed to build the Go executable.")
	exit(1)

build_process = subprocess.run(["go", "build", "-o", "testexec"])
if build_process.returncode != 0:
	print("Failed to build the Go executable.")
	exit(1)

# Run the Go executable
run_process = subprocess.run(["./testexec"])
if run_process.returncode != 0:
	print("Failed to run the Go executable.")
	exit(1)

# Delete the Go executable
try:
	os.remove("testexec")
except OSError as e:
	print("Error: %s : %s" % ("testexec", e.strerror))
	exit(1)

print('go-runtime unit-test ran successfully')