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

	if err := loadShellRC(".shellrc"); err != nil {
		fmt.Println("Unable to open .shellrc")
		os.Exit(1)
	}

	prompt, reader := os.Getenv("PS"), bufio.NewReader(os.Stdin)
	repl(prompt, reader)
}

func loadShellRC(path string) error {

	file, err := os.Open(path)
	if err != nil {
		fmt.Printf("Unable to open %s, please check whether %s exists or not.\n", path, path)
		return err
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {

		line := strings.TrimSpace(scanner.Text())

		if "" == line || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key := parts[0]
			value := strings.Trim(parts[1], "'\"")
			os.Setenv(key, value)
		}

	}

	return nil

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
	parts, err := Split(commandDetails)
	if err != nil {
		fmt.Println("Unable to parse input commands")
	}
	if len(parts) > 1 {
		return parts[0], parts[1:]
	}

	return parts[0], []string{}
}

func processEchoCommand(args []string) {
	fmt.Println(strings.Join(args, SPACE))
}

func processTypeCommand(commandName string) {

	if SHELL_BUILTIN_COMMANDS[commandName] {
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

func executeCommand(commandName string, args []string) {

	if SHELL_BUILTIN_COMMANDS[commandName] {
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
