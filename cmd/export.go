package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/Voltamon/Uca/internal/export"
)

var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "Export the app as a deployable artifact",
	Run: func(cmd *cobra.Command, args []string) {
		includeEnv, _ := cmd.Flags().GetBool("include-env")
		err := export.Run(includeEnv)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	exportCmd.Flags().Bool("include-env", false, "Include .env file in artifact (contains secrets)")
	rootCmd.AddCommand(exportCmd)
}
