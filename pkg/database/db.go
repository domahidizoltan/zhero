//go:build !android
// +build !android

package database

import (
	_ "modernc.org/sqlite"
)

const sqliteDriver = "sqlite"
