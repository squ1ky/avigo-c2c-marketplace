package db

import (
	"errors"
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/squ1ky/avigo-c2c-marketplace/services/notification-service/internal/config"
)

func RunMigrations(cfg *config.DatabaseConfig) error {
	log.Println("Running database migrations...")

	migrationsPath := fmt.Sprintf("file://%s", cfg.MigrationsPath)

	m, err := migrate.New(
		migrationsPath,
		cfg.URL,
	)
	if err != nil {
		return fmt.Errorf("failed to create migration instance: %w", err)
	}
	defer m.Close()

	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			log.Println("No new migrations to apply")
			return nil
		}
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	log.Println("Database migrations completed")
	return nil
}

func RollbackMigrations(cfg *config.DatabaseConfig) error {
	log.Println("Rolling back last migraton...")

	migrationPath := fmt.Sprintf("file://%s", cfg.MigrationsPath)
	m, err := migrate.New(
		migrationPath,
		cfg.URL,
	)
	if err != nil {
		return fmt.Errorf("failed to create migration instance: %w", err)
	}
	defer m.Close()

	if err := m.Steps(-1); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			log.Println("No migrations to rollback")
			return nil
		}
		return fmt.Errorf("failed to rollback migration: %w", err)
	}

	log.Println("Migration rollback completed successfully")
	return nil
}
