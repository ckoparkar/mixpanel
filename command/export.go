package command

import (
	"flag"
	"log"
	"strings"

	"github.com/cskksc/mixpanel/api"
	"github.com/mitchellh/cli"
)

type ExportCommand struct {
	Ui   cli.Ui
	args []string
}

func (c *ExportCommand) Help() string {
	helpText := `Usage: mixpanel export [options]

  Exports mixpanel data.

Options:

  -from=yesterday Start date to extract events.
  -to=yesterday   End date to extract events.
  -format=json    Choose export format between json/csv.
  -event=E        Extract data for only event E.
`
	return strings.TrimSpace(helpText)
}

func (c *ExportCommand) Run(args []string) int {
	c.args = args
	config, err := api.DefaultConfig()
	if err != nil {
		log.Fatal("[ERR] " + err.Error())
	}
	queryOptions, err := c.readQueryOptions(config)
	if err != nil {
		return 1
	}
	client := api.NewClient(*config)
	client.Export(queryOptions)
	return 0
}

func (c *ExportCommand) Synopsis() string {
	return "Exports mixpanel data."
}

// readConfig reads config provided as cmd-line args,
// and merges it with the defaults
func (c *ExportCommand) readQueryOptions(config *api.Config) (*api.QueryOptions, error) {
	cmdFlags := flag.NewFlagSet("export", flag.ContinueOnError)
	cmdFlags.Usage = func() { c.Ui.Output(c.Help()) }

	var cmdQueryOptions api.QueryOptions
	cmdFlags.StringVar(&cmdQueryOptions.FromDate, "from", "", "from date")
	cmdFlags.StringVar(&cmdQueryOptions.ToDate, "to", "", "to date")
	cmdFlags.StringVar(&cmdQueryOptions.Format, "format", "", "data format")
	cmdFlags.StringVar(&cmdQueryOptions.Event, "event", "", "event name")
	if err := cmdFlags.Parse(c.args); err != nil {
		return nil, err
	}

	queryOptions := api.DefaultQueryOptions(config)
	// Not all config would be provided as cmd-line args
	queryOptions = api.MergeQueryOptions(queryOptions, &cmdQueryOptions)
	return queryOptions, nil
}