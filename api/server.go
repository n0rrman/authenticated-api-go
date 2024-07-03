package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/labstack/echo/v4"
	"golang.org/x/oauth2"
)

func main() {

	// Init server & session
	e := echo.New()
	s := session{name: "session"}
	s.init()

	// Routes
	e.GET("/create-session", func(ctx echo.Context) error {
		s.login(ctx)
		return ctx.NoContent(http.StatusOK)
	})

	e.GET("/read-session", func(ctx echo.Context) error {
		auth := s.isAuthenticated(ctx)
		if auth {
			return ctx.String(http.StatusOK, "logged in")
		} else {
			return ctx.String(http.StatusForbidden, "not logged in")
		}
	})

	var verifier string
	conf := &oauth2.Config{
		ClientID:     os.Getenv("GH_CLIENT_ID"),
		ClientSecret: os.Getenv("GH_SECRET"),
		Endpoint: oauth2.Endpoint{
			TokenURL: "https://github.com/login/oauth/access_token",
			AuthURL:  "https://github.com/login/oauth/authorize",
		},
	}

	e.GET("/auth/github/callback", func(ctx echo.Context) error {
		// var code string

		code := ctx.QueryParam("code")

		tok, err := conf.Exchange(context.Background(), code, oauth2.VerifierOption(verifier))
		if err != nil {
			ctx.Redirect(http.StatusNotFound, "/error")

		}

		client := conf.Client(context.Background(), tok)
		resp, err := client.Get("https://api.github.com/user")

		if err != nil {
			ctx.Redirect(http.StatusNotFound, "/next-error")

		}

		body, err := io.ReadAll(resp.Body)

		if err != nil {
			ctx.Redirect(http.StatusNotFound, "/json-error")

		}

		var result map[string]any
		json.Unmarshal(body, &result)
		fmt.Println(strconv.FormatFloat(result["id"].(float64), 'f', -1, 64))
		fmt.Println(result["login"])
		fmt.Println(result["avatar_url"])

		return ctx.String(http.StatusAccepted, "success")
	})

	e.GET("/auth/github", func(ctx echo.Context) error {
		verifier = oauth2.GenerateVerifier()

		url := conf.AuthCodeURL("state", oauth2.S256ChallengeOption(verifier))

		ctx.Redirect(http.StatusFound, url)

		return ctx.String(http.StatusOK, "github callback")
	})

	port := "80"
	fmt.Println("Listening on port " + port)
	e.Logger.Fatal(e.Start(":" + port))
}
