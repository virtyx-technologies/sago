package main

import (
	. "github.com/virtyx-technologies/sago/globals"
	"github.com/virtyx-technologies/sago/stopwatch"
	log "github.com/sirupsen/logrus"
	"regexp"
	"time"
)

var startTime, lastTime time.Time

var (
	rxNotColon = regexp.MustCompile(`[^:]`)
	rxColon    = regexp.MustCompile(`[:]`)
	rxValid    = regexp.MustCompile(`[a-zA-Z0-9:]`)
)

func init() {
	startTime = time.Now()
	lastTime = startTime
	// TODO	[[ -n "$MEASURE_TIME_FILE" ]] && >"$MEASURE_TIME_FILE"
}

func letsRoll(service, node string, port int) int {
	var ret int
	sectionNumber := 1

	if "" == node {
		log.Fatal("NODE doesn't resolve to an IP address", node, ERR_DNSLOOKUP)
	}

	node = nodeIpToProperIp6(node)
	resetHostDependedVars()
	stopwatch.Click("determineRdns")

	SERVER_COUNTER++
	startTlsService := determineService(service, node, port) // STARTTLS service? Other will be determined here too. Returns always 0 or has already exited if fatal error occurred

	// "secret" devel globals --devel:
	//	$doTlsSockets && [[ $TLS_LOW_BYTE -eq 22 ]] && { sslv2Sockets "" "true"; echo $? ; exit $ALLOK
	//$doTlsSockets && [[ $TLS_LOW_BYTE -ne 22 ]] && { tlsSockets "$TLS_LOW_BYTE" "$HEX_CIPHER" "all"; echo $? ; exit $ALLOK
	//$doCipherMatch && { fileoutSectionHeader $sectionNumber false; runCipherMatch ${singleCipher}

	// all top level functions  now following have the prefix "run_"
	fileoutSectionHeader(&sectionNumber, false)
	if Options.GetBool("DoProtocols") {
		ret += runProtocols()
		stopwatch.Click("runProtocols()")
		ret += runNpn()
		stopwatch.Click("runNpn")
		ret += runAlpn()
		stopwatch.Click("runAlpn")
	}

	fileoutSectionHeader(&sectionNumber, true)
	if Options.GetBool("DoGrease") {
		ret += runGrease()
		stopwatch.Click("runGrease")
	}

	fileoutSectionHeader(&sectionNumber, true)
	if Options.GetBool("DoCipherlists") {
		ret += runCipherlists()
		stopwatch.Click("runCipherlists")
	}

	fileoutSectionHeader(&sectionNumber, true)
	if Options.GetBool("DoPfs") {
		ret += runPfs()
		stopwatch.Click("runPfs")
	}

	fileoutSectionHeader(&sectionNumber, true)
	if Options.GetBool("DoServerPreference") {
		ret += runServerPreference()
		stopwatch.Click("runServerPreference")
	}

	fileoutSectionHeader(&sectionNumber, true)
	if Options.GetBool("DoServerDefaults") {
		ret += runServerDefaults()
		stopwatch.Click("runServerDefaults")
	}

	if Options.GetBool("DoHeader") {
		// TODO: refactor this into functions
		fileoutSectionHeader(&sectionNumber, true)
		if SERVICE == "HTTP" {
			ret += runHttpHeader(UrlPath)
			ret += runHttpDate(UrlPath)
			ret += runHsts(UrlPath)
			ret += runHpkp(UrlPath)
			ret += runServerBanner(UrlPath)
			ret += runApplBanner(UrlPath)
			ret += runCookieFlags(UrlPath)
			ret += runSecurityHeaders(UrlPath)
			ret += runRpBanner(UrlPath)
			stopwatch.Click("doHeader")
		} else {
			sectionNumber++
		}
	}

	// vulnerabilities
	if VULN_COUNT <= Options.GetInt("VULN_THRESHLD") && Options.GetBool("DoVulnerabilities") {
		outln()
		prHeadlineln(" Testing vulnerabilities ")
		outln()
	}

	fileoutSectionHeader(&sectionNumber, true)
	if Options.GetBool("DoHeartbleed") {
		ret += runHeartbleed()
		stopwatch.Click("runHeartbleed")
	}
	if Options.GetBool("DoCcsInjection") {
		ret += runCcsInjection()
		stopwatch.Click("runCcsInjection")
	}
	if Options.GetBool("DoTicketbleed") {
		ret += runTicketbleed()
		stopwatch.Click("runTicketbleed")
	}
	if Options.GetBool("DoRobot") {
		ret += runRobot()
		stopwatch.Click("runRobot")
	}
	if Options.GetBool("DoRenego") {
		ret += runRenego()
		stopwatch.Click("runRenego")
	}
	if Options.GetBool("DoCrime") {
		ret += runCrime()
		stopwatch.Click("runCrime")
	}
	if Options.GetBool("DoBreach") {
		runBreach(UrlPath)
	}
	if Options.GetBool("DoSslPoodle") {
		ret += runSslPoodle()
		stopwatch.Click("runSslPoodle")
	}
	if Options.GetBool("DoTlsFallbackScsv") {
		ret += runTlsFallbackScsv()
		stopwatch.Click("runTlsFallbackScsv")
	}
	if Options.GetBool("DoSweet32") {
		ret += runSweet32()
		stopwatch.Click("runSweet32")
	}
	if Options.GetBool("DoFreak") {
		ret += runFreak()
		stopwatch.Click("runFreak")
	}
	if Options.GetBool("DoDrown") {
		ret += runDrown()
		stopwatch.Click("runDrown")
	}
	if Options.GetBool("DoLogjam") {
		ret += runLogjam()
		stopwatch.Click("runLogjam")
	}
	if Options.GetBool("DoBeast") {
		ret += runBeast()
		stopwatch.Click("runBeast")
	}
	if Options.GetBool("DoLucky13") {
		ret += runLucky13()
		stopwatch.Click("runLucky13")
	}
	if Options.GetBool("DoRc4") {
		ret += runRc4()
		stopwatch.Click("runRc4")
	}

	fileoutSectionHeader(&sectionNumber, true)
	if Options.GetBool("DoAllciphers") {
		ret += runAllciphers()
		stopwatch.Click("runAllciphers")
	}
	if Options.GetBool("DoCipherPerProto") {
		ret += runCipherPerProto()
		stopwatch.Click("runCipherPerProto")
	}

	fileoutSectionHeader(&sectionNumber, true)
	if Options.GetBool("DoClientSimulation") {
		ret += runClientSimulation()
		stopwatch.Click("runClientSimulation")
	}

	fileoutSectionFooter(true)

	outln()
	calcScantime()
	datebanner(" Done")

	// reset the failed connect counter as we are finished
	NR_SOCKET_FAIL = 0
	NR_OSSL_FAIL = 0

	// TODO "$MEASURE_TIME" && printf "$1: %${COLUMNS}s\n" "$SCAN_TIME"
	//[[ -e "$MEASURE_TIME_FILE" ]] && echo "Total : $SCAN_TIME " >> "$MEASURE_TIME_FILE"

	return ret
}

