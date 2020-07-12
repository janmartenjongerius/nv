package command

import (
	"bufio"
	"bytes"
	"context"
	"flag"
	"fmt"
	"io"
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

type Command interface {
	Name() string
	Init(args []string) ([]string, error)
	Run(args []string, io IO) int
	Usage(buffer io.Writer)
}

type Commands map[string]*Command

type Runner interface {
	Run(command string, args []string) int
	Register(cmd Command)
	Usage(collapse bool)
}

type Application struct {
	input    *bufio.Reader
	output   *os.File
	error    *os.File
	commands Commands
	context  context.Context
}

func (app Application) Usage(collapse bool) {
	usage := ""
	usages := make(map[string]*bytes.Buffer, len(app.commands))
	names := make([]string, len(app.commands))

	for name, cmd := range app.commands {
		names = append(names, name)

		if collapse {
			buffer := bytes.NewBufferString("")
			(*cmd).Usage(buffer)
			usages[name] = buffer
		}
	}

	sort.Strings(names)

	for _, name := range names {
		usage += fmt.Sprintf("%s\n", name)

		if buffer, ok := usages[name]; ok == true {
			usage += fmt.Sprintf("%s\n", buffer.String())
		}
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
			app.Usage(false)
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

func NewApplication(ctx context.Context, input *os.File, output *os.File, error *os.File) Runner {
	return Application{
		context:  ctx,
		input:    bufio.NewReader(input),
		output:   output,
		error:    error,
		commands: make(Commands),
	}
}

type BaseCommand struct {
	flags *flag.FlagSet
}

func (cmd BaseCommand) Usage(buffer io.Writer) {
	writer := cmd.flags.Output()

	cmd.flags.SetOutput(buffer)
	cmd.flags.Usage()
	cmd.flags.SetOutput(writer)
}

func (cmd BaseCommand) Name() string {
	return cmd.flags.Name()
}

func (cmd BaseCommand) Init(args []string) ([]string, error) {
	e := cmd.flags.Parse(args)
	return cmd.flags.Args(), e
}
