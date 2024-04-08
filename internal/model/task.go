package model

import "time"

type Task struct {
	ID      int64     `db:"id"`
	Date    time.Time `db:"date"`
	Title   string    `db:"title"`
	Comment string    `db:"comment"`
	Repeat  string    `db:"repeat"`
}
