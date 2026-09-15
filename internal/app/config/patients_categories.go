package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	ServiceHost string
	ServicePort int
}

func NewConfig() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath("config")

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	return &Config{
		ServiceHost: viper.GetString("ServiceHost"),
		ServicePort: viper.GetInt("ServicePort"),
	}, nil
}
