package core

import (
	"os/exec"
)

func OfflineUpdate() error {
	if err := OverlayAdd("/usr", true, true); err != nil {
		return err
	}

	EnterRw(true)

	// TODO: this should be done in a more elegant way, using a persistent
	// overlay and fstab entries

	cmd := exec.Command("/usr/libexec/pk-offline-update")
	if err := cmd.Run(); err != nil {
		if err := OverlayRemove("/usr", false, true); err != nil {
			return err
		}
		return err
	}

	if err := OverlayRemove("/usr", true, true); err != nil {
		return err
	}

	EnterRo(true)

	return nil
}
