package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/vanilla-os/almost/core"
)

func stateUsage(*cobra.Command) error {
	fmt.Print(`Description: 
	Manage persistent overlays.

Usage:
	state [options] [command]

Options:
	--help/-h		show this message
	--verbose/-v		enable verbose output
	
Commands:
	new			Create a new state
	rollback [id]		Rollback to a previous state

Examples:
	almost state new
	almost state rollback 1`)
	return nil
}

func NewStateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "state",
		Short: "Manage persistent overlays",
		RunE:  state,
	}
	cmd.SetUsageFunc(stateUsage)
	cmd.Flags().BoolP("verbose", "v", false, "enable verbose output")
	return cmd
}

func state(cmd *cobra.Command, args []string) error {
	if !core.RootCheck(true) {
		return nil
	}

	if len(args) == 0 {
		return fmt.Errorf("missing command")
	}

	// verbose, _ := cmd.Flags().GetBool("verbose")

	switch args[0] {
	case "new":
		return core.StateNew()
	case "rollback":
		if len(args) != 2 {
			return fmt.Errorf("missing command")
		}
		return core.StateRollback(args[1])
	default:
		return fmt.Errorf("unknown command")
	}
}
