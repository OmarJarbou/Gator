package main

import (
	"fmt"
	"github.com/OmarJarbou/Gator/internal/config"
	_ "github.com/lib/pq"
	"os"
)

func main() {
	confg := config.Read()
	stt := state{
		Config: &confg,
	}

	cmds := commands{
		cmdsMap: make(map[string]func(*state, command) error),
	}

	cmds.register("login", handleLogin)

	err := cli(&stt, &cmds, os.Args)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
