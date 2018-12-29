// +build !windows

package options

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetPaths(t *testing.T) {
	assert.Len(t, getConfigFilesLocations(), 3)
}
