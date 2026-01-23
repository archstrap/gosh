package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var (
	ShellBuiltinCommands = map[string]bool{
		"type": true,
		"exit": true,
		"pwd":  true,
		"cd":   true,
	}
)

func main() {

	if err := loadShellRC(".shellrc"); err != nil {
		fmt.Println("Unable to open .shellrc")
		os.Exit(1)
	}

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
