package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/vanilla-os/almost/core"
)

func statusUsage(*cobra.Command) error {
	fmt.Print(`Description: 
	Show information about the current state.

Usage:
	status

Options:
	--help/-h		show this message`)
	return nil
}

func NewStatusCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "status",
		Short: "Show information about the current state",
		RunE:  status,
	}
	cmd.SetUsageFunc(statusUsage)
	return cmd
}

func status(cmd *cobra.Command, args []string) error {
	if !core.RootCheck(true) {
		return nil
	}

	states, _, err := core.StateList()
	if err != nil {
		return err
	}

	if len(states) == 0 {
		fmt.Println("No states found.")
		return nil
	}

	fmt.Println("Current state:", states[len(states)-1])
	fmt.Println(len(states), "states found:")
	for _, state := range states {
		fmt.Println("-", state)
	}

	return nil
}
