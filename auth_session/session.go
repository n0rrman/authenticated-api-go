package auth_session

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
	"github.com/rbcervilla/redisstore/v9"
	"github.com/redis/go-redis/v9"
	"golang.org/x/oauth2"
)

type AuthSession struct {
	name         string
	authVerifier string
	successURL   string
	failedURL    string
	baseURL      string
	expDays      int
	client       *redis.Client
	store        *redisstore.RedisStore

	bnet    *oauth2.Config
	discord *oauth2.Config
	github  *oauth2.Config
	google  *oauth2.Config
	line    *oauth2.Config
	twitch  *oauth2.Config
}

func (s *AuthSession) Init() {
	// Redis Client
	s.client = redis.NewClient(&redis.Options{
		Addr: "session_storage:6379",
	})

	// Default RedisStore
	var err error
	s.store, err = redisstore.NewRedisStore(context.Background(), s.client)
	if err != nil {
		log.Fatal("failed to create redis store: ", err)
	}

	// Cookie options
	s.store.KeyPrefix("session_")
	s.store.Options(sessions.Options{
		Path:     "/",
		Domain:   s.baseURL,
		HttpOnly: true,
		MaxAge:   86400 * s.expDays,
	})

	s.BnetInit()
	s.DiscordInit()
	s.GithubInit()
	s.GoogleInit()
	s.LineInit()
	s.TwitchInit()
}

func (s *AuthSession) SignIn(ctx echo.Context, id string, username string, pictureURL string, method string) {
	sess, _ := s.store.Get(ctx.Request(), s.name)

	sess.Values["authenticated"] = true
	sess.Values["id"] = id
	sess.Values["username"] = username
	sess.Values["pictureURL"] = pictureURL
	sess.Values["method"] = method

	fmt.Println(method, "id:", id, "username:", username, pictureURL, "signed in")
	sess.Save(ctx.Request(), ctx.Response())
}

func (s *AuthSession) SignOut(ctx echo.Context) error {
	sess, _ := s.store.Get(ctx.Request(), s.name)

	sess.Options.MaxAge = -1

	fmt.Println(sess.Values["method"], "id:", sess.Values["id"], "username:", sess.Values["username"], "signed out")
	sess.Save(ctx.Request(), ctx.Response())

	return ctx.String(http.StatusOK, "signing out")
}

func (s *AuthSession) IsAuthenticated(c echo.Context) bool {
	sess, _ := s.store.Get(c.Request(), s.name)
	return sess.Values["authenticated"] != nil
}

func (s *AuthSession) StatusCheck(ctx echo.Context) error {
	auth := s.IsAuthenticated(ctx)
	if auth {
		return ctx.String(http.StatusOK, "signed in")
	} else {
		return ctx.String(http.StatusForbidden, "not signed in")
	}
}
