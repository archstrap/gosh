package main

import (
	"bufio"
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
	bell string = "\x07"
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

	reader := bufio.NewReader(os.Stdin)
	var previousTypedCharacter byte

	for {
		currentTypedCharacter, _ := reader.ReadByte()

		switch currentTypedCharacter {
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

			commandPrefix := command.String()
			builtins := trie.SearchAll(commandPrefix)
			executables := SearchAllExecutable(commandPrefix)

			var combined []string
			AddItems(&combined, &builtins)
			AddItems(&combined, &executables)

			switch len(combined) {
			case 0:
				fmt.Print(bell)
			case 1:
				fmt.Print("\r")
				command.Reset()
				command.WriteString(fmt.Sprintf("%s ", combined[0]))
				fmt.Printf("%s%s", prompt, command.String())
			default:

				slices.Sort(combined)
				t := NewTrie()
				t.InsertAll(combined...)

				if t.LongestCommonPrefix() == command.String() {
					if previousTypedCharacter != '\t' {
						fmt.Print(bell)
					} else {
						fmt.Print("\r")
						fmt.Printf("%s%s", prompt, command.String())
						fmt.Printf("\r\n%s", strings.Join(combined, "  "))
						fmt.Printf("\r\n%s%s", prompt, t.LongestCommonPrefix())
					}
				} else {
					fmt.Printf("\r%s%s", prompt, t.LongestCommonPrefix())
					command.Reset()
					command.WriteString(t.LongestCommonPrefix())
					previousTypedCharacter = '\n'
					continue
				}

			}
		// Handling Enter
		case '\n', '\r':
			fmt.Print("\r\n")
			if command.Len() > 0 {
				// We want commands to run in cooked mode ( Normal ) for proper output formatting
				if err := term.Restore(terminalFd, oldState); err != nil {
					return
				}
				StartCommandExecution(command.String())
				command.Reset()
				// Again making it RAW mode for the next input handling
				if _, err := term.MakeRaw(terminalFd); err != nil {
					return
				}
			}
			fmt.Print(prompt)
		// Handling BackSpace and Del
		case 127, 8:
			if command.Len() > 0 {
				s := command.String()
				command.Reset()
				command.WriteString(s[:len(s)-1])
				fmt.Print("\033[D\033[K")
			}

		default:
			command.WriteByte(currentTypedCharacter)
			fmt.Print(string(currentTypedCharacter))
		}
		previousTypedCharacter = currentTypedCharacter
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

func AddItems(dest *[]string, src *[]string) {
	for _, item := range *src {
		if !slices.Contains(*dest, item) {
			*dest = append(*dest, item)
		}
	}
}
