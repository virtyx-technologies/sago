package stopwatch

import (
	"fmt"
	. "time"
)

type aStopWatch struct {
	start Time
}


var theStopWatch = &aStopWatch{start: Now()}

func  Click(label string) {
	delta := Since(theStopWatch.start)
	theStopWatch.start = Now()
	fmt.Printf("%32s: %14s", label, fmtDuration(delta))
}

// TODO : Surely there must be an easier way than this?
func fmtDuration(d Duration) string {
	x := d // .Round(Nano)
	h := x / Hour
	x -= h * Hour
	m := x / Minute
	x -= m * Minute
	s := x / Second
	n := x - s * Second
	return fmt.Sprintf("%02d:%02d:%02d.%09d", h, m, s, n)
}




