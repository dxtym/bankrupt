package utils

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Driver            string        `mapstructure:"DRIVER"`
	Source            string        `mapstructure:"SOURCE"`
	Address           string        `mapstructure:"ADDRESS"`
	TokenSymmetricKey string        `mapstructure:"TOKEN_SYMMETRIC_KEY"`
	TokenDuration     time.Duration `mapstructure:"TOKEN_DURATION"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()
	if err = viper.ReadInConfig(); err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
