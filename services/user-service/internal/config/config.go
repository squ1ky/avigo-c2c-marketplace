package config

import (
	"errors"
	"fmt"
	"github.com/IBM/sarama"
	"log"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig
	GRPC     GRPCConfig
	Database DatabaseConfig
	JWT      JWTConfig
	Auth     AuthConfig
	Kafka    KafkaConfig
	Cookie   CookieConfig
}

type ServerConfig struct {
	Address string
}

type GRPCConfig struct {
	Address               string
	MaxConnectionIdle     time.Duration
	MaxConnectionAge      time.Duration
	MaxConnectionAgeGrace time.Duration
	Time                  time.Duration
	Timeout               time.Duration
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

type AuthConfig struct {
	ConfirmationCodeExpiry time.Duration
}

type KafkaConfig struct {
	Brokers         []string
	TopicUserEvents string
	RequiredAcks    sarama.RequiredAcks
	CompressionType sarama.CompressionCodec
	DialTimeout     time.Duration
	WriteTimeout    time.Duration
	ReadTimeout     time.Duration
	RetryMax        int
	RetryBackoff    time.Duration
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

	setDefaults()

	if err := viper.ReadInConfig(); err != nil {
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if errors.As(err, &configFileNotFoundError) {
			log.Println("No .env file found, using environment variables and defaults")
		} else {
			log.Printf("Error reading config file: %v", err)
		}
	}

	accessTokenDuration, err := parseDuration("ACCESS_TOKEN_DURATION")
	if err != nil {
		return nil, err
	}

	refreshTokenDuration, err := parseDuration("REFRESH_TOKEN_DURATION")
	if err != nil {
		return nil, err
	}

	confirmationCodeExpiry, err := parseDuration("CONFIRMATION_CODE_EXPIRY")
	if err != nil {
		return nil, err
	}

	cfg := &Config{
		Server: ServerConfig{
			Address: viper.GetString("SERVER_ADDRESS"),
		},
		GRPC: GRPCConfig{
			Address:               viper.GetString("GRPC_ADDRESS"),
			MaxConnectionIdle:     viper.GetDuration("GRPC_MAX_CONNECTION_IDLE"),
			MaxConnectionAge:      viper.GetDuration("GRPC_MAX_CONNECTION_AGE"),
			MaxConnectionAgeGrace: viper.GetDuration("GRPC_MAX_CONNECTION_AGE_GRACE"),
			Time:                  viper.GetDuration("GRPC_KEEPALIVE_TIME"),
			Timeout:               viper.GetDuration("GRPC_KEEPALIVE_TIME"),
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
		Auth: AuthConfig{
			ConfirmationCodeExpiry: confirmationCodeExpiry,
		},
		Kafka: KafkaConfig{
			Brokers:         []string{viper.GetString("KAFKA_BROKERS")}, // need parseBrokers() when we have > 1 broker
			TopicUserEvents: viper.GetString("KAFKA_TOPIC_USER_EVENTS"),
			RequiredAcks:    sarama.RequiredAcks(viper.GetInt("KAFKA_REQUIRED_ACKS")),
			CompressionType: sarama.CompressionCodec(viper.GetInt("KAFKA_COMPRESSION")),
			DialTimeout:     viper.GetDuration("KAFKA_DIAL_TIMEOUT"),
			WriteTimeout:    viper.GetDuration("KAFKA_WRITE_TIMEOUT"),
			ReadTimeout:     viper.GetDuration("KAFKA_READ_TIMEOUT"),
			RetryMax:        viper.GetInt("KAFKA_RETRY_MAX"),
			RetryBackoff:    viper.GetDuration("KAFKA_RETRY_BACKOFF"),
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

func setDefaults() {
	// HTTP Server
	viper.SetDefault("SERVER_ADDRESS", ":8081")

	// gRPC Server
	viper.SetDefault("GRPC_ADDRESS", ":50051")
	viper.SetDefault("GRPC_MAX_CONNECTION_IDLE", 5*time.Minute)
	viper.SetDefault("GRPC_MAX_CONNECTION_AGE", 5*time.Minute)
	viper.SetDefault("GRPC_MAX_CONNECTION_AGE_GRACE", 1*time.Minute)
	viper.SetDefault("GRPC_KEEPALIVE_TIME", 2*time.Hour)
	viper.SetDefault("GRPC_KEEPALIVE_TIMEOUT", 20*time.Second)

	// JWT
	viper.SetDefault("ACCESS_TOKEN_DURATION", "15m")
	viper.SetDefault("REFRESH_TOKEN_DURATION", "168h")
	viper.SetDefault("CONFIRMATION_CODE_EXPIRY", "15m")

	// Cookie
	viper.SetDefault("COOKIE_HTTP_ONLY", true)
	viper.SetDefault("COOKIE_SECURE", false)
	viper.SetDefault("COOKIE_SAME_SITE", "lax")

	// Kafka
	viper.SetDefault("KAFKA_REQUIRED_ACKS", int(sarama.WaitForAll))
	viper.SetDefault("KAFKA_COMPRESSION", int(sarama.CompressionSnappy))
	viper.SetDefault("KAFKA_DIAL_TIMEOUT", 10*time.Second)
	viper.SetDefault("KAFKA_WRITE_TIMEOUT", 10*time.Second)
	viper.SetDefault("KAFKA_READ_TIMEOUT", 10*time.Second)
	viper.SetDefault("KAFKA_RETRY_MAX", 3)
	viper.SetDefault("KAFKA_RETRY_BACKOFF", 100*time.Millisecond)

	// Database
	viper.SetDefault("MIGRATIONS_PATH", "internal/db/migrations")
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
	if cfg.GRPC.Address == "" {
		return errors.New("GRPC_ADDRESS is not set")
	}

	return nil
}
