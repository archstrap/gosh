package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"syscall"
)

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

func SetIO(ioDetails *map[int]*Redirection, cmd Executable) {
	cmd.SetStdin(openFile(GetOrDefault(*ioDetails, syscall.Stdin, NewRedirection("/dev/stdin")), syscall.Stdin, os.O_RDONLY))
	cmd.SetStdout(openFile(GetOrDefault(*ioDetails, syscall.Stdout, NewRedirection("/dev/stdout")), syscall.Stdout, os.O_WRONLY))
	cmd.SetStderr(openFile(GetOrDefault(*ioDetails, syscall.Stderr, NewRedirection("/dev/stderr")), syscall.Stderr, os.O_WRONLY))
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

func isExternal(target string) (bool, string) {

	allPaths, ok := os.LookupEnv("PATH")
	if !ok {
		return false, ""
	}

	for path := range strings.SplitSeq(allPaths, string(os.PathListSeparator)) {
		if ok, fullPath := isExecutable(path, target); ok {
			return true, fullPath
		}
	}

	return false, ""

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
