package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	SPACE string = " "
)

var (
	SHELL_BUILTIN_COMMANDS map[string]bool = map[string]bool{
		"type": true,
		"exit": true,
		"echo": true,
	}
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
		case "type":
			processTypeCommand(args[0])
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

func processTypeCommand(commandName string) {
	_, ok := SHELL_BUILTIN_COMMANDS[commandName]

	if ok {
		fmt.Printf("%s is a shell builtin\n", commandName)
		return
	}

	execuatblePaths, isPathEnvSet := os.LookupEnv("PATH")

	if isPathEnvSet {

		for path := range strings.SplitSeq(execuatblePaths, ":") {
			ok, commandFullPath := isExecutable(path, commandName)
			if ok {
				fmt.Printf("%s is %s\n", commandName, commandFullPath)
				return
			}
		}

	}

	fmt.Printf("%s: not found\n", commandName)

}

func isExecutable(directoryPath string, target string) (bool, string) {

	info, err := os.Stat(directoryPath)
	if err != nil {
		if os.IsNotExist(err) {
			// fmt.Printf("Invalid directory path: %s\n", directoryPath)
		}
		return false, ""
	}

	if !info.IsDir() {
		return false, ""
	}

	matched := false
	fullPath := ""

	filepath.Walk(directoryPath, func(path string, info os.FileInfo, err error) error {

		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		mode := info.Mode()

		if target == info.Name() && mode&0100 != 0 {
			matched = true
			fullPath = path
			return nil
		}

		return nil
	})

	return matched, fullPath
}
