package supervisor

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/Voltamon/Uca/internal/tidy"
)

func Start() error {
	fmt.Println("Running tidy...")
	err := tidy.Run()
	if err != nil {
		return fmt.Errorf("tidy failed: %w", err)
	}

	fmt.Println("Resolving dependencies...")
	err = runGoModTidy()
	if err != nil {
		return fmt.Errorf("dependency resolution failed: %w", err)
	}

	fmt.Println("Building project...")
	err = buildProject()
	if err != nil {
		return fmt.Errorf("build failed: %w", err)
	}

	fmt.Println("Starting dev server...")
	return runProject()
}

func runGoModTidy() error {
	cmd := exec.Command("go", "mod", "tidy")
	cmd.Dir = ".uca"
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func buildProject() error {
	cmd := exec.Command("go", "build", "-o", "server")
	cmd.Dir = ".uca"
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func runProject() error {
	cmd := exec.Command(".uca/server")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
