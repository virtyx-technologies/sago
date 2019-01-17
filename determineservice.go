package main

import (
	log "github.com/sirupsen/logrus"
	. "github.com/virtyx-technologies/sago/globals"
	"net"
	"strings"
)

// arg1: ftp smtp, lmtp, pop3, imap, xmpp, telnet, ldap, postgres, mysql (maybe with trailing s)
func determineService(startTLS string, node *Node) (err error){

	var service, optimalProtocol string

	conn, err := net.Dial("tcp", node.hostPort())
	conn.Close()
	if err != nil {
		log.Error("Failed to connect to ", node, "\n", err.Error())
		return
	}

	if startTLS == "" { // no STARTTLS.
		sclient_auth(node)
		// TODO: Defunct? optimalProtocol = determineOptimalProtocol("https") // TODO : Not sure is it's ok to hard-code this argument
		var ua string
		if Options.GetBool(Sneaky) {
			ua = UA_SNEAKY
		} else {
			ua = UA_STD
		}
		node.GetReq11 = `GET, `+node.Path+` HTTP/1.1
Host:, `+node.Host +`
User-Agent:, `+ua+`
Accept-Encoding: identity
Accept: text/*
Connection: Close

`
		service = serviceDetection(optimalProtocol, node)
	} else { // STARTTLS
		service = startTLS

		var protocol string
		if service == "postgres" {
			protocol = "postgres"
		} else {
			protocol = RxFinalS.ReplaceAllString(service, "") // strip trailing 's' in ftp(s), smtp(s), pop3(s), etc
		}
		switch protocol {
		case "ftp", "smtp", "lmtp", "pop3", "imap", "xmpp", "telnet", "ldap", "postgres", "mysql", "nntp":
			node.startTlsOpts = "-starttls " + protocol
			SNI = ""
			switch protocol {
			case "xmpp":
				// for XMPP, openssl has a problem using -connect, NODEIP:$PORT. thus we use -connect, NODE:$PORT instead!
				// NODEIP = "$NODE"  TODO: figure this out
				if Options.GetString(XmppHost) != "" {
					if !Meta.HasXmpp {
						log.Fatal("Your, OPENSSL does not support the \"-xmpphost\" option", ERR_OSSLBIN)
					}
					node.startTlsOpts = node.startTlsOpts + " -xmpphost, XMPP_HOST" // small hack -- instead of changing calls all over the place
					// see http://xmpp.org/rfcs/rfc3920.html
				} else {
					if node.IsIpv4Addr() {
						// XMPP needs a jabber domainname
						if rDNS := node.Host; rDNS != "" {
							log.Warn(" IP address doesn't work for XMPP, trying PTR record, rDNS")
						} else {
							log.Fatal("No DNS supplied and no PTR record available which I can try for XMPP", ERR_DNSLOOKUP)
						}
					}
				}
			case "postgres":
				// Check if openssl version supports postgres.
				if ! Meta.HasPostgres {
					log.Fatal("Your, OPENSSL does not support the \"-starttls postgres\" option", ERR_OSSLBIN)
				}
			case "mysql":
				// Check if openssl version supports mysql.
				if ! Meta.HasMysql {
					log.Fatal("Your, OPENSSL does not support the \"-starttls mysql\" option", ERR_OSSLBIN)
				}

			case "lmtp":
				// Check if openssl version supports lmtp.
				if ! Meta.HasLmtp {
					log.Fatal("Your, OPENSSL does not support the \"-starttls lmtp\" option", ERR_OSSLBIN)
				}
			case "nntp":
				// Check if openssl version supports lmtp.
				if ! Meta.HasNntp {
					log.Fatal("Your, OPENSSL does not support the \"-starttls nntp\" option", ERR_OSSLBIN)
				}
			}
			stdOut, stdErr, _ :=runOssl("s_client", s_client_options("-connect", "NODEIP:$PORT", PROXY, BugsOpt, node.startTlsOpts), "")
			if stdErr != "" {
				log.Fatal( ", OPENSSL couldn't establish STARTTLS via, protocol to, NODEIP:$PORT", ERR_CONNECT)
			}
			if strings.Contains(stdOut, "Server Temp Key") {
				HAS_DH_BITS = true // FIX //190
			}
			out(" Service set:$CORRECT_SPACES            STARTTLS via ")
			out( strings.ToUpper(protocol))
			if protocol == "mysql" {
				out(" -- attention, this is experimental")
			}
			fileout("service", "INFO", protocol)
			if Options.GetString(XmppHost) != ""{
				out( " (XMPP domain='$XMPP_HOST')")
			}
			outln()

		default:
			log.Fatal("currently only ftp, smtp, lmtp, pop3, imap, xmpp, telnet, ldap, postgres, and mysql allowed", ERR_CMDLINE)

		}
	}
	// tmpfile_handle, {FUNCNAME[0]}.txt
	return
}
