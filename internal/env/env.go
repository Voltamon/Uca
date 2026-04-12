package env

import (
	"bufio"
	"fmt"
	"os"
	"strings"
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
	entries, err := readAll()
	if err != nil {
		return err
	}

	if _, exists := entries[key]; !exists {
		return fmt.Errorf("key %q not found", key)
	}

	delete(entries, key)
	return writeAll(entries)
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
