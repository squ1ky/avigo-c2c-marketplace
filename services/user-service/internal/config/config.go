package config

import (
	"errors"
	"fmt"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
	Kafka    KafkaConfig
	Cookie   CookieConfig
}

type ServerConfig struct {
	Address string
}

type DatabaseConfig struct {
	URL            string
	MigrationsPath string
}

type JWTConfig struct {
	Secret               string
	AccessTokenDuration  time.Duration
	RefreshTokenDuration time.Duration
}

type KafkaConfig struct {
	Brokers         []string
	TopicUserEvents string
}

type CookieConfig struct {
	Domain   string
	Secure   bool
	HTTPOnly bool
	SameSite string
}

func LoadConfig() (*Config, error) {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	accessTokenDuration, err := parseDuration("ACCESS_TOKEN_DURATION")
	if err != nil {
		return nil, err
	}

	refreshTokenDuration, err := parseDuration("REFRESH_TOKEN_DURATION")
	if err != nil {
		return nil, err
	}

	cfg := &Config{
		Server: ServerConfig{
			Address: viper.GetString("SERVER_ADDRESS"),
		},
		Database: DatabaseConfig{
			URL:            viper.GetString("DATABASE_URL"),
			MigrationsPath: viper.GetString("MIGRATIONS_PATH"),
		},
		JWT: JWTConfig{
			Secret:               viper.GetString("JWT_SECRET"),
			AccessTokenDuration:  accessTokenDuration,
			RefreshTokenDuration: refreshTokenDuration,
		},
		Kafka: KafkaConfig{
			Brokers:         []string{viper.GetString("KAFKA_BROKERS")},
			TopicUserEvents: viper.GetString("KAFKA_TOPIC_USER_EVENTS"),
		},
		Cookie: CookieConfig{
			Domain:   viper.GetString("COOKIE_DOMAIN"),
			Secure:   viper.GetBool("COOKIE_SECURE"),
			HTTPOnly: viper.GetBool("COOKIE_HTTP_ONLY"),
			SameSite: viper.GetString("COOKIE_SAME_SITE"),
		},
	}

	if err := validateConfig(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

func parseDuration(key string) (time.Duration, error) {
	durationStr := viper.GetString(key)
	if durationStr == "" {
		return 0, fmt.Errorf("%s is not set", key)
	}

	duration, err := time.ParseDuration(durationStr)
	if err != nil {
		return 0, fmt.Errorf("invalid %s format: %w", key, err)
	}

	return duration, nil
}

func validateConfig(cfg *Config) error {
	if cfg.Database.URL == "" {
		return errors.New("DATABASE_URL is not set")
	}
	if cfg.Database.MigrationsPath == "" {
		return errors.New("MIGRATIONS_PATH is not set")
	}
	if cfg.JWT.Secret == "" {
		return errors.New("JWT_SECRET is not set")
	}
	if len(cfg.Kafka.Brokers) == 0 || cfg.Kafka.Brokers[0] == "" {
		return errors.New("KAFKA_BROKERS is not set")
	}

	return nil
}
