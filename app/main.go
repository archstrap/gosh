package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {

	prompt, reader := "$ ", bufio.NewReader(os.Stdin)
	repl(prompt, reader)
}

func repl(prompt string, reader *bufio.Reader) {

	for {
		fmt.Print(prompt)
		command, err := reader.ReadString('\n')

		if err != nil {
			fmt.Println(err)
		}

		commandName := strings.TrimRight(command, "\n")
		if "exit" == commandName {
			return
		}
		fmt.Printf("%s: command not found\n", commandName)
	}
}
