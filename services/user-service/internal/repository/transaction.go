package repository

import (
	"context"
	"gorm.io/gorm"
)

type contextKey string

const txKey contextKey = "db_tx"

type TransactionManager interface {
	WithinTransaction(ctx context.Context, fn func(ctx context.Context) error) error
}

type transactionManager struct {
	db *gorm.DB
}

func NewTransactionManager(db *gorm.DB) TransactionManager {
	return &transactionManager{db: db}
}

func (tm *transactionManager) WithinTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	return tm.db.Transaction(func(tx *gorm.DB) error {
		txCtx := context.WithValue(ctx, txKey, tx)
		return fn(txCtx)
	})
}

func GetTxFromContext(ctx context.Context) *gorm.DB {
	if tx, ok := ctx.Value(txKey).(*gorm.DB); ok {
		return tx
	}
	return nil
}
