package tests

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/Voltamon/Uca/internal/runtime"
	"github.com/Voltamon/Uca/internal/scaffold"
)

func Run(goOnly bool, tsOnly bool, pyOnly bool) error {
	err := ensureTestDeps()
	if err != nil {
		return fmt.Errorf("failed to ensure test deps: %w", err)
	}

	runAll := !goOnly && !tsOnly && !pyOnly

	var goResult, tsResult, pyResult *TestResult

	if runAll || goOnly {
		goResult = runGoTests()
	}

	if runAll || tsOnly {
		tsResult = runTSTests()
	}

	if runAll || pyOnly {
		pyResult = runPyTests()
	}

	fmt.Println()
	printSummary(goResult, tsResult, pyResult)

	if hasFailures(goResult, tsResult, pyResult) {
		return fmt.Errorf("some tests failed")
	}

	return nil
}

type TestResult struct {
	Name   string
	Passed bool
	Output string
}

func runGoTests() *TestResult {
	fmt.Println("\nRunning Go tests...")
	cmd := exec.Command("go", "test", "./services/...", "-v")
	cmd.Dir = ".uca"
	output, err := cmd.CombinedOutput()

	result := &TestResult{
		Name:   "Go",
		Output: string(output),
		Passed: err == nil,
	}

	fmt.Print(string(output))
	return result
}

func runTSTests() *TestResult {
	fmt.Println("\nRunning TypeScript tests...")

	absNodeBin, err := filepath.Abs(runtime.NodeBin())
	if err != nil {
		return &TestResult{Name: "TypeScript", Passed: false, Output: err.Error()}
	}

	cmd := exec.Command(absNodeBin, "node_modules/.bin/vitest", "run")
	cmd.Dir = ".uca"
	cmd.Env = append(os.Environ())
	output, err := cmd.CombinedOutput()

	result := &TestResult{
		Name:   "TypeScript",
		Output: string(output),
		Passed: err == nil,
	}

	fmt.Print(string(output))
	return result
}

func runPyTests() *TestResult {
	fmt.Println("\nRunning Python tests...")
	cmd := exec.Command(".uca/venv/bin/pytest", "agents/", "-v")
	output, err := cmd.CombinedOutput()

	result := &TestResult{
		Name:   "Python",
		Output: string(output),
		Passed: err == nil,
	}

	fmt.Print(string(output))
	return result
}

func printSummary(results ...*TestResult) {
	fmt.Println("─── Test Summary ───")
	for _, r := range results {
		if r == nil {
			continue
		}
		status := "✓"
		if !r.Passed {
			status = "✗"
		}
		fmt.Printf("  %s %s\n", status, r.Name)
	}
}

func hasFailures(results ...*TestResult) bool {
	for _, r := range results {
		if r != nil && !r.Passed {
			return true
		}
	}
	return false
}

func AddTestStub(name string, testType string) error {
	vars := scaffold.TemplateVars{
		Name: name,
	}

	switch testType {
	case "go":
		dest := "services/" + name + "_test.go"
		if _, err := os.Stat(dest); err == nil {
			return fmt.Errorf("test file %s already exists", dest)
		}
		err := scaffold.CopyTemplate("tests/service_test.go", dest, vars)
		if err != nil {
			return err
		}
		fmt.Println("Created:", dest)
	case "ts":
		dest := "pages/" + name + ".test.tsx"
		if _, err := os.Stat(dest); err == nil {
			return fmt.Errorf("test file %s already exists", dest)
		}
		err := scaffold.CopyTemplate("tests/page_test.tsx", dest, vars)
		if err != nil {
			return err
		}
		fmt.Println("Created:", dest)
	case "py":
		dest := "agents/" + name + "_test.py"
		if _, err := os.Stat(dest); err == nil {
			return fmt.Errorf("test file %s already exists", dest)
		}
		err := scaffold.CopyTemplate("tests/agent_test.py", dest, vars)
		if err != nil {
			return err
		}
		fmt.Println("Created:", dest)
	default:
		return fmt.Errorf("unknown test type %q (use: go, ts, py)", testType)
	}

	return nil
}
