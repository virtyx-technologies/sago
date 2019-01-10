package globals

import (
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"regexp"
)

//////////////////////////////////////////////////////////////////////////////////
//
// Most of the variables in this package should be removed during development.
// They mostly represent environment variables used by the original shell script.
//
///////////////////////////////////////////////////////////////////////////////////

// ########## Definition of error codes
//
const ERR_BASH = 255        // Bash version incorrect
const ERR_CMDLINE = 254     // Cmd line couldn't be parsed
const ERR_FCREATE = 253     // Output file couldn't be created
const ERR_FNAMEPARSE = 252  // Input file couldn't be parsed
const ERR_NOSUPPORT = 251   // Feature requested is not supported
const ERR_OSSLBIN = 250     // Problem with OpenSSL binary
const ERR_DNSBIN = 249      // Problem with DNS lookup binaries
const ERR_OTHERCLIENT = 248 // Other client problem
const ERR_DNSLOOKUP = 247   // Problem with resolving IP addresses or names
const ERR_CONNECT = 246     // Connectivity problem
const ERR_CLUELESS = 245    // Weird state, either though user globals or testssl.sh
const ERR_RESOURCE = 244    // Resources testssl.sh needs couldn't be read
const ERR_CHILD = 242       // Child received a signal from master
const ALLOK = 0             // All is fine


var (
	Options    *viper.Viper
	Flags      *pflag.FlagSet
	configFile string
	configDir  string

	OpenSSL string // The path to the openssl executable
	DataDir string
	Meta    *OpenSslMetadata
)



// ########## Traps! Make sure that temporary files are cleaned up after use in ANY case
// TODO
//trap "cleanup" QUIT EXIT
//trap "child_error" USR1

// ########## Internal definitions
//
const VERSION = "3.0rc3"
const SWCONTACT = "dirk aet testssl globals.Dot sh"

var SWURL string
// TODO egrep -q "dev|rc|beta" <<< "$VERSION" && \
//      SWURL="https://testssl.sh/dev/" ||
//      SWURL="https://testssl.sh/"
const CVS_REL = "$(tail -5 \"$0\" | awk '/dirkw Exp/ { print $4\" \"$5\" \"$6}')"
const CVS_REL_SHORT = "$(tail -5 \"$0\" | awk '/dirkw Exp/ { print $4 }')"

var REL_DATE, GIT_REL, GIT_REL_SHORT string

const PROG_NAME = "$(basename \"$0\")"
const RUN_DIR = "$(dirname \"$0\")"

const UA_STD = "TLS tester from $SWURL"
const UA_SNEAKY = "Mozilla/5.0 (X11; Linux x86_64; rv:52.0) Gecko/20100101 Firefox/52.0"

// ########## End of Initialization part

// Other global vars just being declared here
//
var PRINTF = "" // which external printf to use. Empty presets the internal one, see #1130
var IKNOW_FNAME = false
var FIRST_FINDING = true          // is this the first finding we are outputting to file?
var JSONHEADER = true             // include JSON headers and footers in HTML file, if one is being created
var CSVHEADER = true              // same for CSV
var HTMLHEADER = true             // same for HTML
var SECTION_FOOTER_NEEDED = false // kludge for tracking whether we need to close the JSON section object
var GIVE_HINTS = false            // give an additional info to findings
var SERVER_SIZE_LIMIT_BUG = false // Some servers have either a ClientHello total size limit or a 128 cipher limit (e.g. old ASAs)
var CHILD_MASS_TESTING bool
var HAD_SLEPT = 0
var NR_SOCKET_FAIL = 0  // Counter for socket failures
var NR_OSSL_FAIL = 0    // .. for OpenSSL connects
var NR_HEADER_FAIL = 0  // .. for HTTP_GET
var PROTOS_OFFERED = "" // This keeps which protocol is being offered. See has_server_protocol().
var DETECTED_TLS_VERSION = ""
var TLS_EXTENSIONS = ""

var Targets []string

const NPN_PROTOs = "spdy/4a2,spdy/3,spdy/3.1,spdy/2,spdy/1,http/1.1"

//  alpn_protos needs to be space-separated, not comma-seperated, including odd ones observed @ facebook and others, old ones like h2-17 omitted as they could not be found
const ALPN_PROTOs = "h2 spdy/3.1 http/1.1 h2-fb spdy/1 spdy/2 spdy/3 stun.turn stun.nat-discovery webrtc c-webrtc ftp"

var SESS_RESUMPTION []string
var TEMPDIR = ""
var TMPFILE = ""
var ERRFILE = ""
var CLIENT_AUTH = false
var BugsOpt  = ""
var NO_SSL_SESSIONID = false
var HOSTCERT = "" // File with host certificate, without intermediate certificate
var HEADERFILE = ""
var HEADERVALUE = ""
var HTTP_STATUS_CODE = ""
var DH_GROUP_OFFERED = ""
var DH_GROUP_LEN_P = 0
var KEY_SHARE_EXTN_NR = "33"  // The extension number for key_share was changed from 40 to 51 in TLSv1.3 draft 23.
										// In order to support draft 23 and later in addition to earlier drafts, need to
										// know which extension number to use. Note that it appears that a single
										// ClientHello cannot advertise both draft 23 and later and earlier drafts.
										// Preset may help to deal with STARTTLS + TLS 1.3 draft 23 and later but not earlier.
