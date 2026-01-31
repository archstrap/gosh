package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type Executable interface {
	SetStdin(in io.Reader)
	SetStdout(out io.Writer)
	SetStderr(err io.Writer)
	GetStdin() io.Reader
	GetStdout() io.Writer
	GetStderr() io.Writer
	GetCommandType() string
	Start() error
	Wait() error
}

type ResourceManager struct {
	readers []io.Reader
	writers []io.Writer
}

func (r *ResourceManager) CloseResources() {
	PerformTask(r.readers, func(resource io.Reader) {
		if file, ok := resource.(*os.File); ok && isStandardIoFile(file.Name()) {
			return
		}
		if closer, ok := resource.(io.Closer); ok {
			closer.Close()
		}
	})
	PerformTask(r.writers, func(resource io.Writer) {
		if file, ok := resource.(*os.File); ok && isStandardIoFile(file.Name()) {
			return
		}
		if closer, ok := resource.(io.Closer); ok {
			closer.Close()
		}
	})

}
func (r *ResourceManager) AddReader(reader io.Reader) {
	r.readers = append(r.readers, reader)
}
func (r *ResourceManager) AddWriter(writer io.Writer) {
	r.writers = append(r.writers, writer)
}

func (r *ResourceManager) AddResource(executable Executable) {
	r.AddReader(executable.GetStdin())
	r.AddWriter(executable.GetStdout())
	r.AddWriter(executable.GetStderr())
}

func CreateExecutable(command *Command, r *ResourceManager) (Executable, error) {

	if ShellBuiltinCommands[command.name] {
		builtinCommand := NewBuiltinCommand(command.name, command.args...)
		SetIO(&command.redirections, builtinCommand)
		return builtinCommand, nil
	}

	if ok, path := isExternal(command.name); ok {
		externalCommand := NewExternalCommand(path, command.args...)
		externalCommand.cmd.Args = append([]string{command.name}, command.args...)
		SetIO(&command.redirections, externalCommand)
		r.AddResource(externalCommand)
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
func (e *ExternalCommand) GetStdin() io.Reader {
	return e.cmd.Stdin
}
func (e *ExternalCommand) GetStdout() io.Writer {
	return e.cmd.Stdout
}
func (e *ExternalCommand) GetStderr() io.Writer {
	return e.cmd.Stderr
}

func (e *ExternalCommand) GetCommandType() string {
	return "EXTERNAL"
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

func (b *BuiltinCommand) GetStdin() io.Reader {
	return b.Stdin
}
func (b *BuiltinCommand) GetStdout() io.Writer {
	return b.Stdout
}
func (b *BuiltinCommand) GetStderr() io.Writer {
	return b.Stderr
}

func (b *BuiltinCommand) GetCommandType() string {
	return "BUILTIN"
}

func (b *BuiltinCommand) Start() error {
	go func() {

		defer func() {

			// ignore if standard I/O
			if file, ok := b.Stdin.(*os.File); ok && !isStandardIoFile(file.Name()) {
				file.Close()
			}

			if file, ok := b.Stdout.(*os.File); ok && !isStandardIoFile(file.Name()) {
				file.Close()
			}

			if file, ok := b.Stderr.(*os.File); ok && !isStandardIoFile(file.Name()) {
				file.Close()
			}

		}()

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
	"pwd":     pwdBuiltin,
	"echo":    echoBuiltin,
	"type":    typeBuiltin,
	"exit":    exitBuiltin,
	"cd":      cdBuiltin,
	"history": historyBuiltin,
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
	if path, ok := os.LookupEnv("HISTFILE"); ok {
		history.AppendHistory(path)
	}
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

func historyBuiltin(args []string, stdin io.Reader, stdout io.Writer, stderr io.Writer) error {

	switch len(args) {
	case 0:
		data := history.GetLast(len(history.commands))
		fmt.Fprint(stdout, data)
	case 1:
		arg := args[0]
		count, _ := strconv.Atoi(arg)
		data := history.GetLast(count)
		fmt.Fprint(stdout, data)
	case 2:
		option, file := args[0], args[1]
		switch option {
		case `-r`:
			history.LoadHistory(file)
		case `-w`:
			history.WriteHistory(file)
		case `-a`:
			history.AppendHistory(file)
		}

	}

	return nil
}
