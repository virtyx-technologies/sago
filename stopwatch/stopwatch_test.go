package stopwatch

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestFormat(t *testing.T) {
	assert.NotNil(t, theStopWatch.start)
	now := time.Now()
	x := now.Add(2 * time.Minute).Add(3*time.Second).Add(5 * time.Millisecond)
	d := x.Sub(now)

	s := fmtDuration(d)
	assert.Equal(t, "00:02:03.005000000", s)
}
