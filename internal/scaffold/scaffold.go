package scaffold

import (
    "fmt"
    "os"
)

func InitProject(appName string) error {
    dirs := []string{
        appName,
        appName + "/.uca",
        appName + "/agents",
        appName + "/assets",
        appName + "/modules",
        appName + "/pages",
        appName + "/services",
    }

    for _, dir := range dirs {
        err := os.MkdirAll(dir, 0755)
        if err != nil {
            return fmt.Errorf("failed to create directory %s: %w", dir, err)
        }
        fmt.Println("Created:", dir)
    }

    err := generateManifest(appName)
    if err != nil {
        return fmt.Errorf("failed to generate uca.yaml: %w", err)
    }

    fmt.Println("Generated: uca.yaml")

    err = GenerateFiles(appName, "github/gpt-4o", "8090", "user")
    if err != nil {
        return fmt.Errorf("failed to generate files: %w", err)
    }

    fmt.Println("Generated: starter files")
    return nil
}
