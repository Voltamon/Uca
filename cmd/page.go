package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/Voltamon/Uca/internal/manifest"
)

var pageCmd = &cobra.Command{
	Use:   "page",
	Short: "Manage pages",
}

var pageAddCmd = &cobra.Command{
	Use:   "add [name] [route]",
	Short: "Add a new page",
	Args:  cobra.RangeArgs(1, 2),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		route := "/" + strings.ToLower(name)
		if len(args) == 2 {
			route = args[1]
		}

		err := manifest.AddPage(name, route)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Printf("Added page %q at route %q\n", name, route)
		runTidy()
	},
}

var pageRemoveCmd = &cobra.Command{
	Use:   "remove [name]",
	Short: "Remove a page",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		err := manifest.RemovePage(args[0])
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Printf("Removed page %q\n", args[0])
	},
}

var pageListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all pages",
	Run: func(cmd *cobra.Command, args []string) {
		pages, err := manifest.ListPages()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		if len(pages) == 0 {
			fmt.Println("No pages defined")
			return
		}
		for _, p := range pages {
			fmt.Printf("  %s → %s\n", p.Name, p.Route)
		}
	},
}

func init() {
	pageCmd.AddCommand(pageAddCmd)
	pageCmd.AddCommand(pageRemoveCmd)
	pageCmd.AddCommand(pageListCmd)
	rootCmd.AddCommand(pageCmd)
}
