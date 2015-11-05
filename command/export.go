package command

import (
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/cskksc/mixpanel/api"
	"github.com/mitchellh/cli"
)

type ExportCommand struct {
	Ui   cli.Ui
	args []string
	out  string
}

func (c *ExportCommand) Help() string {
	helpText := `Usage: mixpanel export [options]

  Exports raw dump of mixpanel data for a set of events over a time period.

Options:

  -from=yesterday Start date to extract events.
  -to=yesterday   End date to extract events.
  -format=json    Choose export format between json/csv.
  -event=E        Extract data for only event E.
  -out=STDOUT     Decides where to write the data.
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
	data, err := client.Export(queryOptions)
	if err != nil {
		log.Printf("[ERR] %s", err)
		return 1
	}
	if c.out != "" {
		api.OverwriteFile(c.out, data)
	} else {
		fmt.Println(string(data))
	}
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
	var out string
	cmdFlags.StringVar(&cmdQueryOptions.FromDate, "from", "", "from date")
	cmdFlags.StringVar(&cmdQueryOptions.ToDate, "to", "", "to date")
	cmdFlags.StringVar(&cmdQueryOptions.Format, "format", "", "data format")
	cmdFlags.StringVar(&cmdQueryOptions.Event, "event", "", "event name")
	cmdFlags.StringVar(&out, "out", "", "output destination")
	if err := cmdFlags.Parse(c.args); err != nil {
		return nil, err
	}
	if out != "" {
		c.out = out
	}

	queryOptions := api.DefaultQueryOptions(config)
	// Not all config would be provided as cmd-line args
	queryOptions = api.MergeQueryOptions(queryOptions, &cmdQueryOptions)
	return queryOptions, nil
}
