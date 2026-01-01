package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {

	prompt, reader := "$ ", bufio.NewReader(os.Stdin)
	for {
		repl(prompt, reader)
	}
}

func repl(prompt string, reader *bufio.Reader) {

	fmt.Print(prompt)
	command, err := reader.ReadString('\n')

	if err != nil {
		fmt.Println(err)
	}

	commandName := strings.TrimRight(command, "\n")
	fmt.Printf("%s: command not found\n", commandName)

}