func calcScantime() {

}

func resetHostDependedVars() {
	// TODO should probably be struct members
	TLS_EXTENSIONS = ""
	PROTOS_OFFERED = ""
	OPTIMAL_PROTO = ""
	SERVER_SIZE_LIMIT_BUG = false

}

func nodeIpToProperIp6(nodeIp string) string {
	if isIpv6Addr(nodeIp) {
		if !Options.GetBool("UNBRACKTD_IPV6") {
			nodeIp = "[" + nodeIp + "]"
		}
		// TODO don't think we need this
		// IPv6 addresses are longer, this variable takes care that "further IP" and "Service" is properly aligned
		// len_nodeip = ${//NODEIP}
		// CORRECT_SPACES = "$(printf -- " "'%.s' $(eval "echo {1.."$((len_nodeip - 17))"}"))"
	}
	return nodeIp
}

func isIpv6Addr(s string) bool {

	if len(s) == 0 {
		return false
	}

	// less than 2x ":"
	if len(rxNotColon.ReplaceAllString(s, "")) < 2 {
		return false
	}

	// check on chars allowed:
	if len(rxValid.ReplaceAllString(s, "")) > 0 {
		return false
	}

	return true
}

func fileoutSectionHeader(sectionNumber *int, flag bool) {
	*sectionNumber++
	// TODO
}

func runAlpn() int {
	return 0 // TODO
}

func runNpn() int {
	return 0 // TODO
}

func runProtocols() int {
	return 0 // TODO
}

func runAllciphers() int {
	return 0 // TODO
}
func runApplBanner(UrlPath string) int {
	return 0 // TODO
}
func runBeast() int {
	return 0 // TODO
}
func runBreach(UrlPath string) int {
	return 0 // TODO
}
func runCcsInjection() int {
	return 0 // TODO
}
func runCipherPerProto() int {
	return 0 // TODO
}
func runCipherlists() int {
	return 0 // TODO
}
func runClientSimulation() int {
	return 0 // TODO
}
func runCookieFlags(UrlPath string) int {
	return 0 // TODO
}
func runCrime() int {
	return 0 // TODO
}
func runDrown() int {
	return 0 // TODO
}
func runFreak() int {
	return 0 // TODO
}
func runGrease() int {
	return 0 // TODO
}
func runHeartbleed() int {
	return 0 // TODO
}
func runHpkp(UrlPath string) int {
	return 0 // TODO
}
func runHsts(UrlPath string) int {
	return 0 // TODO
}
func runHttpDate(UrlPath string) int {
	return 0 // TODO
}
func runHttpHeader(UrlPath string) int {
	return 0 // TODO
}
func runLogjam() int {
	return 0 // TODO
}
func runLucky13() int {
	return 0 // TODO
}
func runPfs() int {
	return 0 // TODO
}
func runRc4() int {
	return 0 // TODO
}
func runRenego() int {
	return 0 // TODO
}
func runRobot() int {
	return 0 // TODO
}
func runRpBanner(UrlPath string) int {
	return 0 // TODO
}
func runSecurityHeaders(UrlPath string) int {
	return 0 // TODO
}
func runServerBanner(UrlPath string) int {
	return 0 // TODO
}
func runServerDefaults() int {
	return 0 // TODO
}
func runServerPreference() int {
	return 0 // TODO
}
func runSslPoodle() int {
	return 0 // TODO
}
func runSweet32() int {
	return 0 // TODO
}
func runTicketbleed() int {
	return 0 // TODO
}
func runTlsFallbackScsv() int {
	return 0 // TODO
}
