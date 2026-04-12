package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/Voltamon/Uca/internal/manifest"
)

var serviceCmd = &cobra.Command{
	Use:   "service",
	Short: "Manage services",
}

var serviceAddCmd = &cobra.Command{
	Use:   "add [name] [methods...]",
	Short: "Add a new service",
	Args:  cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		methods := []string{}
		for _, m := range args[1:] {
			methods = append(methods, strings.ToUpper(m))
		}

		err := manifest.AddService(name, methods)
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
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		err := manifest.RemoveService(args[0])
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Printf("Removed service %q\n", args[0])
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
