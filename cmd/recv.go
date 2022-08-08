package cmd

import (
	"context"

	"filegogo/client"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(recvCmd)
	recvCmd.Flags().StringP("share", "s", "https://send.22333.fun", "Share Link")
}

var recvCmd = &cobra.Command{
	Use:   "recv",
	Short: "recv <file>",
	Run: func(cmd *cobra.Command, args []string) {
		config := client.DefaultConfig
		loadConfig(config)

		if share, err := cmd.Flags().GetString("share"); err == nil {
			config.Server = share
		}

		cli, err := client.NewClient(config)
		if err != nil {
			panic(err)
		}
		cli.Recv(context.Background(), args)
	},
}
