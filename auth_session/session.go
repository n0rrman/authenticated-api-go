package auth_session

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
	"github.com/rbcervilla/redisstore/v9"
	"github.com/redis/go-redis/v9"
)

type session struct {
	name         string
	authVerifier string
	successURL   string
	failedURL    string
	expDays      int
	client       *redis.Client
	store        *redisstore.RedisStore
}

func (s *session) init() {
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
		Domain:   os.Getenv("BASE_URL"),
		HttpOnly: true,
		MaxAge:   86400 * s.expDays,
	})
}

func (s *session) signin(ctx echo.Context, id string, username string, pictureURL string, method string) {
	sess, _ := s.store.Get(ctx.Request(), s.name)

	sess.Values["authenticated"] = true
	sess.Values["id"] = id
	sess.Values["username"] = username
	sess.Values["pictureURL"] = pictureURL
	sess.Values["method"] = method

	fmt.Println(method, "id:", id, "username:", username, pictureURL, "signed in")
	sess.Save(ctx.Request(), ctx.Response())
}

func (s *session) signout(ctx echo.Context) error {
	sess, _ := s.store.Get(ctx.Request(), s.name)

	sess.Options.MaxAge = -1

	fmt.Println(sess.Values["method"], "id:", sess.Values["id"], "username:", sess.Values["username"], "signed out")
	sess.Save(ctx.Request(), ctx.Response())

	return ctx.String(http.StatusOK, "signing out")
}

func (s *session) isAuthenticated(c echo.Context) bool {
	sess, _ := s.store.Get(c.Request(), s.name)
	return sess.Values["authenticated"] != nil
}

func (s *session) statusCheck(ctx echo.Context) error {
	auth := s.isAuthenticated(ctx)
	if auth {
		return ctx.String(http.StatusOK, "signed in")
	} else {
		return ctx.String(http.StatusForbidden, "not signed in")
	}
}
