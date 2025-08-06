package main

import (
	"errors"
	"fmt"

	"github.com/OmarJarbou/Gator/internal/config"
)

type commands struct {
	cmdsMap map[string]func(*state, command) error
}

func commandMapping(cmd string, args []string) command {
	var cmnd command
	switch cmd {
	case "login":
		cmnd = command{
			Name:      "login",
			Arguments: args,
		}
	default:
		return cmnd
	}
	return cmnd
}

func (cmds *commands) register(name string, f func(*state, command) error) {
	cmds.cmdsMap[name] = f
}

func (cmds *commands) run(state *state, cmd command) error {
	switch cmd.Name {
	case "login":
		err := cmds.cmdsMap["login"](state, cmd)
		if err != nil {
			return err
		}
	default:
		return errors.New("COMMAND NOT FOUND")
	}
	return nil
}

func handleLogin(state *state, cmd command) error {
	if len(cmd.Arguments) == 0 {
		return errors.New("THE LOGIN HANDLER EXPECTS A SINGLE ARGUMENT, THE USERNAME")
	}
	err := config.SetUser(cmd.Arguments[0], state.Config)
	if err != nil {
		return err
	}
	fmt.Println("user set to", cmd.Arguments[0])
	return nil
}
