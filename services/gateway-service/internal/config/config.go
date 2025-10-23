package config

import (
	"errors"
	"github.com/spf13/viper"
	"time"
)

type Config struct {
	Server   ServerConfig
	JWT      JWTConfig
	Services ServicesConfig
	CORS     CORSConfig
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

type CORSConfig struct {
	AllowOrigins     []string
	AllowMethods     []string
	AllowHeaders     []string
	ExposeHeaders    []string
	AllowCredentials bool
	MaxAge           time.Duration
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
		CORS: CORSConfig{
			AllowOrigins:     viper.GetStringSlice("CORS_ALLOW_ORIGINS"),
			AllowMethods:     viper.GetStringSlice("CORS_ALLOW_METHODS"),
			AllowHeaders:     viper.GetStringSlice("CORS_ALLOW_HEADERS"),
			ExposeHeaders:    viper.GetStringSlice("CORS_EXPOSE_HEADERS"),
			AllowCredentials: viper.GetBool("CORS_ALLOW_CREDENTIALS"),
			MaxAge:           viper.GetDuration("CORS_MAX_AGE"),
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

	// CORS defaults
	viper.SetDefault("CORS_ALLOW_ORIGINS", []string{"http://localhost:3000", "http://localhost:8080"})
	viper.SetDefault("CORS_ALLOW_METHODS", []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"})
	viper.SetDefault("CORS_ALLOW_HEADERS", []string{"Origin", "Content-Type", "Accept", "Authorization", "Cookie"})
	viper.SetDefault("CORS_EXPOSE_HEADERS", []string{"Content-Length", "Set-Cookie"})
	viper.SetDefault("CORS_ALLOW_CREDENTIALS", true)
	viper.SetDefault("CORS_MAX_AGE", 12*time.Hour)
}

func validateConfig(cfg *Config) error {
	if cfg.JWT.Secret == "" {
		return errors.New("JWT_SECRET is required")
	}
	return nil
}
