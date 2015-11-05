package api

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"
)

// Client provides a client to the Mixpanel API
type Client struct {
	config Config
}

// NewClient returns a new client
func NewClient(config Config) *Client {
	return &Client{
		config: config,
	}
}

// Export implements the mixpanel export API, prints data to STDOUT.
func (c *Client) Export(q *QueryOptions) {
	r := c.newRequest("GET", "/export/")
	r.setQueryOptions(q)

	_, resp, err := c.doRequest(r)
	if err != nil {
		log.Fatal(err)
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(b))
}

// newRequest is used to create a new request
func (c *Client) newRequest(method, path string) *request {
	r := &request{
		config: &c.config,
		method: method,
		url: &url.URL{
			Scheme: c.config.Scheme,
			Host:   c.config.Address,
			Path:   "/api/2.0" + path,
		},
		params: make(map[string][]string),
	}
	return r
}

// doRequest runs a request with our client
func (c *Client) doRequest(r *request) (time.Duration, *http.Response, error) {
	req, err := r.toHTTP()
	if err != nil {
		return 0, nil, err
	}
	quit := make(chan int, 0)
	var diff time.Duration
	go func() {
		ticker := time.NewTicker(time.Second)
		for {
			select {
			case <-quit:
				fmt.Println(" Done in: ", diff)
				return
			case <-ticker.C:
				fmt.Print(".")
			}
		}
	}()
	start := time.Now()
	resp, err := c.config.HttpClient.Do(req)
	diff = time.Now().Sub(start)
	quit <- 1
	return diff, resp, err
}
