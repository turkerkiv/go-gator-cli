package main

import (
	"fmt"

	"github.com/turkerkiv/gator/internal/config"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(cfg.CurrentUserName, cfg.DbUrl)
	cfg.SetUser("turker")

	cfg, err = config.Read()
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(cfg.CurrentUserName, cfg.DbUrl)
}
