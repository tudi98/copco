package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "copco",
	Short: "Competitive Programming Companion",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	home, err := homedir.Dir()
	if err != nil {
		log.Fatal(err)
	}

	viper.AddConfigPath(home)
	viper.SetConfigName(".copco")
	viper.SetConfigType("yaml")

	viper.SetDefault("COPCO_COMPILE_COMMAND", map[string]string{
		"command": "g++",
		"args":    "-o main main.cpp",
	})
	viper.SetDefault("COPCO_RUN_COMMAND", map[string]string{
		"command": "./main",
		"args":    "",
	})
	viper.SetDefault("COPCO_PATH", home+string(os.PathSeparator)+"copco")
	viper.SetDefault("COPCO_TEMPLATE", home+string(os.PathSeparator)+"copco"+string(os.PathSeparator)+"template.cpp")
	viper.SetDefault("COPCO_CUSTOM_CODE", home+string(os.PathSeparator)+"copco"+string(os.PathSeparator)+"custom_code")

	_ = viper.SafeWriteConfig()
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal(err)
	}
}
