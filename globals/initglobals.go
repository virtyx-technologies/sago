package globals

import (
	"github.com/virtyx-technologies/sago/util"
	"log"
	"os"
	"path/filepath"
)

// Initialise global variables
func initGlobals() {
	OpenSSL = findOpensslBinary()
	OpenSslMeta = NewOpenSslMetadata()
	DataDir = findDataDir()
	Targets = getTargets()
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

