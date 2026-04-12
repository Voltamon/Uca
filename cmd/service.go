package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/Voltamon/Uca/internal/manifest"
	"github.com/Voltamon/Uca/internal/prompt"
)

var serviceCmd = &cobra.Command{
	Use:   "service",
	Short: "Manage services",
}

var serviceAddCmd = &cobra.Command{
	Use:   "add [name] [methods...]",
	Short: "Add a new service",
	Args:  cobra.ArbitraryArgs,
	Run: func(cmd *cobra.Command, args []string) {
		var name string
		var methods []string
		var err error

		if len(args) >= 1 {
			name = args[0]
		} else {
			name, err = prompt.AskRequired("Service name")
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}

		if len(args) >= 2 {
			for _, m := range args[1:] {
				methods = append(methods, strings.ToUpper(m))
			}
		} else {
			methodStr, err := prompt.AskDefault("Methods (comma separated)", "GET,POST")
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			for _, m := range strings.Split(methodStr, ",") {
				methods = append(methods, strings.ToUpper(strings.TrimSpace(m)))
			}
		}

		err = manifest.AddService(name, methods)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Printf("Added service %q with methods %v\n", name, methods)
		runTidy()
	},
}

var serviceRemoveCmd = &cobra.Command{
	Use:   "remove [name]",
	Short: "Remove a service",
	Args:  cobra.RangeArgs(0, 1),
	Run: func(cmd *cobra.Command, args []string) {
		var name string
		var err error

		if len(args) >= 1 {
			name = args[0]
		} else {
			services, err := manifest.ListServices()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			if len(services) == 0 {
				fmt.Println("No services defined")
				return
			}
			fmt.Println("Available services:")
			for _, s := range services {
				fmt.Printf("  %s %v\n", s.Name, s.Methods)
			}
			name, err = prompt.AskRequired("Service name to remove")
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}

		err = manifest.RemoveService(name)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Printf("Removed service %q\n", name)
	},
}
var serviceListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all services",
	Run: func(cmd *cobra.Command, args []string) {
		services, err := manifest.ListServices()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		if len(services) == 0 {
			fmt.Println("No services defined")
			return
		}
		for _, s := range services {
			fmt.Printf("  %s %v\n", s.Name, s.Methods)
		}
	},
}

func init() {
	serviceCmd.AddCommand(serviceAddCmd)
	serviceCmd.AddCommand(serviceRemoveCmd)
	serviceCmd.AddCommand(serviceListCmd)
	rootCmd.AddCommand(serviceCmd)
}
