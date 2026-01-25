package main

import (
	"bufio"
	"fmt"
	// term "golang.org/x/term"
	"os"
	"strings"
)

var (
	ShellBuiltinCommands = map[string]bool{
		"type": true,
		"exit": true,
		"pwd":  true,
		"cd":   true,
		"echo": true,
	}
)

func main() {

	if err := loadShellRC(); err != nil {
		fmt.Println("Unable to open .shellrc")
		os.Exit(1)
	}

	// terminalFd := int(os.Stdin.Fd())
	// oldState, err := term.MakeRaw(terminalFd)
	// if err != nil {
	// 	panic(err)
	// }
	//
	// defer term.Restore(terminalFd, oldState)
	//
	prompt, reader := os.Getenv("PS"), bufio.NewReader(os.Stdin)
	repl(prompt, reader)
}

func repl(prompt string, reader *bufio.Reader) {

	for {
		fmt.Print(prompt)
		command, err := reader.ReadString('\n')
		if err != nil {
			panic(err)
		}

		commandDetails := Parse(command)

		switch commandDetails.name {
		case "exit":
			return
		case "type":
			processTypeCommand(commandDetails.args[0])
		case "pwd":
			processPwdCommand()
		case "cd":
			processCdCommand(commandDetails.args)
		default:
			executeCommand(commandDetails)
		}

	}
}

func processTypeCommand(commandName string) {

	if ShellBuiltinCommands[commandName] {
		fmt.Printf("%s is a shell builtin\n", commandName)
		return
	}

	executablePaths, isPathEnvSet := os.LookupEnv("PATH")

	if isPathEnvSet {

		for path := range strings.SplitSeq(executablePaths, ":") {
			ok, commandFullPath := isExecutable(path, commandName)
			if ok {
				fmt.Printf("%s is %s\n", commandName, commandFullPath)
				return
			}
		}

	}

	fmt.Printf("%s: not found\n", commandName)

}

func processPwdCommand() {
	pwd, err := os.Getwd()
	if err != nil {
		fmt.Println("Unable to find the present working directory")
	}
	fmt.Println(pwd)
}

func processCdCommand(arg []string) {

	var directory string
	if len(arg) > 0 {
		directory = arg[0]
	} else {
		directory = "~"
	}

	if strings.HasPrefix(directory, "~") {
		homeDirectory := os.Getenv("HOME")
		directory = strings.ReplaceAll(directory, "~", homeDirectory)
	}

	info, err := os.Stat(directory)

	if err != nil || !info.IsDir() {
		fmt.Printf("cd: %s: No such file or directory\n", directory)
	}

	os.Chdir(directory)

}
