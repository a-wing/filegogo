package cmd

import (
	"context"

	"github.com/spf13/cobra"
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
		getClient().Recv(context.Background(), args)
	},
}
