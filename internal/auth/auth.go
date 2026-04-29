package auth

import (
	"encoding/json"
	"fmt"
	"os"
)

const rolesPath = ".uca/roles.json"

const DefaultRole = "user"

type RolesConfig struct {
	Roles []string `json:"roles"`
}

func readRoles() (*RolesConfig, error) {
	cfg := &RolesConfig{}

	data, err := os.ReadFile(rolesPath)
	if err != nil {
		if os.IsNotExist(err) {
			return cfg, nil
		}
		return nil, err
	}

	err = json.Unmarshal(data, cfg)
	return cfg, err
}

func writeRoles(cfg *RolesConfig) error {
	err := os.MkdirAll(".uca", 0755)
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(rolesPath, data, 0644)
}

func EnsureDefaultRole() error {
	cfg, err := readRoles()
	if err != nil {
		return err
	}

	for _, r := range cfg.Roles {
		if r == DefaultRole {
			return nil
		}
	}

	cfg.Roles = append([]string{DefaultRole}, cfg.Roles...)
	return writeRoles(cfg)
}

func AddRole(role string) error {
	if role == DefaultRole {
		return fmt.Errorf("role %q is the default and always exists", role)
	}

	cfg, err := readRoles()
	if err != nil {
		return err
	}

	for _, r := range cfg.Roles {
		if r == role {
			return fmt.Errorf("role %q already exists", role)
		}
	}

	cfg.Roles = append(cfg.Roles, role)
	return writeRoles(cfg)
}

func RemoveRole(role string) error {
	if role == DefaultRole {
		return fmt.Errorf("role %q is the default and cannot be removed", role)
	}

	cfg, err := readRoles()
	if err != nil {
		return err
	}

	newRoles := []string{}
	found := false
	for _, r := range cfg.Roles {
		if r == role {
			found = true
			continue
		}
		newRoles = append(newRoles, r)
	}

	if !found {
		return fmt.Errorf("role %q not found", role)
	}

	cfg.Roles = newRoles
	return writeRoles(cfg)
}

func ListRoles() ([]string, error) {
	cfg, err := readRoles()
	if err != nil {
		return nil, err
	}
	return cfg.Roles, nil
}

func Load() (*RolesConfig, error) {
	return readRoles()
}

func SyncFromConfig(pages []string) error {
	cfg, err := readRoles()
	if err != nil {
		return err
	}

	existing := make(map[string]bool)
	for _, r := range cfg.Roles {
		existing[r] = true
	}

	changed := false
	for _, role := range pages {
		if role == "" {
			continue
		}
		if !existing[role] {
			cfg.Roles = append(cfg.Roles, role)
			existing[role] = true
			changed = true
			fmt.Printf("Auto-registered role %q from uca.yaml\n", role)
		}
	}

	if !existing[DefaultRole] {
		cfg.Roles = append([]string{DefaultRole}, cfg.Roles...)
		changed = true
	}

	if changed {
		return writeRoles(cfg)
	}
	return nil
}
