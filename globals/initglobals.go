package globals

import (
	"bufio"
	"github.com/virtyx-technologies/sago/util"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

// Initialise global variables
func initGlobals() {
	OpenSSL = findOpensslBinary()
	OpenSslMeta = findOpenSslMetadata()
	DataDir = findDataDir()
	// globals.InternalMetrics = Options.GetBool(ParamInternalMetrics)
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

type OpenSslMetadata struct {
	Name string
	Version string
	VerMajor string
	VerMinor string
	VerAppendix string
	Platform string
	BuildDate string
}

func findOpenSslMetadata() *OpenSslMetadata {
	cmdOut, err := exec.Command("openssl", "-a").Output()
	if err != nil {
		log.Fatal("Failed to get openssl metadata")
	}
	scanner := bufio.NewScanner(strings.NewReader(string(cmdOut)))
	var lines [3]string
	for i := range lines {
		scanner.Scan()
		lines[i] = scanner.Text()
	}
	osm := &OpenSslMetadata{}

	words := strings.Fields(lines[0])
	osm.Name = words[0]
	osm.Version = words[1]
	re := regexp.MustCompile(`^(\w+)\.(\w+)\.(.+)$`)
	tokens := re.FindStringSubmatch(osm.Version)
	osm.VerMajor = tokens[1]
	osm.VerMinor = tokens[2]
	osm.VerAppendix = tokens[3]

	words = strings.Fields(lines[1])
	osm.BuildDate = strings.Join(words[1:], " ")

	words = strings.Fields(lines[2])
	osm.Platform = words[1]

	return osm;
}

