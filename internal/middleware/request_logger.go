package middleware

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

// RequestLogger represents request logger middleware.
type (
	RequestLogger struct{}

	responseData struct {
		status int
		size   int
	}

	loggingResponseWriter struct {
		http.ResponseWriter
		responseData *responseData
	}
)

// NewRequestLogger returns a new instance of RequestLogger.
func NewRequestLogger() *RequestLogger {
	return &RequestLogger{}
}

// Write implements the http.ResponseWriter interface.
func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b)
	r.responseData.size += size
	return size, err
}

// WriteHeader implements the http.ResponseWriter interface.
func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
	r.responseData.status = statusCode
}

// Get returns a middleware that logs each HTTP request.
func (r *RequestLogger) Get() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()

			uri := c.Request().RequestURI

			method := c.Request().Method

			duration := time.Since(start)

			responseData := &responseData{}

			lw := loggingResponseWriter{
				ResponseWriter: c.Response().Writer,
				responseData:   responseData,
			}

			c.Response().Writer = &lw

			err := next(c)

			slog.Info("request_logger",
				slog.String("URI", uri),
				slog.String("method", method),
				slog.Duration("duration", duration),
				slog.Int("response_code", responseData.status),
				slog.Int("response_body_size", responseData.size),
			)
			return err
		}
	}
}
