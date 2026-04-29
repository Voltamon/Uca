package tidy

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
	"github.com/Voltamon/Uca/internal/config"
	"github.com/Voltamon/Uca/internal/scaffold"
	"github.com/Voltamon/Uca/internal/runtime"
	"github.com/Voltamon/Uca/internal/schema"
	"github.com/Voltamon/Uca/internal/typegen"
	"github.com/Voltamon/Uca/internal/auth"
	"github.com/Voltamon/Uca/internal/manifest"
)

func Run() (*config.Config, error) {
	data, err := os.ReadFile("uca.yaml")
	if err != nil {
		return nil, fmt.Errorf("uca.yaml not found — are you inside a uca project? %w", err)
	}

	var cfg config.Config
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to parse uca.yaml: %w", err)
	}

	err = validateConfig(&cfg)
	if err != nil {
		return nil, fmt.Errorf("invalid uca.yaml: %w", err)
	}

	fmt.Println("uca.yaml is valid")

	err = reconcileRoles(&cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to reconcile roles: %w", err)
	}

	err = reconcileFiles(&cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to reconcile files: %w", err)
	}

	fmt.Println("reconciliation complete")

	err = generateTypes()
	if err != nil {
		return nil, fmt.Errorf("failed to generate types: %w", err)
	}

	err = ensureGoMod(cfg.App.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to ensure go.mod: %w", err)
	}

	err = generateRolesPackage()
	if err != nil {
		return nil, fmt.Errorf("failed to generate roles package: %w", err)
	}

	err = runtime.EnsureAll()
	if err != nil {
		return nil, fmt.Errorf("failed to ensure runtimes: %w", err)
	}

	err = generateFrontend(&cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to generate frontend: %w", err)
	}

	err = reconcileSchema(&cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to reconcile schema: %w", err)
	}

	err = generateRegistry(&cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to generate registry: %w", err)
	}

	fmt.Println("Generated: .uca/uca/registry.go")

	err = generateUcaPackage(cfg.App.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to generate uca package: %w", err)
	}

	err = generateAgentServer(&cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to generate agent server: %w", err)
	}

	err = generateDevMain(cfg.App.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to generate dev main: %w", err)
	}

	fmt.Println("Generated: .uca/main.go")

	return &cfg, nil
}

func reconcileSchema(cfg *config.Config) error {
	desired := schema.ParseFromConfig(cfg)

	rolesCfg, err := auth.Load()
	if err != nil {
		return fmt.Errorf("failed to load roles: %w", err)
	}

	desired.Roles = rolesCfg.Roles
	desired.DefaultRole = auth.DefaultRole

	current, err := schema.LoadSnapshot()
	if err != nil {
		return fmt.Errorf("failed to load schema snapshot: %w", err)
	}

	changes := schema.Diff(current, desired)

	accepted, err := schema.ApplyChanges(changes)
	if err != nil {
		return fmt.Errorf("failed to apply schema changes: %w", err)
	}

	if !accepted {
		return nil
	}

	return schema.SaveSnapshot(desired)
}

func reconcileFiles(cfg *config.Config) error {
	rolesCfg, _ := auth.Load()
	defaultRole := auth.DefaultRole
	if len(rolesCfg.Roles) > 0 {
		defaultRole = rolesCfg.Roles[0]
	}

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
				Name:        service.Name,
				DefaultRole: defaultRole,
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

	if cfg.App.Port.Frontend == 0 {
		cfg.App.Port.Frontend = 5173
	}
	if cfg.App.Port.Backend == 0 {
		cfg.App.Port.Backend = 8090
	}
	if cfg.App.Port.AI == 0 {
		cfg.App.Port.AI = 8091
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

func reconcileRoles(cfg *config.Config) error {
	roles := []string{}
	for _, page := range cfg.Pages {
		if page.Role != "" {
			roles = append(roles, page.Role)
		}
	}

	err := auth.SyncFromConfig(roles)
	if err != nil {
		return err
	}

	err = auth.EnsureDefaultRole()
	if err != nil {
		return err
	}

	rolesCfg, err := auth.Load()
	if err != nil {
		return err
	}

	if len(rolesCfg.Roles) == 0 {
		return nil
	}

	for i, service := range cfg.Services {
		if service.Name != "User" {
			continue
		}
		if cfg.Services[i].Schema == nil {
			cfg.Services[i].Schema = make(map[string]string)
		}
		roleType := "select:" + strings.Join(rolesCfg.Roles, ",") + " | required"
		cfg.Services[i].Schema["role"] = roleType
		break
	}

	return manifest.SaveConfig(cfg)
}

func generateTypes() error {
    return typegen.GenerateTypes("services", ".uca/types")
}
