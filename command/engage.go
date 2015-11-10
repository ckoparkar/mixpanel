package command

import (
	"flag"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/cskksc/mixpanel/api"
	"github.com/mitchellh/cli"
)

type EngageCommand struct {
	Ui   cli.Ui
	args []string
	out  string
}

func (c *EngageCommand) Help() string {
	helpText := `Usage: mixpanel engage [options]

  Exports data from mixpanel people analytics.

Options:

  -from=yesterday Start date to extract events.
  -out=STDOUT     Decides where to write the data.
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
	var w io.Writer
	if c.out != "" {
		f, err := os.OpenFile(c.out, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0660)
		if err != nil {
			f, _ = ioutil.TempFile(".", "")
			log.Printf("[ERR] Couldnt open file. Encoding to %s.", err, f.Name())
		}
		defer f.Close()
		w = f
	} else {
		w = os.Stdout
	}

	if err := client.Engage(queryOptions, w); err != nil {
		return 1
	}
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
	var out string
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
