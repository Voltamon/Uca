package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/Voltamon/Uca/internal/env"
)

var keysCmd = &cobra.Command{
	Use:   "keys",
	Short: "Manage environment variable keys",
}

var keysAddCmd = &cobra.Command{
	Use:   "add [key] [value]",
	Short: "Add an environment variable",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		err := env.Add(args[0], args[1])
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Printf("Added %s to .uca/.env\n", args[0])
	},
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
		fmt.Printf("Removed %s from .uca/.env\n", args[0])
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
	keysCmd.AddCommand(keysAddCmd)
	keysCmd.AddCommand(keysRemoveCmd)
	keysCmd.AddCommand(keysInfoCmd)
	rootCmd.AddCommand(keysCmd)
}
