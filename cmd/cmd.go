package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "filegogo",
	Short: "Filegogo",
	Long: `Filegogo is a p2p file transport tool
						https://github.com/a-wing/filegogo`,
	PreRun: func(c *cobra.Command, args []string) {
		if verbose, _ := c.Flags().GetBool("verbose"); verbose {
			log.SetReportCaller(true)
			log.SetLevel(log.DebugLevel)
		}
	},
}

func Execute() {
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "verbose output")
	rootCmd.PersistentFlags().StringP("config", "c", "filegogo.toml", "Load configuration from `FILE`")
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
