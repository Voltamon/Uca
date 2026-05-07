package tidy

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"text/template"

	"github.com/Voltamon/Uca/internal/config"
	"github.com/Voltamon/Uca/internal/templates"
)

func generateAutoTests(cfg *config.Config) error {
	err := os.MkdirAll(".uca/tests/autogen/pages", 0755)
	if err != nil {
		return err
	}
	err = os.MkdirAll(".uca/tests/autogen/services", 0755)
	if err != nil {
		return err
	}
	err = os.MkdirAll(".uca/tests/autogen/agents", 0755)
	if err != nil {
		return err
	}

	for _, page := range cfg.Pages {
		err := generatePageAutoTest(page)
		if err != nil {
			return err
		}
	}

	for _, service := range cfg.Services {
		err := generateServiceAutoTest(cfg.App.Name, service)
		if err != nil {
			return err
		}
	}

	for _, agent := range cfg.Agents {
		err := generateAgentAutoTest(agent)
		if err != nil {
			return err
		}
	}

	return nil
}

func generatePageAutoTest(page config.PageConfig) error {
	tmplData, err := templates.FS.ReadFile("tests/autogen/page.tsx.tmpl")
	if err != nil {
		return err
	}

	tmpl, err := template.New("page_test").Parse(string(tmplData))
	if err != nil {
		return err
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, page)
	if err != nil {
		return err
	}

	dest := fmt.Sprintf(".uca/tests/autogen/pages/autogen_%s.test.tsx", page.Name)
	return os.WriteFile(dest, buf.Bytes(), 0644)
}

func generateServiceAutoTest(appName string, service config.ServiceConfig) error {
	tmplData, err := templates.FS.ReadFile("tests/autogen/service.go.tmpl")
	if err != nil {
		return err
	}

	funcMap := template.FuncMap{
		"Lower": strings.ToLower,
	}

	tmpl, err := template.New("service_test").Funcs(funcMap).Parse(string(tmplData))
	if err != nil {
		return err
	}

	data := struct {
		AppName string
		config.ServiceConfig
	}{
		AppName:       appName,
		ServiceConfig: service,
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, data)
	if err != nil {
		return err
	}

	dest := fmt.Sprintf(".uca/tests/autogen/services/autogen_%s_test.go", service.Name)
	return os.WriteFile(dest, buf.Bytes(), 0644)
}

func generateAgentAutoTest(agent config.AgentConfig) error {
	tmplData, err := templates.FS.ReadFile("tests/autogen/agent.py.tmpl")
	if err != nil {
		return err
	}

	tmpl, err := template.New("agent_test").Parse(string(tmplData))
	if err != nil {
		return err
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, agent)
	if err != nil {
		return err
	}

	dest := fmt.Sprintf(".uca/tests/autogen/agents/autogen_%s_test.py", agent.Name)
	return os.WriteFile(dest, buf.Bytes(), 0644)
}
