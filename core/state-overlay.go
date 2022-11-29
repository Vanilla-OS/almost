package core

import (
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"

	"golang.org/x/sys/unix"
)

var (
	stateSourcePath    = "/usr"
	statesPath         = "/var/almost/states"
	statesTrashPath    = "/var/almost/trash"
	stateMountUnitPath = "/etc/systemd/system/usr.mount"
	stateMountUnitName = "usr.mount"
)

func init() {
	if !RootCheck(false) {
		return
	}

	if err := os.MkdirAll(statesPath, 0755); err != nil {
		fmt.Println(err)
		return
	}

	if err := os.MkdirAll(statesTrashPath, 0755); err != nil {
		fmt.Println(err)
		return
	}
}

func StateNew() error {
	// preparing a new folder in /var/almost/states for the new state
	stateId := StateNextId()
	statePath := fmt.Sprintf("%s/%s", statesPath, stateId)
	if err := os.MkdirAll(statePath, 0755); err != nil {
		fmt.Println("Error creating new folder for state with Id:", stateId, err)
		return err
	}

	for _, dir := range []string{"data", "temp"} {
		if err := os.MkdirAll(fmt.Sprintf("%s/%s", statePath, dir), 0755); err != nil {
			fmt.Println("Error creating", dir, "directory:", err)
			return err
		}
	}

	// creating a new overlay in the just created folder
	if err := unix.Mount("overlay", stateSourcePath, "overlay", 0, fmt.Sprintf("lowerdir=%s,upperdir=%s/data,workdir=%s/temp", stateSourcePath, statePath, statePath)); err != nil {
		fmt.Println("Error creating new overlay for state with Id:", stateId, err)
		//fmt.Println("command was:", fmt.Sprintf("mount -t overlay overlay %s -o lowerdir=%s,upperdir=%s/data,workdir=%s/temp", stateSourcePath, stateSourcePath, statePath, statePath))
		return err
	}

	// to avoid creating new timelines during a travel to the future, we
	// need to empty the trash so that the new state lives in the same
	// timeline as the older ones
	StateEmptyTrash()

	// at this point the new state is mounted and ready to be used in the
	// transaction, now we are going to re-generate the fstab file to
	// include it, so that it will be mounted at boot
	StateMountUnitRegenerate()

	fmt.Println("New state created with Id:", stateId)

	return nil
}

func StateList() ([]string, []string, error) {
	states := []string{}
	trashedStates := []string{}

	// listing all states
	files, err := os.ReadDir(statesPath)
	if err != nil {
		return states, trashedStates, err
	}

	for _, file := range files {
		if file.IsDir() {
			if _, err := strconv.Atoi(file.Name()); err == nil {
				states = append(states, file.Name())
			}
		}
	}

	// listing all trashed states
	files, err = os.ReadDir(statesTrashPath)
	if err != nil {
		return states, trashedStates, err
	}

	for _, file := range files {
		if file.IsDir() {
			if _, err := strconv.Atoi(file.Name()); err == nil {
				trashedStates = append(trashedStates, file.Name())
			}
		}
	}

	return states, trashedStates, nil
}

func StateNextId() string {
	// predicting the next state id based on the highest id
	states, _, err := StateList()
	if err != nil || len(states) == 0 {
		return "0"
	}

	highestId := 0
	for _, state := range states {
		id, err := strconv.Atoi(state)
		if err != nil {
			continue
		}

		if id > highestId {
			highestId = id
		}
	}

	return fmt.Sprintf("%d", highestId+1)
}

func StateStatus() error {
	states, trashedStates, err := StateList()
	if err != nil {
		return err
	}

	fmt.Println("States")
	fmt.Println("----------------")
	for _, state := range states {
		fmt.Println("[" + state + "]")
	}

	fmt.Println("Trashed States")
	fmt.Println("----------------")
	for _, state := range trashedStates {
		fmt.Println("[" + state + "]")
	}

	return nil
}

func StateTrash(id string) error {
	// checking if the state exists
	statePath := fmt.Sprintf("%s/%s", statesPath, id)
	if _, err := os.Stat(statePath); os.IsNotExist(err) {
		return fmt.Errorf("state with Id %s does not exist", id)
	}

	// unmounting the state
	if err := unix.Unmount(statePath, 0); err != nil {
		return err
	}

	// moving the state to trash
	if err := os.Rename(statePath, fmt.Sprintf("%s/%s", statesTrashPath, id)); err != nil {
		return err
	}

	StateMountUnitRegenerate()

	fmt.Println("State", id, "trashed")

	return nil
}

func StateEmptyTrash() error {
	if err := os.RemoveAll(statesTrashPath); err != nil {
		return err
	}

	if err := os.MkdirAll(statesTrashPath, 0755); err != nil {
		return err
	}

	return nil
}

func StateMountUnitRegenerate() error {
	/*
		This function is responsible for generating the mount unit file for the
		new state three. Once the state is created, it is mounted in real time
		so that it can be used in the transaction. However, the state needs to
		be mounted at boot time as well, so that it can be used in the next
		session, and that is what this function does.

		Mount unit reference:
		[Unit]
		Description=Almost Overlay for /usr
		Documentation=https://documentation.vanillaos.org
		Before=systemd-remount-fs.service
		Wants=systemd-remount-fs.service

		[Mount]
		What=overlay
		Where=/usr
		Type=overlay
		Options=auto,lowerdir=/usr,upperdir=/var/almost/states/0/data,workdir=/var/almost/states/0/temp

		[Install]
		WantedBy=systemd-remount-fs.service
	*/

	// preparing the list of states, sorted by if from lowest to highest
	// so that the lower states are mounted first to keep the three
	// consistent
	states, _, err := StateList()
	if err != nil {
		return err
	}
	sort.Slice(states, func(i, j int) bool {
		return states[i] < states[j]
	})

	// preparing values for the Options field
	lowerdir := fmt.Sprintf("%s:", stateSourcePath)
	for _, state := range states[:len(states)-1] {
		lowerdir += fmt.Sprintf("%s/%s/data:", statesPath, state)
	}

	lowerdir = strings.TrimSuffix(lowerdir, ":")
	upperdir := fmt.Sprintf("%s/%s/data", statesPath, states[len(states)-1])
	workdir := fmt.Sprintf("%s/%s/temp", statesPath, states[len(states)-1])
	optionsField := fmt.Sprintf("auto,lowerdir=%s,upperdir=%s,workdir=%s", lowerdir, upperdir, workdir)

	// preparing and writing the new mount unit
	newSysUnit := `[Unit]
Description=Almost Overlay for ` + stateSourcePath + `
Documentation=https://documentation.vanillaos.org
After=systemd-remount-fs.service
Wants=systemd-remount-fs.service

[Mount]
What=overlay
Where=` + stateSourcePath + `
Type=overlay
Options=` + optionsField + `

[Install]
WantedBy=systemd-remount-fs.service`

	if err := os.WriteFile(stateMountUnitPath, []byte(newSysUnit), 0644); err != nil {
		return err
	}

	// reloading the systemd daemon to make systemd aware of the new unit
	if err := exec.Command("systemctl", "daemon-reload").Run(); err != nil {
		return err
	}

	// enabling the new unit so that it is mounted at boot time
	if err := exec.Command("systemctl", "enable", stateMountUnitName).Run(); err != nil {
		return err
	}

	fmt.Println("State mount unit regenerated")
	return nil
}

func StateRollback(id string) error {
	// unmount all states after the given id, plus the given id
	// move the states to trash
	// regenerate the mount unit
	// request a reboot
	return nil
}
