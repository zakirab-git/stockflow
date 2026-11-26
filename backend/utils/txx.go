package utils

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

type TransactionxManagerInterface interface {
	WithTxM(ctx context.Context, fn TxxFn) error
}

type TransactionxManager struct {
	db *sqlx.DB
}

type Getter interface {
	sqlx.ExtContext
	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	NamedExecContext(ctx context.Context, query string, arg interface{}) (sql.Result, error)
}

type KeyTx struct{}

func NewTxM(db *sqlx.DB) *TransactionxManager {
	return &TransactionxManager{
		db: db,
	}
}

type TxxFn func(context.Context) error

func (txM *TransactionxManager) WithTxM(ctx context.Context, fn TxxFn) (err error) {
	tx, err := txM.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	ctxWithTx := context.WithValue(ctx, KeyTx{}, tx)

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		} else if err != nil {
			_ = tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	err = fn(ctxWithTx)
	return err
}

func ExtractTx(ctx context.Context) *sqlx.Tx {
	if tx, ok := ctx.Value(KeyTx{}).(*sqlx.Tx); ok {
		return tx
	}

	return nil
}
