package cmd

import (
	"filegogo/server"
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
		config := &server.Config{
			Http: &httpd.Config{
				Listen:    "0.0.0.0:8080",
				RoomAlive: 1024,
				RoomCount: 10000,
				SubFolder: "",
			},
		}

		loadConfig(config)
		server.Run(config)
	},
}
