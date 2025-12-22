// Package db is for managing databases
package database

import (
	"context"
	"database/sql"
	"errors"
)

type txKey = struct{}

var (
	db                     *sql.DB
	ErrTransactionNotFound = errors.New("transaction not found")
)

func InitSqliteDB(dbFile string) error {
	d, err := sql.Open(sqliteDriver, dbFile)
	if err != nil {
		return err
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
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	ctx = context.WithValue(ctx, txKey{}, tx)
	if err := fn(ctx); err != nil {
		_ = tx.Rollback()
		return err
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
			return err
		}
	}
	return nil
}
