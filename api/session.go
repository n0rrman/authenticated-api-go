package main

import (
	"context"
	"log"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
	"github.com/rbcervilla/redisstore/v9"
	"github.com/redis/go-redis/v9"
)

type session struct {
	name   string
	client *redis.Client
	store  *redisstore.RedisStore
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
		Domain:   "localhost",
		HttpOnly: true,
		// MaxAge: 86400 * 20,
		MaxAge: 10,
	})
}

func (s *session) login(c echo.Context) {
	sess, _ := s.store.Get(c.Request(), s.name)
	sess.Values["authenticated"] = true
	sess.Save(c.Request(), c.Response())
}

func (s *session) isAuthenticated(c echo.Context) bool {
	sess, _ := s.store.Get(c.Request(), s.name)
	return sess.Values["authenticated"] != nil
}
