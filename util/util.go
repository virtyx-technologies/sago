package util

import "os/exec"

func IsOnPath(name string) bool {
	var (
		command = "type"
		arg1 = "-p"
	)

	cmd := exec.Command(command, arg1, name)
	err := cmd.Run()
	return err == nil
}

