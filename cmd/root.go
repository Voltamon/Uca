package cmd

import (
    "fmt"
    "os"

    "github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
    Use:   "uca",
    Short: "Uca - A polyglot microframework",
    Long:  "Uca is a hermetic, polyglot microframework for building full-stack apps with AI agents.",
}

func Execute() {
    if err := rootCmd.Execute(); err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
}
