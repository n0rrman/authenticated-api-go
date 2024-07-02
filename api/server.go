package main

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

func main() {

	// Init server & session
	e := echo.New()
	s := session{name: "session"}
	s.init()

	// Routes
	e.GET("/create-session", func(c echo.Context) error {
		s.login(c)
		return c.NoContent(http.StatusOK)
	})

	e.GET("/read-session", func(c echo.Context) error {
		auth := s.isAuthenticated(c)
		if auth {
			return c.String(http.StatusOK, "logged in")
		} else {
			return c.String(http.StatusForbidden, "not logged in")
		}
	})

	// Start server
	port := "80"
	fmt.Println("Listening on port " + port)
	e.Logger.Fatal(e.Start(":" + port))
}
