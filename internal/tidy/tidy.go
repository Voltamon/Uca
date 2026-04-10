package tidy

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
	"github.com/Voltamon/Uca/internal/config"
	"github.com/Voltamon/Uca/internal/scaffold"
)

func Run() error {
	data, err := os.ReadFile("uca.yaml")
	if err != nil {
		return fmt.Errorf("uca.yaml not found — are you inside a uca project? %w", err)
	}

	var cfg config.Config
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return fmt.Errorf("failed to parse uca.yaml: %w", err)
	}

	err = validateConfig(&cfg)
	if err != nil {
		return fmt.Errorf("invalid uca.yaml: %w", err)
	}

	fmt.Println("uca.yaml is valid")

	err = reconcileFiles(&cfg)
	if err != nil {
		return fmt.Errorf("failed to reconcile files: %w", err)
	}

	fmt.Println("reconciliation complete")
	return nil
}

func reconcileFiles(cfg *config.Config) error {
	for _, page := range cfg.Pages {
		path := "pages/" + page.Name + ".tsx"
		if _, err := os.Stat(path); os.IsNotExist(err) {
			err = scaffold.CopyTemplate("pages/default.tsx", path, scaffold.TemplateVars{
				Name: page.Name,
			})
			if err != nil {
				return fmt.Errorf("failed to create %s: %w", path, err)
			}
			fmt.Println("Created:", path)
		}
	}

	for _, service := range cfg.Services {
		path := "services/" + service.Name + ".go"
		if _, err := os.Stat(path); os.IsNotExist(err) {
			err = scaffold.CopyTemplate("services/default.go", path, scaffold.TemplateVars{
				Name: service.Name,
			})
			if err != nil {
				return fmt.Errorf("failed to create %s: %w", path, err)
			}
			fmt.Println("Created:", path)
		}
	}

	for _, agent := range cfg.Agents {
		path := "agents/" + agent.Name + ".py"
		if _, err := os.Stat(path); os.IsNotExist(err) {
			err = scaffold.CopyTemplate("agents/default.py", path, scaffold.TemplateVars{
				Name:  agent.Name,
				Model: agent.Model,
			})
			if err != nil {
				return fmt.Errorf("failed to create %s: %w", path, err)
			}
			fmt.Println("Created:", path)
		}
	}

	return nil
}

func validateConfig(cfg *config.Config) error {
	if cfg.App.Name == "" {
		return fmt.Errorf("app.name is required")
	}

	if cfg.App.Version == "" {
		return fmt.Errorf("app.version is required")
	}

	for _, page := range cfg.Pages {
		if page.Name == "" {
			return fmt.Errorf("every page must have a name")
		}
		if page.Route == "" {
			return fmt.Errorf("page %q must have a route", page.Name)
		}
	}

	for _, service := range cfg.Services {
		if service.Name == "" {
			return fmt.Errorf("every service must have a name")
		}
		if len(service.Methods) == 0 {
			return fmt.Errorf("service %q must have at least one method", service.Name)
		}
	}

	for _, agent := range cfg.Agents {
		if agent.Name == "" {
			return fmt.Errorf("every agent must have a name")
		}
		if agent.Model == "" {
			return fmt.Errorf("agent %q must have a model", agent.Name)
		}
	}

	return nil
}
