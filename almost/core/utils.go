package core

import (
	"fmt"
	"os"
)

var almostDir = "/etc/almost"

func init() {
	if !RootCheck(false) {
		return
	}
	if _, err := os.Stat(almostDir); os.IsNotExist(err) {
		if err := os.Mkdir(almostDir, 0755); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
}

func RootCheck(display bool) bool {
	if os.Geteuid() != 0 {
		if display {
			fmt.Println("You must be root to run this command")
		}
		return false
	}
	return true
}
