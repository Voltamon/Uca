package manifest

import (
	"fmt"

	"github.com/Voltamon/Uca/internal/config"
)

func AddService(name string, methods []string) error {
	cfg, err := readConfig()
	if err != nil {
		return err
	}

	for _, s := range cfg.Services {
		if s.Name == name {
			return fmt.Errorf("service %q already exists", name)
		}
	}

	cfg.Services = append(cfg.Services, config.ServiceConfig{
		Name:    name,
		Methods: methods,
	})

	return writeConfig(cfg)
}

func RemoveService(name string) error {
	cfg, err := readConfig()
	if err != nil {
		return err
	}

	newServices := []config.ServiceConfig{}
	found := false
	for _, s := range cfg.Services {
		if s.Name == name {
			found = true
			continue
		}
		newServices = append(newServices, s)
	}

	if !found {
		return fmt.Errorf("service %q not found", name)
	}

	cfg.Services = newServices
	return writeConfig(cfg)
}

func ListServices() ([]config.ServiceConfig, error) {
	cfg, err := readConfig()
	if err != nil {
		return nil, err
	}
	return cfg.Services, nil
}
