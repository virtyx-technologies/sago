package util

import (
	"bufio"
	"os/exec"
	"strconv"
	"strings"
)


// TODO : need Windows implementation
func IsOnPath(name string) (bool, string) {
	path, err := exec.LookPath(name)
	return err == nil, string(path)
}

func Atoi(in string, defaultValue int) int {
	if out, err := strconv.Atoi(in); err != nil {
		return defaultValue
	} else {
		return out
	}
}

func Head(s string) string {
	scanner := bufio.NewScanner(strings.NewReader(s))
	i := 10
	var out string
	for scanner.Scan() {
		out += scanner.Text()
		out += "\n"
		i--
		if i == 0 {
			break
		}
	}

	return out
}
