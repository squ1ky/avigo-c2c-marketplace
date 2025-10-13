package main

import (
	"github.com/gin-gonic/gin"
	"github.com/squ1ky/avigo-c2c-marketplace/services/user-service/internal/auth"
	"github.com/squ1ky/avigo-c2c-marketplace/services/user-service/internal/handler"
	"github.com/squ1ky/avigo-c2c-marketplace/services/user-service/internal/kafka"
	"github.com/squ1ky/avigo-c2c-marketplace/services/user-service/internal/middleware"
	"github.com/squ1ky/avigo-c2c-marketplace/services/user-service/internal/repository"
	"github.com/squ1ky/avigo-c2c-marketplace/services/user-service/internal/service"
	"github.com/squ1ky/avigo-c2c-marketplace/services/user-service/internal/validation"
	"log"
	"log/slog"
	"os"

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

	engine := gin.Default()

	engine.Use(middleware.ErrorHandler())

	router := handler.NewRouter(authHandler)
	router.SetupRoutes(engine)

	if err := engine.Run(cfg.Server.Address); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
