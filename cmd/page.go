package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/Voltamon/Uca/internal/manifest"
"github.com/Voltamon/Uca/internal/prompt"
)

var pageCmd = &cobra.Command{
	Use:   "page",
	Short: "Manage pages",
}

var pageAddCmd = &cobra.Command{
	Use:   "add [name] [route]",
	Short: "Add a new page",
	Args:  cobra.RangeArgs(0, 2),
	Run: func(cmd *cobra.Command, args []string) {
		var name, route string
		var err error

		if len(args) >= 1 {
			name = args[0]
		} else {
			name, err = prompt.AskRequired("Page name")
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}

		if len(args) >= 2 {
			route = args[1]
		} else {
			route, err = prompt.AskDefault("Route", "/"+strings.ToLower(name))
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}

		err = manifest.AddPage(name, route)
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
	Args:  cobra.RangeArgs(0, 1),
	Run: func(cmd *cobra.Command, args []string) {
		var name string
		var err error

		if len(args) >= 1 {
			name = args[0]
		} else {
			pages, err := manifest.ListPages()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			if len(pages) == 0 {
				fmt.Println("No pages defined")
				return
			}
			fmt.Println("Available pages:")
			for _, p := range pages {
				fmt.Printf("  %s → %s\n", p.Name, p.Route)
			}
			name, err = prompt.AskRequired("Page name to remove")
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}

		err = manifest.RemovePage(name)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Printf("Removed page %q\n", name)
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
