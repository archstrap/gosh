package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

type Executable interface {
	SetStdin(in io.Reader)
	SetStdout(out io.Writer)
	SetStderr(err io.Writer)
	Start() error
	Wait() error
}

func CreateExecutable(command *Command) (Executable, error) {

	if ShellBuiltinCommands[command.name] {
		builtinCommand := NewBuiltinCommand(command.name, command.args...)
		SetIO(&command.redirections, builtinCommand)
		return builtinCommand, nil
	}

	if ok, path := isExternal(command.name); ok {
		externalCommand := NewExternalCommand(path, command.args...)
		externalCommand.cmd.Args = append([]string{command.name}, command.args...)
		SetIO(&command.redirections, externalCommand)
		return externalCommand, nil
	}

	return nil, fmt.Errorf("%s: not found", command.name)
}

// ExternalCommand
type ExternalCommand struct {
	cmd *exec.Cmd
}

func NewExternalCommand(name string, args ...string) *ExternalCommand {
	return &ExternalCommand{
		cmd: exec.Command(name, args...),
	}
}

func (e *ExternalCommand) SetStdin(in io.Reader) {
	e.cmd.Stdin = in
}
func (e *ExternalCommand) SetStdout(out io.Writer) {
	e.cmd.Stdout = out
}
func (e *ExternalCommand) SetStderr(err io.Writer) {
	e.cmd.Stderr = err
}
func (e *ExternalCommand) Start() error {
	return e.cmd.Start()
}
func (e *ExternalCommand) Wait() error {
	return e.cmd.Wait()
}

// BuiltinCommand

type BuiltinCommand struct {
	Name   string
	Args   []string
	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer
	Done   chan error
}

func NewBuiltinCommand(name string, args ...string) *BuiltinCommand {
	return &BuiltinCommand{
		Name: name,
		Args: args,
		Done: make(chan error),
	}
}

func (b *BuiltinCommand) SetStdin(in io.Reader) {
	b.Stdin = in
}
func (b *BuiltinCommand) SetStdout(out io.Writer) {
	b.Stdout = out
}
func (b *BuiltinCommand) SetStderr(err io.Writer) {
	b.Stderr = err
}
func (b *BuiltinCommand) Start() error {
	go func() {
		executeFn := BuiltinRegistry[b.Name]
		b.Done <- executeFn(b.Args, b.Stdin, b.Stdout, b.Stderr)
	}()
	return nil
}
func (b *BuiltinCommand) Wait() error {
	return <-b.Done
}

// BuiltinRegistry

type BuiltinFunc func(args []string, stdin io.Reader, stdout io.Writer, stderr io.Writer) error

var BuiltinRegistry = map[string]BuiltinFunc{
	"pwd":  pwdBuiltin,
	"echo": echoBuiltin,
	"type": typeBuiltin,
	"exit": exitBuiltin,
	"cd":   cdBuiltin,
}

// pwd pwdBuiltin
func pwdBuiltin(args []string, stdin io.Reader, stdout io.Writer, stderr io.Writer) error {
	pwd, err := os.Getwd()
	if err != nil {
		fmt.Fprintln(stderr, "pwd: unable to get current working directory")
		return err
	}

	fmt.Fprintln(stdout, pwd)

	return nil
}

// echo builtin
func echoBuiltin(args []string, stdin io.Reader, stdout io.Writer, stderr io.Writer) error {
	data := strings.Join(args, " ")
	fmt.Fprintln(stdout, data)
	return nil
}

// type builtin
func typeBuiltin(args []string, stdin io.Reader, stdout io.Writer, stderr io.Writer) error {

	for _, command := range args {
		commandName := strings.TrimSpace(command)
		if ShellBuiltinCommands[commandName] {
			fmt.Fprintf(stdout, "%s is a shell builtin\n", commandName)
		} else if ok, path := isExternal(commandName); ok {
			fmt.Fprintf(stdout, "%s is %s\n", commandName, path)
		} else {
			fmt.Fprintf(stderr, "%s: not found\n", commandName)
		}
	}
	return nil
}

// exit builtin
func exitBuiltin(args []string, stdin io.Reader, stdout io.Writer, stderr io.Writer) error {
	os.Exit(0)
	return nil
}

// cdBuiltin
func cdBuiltin(args []string, stdin io.Reader, stdout io.Writer, stderr io.Writer) error {

	var directory string
	if len(args) > 0 {
		directory = args[0]
	} else {
		directory = "$HOME"
	}

	if strings.HasPrefix(directory, "~") {
		directory = strings.ReplaceAll(directory, "~", "$HOME")
	}

	directory = os.ExpandEnv(directory)

	info, err := os.Stat(directory)
	if err != nil || !info.IsDir() {
		errorMessage := fmt.Sprintf("cd: %s: No such file or directory", directory)
		fmt.Fprintln(stderr, errorMessage)
		return fmt.Errorf(errorMessage)
	}
	return os.Chdir(directory)
}
