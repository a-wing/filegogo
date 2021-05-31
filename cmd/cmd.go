package cmd

import (
	"fmt"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	// Used for flags.
	cfgFile     string
	userLicense string

	rootCmd = &cobra.Command{
		Use:   "filegogo",
		Short: "a p2p file transfer tool",
		Long:  `A p2p file transfer tool that can be used in the webrtc p2p`,
	}
)

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.config/filegogo.toml)")
	rootCmd.PersistentFlags().StringP("addr", "a", "", "Signal Server Address (default is https://send.22333.fun)")
	rootCmd.PersistentFlags().StringVarP(&userLicense, "license", "l", "", "name of license for the project")
	rootCmd.PersistentFlags().Bool("viper", true, "use Viper for configuration")
	//viper.BindPFlag("author", rootCmd.PersistentFlags().Lookup("author"))
	viper.BindPFlag("address", rootCmd.PersistentFlags().Lookup("addr"))
	viper.BindPFlag("useViper", rootCmd.PersistentFlags().Lookup("viper"))
	//viper.SetDefault("author", "NAME HERE <EMAIL ADDRESS>")
	viper.SetDefault("address", "http://localhost:8033")
	viper.SetDefault("license", "MIT")
}

func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		cobra.CheckErr(err)

		viper.SetConfigName("filegogo")
		viper.SetConfigType("toml")

		viper.AddConfigPath(".")
		viper.AddConfigPath(home)
		viper.AddConfigPath(home + "/.config/")
		viper.AddConfigPath("/etc/filegogo/")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
