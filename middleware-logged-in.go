package main

import (
	"context"
	"errors"

	"github.com/OmarJarbou/Gator/internal/database"
)

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	return func(s *state, cmd command) error {
		currentUserName := s.Config.CurrentUserName
		currentUser, err := s.DBQueries.GetUser(context.Background(), currentUserName)
		if err != nil {
			return errors.New("FAILED TO GET CURRENT USER: " + err.Error())
		}
		err2 := handler(s, cmd, currentUser)
		if err2 != nil {
			return err2
		}
		return nil
	}
}
