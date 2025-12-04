package pgrepo

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type key string

const (
	txKey key = "tx"
)

type TransactionManager interface {
	WithinTransaction(ctx context.Context, tFunc func(ctx context.Context) error) error
}

type txManager struct {
	db *sqlx.DB
}

func NewTransactionManager(db *sqlx.DB) TransactionManager {
	return &txManager{db: db}
}

func (m *txManager) WithinTransaction(ctx context.Context, tFunc func(ctx context.Context) error) error {
	tx, err := m.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("cannot begin transaction: %w", err)
	}

	ctxWithTx := context.WithValue(ctx, txKey, tx)

	if err := tFunc(ctxWithTx); err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("cannot commit transaction: %w", err)
	}
	return nil
}

func injectTx(ctx context.Context) *sqlx.Tx {
	if tx, ok := ctx.Value(txKey).(*sqlx.Tx); ok {
		return tx
	}
	return nil
}

type SQLQueryer interface {
	sqlx.ExtContext
	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	NamedExecContext(ctx context.Context, query string, arg interface{}) (sql.Result, error)
}
