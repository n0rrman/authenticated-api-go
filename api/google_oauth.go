package main

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"golang.org/x/oauth2"
)

var google *oauth2.Config = &oauth2.Config{
	ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
	ClientSecret: os.Getenv("GOOGLE_SECRET"),
	RedirectURL:  os.Getenv("BASE_URL") + "/auth/google/callback",
	Scopes:       []string{"profile", "email"},
	Endpoint: oauth2.Endpoint{
		TokenURL: "https://oauth2.googleapis.com/token",
		AuthURL:  "https://accounts.google.com/o/oauth2/v2/auth",
	},
}

func (s *session) googleAuth(ctx echo.Context) error {
	s.authVerifier = oauth2.GenerateVerifier()
	url := google.AuthCodeURL("state", oauth2.S256ChallengeOption(s.authVerifier))

	return ctx.Redirect(http.StatusFound, url)
}

func (s *session) googleCallback(ctx echo.Context) error {
	// Get token
	tok, err := google.Exchange(
		context.Background(),
		ctx.QueryParam("code"),
		oauth2.VerifierOption(s.authVerifier),
	)
	if err != nil {
		return ctx.Redirect(http.StatusMovedPermanently, s.failedURL)
	}

	// Get data and read body
	client := google.Client(context.Background(), tok)
	resp, _ := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
	body, _ := io.ReadAll(resp.Body)

	// Get vars from body
	var result map[string]string
	json.Unmarshal(body, &result)
	id := result["sub"]
	username := result["email"]
	picture := result["picture"]
	method := "google"

	// Login user
	s.signin(ctx, id, username, picture, method)
	return ctx.Redirect(http.StatusMovedPermanently, s.successURL)
}
