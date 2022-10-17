package main

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/vanilla-os/almost/cmd"
)

var (
	Version = "1.1.1"
)

func help(cmd *cobra.Command, args []string) {
	fmt.Print(`Usage: 
almost [options] [command]

Options:
	--help/-h		show this message
	--verbose/-v		show more verbosity
	--version/-V		show version

Commands:
	enter			set the filesystem as ro or rw until reboot
	config			show the current configuration
	check			check whether the filesystem is read-only or read-write
	run			runs a command in read-write mode and returns to read-only mode after the command exits
	shell			runs a shell in read-write mode and returns to read-only mode after the shell exits
	overlay			overlay a directory
`)
}

func newAlmostCommand() *cobra.Command {
	return &cobra.Command{
		Use:     "almost",
		Short:   "Almost provides a simple way to set the filesystem as read-only or read-write",
		Version: Version,
	}
}

func main() {
	rootCmd := newAlmostCommand()
	rootCmd.AddCommand(cmd.NewEnterCommand())
	rootCmd.AddCommand(cmd.NewConfigCommand())
	rootCmd.AddCommand(cmd.NewCheckCommand())
	rootCmd.AddCommand(cmd.NewOverlayCommand())
	rootCmd.AddCommand(cmd.NewRunCommand())
	rootCmd.AddCommand(cmd.NewShellCommand())
	rootCmd.SetHelpFunc(help)
	rootCmd.Execute()
}
