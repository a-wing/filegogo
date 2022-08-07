package cmd

import (
	"fmt"

	"filegogo/version"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of gosd",
	Long:  `All software has versions. This is filegogo's`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(version.Version + " " + version.Date)
	},
}
