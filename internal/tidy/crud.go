package tidy

import (
	"bytes"
	"fmt"
	"os"
	"text/template"

	"github.com/Voltamon/Uca/internal/templates"
)

func generateCRUDPackage(appName string) error {
	tmplData, err := templates.FS.ReadFile("uca/crud.go.tmpl")
	if err != nil {
		return fmt.Errorf("failed to read crud template: %w", err)
	}

	tmpl, err := template.New("crud").Parse(string(tmplData))
	if err != nil {
		return fmt.Errorf("failed to parse crud template: %w", err)
	}

	data := struct {
		AppName string
	}{
		AppName: appName,
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, data)
	if err != nil {
		return fmt.Errorf("failed to execute crud template: %w", err)
	}

	err = os.MkdirAll(".uca/uca", 0755)
	if err != nil {
		return fmt.Errorf("failed to create .uca/uca directory: %w", err)
	}

	return os.WriteFile(".uca/uca/crud.go", buf.Bytes(), 0644)
}
