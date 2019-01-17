package main

import (
	. "github.com/virtyx-technologies/sago/globals"
	"github.com/virtyx-technologies/sago/stopwatch"
	log "github.com/sirupsen/logrus"
)

func main() {

	stopwatch.Click("start")

	// html_header() needs to be called early! Otherwise if html_out() is called before html_header() and the
	// command line contains --htmlfile <htmlfile> or --html, it'll make problems with html output, see #692.
	// json_header and csv_header could be called later but for context reasons we'll leave it here
	htmlHeader()
	jsonHeader()
	csvHeader()
	// see #705, we need to source TLS_DATA_FILE here instead of in get_install_dir(), see #705
	loadTlsVars() // See ./etc/tls_data.txt
	setColorFunctions()
	stopwatch.Click("parse")
	LoadCiphers()
	stopwatch.Click("LoadCiphers")
	mybanner()
	checkProxy()
	check4opensslOldfarts()
	checkBsdMount()

	if Options.GetBool(PrintCiphers) {
		prettyPrintLocal("$PATTERN2SHOW") // TODO - print local ciphers, see -V option
		return
	}
	fileoutBanner()

	// Mass testing means reading multiple command lines from the file specified by --file
	if Options.GetBool(DoMassTesting) {
		if Options.GetString("MASS-TESTING-MODE") == "parallel" {
			runMassTestingParallel()
		} else {
			runMassTesting()
		}
		return
	}
	htmlBanner()

	// TODO: there shouldn't be the need for a special case for --mx, only the ip addresses we would need upfront and the do-parser
	if Options.GetBool(DoMxAllIps) {
		if 1 == queryGlobals() { // if we have just 1x "do_*" --> we do a standard run -- otherwise just the one specified
			setScanningDefaults()
		}
		runMxAllIps(URI, DefaultPort) // we should reduce run_mx_all_ips to the stuff necessary as ~15 lines later we have similar code
		return
	}

	// Main loop
	for _, target := range Targets {
		letsRoll(target)
	}
	return
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
	if Meta.VerMajor < 1 {
		log.Fatal("Versions of openssl older than 1.0 are not supported")
	}
}

func checkProxy() {
	// TODO
}

func mybanner() {
	// TODO Banner for ToS, etc
}

func setColorFunctions() {
	// TODO
}

func loadTlsVars() {
	// TODO
}

