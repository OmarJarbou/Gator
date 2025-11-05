package main

import (
	"github.com/OmarJarbou/Gator/internal/config"
	"github.com/OmarJarbou/Gator/internal/database"
)

type state struct {
	DBQueries *database.Queries
	Config    *(config.Config)
}
