package main

import (
	"os"

	"github.com/cskksc/mixpanel/command"
	"github.com/mitchellh/cli"
)

// Commands is the mapping of all the available mixpanel commands.
var Commands map[string]cli.CommandFactory

func init() {
	ui := &cli.BasicUi{Writer: os.Stdout}
	Commands = map[string]cli.CommandFactory{
		"version": func() (cli.Command, error) {
			return &command.VersionCommand{
				Version: Version,
				Ui:      ui,
			}, nil
		},
		"export": func() (cli.Command, error) {
			return &command.ExportCommand{
				Ui: ui,
			}, nil
		},
	}
}
