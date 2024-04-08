package service

import "github.com/msmkdenis/go-task-tracker/internal/model"

type TaskRepository interface {
	Insert(task model.Task) (int64, error)
}

type TaskService struct {
	Repository TaskRepository
}

func NewTaskService(repository TaskRepository) *TaskService {
	return &TaskService{
		Repository: repository,
	}
}

func (t *TaskService) AddTask(task model.Task) (int64, error) {
	return t.Repository.Insert(task)
}