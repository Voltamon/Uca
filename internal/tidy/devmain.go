package tidy

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/Voltamon/Uca/internal/config"
	"github.com/Voltamon/Uca/internal/scaffold"
)

func generateDevMain(cfg *config.Config) error {
	return scaffold.CopyTemplate("uca/main.go", ".uca/main.go", scaffold.TemplateVars{
		AppName:     cfg.App.Name,
		BackendPort: fmt.Sprintf("%d", cfg.App.Port.Backend),
		AIPort:      fmt.Sprintf("%d", cfg.App.Port.AI),
	})
}

func generateUcaPackage(appName string) error {
	files := []struct {
		src  string
		dest string
	}{
		{"uca/migrations.go", ".uca/uca/migrations.go"},
		{"uca/types.go", ".uca/uca/types.go"},
	}

	for _, f := range files {
		err := scaffold.CopyTemplate(f.src, f.dest, scaffold.TemplateVars{
			AppName: appName,
		})
		if err != nil {
			return fmt.Errorf("failed to generate %s: %w", f.dest, err)
		}
	}

	fmt.Println("Generated: .uca/uca/migrations.go")
	fmt.Println("Generated: .uca/uca/types.go")
	return nil
}

func ensureGoMod(appName string) error {
	err := os.MkdirAll(".uca", 0755)
	if err != nil {
		return fmt.Errorf("failed to create .uca directory: %w", err)
	}

	if _, err := os.Stat(".uca/go.mod"); os.IsNotExist(err) {
		cmd := exec.Command("go", "mod", "init", appName)
		cmd.Dir = ".uca"
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err = cmd.Run()
		if err != nil {
			return fmt.Errorf("failed to initialize go.mod: %w", err)
		}

		cmd = exec.Command("go", "get", "github.com/pocketbase/pocketbase@v0.22.11")
		cmd.Dir = ".uca"
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err = cmd.Run()
		if err != nil {
			return fmt.Errorf("failed to get pocketbase: %w", err)
		}

		fmt.Println("Generated: .uca/go.mod")
	}

	cmd := exec.Command("go", "mod", "tidy")
	cmd.Dir = ".uca"
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to run go mod tidy: %w", err)
	}

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
