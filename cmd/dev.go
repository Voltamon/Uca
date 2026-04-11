package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/Voltamon/Uca/internal/supervisor"
)

var devCmd = &cobra.Command{
	Use:   "dev",
	Short: "Start the development server",
	Run: func(cmd *cobra.Command, args []string) {
		err := supervisor.Start()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(devCmd)
}
