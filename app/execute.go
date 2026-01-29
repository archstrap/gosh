package main

import (
	"fmt"
	"os"
	"strings"
)

func ExecuteCommand(input string) {
	commands := Parse(strings.TrimSpace(input))
	if len(commands) == 0 {
		fmt.Printf("%s: not found\n", input)
		return
	}

	var executables []Executable
	for _, command := range commands {
		executable, err := CreateExecutable(command)
		if err != nil {
			fmt.Println(err)
			return
		}
		executables = append(executables, executable)
	}

	pipeCount := len(executables) - 1
	var openFiles []*os.File

	for i := range pipeCount {
		r, w, _ := os.Pipe()
		executables[i].SetStdout(w)
		executables[i+1].SetStdin(r)
		openFiles = append(openFiles, r, w)
	}

	PerformTask(executables, func(e Executable) {
		e.Start()
	})

	Do(openFiles, func(file *os.File) {
		file.Close()
	})

	PerformTask(executables, func(e Executable) {
		e.Wait()
	})
}
