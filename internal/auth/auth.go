package auth

import (
	"encoding/json"
	"fmt"
	"os"
)

const rolesPath = ".uca/roles.json"

type RolesConfig struct {
	Roles   []string `json:"roles"`
	Default string   `json:"default"`
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

func AddRole(role string, isDefault bool) error {
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

	if isDefault || cfg.Default == "" {
		cfg.Default = role
	}

	return writeRoles(cfg)
}

func RemoveRole(role string) error {
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

	if cfg.Default == role {
		if len(newRoles) > 0 {
			cfg.Default = newRoles[0]
		} else {
			cfg.Default = ""
		}
	}

	return writeRoles(cfg)
}

func ListRoles() ([]string, error) {
	cfg, err := readRoles()
	if err != nil {
		return nil, err
	}
	return cfg.Roles, nil
}

func GetDefault() (string, error) {
	cfg, err := readRoles()
	if err != nil {
		return "", err
	}
	return cfg.Default, nil
}

func Load() (*RolesConfig, error) {
	return readRoles()
}
