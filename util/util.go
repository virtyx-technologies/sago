package util

import "os/exec"


// TODO : Add .exe suffix on Windows
func IsOnPath(name string) (bool, string) {
	var (
		command = "type"
		arg1 = "-p"
	)

	cmd := exec.Command(command, arg1, name)
	path, err := cmd.Output()
	return err == nil, string(path)
}

