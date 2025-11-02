package config

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Postgres PostgresConfig `mapstructure:"postgres"`
	MongoDB  MongoConfig    `mapstructure:"mongodb"`
	S3       S3Config       `mapstructure:"s3"`
	Kafka    KafkaConfig    `mapstructure:"kafka"`
}

type ServerConfig struct {
	Address      string        `mapstructure:"address"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
	LogLevel     string        `mapstructure:"log_level"`
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
	v := viper.New()

	v.SetConfigFile(".env")
	v.SetConfigType("env")

	v.AutomaticEnv()

	setDefaults(v)

	if err := v.ReadInConfig(); err != nil {
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if errors.As(err, &configFileNotFoundError) {
			log.Println("⚠️  .env file not found, using environment variables and defaults")
		} else {
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
	}

	cfg := &Config{}

	if err := v.Unmarshal(cfg); err != nil {
		return nil, fmt.Errorf("unable to unmarshal config: %w", err)
	}

	if brokersStr := v.GetString("KAFKA_BROKERS"); brokersStr != "" {
		cfg.Kafka.Brokers = parseCommaSeparated(brokersStr)
	}

	if allowedTypesStr := v.GetString("S3_ALLOWED_TYPES"); allowedTypesStr != "" {
		cfg.S3.AllowedTypes = parseCommaSeparated(allowedTypesStr)
	}

	if err := validateConfig(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

func setDefaults(v *viper.Viper) {

	v.SetDefault("SERVER_ADDRESS", ":8083")
	v.SetDefault("SERVER_READ_TIMEOUT", 15*time.Second)
	v.SetDefault("SERVER_WRITE_TIMEOUT", 15*time.Second)
	v.SetDefault("LOG_LEVEL", "info")

	v.SetDefault("POSTGRES_HOST", "localhost")
	v.SetDefault("POSTGRES_PORT", 5432)
	v.SetDefault("POSTGRES_USER", "avigo_listings")
	v.SetDefault("POSTGRES_DB", "avigo_listings")
	v.SetDefault("POSTGRES_SSLMODE", "disable")
	v.SetDefault("POSTGRES_MAX_OPEN_CONNS", 25)
	v.SetDefault("POSTGRES_MAX_IDLE_CONNS", 5)
	v.SetDefault("POSTGRES_CONN_MAX_LIFETIME", 5*time.Minute)
	v.SetDefault("POSTGRES_MIGRATIONS_PATH", "migrations")

	v.SetDefault("MONGODB_URI", "mongodb://localhost:27017")
	v.SetDefault("MONGODB_DATABASE", "listings_metadata")
	v.SetDefault("MONGODB_CONNECT_TIMEOUT", 10*time.Second)
	v.SetDefault("MONGODB_MAX_POOL_SIZE", 100)
	v.SetDefault("MONGODB_MIN_POOL_SIZE", 10)
	v.SetDefault("MONGODB_TAGS_COLLECTION", "tags")
	v.SetDefault("MONGODB_CATEGORIES_COLLECTION", "categories")
	v.SetDefault("MONGODB_METADATA_COLLECTION", "listing_metadata")

	v.SetDefault("S3_REGION", "us-east-1")
	v.SetDefault("S3_BUCKET", "avigo-listings-media")
	v.SetDefault("S3_USE_SSL", true)
	v.SetDefault("S3_MAX_FILE_SIZE", 10*1024*1024) // 10MB
	v.SetDefault("S3_ALLOWED_TYPES", "image/jpeg,image/png,image/webp,video/mp4")

	v.SetDefault("KAFKA_BROKERS", "localhost:9092")
	v.SetDefault("KAFKA_TOPIC_LISTINGS_EVENTS", "listing.events")
	v.SetDefault("KAFKA_GROUP_ID", "listing-service")
	v.SetDefault("KAFKA_REQUIRED_ACKS", -1)
	v.SetDefault("KAFKA_COMPRESSION", "snappy")
	v.SetDefault("KAFKA_WRITE_TIMEOUT", 10*time.Second)
	v.SetDefault("KAFKA_RETRY_MAX", 3)
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

func parseCommaSeparated(input string) []string {
	parts := strings.Split(input, ",")
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}
