package tidy

import (
	"bytes"
	"fmt"
	"os"
	"text/template"

	"github.com/Voltamon/Uca/internal/config"
	"github.com/Voltamon/Uca/internal/templates"
)

func generateRegistry(cfg *config.Config) error {
	tmplData, err := templates.FS.ReadFile("uca/registry.go.tmpl")
	if err != nil {
		return fmt.Errorf("failed to read registry template: %w", err)
	}

	tmpl, err := template.New("registry").Parse(string(tmplData))
	if err != nil {
		return fmt.Errorf("failed to parse registry template: %w", err)
	}

	data := struct {
		AppName  string
		Services []config.ServiceConfig
	}{
		AppName:  cfg.App.Name,
		Services: cfg.Services,
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, data)
	if err != nil {
		return fmt.Errorf("failed to execute registry template: %w", err)
	}

	err = os.MkdirAll(".uca/uca", 0755)
	if err != nil {
		return fmt.Errorf("failed to create .uca/uca directory: %w", err)
	}

	return os.WriteFile(".uca/uca/registry.go", buf.Bytes(), 0644)
}
