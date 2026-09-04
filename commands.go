package main

import (
	"errors"
	"fmt"
)

type commands struct {
	commands map[string]func(*state, command) error
}

type command struct {
	name string
	args []string
}

func newCommands() commands {
	return commands{
		commands: make(map[string]func(*state, command) error),
	}
}

func (cs *commands) run(s *state, cmd command) error {
	if val, ok := cs.commands[cmd.name]; ok {
		if err := val(s, cmd); err != nil {
			return err
		}
		return nil
	} else {
		return errors.New("command not found")
	}
}

func (cs *commands) register(name string, f func(*state, command) error) {

	if _, ok := cs.commands[name]; !ok {
		cs.commands[name] = f
	} else {
		fmt.Println("command already exists, did nothing")
	}
}
