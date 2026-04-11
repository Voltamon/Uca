package tidy

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/Voltamon/Uca/internal/scaffold"
)

func generateDevMain(appName string) error {
	return scaffold.CopyTemplate("uca/main.go", ".uca/main.go", scaffold.TemplateVars{
		AppName: appName,
	})
}

func ensureGoMod(appName string) error {
	err := os.MkdirAll(".uca", 0755)
	if err != nil {
		return fmt.Errorf("failed to create .uca directory: %w", err)
	}

	if _, err := os.Stat(".uca/go.mod"); err == nil {
		return nil
	}

	cmd := exec.Command("go", "mod", "init", appName)
	cmd.Dir = ".uca"
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to initialize go.mod: %w", err)
	}

	fmt.Println("Generated: .uca/go.mod")

	symlinks := []struct {
		src  string
		dest string
	}{
		{"../services", ".uca/services"},
		{"../agents", ".uca/agents"},
		{"../pages", ".uca/pages"},
		{"../modules", ".uca/modules"},
	}

	for _, s := range symlinks {
		if _, err := os.Lstat(s.dest); err == nil {
			continue
		}
		err := os.Symlink(s.src, s.dest)
		if err != nil {
			return fmt.Errorf("failed to create symlink %s: %w", s.dest, err)
		}
		fmt.Println("Linked:", s.dest)
	}

	return nil
}
