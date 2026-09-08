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

func handlerAddFeed(s *state, cmd command, usr database.User) error {
	if len(cmd.args) < 2 {
		return errors.New("addFeed expects a two arguments, name of feed and url of feed")
	}

	name := cmd.args[0]
	url := cmd.args[1]
	params := database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      name,
		Url:       url,
		UserID:    usr.ID,
	}
	feed, err := s.db.CreateFeed(context.Background(), params)
	if err != nil {
		return err
	}

	cmd.args = cmd.args[1:]
	if err := handlerFollow(s, cmd, usr); err != nil {
		return err
	}

	fmt.Println(feed)
	return nil
}

func handlerListFeeds(s *state, cmd command) error {
	feeds, err := s.db.ListFeeds(context.Background())
	if err != nil {
		return err
	}

	for _, feed := range feeds {
		usr, err := s.db.GetUserById(context.Background(), feed.UserID)
		if err != nil {
			return nil
		}

		fmt.Printf("* %s - %s - %s\n", feed.Name, feed.Url, usr.Name)
	}
	return nil
}

func handlerFollow(s *state, cmd command, usr database.User) error {
	if len(cmd.args) < 1 {
		return errors.New("follow expects a 1 arguments, the url.")
	}
	url := cmd.args[0]

	feed, err := s.db.GetFeed(context.Background(), url)
	if err != nil {
		return err
	}

	params := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    usr.ID,
		FeedID:    feed.ID,
	}

	follows, err := s.db.CreateFeedFollow(context.Background(), params)
	if err != nil {
		return err
	}
	fmt.Printf("%s - %s\n", follows.FeedName, follows.UserName)
	return nil
}

func handlerFollowing(s *state, cmd command, usr database.User) error {
	follows, err := s.db.GetFeedFollowsForUser(context.Background(), usr.ID)
	if err != nil {
		return err
	}

	for _, val := range follows {
		fmt.Printf("%s %s\n", val.UserName, val.FeedName)
	}
	return nil
}

func handlerUnfollow(s *state, cmd command, usr database.User) error {
	if len(cmd.args) < 1 {
		return errors.New("unfollow expects a 1 arguments, the url.")
	}
	url := cmd.args[0]

	feed, err := s.db.GetFeed(context.Background(), url)
	if err != nil {
		return err
	}

	params := database.UnfollowFeedParams{
		UserID: usr.ID,
		FeedID: feed.ID,
	}

	if err := s.db.UnfollowFeed(context.Background(), params); err != nil {
		return err
	}

	fmt.Println("unfollow successful")
	return nil
}
