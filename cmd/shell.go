package cmd

import (
	"fmt"
	"os/exec"

	"github.com/spf13/cobra"
	"github.com/vanilla-os/almost/core"
)

func shellUsage(*cobra.Command) error {
	fmt.Print(`Description: 
	Runs a shell in read-write mode and returns to read-only mode after the shell exits.

Usage:
	shell [command]

Options:
	--help/-h		show this message

Examples:
	almost shell
`)
	return nil
}

func NewShellCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "shell",
		Short: "Runs a shell in read-write mode and returns to read-only mode after the shell exits.",
		RunE:  shell,
	}
	cmd.SetUsageFunc(shellUsage)
	return cmd
}

func shell(cmd *cobra.Command, args []string) error {
	if !core.RootCheck(true) {
		return nil
	}

	core.EnterRw(false)
	fmt.Println("\033[33m⚠ WARNING: You are now in read-write mode.")
	fmt.Println("Any changes you make will be saved to the root filesystem and will persist after you exit.")
	fmt.Println("Use the `exit` command to return to read-only mode once you are done.\033[0m")

	c := exec.Command("su", core.CurrentUser())
	c.Stdout = cmd.OutOrStdout()
	c.Stderr = cmd.OutOrStderr()
	c.Stdin = cmd.InOrStdin()
	err := c.Run()

	if err != nil {
		fmt.Println(err)
	}

	core.EnterRo(false)
	fmt.Println("\033[32m✓ You are now in read-only mode.\033[0m")

	return nil
}
