package env

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
	"github.com/Voltamon/Uca/internal/config"
)

const envPath = ".uca/.env"

func Add(key string, value string) error {
	entries, err := readAll()
	if err != nil {
		return err
	}

	entries[key] = value
	return writeAll(entries)
}

func Remove(key string) error {
	err := removeKeyFromYaml(key)
	if err != nil {
		return err
	}

	entries, err := readAll()
	if err != nil {
		return err
	}

	if _, exists := entries[key]; !exists {
		return fmt.Errorf("key %q not found in .env", key)
	}

	delete(entries, key)
	return writeAll(entries)
}

func removeKeyFromYaml(key string) error {
	data, err := os.ReadFile("uca.yaml")
	if err != nil {
		return fmt.Errorf("uca.yaml not found — are you inside a uca project?")
	}

	var cfg config.Config
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return fmt.Errorf("failed to parse uca.yaml")
	}

	newKeys := []string{}
	for _, k := range cfg.App.Keys {
		if k != key {
			newKeys = append(newKeys, k)
		}
	}

	cfg.App.Keys = newKeys

	updated, err := yaml.Marshal(&cfg)
	if err != nil {
		return fmt.Errorf("failed to marshal uca.yaml")
	}

	return os.WriteFile("uca.yaml", updated, 0644)
}

func Info() ([]string, error) {
	entries, err := readAll()
	if err != nil {
		return nil, err
	}

	keys := make([]string, 0, len(entries))
	for k := range entries {
		keys = append(keys, k)
	}
	return keys, nil
}

func Load() (map[string]string, error) {
	return readAll()
}

func readAll() (map[string]string, error) {
	entries := make(map[string]string)

	f, err := os.Open(envPath)
	if err != nil {
		if os.IsNotExist(err) {
			return entries, nil
		}
		return nil, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			entries[parts[0]] = parts[1]
		}
	}

	return entries, scanner.Err()
}

func writeAll(entries map[string]string) error {
	err := os.MkdirAll(".uca", 0755)
	if err != nil {
		return err
	}

	f, err := os.Create(envPath)
	if err != nil {
		return err
	}
	defer f.Close()

	for k, v := range entries {
		_, err := fmt.Fprintf(f, "%s=%s\n", k, v)
		if err != nil {
			return err
		}
	}

	return nil
}

func EnsureKeyDeclared(key string) error {
	data, err := os.ReadFile("uca.yaml")
	if err != nil {
		return fmt.Errorf("uca.yaml not found — are you inside a uca project?")
	}

	var cfg config.Config
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return fmt.Errorf("failed to parse uca.yaml")
	}

	for _, k := range cfg.App.Keys {
		if k == key {
			return nil
		}
	}

	cfg.App.Keys = append(cfg.App.Keys, key)

	updated, err := yaml.Marshal(&cfg)
	if err != nil {
		return fmt.Errorf("failed to marshal uca.yaml")
	}

	err = os.WriteFile("uca.yaml", updated, 0644)
	if err != nil {
		return fmt.Errorf("failed to update uca.yaml")
	}

	fmt.Printf("Added %s to uca.yaml keys\n", key)
	return nil
}
