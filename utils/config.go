package utils

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Environment       string        `mapstructure:"ENVIRONMENT"`
	Driver            string        `mapstructure:"DRIVER"`
	Source            string        `mapstructure:"SOURCE"`
	MigrateURL        string        `mapstructure:"MIGRATE_URL"`
	HTTPAddress       string        `mapstructure:"HTTP_ADDRESS"`
	GRPCAddress       string        `mapstructure:"GRPC_ADDRESS"`
	TokenSymmetricKey string        `mapstructure:"TOKEN_SYMMETRIC_KEY"`
	TokenDuration     time.Duration `mapstructure:"TOKEN_DURATION"`
	RefreshDuration   time.Duration `mapstructure:"REFRESH_DURATION"`
	RedisAddress      string        `mapstructure:"REDIS_ADDRESS"`
	SenderName 	  	  string        `mapstructure:"SENDER_NAME"`
	SenderAddress     string        `mapstructure:"SENDER_ADDRESS"`
	SenderPassword    string        `mapstructure:"SENDER_PASSWORD"`
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
