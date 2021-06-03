package cmd

import (
	"context"

	"filegogo/client"
	"filegogo/client/qrcode"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	rootCmd.AddCommand(recvCmd)
}

var recvCmd = &cobra.Command{
	Use:   "recv <file>",
	Short: "Recv File",
	Long:  `Recv File or path. if not set, use raw filename`,
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		qrconfig := &qrcode.Config{}
		viper.Sub("qrcode").Unmarshal(qrconfig)

		config := &client.ClientConfig{
			QRcodeConfig: qrconfig,
		}
		viper.Sub("recv").Unmarshal(config)

		client := &client.Client{
			Config: config,
			Server: viper.GetString("address"),
		}
		client.Recv(context.Background(), args)
	},
}
