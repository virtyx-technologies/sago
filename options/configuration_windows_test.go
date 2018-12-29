// +build windows

package options

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetPaths(t *testing.T) {
	assert.Len(t, getFilesLocations(), 2)
}

func TestGetDefault(t *testing.T) {
	assert.Contains(t, defaultFile(), "AppData")
}
