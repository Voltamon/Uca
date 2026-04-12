package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/Voltamon/Uca/internal/manifest"
)

var agentCmd = &cobra.Command{
	Use:   "agent",
	Short: "Manage agents",
}

var agentAddCmd = &cobra.Command{
	Use:   "add [name] [model]",
	Short: "Add a new agent",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		err := manifest.AddAgent(args[0], args[1])
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Printf("Added agent %q with model %q\n", args[0], args[1])
		runTidy()
	},
}

var agentRemoveCmd = &cobra.Command{
	Use:   "remove [name]",
	Short: "Remove an agent",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		err := manifest.RemoveAgent(args[0])
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Printf("Removed agent %q\n", args[0])
	},
}

var agentListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all agents",
	Run: func(cmd *cobra.Command, args []string) {
		agents, err := manifest.ListAgents()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		if len(agents) == 0 {
			fmt.Println("No agents defined")
			return
		}
		for _, a := range agents {
			fmt.Printf("  %s (%s)\n", a.Name, a.Model)
		}
	},
}

func init() {
	agentCmd.AddCommand(agentAddCmd)
	agentCmd.AddCommand(agentRemoveCmd)
	agentCmd.AddCommand(agentListCmd)
	rootCmd.AddCommand(agentCmd)
}
