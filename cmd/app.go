package cmd

import (
	"os"
	"strings"

	"github.com/mitchellh/cli"

	// These are lua wrapped code that will be made accessible to lua
	_ "github.com/tanema/amore/gfx/wrap"
	_ "github.com/tanema/amore/input"

	"github.com/tanema/amore/runtime"
)

var App = &cli.CLI{
	Name:                       "moony",
	Args:                       os.Args[1:],
	Commands:                   commands,
	Autocomplete:               true,
	AutocompleteNoDefaultFlags: true,
}

var ui = &cli.BasicUi{
	Reader:      os.Stdin,
	Writer:      os.Stdout,
	ErrorWriter: os.Stderr,
}

var commands = map[string]cli.CommandFactory{
	"":    func() (cli.Command, error) { return &runCommand{ui: ui}, nil },
	"run": func() (cli.Command, error) { return &runCommand{ui: ui}, nil },
}

type runCommand struct {
	ui cli.Ui
}

func (run *runCommand) Run(args []string) int {
	if err := runtime.Run("main.lua"); err != nil {
		run.ui.Error(err.Error())
		return 1
	}
	return 0
}

func (run *runCommand) Synopsis() string {
	return ""
}

func (run *runCommand) Help() string {
	helpText := `
Usage: moony
Run your program yo
Options:
  -h, --help  show this help
`
	return strings.TrimSpace(helpText)
}
