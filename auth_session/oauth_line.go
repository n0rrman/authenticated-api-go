package auth_session

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"golang.org/x/oauth2"
)

func (s *AuthSession) LineInit() {
	s.line = &oauth2.Config{
		ClientID:     os.Getenv("LINE_CLIENT_ID"),
		ClientSecret: os.Getenv("LINE_SECRET"),
		RedirectURL:  s.baseURL + "/auth/line/callback",
		Scopes:       []string{"profile"},
		Endpoint: oauth2.Endpoint{
			TokenURL: "https://api.line.me/oauth2/v2.1/token",
			AuthURL:  "https://access.line.me/oauth2/v2.1/authorize",
		},
	}
}

func (s *AuthSession) LineAuth(ctx echo.Context) error {
	s.authVerifier = oauth2.GenerateVerifier()
	url := s.line.AuthCodeURL("state", oauth2.S256ChallengeOption(s.authVerifier))

	return ctx.Redirect(http.StatusFound, url)
}

func (s *AuthSession) LineCallback(ctx echo.Context) error {
	// Get token
	tok, err := s.line.Exchange(
		context.Background(),
		ctx.QueryParam("code"),
		oauth2.VerifierOption(s.authVerifier),
	)
	if err != nil {
		return ctx.Redirect(http.StatusMovedPermanently, s.failedURL)
	}

	// Get data and read body
	client := s.line.Client(context.Background(), tok)
	resp, _ := client.Get("https://api.line.me/v2/profile")
	body, _ := io.ReadAll(resp.Body)

	// Get vars from body
	var result map[string]string
	json.Unmarshal(body, &result)
	id := result["userId"]
	username := result["displayName"]
	picture := result["pictureUrl"]
	method := "line"

	// // Login user
	s.SignIn(ctx, id, username, picture, method)
	return ctx.Redirect(http.StatusMovedPermanently, s.successURL)
}
