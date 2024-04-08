package repository

import (
	"database/sql"

	"github.com/msmkdenis/go-task-tracker/internal/model"
)

func (t *SQLiteTaskRepository) SelectAllByTitle(title string) ([]model.Task, error) {
	tasks := make([]model.Task, 0)
	rows, err := t.db.Query(
		`
		select
			id, 
		    date, 
		    title, 
		    comment, 
		    repeat
		from scheduler
		where title like :title or comment like :title
		order by date
		limit :limit
		`,
		sql.Named("title", title),
		sql.Named("limit", taskLoadLimit))
	if err != nil {
		return nil, err
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}
	defer rows.Close()

	for rows.Next() {
		var task model.Task
		if errScan := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat); errScan != nil {
			return nil, errScan
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}
