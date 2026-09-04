package main

import (
	"fmt"
	"os"

	"github.com/turkerkiv/gator/internal/config"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		fmt.Println(err)
	}

	appState := state{
		config: &cfg,
	}

	commands := newCommands()
	commands.register("login", handlerLogin)

	if len(os.Args) < 2 {
		fmt.Println("there must be at least a command and potentially arguments")
		os.Exit(1)
	}
	userCommand := command{
		name: os.Args[1],
		args: os.Args[2:],
	}

	if err := commands.run(&appState, userCommand); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// cfg.SetUser("turker")

	// cfg, err = config.Read()
	// if err != nil {
	// 	fmt.Println(err)
	// }

	// fmt.Println(cfg.CurrentUserName, cfg.DbUrl)

}
