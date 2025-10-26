package config

import (
	"errors"
	"strings"
	"time"

	"github.com/spf13/viper"
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
			AllowOrigins:     parseCSV(viper.GetString("CORS_ALLOW_ORIGINS")),
			AllowMethods:     parseCSV(viper.GetString("CORS_ALLOW_METHODS")),
			AllowHeaders:     parseCSV(viper.GetString("CORS_ALLOW_HEADERS")),
			ExposeHeaders:    parseCSV(viper.GetString("CORS_EXPOSE_HEADERS")),
			AllowCredentials: viper.GetBool("CORS_ALLOW_CREDENTIALS"),
			MaxAge:           viper.GetDuration("CORS_MAX_AGE"),
		},
	}

	if err := validateConfig(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

// parseCSV parses string "a,b,c" to []string{"a", "b", "c"}
func parseCSV(s string) []string {
	if s == "" {
		return []string{}
	}
	parts := strings.Split(s, ",")
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

func setDefaults() {
	viper.SetDefault("GATEWAY_ADDRESS", ":8080")
	viper.SetDefault("USER_SERVICE_URL", "http://user-service:8081")
	viper.SetDefault("LISTING_SERVICE_URL", "http://listing-service:8082")

	// CORS defaults
	viper.SetDefault("CORS_ALLOW_ORIGINS", "http://localhost:3000")
	viper.SetDefault("CORS_ALLOW_METHODS", "GET,POST,PUT,PATCH,DELETE,OPTIONS")
	viper.SetDefault("CORS_ALLOW_HEADERS", "Origin,Content-Type,Accept,Authorization,Cookie")
	viper.SetDefault("CORS_EXPOSE_HEADERS", "Content-Length,Set-Cookie")
	viper.SetDefault("CORS_ALLOW_CREDENTIALS", true)
	viper.SetDefault("CORS_MAX_AGE", 12*time.Hour)
}

func validateConfig(cfg *Config) error {
	if cfg.JWT.Secret == "" {
		return errors.New("JWT_SECRET is required")
	}
	return nil
}
