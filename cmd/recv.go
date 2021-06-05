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
		if viper.IsSet("qrcode") {
			viper.Sub("qrcode").Unmarshal(qrconfig)
		}

		config := &client.ClientConfig{
			ShowQRcode:   viper.GetBool("show-qrcode"),
			ShowProgress: viper.GetBool("show-progress"),
			QRcodeConfig: qrconfig,
		}
		viper.Unmarshal(config)

		c, _ := client.NewClient(config)
		c.Recv(context.Background(), args)
	},
}
