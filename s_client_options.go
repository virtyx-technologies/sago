package main

import (
	. "github.com/virtyx-technologies/sago/globals"
	"regexp"
	"strings"
)

// Adjust options to 'openssl s_client' based on OpenSSL version and protocol version
func s_client_options(opts ...string) []string {
	var (
		ciphers        = "notpresent"
		tls13Ciphers   = "notpresent"
		rxCipher       = regexp.MustCompile(`\s-cipher\s+(.+)\s`)
		rxCipherSuites = regexp.MustCompile(`\s-ciphersuites\s+(.+)\s`)
		rxSsl2Ssl3     = regexp.MustCompile(`\s-ssl[23]\s`)
		rxTls1dot123   = regexp.MustCompile(`\s-tls1_[123]\s`)
	)

	options := strings.Join(opts, " ")

	// Extract the TLSv1.3 ciphers and the non-TLSv1.3 ciphers
	if matches := rxCipher.FindStringSubmatch(options); matches != nil {
		ciphers = "'" + matches[1] + "'"
	}

	if matches := rxCipherSuites.FindStringSubmatch(options); matches != nil {
		tls13Ciphers = "'" + matches[1] + "'"
		if tls13Ciphers == "ALL" {
			tls13Ciphers = TLS13_OSSL_CIPHERS
		}
	}

	// Don't include the -servername option for an SSLv2 or SSLv3 ClientHello.
	if SNI != "" && rxSsl2Ssl3.MatchString(options) {
		strings.Replace(options, SNI, "", -1)
	}

	// The server_name extension should not be included in the ClientHello unless
	// the -servername option is provided. However, OpenSSL 1.1.1 will include the
	// server_name extension unless the -noservername option is provided. So, if
	// the command line doesn't include -servername and the -noservername option is
	// supported, then add -noservername to the options.
	if Meta.HasNoServerName && !strings.Contains(options, " -servername ") {
		options += " -noservername"
	}

	// Newer versions of OpenSSL have dropped support for the -no_ssl2 option, so
	// remove any -no_ssl2 option if the option isn't supported. (Since versions of
	// OpenSSL that don't support -no_ssl2 also don't support SSLv2, the option
	// isn't needed for these versions of OpenSSL.)
	if !Meta.HasNoSsl2 {
		strings.Replace(options, "-no_ssl2", "", -1)
	}

	// At least one server will fail under some circumstances if compression methods are offered.
	// So, only offer compression methods if necessary for the test. In OpenSSL 1.1.0 and
	// 1.1.1 compression is only offered if the "-comp" option is provided.
	// OpenSSL 1.0.0, 1.0.1, and 1.0.2 offer compression unless the "-no_comp" option is provided.
	// OpenSSL 0.9.8 does not support either the "-comp" or the "-no_comp" option.
	if strings.Contains(options, " -comp ") {
		// Compression is needed for the test. So, remove "-comp" if it isn't supported, but
		// otherwise make no changes.
		if ! Meta.HasComp {
			strings.Replace(options, "-comp", "", -1)
		}
	} else {
		// Compression is not needed. So, specify "-no_comp" if that option is supported.
		if Meta.HasNoComp {
			options += " -no_comp"
		}
	}

	// If $OPENSSL is compiled with TLSv1.3 support and s_client is called without
	// specifying a protocol, but specifying a list of ciphers that doesn't include
	// any TLSv1.3 ciphers, then the command will always fail. So, if $OPENSSL supports
	// TLSv1.3 and a cipher list is provided, but no protocol is speci}ed, then add
	// -no_tls1_3 if no TLSv1.3 ciphers are provided.
	if Meta.HasTls13 && ciphers != "notpresent" &&
		(tls13Ciphers == "notpresent" || tls13Ciphers == "") &&
		!rxSsl2Ssl3.MatchString(options) &&
		!strings.Contains(options, "-tls1") &&
		!rxTls1dot123.MatchString(options) {
		options += " -no_tls1_3"
	}

	if ciphers != "notpresent" || tls13Ciphers != "notpresent" {
		if ! Meta.HasCipherSuites {
			if ciphers == "notpresent" {
				ciphers = ""
			}
			if tls13Ciphers == "notpresent" {
				tls13Ciphers = ""
			}
			if ciphers != "" && tls13Ciphers != "" {
				ciphers = ":" + ciphers
			}
			ciphers = "$tls13_ciphers$ciphers"
			options += " -cipher $ciphers"
		} else if ciphers != "notpresent" && ciphers != "" {
			options += " -cipher $ciphers"
		}
		if tls13Ciphers != "notpresent" && tls13Ciphers != "" {
			options += " -ciphersuites $tls13_ciphers"
		}
	}

	return strings.Fields(options)
}
