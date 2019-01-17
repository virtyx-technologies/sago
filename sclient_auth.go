package main

import (
	. "github.com/virtyx-technologies/sago/globals"
	"strings"
)
// <original-comment>
//    this is only being called from determine_optimal_proto in order to check whether we have a server
//    with client authentication, a server with no SSL session ID switched off
// </original-comment>
//
// Keeping this for side-effects: setting CLIENT_AUTH and NO_SSL_SESSIONID
//
func sclient_auth(node *Node) bool {

	stdOut, _, exitStatus :=runOssl("s_client",
		s_client_options( node.OptimalProtocol, BugsOpt,  "-connect", node.hostPort(), PROXY,  "-msg", "-starttls", "http"),
		"")
	if exitStatus == 0 {
		return true
	}

	// no client auth (CLIENT_AUTH=false is preset globally)
	if strings.Contains(stdOut, "Master-Key: ") { // connect succeeded
		if strings.HasPrefix(stdOut, "<<< .*CertificateRequest") { // CertificateRequest message in -msg
			node.features["CLIENT_AUTH"] = true
			return true
		}
		if !strings.Contains(stdOut, "Session-ID: ") { // probably no SSL session
			if 2 == strings.Count(stdOut, "CERTIFICATE") { // do another sanity check to be sure
				node.features["CLIENT_AUTH"] = false
				node.features["NO_SSL_SESSIONID"] = true // NO_SSL_SESSIONID is preset globally to false for all other cases
				return true
			}
		}
	}
	// what's left now is: master key empty, handshake returned not successful, session ID empty --> not successful
	return false
}
