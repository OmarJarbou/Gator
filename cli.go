package main

import (
	"fmt"
)

func cli(stt *state, cmds *commands, args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("YOU SHOULD ENTER TWO KEYWORDS AT LEAST! (PROGRAM + COMMAND NAME)")
	}
	cmd := cmds.commandMapping(args[1], args[2:])
	err := cmds.run(stt, cmd)
	if err != nil {
		return err
	}
	return nil
}
