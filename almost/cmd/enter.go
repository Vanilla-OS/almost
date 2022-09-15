package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/vanilla-os/almost/core"
)

func enterUsage(*cobra.Command) error {
	fmt.Print(`Description: 
Set the filesystem as read-only or read-write until reboot.

Setting the filesystem as read-write mode may consist of a security risk, be
careful when using this command.

Usage:
enter [options] [command]

Options:
	--help/-h		show this message
	--verbose/-v		verbose output

Commands:
	ro			set the filesystem as read-only
	rw			set the filesystem as read-write
	default			set the filesystem as defined in the configuration file

Examples:
	almost enter ro
	almost enter rw
`)	
	return nil
}

func NewEnterCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "enter",
		Short: "Enter sets the filesystem as read-only or read-write until reboot",
		RunE: enter,
	}
	cmd.SetUsageFunc(enterUsage)
	cmd.Flags().BoolP("verbose", "v", false, "verbose output")
	cmd.Flags().BoolP("on-persistent", "p", false, "used by systemd to enter in default mode only if the persistent mode is enabled")
	return cmd
}

func enter(cmd *cobra.Command, args []string) error {
	if !core.RootCheck(true) {
		return nil
	}
	
	if len(args) == 0 {
		return fmt.Errorf("missing command")
	}

	verbose, _ := cmd.Flags().GetBool("verbose")
	on_persistent, _ := cmd.Flags().GetBool("on-persistent")

	switch args[0] {
	case "ro":
		return core.EnterRo(verbose)
	case "rw":
		return core.EnterRw(verbose)
	case "default":
		return core.EnterDefault(verbose, on_persistent)
	default:
		return fmt.Errorf("unknown command: %s", args[0])
	}
}