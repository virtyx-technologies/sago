package main

import "os"

// ########## Definition of error codes
//
const ERR_BASH =255            // Bash version incorrect
const ERR_CMDLINE =254         // Cmd line couldn't be parsed
const ERR_FCREATE =253         // Output file couldn't be created
const ERR_FNAMEPARSE =252      // Input file couldn't be parsed
const ERR_NOSUPPORT =251       // Feature requested is not supported
const ERR_OSSLBIN =250         // Problem with OpenSSL binary
const ERR_DNSBIN =249          // Problem with DNS lookup binaries
const ERR_OTHERCLIENT =248     // Other client problem
const ERR_DNSLOOKUP =247       // Problem with resolving IP addresses or names
const ERR_CONNECT =246         // Connectivity problem
const ERR_CLUELESS =245        // Weird state, either though user options or testssl.sh
const ERR_RESOURCE =244        // Resources testssl.sh needs couldn't be read
const ERR_CHILD =242           // Child received a signal from master
const ALLOK =0                 // All is fine


// ########## Traps! Make sure that temporary files are cleaned up after use in ANY case
// TODO
//trap "cleanup" QUIT EXIT
//trap "child_error" USR1


// ########## Internal definitions
//
const VERSION ="3.0rc3"
const SWCONTACT ="dirk aet testssl dot sh"
var SWURL string
// egrep -q "dev|rc|beta" <<< "$VERSION" && \
//      SWURL="https://testssl.sh/dev/" ||
//      SWURL="https://testssl.sh/"
const CVS_REL ="$(tail -5 \"$0\" | awk '/dirkw Exp/ { print $4\" \"$5\" \"$6}')"
const CVS_REL_SHORT ="$(tail -5 \"$0\" | awk '/dirkw Exp/ { print $4 }')"
var REL_DATE, GIT_REL, GIT_REL_SHORT string
const PROG_NAME ="$(basename \"$0\")"
const RUN_DIR ="$(dirname \"$0\")"
var TESTSSL_INSTALL_DIR = "${TESTSSL_INSTALL_DIR:-\"\"}"  // If you run testssl.sh and it doesn't find it necessary file automagically set TESTSSL_INSTALL_DIR
var CA_BUNDLES_PATH = "${CA_BUNDLES_PATH:-\"\"}"          // You can have your stores some place else
var ADDITIONAL_CA_FILES = "${ADDITIONAL_CA_FILES:-\"\"}"  // single file with a CA in PEM format or comma separated lists of them
var CIPHERS_BY_STRENGTH_FILE = ""
var TLS_DATA_FILE = ""                                  // mandatory file for socket-based handshakes
var OPENSSL_LOCATION = ""
var HNAME = "$(hostname)"

var CMDLINE string
var CMDLINE_ARRAY = os.Args                // When performing mass testing, the child processes need to be sent the
var MASS_TESTING_CMDLINE []string                   // command line in the form of an array (see #702 and http://mywiki.wooledge.org/BashFAQ/050).


// ########## Defining (and presetting) variables which can be changed
//
//  Following variables make use of $ENV and can be used like "OPENSSL=<myprivate_path_to_openssl> ./testssl.sh <URI>"
var OPENSSL, OPENSSL_TIMEOUT string
var PHONE_OUT bool           // Whether testssl can retrieve CRLs and OCSP
var FAST_SOCKET bool       // EXPERIMENTAL feature to accelerate sockets -- DO NOT USE it for production
var COLOR = 2                       // 3: Extra color (ciphers, curves), 2: Full color, 1: B/W only 0: No ESC at all
var COLORBLIND bool         // if true, swap blue and green in the output
var SHOW_EACH_C bool       // where individual ciphers are tested show just the positively ones tested
var SHOW_SIGALGO bool     // "secret" switch whether testssl.sh shows the signature algorithm for -E / -e
var SNEAKY bool                 // is the referer and useragent we leave behind just usual?
var QUIET bool                   // don't output the banner. By doing this you acknowledge usage term appearing in the banner
var SSL_NATIVE bool         // we do per default bash sockets where possible "true": switch back to "openssl native"
var ASSUME_HTTP bool       // in seldom cases (WAF, old servers, grumpy SSL) service detection fails. "True" enforces HTTP checks
var BUGS = "TODO"                        // -bugs option from openssl, needed for some BIG IP F5
var WARNINGS = "off"                // can be either off or batch
var DEBUG = 0                       // 1: normal output the files in /tmp/ are kept for further debugging purposes
// 2: list more what's going on , also lists some errors of connections
// 3: slight hexdumps + other info,
// 4: display bytes sent via sockets
// 5: display bytes received via sockets
// 6: whole 9 yards
var FAST bool                     // preference: show only first cipher, run_allciphers with openssl instead of sockets
var WIDE bool                     // whether to display for some options just ciphers or a table w hexcode/KX,Enc,strength etc.
var MASS_TESTING_MODE = "serial"    // can be serial or parallel. Subject to change
var LOGFILE string                // logfile if used
var JSONFILE string              // jsonfile if used
var CSVFILE string                // csvfile if used
var HTMLFILE string              // HTML if used
var FNAME string                      // file name to read commands from
var FNAME_PREFIX string        // output filename prefix, see --outprefix
var APPEND bool                 // append to csv/json file instead of overwriting it
var NODNS  string    // If unset it does all DNS lookups per default. "min" only for hosts or "none" at all
var HAS_IPv6=false             // if you have OpenSSL with IPv6 support AND IPv6 networking set it to yes
var ALL_CLIENTS bool       // do you want to run all client simulation form all clients supplied by SSLlabs?
var OFFENSIVE =true            // do you want to include offensive vulnerability tests which may cause blocking by an IDS?

