/*
TODO:
	> Split up in separate files
		> command.go
		> application.go
		> descriptor.go
	> Document symbols
	> Write unit tests

Final notes:
	This package DOES NOT hold a stable public API and is meant for internal purposes only.
*/
package command

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"os"
	"sort"
)

type IO interface {
	ReadLn() string
	Write(message string)
	WriteLn(message string)
	WriteError(message string)
	WriteErrorLn(message string)
}

type Describable interface {
	Name() string
	Description() string
}

type Describers []Describable

func (d Describers) Len() int {
	return len(d)
}

func (d Describers) Less(i, j int) bool {
	return d[1].Name() < d[j].Name()
}

func (d Describers) Swap(i, j int) {
	d[i], d[j] = d[j], d[i]
}

type descriptor struct {
	Describable
	name        string
	description string
}

func (d *descriptor) Name() string {
	return d.name
}

func (d *descriptor) SetName(name string) {
	d.name = name
}

func (d *descriptor) Description() string {
	return d.description
}

func (d *descriptor) SetDescription(description string) {
	d.description = description
}

type Command interface {
	Describable
	Init(args []string) ([]string, error)
	Run(args []string, io IO) int
}

type Commands map[string]*Command

type Runner interface {
	Run(command string, args []string) int
	Register(cmd Command)
	Usage()
}

type Application struct {
	descriptor
	IO
	Runner
	input    *bufio.Reader
	output   *os.File
	error    *os.File
	commands Commands
	context  context.Context
}

func (app Application) Usage() {
	usage := fmt.Sprintf("%s\n\n", app.Description())
	usage += fmt.Sprintf("Usage:\n\n\t%s <command> [arguments]\n\n", app.Name())
	usage += "The commands are:\n\n"
	commands := make(Describers, 0)

	for _, cmd := range app.commands {
		commands = append(commands, *cmd)
	}

	sort.Sort(commands)

	for _, command := range commands {
		usage += fmt.Sprintf(
			"\t%-10s %s\n",
			command.Name(),
			command.Description(),
		)
	}

	app.Write(usage)
}

func (app Application) Run(command string, args []string) int {
	select {
	case <-app.context.Done():
		app.WriteErrorLn("Aborted.")
		return 3
	default:
		cmd, ok := app.commands[command]

		if ok == false {
			app.Usage()
			return 1
		}

		args, e := (*cmd).Init(args)

		if e != nil {
			app.WriteErrorLn(e.Error())
			return 2
		}

		return (*cmd).Run(args[0:], app)
	}
}

func (app Application) ReadLn() string {
	var (
		result []byte
		e      error
	)

	if result, _, e = app.input.ReadLine(); e != nil {
		panic(e)
	}

	return string(result)
}

func (app Application) Write(message string) {
	if _, e := app.output.WriteString(message); e != nil {
		panic(e)
	}
}

func (app Application) WriteLn(message string) {
	app.Write(message + "\n")
}

func (app Application) WriteError(message string) {
	if _, e := app.error.WriteString(message); e != nil {
		panic(e)
	}
}

func (app Application) WriteErrorLn(message string) {
	app.WriteError(message + "\n")
}

func (app Application) Register(cmd Command) {
	app.commands[cmd.Name()] = &cmd
}

func ExtractCommandArgs(input []string) (cmd string, args []string) {
	if len(input) > 0 {
		cmd = input[0]
	}

	if len(input) > 1 {
		args = input[1:]
	}

	return cmd, args
}

func NewApplication(ctx context.Context, input *os.File, output *os.File, error *os.File) *Application {
	return &Application{
		descriptor: descriptor{
			name:        os.Args[0],
			description: "",
		},
		context:  ctx,
		input:    bufio.NewReader(input),
		output:   output,
		error:    error,
		commands: make(Commands),
	}
}

type BaseCommand struct {
	Command
	flags       *flag.FlagSet
	description string
}

func (cmd BaseCommand) Name() string {
	return cmd.flags.Name()
}

func (cmd BaseCommand) Description() string {
	return cmd.description
}

func (cmd BaseCommand) Init(args []string) ([]string, error) {
	e := cmd.flags.Parse(args)
	return cmd.flags.Args(), e
}
