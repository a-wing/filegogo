package cmd

import (
	"encoding/json"
	"os"

	"github.com/BurntSushi/toml"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
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

func Execute() {
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "verbose output")
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

// Priority:
// - ./
// - /etc/
func loadConfig(cfg interface{}) {
	const name = "filegogo.toml"
	if f, err := os.Open("/etc/" + name); err == nil {
		toml.DecodeReader(f, cfg)
	}

	if f, err := os.Open(name); err == nil {
		toml.DecodeReader(f, cfg)
	}

	data, _ := json.Marshal(cfg)

	log.Debugf("%+s\n", data)
}
