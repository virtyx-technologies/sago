package main

import "testing"
import "github.com/stretchr/testify/assert"

func TestSimple(t *testing.T) {
	stdOut, stdErr, rc :=runOssl("verify", []string{"-help"}, "")
   assert.Empty(t, stdOut)
   assert.NotEmpty(t, stdErr) // Help is written to stdErr
	assert.Equal(t, 0, rc)
}
