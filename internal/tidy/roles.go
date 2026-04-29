package tidy

import (
	"fmt"
	"os"
	"strings"

	"github.com/Voltamon/Uca/internal/auth"
)

func generateRolesPackage() error {
	rolesCfg, err := auth.Load()
	if err != nil {
		return err
	}

	if len(rolesCfg.Roles) == 0 {
		return nil
	}

	err = os.MkdirAll(".uca/uca/roles", 0755)
	if err != nil {
		return err
	}

	err = generateGoRoles(rolesCfg.Roles)
	if err != nil {
		return err
	}

	err = os.MkdirAll(".uca/roles", 0755)
	if err != nil {
		return err
	}

	err = generateTSRoles(rolesCfg.Roles)
	if err != nil {
		return err
	}

	return nil
}

func generateGoRoles(roles []string) error {
	var sb strings.Builder

	sb.WriteString("package roles\n\n")
	sb.WriteString("const (\n")

	for _, role := range roles {
		constName := strings.Title(strings.ToLower(role))
		sb.WriteString(fmt.Sprintf("\t%s = %q\n", constName, role))
	}

	sb.WriteString(")\n\n")
	sb.WriteString("var All = []string{\n")

	for _, role := range roles {
		sb.WriteString(fmt.Sprintf("\t%q,\n", role))
	}

	sb.WriteString("}\n\n")
	sb.WriteString(fmt.Sprintf("const Default = %q\n", auth.DefaultRole))

	return os.WriteFile(".uca/uca/roles/roles.go", []byte(sb.String()), 0644)
}

func generateTSRoles(roles []string) error {
	var sb strings.Builder

	sb.WriteString("export const roles = {\n")

	for _, role := range roles {
		constName := strings.ToLower(role)
		sb.WriteString(fmt.Sprintf("\t%s: %q,\n", constName, role))
	}

	sb.WriteString("} as const\n\n")
	sb.WriteString("export type Role = typeof roles[keyof typeof roles]\n\n")
	sb.WriteString(fmt.Sprintf("export const defaultRole = %q\n", auth.DefaultRole))

	return os.WriteFile(".uca/roles/index.ts", []byte(sb.String()), 0644)
}
