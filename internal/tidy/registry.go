package tidy

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"text/template"

	"github.com/Voltamon/Uca/internal/config"
	"github.com/Voltamon/Uca/internal/templates"
)

func generateRegistry(cfg *config.Config) error {
	implemented, err := findImplementedFunctions()
	if err != nil {
		return fmt.Errorf("failed to find implemented functions: %w", err)
	}

	tmplData, err := templates.FS.ReadFile("uca/registry.go.tmpl")
	if err != nil {
		return fmt.Errorf("failed to read registry template: %w", err)
	}

	tmpl, err := template.New("registry").Parse(string(tmplData))
	if err != nil {
		return fmt.Errorf("failed to parse registry template: %w", err)
	}

	type AgentEntry struct {
		Name string
		Port string
	}

	agents := make([]AgentEntry, len(cfg.Agents))
	for i, a := range cfg.Agents {
		agents[i] = AgentEntry{
			Name: a.Name,
			Port: fmt.Sprintf("%d", cfg.App.Port.AI),
		}
	}

	type MethodInfo struct {
		Name     string
		IsCustom bool
	}

	type ServiceInfo struct {
		Name    string
		Methods []MethodInfo
	}

	services := make([]ServiceInfo, len(cfg.Services))
	for i, s := range cfg.Services {
		services[i] = ServiceInfo{
			Name: s.Name,
		}
		for _, m := range s.Methods {
			funcName := s.Name + m
			services[i].Methods = append(services[i].Methods, MethodInfo{
				Name:     m,
				IsCustom: implemented[funcName],
			})
		}
	}

	data := struct {
		AppName  string
		Services []ServiceInfo
		Agents   []AgentEntry
	}{
		AppName:  cfg.App.Name,
		Services: services,
		Agents:   agents,
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, data)
	if err != nil {
		return fmt.Errorf("failed to execute registry template: %w", err)
	}

	err = os.MkdirAll(".uca/registry", 0755)
	if err != nil {
		return fmt.Errorf("failed to create .uca/registry directory: %w", err)
	}

	// Remove old registry file if it exists in the wrong place
	os.Remove(".uca/uca/registry.go")

	return os.WriteFile(".uca/registry/registry.go", buf.Bytes(), 0644)
}

func findImplementedFunctions() (map[string]bool, error) {
	implemented := make(map[string]bool)

	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, "services", func(fi os.FileInfo) bool {
		return filepath.Ext(fi.Name()) == ".go"
	}, 0)

	if err != nil {
		if os.IsNotExist(err) {
			return implemented, nil
		}
		return nil, err
	}

	for _, pkg := range pkgs {
		for _, file := range pkg.Files {
			for _, decl := range file.Decls {
				if fn, ok := decl.(*ast.FuncDecl); ok {
					if fn.Recv == nil { // Global function
						implemented[fn.Name.Name] = true
					}
				}
			}
		}
	}

	return implemented, nil
}
