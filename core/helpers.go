package core

import (
	"os/exec"
)

func CurrentUser() string {
	cmd := exec.Command("logname")
	out, err := cmd.Output()
	if err != nil {
		return ""
	}
	user := string(out)
	return user[:len(user)-1]
}
