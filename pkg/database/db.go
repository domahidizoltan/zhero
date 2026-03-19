// Package database is for managing databases
package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

type txKey = struct{}

var (
	db                     *sql.DB
	ErrTransactionNotFound = errors.New("transaction not found")
	ErrTransaction         = errors.New("transaction failed")
	ErrDBConnection        = errors.New("database connection failed")
	ErrDBMigration         = errors.New("database migration failed")
)

func InitSqliteDB(dbFile string) error {
	d, err := sql.Open(sqliteDriver, dbFile)
	if err != nil {
		return fmt.Errorf("SQLite %w: %w", ErrDBConnection, err)
	}
	db = d
	return nil
}

func GetDB() *sql.DB {
	return db
}

func GetTx(ctx context.Context) *sql.Tx {
	if tx, found := ctx.Value(txKey{}).(*sql.Tx); found {
		return tx
	}
	return nil
}

func InTx(ctx context.Context, fn func(ctx context.Context) error) error {
	db := GetDB()
	if tx := GetTx(ctx); tx != nil {
		return fn(ctx)
	}

	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("%w: %w", ErrTransaction, err)
	}

	ctx = context.WithValue(ctx, txKey{}, tx)
	if err := fn(ctx); err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("%w: %w", ErrTransaction, err)
	}

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
	return nil
}

func Migrate(db *sql.DB, scripts []string) error {
	for _, s := range scripts {
		if _, err := db.Exec(s); err != nil {
			return fmt.Errorf("%w: %w", ErrDBMigration, err)
		}
	}
	return nil
}
