package transport

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/msmkdenis/go-task-tracker/internal/model"

	"github.com/labstack/echo/v4"
)

func (h *TaskHandlers) CompleteTaskByID(c echo.Context) error {
	id, err := strconv.ParseInt(c.QueryParam("id"), 10, 64)
	if err != nil {
		slog.Info("failed to parse id", slog.String("error", err.Error()))
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid id"})
	}

	task, err := h.Service.GetTaskByID(id)
	if errors.Is(err, model.ErrTaskNotFound) {
		slog.Info(fmt.Sprintf("task %d not found", id), slog.String("error", err.Error()))
		return c.JSON(http.StatusNotFound, map[string]string{"error": "task not found"})
	}

	if err != nil {
		slog.Error("failed to load task", slog.String("error", err.Error()))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to load task"})
	}

	if len(task.Repeat) == 0 {
		err = h.Service.DeleteTaskByID(id)
		if err != nil {
			slog.Error("failed to delete done task", slog.String("error", err.Error()))
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to delete done task"})
		}
	}

	if len(task.Repeat) != 0 {
		taskDate, err := time.Parse("20060102", task.Date)
		if err != nil {
			slog.Info("failed to parse date", slog.String("error", err.Error()))
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "internal error"})
		}

		nextDate := calcNextDate(time.Now(), taskDate, task.Repeat)
		task.Date = nextDate.Format("20060102")

		err = h.Service.UpdateTaskByID(task)
		if err != nil {
			slog.Error("failed to update done task", slog.String("error", err.Error()))
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to update done task"})
		}
	}

	return c.JSON(http.StatusOK, struct{}{})
}