// ########## Tuning vars which cannot be set by a cmd line switch. Use instead e.g "HEADER_MAXSLEEP=10 ./testssl.sh <your_args_here>"
//
var EXPERIMENTAL bool     // a development hook which allows us to disable code
var PROXY_WAIT int            // waiting at max 20 seconds for socket reply through proxy
var DNS_VIA_PROXY = true    // do DNS lookups via proxy. --ip=proxy reverses this
var IGN_OCSP_PROXY bool // Also when --proxy is supplied it is ignored when testing for revocation via OCSP via --phone-out
var HEADER_MAXSLEEP int   // we wait this long before killing the process to retrieve a service banner / http header
var MAX_SOCKET_FAIL int   // If this many failures for TCP socket connects are reached we terminate
var MAX_OSSL_FAIL int       // If this many failures for s_client connects are reached we terminate
var MAX_HEADER_FAIL int   // If this many failures for HTTP GET are encountered we terminate
var MAX_WAITSOCK int        // waiting at max 10 seconds for socket reply. There shouldn't be any reason to change this.
var CCS_MAX_WAITSOCK int // for the two CCS payload (each). There shouldn't be any reason to change this.
var HEARTBLEED_MAX_WAITSOCK int      // for the heartbleed payload. There shouldn't be any reason to change this.
var STARTTLS_SLEEP int    // max time wait on a socket for STARTTLS. MySQL has a fixed value of 1 which can't be overwritten (#914)
var FAST_STARTTLS = true    // at the cost of reliability decrease the handshakes for STARTTLS
var USLEEP_SND = 0.1           // sleep time for general socket send
var USLEEP_REC = 0.2           // sleep time for general socket receive
var HSTS_MIN int               // >179 days is ok for HSTS
var HPKP_MIN int                // >=30 days should be ok for HPKP_MIN, practical hints?
var DAYS2WARN1 = 60            // days to warn before cert expires, threshold 1
var DAYS2WARN2 = 30            // days to warn before cert expires, threshold 2
var VULN_THRESHLD int       // if vulnerabilities to check >$VULN_THRESHLD we DON'T show a separate header line in the output each vuln. check
var UNBRACKTD_IPV6 = false // some versions of OpenSSL (like Gentoo) don't support [bracketed] IPv6 addresses
var NO_ENGINE bool           // if there are problems finding the (external) openssl engine set this to true
const CLIENT_MIN_PFS =5             // number of ciphers needed to run a test for PFS
var CAPATH = "${CAPATH:-/etc/ssl/certs/}"     // Does nothing yet (FC has only a CA bundle per default, ==> openssl version -d)
var GOOD_CA_BUNDLE = ""                       // A bundle of CA certificates that can be used to validate the server's certificate
var CERTIFICATE_LIST_ORDERING_PROBLEM = false // Set to true if server sends a certificate list that contains a certificate
// that does not certify the one immediately preceding it. (See RFC 8446, Section 4.4.2)
var STAPLED_OCSP_RESPONSE = ""
var HAS_DNS_SANS = false                      // Whether the certificate includes a subjectAltName extension with a DNS name or an application-specific identifier type.
var MEASURE_TIME bool
var MEASURE_TIME_FILE string
var DISPLAY_CIPHERNAMES = "openssl"           // display OpenSSL ciphername (but both OpenSSL and RFC ciphernames in wide mode)
const UA_STD ="TLS tester from $SWURL"
const UA_SNEAKY ="Mozilla/5.0 (X11; Linux x86_64; rv:52.0) Gecko/20100101 Firefox/52.0"

