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
		fmt.Println("FAILED TO OPEN DATABASE:", err)
		os.Exit(1)
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
	cmds.register("register", handleRegister)
	cmds.register("reset", handleReset)
	cmds.register("users", handleListUsers)
	cmds.register("agg", handleAggregate)
	cmds.register("addfeed", middlewareLoggedIn(handleAddFeed))
	cmds.register("feeds", handleListFeeds)
	cmds.register("follow", middlewareLoggedIn(handleFollowFeed))
	cmds.register("following", middlewareLoggedIn(handleListFollowing))
	cmds.register("unfollow", middlewareLoggedIn(handleUnfollowFeed))

	err2 := cli(&stt, &cmds, os.Args)
	if err2 != nil {
		fmt.Println(err2)
		os.Exit(1)
	}
}
