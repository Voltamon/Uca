package manifest

import (
	"fmt"

	"github.com/Voltamon/Uca/internal/config"
)

func AddPage(name string, route string) error {
	cfg, err := readConfig()
	if err != nil {
		return err
	}

	for _, p := range cfg.Pages {
		if p.Name == name {
			return fmt.Errorf("page %q already exists", name)
		}
	}

	cfg.Pages = append(cfg.Pages, config.PageConfig{
		Name:  name,
		Route: route,
	})

	return writeConfig(cfg)
}

func RemovePage(name string) error {
	cfg, err := readConfig()
	if err != nil {
		return err
	}

	newPages := []config.PageConfig{}
	found := false
	for _, p := range cfg.Pages {
		if p.Name == name {
			found = true
			continue
		}
		newPages = append(newPages, p)
	}

	if !found {
		return fmt.Errorf("page %q not found", name)
	}

	cfg.Pages = newPages
	return writeConfig(cfg)
}

func ListPages() ([]config.PageConfig, error) {
	cfg, err := readConfig()
	if err != nil {
		return nil, err
	}
	return cfg.Pages, nil
}
