package repository

import (
	"database/sql"
	"errors"

	"github.com/msmkdenis/go-task-tracker/internal/model"
)

func (t *SQLiteTaskRepository) SelectByID(id int64) (model.Task, error) {
	row := t.db.QueryRow(
		`
		select
			id, 
		    date, 
		    title, 
		    comment, 
		    repeat
		from scheduler
		where id = :id
		`,
		sql.Named("id", id))

	if errors.Is(row.Err(), sql.ErrNoRows) {
		return model.Task{}, model.ErrTaskNotFound
	}

	if row.Err() != nil {
		return model.Task{}, row.Err()
	}

	var task model.Task
	if errScan := row.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat); errScan != nil {
		return model.Task{}, errScan
	}

	return task, nil
}
