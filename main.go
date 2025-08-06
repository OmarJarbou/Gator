package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/OmarJarbou/Gator/internal/config"
	"github.com/OmarJarbou/Gator/internal/database"
	_ "github.com/lib/pq"
)

func main() {
	confg := config.Read()

	db, err := sql.Open("postgres", confg.DbURL)
	if err != nil {
		fmt.Errorf("failed to open database: %v", err)
	}

	dbQueries := database.New(db)

	stt := state{
		DBQueries: dbQueries,
		Config:    &confg,
	}

	cmds := commands{
		cmdsMap: make(map[string]func(*state, command) error),
	}

	cmds.register("login", handleLogin)

	err2 := cli(&stt, &cmds, os.Args)
	if err2 != nil {
		fmt.Println(err2)
		os.Exit(1)
	}
}
