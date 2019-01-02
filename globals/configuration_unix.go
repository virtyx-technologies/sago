// +build !windows

package globals

import (
	"os"
	"os/user"
)

func getConfigFilesLocations() (dirs []string) {
	home := os.Getenv("HOME")
	usr, err := user.Current()

	if err == nil && usr != nil {
		dirs = append(dirs, usr.HomeDir)
	}

	if home != "" {
		dirs = append(dirs, home)
	}
	dirs = append(dirs, ".")
	return
}
