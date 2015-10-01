package main

import (
	"log"

	"github.com/cskksc/mixpanel/api"
)

func main() {
	config, err := api.DefaultConfig()
	if err != nil {
		log.Fatal(err)
	}
	q := &api.QueryOptions{
		Key:      config.Key,
		Secret:   config.Secret,
		Expire:   "1445934932",
		FromDate: "2015-01-02",
		ToDate:   "2015-01-02",
		Format:   "json",
	}
	client := api.NewClient(*config)
	client.Export(q)
}
