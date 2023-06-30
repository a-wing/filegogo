package cmd

import (
	"encoding/json"
	"os"

	"github.com/BurntSushi/toml"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	// Used for flags.
	cfgFile = ""
	rootCmd = &cobra.Command{
		Use:   "filegogo",
		Short: "Filegogo",
		Long: `Filegogo is a p2p file transport tool
						https://github.com/a-wing/filegogo`,
		PersistentPreRun: func(c *cobra.Command, args []string) {
			if verbose, _ := c.Flags().GetBool("verbose"); verbose {
				log.SetReportCaller(true)
				log.SetLevel(log.DebugLevel)
			}
		},
	}
)

func Execute() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is /etc/filegogo.toml)")
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "verbose output")
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

// Exists determine whether the file exists
func Exists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

// Priority:
// - --config <filegogo.toml>
// - ./
// - /etc/
func loadConfig(cfg interface{}) {
	const name = "filegogo.toml"
	if f, err := os.Open("/etc/" + name); err == nil {
		log.Info("read config file: /etc/", name)
		toml.DecodeReader(f, cfg)
		f.Close()
	}

	if f, err := os.Open(name); err == nil {
		log.Info("read config file: ", name)
		toml.DecodeReader(f, cfg)
		f.Close()
	}

	if Exists(cfgFile) {
		if f, err := os.Open(cfgFile); err == nil {
			log.Info("read config file: ", cfgFile)
			toml.DecodeReader(f, cfg)
			f.Close()
		}
	}

	data, _ := json.Marshal(cfg)

	log.Debugf("%+s\n", data)
}
