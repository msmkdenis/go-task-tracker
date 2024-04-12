package signin

import (
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/msmkdenis/go-task-tracker/internal/transport/dto"
)

func (h *Handlers) SignIn(c echo.Context) error {
	var signInRequest dto.PostSignInRequest
	err := c.Bind(&signInRequest)
	if err != nil {
		slog.Error("failed to bind sign in request", slog.String("error", err.Error()))
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "failed to password"})
	}

	if signInRequest.Password != h.password {
		slog.Info("wrong password")
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "wrong password"})
	}

	token, err := h.jwtManager.BuildJWTString()
	if err != nil {
		slog.Error("failed to build jwt", slog.String("error", err.Error()))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to build token"})
	}

	return c.JSON(http.StatusOK, map[string]string{"token": token})
}
