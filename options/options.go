package options

import (
	"fmt"
	"github.com/virtyx-technologies/sago/util"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var (
	Options    *viper.Viper
	Flags      *pflag.FlagSet
	configFile string
	configDir  string

	OpenSSL, DataDir string
)


const (
	configFileName = "sago-config.yaml"

	// Flag names
   DoDisplayOnly = "DoDisplayOnly"
   DoMassTesting = "DoMassTesting"
   DoMxAllIps    = "DoMxAllIps"

	Version = "version"
	Help = "help"
	Target = "target"
	OpenSslFile = "openssl-file"
	InstallDir = "install-dir"
)

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
	fmt.Println("Using config file", configFile)
	v.SetConfigFile(configFile)
	v.ReadInConfig()

	Options = v
	Flags = cmd

	initGlobals()
}

func addFlags(fs *pflag.FlagSet) { // TODO add real flags
	// Actions
	fs.Bool(Version, false, "Print version & exit ")
	fs.Bool(Help, false, "Display help & exit")
	// Configuration
	fs.String(Target, "", "Comma-separated list of IPs and/or Hosts")   // TODO
	fs.String(OpenSslFile, "", "full path to OpenSSL executable")   // TODO
	fs.String("log", "info", "Level of logging ")   // TODO
	fs.Bool(DoDisplayOnly, false, "TODO")
	fs.Bool(DoMassTesting, false, "TODO")
	fs.Bool(DoMxAllIps, false, "TODO")
}

func PrintDefaults() {
	Flags.PrintDefaults()
}

// Initialise global variables
func initGlobals() {
	OpenSSL = findOpensslBinary()
	DataDir = findDataDir()
	// globals.InternalMetrics = Options.GetBool(ParamInternalMetrics)
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



