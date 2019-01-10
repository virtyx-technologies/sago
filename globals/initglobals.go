package globals

import (
	"github.com/virtyx-technologies/sago/util"
	log "github.com/sirupsen/logrus"
	"os"
	"path/filepath"
)

// Initialise global variables
func initGlobals() {
	OpenSSL = findOpensslBinary()
	Meta = NewOpenSslMetadata()
	DataDir = findDataDir()
	Targets = getTargets()

	getReleaseInfo()

	if Options.GetBool(Bugs) {
		BugsOpt  = "-bugs"
	}
}

func getTargets() []string {
	s := Options.GetString(Target)
	if s == "" {
		log.Fatal("No target hosts specified")
	}
	return RxCommas.Split(s, -1)
}

func findOpensslBinary() string {
	const openssl = "openssl"
	if found, path := util.IsOnPath(openssl); found {
		Options.SetDefault(OpenSslFile, path)
	}
	var path string
	if path = Options.GetString(OpenSslFile); path == "" {
		log.Fatal("Cannot locate openssl")
	}
	return path
}

func findDataDir() string {
	installDir := filepath.Dir(os.Args[0])
	Options.SetDefault(InstallDir, installDir)
	var path string
	if path = Options.GetString(InstallDir); path == "" {
		log.Fatal("Cannot locate InstallDir")
	}
	return path
}

func getReleaseInfo() { // TODO finish this
	if isDevBuild() {
		GIT_REL = "$(git log --format='%h %ci' -1 2>/dev/null | awk '{ print $1\" \"$2\" \"$3 }')"
		GIT_REL_SHORT = "$(git log --format='%h %ci' -1 2>/dev/null | awk '{ print $1 }')"
		REL_DATE = "$(git log --format='%h %ci' -1 2>/dev/null | awk '{ print $2 }')"
	} else {
		REL_DATE = "$(tail -5 \"$0\" | awk '/dirkw Exp/ { print $5 }')"
	}

	// TODO HSTS_MIN = HSTS_MIN * 86400 // correct to seconds
	//HPKP_MIN = HPKP_MIN * 86400 // correct to seconds
	//
	//if MEASURE_TIME_FILE != "" {
	//	MEASURE_TIME = true
	//}

}

// TODO
func isDevBuild() bool {
	return true
}
