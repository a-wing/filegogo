package cmd

import (
	"context"

	"filegogo/client"
	"filegogo/client/qrcode"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	rootCmd.AddCommand(sendCmd)
}

var sendCmd = &cobra.Command{
	Use:   "send <file>",
	Short: "Send File",
	Long:  `Send File or path`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		qrconfig := &qrcode.Config{}
		viper.Sub("qrcode").Unmarshal(qrconfig)

		config := &client.ClientConfig{
			QRcodeConfig: qrconfig,
		}
		viper.Sub("send").Unmarshal(config)

		client := &client.Client{
			Config: config,
			Server: viper.GetString("address"),
		}
		client.Send(context.Background(), args)
	},
}
