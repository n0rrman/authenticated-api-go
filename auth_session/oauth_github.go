package auth_session

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/labstack/echo/v4"
	"golang.org/x/oauth2"
)

func (s *AuthSession) GithubInit() {
	s.github = &oauth2.Config{
		ClientID:     os.Getenv("GH_CLIENT_ID"),
		ClientSecret: os.Getenv("GH_SECRET"),
		Endpoint: oauth2.Endpoint{
			TokenURL: "https://github.com/login/oauth/access_token",
			AuthURL:  "https://github.com/login/oauth/authorize",
		},
	}
}

func (s *AuthSession) GithubAuth(ctx echo.Context) error {
	s.authVerifier = oauth2.GenerateVerifier()
	url := s.github.AuthCodeURL("state", oauth2.S256ChallengeOption(s.authVerifier))

	return ctx.Redirect(http.StatusFound, url)
}

func (s *AuthSession) GithubCallback(ctx echo.Context) error {
	// Get token
	tok, err := s.github.Exchange(
		context.Background(),
		ctx.QueryParam("code"),
		oauth2.VerifierOption(s.authVerifier),
	)
	if err != nil {
		return ctx.Redirect(http.StatusMovedPermanently, s.failedURL)
	}

	// Get data and read body
	client := s.github.Client(context.Background(), tok)
	resp, _ := client.Get("https://api.github.com/user")
	body, _ := io.ReadAll(resp.Body)

	// Get vars from body
	var result map[string]any
	json.Unmarshal(body, &result)
	id := strconv.FormatFloat(result["id"].(float64), 'f', -1, 64)
	username := result["login"].(string)
	picture := result["avatar_url"].(string)
	method := "github"

	// Login user
	s.SignIn(ctx, id, username, picture, method)
	return ctx.Redirect(http.StatusMovedPermanently, s.successURL)
}
