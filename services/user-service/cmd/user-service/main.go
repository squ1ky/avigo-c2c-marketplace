package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/squ1ky/avigo-c2c-marketplace/services/user-service/internal/config"
	"github.com/squ1ky/avigo-c2c-marketplace/services/user-service/internal/db"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Config load error: %v", err)
	}

	pool, err := db.NewPgPool(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}
	defer pool.Close()

	router := gin.Default()

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	if err := router.Run(cfg.ServerAddress); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
