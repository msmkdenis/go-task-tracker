package dto

import (
	"strconv"

	"github.com/msmkdenis/go-task-tracker/internal/model"
)

type PostTaskRequest struct {
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

type PutTaskRequest struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

type TaskResponse struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

func ToTaskResponse(t model.Task) TaskResponse {
	return TaskResponse{
		ID:      strconv.FormatInt(t.ID, 10),
		Date:    t.Date,
		Title:   t.Title,
		Comment: t.Comment,
		Repeat:  t.Repeat,
	}
}
