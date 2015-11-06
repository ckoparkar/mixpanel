package command

import (
	"flag"
	"log"
	"strings"

	"github.com/cskksc/mixpanel/api"
	"github.com/mitchellh/cli"
)

type EngageCommand struct {
	Ui   cli.Ui
	args []string
}

func (c *EngageCommand) Help() string {
	helpText := `Usage: mixpanel engage [options]

  Exports data from mixpanel people analytics.

Options:

  -from=yesterday Start date to extract events.
`
	return strings.TrimSpace(helpText)
}

func (c *EngageCommand) Run(args []string) int {
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
	client.Engage(queryOptions)
	return 0
}

func (c *EngageCommand) Synopsis() string {
	return "Exports mixpanel data."
}

// readConfig reads config provided as cmd-line args,
// and merges it with the defaults
func (c *EngageCommand) readQueryOptions(config *api.Config) (*api.QueryOptions, error) {
	cmdFlags := flag.NewFlagSet("engage", flag.ContinueOnError)
	cmdFlags.Usage = func() { c.Ui.Output(c.Help()) }

	var cmdQueryOptions api.QueryOptions
	//cmdFlags.StringVar(&cmdQueryOptions.FromDate, "from", "", "from date")
	if err := cmdFlags.Parse(c.args); err != nil {
		return nil, err
	}

	queryOptions := api.DefaultQueryOptions(config)
	// Not all config would be provided as cmd-line args
	queryOptions = api.MergeQueryOptions(queryOptions, &cmdQueryOptions)
	return queryOptions, nil
}
