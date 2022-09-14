package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/vanilla-os/almost/core"
)

func configUsage(*cobra.Command) error {
	fmt.Print(`Description: 
Manage and show the current configuration.

Usage:
almost config

Options:
	--help/-h		show this message
	
Commands:
	set [key] [value]	set a configuration value

Examples:
	almost config
	almost config set Almost::DefaultMode 1
`)
	return nil
}

func NewCmdConfig() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Show the current configuration",
		RunE:  config,
	}
	cmd.SetUsageFunc(configUsage)
	cmd.AddCommand(CmdConfigSet())
	return cmd
}

func config(cmd *cobra.Command, args []string) error {
	if !core.RootCheck(true) {
		return nil
	}
	return core.Show()
}

func CmdConfigSet() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set",
		Short: "Set a configuration value",
		RunE:  configSet,
	}
	return cmd
}

func configSet(cmd *cobra.Command, args []string) error {
	if !core.RootCheck(true) {
		return nil
	}

	if len(args) < 2 {
		return fmt.Errorf("missing key or value")
	}
	return core.Set(args[0], args[1])
}
