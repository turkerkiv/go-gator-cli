package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
	"github.com/turkerkiv/gator/internal/config"
	"github.com/turkerkiv/gator/internal/database"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		fmt.Println(err)
	}

	db, err2 := sql.Open("postgres", cfg.DbUrl)
	if err2 != nil {
		fmt.Println(err2)
	}
	dbQueries := database.New(db)

	appState := state{
		db:     dbQueries,
		config: &cfg,
	}

	commands := newCommands()
	commands.register("login", handlerLogin)
	commands.register("register", handlerRegister)
	commands.register("reset", handlerReset)
	commands.register("users", handlerListUsers)
	commands.register("agg", handlerAgg)
	commands.register("addfeed", middlewareLoggedIn(handlerAddFeed))
	commands.register("feeds", handlerListFeeds)
	commands.register("follow", middlewareLoggedIn(handlerFollow))
	commands.register("following", middlewareLoggedIn(handlerFollowing))

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
