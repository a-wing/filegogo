package cmd

import (
	"fmt"

	"filegogo/server"
	"filegogo/server/httpd"

	homedir "github.com/mitchellh/go-homedir"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	// Used for flags.
	cfgFile     string
	userLicense string

	rootCmd = &cobra.Command{
		Use:   "server",
		Short: "server",
		Long:  `websocket broker server`,
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			config := &server.Config{
				Http: &httpd.Config{
					RoomAlive: 1024,
					RoomCount: 10000,
				},
			}
			viper.Unmarshal(config)
			server.Run(config)
		},
	}
)

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.config/filegogo-server.toml)")
	rootCmd.Flags().StringP("listen", "l", "0.0.0.0:8033", "set server listen address and port")

	rootCmd.PersistentFlags().StringP("level", "", "info", "log level")
	viper.BindPFlags(rootCmd.PersistentFlags())

	viper.SetDefault("server", "localhost:8033")
	viper.SetDefault("level", "info")

	// server
	viper.SetDefault("listen", "0.0.0.0:8033")
	viper.SetDefault("browser", "./config.json")
}

func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		cobra.CheckErr(err)

		viper.SetConfigName("filegogo-server")
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

	//log.SetReportCaller(true)
	if level, err := log.ParseLevel(viper.GetString("level")); err != nil {
		fmt.Println(err)
	} else {
		log.SetLevel(level)
	}
	log.SetFormatter(&log.TextFormatter{
		//FullTimestamp: true,
	})
}