package globals

import (
	"bufio"
	"github.com/virtyx-technologies/sago/util"
	"io"
	"log"
	"os/exec"
	"regexp"
	"strings"
)

type OpenSslMetadata struct {
	Name        string
	Version     string
	VerMajor    int
	VerMinor    int
	VerAppendix string
	Platform    string
	BuildDate   string

	// Capabilities
	HasSsl2         bool
	HasSsl3         bool
	HasTls13        bool
	HasNoSsl2       bool
	HasNoServerName bool
	HasCipherSuites bool
	HasComp         bool
	HasNoComp       bool

	HasAlpn         bool
	HasNpn          bool
	HasFallbackScsv bool
	HasProxy        bool
	HasXmpp         bool

	HasPostgres bool
	HasMysql    bool
	HasLmtp     bool
	HasNntp     bool
	HasIrc      bool

	HasPkey   bool
	HasPkUtil bool

	HasChaCha20  bool
	HasAes128Gcm bool
	HasAes256Gcm bool
}

func NewOpenSslMetadata() *OpenSslMetadata {
	osm := &OpenSslMetadata{}
	osm.getVersions()
	osm.getCapabilities()
	return osm
}

func (osm *OpenSslMetadata) getVersions() {
	cmdOut, err := exec.Command(OpenSSL, "version", "-a").Output()
	checkErr("Failed to get openssl metadata", err)
	scanner := bufio.NewScanner(strings.NewReader(string(cmdOut)))
	var lines [3]string
	for i := range lines {
		scanner.Scan()
		lines[i] = scanner.Text()
	}

	words := strings.Fields(lines[0])
	osm.Name = words[0]
	osm.Version = words[1]
	re := regexp.MustCompile(`^(\d+)\.(\d+)\.(.+)$`)
	tokens := re.FindStringSubmatch(osm.Version)
	osm.VerMajor = util.Atoi(tokens[1], -1)
	osm.VerMinor = util.Atoi(tokens[2], -1)
	osm.VerAppendix = tokens[3]

	words = strings.Fields(lines[1])
	osm.BuildDate = strings.Join(words[1:], " ")

	words = strings.Fields(lines[2])
	osm.Platform = words[1]
}

func checkErr(msg string, err error) {
	if err != nil {
		log.Fatal(msg, err.Error())
	}
}

// TODO : This method runs openssl many times just to gather information about openssl itself.
// Maybe find a way to gather this information as required, rather than during startup
func (osm *OpenSslMetadata) getCapabilities() {
	cmdOut, err := exec.Command(OpenSSL, "s_client", "-help").CombinedOutput()
	checkErr("Failed to run openssl -help", err)
	s := string(cmdOut)
	osm.HasAlpn         = strings.Contains(s, "-alpn")
	osm.HasNpn          = strings.Contains(s, "-nextprotoneg")
	osm.HasFallbackScsv = strings.Contains(s, "-fallback_scsv")
	osm.HasProxy        = strings.Contains(s, "-proxy")
	osm.HasXmpp         = strings.Contains(s, "-xmpp")

	cmdOut, err = exec.Command(OpenSSL, "s_client", "-starttls", "foo").CombinedOutput()
	checkErr("Failed to run openssl -starttls", err)
	s = string(cmdOut)
	osm.HasAlpn         = strings.Contains(s, "postgres")
	osm.HasNpn          = strings.Contains(s, "mysql")
	osm.HasFallbackScsv = strings.Contains(s, "lmtp")
	osm.HasProxy        = strings.Contains(s, "nntp")
	osm.HasXmpp         = strings.Contains(s, "irc")

	osm.HasSsl2 = checkOption("-ssl2")
	osm.HasSsl3 = checkOption("-ssl3")
	osm.HasTls13 = checkOption("-tls1_3")
	osm.HasNoSsl2 = checkOption("-no_ssl2")
	osm.HasNoServerName = checkOption("-noservername")
	osm.HasCipherSuites = checkOption("-ciphersuites")
	osm.HasComp = checkOption("-comp")
	osm.HasNoComp = checkOption("-no_comp")

	_, err = exec.Command(OpenSSL, "pkey", "-help").CombinedOutput()
	osm.HasPkey = err == nil

	_, err = exec.Command(OpenSSL, "pkeyutl").CombinedOutput()
	osm.HasPkUtil = err == nil

	osm.HasChaCha20 = checkEncoding("-chacha20", "12345678901234567890123456789012", "01000000123456789012345678901234")
	osm.HasAes128Gcm = checkEncoding("-aes-128-gcm", "0123456789abcdef0123456789abcdef", "0123456789abcdef01234567")
	osm.HasAes256Gcm = checkEncoding("-aes-256-gcm", "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef", "0123456789abcdef01234567")
}

func checkEncoding(flag, k, iv string) bool {
	cmd := exec.Command(OpenSSL, "enc", flag, "-K", k, "-iv", iv)
	stdin, _ := cmd.StdinPipe()
	go func() {
		defer stdin.Close()
		io.WriteString(stdin, "test")
	}()

	_, err := cmd.CombinedOutput()
	return err == nil
}

func checkOption(opt string) bool {
	cmdOut, _ := exec.Command(OpenSSL, "s_client", opt, "-connect",  "x").CombinedOutput()
	return !strings.Contains(string(cmdOut), "unknown option")
}
