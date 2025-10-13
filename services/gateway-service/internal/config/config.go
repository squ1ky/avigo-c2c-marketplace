package config

import (
	"errors"
	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig
	JWT      JWTConfig
	Services ServicesConfig
}

type ServerConfig struct {
	Address string
}

type JWTConfig struct {
	Secret string
}

type ServicesConfig struct {
	UserServiceURL    string
	ListingServiceURL string
}

func LoadConfig() (*Config, error) {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	setDefaults()

	_ = viper.ReadInConfig()

	cfg := &Config{
		Server: ServerConfig{
			Address: viper.GetString("GATEWAY_ADDRESS"),
		},
		JWT: JWTConfig{
			Secret: viper.GetString("JWT_SECRET"),
		},
		Services: ServicesConfig{
			UserServiceURL:    viper.GetString("USER_SERVICE_URL"),
			ListingServiceURL: viper.GetString("LISTING_SERVICE_URL"),
		},
	}

	if err := validateConfig(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

func setDefaults() {
	viper.SetDefault("SERVER_ADDRESS", ":8080")
	viper.SetDefault("USER_SERVICE_URL", "http://user-service:8081")
	viper.SetDefault("LISTING_SERVICE_URL", "http://listing-service:8082")
}

func validateConfig(cfg *Config) error {
	if cfg.JWT.Secret == "" {
		return errors.New("JWT_SECRET is required")
	}
	return nil
}
