package globals

import (
	"os"
	"path/filepath"
)

const (
	configDirectory      = ".sago"
	configDirPermissions = 0700
)

// Find a path to a previously created `.sago` directory. This takes the known
// places the directory may exist and checks for it's existence.
// If the directory is not found, attempt to create it.
func locateConfigDir() (string, error) {
	paths := getConfigFilesLocations()
	for _, path := range paths {
		fullPath := filepath.Join(path, configDirectory)
		if exists, err := pathExists(fullPath); exists {
			return fullPath, nil
		} else if err != nil {
			return "", err
		}
	}

	var err error
	for _, path := range paths {
		fullPath := filepath.Join(path, )
		err = os.MkdirAll(fullPath, configDirPermissions)
		if err == nil {
			return fullPath, nil
		}
	}
	return "", err
}

func pathExists(fullPath string) (bool, error) {
	_, err := os.Stat(fullPath)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
