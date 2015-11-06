package api

import (
	"fmt"
	"io/ioutil"
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
func (c *Client) Export(q *QueryOptions) ([]byte, error) {
	r := c.newRequest("GET", "/export/")
	r.setQueryOptions(q)

	_, resp, err := c.doRequest(r)
	if err != nil {
		return nil, err
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func (c *Client) Engage(q *QueryOptions) error {
	c.config.Address = "mixpanel.com"
	r := c.newRequest("GET", "/engage")
	r.setQueryOptions(q)

	_, resp, err := c.doRequest(r)
	if err != nil {
		return err
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	fmt.Println(string(b))
	return nil
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
