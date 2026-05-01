package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/Voltamon/Uca/internal/logs"
)

var logsCmd = &cobra.Command{
	Use:   "logs",
	Short: "View development logs",
}

var logsViewCmd = &cobra.Command{
	Use:   "view",
	Short: "View recent logs",
	Run: func(cmd *cobra.Command, args []string) {
		source, _ := cmd.Flags().GetString("source")
		level, _ := cmd.Flags().GetString("level")
		lines, _ := cmd.Flags().GetInt("lines")

		err := logs.View(source, level, lines)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

var logsTailCmd = &cobra.Command{
	Use:   "tail",
	Short: "Tail live logs",
	Run: func(cmd *cobra.Command, args []string) {
		source, _ := cmd.Flags().GetString("source")
		level, _ := cmd.Flags().GetString("level")

		err := logs.Tail(source, level)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

var logsClearCmd = &cobra.Command{
	Use:   "clear",
	Short: "Clear all logs",
	Run: func(cmd *cobra.Command, args []string) {
		err := logs.Clear()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	logsViewCmd.Flags().String("source", "", "Filter by source (server, vite, agent)")
	logsViewCmd.Flags().String("level", "", "Filter by level (error, warn, info)")
	logsViewCmd.Flags().Int("lines", 100, "Number of lines to show")

	logsTailCmd.Flags().String("source", "", "Filter by source (server, vite, agent)")
	logsTailCmd.Flags().String("level", "", "Filter by level (error, warn, info)")

	logsCmd.AddCommand(logsViewCmd)
	logsCmd.AddCommand(logsTailCmd)
	logsCmd.AddCommand(logsClearCmd)
	rootCmd.AddCommand(logsCmd)
}
