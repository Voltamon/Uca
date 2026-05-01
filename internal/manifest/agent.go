package manifest

import (
	"fmt"

	"github.com/Voltamon/Uca/internal/config"
)

func AddAgent(name string, model string) error {
	cfg, err := readConfig()
	if err != nil {
		return err
	}

	for _, a := range cfg.Agents {
		if a.Name == name {
			return fmt.Errorf("agent %q already exists", name)
		}
	}

	cfg.Agents = append(cfg.Agents, config.AgentConfig{
		Name:    name,
		Model:   model,
		Timeout: 30,
	})

	return writeConfig(cfg)
}

func RemoveAgent(name string) error {
	cfg, err := readConfig()
	if err != nil {
		return err
	}

	newAgents := []config.AgentConfig{}
	found := false
	for _, a := range cfg.Agents {
		if a.Name == name {
			found = true
			continue
		}
		newAgents = append(newAgents, a)
	}

	if !found {
		return fmt.Errorf("agent %q not found", name)
	}

	cfg.Agents = newAgents
	return writeConfig(cfg)
}

func ListAgents() ([]config.AgentConfig, error) {
	cfg, err := readConfig()
	if err != nil {
		return nil, err
	}
	return cfg.Agents, nil
}
