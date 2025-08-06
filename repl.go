package main

import (
	"bufio"
	"fmt"
	"os"
)

func repl(stt *state, cmds *commands) {
	scanner := bufio.NewScanner(os.Stdin)

	for nextScan(scanner) {
		full_line := scanner.Text()
		line_words := cleanInput(full_line)
		if len(line_words) < 2 {
			fmt.Println("You should enter two keywords at least! (gator + command name)")
			continue
		}
		cmd := commandMapping(line_words[1], line_words[2:])
		err := cmds.run(stt, cmd)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func nextScan(scanner *bufio.Scanner) bool {
	fmt.Print("Gator> ")
	return scanner.Scan()
}

func cleanInput(input string) []string {
	cleaned_input := []string{}
	temp := ""
	runes := []rune(input)
	for i := 0; i < len(runes); i++ {
		if runes[i] != rune(' ') {
			temp += string(runes[i])
			// access next char:
			if (i+1 == len(runes)) || (runes[i+1] == rune(' ')) {
				cleaned_input = append(cleaned_input, temp)
				temp = ""
			}
		} else {
			temp = ""
		}
	}
	return cleaned_input
}
