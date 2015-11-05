package main

import (
	"fmt"
	"os"

	"github.com/mitchellh/cli"
)

func main() {
	args := os.Args[1:]

	// Get the command line args. We shortcut "--version" and "-v" to
	// just show the version.
	for _, arg := range args {
		if arg == "-v" || arg == "--version" {
			newArgs := make([]string, 1)
			newArgs[0] = "version"
			args = newArgs
			break
		}
	}

	cli := &cli.CLI{
		Args:     args,
		Commands: Commands,
		HelpFunc: cli.BasicHelpFunc("mixpanel"),
	}
	exitCode, err := cli.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error executing CLI: %s\n", err.Error())
		os.Exit(1)
	}
	os.Exit(exitCode)

	// config, err := api.DefaultConfig()
	// if err != nil {
	//	log.Fatal(err)
	// }
	// q := &api.QueryOptions{
	//	Key:      config.Key,
	//	Secret:   config.Secret,
	//	Expire:   "1445934932",
	//	FromDate: "2015-01-02",
	//	ToDate:   "2015-01-02",
	//	Format:   "json",
	// }
	// client := api.NewClient(*config)
	// client.Export(q)
}
