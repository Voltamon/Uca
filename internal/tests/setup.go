package tests

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/Voltamon/Uca/internal/runtime"
	"github.com/Voltamon/Uca/internal/scaffold"
)

func EnsureTestDeps(goOnly bool, tsOnly bool, pyOnly bool) error {
	runAll := !goOnly && !tsOnly && !pyOnly

	if runAll || tsOnly {
		err := ensureVitest()
		if err != nil {
			return err
		}
		err = ensureVitestConfig()
		if err != nil {
			return err
		}
	}

	if runAll || pyOnly {
		err := ensurePytest()
		if err != nil {
			return err
		}
	}

	return nil
}

func ensureVitest() error {
	vitestBin := filepath.Join(".uca", "node_modules", ".bin", "vitest")
	if _, err := os.Stat(vitestBin); err == nil {
		return nil
	}

	err := runtime.EnsureNode()
	if err != nil {
		return fmt.Errorf("failed to ensure node runtime: %w", err)
	}

	fmt.Println("Installing vitest...")

	absNodeBin, err := filepath.Abs(runtime.NodeBin())
	if err != nil {
		return err
	}

	cmd := exec.Command(absNodeBin, runtime.NpmBin(), "install", "--save-dev", "vitest", "@testing-library/preact", "jsdom", "--prefix", ".uca")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func ensurePytest() error {
	pytestBin := ".uca/venv/bin/pytest"
	if _, err := os.Stat(pytestBin); err == nil {
		return nil
	}

	err := runtime.EnsurePython()
	if err != nil {
		return fmt.Errorf("failed to ensure python runtime: %w", err)
	}

	fmt.Println("Installing pytest...")
	cmd := exec.Command(".uca/venv/bin/pip", "install", "--quiet", "pytest")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func ensureVitestConfig() error {
	err := scaffold.CopyTemplate("uca/vitest.config.ts", ".uca/vitest.config.ts", scaffold.TemplateVars{})
	if err != nil {
		return err
	}
	return scaffold.CopyTemplate("tests/vitest.setup.ts", ".uca/vitest.setup.ts", scaffold.TemplateVars{})
}
