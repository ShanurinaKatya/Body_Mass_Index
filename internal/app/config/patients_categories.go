package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/spf13/viper"
)

type Config struct {
	ServiceHost string
	ServicePort int

	JWT   JWTConfig
	Redis RedisConfig
}

type JWTConfig struct {
	SigningMethod jwt.SigningMethod
	Token         string
	ExpiresIn     time.Duration
}

type RedisConfig struct {
	Host        string
	Port        int
	User        string
	Password    string
	DialTimeout time.Duration
	ReadTimeout time.Duration
}

func NewConfig() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath("config")

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	expiresIn, err := time.ParseDuration(envOr("JWT_EXPIRES_IN", "1h"))
	if err != nil {
		return nil, fmt.Errorf("JWT_EXPIRES_IN must be a duration: %w", err)
	}

	redisPort, err := strconv.Atoi(envOr("REDIS_PORT", "6379"))
	if err != nil {
		return nil, fmt.Errorf("REDIS_PORT must be an int: %w", err)
	}

	return &Config{
		ServiceHost: viper.GetString("ServiceHost"),
		ServicePort: viper.GetInt("ServicePort"),
		JWT: JWTConfig{
			SigningMethod: jwt.SigningMethodHS256,
			Token:         envOr("JWT_TOKEN", "test"),
			ExpiresIn:     expiresIn,
		},
		Redis: RedisConfig{
			Host:        envOr("REDIS_HOST", "0.0.0.0"),
			Port:        redisPort,
			User:        envOr("REDIS_USER", ""),
			Password:    envOr("REDIS_PASSWORD", ""),
			DialTimeout: 10 * time.Second,
			ReadTimeout: 10 * time.Second,
		},
	}, nil
}

func envOr(key, def string) string {
	if v, ok := os.LookupEnv(key); ok {
		return v
	}
	return def
}
