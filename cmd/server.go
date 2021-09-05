package cmd

import (
	"filegogo/server"

	"github.com/pion/webrtc/v3"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	rootCmd.AddCommand(serverCmd)

	serverCmd.Flags().StringP("listen", "l", "0.0.0.0:8033", "set server listen address and port")
}

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "server",
	Long:  `websocket broker server`,
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		// IcsServers
		iceservers := &webrtc.Configuration{}
		viper.Unmarshal(iceservers)

		server.Run(&server.Config{
			Server:     viper.GetString("listen"),
			IcsServers: iceservers,
		})
	},
}