var BAD_SERVER_HELLO_CIPHER = false // reserved for cases where a ServerHello doesn't contain a cipher offered in the ClientHello
var GOST_STATUS_PROBLEM = false
var PATTERN2SHOW = ""
var SOCK_REPLY_FILE = ""
var NW_STR = ""
var LEN_STR = ""
var SNI = ""
var POODLE = ""    // keep vulnerability status for TLS_FALLBACK_SCSV
var OSSL_NAME = "" // openssl name, in case of LibreSSL it's LibreSSL
var OSSL_VER = ""  // openssl version, will be auto-determined
var OSSL_VER_MAJOR = 0
var OSSL_VER_MINOR = 0
var OSSL_VER_APPENDIX = "none"
var CLIENT_PROB_NO = 1
var HAS_DH_BITS bool // initialize openssl variables
var OSSL_SUPPORTED_CURVES = ""
var Port = 443 // unless otherwise auto-determined, see below
var NODE = ""
var NODEIP = ""
var rDNS = ""
var CORRECT_SPACES = "" // Used for IPv6 and proper output formatting
var IP46ADDRs = ""
var LOCAL_A = false    // Does the $NODEIP come from /etc/hosts?
var LOCAL_AAAA = false // Does the IPv6 IP come from /etc/hosts?
var XMPP_HOST = ""
var PROXYIP = ""   // $PROXYIP:$PROXPORT is your proxy if --proxy is defined ...
var PROXYPORT = "" // ... and openssl has proxy support
var PROXY = ""     // Once check_proxy() executed it contains $PROXYIP:$PROXPORT
var VULN_COUNT = 0
var SERVICE = "" // Is the server running an HTTP server, SMTP, POP or IMAP?
var URI = ""
var CERT_FINGERPRINT_SHA2 = ""
var RSA_CERT_FINGERPRINT_SHA2 = ""
var STARTTLS_PROTOCOL = ""
var OPTIMAL_PROTO = "" // Need this for IIS6 (sigh) + OpenSSL 1.0.2, otherwise some handshakes will fail see
// https://github.com/PeterMosmans/openssl/issues/19#issuecomment-100897892
var STARTTLS_OPTIMAL_PROTO = "" // Same for STARTTLS, see https://github.com/drwetter/testssl.sh/issues/188
var TLS_TIME = ""               // To keep the value of TLS server timestamp
var TLS_NOW = ""                // Similar
var TLS_DIFFTIME_SET = false    // Tells TLS functions to measure the TLS difftime or not
var NOW_TIME = ""
var HTTP_TIME = ""
var GET_REQ11 = ""
var START_TIME = 0     // time in epoch when the action started
var END_TIME = 0       // .. ended
var SCAN_TIME = 0      // diff of both: total scan time
var LAST_TIME = 0      // only used for performance measurements (MEASURE_TIME=true)
var SERVER_COUNTER = 0 // Counter for multiple servers

var TLS_LOW_BYTE = "" // For "secret" development stuff, see -q below
var HEX_CIPHER = ""   // "
var doDisplayOnly = false
var doMassTesting = false
var doMxAllIps = false

// ########## Global variables for parallel mass testing
//
const PARALLEL_SLEEP = 1 // Time to sleep after starting each test
var MAX_WAIT_TEST int    // Maximum time (in seconds) to wait for a test to complete
var MAX_PARALLEL int     // Maximum number of tests to run in parallel
// This value may be made larger on systems with faster processors
var PARALLEL_TESTING_PID []int        // process id for each child test (or 0 to indicate test has already completed)
var PARALLEL_TESTING_CMDLINE []string // command line for each child test
var NR_PARALLEL_TESTS int             // number of parallel tests run
var NEXT_PARALLEL_TEST_TO_FINISH = 0  // number of parallel tests that have completed and have been processed
var FIRST_JSON_OUTPUT = true          // true if no output has been added to $JSONFILE yet.

// ########## Severity functions and globals
//
var INFO = 0
var OK = 0
var LOW = 1
var MEDIUM = 2
var HIGH = 3
var CRITICAL = 4
var SEVERITY_LEVEL = 0

// Misc regular expressions
var RxCommas   = regexp.MustCompile(`\s*,\s*`) // Used to split string at comma with or without spaces
var RxFinalS   = regexp.MustCompile(`s$`)      // Used to remove trailing 's'
var RxFinalDot = regexp.MustCompile(`\.$`)      // Used to remove trailing '.'


var UrlPath string


