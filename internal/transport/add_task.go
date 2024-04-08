package transport

import (
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/msmkdenis/go-task-tracker/internal/model"
	"github.com/msmkdenis/go-task-tracker/internal/transport/dto"
)

func (h *TaskHandlers) AddTask(c echo.Context) error {
	var taskRequest dto.PostTaskRequest
	err := c.Bind(&taskRequest)
	if err != nil {
		slog.Error("failed to bind task", slog.String("error", err.Error()))
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "failed to read body"})
	}

	if len(taskRequest.Title) == 0 {
		slog.Info("task title is empty")
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "title is required"})
	}

	if len(taskRequest.Repeat) != 0 && !repeatFormat(taskRequest.Repeat) {
		slog.Info("repeat format is wrong")
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "wrong repeat format"})
	}

	date := time.Now()
	if len(taskRequest.Date) != 0 {
		dateRequest, errParseDate := time.Parse("20060102", taskRequest.Date)
		if errParseDate != nil {
			slog.Info("failed to parse date", slog.String("error", errParseDate.Error()))
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "date format must be: 20060102"})
		}
		date = dateRequest
	}

	if date.Format("20060102") < time.Now().Format("20060102") {
		if len(taskRequest.Repeat) == 0 {
			date = time.Now()
		} else {
			date = calcNextDate(time.Now(), date, taskRequest.Repeat)
		}
	}

	task := model.Task{
		Date:    date.Format("20060102"),
		Title:   taskRequest.Title,
		Comment: taskRequest.Comment,
		Repeat:  taskRequest.Repeat,
	}

	id, err := h.Service.AddTask(task)
	if err != nil {
		slog.Error("failed to add task", slog.String("error", err.Error()))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to add task"})
	}

	return c.JSON(http.StatusOK, map[string]string{"id": strconv.FormatInt(id, 10)})
}
