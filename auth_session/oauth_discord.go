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

func (s *AuthSession) DiscordInit() {
	s.discord = &oauth2.Config{
		ClientID:     os.Getenv("DISCORD_CLIENT_ID"),
		ClientSecret: os.Getenv("DISCORD_SECRET"),
		RedirectURL:  s.baseURL + "/auth/discord/callback",
		Scopes:       []string{"identify"},
		Endpoint: oauth2.Endpoint{
			TokenURL: "https://discord.com/api/oauth2/token",
			AuthURL:  "https://discord.com/oauth2/authorize",
		},
	}
}

func (s *AuthSession) DiscordAuth(ctx echo.Context) error {
	s.authVerifier = oauth2.GenerateVerifier()
	url := s.discord.AuthCodeURL("state", oauth2.S256ChallengeOption(s.authVerifier))

	return ctx.Redirect(http.StatusFound, url)
}

func (s *AuthSession) DiscordCallback(ctx echo.Context) error {
	// Get token
	tok, err := s.discord.Exchange(
		context.Background(),
		ctx.QueryParam("code"),
		oauth2.VerifierOption(s.authVerifier),
	)
	if err != nil {
		return ctx.Redirect(http.StatusMovedPermanently, s.failedURL)
	}

	// Get data and read body
	client := s.discord.Client(context.Background(), tok)
	resp, _ := client.Get("https://discord.com/api/users/@me")
	body, _ := io.ReadAll(resp.Body)

	// Get vars from body
	var result map[string]string
	json.Unmarshal(body, &result)
	id := result["id"]
	username := result["username"]
	picture := "https://cdn.discordapp.com/avatars/" + id + "/" + result["avatar"]
	method := "discord"

	// Login user
	s.SignIn(ctx, id, username, picture, method)
	return ctx.Redirect(http.StatusMovedPermanently, s.successURL)
}
