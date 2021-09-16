package cmd

import (
	"filegogo/server"
	"filegogo/server/httpd"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	// Used for flags.
	cfgFile string

	rootCmd = &cobra.Command{
		Use:   "server",
		Short: "server",
		Long:  `websocket broker server`,
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			config := &server.Config{
				Http: &httpd.Config{
					Listen:    "0.0.0.0:8080",
					RoomAlive: 1024,
					RoomCount: 10000,
				},
			}
			viper.Unmarshal(config)
			server.Run(config)
		},
	}
)
