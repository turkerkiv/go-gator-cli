package main

import (
	"github.com/turkerkiv/gator/internal/config"
	"github.com/turkerkiv/gator/internal/database"
)

type state struct {
	db     *database.Queries
	config *config.Config
}
