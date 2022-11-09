package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/vanilla-os/almost/core"
)

func overlayUsage(*cobra.Command) error {
	fmt.Print(`Description: 
		Overlay a directory to make it mutable and being able to edit its contents without modifying the originals
		
		Usage:
		overlay [options] [command] [directory]
		
		Options:
			--help/-h		show this message
			--verbose/-v		enable verbose output
			
		Commands:
			new [directory]			Overlay a directory
			commit [directory]		Commit the changes
			discard [directory]		Discard the changes
			list					List the active overlays
		
		Examples:
			# almost overlay new /etc/cute-path
			# almost overlay commit /etc/cute-path
			# almost overlay discard /etc/cute-path
			# almost overlay list`)
	return nil
}

func NewOverlayCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "overlay",
		Short: "Overlay a directory",
		RunE:  overlay,
	}
	cmd.SetUsageFunc(overlayUsage)
	cmd.Flags().BoolP("verbose", "v", false, "enable verbose output")
	return cmd
}

func overlay(cmd *cobra.Command, args []string) error {
	if !core.RootCheck(true) {
		return nil
	}

	if len(args) == 0 {
		return fmt.Errorf("missing command")
	}

	verbose, _ := cmd.Flags().GetBool("verbose")

	switch args[0] {
	case "new":
		if len(args) != 2 {return fmt.Errorf("missing command")}
		return core.OverlayAdd(args[1], false, verbose)
	case "commit":
		if len(args) != 2 {return fmt.Errorf("missing command")}
		return core.OverlayRemove(args[1], true, verbose)
	case "discard":
		if len(args) != 2 {return fmt.Errorf("missing command")}
		return core.OverlayRemove(args[1], false, verbose)
		
	case "list":
		return listOverlays()
	default:
		return fmt.Errorf("unknown command")
	}
}

func listOverlays() error {
	overlays := core.OverlayList()
	count := len(overlays)

	if count == 0 {
		fmt.Println("No overlays found")
		return nil
	}

	fmt.Printf("Found %d overlay(s):\n", count)

	for path, workdir := range overlays {
		fmt.Printf("%s -> %s\n", path, workdir)
	}
	return nil
}
