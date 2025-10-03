-include .env
export

MIGRATE=migrate
MIGRATIONS_DIR=$(MIGRATIONS_PATH)

ifndef MIGRATIONS_DIR
	MIGRATIONS_DIR=internal/db/migrations
endif

.PHONY: help migrate-create migrate-up migrate-down migrate-version migrate-status

.DEFAULT_GOAL := help

help:
	@echo === Migration Commands ===
	@echo.
	@echo   make migrate-create NAME=^<name^>  - Create new migration
	@echo   make migrate-up                    - Apply all migrations
	@echo   make migrate-down                  - Rollback last migration
	@echo   make migrate-version               - Show current version
	@echo   make migrate-status                - Show migration status
	@echo.

migrate-create:
ifndef NAME
	$(error NAME not specified. Use: make migrate-create NAME=migration_name)
endif
	$(MIGRATE) create -ext sql -dir $(MIGRATIONS_DIR) -seq $(NAME)
	@echo Migration created in $(MIGRATIONS_DIR)

migrate-up:
	@echo Applying migrations...
	$(MIGRATE) -path $(MIGRATIONS_DIR) -database "$(DATABASE_URL)" -verbose up
	@echo Migrations applied

migrate-down:
	@echo Rolling back last migration...
	$(MIGRATE) -path $(MIGRATIONS_DIR) -database "$(DATABASE_URL)" -verbose down 1
	@echo Migration rolled back

migrate-version:
	@echo Current version:
	@$(MIGRATE) -path $(MIGRATIONS_DIR) -database "$(DATABASE_URL)" version

migrate-status:
	@echo Migration version:
	@$(MIGRATE) -path $(MIGRATIONS_DIR) -database "$(DATABASE_URL)" version
	@echo.
	@echo Migration files:
	@dir $(MIGRATIONS_DIR) 2>nul || ls -la $(MIGRATIONS_DIR) 2>/dev/null
