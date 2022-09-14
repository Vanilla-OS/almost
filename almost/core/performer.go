package core

import (
	"fmt"
	"os/exec"
)

var managed_paths = []string{
	"/bin",
	"/lib",
	"/lib64",
	"/sbin",
	"/usr",
}

func EnterRo(verbose bool) error {
	if !RootCheck(false) {
		return nil
	}

	fmt.Println("Locking system..")

	for _, path := range managed_paths {
		if verbose {
			fmt.Println("Processing: ", path)
		}
		SetImmutableFlag(path, verbose, 0, false)
	}

	fmt.Println("System is now locked.")
	Set("Almost::CurrentMode", "0")
	return nil
}

func EnterRw(verbose bool) error {
	if !RootCheck(false) {
		return nil
	}
	
	fmt.Println("Unlocking system..")

	for _, path := range managed_paths {
		if verbose {
			fmt.Println("Processing: ", path)
		}
		SetImmutableFlag(path, verbose, 1, false)
	}

	fmt.Println("System is now unlocked.")
	Set("Almost::CurrentMode", "1")
	return nil
}

func EnterDefault(verbose bool) error {
	if !RootCheck(false) {
		return nil
	}
	
	config, _ := Get("Almost::DefaultMode")
	if config == "0" {
		EnterRo(verbose)
	} else {
		EnterRw(verbose)
	}
	return nil
}

func SetImmutableFlag(path string, verbose bool, state int, ifDiff bool) error {	
	if verbose {
		fmt.Println("Processing: ", path)
	}

	if ifDiff {
		// Here we check if the file already respects the Almost::CurrentMode
		// if it does, we skip it. This is useful to predict and restore the
		// prior state when taking temporary ownership.
		current := GetImmutableFlag(path)
		currentConfig, _ := Get("Almost::CurrentMode")
		if fmt.Sprint(current) == currentConfig {
			return nil
		}
	}

	if state == 1 {
		err := exec.Command("chattr", "-R", "-i", "-f", path).Run()
		if err != nil && verbose {
			fmt.Println("Error while removing immutable flag: ", err)
		}
	} else {
		err := exec.Command("chattr", "-R", "+i", "-f", path).Run()
		if err != nil && verbose {
			fmt.Println("Error while setting immutable flag: ", err)
		}
	}
	return nil
}

func GetImmutableFlag(path string) (int) {
	out, err := exec.Command("lsattr", path).Output()
	if err != nil {
		return 0
	}
	if string(out[4]) == "i" {
		return 1
	}
	return 0
}
