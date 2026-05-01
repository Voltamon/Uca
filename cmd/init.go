package cmd

import (
    "fmt"
    "os"

    "github.com/spf13/cobra"
    "github.com/Voltamon/Uca/internal/scaffold"
)

var initCmd = &cobra.Command{
    Use:   "init [app-name]",
    Short: "Scaffold a new Uca project",
    Args:  cobra.ExactArgs(1),
    Run: func(cmd *cobra.Command, args []string) {
        appName := args[0]
        err := scaffold.InitProject(appName)
        if err != nil {
            fmt.Println(err)
            os.Exit(1)
        }
    },
}

func init() {
    rootCmd.AddCommand(initCmd)
}
