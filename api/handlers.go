package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func devFailed(ctx echo.Context) error {
	return ctx.String(http.StatusBadGateway, "failed")
}

func devSuccess(ctx echo.Context) error {
	return ctx.String(http.StatusAccepted, "success")
}

func (s *session) statusCheck(ctx echo.Context) error {
	auth := s.isAuthenticated(ctx)
	if auth {
		return ctx.String(http.StatusOK, "signed in")
	} else {
		return ctx.String(http.StatusForbidden, "not signed in")
	}
}