// ########## Initialization part, further global vars just being declared here
//
var PRINTF = ""                               // which external printf to use. Empty presets the internal one, see #1130
var IKNOW_FNAME = false
var FIRST_FINDING = true                      // is this the first finding we are outputting to file?
var JSONHEADER = true                         // include JSON headers and footers in HTML file, if one is being created
var CSVHEADER = true                          // same for CSV
var HTMLHEADER = true                         // same for HTML
var SECTION_FOOTER_NEEDED = false             // kludge for tracking whether we need to close the JSON section object
var GIVE_HINTS = false                        // give an additional info to findings
var SERVER_SIZE_LIMIT_BUG = false             // Some servers have either a ClientHello total size limit or a 128 cipher limit (e.g. old ASAs)
var CHILD_MASS_TESTING bool
var HAD_SLEPT = 0
var NR_SOCKET_FAIL = 0                        // Counter for socket failures
var NR_OSSL_FAIL = 0                          // .. for OpenSSL connects
var NR_HEADER_FAIL = 0                        // .. for HTTP_GET
var PROTOS_OFFERED = ""                       // This keeps which protocol is being offered. See has_server_protocol().
var DETECTED_TLS_VERSION = ""
var TLS_EXTENSIONS = ""
const NPN_PROTOs ="spdy/4a2,spdy/3,spdy/3.1,spdy/2,spdy/1,http/1.1"
//  alpn_protos needs to be space-separated, not comma-seperated, including odd ones observed @ facebook and others, old ones like h2-17 omitted as they could not be found
const ALPN_PROTOs ="h2 spdy/3.1 http/1.1 h2-fb spdy/1 spdy/2 spdy/3 stun.turn stun.nat-discovery webrtc c-webrtc ftp"
var SESS_RESUMPTION []string
var TEMPDIR = ""
var TMPFILE = ""
var ERRFILE = ""
var CLIENT_AUTH = false
var NO_SSL_SESSIONID = false
var HOSTCERT = ""                             // File with host certificate, without intermediate certificate
var HEADERFILE = ""
var HEADERVALUE = ""
var HTTP_STATUS_CODE = ""
var DH_GROUP_OFFERED = ""
var DH_GROUP_LEN_P = 0
var KEY_SHARE_EXTN_NR = "33"                  // The extension number for key_share was changed from 40 to 51 in TLSv1.3 draft 23.
// In order to support draft 23 and later in addition to earlier drafts, need to
// know which extension number to use. Note that it appears that a single
// ClientHello cannot advertise both draft 23 and later and earlier drafts.
// Preset may help to deal with STARTTLS + TLS 1.3 draft 23 and later but not earlier.
var BAD_SERVER_HELLO_CIPHER = false           // reserved for cases where a ServerHello doesn't contain a cipher offered in the ClientHello
var GOST_STATUS_PROBLEM = false
var PATTERN2SHOW = ""
var SOCK_REPLY_FILE = ""
var NW_STR = ""
var LEN_STR = ""
var SNI = ""
var POODLE = ""                               // keep vulnerability status for TLS_FALLBACK_SCSV
var OSSL_NAME = ""                            // openssl name, in case of LibreSSL it's LibreSSL
var OSSL_VER = ""                             // openssl version, will be auto-determined
var OSSL_VER_MAJOR = 0
var OSSL_VER_MINOR = 0
var OSSL_VER_APPENDIX = "none"
var CLIENT_PROB_NO = 1
var HAS_DH_BITS bool       // initialize openssl variables
var OSSL_SUPPORTED_CURVES = ""
var HAS_SSL2 = false
var HAS_SSL3 = false
var HAS_TLS13 = false
var HAS_PKUTIL = false
var HAS_PKEY = false
var HAS_NO_SSL2 = false
var HAS_NOSERVERNAME = false
var HAS_CIPHERSUITES = false
var HAS_COMP = false
var HAS_NO_COMP = false
var HAS_ALPN = false
var HAS_NPN = false
var HAS_FALLBACK_SCSV = false
var HAS_PROXY = false
var HAS_XMPP = false
var HAS_POSTGRES = false
var HAS_MYSQL = false
var HAS_LMTP = false
var HAS_NNTP = false
var HAS_IRC = false
var HAS_CHACHA20 = false
var HAS_AES128_GCM = false
var HAS_AES256_GCM = false
var PORT = 443                                // unless otherwise auto-determined, see below
var NODE = ""
var NODEIP = ""
var rDNS=""
var CORRECT_SPACES = ""                       // Used for IPv6 and proper output formatting
var IPADDRs []string
var IP46ADDRs=""
var LOCAL_A = false                           // Does the $NODEIP come from /etc/hosts?
var LOCAL_AAAA = false                        // Does the IPv6 IP come from /etc/hosts?
var XMPP_HOST = ""
var PROXYIP = ""                              // $PROXYIP:$PROXPORT is your proxy if --proxy is defined ...
var PROXYPORT = ""                            // ... and openssl has proxy support
var PROXY = ""                                // Once check_proxy() executed it contains $PROXYIP:$PROXPORT
var VULN_COUNT = 0
var SERVICE = ""                              // Is the server running an HTTP server, SMTP, POP or IMAP?
var URI = ""
var CERT_FINGERPRINT_SHA2 = ""
var RSA_CERT_FINGERPRINT_SHA2 = ""
var STARTTLS_PROTOCOL = ""
var OPTIMAL_PROTO = ""                        // Need this for IIS6 (sigh) + OpenSSL 1.0.2, otherwise some handshakes will fail see
// https://github.com/PeterMosmans/openssl/issues/19#issuecomment-100897892
var STARTTLS_OPTIMAL_PROTO = ""               // Same for STARTTLS, see https://github.com/drwetter/testssl.sh/issues/188
var TLS_TIME = ""                             // To keep the value of TLS server timestamp
var TLS_NOW = ""                              // Similar
var TLS_DIFFTIME_SET = false                  // Tells TLS functions to measure the TLS difftime or not
var NOW_TIME = ""
var HTTP_TIME = ""
var GET_REQ11 = ""
var START_TIME = 0                            // time in epoch when the action started
var END_TIME = 0                              // .. ended
var SCAN_TIME = 0                             // diff of both: total scan time
var LAST_TIME = 0                             // only used for performance measurements (MEASURE_TIME=true)
var SERVER_COUNTER = 0                        // Counter for multiple servers

