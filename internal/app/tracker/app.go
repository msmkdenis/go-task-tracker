package tracker

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/msmkdenis/go-task-tracker/internal/config"
	"github.com/msmkdenis/go-task-tracker/internal/middleware"
	"github.com/msmkdenis/go-task-tracker/internal/repository"
	"github.com/msmkdenis/go-task-tracker/internal/service"
	"github.com/msmkdenis/go-task-tracker/internal/storage"
	"github.com/msmkdenis/go-task-tracker/internal/transport/signin"
	"github.com/msmkdenis/go-task-tracker/internal/transport/task"
	"github.com/msmkdenis/go-task-tracker/pkg/jwtgen"
	"github.com/msmkdenis/go-task-tracker/tests"
)

func Run(quitSignal chan os.Signal) {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{AddSource: true}))
	slog.SetDefault(logger)

	cfg := config.New()

	database := storage.NewSQLiteDB()
	migrations := storage.NewMigrations(database)
	err := migrations.Up()
	if err != nil {
		slog.Error("failed to run migrations", "error", err)
		os.Exit(1)
	}

	jwtManager := jwtgen.NewJWTManager(cfg.TokenName, cfg.Secret, cfg.TokenTTL, cfg.Salt)
	jwtAuth := middleware.InitJWTAuth(jwtManager, cfg.Secret)

	taskRepository := repository.NewSQLiteTaskRepository(database)
	taskService := service.NewTaskService(taskRepository)

	requestLogger := middleware.NewRequestLogger()

	e := echo.New()
	e.Use(requestLogger.Get())
	e.Static("/", "web")

	task.NewHandlers(e, taskService, jwtAuth)
	signin.NewHandlers(e, jwtManager, cfg.Secret)

	httpServerCtx, httpServerStopCtx := context.WithCancel(context.Background())

	go func() {
		slog.Info("staring server", slog.Int("localhost:", tests.Port))
		errStart := e.Start(cfg.URLServer)
		if errStart != nil && !errors.Is(errStart, http.ErrServerClosed) {
			slog.Error("failed to start server", slog.String("error", errStart.Error()))
			os.Exit(1)
		}
	}()

	quit := make(chan struct{})
	go func() {
		<-quitSignal
		close(quit)
	}()

	go func() {
		<-quit

		// Shutdown signal with grace period of 10 seconds
		shutdownCtx, cancel := context.WithTimeout(httpServerCtx, 10*time.Second)
		defer cancel()

		go func() {
			<-shutdownCtx.Done()
			if errors.Is(shutdownCtx.Err(), context.DeadlineExceeded) {
				slog.Error("graceful shutdown timed out.. forcing exit.")
				os.Exit(1)
			}
		}()

		// Trigger graceful shutdown
		slog.Info("Initiating graceful shutdown")
		if err := e.Shutdown(shutdownCtx); err != nil {
			slog.Error("Failed to shutdown server", slog.String("error", err.Error()))
			os.Exit(1)
		}
		httpServerStopCtx()
	}()

	<-httpServerCtx.Done()
}
