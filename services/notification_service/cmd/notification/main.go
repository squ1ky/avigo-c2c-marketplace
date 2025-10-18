package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/squ1ky/avigo-c2c-marketplace/services/notification-service/internal/adapter/email"
	"github.com/squ1ky/avigo-c2c-marketplace/services/notification-service/internal/adapter/kafka"
	"github.com/squ1ky/avigo-c2c-marketplace/services/notification-service/internal/config"
	database "github.com/squ1ky/avigo-c2c-marketplace/services/notification-service/internal/db"
	httpDelivery "github.com/squ1ky/avigo-c2c-marketplace/services/notification-service/internal/handler/http"
	"github.com/squ1ky/avigo-c2c-marketplace/services/notification-service/internal/repository/postgres"
	"github.com/squ1ky/avigo-c2c-marketplace/services/notification-service/internal/usecase"
)

func main() {
	log.Println("🚀 Starting Notification Service...")

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Config load error: %v", err)
	}

	db, err := database.NewPostgresDB(&cfg.Database)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	log.Println("Connected to PostgreSQL")

	if err := database.RunMigrations(&cfg.Database); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}
	log.Println("Database migrations completed")

	emailSender := email.NewSMTPSender(&cfg.Email)
	templateRenderer, err := email.NewHTMLTemplateRenderer(&cfg.Email)
	if err != nil {
		log.Fatalf("Failed to initialize template renderer: %v", err)
	}
	log.Println("Email infrastructure initialized")

	notificationRepo := postgres.NewNotificationRepository(db)
	log.Println("Repository initialized")

	emailUseCase := usecase.NewEmailNotificationUseCase(
		notificationRepo,
		emailSender,
		templateRenderer,
	)
	log.Println("Use case initialized")

	kafkaHandler := kafka.NewNotificationHandler(emailUseCase)
	kafkaConsumer := kafka.NewConsumer(&cfg.Kafka, kafkaHandler)
	log.Println("Kafka consumer initialized")

	router := httpDelivery.SetupRouter()
	srv := &http.Server{
		Addr:         cfg.Server.Address,
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	log.Println("HTTP server initialized")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		log.Println("Starting Kafka consumer")
		if err := kafkaConsumer.Start(ctx); err != nil {
			log.Printf("Kafka consumer stopped: %v", err)
		}
	}()

	go func() {
		log.Printf("Starting HTTP server on %s", cfg.Server.Address)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	<-quit
	log.Println("\n⏳ Shutting down gracefully...")
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	log.Println("⏸️  Stopping HTTP server...")
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("⚠️  HTTP server shutdown error: %v", err)
	}

	defer func() {
		sqlDB, _ := db.DB()
		if sqlDB != nil {
			sqlDB.Close()
		}
	}()
}
