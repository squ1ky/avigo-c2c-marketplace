package config

import (
	"errors"

	"github.com/spf13/viper"
)

type Config struct {
	ServerAddress  string
	DatabaseURL    string
	MigrationsPath string
	JWTSecret      string
	LogLevel       string
}

func LoadConfig() (*Config, error) {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	cfg := &Config{
		ServerAddress:  viper.GetString("SERVER_ADDRESS"),
		DatabaseURL:    viper.GetString("DATABASE_URL"),
		MigrationsPath: viper.GetString("MIGRATIONS_PATH"),
		JWTSecret:      viper.GetString("JWT_SECRET"),
		LogLevel:       viper.GetString("LOG_LEVEL"),
	}

	if cfg.DatabaseURL == "" {
		return nil, errors.New("DATABASE_URL is not set")
	}
	if cfg.MigrationsPath == "" {
		return nil, errors.New("MIGRATIONS_PATH is not set")
	}
	if cfg.JWTSecret == "" {
		return nil, errors.New("JWT_SECRET is not set")
	}
	if cfg.ServerAddress == "" {
		cfg.ServerAddress = ":8081"
	}

	return cfg, nil
}
