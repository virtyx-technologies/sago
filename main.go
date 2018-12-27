package main

import ("github.com/virtyx-technologies/sago/stopwatch")

func main() {
	initGlobals()

	stopwatch.Click("start")

	checkBaseRequirements() // needs to come after $do_html is defined
	parseCmdLine("$@")
	// html_header() needs to be called early! Otherwise if html_out() is called before html_header() and the
	// command line contains --htmlfile <htmlfile> or --html, it'll make problems with html output, see #692.
	// json_header and csv_header could be called later but for context reasons we'll leave it here
	htmlHeader()
	jsonHeader()
	csvHeader()
	getInstallDir()
	// see #705, we need to source TLS_DATA_FILE here instead of in get_install_dir(), see #705
	loadTlsVars() // See ./etc/tls_data.txt
	setColorFunctions()
	maketempf()
	findOpensslBinary()
	choosePrintf()
	prepareDebug()
	stopwatch.Click("parse")
	prepareArrays()
	stopwatch.Click("prepare_arrays")
	mybanner()
	checkProxy()
	check4opensslOldfarts()
	checkBsdMount()

	if doDisplayOnly {
		prettyPrintLocal("$PATTERN2SHOW")
		return
	}
	fileoutBanner()

	if doMassTesting {
		prepareLogging()
		if MASS_TESTING_MODE == "parallel" {
			runMassTestingParallel()
		} else {
			runMassTesting()
		}
		return
	}
	htmlBanner()

	// TODO: there shouldn't be the need for a special case for --mx, only the ip addresses we would need upfront and the do-parser
	if doMxAllIps {
		if 1 == queryGlobals() { // if we have just 1x "do_*" --> we do a standard run -- otherwise just the one specified
			setScanningDefaults()
		}
		runMxAllIps(URI, PORT) // we should reduce run_mx_all_ips to the stuff necessary as ~15 lines later we have similar code
		return
	}

	if NODE == "" {
		parseHnPort(URI) // NODE, URL_PATH, PORT, IPADDRs and IP46ADDR is set now
	}
	prepareLogging()

	if ! determineIpAddresses() {
		fatal("No IP address could be determined", ERR_DNSLOOKUP)
	}
	if len(IPADDRs) > 1 { // we have more than one ipv4 address to check
		prBold("Testing all IPv4 addresses (port $PORT): ")
		outLine(IPADDRs)
		TERM_WIDTH := 100
		for _, ip := range IPADDRs {
			drawLine("-", (TERM_WIDTH * 2 / 3))
			outLine(nil)
			NODEIP = ip
			letsRoll("${STARTTLS_PROTOCOL}", ip)
			//  TODO : RET = $((RET + $?)) // RET value per IP address
		}
		drawLine("-", (TERM_WIDTH * 2 / 3))
		outLine()
		prBold("Done testing now all IP addresses (on port $PORT): ")
		outLine("$IPADDRs")
	} else { // Just 1x ip4v to check, applies also if CMDLINE_IP was supplied
		NODEIP = IPADDRs[0]
		letsRoll("${STARTTLS_PROTOCOL}", NODEIP)
		// RET=$?
	}
	return
}


func drawLine(s string, i int) {
	// TODO
}

func outLine(strings ...interface{}) {
	// TODO
}

func prBold(s string) {
	// TODO
}

func fatal(s string, i int) {
	// TODO
}

func determineIpAddresses() bool {
	// TODO
	return false
}

func parseHnPort(s string) {
	// TODO
}

func runMxAllIps(s string, i int) {
	// TODO
}

func setScanningDefaults() {
	// TODO
}

func queryGlobals() int {
	// TODO
	return 0
}

func htmlBanner() {
	// TODO
}

func runMassTesting() {
	// TODO
}

func runMassTestingParallel() {
	// TODO
}

func prepareLogging() {
	// TODO
}

func fileoutBanner() {
	// TODO
}

func prettyPrintLocal(s string) {
	// TODO
}

func checkBsdMount() {
	// TODO
}

func check4opensslOldfarts() {
	// TODO
}

func checkProxy() {
	// TODO
}

func mybanner() {
	// TODO
}

func prepareArrays() {
	// TODO
}

func prepareDebug() {
	// TODO
}

func choosePrintf() {
	// TODO
}

func findOpensslBinary() {
	// TODO
}

func maketempf() {
	// TODO
}

func setColorFunctions() {
	// TODO
}

func loadTlsVars() {
	// TODO
}

func getInstallDir() {
	// TODO
}

func csvHeader() {
	// TODO
}

func jsonHeader() {
	// TODO
}

func htmlHeader() {
	// TODO
}

func parseCmdLine(s string) {
	// TODO
}

func checkBaseRequirements() {
	// TODO
}