var TLS_LOW_BYTE = ""                         // For "secret" development stuff, see -q below
var HEX_CIPHER = ""                           // "
var doDisplayOnly = false
var doMassTesting = false
var doMxAllIps = false


// ########## Global variables for parallel mass testing
//
const PARALLEL_SLEEP =1               // Time to sleep after starting each test
var MAX_WAIT_TEST int      // Maximum time (in seconds) to wait for a test to complete
var MAX_PARALLEL int          // Maximum number of tests to run in parallel
// This value may be made larger on systems with faster processors
var PARALLEL_TESTING_PID []int     // process id for each child test (or 0 to indicate test has already completed)
var PARALLEL_TESTING_CMDLINE []string    // command line for each child test
var NR_PARALLEL_TESTS int            // number of parallel tests run
var NEXT_PARALLEL_TEST_TO_FINISH=0 // number of parallel tests that have completed and have been processed
var FIRST_JSON_OUTPUT=true            // true if no output has been added to $JSONFILE yet.


// ########## Cipher suite information
//
var (
	TLS_NR_CIPHERS int
	TLS_CIPHER_HEXCODE string
	TLS_CIPHER_OSSL_NAME string
	TLS_CIPHER_RFC_NAME string
	TLS_CIPHER_SSLVERS string
	TLS_CIPHER_KX string
	TLS_CIPHER_AUTH string
	TLS_CIPHER_ENC string
	TLS_CIPHER_EXPORT string
	TLS_CIPHER_OSSL_SUPPORTED string
)
const TLS13_OSSL_CIPHERS="TLS_AES_128_GCM_SHA256:TLS_AES_256_GCM_SHA384:TLS_CHACHA20_POLY1305_SHA256:TLS_AES_128_CCM_SHA256:TLS_AES_128_CCM_8_SHA256"

// ########## Severity functions and globals
//
var INFO = 0
var OK = 0
var LOW = 1
var MEDIUM = 2
var HIGH = 3
var CRITICAL = 4
var SEVERITY_LEVEL = 0

func initGlobals() {
	if isDevBuild() {
	GIT_REL = "$(git log --format='%h %ci' -1 2>/dev/null | awk '{ print $1\" \"$2\" \"$3 }')"
	GIT_REL_SHORT = "$(git log --format='%h %ci' -1 2>/dev/null | awk '{ print $1 }')"
	REL_DATE = "$(git log --format='%h %ci' -1 2>/dev/null | awk '{ print $2 }')"
	} else {
	REL_DATE ="$(tail -5 \"$0\" | awk '/dirkw Exp/ { print $5 }')"
	}

	HSTS_MIN=HSTS_MIN * 86400     // correct to seconds
	HPKP_MIN=HPKP_MIN * 86400     // correct to seconds

	if MEASURE_TIME_FILE != "" {
		MEASURE_TIME=true
	}

}

// TODO
func isDevBuild() bool {
	return true
}
