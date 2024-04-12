package task

import (
	"github.com/labstack/echo/v4"

	"github.com/msmkdenis/go-task-tracker/internal/middleware"
	"github.com/msmkdenis/go-task-tracker/internal/model"
)

type Service interface {
	AddTask(task model.Task) (int64, error)
	GetTasks() ([]model.Task, error)
	GetTasksByDate(date string) ([]model.Task, error)
	GetTasksByTitle(title string) ([]model.Task, error)
	GetTaskByID(id int64) (model.Task, error)
	UpdateTaskByID(task model.Task) error
	DeleteTaskByID(id int64) error
}

type Handlers struct {
	Service Service
	jwtAuth *middleware.JWTAuth
}

func NewHandlers(e *echo.Echo, service Service, auth *middleware.JWTAuth) *Handlers {
	handler := &Handlers{
		Service: service,
		jwtAuth: auth,
	}

	e.POST("/api/task", handler.AddTask, handler.jwtAuth.JWTAuth())
	e.POST("/api/task/done", handler.CompleteTaskByID, handler.jwtAuth.JWTAuth())

	e.PUT("/api/task", handler.UpdateTaskByID, handler.jwtAuth.JWTAuth())

	e.GET("/api/tasks", handler.LoadTasks, handler.jwtAuth.JWTAuth())
	e.GET("/api/task", handler.LoadTaskByID, handler.jwtAuth.JWTAuth())
	e.GET("/api/nextdate", handler.CalculateNextDate)

	e.DELETE("/api/task", handler.DeleteTaskByID, handler.jwtAuth.JWTAuth())

	return handler
}
