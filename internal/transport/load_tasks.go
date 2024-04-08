package transport

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/msmkdenis/go-task-tracker/internal/model"
	"github.com/msmkdenis/go-task-tracker/internal/transport/dto"
)

func (h *TaskHandlers) LoadTasks(c echo.Context) error {
	var (
		tasks []model.Task
		err   error
	)

	searchParam := c.QueryParam("search")
	if len(searchParam) != 0 {
		date, errParse := time.Parse("02.01.2006", searchParam)
		if errParse != nil {
			searchParam = "%" + searchParam + "%"
			tasks, err = h.Service.GetTasksByTitle(searchParam)
		} else {
			tasks, err = h.Service.GetTasksByDate(date.Format("20060102"))
		}
	} else {
		tasks, err = h.Service.GetTasks()
	}

	if err != nil {
		slog.Error("failed to load tasks", slog.String("error", err.Error()))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to load tasks"})
	}

	taskResponse := make([]dto.TaskResponse, 0, len(tasks))
	for _, task := range tasks {
		taskResponse = append(taskResponse, dto.ToTaskResponse(task))
	}

	return c.JSON(http.StatusOK, map[string][]dto.TaskResponse{"tasks": taskResponse})
}
