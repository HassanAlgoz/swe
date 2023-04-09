package config

import (
	"log"
	"strings"

	"github.com/spf13/viper"
)

var configPath = "etc/app/"

func SetupConfig() {
	viper.SetConfigType("yaml")
	viper.AddConfigPath(configPath)

	viper.SetConfigName("conf") // filename (without extension)
	err := viper.MergeInConfig()
	if err != nil {
		log.Printf("cannot load config file: %s \n", err)
	}

	viper.AutomaticEnv()
	viper.SetEnvPrefix("app")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
}
