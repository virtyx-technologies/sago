package main

import (
	. "github.com/virtyx-technologies/sago/globals"
	. "github.com/virtyx-technologies/sago/util"
	"regexp"
	"strings"
)

var (
	http = regexp.MustCompile(`^HTTP/`)
	smtp = regexp.MustCompile(`SMTP|ESMTP|Exim|IdeaSmtpServer|Kerio Connect|Postfix`)
	pop  = regexp.MustCompile(`POP|Gpop|MailEnable POP3 Server|OK Dovecot|Cyrus POP3`)
	imap = regexp.MustCompile(`IMAP|IMAP4|Cyrus IMAP4IMAP4rev1|IMAP4REV1|Gimap`)
	ftp  = regexp.MustCompile(`FTP`)
	xmpp = regexp.MustCompile(`jabber|xmpp`)
	nntp = regexp.MustCompile(`Jive News|InterNetNews|NNRP|INN|Kerio Connect|NNTP Service|Kerio MailServer|NNTP server`)
)

// determines whether the port has an HTTP service running or not (plain TLS, no STARTTLS)
// arg1 could be the protocol determined as "working". IIS6 needs that
func service_detection(protocol string, node string, port int) string {

	var service string

	nodePort := node + ":" + string(port)

	if ! CLIENT_AUTH {
		// SNI is not standardized for !HTTPS but fortunately for other protocols s_client doesn't seem to care
		stdOut, _ := runOssl("s_client", s_client_options(protocol, "-quiet", BugsOpt, "-connect", nodePort, PROXY, SNI), GET_REQ11)
		head := Head(stdOut)

		switch {
		case http.MatchString(head):
			service = "HTTP"
			if strings.Contains(head, "MongoDB") {
				// MongoDB port 27017 will respond to a GET request with a mocked HTTP response
				service = "MongoDB"
			}
		case smtp.MatchString(head):
			service = "SMTP"
		case pop.MatchString(head):
			service = "POP"
		case imap.MatchString(head):
			service = "IMAP"
		case ftp.MatchString(head):
			service = "FTP"
		case xmpp.MatchString(head):
			service = "XMPP"
		case nntp.MatchString(head):
			service = "NNTP"
		}

		switch service {
		case "HTTP":
			// out " $SERVICE"
			// fileout "${jsonID}" "INFO" "$SERVICE"

		case "IMAP", "POP", "SMTP", "NNTP", "MongoDB":
			// out " $SERVICE, thus skipping HTTP specific checks"
			// fileout "${jsonID}" "INFO" "$SERVICE, thus skipping HTTP specific checks"

		default:
			if CLIENT_AUTH {
				//out " certificate-based authentication => skipping all HTTP checks"
				//echo "certificate-based authentication => skipping all HTTP checks" >$TMPFILE
				//fileout "${jsonID}" "INFO" "certificate-based authentication => skipping all HTTP checks"
			} else {
				// out " Couldn't determine what's running on port $PORT"
				if Options.GetBool(AssumeHttp) {
					service = "HTTP"
					//out " -- ASSUME_HTTP set though"
					//fileout "${jsonID}" "DEBUG" "Couldn't determine service -- ASSUME_HTTP set"
				} else {
					//out ", assuming no HTTP service => skipping all HTTP checks"
					//fileout "${jsonID}" "DEBUG" "Couldn't determine service, skipping all HTTP checks"
				}
			}

		}

		//outln "\n"
		//tmpfile_handle ${FUNCNAME[0]}.txt
	}
	return service
}
