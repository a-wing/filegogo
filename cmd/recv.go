package cmd

import (
	"context"

	"filegogo/libfgg"

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
		fgg := &libfgg.Fgg{
			Server: viper.GetString("address"),
		}
		fgg.Recv(context.Background(), args)
	},
}
