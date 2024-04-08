package transport

import (
	"github.com/labstack/echo/v4"

	"github.com/msmkdenis/go-task-tracker/internal/model"
)

type TaskService interface {
	AddTask(task model.Task) (int64, error)
	GetTasks() ([]model.Task, error)
	GetTasksByDate(date string) ([]model.Task, error)
	GetTasksByTitle(title string) ([]model.Task, error)
}

type TaskHandlers struct {
	Service TaskService
}

func NewTaskHandlers(e *echo.Echo, service TaskService) *TaskHandlers {
	handler := &TaskHandlers{
		Service: service,
	}

	e.GET("/api/nextdate", handler.CalculateNextDate)
	e.POST("/api/task", handler.AddTask)
	e.GET("/api/tasks", handler.LoadTasks)

	return handler
}
