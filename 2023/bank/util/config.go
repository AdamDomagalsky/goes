package util

import (
	"github.com/spf13/viper"
)

type Config struct {
	DBDriver       string `mapstructure:"DATABASE_DRVIER"`
	DATABASE_URL   string `mapstructure:"DATABASE_URL"`
	SERVER_API_URL string `mapstructure:"SERVER_API_URL"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName(".env")
	viper.SetConfigType("env")

	viper.AutomaticEnv()
	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
