// Package test provides E2E tests for the Go host compiler.
// Orchestration: run metaffi to generate host code, prepare output/test with go.mod and entity tests, run go test there.
package test

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
)

const (
	idlRelPath     = "sdk/test_modules/guest_modules/test/xllr.test.idl.json"
	pluginDir      = "test"
	pluginNameWin  = "xllr.test.dll"
	pluginNameUnix = "libxllr.test.so"
	pluginNameMac  = "libxllr.test.dylib"
	outputDir      = "output"
	outputTestDir  = "output/test"
	// Go host compiler writes outputDir/moduleName/idl_source_dots_to_underscore_MetaFFIHost.go
	generatedFile = "xllr_test_MetaFFIHost.go"
)

func requireEnv(t *testing.T, name string) string {
	t.Helper()
	v := os.Getenv(name)
	if v == "" {
		t.Fatalf("%s must be set (test fails if unset)", name)
	}
	return v
}

func getPluginName() string {
	switch runtime.GOOS {
	case "windows":
		return pluginNameWin
	case "darwin":
		return pluginNameMac
	default:
		return pluginNameUnix
	}
}

// TestHostCompilerE2E runs the full E2E: generate host code, prepare output/test, run go test there.
func TestHostCompilerE2E(t *testing.T) {
	sourceRoot := requireEnv(t, "METAFFI_SOURCE_ROOT")
	metaffiHome := requireEnv(t, "METAFFI_HOME")

	idlPath := filepath.Join(sourceRoot, filepath.FromSlash(idlRelPath))
	if _, err := os.Stat(idlPath); err != nil {
		t.Fatalf("IDL file missing: %s: %v", idlPath, err)
	}

	pluginPath := filepath.Join(metaffiHome, pluginDir, getPluginName())
	if _, err := os.Stat(pluginPath); err != nil {
		t.Fatalf("xllr.test plugin missing: %s: %v", pluginPath, err)
	}

	testDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}

	// Create output directory (clean if exists)
	outputBase := filepath.Join(testDir, outputDir)
	_ = os.RemoveAll(outputBase)
	if err := os.MkdirAll(outputBase, 0755); err != nil {
		t.Fatalf("mkdir output: %v", err)
	}

	// Run metaffi -c --idl <idl> -h go (output goes under output_base, creates output/test/)
	cmd := exec.Command("metaffi", "-c", "--idl", idlPath, "-h", "go")
	cmd.Dir = outputBase
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("metaffi compiler failed: %v\n%s", err, out)
	}

	generatedPath := filepath.Join(testDir, outputTestDir, generatedFile)
	if _, err := os.Stat(generatedPath); err != nil {
		t.Fatalf("generated file missing after compile: %s: %v", generatedPath, err)
	}

	// Prepare output/test for go test: go.mod + copy static entity test
	outTestDir := filepath.Join(testDir, outputTestDir)
	goModPath := filepath.Join(outTestDir, "go.mod")
	goModContent := "module test\n\n" +
		"go 1.21\n\n" +
		"require (\n" +
		"\tgithub.com/MetaFFI/sdk/api/go v0.0.0\n" +
		"\tgithub.com/MetaFFI/sdk/idl_entities/go v0.0.0\n" +
		")\n\n" +
		"replace (\n" +
		"\tgithub.com/MetaFFI/sdk/api/go => " + filepath.ToSlash(sourceRoot) + "/sdk/api/go\n" +
		"\tgithub.com/MetaFFI/sdk/idl_entities/go => " + filepath.ToSlash(sourceRoot) + "/sdk/idl_entities/go\n" +
		")\n"
	if err := os.WriteFile(goModPath, []byte(goModContent), 0644); err != nil {
		t.Fatalf("write go.mod: %v", err)
	}

	// Copy static entity test into output/test (strip //go:build ignore so it builds there)
	staticEntityTest := filepath.Join(testDir, "static", "e2e_entity_test.go")
	entityTestBytes, err := os.ReadFile(staticEntityTest)
	if err != nil {
		t.Fatalf("read e2e_entity_test.go: %v", err)
	}
	// Remove the leading "//go:build ignore\n" line so the file is built in output/test
	const buildIgnore = "//go:build ignore\n"
	if len(entityTestBytes) >= len(buildIgnore) && string(entityTestBytes[:len(buildIgnore)]) == buildIgnore {
		entityTestBytes = entityTestBytes[len(buildIgnore):]
	}
	destEntityTest := filepath.Join(outTestDir, "e2e_entity_test.go")
	if err := os.WriteFile(destEntityTest, entityTestBytes, 0644); err != nil {
		t.Fatalf("write e2e_entity_test.go: %v", err)
	}

	// Run go test from a temp dir outside the repo so Go uses our go.mod (not the parent compiler module).
	runDir, err := os.MkdirTemp("", "go_host_compiler_e2e_")
	if err != nil {
		t.Fatalf("mkdir temp: %v", err)
	}
	defer os.RemoveAll(runDir)
	// Copy files so the temp dir is a self-contained module
	for _, name := range []string{generatedFile, "go.mod", "e2e_entity_test.go"} {
		src := filepath.Join(outTestDir, name)
		dst := filepath.Join(runDir, name)
		data, err := os.ReadFile(src)
		if err != nil {
			t.Fatalf("read %s: %v", name, err)
		}
		if err := os.WriteFile(dst, data, 0644); err != nil {
			t.Fatalf("write %s: %v", name, err)
		}
	}

	// go.mod may need tidying (versions/transitive deps); run in temp dir so we don't touch repo
	tidy := exec.Command("go", "mod", "tidy")
	tidy.Dir = runDir
	tidy.Env = os.Environ()
	if out, err := tidy.CombinedOutput(); err != nil {
		t.Fatalf("go mod tidy in temp dir: %v\n%s", err, out)
	}

	runTest := exec.Command("go", "test", "-v", "-count=1", ".")
	runTest.Dir = runDir
	runTest.Env = append(os.Environ(), "METAFFI_TEST_PLUGIN_PATH="+pluginPath)
	runTest.Stdout = os.Stdout
	runTest.Stderr = os.Stderr
	if err := runTest.Run(); err != nil {
		// Dump generated file for debugging
		genData, _ := os.ReadFile(filepath.Join(runDir, generatedFile))
		lines := splitLines(string(genData))
		snippet := ""
		// Dump lines around errors
		for _, ln := range []int{1047, 1048, 1049, 1050, 1051, 1052, 1053, 1054, 1055, 1056, 1057, 1058, 1059, 1060, 1065, 1066, 1067, 1068, 1069, 1070, 1071, 1072, 1073, 1074, 1075, 1076, 1080, 1081, 1082, 1083, 1084, 1085, 1086, 1087, 1088, 1089, 1090, 1091, 1092, 1093, 1094, 1095, 1096, 1097, 1098, 1099, 1100, 1101, 1102, 1103, 1104, 1105} {
			if ln > 0 && ln <= len(lines) {
				snippet += fmt.Sprintf("  L%d: %s\n", ln, lines[ln-1])
			}
		}
		t.Fatalf("go test in temp dir failed: %v\nError lines in %s:\n%s", err, generatedFile, snippet)
	}
}

func splitLines(s string) []string {
	var lines []string
	start := 0
	for i := 0; i < len(s); i++ {
		if s[i] == '\n' {
			lines = append(lines, s[start:i])
			start = i + 1
		}
	}
	if start < len(s) {
		lines = append(lines, s[start:])
	}
	return lines
}
