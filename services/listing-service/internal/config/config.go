package config

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Server      ServerConfig   `mapstructure:"server"`
	UserService UserService    `mapstructure:"user_service"`
	Postgres    PostgresConfig `mapstructure:"postgres"`
	MongoDB     MongoConfig    `mapstructure:"mongodb"`
	S3          S3Config       `mapstructure:"s3"`
	Kafka       KafkaConfig    `mapstructure:"kafka"`
}

type ServerConfig struct {
	Address  string `mapstructure:"address"`
	LogLevel string `mapstructure:"log_level"`
}

type UserService struct {
	Address string `mapstructure:"user_service_address"`
}

type PostgresConfig struct {
	Host            string        `mapstructure:"host"`
	Port            int           `mapstructure:"port"`
	User            string        `mapstructure:"user"`
	Password        string        `mapstructure:"password"`
	Database        string        `mapstructure:"database"`
	SSLMode         string        `mapstructure:"sslmode"`
	MaxOpenConns    int           `mapstructure:"max_open_conns"`
	MaxIdleConns    int           `mapstructure:"max_idle_conns"`
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime"`
	MigrationsPath  string        `mapstructure:"migrations_path"`
	URL             string        `mapstructure:"url"`
}

func (p PostgresConfig) DSN() string {
	if p.URL != "" {
		return p.URL
	}
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		p.Host, p.Port, p.User, p.Password, p.Database, p.SSLMode,
	)
}

type MongoConfig struct {
	URI                  string        `mapstructure:"uri"`
	Database             string        `mapstructure:"database"`
	ConnectTimeout       time.Duration `mapstructure:"connect_timeout"`
	MaxPoolSize          uint64        `mapstructure:"max_pool_size"`
	MinPoolSize          uint64        `mapstructure:"min_pool_size"`
	TagsCollection       string        `mapstructure:"tags_collection"`
	CategoriesCollection string        `mapstructure:"categories_collection"`
	MetadataCollection   string        `mapstructure:"metadata_collection"`
}

type S3Config struct {
	Region          string   `mapstructure:"region"`
	Bucket          string   `mapstructure:"bucket"`
	AccessKeyID     string   `mapstructure:"access_key_id"`
	SecretAccessKey string   `mapstructure:"secret_access_key"`
	Endpoint        string   `mapstructure:"endpoint"`
	UseSSL          bool     `mapstructure:"use_ssl"`
	MaxFileSize     int64    `mapstructure:"max_file_size"`
	AllowedTypes    []string `mapstructure:"allowed_types"`
	PublicURL       string   `mapstructure:"public_url"`
}

type KafkaConfig struct {
	Brokers             []string      `mapstructure:"brokers"`
	TopicListingsEvents string        `mapstructure:"topic_listings_events"`
	GroupID             string        `mapstructure:"group_id"`
	RequiredAcks        int           `mapstructure:"required_acks"`
	Compression         string        `mapstructure:"compression"`
	WriteTimeout        time.Duration `mapstructure:"write_timeout"`
	RetryMax            int           `mapstructure:"retry_max"`
}

