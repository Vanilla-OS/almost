package core

import (
	"fmt"
	"os"
)

var smartLockpath = "/etc/almost-smartlock"

func init() {
	if !RootCheck(false) {
		return
	}

	entryPoint, _ := Get("Almost::PkgManager::EntryPoint")
	newEntryPoint := entryPoint + "-almost"

	SetImmutableFlag(entryPoint, false, 1, false)

	if _, err := os.Stat(smartLockpath); os.IsNotExist(err) {
		f, err := os.Create(smartLockpath)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		f.WriteString(`#!/bin/sh
if [ "id -u" != 0 ]; then
	echo "You must be root to use the package manager."
	exit 1
fi
if sudo almost check | grep -q 'Mode: ro'; then
	echo "The system is locked, the package manager is disabled. Use apx instead or enter in rw mode."
else
` + newEntryPoint + ` $@
fi
`)
		f.Chmod(0755)
		f.Close()
	}

	if _, err := os.Stat(newEntryPoint); os.IsNotExist(err) {
		if _, err := os.Stat(entryPoint); !os.IsNotExist(err) {
			fmt.Println("Creating symlink...")
			err := os.Rename(entryPoint, newEntryPoint)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			os.Symlink(smartLockpath, entryPoint)
		}
	}

	SetImmutableFlag(entryPoint, false, 0, true)
}
