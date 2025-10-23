package config

import (
	"errors"
	"log"
	"time"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Kafka    KafkaConfig
	Email    EmailConfig
}

type ServerConfig struct {
	Address  string
	LogLevel string
}

type DatabaseConfig struct {
	URL            string
	MigrationsPath string
}

type KafkaConfig struct {
	Brokers         []string
	TopicUserEvents string
	GroupID         string
	MinBytes        int
	MaxBytes        int
}

type EmailConfig struct {
	Provider      string // "smtp"
	SMTPHost      string
	SMTPPort      int
	SMTPUsername  string
	SMTPPassword  string
	FromName      string
	FromAddress   string
	TemplatesPath string
}

func LoadConfig() (*Config, error) {

	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	setDefaults()

	cfg := &Config{
		Server: ServerConfig{
			Address:  viper.GetString("SERVER_ADDRESS"),
			LogLevel: viper.GetString("LOG_LEVEL"),
		},
		Database: DatabaseConfig{
			URL:            viper.GetString("DATABASE_URL"),
			MigrationsPath: viper.GetString("MIGRATIONS_PATH"),
		},
		Kafka: KafkaConfig{
			Brokers:         []string{viper.GetString("KAFKA_BROKERS")},
			TopicUserEvents: viper.GetString("KAFKA_TOPIC_USER_EVENTS"),
			GroupID:         viper.GetString("KAFKA_GROUP_ID"),
			MinBytes:        viper.GetInt("KAFKA_MIN_BYTES"),
			MaxBytes:        viper.GetInt("KAFKA_MAX_BYTES"),
		},
		Email: EmailConfig{
			Provider:      viper.GetString("EMAIL_PROVIDER"),
			SMTPHost:      viper.GetString("EMAIL_SMTP_HOST"),
			SMTPPort:      viper.GetInt("EMAIL_SMTP_PORT"),
			SMTPUsername:  viper.GetString("EMAIL_SMTP_USERNAME"),
			SMTPPassword:  viper.GetString("EMAIL_SMTP_PASSWORD"),
			FromName:      viper.GetString("EMAIL_FROM_NAME"),
			FromAddress:   viper.GetString("EMAIL_FROM_ADDRESS"),
			TemplatesPath: viper.GetString("EMAIL_TEMPLATES_PATH"),
		},
	}

	if err := validateConfig(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

func setDefaults() {
	viper.SetDefault("SERVER_ADDRESS", ":8082")
	viper.SetDefault("LOG_LEVEL", "info")
	viper.SetDefault("MIGRATIONS_PATH", "migrations")
	viper.SetDefault("KAFKA_GROUP_ID", "notification-service")
	viper.SetDefault("KAFKA_DIAL_TIMEOUT", 10*time.Second)
	viper.SetDefault("KAFKA_READ_TIMEOUT", 10*time.Second)
	viper.SetDefault("KAFKA_MIN_BYTES", 10240)    // 10KB
	viper.SetDefault("KAFKA_MAX_BYTES", 10485760) // 10MB
	viper.SetDefault("EMAIL_PROVIDER", "smtp")
	viper.SetDefault("EMAIL_SMTP_PORT", 587)
	viper.SetDefault("EMAIL_TEMPLATES_PATH", "assets/templates/email")
}

func validateConfig(cfg *Config) error {
	if cfg.Database.URL == "" {
		return errors.New("DATABASE_URL is required")
	}

	if len(cfg.Kafka.Brokers) == 0 || cfg.Kafka.Brokers[0] == "" {
		return errors.New("KAFKA_BROKERS is required")
	}

	if cfg.Kafka.TopicUserEvents == "" {
		return errors.New("KAFKA_TOPIC_USER_EVENTS is required")
	}

	if cfg.Email.Provider == "smtp" {
		if cfg.Email.SMTPHost == "" {
			return errors.New("EMAIL_SMTP_HOST is required for SMTP provider")
		}
		if cfg.Email.SMTPUsername == "" {
			return errors.New("EMAIL_SMTP_USERNAME is required for SMTP provider")
		}
		if cfg.Email.SMTPPassword == "" {
			return errors.New("EMAIL_SMTP_PASSWORD is required for SMTP provider")
		}
	}

	if cfg.Email.FromAddress == "" {
		return errors.New("EMAIL_FROM_ADDRESS is required")
	}

	return nil
}
