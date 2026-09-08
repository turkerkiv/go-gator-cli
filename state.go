package main

import (
	"github.com/turkerkiv/go-gator-cli/internal/config"
	"github.com/turkerkiv/go-gator-cli/internal/database"
)

type state struct {
	db     *database.Queries
	config *config.Config
}
