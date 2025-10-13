package main

import (
	"context"
	"errors"
	"github.com/squ1ky/avigo-c2c-marketplace/services/gateway-service/internal/config"
	"github.com/squ1ky/avigo-c2c-marketplace/services/gateway-service/internal/router"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	r := router.SetupRouter(cfg)

	srv := &http.Server{
		Addr:    cfg.Server.Address,
		Handler: r,
	}

	// Graceful shutdown
	go func() {
		log.Printf("Gateway listening on %s", cfg.Server.Address)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down gateway...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Gateway forced to shutdown:", err)
	}

	log.Println("Gateway stopped")
}
