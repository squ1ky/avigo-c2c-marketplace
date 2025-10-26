package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/squ1ky/avigo-c2c-marketplace/services/user-service/internal/auth"
	grpcserver "github.com/squ1ky/avigo-c2c-marketplace/services/user-service/internal/grpc"
	"github.com/squ1ky/avigo-c2c-marketplace/services/user-service/internal/handler"
	"github.com/squ1ky/avigo-c2c-marketplace/services/user-service/internal/kafka"
	"github.com/squ1ky/avigo-c2c-marketplace/services/user-service/internal/middleware"
	"github.com/squ1ky/avigo-c2c-marketplace/services/user-service/internal/repository"
	"github.com/squ1ky/avigo-c2c-marketplace/services/user-service/internal/service"
	"github.com/squ1ky/avigo-c2c-marketplace/services/user-service/internal/validation"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/squ1ky/avigo-c2c-marketplace/services/user-service/internal/config"
	"github.com/squ1ky/avigo-c2c-marketplace/services/user-service/internal/db"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Config load error: %v", err)
	}

	if err := db.RunMigrations(cfg.Database.URL, cfg.Database.MigrationsPath); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	database, err := db.NewGormDB(cfg.Database.URL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	kafkaProducer, err := kafka.NewProducer(cfg.Kafka, logger)
	if err != nil {
		log.Fatalf("Failed to create Kafka producer: %v", err)
	}
	defer kafkaProducer.Close()

	validator := validation.NewValidator()
	jwtManager := auth.NewJWTManager(&cfg.JWT)
	txManager := repository.NewTransactionManager(database)
	userRepo := repository.NewUserRepository(database)

	authService := service.NewAuthService(
		userRepo,
		txManager,
		validator,
		jwtManager,
		kafkaProducer,
		cfg.Auth.ConfirmationCodeExpiry,
	)

	authHandler := handler.NewAuthHandler(authService, cfg.Cookie, cfg.JWT)

	// HTTP (Gin)
	engine := gin.Default()
	engine.Use(middleware.ErrorHandler())
	router := handler.NewRouter(authHandler)
	router.SetupRoutes(engine)

	// gRPC
	grpcSrv, err := grpcserver.NewServer(cfg.GRPC, userRepo)
	if err != nil {
		log.Fatalf("Failed to create GRPC server: %v", err)
	}

	errChan := make(chan error, 2)

	go func() {
		log.Printf("HTTP server listening on %s", cfg.Server.Address)
		if err := engine.Run(cfg.Server.Address); err != nil {
			errChan <- fmt.Errorf("HTTP server error: %w", err)
		}
	}()

	go func() {
		if err := grpcSrv.Start(); err != nil {
			errChan <- fmt.Errorf("gRPC server error: %w", err)
		}
	}()

	// Graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-errChan:
		log.Fatalf("Server error: %v", err)
	case sig := <-quit:
		log.Printf("Received signal %v, shutting down...", sig)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	shutdownDone := make(chan struct{})
	go func() {
		grpcSrv.GracefulStop()
		close(shutdownDone)
	}()

	select {
	case <-shutdownDone:
		log.Println("Graceful shutdown completed")
	case <-ctx.Done():
		log.Println("Shutdown timeout, forcing stop")
		grpcSrv.Stop()
	}
}
