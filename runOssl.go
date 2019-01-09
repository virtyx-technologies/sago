package main

import (
	. "github.com/virtyx-technologies/sago/globals"
	"io/ioutil"
	"os/exec"
)

func runOssl(command string, options []string, stdin string) (stdout, stderr string){

	var args []string
	args = append(args, command)
	args = append(args, options...)

	cmd := exec.Command(OpenSSL, args...)
	in, _ := cmd.StdinPipe()
	out, _ := cmd.StdoutPipe()
	err, _ := cmd.StderrPipe()
	in.Write([]byte(stdin))
	in.Close()
	outBytes, _ := ioutil.ReadAll(out)
	errBytes, _ := ioutil.ReadAll(err)
	cmd.Wait()

	return string(outBytes), string(errBytes)
}

func x () {
	runOssl("s_client", s_client_options( "$1 -quiet $BUGS -connect $NODEIP:$PORT $PROXY $SNI"), "$GET_REQ11")
}
// printf  | $OPENSSL   >$TMPFILE 2>$ERRFILE &
