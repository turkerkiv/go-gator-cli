package main

import (
	"context"

	"github.com/turkerkiv/gator/internal/database"
)

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	return func(s *state, cmd command) error {
		usr, err := s.db.GetUser(context.Background(), s.config.CurrentUserName)
		if err != nil {
			return err
		}

		if err := handler(s, cmd, usr); err != nil {
			return err
		}
		return nil
	}
}
