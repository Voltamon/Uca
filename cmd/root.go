package cmd

import (
    "fmt"
    "os"

    "github.com/spf13/cobra"
   "github.com/Voltamon/Uca/internal/tidy"
)

var rootCmd = &cobra.Command{
    Use:   "uca",
    Short: "Uca - A polyglot microframework",
    Long:  "Uca is a hermetic, polyglot microframework for building full-stack apps with AI agents.",
    CompletionOptions: cobra.CompletionOptions{
        DisableDefaultCmd: true,
    },
} 

func runTidy() {
	_, err := tidy.Run()
	if err != nil {
		fmt.Println("tidy failed:", err)
		os.Exit(1)
	}
}

func Execute() {
    if err := rootCmd.Execute(); err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
}
