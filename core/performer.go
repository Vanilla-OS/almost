package core

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
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

func EnterDefault(verbose bool, on_persistent bool) error {
	if !RootCheck(false) {
		return nil
	}

	confDefault, _ := Get("Almost::DefaultMode")
	confPersist, _ := Get("Almost::PersistModeStatus")

	if on_persistent {
		// this is being called by the systemd unit on shutdown
		// here we check for offline updates, then set the rw mode
		// to allow PackageKit install them on next boot
		if PackageKitUpdatePrepared() || PackageKitUpgradePrepared() {
			fmt.Println("Offline updates found! Entering rw mode..")
			return EnterRw(verbose)
		}
		// with no updates found, we skip switching mode if the user
		// disabled the persistent mode
		if confPersist == "1" {
			fmt.Println("Persistent mode is disabled, nothing to do.")
			return nil
		}
	}

	if confDefault == "0" {
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

	files, _ := filepath.Glob(filepath.Join(path, "*"))

	for _, file := range files {
		fi, err := os.OpenFile(file, os.O_WRONLY, 0755)

		if err != nil {
			if verbose {
				fmt.Printf("Error opening %s: %v, try opening as a directory", file, err)
			}

			fi, err = os.Open(file)
			if err != nil {
				if verbose {
					fmt.Printf("Skipping %s: %s", file, err.Error())
					continue
				}
			}
		}

		if state == 0 {
			err := SetAttr(fi, FS_IMMUTABLE_FL)
			if err != nil && verbose {
				fmt.Println("Error while removing immutable flag: ", err)
			}
		} else {
			err := UnsetAttr(fi, FS_IMMUTABLE_FL)
			if err != nil && verbose {
				fmt.Println("Error while setting immutable flag: ", err)
			}
		}
	}

	return nil
}

func GetImmutableFlag(path string) int {
	out, err := exec.Command("lsattr", path).Output()
	if err != nil {
		return 0
	}
	if string(out[4]) == "i" {
		return 1
	}
	return 0
}
