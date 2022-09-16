package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/vanilla-os/almost/cmd"
)

var (
	Version = "0.0.2"
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
	overlay			overlay a directory
`)
}

func newAlmostCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "almost",
		Short: "Almost provides a simple way to set the filesystem as read-only or read-write",
		Version: Version,
	}
}

func init() {
	if _, err := os.Stat("/etc/almost"); os.IsNotExist(err) {
		os.Mkdir("/etc/almost", 0755)
	}
}

func main() {
	rootCmd := newAlmostCommand()
	rootCmd.AddCommand(cmd.NewEnterCommand())
	rootCmd.AddCommand(cmd.NewConfigCommand())
	rootCmd.AddCommand(cmd.NewCheckCommand())
	rootCmd.AddCommand(cmd.NewOverlayCommand())
	rootCmd.SetHelpFunc(help)
	rootCmd.Execute()
}