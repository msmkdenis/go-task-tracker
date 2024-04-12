package storage

import (
	"database/sql"
	"log/slog"
	"os"
)

type SQLiteDB struct {
	DB *sql.DB
}

func NewSQLiteDB(dbFile string) *SQLiteDB {
	db, err := sql.Open("sqlite", dbFile)
	if err != nil {
		slog.Error("failed to open database", "error", err)
		os.Exit(1)
	}

	return &SQLiteDB{
		DB: db,
	}
}
