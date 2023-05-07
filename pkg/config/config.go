package config

import (
	"log"
	"strings"

	"github.com/spf13/viper"
)

var configPath = "etc/"

func SetupConfig() {
	viper.SetConfigType("yaml")
	viper.AddConfigPath(configPath)

	viper.SetConfigName("config") // filename (without extension)
	err := viper.MergeInConfig()
	if err != nil {
		log.Fatalf("cannot load config file: %s \n", err)
	}

	viper.AutomaticEnv()
	viper.SetEnvPrefix("app")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.WatchConfig() // This makes feature flagging possible at runtime (see the middleware)
}

func SetupTestConfig() {
	viper.SetConfigType("yaml")
	viper.AddConfigPath(configPath)

	viper.SetConfigName("config") // filename (without extension)
	err := viper.MergeInConfig()
	if err != nil {
		log.Fatalf("cannot load config file: %s \n", err)
	}

	viper.AutomaticEnv()
	viper.SetEnvPrefix("app")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.Set("app.env", "test")
}
