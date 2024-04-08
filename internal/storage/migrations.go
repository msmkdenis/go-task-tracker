package storage

import (
	"database/sql"
	"embed"
	"log/slog"
	"os"

	"github.com/pressly/goose/v3"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS //

type Migrations struct {
	db *sql.DB
}

func NewMigrations(db *SQLiteDB) *Migrations {
	err := goose.SetDialect("sqlite3")
	if err != nil {
		slog.Error("failed to set dialect", "error", err)
		os.Exit(1)
	}
	goose.SetBaseFS(embedMigrations)

	return &Migrations{
		db: db.DB,
	}
}

func (m *Migrations) Up() error {
	return goose.Up(m.db, "migrations")
}
