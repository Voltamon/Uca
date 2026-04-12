package manifest

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
	"github.com/Voltamon/Uca/internal/config"
)

func readConfig() (*config.Config, error) {
	data, err := os.ReadFile("uca.yaml")
	if err != nil {
		return nil, fmt.Errorf("uca.yaml not found — are you inside a uca project?")
	}

	var cfg config.Config
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to parse uca.yaml: %w", err)
	}

	return &cfg, nil
}

func writeConfig(cfg *config.Config) error {
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("failed to marshal uca.yaml: %w", err)
	}

	return os.WriteFile("uca.yaml", data, 0644)
}
