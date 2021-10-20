package cmd

import (
	"context"

	"github.com/spf13/cobra"
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
		getClient().Send(context.Background(), args)
	},
}
