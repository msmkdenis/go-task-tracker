package repository

import (
	"database/sql"

	"github.com/msmkdenis/go-task-tracker/internal/model"
)

func (t *SQLiteTaskRepository) UpdateByID(task model.Task) error {
	res, err := t.db.Exec(
		`
		update scheduler
		set 
		    date = :date, 
		    title = :title, 
		    comment = :comment, 
		    repeat = :repeat
		where id = :id
		`,
		sql.Named("id", task.ID),
		sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat))
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if count == 0 {
		return model.ErrTaskNotFound
	}

	return nil
}
