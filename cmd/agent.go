package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/Voltamon/Uca/internal/manifest"
"github.com/Voltamon/Uca/internal/prompt"
)

var agentCmd = &cobra.Command{
	Use:   "agent",
	Short: "Manage agents",
}

var agentAddCmd = &cobra.Command{
	Use:   "add [name] [model]",
	Short: "Add a new agent",
	Args:  cobra.RangeArgs(0, 2),
	Run: func(cmd *cobra.Command, args []string) {
		var name, model string
		var err error

		if len(args) >= 1 {
			name = args[0]
		} else {
			name, err = prompt.AskRequired("Agent name")
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}

		if len(args) >= 2 {
			model = args[1]
		} else {
			model, err = prompt.AskDefault("Model", "gpt-4o-mini")
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}

		err = manifest.AddAgent(name, model)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Printf("Added agent %q with model %q\n", name, model)
		runTidy()
	},
}

var agentRemoveCmd = &cobra.Command{
	Use:   "remove [name]",
	Short: "Remove an agent",
	Args:  cobra.RangeArgs(0, 1),
	Run: func(cmd *cobra.Command, args []string) {
		var name string
		var err error

		if len(args) >= 1 {
			name = args[0]
		} else {
			agents, err := manifest.ListAgents()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			if len(agents) == 0 {
				fmt.Println("No agents defined")
				return
			}
			fmt.Println("Available agents:")
			for _, a := range agents {
				fmt.Printf("  %s (%s)\n", a.Name, a.Model)
			}
			name, err = prompt.AskRequired("Agent name to remove")
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}

		err = manifest.RemoveAgent(name)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Printf("Removed agent %q\n", name)
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
