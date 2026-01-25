package main

import (
	"fmt"
	"maps"
	"os"
	"slices"
	"strings"

	"golang.org/x/term"
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

	terminalFd := int(os.Stdin.Fd())
	oldState, err := term.MakeRaw(terminalFd)
	if err != nil {
		panic(err)
	}

	defer term.Restore(terminalFd, oldState)

	prompt := os.Getenv("PS")
	repl(prompt, terminalFd, oldState)
}

func repl(prompt string, terminalFd int, oldState *term.State) {
	trie := NewTrie()
	keys := slices.Collect(maps.Keys(ShellBuiltinCommands))
	trie.InsertAll(keys...)
	var command strings.Builder

	fmt.Print(prompt) // Print prompt once at start

	for {
		buf := make([]byte, 1)
		if _, err := os.Stdin.Read(buf); err != nil {
			fmt.Println(err)
			continue
		}

		switch buf[0] {
		// handling Ctrl + c
		case 3:
			fmt.Print("\r\n")
			command.Reset()
			fmt.Print(prompt)
		// handling Ctrl + d
		case 4:
			fmt.Print("\r\n")
			return
		// Handling tab
		case '\t':
			fmt.Print("\r\n")
			suggestions := trie.SearchAll(command.String())
			if len(suggestions) == 1 {
				command.Reset()
				command.WriteString(fmt.Sprintf("%s ", suggestions[0]))
				fmt.Printf("%s%s ", prompt, suggestions[0])
			}
		// Handling Enter
		case '\n', '\r':
			fmt.Print("\r\n")
			if command.Len() > 0 {
				if err := term.Restore(terminalFd, oldState); err != nil {
					return
				}
				StartCommandExecution(command.String())
				command.Reset()
				if _, err := term.MakeRaw(terminalFd); err != nil {
					return
				}
			}
			fmt.Print(prompt)
		// Handling BackSpace and Del
		case 127, 8:

		default:
			command.WriteByte(buf[0])
			fmt.Print(string(buf[0]))
		}
	}
}
func StartCommandExecution(command string) {

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
