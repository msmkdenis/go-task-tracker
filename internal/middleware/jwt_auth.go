package middleware

import (
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/msmkdenis/go-task-tracker/pkg/jwtgen"
)

type JWTAuth struct {
	jwtManager *jwtgen.JWTManager
	secret     string
}

func InitJWTAuth(jwtManager *jwtgen.JWTManager, secret string) *JWTAuth {
	return &JWTAuth{
		jwtManager: jwtManager,
		secret:     secret,
	}
}

func (j *JWTAuth) JWTAuth() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if j.secret == "" {
				return next(c)
			}
			cookie, err := c.Request().Cookie(j.jwtManager.TokenName)
			if err != nil {
				slog.Info("no cookie", slog.String("error", err.Error()))
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "cookie not found"})
			}
			err = j.jwtManager.IsValid(cookie.Value)
			if err != nil {
				slog.Info("wrong password", slog.String("error", err.Error()))
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "wrong password"})
			}
			return next(c)
		}
	}
}
