package tidy

import (
	"fmt"

	"github.com/Voltamon/Uca/internal/config"
	"github.com/Voltamon/Uca/internal/scaffold"
)

func generateAgentServer(cfg *config.Config) error {
	if len(cfg.Agents) == 0 {
		return nil
	}

	model := cfg.Agents[0].Model

	err := scaffold.CopyTemplate("agents/server.py", ".uca/venv/server.py", scaffold.TemplateVars{
		Model: model,
	})
	if err != nil {
		return fmt.Errorf("failed to generate agent server: %w", err)
	}

	fmt.Println("Generated: .uca/venv/server.py")
	return nil
}
