package tidy

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"text/template"

	"github.com/Voltamon/Uca/internal/config"
	"github.com/Voltamon/Uca/internal/scaffold"
	"github.com/Voltamon/Uca/internal/templates"
)

func generateFrontend(cfg *config.Config) error {
	err := generateMainTsx(cfg)
	if err != nil {
		return err
	}

	files := []struct {
		src  string
		dest string
	}{
		{"uca/package.json", ".uca/package.json"},
		{"uca/vite.config.mts", ".uca/vite.config.mts"},
		{"uca/index.html", ".uca/index.html"},
	}

	for _, f := range files {
		err := scaffold.CopyTemplate(f.src, f.dest, scaffold.TemplateVars{
			AppName:      cfg.App.Name,
			BackendPort:  fmt.Sprintf("%d", cfg.App.Port.Backend),
			FrontendPort: fmt.Sprintf("%d", cfg.App.Port.Frontend),
			AIPort:       fmt.Sprintf("%d", cfg.App.Port.AI),
		})
		if err != nil {
			return fmt.Errorf("failed to generate %s: %w", f.dest, err)
		}
		fmt.Println("Generated:", f.dest)
	}

	err = installFrontendDeps()
	if err != nil {
		return err
	}

	err = installPythonDeps()
	if err != nil {
		return err
	}

	return nil
}

func generateMainTsx(cfg *config.Config) error {
	tmplData, err := templates.FS.ReadFile("uca/main.tsx.tmpl")
	if err != nil {
		return fmt.Errorf("failed to read main.tsx template: %w", err)
	}

	tmpl, err := template.New("main.tsx").Parse(string(tmplData))
	if err != nil {
		return fmt.Errorf("failed to parse main.tsx template: %w", err)
	}

	data := struct {
		Pages []config.PageConfig
	}{
		Pages: cfg.Pages,
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, data)
	if err != nil {
		return fmt.Errorf("failed to execute main.tsx template: %w", err)
	}

	return os.WriteFile(".uca/main.tsx", buf.Bytes(), 0644)
}

func installFrontendDeps() error {
	if _, err := os.Stat(".uca/node_modules"); err == nil {
		return nil
	}

	fmt.Println("Installing frontend dependencies...")
	cmd := exec.Command("npm", "install", "--silent")
	cmd.Dir = ".uca"
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func installPythonDeps() error {
	if _, err := os.Stat(".uca/venv"); err == nil {
		return nil
	}

	fmt.Println("Creating Python virtual environment...")
	cmd := exec.Command("python3", "-m", "venv", "venv")
	cmd.Dir = ".uca"
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to create venv: %w", err)
	}

	fmt.Println("Installing Python dependencies...")
	cmd = exec.Command("venv/bin/pip", "install", "--quiet", "smolagents[litellm]", "httpx", "litellm")
	cmd.Dir = ".uca"
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to install python deps: %w", err)
	}

	fmt.Println("Python dependencies installed")
	return nil
}
