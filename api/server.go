package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {

	// Init server & session handler
	e := echo.New()
	s := session{
		name:       "session",
		successURL: "/auth/status",
		failedURL:  "/dev/failed",
		expDays:    7,
	}
	s.init()

	// Middleware
	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.CORS())

	// Routes
	//-- Session
	e.GET("/auth/status", s.statusCheck)
	e.GET("/auth/signout", s.signout)
	//-- OAuth
	e.GET("/auth/github", s.githubAuth)
	e.GET("/auth/github/callback", s.githubCallback)
	e.GET("/auth/google", s.googleAuth)
	e.GET("/auth/google/callback", s.googleCallback)
	//-- Dev testing
	e.GET("/dev/success", devSuccess)
	e.GET("/dev/failed", devFailed)

	// Start server
	e.Logger.Fatal(e.Start(":80"))
}
