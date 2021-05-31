package cmd

import (
	"context"
	"fmt"

	"filegogo/libfgg"

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
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(args)
		fgg := &libfgg.Fgg{
			Server: viper.GetString("address"),
		}
		fgg.Send(context.Background(), args)
	},
}
