package tests

import (
	"os"
	"testing"
)

func TestAddTestStub(t *testing.T) {
	// Create a temporary directory for the test
	tmpDir, err := os.MkdirTemp("", "uca-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Change working directory to temp dir
	oldWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get current wd: %v", err)
	}
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("failed to change wd: %v", err)
	}
	defer os.Chdir(oldWd)

	// Test Go stub
	err = AddTestStub("MyService", "go")
	if err != nil {
		t.Errorf("failed to add go test stub: %v", err)
	}
	goFile := "services/MyService_test.go"
	if _, err := os.Stat(goFile); os.IsNotExist(err) {
		t.Errorf("go test stub file was not created: %s", goFile)
	}

	// Test TS stub
	err = AddTestStub("MyPage", "ts")
	if err != nil {
		t.Errorf("failed to add ts test stub: %v", err)
	}
	tsFile := "pages/MyPage.test.tsx"
	if _, err := os.Stat(tsFile); os.IsNotExist(err) {
		t.Errorf("ts test stub file was not created: %s", tsFile)
	}

	// Test Python stub
	err = AddTestStub("MyAgent", "py")
	if err != nil {
		t.Errorf("failed to add py test stub: %v", err)
	}
	pyFile := "agents/MyAgent_test.py"
	if _, err := os.Stat(pyFile); os.IsNotExist(err) {
		t.Errorf("py test stub file was not created: %s", pyFile)
	}

	// Test duplicate error
	err = AddTestStub("MyService", "go")
	if err == nil {
		t.Error("expected error when adding duplicate go test stub, but got nil")
	}

	// Test unknown type
	err = AddTestStub("Other", "invalid")
	if err == nil {
		t.Error("expected error when adding test stub with invalid type, but got nil")
	}
}

func TestRun_NoProject(t *testing.T) {
	// Create a temporary directory for the test
	tmpDir, err := os.MkdirTemp("", "uca-test-run-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Change working directory to temp dir
	oldWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get current wd: %v", err)
	}
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("failed to change wd: %v", err)
	}
	defer os.Chdir(oldWd)

	// tests.Run() should fail because uca.yaml is missing (via tidy.Run())
	err = Run(false, false, false)
	if err == nil {
		t.Error("expected error when running tests in a non-uca project, but got nil")
	}
}

func TestRun_Go(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	// Create a temporary directory for the test
	tmpDir, err := os.MkdirTemp("", "uca-test-go-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Change working directory to temp dir
	oldWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get current wd: %v", err)
	}
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("failed to change wd: %v", err)
	}
	defer os.Chdir(oldWd)

	// Setup a minimal uca.yaml
	ucaYaml := `app:
  name: testapp
  version: 0.1.0
`
	if err := os.WriteFile("uca.yaml", []byte(ucaYaml), 0644); err != nil {
		t.Fatalf("failed to write uca.yaml: %v", err)
	}

	// Add a Go test stub
	if err := AddTestStub("Ping", "go"); err != nil {
		t.Fatalf("failed to add go test stub: %v", err)
	}

	// Modify the stub to be a simple passing test
	goTestContent := `package services

import "testing"

func TestPing(t *testing.T) {
	// Simple passing test
}
`
	if err := os.WriteFile("services/Ping_test.go", []byte(goTestContent), 0644); err != nil {
		t.Fatalf("failed to update go test stub: %v", err)
	}

	// Run Go tests
	// Note: This will call tidy.Run(), which might take time if it downloads Node/Python.
	// In a real environment, we might want to skip those if they are not needed.
	// But since we are testing 'uca test', we want to see it work.
	err = Run(true, false, false)
	if err != nil {
		t.Errorf("expected tests to pass, but got error: %v", err)
	}
}

func TestRun_GoFailure(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	// Create a temporary directory for the test
	tmpDir, err := os.MkdirTemp("", "uca-test-go-fail-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Change working directory to temp dir
	oldWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get current wd: %v", err)
	}
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("failed to change wd: %v", err)
	}
	defer os.Chdir(oldWd)

	// Setup a minimal uca.yaml
	ucaYaml := `app:
  name: failapp
  version: 0.1.0
`
	if err := os.WriteFile("uca.yaml", []byte(ucaYaml), 0644); err != nil {
		t.Fatalf("failed to write uca.yaml: %v", err)
	}

	// Add a failing Go test
	if err := os.MkdirAll("services", 0755); err != nil {
		t.Fatalf("failed to create services dir: %v", err)
	}
	goTestContent := `package services

import "testing"

func TestFail(t *testing.T) {
	t.Fatal("intended failure")
}
`
	if err := os.WriteFile("services/fail_test.go", []byte(goTestContent), 0644); err != nil {
		t.Fatalf("failed to write failing go test: %v", err)
	}

	// Run Go tests
	err = Run(true, false, false)
	if err == nil {
		t.Error("expected tests to fail, but got no error")
	}
}
