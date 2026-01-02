package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

const (
	SPACE string = " "
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

		commandName, args := SplitCommandDetails(command)

		switch commandName {
		case "exit":
			return
		case "echo":
			processEchoCommand(args)
		default:
			fmt.Printf("%s: command not found\n", commandName)
		}

	}
}

func SplitCommandDetails(commandDetails string) (string, []string) {
	parts := strings.Split(strings.TrimRight(commandDetails, "\n"), SPACE)
	if len(parts) > 1 {
		return parts[0], parts[1:]
	}

	return parts[0], []string{}
}

func processEchoCommand(args []string) {
	fmt.Println(strings.Join(args, SPACE))
}
