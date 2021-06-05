package cmd

import (
	"filegogo/client"
	"filegogo/client/qrcode"

	"github.com/pion/webrtc/v3"

	"github.com/spf13/viper"
)

func getClient() *client.Client {
	// QRcodeConfig
	qrconfig := &qrcode.Config{}
	if viper.IsSet("qrcode") {
		viper.Sub("qrcode").Unmarshal(qrconfig)
	}

	// IcsServers
	iceservers := &webrtc.Configuration{}
	viper.Unmarshal(iceservers)

	config := &client.ClientConfig{
		ShowQRcode:   viper.GetBool("show-qrcode"),
		ShowProgress: viper.GetBool("show-progress"),
		IcsServers:   iceservers,
		QRcodeConfig: qrconfig,
	}
	viper.Unmarshal(config)
	c, _ := client.NewClient(config)
	return c
}
