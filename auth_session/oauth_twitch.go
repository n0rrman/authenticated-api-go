package auth_session

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"golang.org/x/oauth2"
)

func (s *AuthSession) TwitchInit() {
	s.twitch = &oauth2.Config{
		ClientID:     os.Getenv("LINE_CLIENT_ID"),
		ClientSecret: os.Getenv("LINE_SECRET"),
		RedirectURL:  s.baseURL + "/auth/twitch/callback",
		Scopes:       []string{"profile"},
		Endpoint: oauth2.Endpoint{
			TokenURL: "https://id.twitch.tv/oauth2/token",
			AuthURL:  "https://id.twitch.tv/oauth2/authorize",
		},
	}
}

func (s *AuthSession) TwitchAuth(ctx echo.Context) error {
	s.authVerifier = oauth2.GenerateVerifier()
	url := s.twitch.AuthCodeURL("state", oauth2.S256ChallengeOption(s.authVerifier))

	return ctx.Redirect(http.StatusFound, url)
}

func (s *AuthSession) TwitchCallback(ctx echo.Context) error {
	// Get token
	tok, err := s.twitch.Exchange(
		context.Background(),
		ctx.QueryParam("code"),
		oauth2.VerifierOption(s.authVerifier),
	)
	if err != nil {
		return ctx.Redirect(http.StatusMovedPermanently, s.failedURL)
	}

	// Get data and read body
	client := s.twitch.Client(context.Background(), tok)
	resp, _ := client.Get("")
	body, _ := io.ReadAll(resp.Body)

	fmt.Println(string(body))

	// // Get vars from body
	// var result map[string]string
	// json.Unmarshal(body, &result)
	// id := result["userId"]
	// username := result["displayName"]
	// picture := result["pictureUrl"]
	// method := "line"

	// // // Login user
	// s.signin(ctx, id, username, picture, method)
	return ctx.Redirect(http.StatusMovedPermanently, s.successURL)
}
