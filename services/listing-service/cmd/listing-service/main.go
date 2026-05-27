package main

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/squ1ky/avigo-c2c-marketplace/services/listing-service/internal/grpc/client/user"
	"github.com/squ1ky/avigo-c2c-marketplace/services/listing-service/internal/handler"
	"github.com/squ1ky/avigo-c2c-marketplace/services/listing-service/internal/llm"
	"github.com/squ1ky/avigo-c2c-marketplace/services/listing-service/internal/middleware"
	mngrepo "github.com/squ1ky/avigo-c2c-marketplace/services/listing-service/internal/repository/mongo"
	pgrepo "github.com/squ1ky/avigo-c2c-marketplace/services/listing-service/internal/repository/postgres"
	"github.com/squ1ky/avigo-c2c-marketplace/services/listing-service/internal/repository/s3"
	"github.com/squ1ky/avigo-c2c-marketplace/services/listing-service/internal/service"
	"log"
	"net/http"
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

	log.Println("Initializing MinIO repository...")
	minioRepo, err := s3.NewMediaStorage(cfg.S3)
	if err != nil {
		log.Fatalf("Failed to initialize MinIO: %v", err)
	}
	log.Printf("MinIO repository initialized (bucket: %s)", cfg.S3.Bucket)

	log.Println("Initializing Kafka Producer...")
	producer, err := kafka.NewProducer(cfg.Kafka)
	if err != nil {
		log.Fatalf("Failed to initialize Kafka producer: %v", err)
	}
	log.Printf("Kafka producer initialized (topic: %s)", cfg.Kafka.TopicListingsEvents)

	// gRPC (to user-service)
	log.Printf("Connecting to User Service gRPC at %s...", cfg.UserService.Address)
	userClient, err := user.NewClient(cfg.UserService.Address)
	if err != nil {
		log.Fatalf("Failed to initialize User Service client: %v", err)
	}
	log.Println("User Service gRPC client initialized")

	txManager := pgrepo.NewTransactionManager(postgresDB)

	listingRepo := pgrepo.NewListingRepository(postgresDB)
	mediaRepo := pgrepo.NewMediaRepository(postgresDB)
	orderRepo := pgrepo.NewOrderRepository(postgresDB)
	reviewRepo := pgrepo.NewReviewRepository(postgresDB)
	charsRepo := mngrepo.NewCharacteristicsRepository(mongoDB)

	mediaSvc := service.NewMediaService(minioRepo, mediaRepo, cfg.S3)
	geminiClient := llm.NewGeminiClient(cfg.Gemini)
	tagSuggestionSvc := service.NewTagSuggestionService(geminiClient, listingRepo)
	listingSvc := service.NewListingService(
		listingRepo,
		mediaRepo,
		charsRepo,
		minioRepo,
		txManager,
		cfg.S3,
		userClient,
	)
	orderSvc := service.NewOrderService(orderRepo, listingRepo, txManager, userClient)
	reviewSvc := service.NewReviewService(reviewRepo, orderRepo)

	h := handler.NewHandler(listingSvc, tagSuggestionSvc, orderSvc, reviewSvc, mediaSvc)

	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(gin.Logger())
	router.Use(middleware.ErrorHandler())

	api := router.Group("/api")
	h.Init(api)

	srv := &http.Server{
		Addr:    cfg.Server.Address,
		Handler: router,
	}

	gcStop := make(chan struct{})
	go func() {
		log.Println("Starting GC worker...")
		ticker := time.NewTicker(1 * time.Hour)
		defer ticker.Stop()
		for {
			select {
			case <-gcStop:
				return
			case <-ticker.C:
				deleted, err := mediaRepo.DeleteExpiredTemp(context.Background(), 24*time.Hour)
				if err != nil {
					log.Printf("GC Error: %v", err)
					continue
				}
				if len(deleted) > 0 {
					keys := make([]string, len(deleted))
					for i, m := range deleted {
						keys[i] = m.S3Key
					}

					if err := minioRepo.DeleteFiles(context.Background(), keys); err != nil {
						log.Printf("GC S3 Error: %v", err)
					} else {
						log.Printf("GC: Cleaned %d expired files", len(deleted))
					}
				}
			}
		}
	}()

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	<-quit
	log.Println("\nShutdown signal received...")

	shutdownCtx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	log.Println("Stopping HTTP server...")
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("Server Force Shutdown: %v", err)
	}

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
