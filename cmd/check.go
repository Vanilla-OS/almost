package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/vanilla-os/almost/core"
)

func checkUsage(*cobra.Command) error {
	fmt.Print(`Description: 
	Check whether the filesystem is read-only or read-write.

Usage:
	check [options] [command]

Options:
	--help/-h		show this message

Examples:
	almost check
`)
	return nil
}

func NewCheckCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "check",
		Short: "Check whether the filesystem is read-only or read-write",
		RunE:  check,
	}
	cmd.SetUsageFunc(checkUsage)
	return cmd
}

func check(cmd *cobra.Command, args []string) error {
	if !core.RootCheck(true) {
		return nil
	}

	mode, err := core.Get("Almost::CurrentMode")
	if err != nil {
		return err
	}
	if mode == "0" {
		fmt.Println("Mode: ro")
		fmt.Println("System is read-only")
	} else {
		fmt.Println("Mode: rw")
		fmt.Println("System is read-write")
	}
	return nil
}
