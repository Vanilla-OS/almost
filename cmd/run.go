package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
	"github.com/vanilla-os/almost/core"
)

func runUsage(*cobra.Command) error {
	fmt.Print(`Description: 
Runs a command in read-write mode and returns to read-only mode after the command exits.

Usage:
run [command]

Options:
	--help/-h		show this message
	--verbose/-v		verbose output

Examples:
	almost run ls
`)
	return nil
}

func NewRunCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "run",
		Short: "Runs a command in read-write mode and returns to read-only mode after the command exits.",
		RunE:  run,
	}
	cmd.SetUsageFunc(runUsage)
	cmd.Flags().BoolP("verbose", "v", false, "verbose output")
	cmd.Flags().SetInterspersed(false)
	return cmd
}

func run(cmd *cobra.Command, args []string) error {
	if !core.RootCheck(true) {
		return nil
	}

	if len(args) == 0 {
		return fmt.Errorf("missing command")
	}

	verbose, _ := cmd.Flags().GetBool("verbose")

	fmt.Println("Running command in read-write mode...")
	core.EnterRw(verbose)

	c := exec.Command(args[0], args[1:]...)
	c.Env = os.Environ()
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	c.Stdin = os.Stdin
	err := c.Run()

	if err != nil {
		fmt.Println(err)
	}

	core.EnterRo(verbose)
	return nil
}
