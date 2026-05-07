package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/Voltamon/Uca/internal/manifest"
	"github.com/Voltamon/Uca/internal/prompt"
	"github.com/Voltamon/Uca/internal/deps"
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

var pageDepsCmd = &cobra.Command{
	Use:   "deps",
	Short: "Manage page dependencies",
}

var pageDepsAddCmd = &cobra.Command{
	Use:   "add [package]",
	Short: "Add a new npm dependency for pages",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		pkg := args[0]
		if err := deps.AddPagesDep(pkg); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		d, _ := deps.Load()
		d.Pages[pkg] = "latest"
		deps.Save(d)
	},
}

func init() {
	pageDepsCmd.AddCommand(pageDepsAddCmd)
	pageCmd.AddCommand(pageDepsCmd)
	pageCmd.AddCommand(pageAddCmd)
	pageCmd.AddCommand(pageRemoveCmd)
	pageCmd.AddCommand(pageListCmd)
	rootCmd.AddCommand(pageCmd)
}
