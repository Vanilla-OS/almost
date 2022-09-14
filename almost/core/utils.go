package core

import (
	"fmt"
	"os"
)

func RootCheck() bool {
	if os.Geteuid() != 0 {
		fmt.Println("You must be root to run this command")
		return false
	}
	return true
}
