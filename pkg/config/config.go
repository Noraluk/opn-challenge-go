package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	Omise omise `mapstructure:"omise"`
}

type omise struct {
	PublicKey string `mapstructure:"public_key"`
	SecretKey string `mapstructure:"secret_key"`
}

var config Config

func Init() error {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	if err := viper.Unmarshal(&config); err != nil {
		return err
	}

	return nil
}

func GetConfig() Config {
	return config
}
