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

var bnet *oauth2.Config = &oauth2.Config{
	ClientID:     os.Getenv("BNET_CLIENT_ID"),
	ClientSecret: os.Getenv("BNET_SECRET"),
	RedirectURL:  os.Getenv("BASE_URL") + "/auth/bnet/callback",
	Endpoint: oauth2.Endpoint{
		AuthURL:  "https://oauth.battle.net/authorize",
		TokenURL: "https://oauth.battle.net/token",
	},
}

func (s *session) bnetAuth(ctx echo.Context) error {
	s.authVerifier = oauth2.GenerateVerifier()
	url := bnet.AuthCodeURL("state", oauth2.S256ChallengeOption(s.authVerifier))

	return ctx.Redirect(http.StatusFound, url)
}

func (s *session) bnetCallback(ctx echo.Context) error {
	// Get token
	tok, err := bnet.Exchange(
		context.Background(),
		ctx.QueryParam("code"),
		oauth2.VerifierOption(s.authVerifier),
	)

	if err != nil {
		return ctx.Redirect(http.StatusMovedPermanently, s.failedURL)
	}

	// Get data and read body
	client := bnet.Client(context.Background(), tok)
	resp, _ := client.Get("https://battle.net/oauth/userinfo")
	body, _ := io.ReadAll(resp.Body)

	// Get vars from body
	var result map[string]string
	json.Unmarshal(body, &result)
	id := result["id"]
	username := result["battletag"]
	picture := ""
	method := "battle.net"

	// // Login user
	s.signin(ctx, id, username, picture, method)
	return ctx.Redirect(http.StatusMovedPermanently, s.successURL)
}
