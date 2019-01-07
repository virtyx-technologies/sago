package main

import (
	. "github.com/virtyx-technologies/sago/globals"
	"log"
	"net"
)

// arg1: ftp smtp, lmtp, pop3, imap, xmpp, telnet, ldap, postgres, mysql (maybe with trailing s)
func determine_service(service, node string, port int) {

	conn, err := net.Dial("tcp", net.JoinHostPort(node, string(port))
	conn.Close()
	if err != nil {
		// TODO handle error
	}

	if service == "" { // no STARTTLS.
		determine_optimal_proto(service)
		if SNEAKY {
			ua = UA_SNEAKY
		} else {
			ua = UA_STD
		}
		GET_REQ11 = "GET, URL_PATH HTTP/1.1\r\nHost:, NODE\r\nUser-Agent:, ua\r\nAccept-Encoding: identity\r\nAccept: text/*\r\nConnection: Close\r\n\r\n"
		// returns always 0:
		service_detection(OPTIMAL_PROTO)
	} else { // STARTTLS
		var protocol string
		if service == "postgres" {
			protocol = "postgres"
		} else {
			protocol = RxFinalS.ReplaceAllString(service, "") // strip trailing 's' in ftp(s), smtp(s), pop3(s), etc
		}
		var STARTTLS, SNI string
		switch protocol {
		case "ftp", "smtp", "lmtp", "pop3", "imap", "xmpp", "telnet", "ldap", "postgres", "mysql", "nntp":
			STARTTLS = "-starttls " + protocol
			SNI = ""
			switch protocol {
			case "xmpp":
				// for XMPP, openssl has a problem using -connect, NODEIP:$PORT. thus we use -connect, NODE:$PORT instead!
				NODEIP = "$NODE"
				if XmppHost != "" {
					if !HAS_XMPP {
						log.Fatal("Your, OPENSSL does not support the \"-xmpphost\" option", ERR_OSSLBIN)
					}
					STARTTLS = "$STARTTLS -xmpphost, XMPP_HOST" // small hack -- instead of changing calls all over the place
					// see http://xmpp.org/rfcs/rfc3920.html
				} else {
					if isIpv4Addr(node) {
						// XMPP needs a jabber domainname
						if [[ -n "$rDNS" ]] {
			prln_warning " IP address doesn't work for XMPP, trying PTR record, rDNS"
			// remove trailing .
			NODE =, {rDNS%%.}
			NODEIP =, {rDNS%%.}
			}else{
			log.Fatal( "No DNS supplied and no PTR record available which I can try for XMPP", ERR_DNSLOOKUP)
			}
			}
		}
case "postgres":
// Check if openssl version supports postgres.
if ! "$HAS_POSTGRES"; then
log.Fatal( "Your, OPENSSL does not support the \"-starttls postgres\" option", ERR_OSSLBIN)
}
case "mysql":
// Check if openssl version supports mysql.
if ! "$HAS_MYSQL"; then
log.Fatal( "Your, OPENSSL does not support the \"-starttls mysql\" option", ERR_OSSLBIN)
}
case "lmtp":
// Check if openssl version supports lmtp.
if ! "$HAS_LMTP"; then
log.Fatal( "Your, OPENSSL does not support the \"-starttls lmtp\" option", ERR_OSSLBIN)
}
case "nntp":
// Check if openssl version supports lmtp.
if ! "$HAS_NNTP"; then
log.Fatal( "Your, OPENSSL does not support the \"-starttls nntp\" option", ERR_OSSLBIN)
}
}
/*
$OPENSSL s_client, (s_client_options "-connect, NODEIP:$PORT, PROXY, BUGS, STARTTLS") 2>$ERRFILE >$TMPFILE </dev/null
if [[, ? -ne 0 ]]; then
debugme cat, TMPFILE | head -25
outln
log.Fatal( ", OPENSSL couldn't establish STARTTLS via, protocol to, NODEIP:$PORT", ERR_CONNECT)
}
grep -q '^Server Temp Key', TMPFILE && HAS_DH_BITS = true // FIX //190
out " Service set:$CORRECT_SPACES            STARTTLS via "
out "$(toupper "$protocol")"
[[ "$protocol" == mysql ]] && out " -- attention, this is experimental"
fileout "service" "INFO" "$protocol"
[[ -n "$XMPP_HOST" ]] && out " (XMPP domain=\'$XMPP_HOST\')"
outln
;;
*/
default:   outln
log.Fatal( "momentarily only ftp, smtp, lmtp, pop3, imap, xmpp, telnet, ldap, postgres, and mysql allowed", ERR_CMDLINE)

}
}
// tmpfile_handle, {FUNCNAME[0]}.txt
return 0 // OPTIMAL_PROTO, GET_REQ*/HEAD_REQ* is set now
}
