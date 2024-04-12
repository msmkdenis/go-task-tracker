package repository

import (
	"database/sql"

	"github.com/msmkdenis/go-task-tracker/internal/model"
)

func (t *SQLiteTaskRepository) DeleteByID(id int64) error {
	res, err := t.db.Exec(
		`
		delete from scheduler
		where id = :id
		`,
		sql.Named("id", id))
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
