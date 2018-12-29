package options

import (
	"fmt"
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
)

const (
	configFileName = "sago.json"
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
	fs.Bool("version", false, "Print version & exit ")
	fs.Bool("help", false, "Display help & exit")
	// Configuration
	fs.String("log", "info", "Level of agent logging ")
	fs.Bool("system-log", true, "Use system logging ")
}

func PrintDefaults() {
	Flags.PrintDefaults()
}

// Initialise global variables
func initGlobals() {
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
