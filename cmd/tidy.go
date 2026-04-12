package cmd

import (
    "fmt"
    "os"

    "github.com/spf13/cobra"
    "github.com/Voltamon/Uca/internal/tidy"
)

var tidyCmd = &cobra.Command{
    Use:   "tidy",
    Short: "Reconcile project with uca.yaml",
    Run: func(cmd *cobra.Command, args []string) {
        _, err := tidy.Run()
        if err != nil {
            fmt.Println(err)
            os.Exit(1)
        }
    },
}

func init() {
    rootCmd.AddCommand(tidyCmd)
}
