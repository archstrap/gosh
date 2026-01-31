package main

import (
	"fmt"
	"os"
	"strings"
)

func ExecuteCommand(input string) {
	commandInput := os.ExpandEnv(strings.TrimSpace(input))
	commands := Parse(commandInput)
	if len(commands) == 0 {
		fmt.Printf("%s: not found\n", input)
		return
	}

	var executables []Executable
	resourceManager := &ResourceManager{}

	defer func() {
		resourceManager.CloseResources()
	}()

	for _, command := range commands {
		executable, err := CreateExecutable(command, resourceManager)
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
		// Current Command Will Write at the Pipe's Write end
		executables[i].SetStdout(w)
		// Next Command Will Read from the Pipe's Read end
		executables[i+1].SetStdin(r)

		if executables[i].GetCommandType() == "EXTERNAL" {
			openFiles = append(openFiles, w)
		}

		if executables[i+1].GetCommandType() == "EXTERNAL" {
			openFiles = append(openFiles, r)
		}

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
