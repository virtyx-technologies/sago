package main

import (
	. "github.com/virtyx-technologies/sago/globals"
	"io/ioutil"
	"log"
	"os/exec"
	"syscall"
)

func runOssl(command string, options []string, stdin string) (stdout, stderr string, exitStatus int){

	var args []string
	args = append(args, command)
	args = append(args, options...)

	cmd := exec.Command(OpenSSL, args...)
	inPipe, _  := cmd.StdinPipe()
	outPipe, _ := cmd.StdoutPipe()
	errPipe, _ := cmd.StderrPipe()
	inPipe.Write([]byte(stdin))
	inPipe.Close()

	if err := cmd.Start(); err != nil {
		log.Fatal("Failed to run command ", cmd, err)
	}
	outBytes, _ := ioutil.ReadAll(outPipe)
	errBytes, _ := ioutil.ReadAll(errPipe)

	err := cmd.Wait()
	if err != nil {
		// Did the command fail because of an unsuccessful exit code
		if exitError, ok := err.(*exec.ExitError); ok {
			exitStatus = exitError.Sys().(syscall.WaitStatus).ExitStatus()
		} else {
			log.Fatal("Command failed ", cmd, err)
		}
	}
	return string(outBytes), string(errBytes), exitStatus
}

