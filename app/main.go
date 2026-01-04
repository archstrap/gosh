package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
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
		"pwd":  true,
		"cd":   true,
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
		case "pwd":
			processPwdCommand()
		case "cd":
			processCdCommand(args)
		default:
			executeCommand(commandName, args)
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

func processPwdCommand() {
	pwd, err := os.Getwd()
	if err != nil {
		fmt.Println("Unable to find the present working directory")
	}
	fmt.Println(pwd)
}

func processCdCommand(arg []string) {
	directory := arg[0]
	// handle absolute path

	info, err := os.Stat(directory)

	if err != nil || !info.IsDir() {
		fmt.Printf("cd: %s: No such file or directory\n", directory)
	}

	os.Chdir(directory)

}

func executeCommand(commandName string, args []string) {

	_, ok := SHELL_BUILTIN_COMMANDS[commandName]

	if ok {
		return
	}

	execuatblePaths, isPathEnvSet := os.LookupEnv("PATH")

	if isPathEnvSet {

		for path := range strings.SplitSeq(execuatblePaths, ":") {
			ok, commandFullPath := isExecutable(path, commandName)
			if ok {

				combinedArgs := append([]string{commandName}, args...)

				cmd := exec.Command(commandFullPath, args...)
				cmd.Args = combinedArgs
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr
				cmd.Stdin = os.Stdin

				err := cmd.Run()
				if err != nil {
					fmt.Println(err)
				}

				// output, err := cmd.Output()
				// if err != nil {
				// 	fmt.Println(err)
				// 	return
				// }
				//
				// fmt.Println(string(output))
				return
			}
		}

	}

	fmt.Printf("%s: command not found\n", commandName)

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
