package main

import (
	"github.com/OmarJarbou/Gator/internal/config"
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

	repl(&stt, &cmds)
}
