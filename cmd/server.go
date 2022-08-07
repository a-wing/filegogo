package cmd

import (
	"filegogo/server"
	"filegogo/server/httpd"

	"github.com/BurntSushi/toml"
	log "github.com/sirupsen/logrus"
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
			},
		}

		if cPath, err := cmd.Flags().GetString("config"); err == nil {
			toml.DecodeFile(cPath, config)
			log.Debugln(config)
		}

		server.Run(config)
	},
}
