package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/Voltamon/Uca/internal/auth"
	"github.com/Voltamon/Uca/internal/prompt"
)

var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Manage roles and access control",
}

var authAddCmd = &cobra.Command{
	Use:   "add [role]",
	Short: "Add a new role",
	Args:  cobra.RangeArgs(0, 1),
	Run: func(cmd *cobra.Command, args []string) {
		var role string
		var err error

		if len(args) >= 1 {
			role = args[0]
		} else {
			role, err = prompt.AskRequired("Role name")
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}

		err = auth.AddRole(role)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		fmt.Printf("Added role %q\n", role)
		runTidy()
	},
}

var authRemoveCmd = &cobra.Command{
	Use:   "remove [role]",
	Short: "Remove a role",
	Args:  cobra.RangeArgs(0, 1),
	Run: func(cmd *cobra.Command, args []string) {
		var role string
		var err error

		if len(args) >= 1 {
			role = args[0]
		} else {
			roles, err := auth.ListRoles()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			if len(roles) == 0 {
				fmt.Println("No roles defined")
				return
			}
			fmt.Println("Available roles:")
			for _, r := range roles {
				if r == auth.DefaultRole {
					fmt.Printf("  %s (default, protected)\n", r)
				} else {
					fmt.Printf("  %s\n", r)
				}
			}
			role, err = prompt.AskRequired("Role name to remove")
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}

		err = auth.RemoveRole(role)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		fmt.Printf("Removed role %q\n", role)
		runTidy()
	},
}

var authListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all roles",
	Run: func(cmd *cobra.Command, args []string) {
		roles, err := auth.ListRoles()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		if len(roles) == 0 {
			fmt.Println("No roles defined")
			return
		}
		for _, r := range roles {
			if r == auth.DefaultRole {
				fmt.Printf("  %s (default)\n", r)
			} else {
				fmt.Printf("  %s\n", r)
			}
		}
	},
}

func init() {
	authCmd.AddCommand(authAddCmd)
	authCmd.AddCommand(authRemoveCmd)
	authCmd.AddCommand(authListCmd)
	rootCmd.AddCommand(authCmd)
}
