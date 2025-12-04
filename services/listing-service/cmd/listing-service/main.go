package main

import (
	"context"
	"github.com/squ1ky/avigo-c2c-marketplace/services/listing-service/internal/repository/s3"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/squ1ky/avigo-c2c-marketplace/services/listing-service/internal/config"
	mongodb "github.com/squ1ky/avigo-c2c-marketplace/services/listing-service/internal/db/mongo"
	"github.com/squ1ky/avigo-c2c-marketplace/services/listing-service/internal/db/postgres"
	"github.com/squ1ky/avigo-c2c-marketplace/services/listing-service/internal/kafka"
)

func main() {
	log.Println("Starting Listing Service...")

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	log.Println("Configuration loaded successfully")
	log.Printf("Server will run on %s", cfg.Server.Address)

	ctx := context.Background()

	log.Println("Connecting to PostgreSQL...")
	postgresDB, err := postgres.NewPostgresDB(&cfg.Postgres)
	if err != nil {
		log.Fatalf("Failed to connect to PostgreSQL: %v", err)
	}
	log.Println("PostgreSQL connected successfully")

	if err := postgres.RunMigrations(&cfg.Postgres); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}
	log.Println("Database migrations completed")

	log.Println("Connection to MongoDB...")
	mongoClient, err := mongodb.NewMongoDB(&cfg.MongoDB)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	mongoDB := mongodb.GetDatabase(mongoClient, cfg.MongoDB.Database)
	log.Println("MongoDB connected successfully")
	_ = mongoDB

	log.Println("Initializing MinIO repository...")
	minioRepo, err := s3.NewMediaStorage(cfg.S3)
	if err != nil {
		log.Fatalf("Failed to initialize MinIO: %v", err)
	}
	_ = minioRepo
	log.Printf("MinIO repository initialized (bucket: %s)", cfg.S3.Bucket)

	log.Println("Initializing Kafka Producer...")
	producer, err := kafka.NewProducer(cfg.Kafka)
	if err != nil {
		log.Fatalf("Failed to initialize Kafka producer: %v", err)
	}
	log.Printf("Kafka producer initialized (topic: %s)", cfg.Kafka.TopicListingsEvents)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	<-quit
	log.Println("\nShutdown signal received...")
	log.Println("Initiating graceful shutdown...")

	shutdownCtx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	log.Println("Flushing Kafka messages...")
	if err := producer.Close(); err != nil {
		log.Printf("Error closing Kafka producer: %v", err)
	} else {
		log.Println("Kafka producer closed")
	}

	log.Println("Disconnecting from MongoDB...")
	if err := mongoClient.Disconnect(shutdownCtx); err != nil {
		log.Printf("Error disconnecting from MongoDB: %v", err)
	} else {
		log.Println("MongoDB disconnected")
	}

	log.Println("Closing PostgreSQL connection...")
	if err := postgresDB.Close(); err != nil {
		log.Printf("Error closing PostgreSQL connection: %v", err)
	} else {
		log.Println("PostgreSQL connection closed")
	}

	log.Println("\n" + strings.Repeat("=", 60))
	log.Println("Listing Service stopped gracefully")
	log.Println(strings.Repeat("=", 60))
}
