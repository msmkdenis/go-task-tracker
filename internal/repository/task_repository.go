package repository

import (
	"database/sql"

	"github.com/msmkdenis/go-task-tracker/internal/storage"
)

const (
	taskLoadLimit = 50
)

type SQLiteTaskRepository struct {
	db *sql.DB
}

func NewSQLiteTaskRepository(db *storage.SQLiteDB) *SQLiteTaskRepository {
	return &SQLiteTaskRepository{
		db: db.DB,
	}
}