func Load() (*Config, error) {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	setDefaults()

	if err := viper.ReadInConfig(); err != nil {
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if errors.As(err, &configFileNotFoundError) {
			log.Println("⚠️  .env file not found, using environment variables and defaults")
		} else {
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
	}

	cfg := &Config{
		ServerConfig{
			Address:  viper.GetString("SERVER_ADDRESS"),
			LogLevel: viper.GetString("LOG_LEVEL"),
		},
		UserService{
			Address: viper.GetString("USER_SERVICE_ADDRESS"),
		},
		PostgresConfig{
			Host:            viper.GetString("SERVER_ADDRESS"),
			Port:            viper.GetInt("POSTGRES_PORT"),
			User:            viper.GetString("POSTGRES_USER"),
			Password:        viper.GetString("POSTGRES_PASSWORD"),
			Database:        viper.GetString("POSTGRES_DB"),
			SSLMode:         viper.GetString("POSTGRES_SSLMODE"),
			MaxOpenConns:    viper.GetInt("POSTGRES_MAX_OPEN_CONNS"),
			MaxIdleConns:    viper.GetInt("POSTGRES_MAX_IDLE_CONNS"),
			ConnMaxLifetime: viper.GetDuration("POSTGRES_CONN_MAX_LIFETIME"),
			MigrationsPath:  viper.GetString("POSTGRES_MIGRATIONS_PATH"),
			URL:             viper.GetString("DATABASE_URL"),
		},
		MongoConfig{
			URI:                  viper.GetString("MONGODB_URI"),
			Database:             viper.GetString("MONGODB_DATABASE"),
			ConnectTimeout:       viper.GetDuration("MONGODB_CONNECT_TIMEOUT"),
			MaxPoolSize:          viper.GetUint64("MONGODB_MAX_POOL_SIZE"),
			MinPoolSize:          viper.GetUint64("MONGODB_MIN_POOL_SIZE"),
			TagsCollection:       viper.GetString("MONGODB_TAGS_COLLECTION"),
			CategoriesCollection: viper.GetString("MONGODB_CATEGORIES_COLLECTION"),
			MetadataCollection:   viper.GetString("MONGODB_METADATA_COLLECTION"),
		},
		S3Config{
			Region:          viper.GetString("S3_REGION"),
			Bucket:          viper.GetString("S3_BUCKET"),
			AccessKeyID:     viper.GetString("S3_ACCESS_KEY_ID"),
			SecretAccessKey: viper.GetString("S3_SECRET_ACCESS_KEY"),
			Endpoint:        viper.GetString("S3_ENDPOINT"),
			UseSSL:          viper.GetBool("S3_USE_SSL"),
			MaxFileSize:     viper.GetInt64("S3_MAX_FILE_SIZE"),
			AllowedTypes:    viper.GetStringSlice("S3_ALLOWED_TYPES"),
			PublicURL:       viper.GetString("S3_PUBLIC_URL"),
		},
		KafkaConfig{
			Brokers:             []string{viper.GetString("KAFKA_BROKERS")},
			TopicListingsEvents: viper.GetString("KAFKA_TOPIC_LISTINGS_EVENTS"),
			GroupID:             viper.GetString("KAFKA_GROUP_ID"),
			RequiredAcks:        viper.GetInt("KAFKA_REQUIRED_ACKS"),
			Compression:         viper.GetString("KAFKA_COMPRESSION"),
			WriteTimeout:        viper.GetDuration("KAFKA_WRITE_TIMEOUT"),
			RetryMax:            viper.GetInt("KAFKA_RETRY_MAX"),
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

	viper.SetDefault("USER_SERVICE_ADDRESS", "user-service:50051")

	viper.SetDefault("SERVER_ADDRESS", "localhost")
	viper.SetDefault("POSTGRES_PORT", 5432)
	viper.SetDefault("POSTGRES_USER", "avigo_listings")
	viper.SetDefault("POSTGRES_DB", "avigo_listings")
	viper.SetDefault("POSTGRES_SSLMODE", "disable")
	viper.SetDefault("POSTGRES_MAX_OPEN_CONNS", 25)
	viper.SetDefault("POSTGRES_MAX_IDLE_CONNS", 5)
	viper.SetDefault("POSTGRES_CONN_MAX_LIFETIME", 5*time.Minute)
	viper.SetDefault("POSTGRES_MIGRATIONS_PATH", "migrations")

	viper.SetDefault("MONGODB_URI", "mongodb://localhost:27017")
	viper.SetDefault("MONGODB_DATABASE", "listings_metadata")
	viper.SetDefault("MONGODB_CONNECT_TIMEOUT", 10*time.Second)
	viper.SetDefault("MONGODB_MAX_POOL_SIZE", 100)
	viper.SetDefault("MONGODB_MIN_POOL_SIZE", 10)
	viper.SetDefault("MONGODB_TAGS_COLLECTION", "tags")
	viper.SetDefault("MONGODB_CATEGORIES_COLLECTION", "categories")
	viper.SetDefault("MONGODB_METADATA_COLLECTION", "listing_metadata")

	viper.SetDefault("S3_REGION", "us-east-1")
	viper.SetDefault("S3_BUCKET", "avigo-listings-media")
	viper.SetDefault("S3_USE_SSL", true)
	viper.SetDefault("S3_MAX_FILE_SIZE", 10*1024*1024) // 10MB
	viper.SetDefault("S3_ALLOWED_TYPES", "image/jpeg,image/png,image/webp,video/mp4")

	viper.SetDefault("KAFKA_BROKERS", "localhost:9092")
	viper.SetDefault("KAFKA_TOPIC_LISTINGS_EVENTS", "listing.events")
	viper.SetDefault("KAFKA_GROUP_ID", "listing-service")
	viper.SetDefault("KAFKA_REQUIRED_ACKS", -1)
	viper.SetDefault("KAFKA_COMPRESSION", "snappy")
	viper.SetDefault("KAFKA_WRITE_TIMEOUT", 10*time.Second)
	viper.SetDefault("KAFKA_RETRY_MAX", 3)
}

func validateConfig(cfg *Config) error {
	if cfg.Postgres.Password == "" && cfg.Postgres.URL == "" {
		return errors.New("POSTGRES_PASSWORD or DATABASE_URL is required")
	}

	if cfg.S3.AccessKeyID == "" {
		return errors.New("S3_ACCESS_KEY_ID is required")
	}

	if cfg.S3.SecretAccessKey == "" {
		return errors.New("S3_SECRET_ACCESS_KEY is required")
	}

	if len(cfg.Kafka.Brokers) == 0 {
		return errors.New("KAFKA_BROKERS is required")
	}

	if cfg.Kafka.TopicListingsEvents == "" {
		return errors.New("KAFKA_TOPIC_LISTINGS_EVENTS is required")
	}

	return nil
}
