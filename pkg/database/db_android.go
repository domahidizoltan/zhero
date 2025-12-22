//go:build android
// +build android

package database

import (
	_ "github.com/mattn/go-sqlite3"
)

const sqliteDriver = "sqlite3"
