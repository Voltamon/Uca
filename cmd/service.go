package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/Voltamon/Uca/internal/manifest"
	"github.com/Voltamon/Uca/internal/prompt"
	"github.com/Voltamon/Uca/internal/deps"
)

var serviceCmd = &cobra.Command{
	Use:   "service",
	Short: "Manage services",
}

var serviceAddCmd = &cobra.Command{
	Use:   "add [name]",
	Short: "Add a new service",
	Args:  cobra.RangeArgs(0, 1),
	Run: func(cmd *cobra.Command, args []string) {
		var name string
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

		err = manifest.AddService(name, []string{"GET", "POST", "PUT", "DELETE"})
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Printf("Added service %q\n", name)
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
				fmt.Printf("  %s\n", s.Name)
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

var serviceDepsCmd = &cobra.Command{
	Use:   "deps",
	Short: "Manage service dependencies",
}

var serviceDepsAddCmd = &cobra.Command{
	Use:   "add [package]",
	Short: "Add a new go module for services",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		pkg := args[0]
		if err := deps.AddServicesDep(pkg); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		d, _ := deps.Load()
		d.Services[pkg] = "latest"
		deps.Save(d)
	},
}

func init() {
	serviceDepsCmd.AddCommand(serviceDepsAddCmd)
	serviceCmd.AddCommand(serviceDepsCmd)
	serviceCmd.AddCommand(serviceAddCmd)
	serviceCmd.AddCommand(serviceRemoveCmd)
	serviceCmd.AddCommand(serviceListCmd)
	rootCmd.AddCommand(serviceCmd)
}
