package main

import (
	log "github.com/sirupsen/logrus"
	. "github.com/virtyx-technologies/sago/globals"
	"github.com/virtyx-technologies/sago/stopwatch"
	"github.com/virtyx-technologies/sago/util"
	"net"
	"regexp"
	"strings"
	"time"
)

var startTime, lastTime time.Time

func init() {
	startTime = time.Now()
	lastTime = startTime
	// TODO	[[ -n "$MEASURE_TIME_FILE" ]] && >"$MEASURE_TIME_FILE"
}

func letsRoll(target string) (ret int) {
	sectionNumber := 1

	service := Options.GetString(StartTls)
	node, err  := parseTarget(target)
	if err != nil {
		return
	}

	resetHostDependedVars()
	stopwatch.Click("determineRdns")
	SERVER_COUNTER++

	err = determineService(service, node)
	if err != nil {
		return
	}

	// "secret" devel globals --devel:
	//	$doTlsSockets && [[ $TLS_LOW_BYTE -eq 22 ]] && { sslv2Sockets "" "true"; echo $? ; exit $ALLOK
	//$doTlsSockets && [[ $TLS_LOW_BYTE -ne 22 ]] && { tlsSockets "$TLS_LOW_BYTE" "$HEX_CIPHER" "all"; echo $? ; exit $ALLOK
	//$doCipherMatch && { fileoutSectionHeader $sectionNumber false; runCipherMatch ${singleCipher}

	fileoutSectionHeader(&sectionNumber, false)
	if Options.GetBool("DoProtocols") {
		ret += runProtocols(node)
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



func runAlpn() int {
	return 0 // TODO
}

func runNpn() int {
	return 0 // TODO
}

func runProtocols(node *Node) int {
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

func parseTarget(target string) (node *Node, err error) {
	var (
		x, p, scheme, host string
		port int
		ip net.IP
	)

	rxScheme := regexp.MustCompile(`^(\w+)://`)
	if match := rxScheme.FindStringSubmatch(target); match != nil {
		scheme = match[1]
		strings.Replace(target, match[0], "", 1)
	}
	x, p, err = net.SplitHostPort(target)
	if err != nil {
		if strings.Contains(err.Error(), "missing port") {
			port = DefaultPort
			x = target
		} else {
			log.Error("Failed to parse target ", target)
			return
		}
	}

	port = util.Atoi(p, DefaultPort)

	ip = net.ParseIP(x)
	if ip == nil {
		host = x
		var addrs []net.IP
		addrs, err = net.LookupIP(host)
		if err != nil {
			return
		}
		ip = addrs[0]
	} else {
		var names []string
		names, err = net.LookupAddr(ip.String())
		if err != nil {
			return
		}
		host = names[0]
	}


	return &Node{Scheme: scheme, Host: host, Ip: ip, Port: port}, nil
}
