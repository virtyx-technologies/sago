package globals

import (
	"fmt"
	"github.com/DataDog/datadog-agent/pkg/util/log"
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	configFileName = "sago-config.yaml"

	// Flag names
   PrintCiphers  = "print-ciphers" // was -V
   DoMassTesting = "DoMassTesting" // was --file
   DoMxAllIps    = "DoMxAllIps"
	XmppHost      = "xmpphost"
	Sneaky = "Sneaky"
	AssumeHttp = "assume-http"
	Bugs  = "bugs"

	Version = "version"
	Help = "help"
	Target = "target"
	OpenSslFile = "openssl-file"
	InstallDir = "install-dir"
)

var logLevel string

func init() {

	v := viper.New()
	v.SetEnvKeyReplacer(strings.NewReplacer("-", "_" ))
	v.SetTypeByDefaultValue(true)
	v.AutomaticEnv()

	// Look for config file path in environment
	configFile = v.GetString("config-file")
	cmd := pflag.NewFlagSet(os.Args[0], pflag.ContinueOnError)
	addFlags(cmd)

	v.BindPFlags(cmd)
	// Ignore errors; CommandLine is set for ExitOnError.
	cmd.Parse(os.Args[1:])
	setLogLevel()

	var err error
	if configFile == "" {
		configDir, err = locateConfigDir()
		if err != nil {
			panic(err)
		}
		configFile = filepath.Join(configDir, configFileName)
	} else {
		s, _ := filepath.Abs(configFile)
		configDir = filepath.Dir(s)
	}
	log.Info("Using config file", configFile)
	v.SetConfigFile(configFile)
	v.ReadInConfig()

	Options = v
	Flags = cmd

	initGlobals()
}

func setLogLevel() {
  lvl, err := logrus.ParseLevel(logLevel)
  fmt.Println(err.Error())
  logrus.SetLevel(lvl)
}

func addFlags(fs *pflag.FlagSet) { // TODO add real flags
	// Actions
	fs.Bool(Version, false, "Print version & exit ")
	fs.Bool(Help, false, "Display help & exit")
	fs.Bool(PrintCiphers, false, "Print local ciphers & exit")
	// Configuration
	fs.StringVar(&logLevel, "log", "info", "Level of logging ")
	fs.String(Target, "", "Comma-separated list of IPs and/or Hosts")
	fs.String(OpenSslFile, "", "full path to OpenSSL executable")
	fs.String(XmppHost, "", "Supplies the XML stream 'to-domain' for STARTTLS enabled XMPP")
	fs.Bool(DoMassTesting, false, "TODO")
	fs.Bool(DoMxAllIps, false, "TODO")
	fs.Bool(Sneaky, false, "Use 'sneaky' User Agent")
	fs.Bool(AssumeHttp, false, "in rare cases (WAF, old servers, grumpy SSL) service detection fails. 'True' enforces HTTP checks")
	fs.Bool(Bugs, false, "Use '-bugs' option for openssl, needed for some BIG IP F5")
}

func PrintDefaults() {
	Flags.PrintDefaults()
}

func ConfigDir() string {
	return configDir
}

func ApiKey() string {
	apiKey := os.Getenv("SAGO_API_KEY")
	if apiKey == "" {
		apiKey = Options.GetString("apikey")
	}
	return apiKey
}

func ConfigFile() string {
	return configFile
}


