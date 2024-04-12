package task

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"github.com/msmkdenis/go-task-tracker/internal/model"
)

func (h *Handlers) DeleteTaskByID(c echo.Context) error {
	id, err := strconv.ParseInt(c.QueryParam("id"), 10, 64)
	if err != nil {
		slog.Info("failed to parse id", slog.String("error", err.Error()))
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid id"})
	}

	err = h.Service.DeleteTaskByID(id)
	if errors.Is(err, model.ErrTaskNotFound) {
		slog.Info(fmt.Sprintf("task %d not found", id), slog.String("error", err.Error()))
		return c.JSON(http.StatusNotFound, map[string]string{"error": "task not found"})
	}

	if err != nil {
		slog.Error("failed to load task", slog.String("error", err.Error()))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to load task"})
	}

	return c.JSON(http.StatusOK, struct{}{})
}
