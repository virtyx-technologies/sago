package main

import (
	"github.com/virtyx-technologies/sago/stopwatch"
	"log"
	"time"
)

var startTime, lastTime time.Time

func init() {
	startTime=time.Now()
	lastTime=startTime
// TODO	[[ -n "$MEASURE_TIME_FILE" ]] && >"$MEASURE_TIME_FILE"
}

func letsRoll(service, nodeIp string) int {
	var ret int
	sectionNumber := 1

stopwatch.Click("initialized")

	if "" == nodeIp {
		log.Fatal("NODE doesn't resolve to an IP address", NODEIP, ERR_DNSLOOKUP)
	}
	nodeIpToProperIp6()
	resetHostDependedVars()
	determineRdns()                // Returns always zero or has already exited if fatal error occurred
	stopwatch.Click("determineRdns")

	SERVER_COUNTER++
	determineService(service)       // STARTTLS service? Other will be determined here too. Returns always 0 or has already exited if fatal error occurred

	// "secret" devel options --devel:
//	$doTlsSockets && [[ $TLS_LOW_BYTE -eq 22 ]] && { sslv2Sockets "" "true"; echo $? ; exit $ALLOK
//$doTlsSockets && [[ $TLS_LOW_BYTE -ne 22 ]] && { tlsSockets "$TLS_LOW_BYTE" "$HEX_CIPHER" "all"; echo $? ; exit $ALLOK
//$doCipherMatch && { fileoutSectionHeader $sectionNumber false; runCipherMatch ${singleCipher}

// all top level functions  now following have the prefix "run_"
fileoutSectionHeader(&sectionNumber, false)
	if doProtocols {
ret += runProtocols()
stopwatch.Click("runProtocols()")
ret += runNpn()
stopwatch.Click("runNpn")
ret += runAlpn()
stopwatch.Click("runAlpn")
}

fileoutSectionHeader(&sectionNumber, true)
if doGrease {
	ret += runGrease()
stopwatch.Click("runGrease")
}

fileoutSectionHeader(&sectionNumber, true)
if doCipherlists { ret += runCipherlists()
stopwatch.Click("runCipherlists")
}

fileoutSectionHeader(&sectionNumber, true)
if doPfs { ret += runPfs()
stopwatch.Click("runPfs")
}

fileoutSectionHeader(&sectionNumber, true)
if doServerPreference { ret += runServerPreference()
stopwatch.Click("runServerPreference")
}

fileoutSectionHeader(&sectionNumber, true)
if doServerDefaults { ret += runServerDefaults()
stopwatch.Click("runServerDefaults")
}

if doHeader {
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
if VULN_COUNT <= VULN_THRESHLD  && doVulnerabilities {
outln()
prHeadlineln( " Testing vulnerabilities ")
outln()
}

fileoutSectionHeader(&sectionNumber, true)
if doHeartbleed { ret += runHeartbleed()
stopwatch.Click("runHeartbleed")
}
if doCcsInjection { ret += runCcsInjection()
stopwatch.Click("runCcsInjection")
}
if doTicketbleed { ret += runTicketbleed()
stopwatch.Click("runTicketbleed")
}
if doRobot { ret += runRobot()
stopwatch.Click("runRobot")
}
if doRenego { ret += runRenego()
stopwatch.Click("runRenego")
}
if doCrime { ret += runCrime()
stopwatch.Click("runCrime")
}
if doBreach { runBreach(UrlPath)
}
if doSslPoodle { ret += runSslPoodle()
stopwatch.Click("runSslPoodle")
}
if doTlsFallbackScsv { ret += runTlsFallbackScsv()
stopwatch.Click("runTlsFallbackScsv")
}
if doSweet32 { ret += runSweet32()
stopwatch.Click("runSweet32")
}
if doFreak { ret += runFreak()
stopwatch.Click("runFreak")
}
if doDrown { ret += runDrown()
stopwatch.Click("runDrown")
}
if doLogjam { ret += runLogjam()
stopwatch.Click("runLogjam")
}
if doBeast { ret += runBeast()
stopwatch.Click("runBeast")
}
if doLucky13 { ret += runLucky13()
stopwatch.Click("runLucky13")
}
if doRc4 { ret += runRc4()
stopwatch.Click("runRc4")
}

fileoutSectionHeader(&sectionNumber, true)
if doAllciphers { ret += runAllciphers()
stopwatch.Click("runAllciphers")
}
if doCipherPerProto { ret += runCipherPerProto()
stopwatch.Click("runCipherPerProto")
}

fileoutSectionHeader(&sectionNumber, true)
if doClientSimulation { ret += runClientSimulation()
stopwatch.Click("runClientSimulation")
}

fileoutSectionFooter(true)

outln()
calcScantime()
datebanner( " Done")

// reset the failed connect counter as we are finished
NR_SOCKET_FAIL=0
NR_OSSL_FAIL=0

// TODO "$MEASURE_TIME" && printf "$1: %${COLUMNS}s\n" "$SCAN_TIME"
//[[ -e "$MEASURE_TIME_FILE" ]] && echo "Total : $SCAN_TIME " >> "$MEASURE_TIME_FILE"

return ret
}

func calcScantime() {

}

func determineService(s string) {
	// TODO
}

func determineRdns() {
	// TODO
}

func resetHostDependedVars() {
	// TODO
}

func nodeIpToProperIp6() {
	// TODO
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

