package main

import (
	"errors"
	"fmt"
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) < 1 {
		return errors.New("login expects a single argument, the username")
	}

	if err := s.config.SetUser(cmd.args[0]); err != nil {
		return err
	}

	fmt.Println("user successfully set")
	return nil
}
