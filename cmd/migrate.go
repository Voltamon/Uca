package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/Voltamon/Uca/internal/tidy"
)

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Generate configuration templates for the current project",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := tidy.Run()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		err = tidy.GenerateDepsJson()
		if err != nil {
			fmt.Printf("Failed to generate deps.json: %v\n", err)
		} else {
			fmt.Println("Generated: deps.json")
		}

		err = tidy.GenerateEnvUca(cfg)
		if err != nil {
			fmt.Printf("Failed to generate .env.uca: %v\n", err)
		} else {
			fmt.Println("Generated: .env.uca")
		}

		fmt.Println("\nMigration complete. You can now edit these files and run 'uca tidy' to sync.")
	},
}

func init() {
	rootCmd.AddCommand(migrateCmd)
}
