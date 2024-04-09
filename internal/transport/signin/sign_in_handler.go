package signin

import (
	"github.com/labstack/echo/v4"

	"github.com/msmkdenis/go-task-tracker/pkg/jwtgen"
)

type Handlers struct {
	jwtManager *jwtgen.JWTManager
	password   string
}

func NewHandlers(e *echo.Echo, jwtManager *jwtgen.JWTManager, password string) *Handlers {
	handler := &Handlers{
		jwtManager: jwtManager,
		password:   password,
	}

	e.POST("/api/signin", handler.SignIn)

	return handler
}
