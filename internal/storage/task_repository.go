package storage

import (
	"database/sql"

	"github.com/msmkdenis/go-task-tracker/internal/model"
)

type SQLiteTaskRepository struct {
	db *sql.DB
}

func NewSQLiteTaskRepository(db *SQLiteDB) *SQLiteTaskRepository {
	return &SQLiteTaskRepository{
		db: db.db,
	}
}

func (t *SQLiteTaskRepository) Insert(task model.Task) (int64, error) {
	res, err := t.db.Exec(
		`
		insert into scheduler 
		    (date, title, comment, repeat) 
		values 
		    (:date, :title, :comment, :repeat)
		`,
		sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat))
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}
