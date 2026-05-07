package tidy

import (
	"bytes"
	"fmt"
	"os"
	"text/template"

	"github.com/Voltamon/Uca/internal/config"
	"github.com/Voltamon/Uca/internal/schema"
	"github.com/Voltamon/Uca/internal/templates"
)

func generateUIPackage(cfg *config.Config) error {
	// Generate base UI exports
	data, err := templates.FS.ReadFile("uca/ui/index.ts")
	if err != nil {
		return fmt.Errorf("failed to read ui template: %w", err)
	}

	err = os.MkdirAll(".uca/ui", 0755)
	if err != nil {
		return err
	}

	err = os.WriteFile(".uca/ui/index.ts", data, 0644)
	if err != nil {
		return err
	}

	// Generate srv module
	srvTmplData, err := templates.FS.ReadFile("uca/ui/srv.ts.tmpl")
	if err != nil {
		return fmt.Errorf("failed to read srv template: %w", err)
	}

	srvTmpl, err := template.New("srv").Parse(string(srvTmplData))
	if err != nil {
		return fmt.Errorf("failed to parse srv template: %w", err)
	}

	var srvBuf bytes.Buffer
	err = srvTmpl.Execute(&srvBuf, cfg)
	if err != nil {
		return fmt.Errorf("failed to execute srv template: %w", err)
	}

	err = os.MkdirAll(".uca/ui/srv", 0755)
	if err != nil {
		return err
	}

	err = os.WriteFile(".uca/ui/srv/index.ts", srvBuf.Bytes(), 0644)
	if err != nil {
		return err
	}

	// Generate ai module
	aiTmplData, err := templates.FS.ReadFile("uca/ui/ai.ts.tmpl")
	if err != nil {
		return fmt.Errorf("failed to read ai template: %w", err)
	}

	aiTmpl, err := template.New("ai").Parse(string(aiTmplData))
	if err != nil {
		return fmt.Errorf("failed to parse ai template: %w", err)
	}

	var aiBuf bytes.Buffer
	err = aiTmpl.Execute(&aiBuf, cfg)
	if err != nil {
		return fmt.Errorf("failed to execute ai template: %w", err)
	}

	err = os.MkdirAll(".uca/ui/ai", 0755)
	if err != nil {
		return err
	}

	err = os.WriteFile(".uca/ui/ai/index.ts", aiBuf.Bytes(), 0644)
	if err != nil {
		return err
	}

	// Generate TypeScript types
	typesTmplData, err := templates.FS.ReadFile("uca/ui/types.ts.tmpl")
	if err != nil {
		return fmt.Errorf("failed to read types template: %w", err)
	}

	typesTmpl, err := template.New("types").Parse(string(typesTmplData))
	if err != nil {
		return fmt.Errorf("failed to parse types template: %w", err)
	}

	schemaObj := schema.ParseFromConfig(cfg)
	var typesBuf bytes.Buffer
	err = typesTmpl.Execute(&typesBuf, schemaObj)
	if err != nil {
		return fmt.Errorf("failed to execute types template: %w", err)
	}

	err = os.WriteFile(".uca/ui/types.ts", typesBuf.Bytes(), 0644)
	if err != nil {
		return err
	}

	// Update index.ts to export types and services
	// The exports are now handled in the template itself
	fmt.Println("Generated: .uca/ui/index.ts, .uca/ui/srv/index.ts and .uca/ui/ai/index.ts")
	return nil
}
