package main

import (
	"errors"
	"fmt"

	"context"
	"time"

	"github.com/OmarJarbou/Gator/internal/config"
	"github.com/OmarJarbou/Gator/internal/database"
	"github.com/google/uuid"
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
	case "register":
		cmnd = command{
			Name:      "register",
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
	_, ok := cmds.cmdsMap[cmd.Name]
	if !ok {
		return errors.New("COMMAND NOT FOUND")
	}
	err := cmds.cmdsMap[cmd.Name](state, cmd)
	if err != nil {
		return err
	}
	return nil
}

func handleLogin(state *state, cmd command) error {
	if len(cmd.Arguments) == 0 || len(cmd.Arguments) > 1 {
		return errors.New("THE LOGIN HANDLER EXPECTS A SINGLE ARGUMENT, THE USERNAME")
	}
	_, err := state.DBQueries.GetUser(context.Background(), cmd.Arguments[0])
	if err != nil {
		return errors.New("USER NOT FOUND: " + err.Error())
	}
	err2 := config.SetUser(cmd.Arguments[0], state.Config)
	if err2 != nil {
		return err2
	}
	fmt.Println("User set to", cmd.Arguments[0])
	return nil
}

func handleRegister(state *state, cmd command) error {
	if len(cmd.Arguments) == 0 || len(cmd.Arguments) > 1 {
		return errors.New("THE REGISTER HANDLER EXPECTS A SINGLE ARGUMENT, THE USERNAME")
	}
	_, err := state.DBQueries.GetUser(context.Background(), cmd.Arguments[0])
	if err == nil {
		return errors.New("USER WITH THIS NAME '" + cmd.Arguments[0] + "' ALREADY EXISTS")
	}

	user := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.Arguments[0],
	}

	usr, err2 := state.DBQueries.CreateUser(context.Background(), user)
	if err2 != nil {
		return errors.New("FAILED TO CREATE USER: " + err2.Error())
	}

	state.Config.CurrentUserName = cmd.Arguments[0]
	fmt.Println("User", usr.Name, "created successfully!")
	fmt.Println("Created at:", usr.CreatedAt)

	err3 := config.SetUser(cmd.Arguments[0], state.Config)
	if err3 != nil {
		return err3
	}
	fmt.Println("User set to", cmd.Arguments[0])

	return nil
}

func handleReset(state *state, cmd command) error {
	return nil
}
