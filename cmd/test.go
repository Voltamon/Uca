package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/Voltamon/Uca/internal/tests"
)

var testCmd = &cobra.Command{
	Use:   "test",
	Short: "Run all tests",
	Run: func(cmd *cobra.Command, args []string) {
		goOnly, _ := cmd.Flags().GetBool("go")
		tsOnly, _ := cmd.Flags().GetBool("ts")
		pyOnly, _ := cmd.Flags().GetBool("py")

		err := tests.Run(goOnly, tsOnly, pyOnly)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

var testAddCmd = &cobra.Command{
	Use:   "add [name] [type]",
	Short: "Generate a test stub (type: go, ts, py)",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		err := tests.AddTestStub(args[0], args[1])
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	testCmd.Flags().Bool("go", false, "Run only Go tests")
	testCmd.Flags().Bool("ts", false, "Run only TypeScript tests")
	testCmd.Flags().Bool("py", false, "Run only Python tests")
	testCmd.AddCommand(testAddCmd)
	rootCmd.AddCommand(testCmd)
}
