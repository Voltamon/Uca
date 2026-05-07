package cmd

import (
	"fmt"
	"os"
	"github.com/spf13/cobra"
	"github.com/Voltamon/Uca/internal/prompt"
	"github.com/Voltamon/Uca/internal/deps"
)

var depsFile string

var depsCmd = &cobra.Command{
	Use:   "deps",
	Short: "Manage dependencies",
}

var depsAddCmd = &cobra.Command{
	Use:   "add [package]",
	Short: "Add a new dependency",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if depsFile != "" {
			err := deps.SyncFromFile(depsFile)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			return
		}

		var pkg string
		var err error
		if len(args) == 1 {
			pkg = args[0]
		} else {
			pkg, err = prompt.AskRequired("Package name")
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}

		choice, err := prompt.AskChoice("Where would you like to add this dependency?", []string{"Pages (npm)", "Services (go)", "Agents (pip)"})
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		d, _ := deps.Load()

		switch choice {
		case "Pages (npm)":
			if err := deps.AddPagesDep(pkg); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			d.Pages[pkg] = "latest"
		case "Services (go)":
			if err := deps.AddServicesDep(pkg); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			d.Services[pkg] = "latest"
		case "Agents (pip)":
			if err := deps.AddAgentsDep(pkg); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			d.Agents[pkg] = "latest"
		}

		deps.Save(d)
	},
}

func init() {
	depsAddCmd.Flags().StringVarP(&depsFile, "file", "r", "", "Read dependencies from file")
	depsCmd.AddCommand(depsAddCmd)
	rootCmd.AddCommand(depsCmd)
}
