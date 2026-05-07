package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/Voltamon/Uca/internal/env"
	"github.com/Voltamon/Uca/internal/prompt"
)

var keysFile string

var keysCmd = &cobra.Command{
	Use:   "keys",
	Short: "Manage environment variable keys",
}

var keysAddCmd = &cobra.Command{
	Use:   "add [key] [value]",
	Short: "Add an environment variable (or sync from a template with -r)",
	Args:  cobra.MaximumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		if keysFile != "" {
			syncKeys(keysFile)
			return
		}

		if len(args) < 2 {
			fmt.Println("Usage: uca keys add [key] [value] OR uca keys add -r .env.uca")
			os.Exit(1)
		}

		key := args[0]
		val := args[1]

		err := env.EnsureKeyDeclared(key)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		err = env.Add(key, val)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Printf("Added %s to .env\n", key)
	},
}

func syncKeys(path string) {
	file, err := os.Open(path)
	if err != nil {
		fmt.Printf("Failed to open template file %s: %v\n", path, err)
		os.Exit(1)
	}
	defer file.Close()

	// Load existing .env keys
	existing := make(map[string]bool)
	envFile, err := os.Open(".env")
	if err == nil {
		scanner := bufio.NewScanner(envFile)
		for scanner.Scan() {
			line := scanner.Text()
			if strings.Contains(line, "=") {
				key := strings.Split(line, "=")[0]
				existing[strings.TrimSpace(key)] = true
			}
		}
		envFile.Close()
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.Split(line, "=")
		key := strings.TrimSpace(parts[0])

		if existing[key] || os.Getenv(key) != "" {
			continue
		}

		val, err := prompt.AskRequired(fmt.Sprintf("Enter value for %s", key))
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		if val != "" {
			err = env.Add(key, val)
			if err != nil {
				fmt.Printf("Failed to save %s: %v\n", key, err)
			} else {
				fmt.Printf("Saved %s to .env\n", key)
			}
		}
	}
}

var keysRemoveCmd = &cobra.Command{
	Use:   "remove [key]",
	Short: "Remove an environment variable",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		err := env.Remove(args[0])
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Printf("Removed %s from .env\n", args[0])
	},
}

var keysInfoCmd = &cobra.Command{
	Use:   "info",
	Short: "List all environment variable keys",
	Run: func(cmd *cobra.Command, args []string) {
		keys, err := env.Info()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		if len(keys) == 0 {
			fmt.Println("No environment variables set")
			return
		}
		for _, k := range keys {
			fmt.Println(" •", k)
		}
	},
}

func init() {
	keysAddCmd.Flags().StringVarP(&keysFile, "file", "r", "", "Read keys from template file")
	keysCmd.AddCommand(keysAddCmd)
	keysCmd.AddCommand(keysRemoveCmd)
	keysCmd.AddCommand(keysInfoCmd)
	rootCmd.AddCommand(keysCmd)
}
