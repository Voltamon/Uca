package scaffold

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/Voltamon/Uca/internal/templates"
)

type TemplateVars struct {
	AppName string
	Name    string
	Model   string
	AIPort  string
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

	err = os.MkdirAll(filepath.Dir(dest), 0755)
	if err != nil {
		return err
	}

	return os.WriteFile(dest, []byte(content), 0644)
}

func GenerateFiles(appName string, model string, aiPort string) error {
	vars := TemplateVars{
		AppName: appName,
		Model:   model,
		AIPort:  aiPort,
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
