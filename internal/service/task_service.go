package service

import "github.com/msmkdenis/go-task-tracker/internal/model"

type TaskRepository interface {
	Insert(task model.Task) (int64, error)
	SelectAll() ([]model.Task, error)
	SelectAllByTitle(title string) ([]model.Task, error)
	SelectAllByDate(date string) ([]model.Task, error)
	SelectByID(id int64) (model.Task, error)
	UpdateByID(task model.Task) error
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

func (t *TaskService) GetTasks() ([]model.Task, error) {
	return t.Repository.SelectAll()
}

func (t *TaskService) GetTasksByDate(date string) ([]model.Task, error) {
	return t.Repository.SelectAllByDate(date)
}

func (t *TaskService) GetTasksByTitle(title string) ([]model.Task, error) {
	return t.Repository.SelectAllByTitle(title)
}

func (t *TaskService) GetTaskByID(id int64) (model.Task, error) {
	return t.Repository.SelectByID(id)
}

func (t *TaskService) UpdateTaskByID(task model.Task) error {
	return t.Repository.UpdateByID(task)
}
