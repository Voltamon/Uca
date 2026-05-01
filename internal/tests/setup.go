package tests

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/Voltamon/Uca/internal/runtime"
	"github.com/Voltamon/Uca/internal/scaffold"
)

func ensureTestDeps() error {
	err := ensureVitest()
	if err != nil {
		return err
	}

	err = ensurePytest()
	if err != nil {
		return err
	}

	err = ensureVitestConfig()
	if err != nil {
		return err
	}

	return nil
}

func ensureVitest() error {
	vitestBin := filepath.Join(".uca", "node_modules", ".bin", "vitest")
	if _, err := os.Stat(vitestBin); err == nil {
		return nil
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

	fmt.Println("Installing pytest...")
	cmd := exec.Command(".uca/venv/bin/pip", "install", "--quiet", "pytest")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func ensureVitestConfig() error {
	return scaffold.CopyTemplate("uca/vitest.config.ts", ".uca/vitest.config.ts", scaffold.TemplateVars{})
}
