package util

import (
	"fmt"

	"time"

	"github.com/spf13/viper"
)

type Config struct {
	DATABASE_DRVIER   string `mapstructure:"DATABASE_DRVIER"`
	DATABASE_HOST     string `mapstructure:"DATABASE_HOST"`
	DATABASE_NAME     string `mapstructure:"DATABASE_NAME"`
	DATABASE_PASSWORD string `mapstructure:"DATABASE_PASSWORD"`
	DATABASE_PORT     string `mapstructure:"DATABASE_PORT"`
	DATABASE_USERNAME string `mapstructure:"DATABASE_USERNAME"`

	SERVER_API_URL        string        `mapstructure:"SERVER_API_URL"`
	SYMMETRIC_KEY         string        `mapstructure:"SYMMETRIC_KEY"`
	ACCESS_TOKEN_DURATION time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
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

func DbURL(config Config) string {
	return fmt.Sprintf("%s://%s:%s@%s:%s/%s?sslmode=disable",
		config.DATABASE_DRVIER,
		config.DATABASE_USERNAME,
		config.DATABASE_PASSWORD,
		config.DATABASE_HOST,
		config.DATABASE_PORT,
		config.DATABASE_NAME)
}
