package scaffold

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/Voltamon/Uca/internal/templates"
)

type TemplateVars struct {
	AppName     string
	Name        string
	Model       string
	AIPort      string
	BackendPort string
	FrontendPort string
	DefaultRole string
}

func CopyTemplate(src string, dest string, vars TemplateVars) error {
	data, err := templates.FS.ReadFile(src)
	if err != nil {
		return err
	}

	content := string(data)
	content = strings.ReplaceAll(content, "{{APP_NAME}}", vars.AppName)
	content = strings.ReplaceAll(content, "{{NAME}}", vars.Name)
	content = strings.ReplaceAll(content, "{{MODEL}}", vars.Model)
	content = strings.ReplaceAll(content, "{{AI_PORT}}", vars.AIPort)
	content = strings.ReplaceAll(content, "{{BACKEND_PORT}}", vars.BackendPort)
	content = strings.ReplaceAll(content, "{{FRONTEND_PORT}}", vars.FrontendPort)
	content = strings.ReplaceAll(content, "{{DEFAULT_ROLE}}", vars.DefaultRole)

	err = os.MkdirAll(filepath.Dir(dest), 0755)
	if err != nil {
		return err
	}

	return os.WriteFile(dest, []byte(content), 0644)
}

func GenerateFiles(appName string, model string, aiPort string, defaultRole string) error {
	vars := TemplateVars{
		AppName:     appName,
		Model:       model,
		AIPort:      aiPort,
		DefaultRole: defaultRole,
	}

	files := []struct {
		src  string
		dest string
	}{
		{"pages/Welcome.tsx", appName + "/pages/Welcome.tsx"},
		{"pages/Chat.tsx", appName + "/pages/Chat.tsx"},
		{"services/User.go", appName + "/services/User.go"},
		{"services/History.go", appName + "/services/History.go"},
		{"agents/Assistant.py", appName + "/agents/Assistant.py"},
	}

	for _, f := range files {
		err := CopyTemplate(f.src, f.dest, vars)
		if err != nil {
			return err
		}
	}

	return nil
}
