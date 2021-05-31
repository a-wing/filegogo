package cmd

import (
	"context"
	"fmt"

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
	Long:  `All software has versions. This is Hugo's`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(args)
		fgg := &libfgg.Fgg{
			Server: viper.GetString("address"),
		}
		fgg.Recv(context.Background(), args)
	},
}
