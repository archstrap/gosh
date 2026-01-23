package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

func executeCommand(commandDetails Command) {

	commandName := commandDetails.name
	args := commandDetails.args
	redirections := commandDetails.redirections

	if ShellBuiltinCommands[commandName] {
		return
	}

	executablePaths, isPathEnvSet := os.LookupEnv("PATH")

	if !isPathEnvSet {
		fmt.Printf("%s: command not found\n", commandName)
		return
	}
	for path := range strings.SplitSeq(executablePaths, ":") {
		ok, commandFullPath := isExecutable(path, commandName)
		if !ok {
			continue
		}
		cmd := exec.Command(commandFullPath, args...)
		setIO(&redirections, cmd)

		cmd.Run()
		return
	}

	fmt.Printf("%s: command not found\n", commandName)

}

func setIO(ioDetails *map[int]*Redirection, cmd *exec.Cmd) {
	cmd.Stdin = openFile(GetOrDefault(*ioDetails, syscall.Stdin, NewRedirection("/dev/stdin")), syscall.Stdin, os.O_RDONLY)
	cmd.Stdout = openFile(GetOrDefault(*ioDetails, syscall.Stdout, NewRedirection("/dev/stdout")), syscall.Stdout, os.O_WRONLY)
	cmd.Stderr = openFile(GetOrDefault(*ioDetails, syscall.Stderr, NewRedirection("/dev/stderr")), syscall.Stderr, os.O_WRONLY)
}

func openFile(r *Redirection, defaultFd int, mode int) *os.File {

	if r.fileName == "/dev/stdin" || r.fileName == "/dev/stdout" || r.fileName == "/dev/stderr" {
		return os.NewFile(uintptr(defaultFd), r.fileName)
	}

	flags := mode | os.O_CREATE

	if r.appendOnly {
		flags |= os.O_APPEND
	} else if mode == os.O_WRONLY {
		flags |= os.O_TRUNC
	}

	fd, err := syscall.Open(r.fileName, flags, 0644)

	if err != nil {
		fmt.Println(err)
	}

	return os.NewFile(uintptr(fd), r.fileName)
}
