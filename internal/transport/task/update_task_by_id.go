package task

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/msmkdenis/go-task-tracker/internal/model"
	"github.com/msmkdenis/go-task-tracker/internal/transport/dto"
)

func (h *Handlers) UpdateTaskByID(c echo.Context) error {
	var taskUpdateRequest dto.PutTaskRequest
	err := c.Bind(&taskUpdateRequest)
	if err != nil {
		slog.Error("failed to bind task", slog.String("error", err.Error()))
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "failed to read body"})
	}

	if len(taskUpdateRequest.ID) == 0 {
		slog.Info("task id is empty")
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "id is required"})
	}

	id, errParseID := strconv.Atoi(taskUpdateRequest.ID)
	if errParseID != nil {
		slog.Info("failed to parse id", slog.String("error", errParseID.Error()))
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "wrong id format"})
	}

	if len(taskUpdateRequest.Title) == 0 {
		slog.Info("task title is empty")
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "title is required"})
	}

	if len(taskUpdateRequest.Repeat) != 0 && !repeatFormat(taskUpdateRequest.Repeat) {
		slog.Info("repeat format is wrong")
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "wrong repeat format"})
	}

	date := time.Now()
	if len(taskUpdateRequest.Date) != 0 {
		dateRequest, errParseDate := time.Parse("20060102", taskUpdateRequest.Date)
		if errParseDate != nil {
			slog.Info("failed to parse date", slog.String("error", errParseDate.Error()))
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "date format must be: 20060102"})
		}
		date = dateRequest
	}

	if date.Format("20060102") < time.Now().Format("20060102") {
		if len(taskUpdateRequest.Repeat) == 0 {
			date = time.Now()
		} else {
			date = calcNextDate(time.Now(), date, taskUpdateRequest.Repeat)
		}
	}

	task := model.Task{
		ID:      int64(id),
		Date:    date.Format("20060102"),
		Title:   taskUpdateRequest.Title,
		Comment: taskUpdateRequest.Comment,
		Repeat:  taskUpdateRequest.Repeat,
	}

	fmt.Println(task)
	err = h.Service.UpdateTaskByID(task)
	if errors.Is(err, model.ErrTaskNotFound) {
		slog.Info("task not found", slog.String("error", err.Error()))
		return c.JSON(http.StatusNotFound, map[string]string{"error": "task not found"})
	}

	if err != nil {
		slog.Error("failed to update task", slog.String("error", err.Error()))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to update task"})
	}

	return c.JSON(http.StatusOK, struct{}{})
}
