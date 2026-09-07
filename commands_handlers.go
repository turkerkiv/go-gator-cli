package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/turkerkiv/gator/internal/database"
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) < 1 {
		return errors.New("login expects a single argument, the username")
	}

	if _, err := s.db.GetUser(context.Background(), cmd.args[0]); err != nil {
		return fmt.Errorf("User not found: %v", err)
	}

	if err := s.config.SetUser(cmd.args[0]); err != nil {
		return err
	}

	fmt.Println("user successfully set")
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) < 1 {
		return errors.New("register expects a single argument, the username")
	}

	if _, err := s.db.GetUser(context.Background(), cmd.args[0]); err == nil {
		return fmt.Errorf("User already exists: %v", err)
	}

	params := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.args[0],
	}

	if _, err := s.db.CreateUser(context.Background(), params); err != nil {
		return err
	}

	s.config.SetUser(params.Name)
	fmt.Println("User created: " + s.config.CurrentUserName)
	return nil
}

func handlerReset(s *state, cmd command) error {
	if err := s.db.TruncateTable(context.Background()); err != nil {
		return err
	}
	fmt.Println("Table data deleted")
	return nil
}

func handlerListUsers(s *state, cmd command) error {
	users, err := s.db.ListUsers(context.Background())
	if err != nil {
		return err
	}

	for _, user := range users {
		currentUsrStr := ""
		if user.Name == s.config.CurrentUserName {
			currentUsrStr = "(current)"
		}
		fmt.Printf("* %s %s\n", user.Name, currentUsrStr)
	}
	return nil
}

func handlerAgg(s *state, cmd command) error {
	feed, err := fetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	if err != nil {
		return err
	}

	fmt.Println(feed)
	return nil
}
