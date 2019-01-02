// +build windows

package globals

import (
	"fmt"
	"os/user"
)

func getConfigFilesLocations() []string {
	usr, _ := user.Current()

	return []string{
		fmt.Sprintf("%s\\AppData\\Virtyx\\", usr.HomeDir),
		".",
	}
}
