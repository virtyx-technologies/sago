package options

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/virtyx-technologies/agent-go/agent/logging"
)

// Removes all Config Files and ENV variables
func setup() {
	os.Clearenv()
	logging.HideLogs()
	for _, path := range getConfigFilesLocations() {
		folder := filepath.Join(path, virtyxAgentDirectory)
		os.RemoveAll(folder)
	}
}

func TestVirtyxDirectory(t *testing.T) {
	setup()
	path, err := locateConfigDir()
	assert.NoError(t, err)
	exists, err := pathExists(path)
	assert.True(t, exists)
}
