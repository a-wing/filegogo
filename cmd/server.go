package cmd

import (
	"os"

	"filegogo/server"
	"filegogo/server/config"
	"filegogo/server/httpd"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(serverCmd)
}

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "websocket broker server",
	Long:  `webapp, websocket, iceServer Server`,
	Run: func(cmd *cobra.Command, args []string) {
		config := &config.Config{
			Http: &httpd.Config{
				Listen:      "0.0.0.0:8080",
				PathPrefix:  "",
				StoragePath: "data",
			},
		}

		loadConfig(config)

		// Override `Http.Listen` to 0.0.0.0:$PORT (Automatic configuration for PaaS).
		if port := os.Getenv("PORT"); port != "" {
			config.Http.Listen = ":" + port
		}

		server.Run(config)
	},
}
