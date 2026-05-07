package tidy

import (
	"bytes"
	"fmt"
	"os"
	"text/template"

	"github.com/Voltamon/Uca/internal/config"
	"github.com/Voltamon/Uca/internal/scaffold"
	"github.com/Voltamon/Uca/internal/templates"
)

func generateAgentPythonPackage(cfg *config.Config) error {
	// 1. Create __init__.py
	err := os.MkdirAll(".uca/uca", 0755)
	if err != nil {
		return err
	}
	err = os.WriteFile(".uca/uca/__init__.py", []byte("# Uca Python Package\n"), 0644)
	if err != nil {
		return err
	}

	// 2. Generate ai.py
	err = scaffold.CopyTemplate("uca/python/ai.py", ".uca/uca/ai.py", scaffold.TemplateVars{})
	if err != nil {
		return fmt.Errorf("failed to generate ai.py: %w", err)
	}

	// 3. Generate srv.py from template
	tmplData, err := templates.FS.ReadFile("uca/python/srv.py.tmpl")
	if err != nil {
		return fmt.Errorf("failed to read srv template: %w", err)
	}

	tmpl, err := template.New("srv").Parse(string(tmplData))
	if err != nil {
		return fmt.Errorf("failed to parse srv template: %w", err)
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, cfg)
	if err != nil {
		return fmt.Errorf("failed to execute srv template: %w", err)
	}

	err = os.WriteFile(".uca/uca/srv.py", buf.Bytes(), 0644)
	if err != nil {
		return err
	}

	fmt.Println("Generated: .uca/uca/ai.py and .uca/uca/srv.py")
	return nil
}
