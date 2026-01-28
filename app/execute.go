package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
)

func ConvertToExecutableCommand(commandDetails *Command) *exec.Cmd {
	command := commandDetails.name
	args := commandDetails.args
	redirections := commandDetails.redirections

	// fmt.Printf("Name: %s, args: %v, io-redirections: %v\n", command, args, redirections)

	executablePaths, set := os.LookupEnv("PATH")
	if !set {
		return nil
	}

	var cmd *exec.Cmd

	for path := range strings.SplitSeq(executablePaths, ":") {

		ok, path := isExecutable(path, command)

		if !ok {
			continue
		}
		cmd = exec.Command(path, args...)
		setIO(&redirections, cmd)
		cmd.Args = append([]string{command}, args...)
		break
	}

	return cmd
}

func executeCommand(commandDetails Command) {

	if cmd := ConvertToExecutableCommand(&commandDetails); cmd != nil {
		cmd.Run()
		return
	}

	fmt.Printf("%s: command not found\n", commandDetails.name)

}

func isExecutable(directoryPath string, target string) (bool, string) {

	fullPath := filepath.Join(directoryPath, target)
	stat, err := os.Stat(fullPath)

	if err != nil {
		return false, ""
	}

	if stat.Mode()&0100 != 0 {
		return true, fullPath
	}

	return false, ""
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

func SearchAllExecutable(commandPrefix string) []string {

	allPaths, ok := os.LookupEnv("PATH")
	if !ok {
		return nil
	}

	var executableCommmands []string

	for path := range strings.SplitSeq(allPaths, string(os.PathListSeparator)) {

		files, err := os.ReadDir(path)
		if err != nil {
			continue
		}

		for _, file := range files {
			if isExecutableFile(path, file.Name(), commandPrefix) {
				executableCommmands = append(executableCommmands, file.Name())
			}
		}
	}

	return executableCommmands

}

func isExecutableFile(path string, file string, commandPrefix string) bool {

	if !strings.HasPrefix(file, commandPrefix) {
		return false
	}

	fileInfo, err := os.Stat(filepath.Join(path, file))

	if err != nil {
		return false
	}

	if fileInfo.IsDir() {
		return false
	}

	return fileInfo.Mode()&0100 != 0
}

func StartCommandExecution(command string) {

	commands := Parse(strings.TrimSpace(command)) // it will be returning an []*Command

	if len(commands) == 0 {
		fmt.Printf("%s: command not found\n", command)
		return
	}

	if len(commands) == 1 {
		commandDetails := commands[0]
		switch commandDetails.name {
		case "exit":
			os.Exit(0)
			return
		case "type":
			processTypeCommand(commandDetails.args[0])
		case "pwd":
			processPwdCommand()
		case "cd":
			processCdCommand(commandDetails.args)
		default:
			executeCommand(*commandDetails)
		}
		return
	}

	var rc []*exec.Cmd

	for _, commandDetails := range commands {
		rc = append(rc, ConvertToExecutableCommand(commandDetails))
	}

	pipeCount := len(commands) - 1
	var openFiles []*os.File

	for i := range pipeCount {
		curr := rc[i]
		next := rc[i+1]

		r, w, _ := os.Pipe()

		curr.Stdout = w
		next.Stdin = r

		openFiles = append(openFiles, r, w)
	}

	Do(rc, func(command *exec.Cmd) {
		command.Start()
	})

	Do(openFiles, func(file *os.File) {
		file.Close()
	})

	Do(rc, func(command *exec.Cmd) {
		command.Wait()
	})

}

func CloseResources(openFiles []*os.File) {
	for _, openFile := range openFiles {
		if openFile != nil {
			openFile.Close()
		}
	}
}

func Do[K any](items []*K, callBack func(item *K)) {
	for _, item := range items {
		callBack(item)
	}
}
