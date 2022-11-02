package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/vanilla-os/almost/core"
)

func offlineUpdateUsage(*cobra.Command) error {
	fmt.Print(`Description: 
Performs an offline update of the system. This command is intended to be run
from the package-offline-update service, and should not be run manually.

Usage:
offline-update [command]

Options:
	--help/-h		show this message
`)
	return nil
}

func NewOfflineUpdateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "offline-update",
		Short: "Performs an offline update of the system.",
		RunE:  offlineUpdate,
	}
	cmd.SetUsageFunc(offlineUpdateUsage)
	return cmd
}

func offlineUpdate(cmd *cobra.Command, args []string) error {
	if !core.RootCheck(true) {
		return nil
	}

	err := core.OfflineUpdate()
	if err != nil {
		return err
	}

	return nil
}
