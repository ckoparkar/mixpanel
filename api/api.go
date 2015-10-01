package api

import (
	"bytes"
	"crypto/md5"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"sort"
	"time"
)

type Config struct {
	// Scheme is the URI scheme for the Consul server.
	Scheme string

	// Address is URI of mixpanel server.
	Address string

	// Key is the API key for mixpanel access.
	Key string

	// Secret is the API secret for mixpanel access.
	Secret string

	// HttpClient is the client to use. Default will be used if not provided.
	HttpClient *http.Client
}

func DefaultConfig() (*Config, error) {
	// Error out if API credentials not found
	key := os.Getenv("MIXPANEL_API_KEY")
	secret := os.Getenv("MIXPANEL_SECRET")
	if key == "" || secret == "" {
		return nil, errors.New("Mixpanel API credentials not found.")
	}

	return &Config{
		Scheme:     "http",
		Address:    "data.mixpanel.com",
		Key:        key,
		Secret:     secret,
		HttpClient: http.DefaultClient,
	}, nil
}

// Describe the various query options required for different requests
// the `json` struct tag value, decides name used in HTTP request
type QueryOptions struct {
	// Key is the API key for mixpanel access.
	Key string `json:"api_key"`

	// Secret is the API key for mixpanel access.
	Secret string `json:"api_secret"`

	// Expire is the UTC time in seconds;
	// used to expire an API request.
	Expire string `json:"expire"`

	// Signature for the method call, more documentation here
	// https://mixpanel.com/docs/api-documentation/data-export-api#auth-implementation .
	Sig string `json:"sig"`

	// FromDate sets the start date for export API.
	FromDate string `json:"from_date"`

	// ToDate sets the start date for export API.
	ToDate string `json:"to_date"`

	// Format describes the response format.
	Format string `json:"format"`

	// Event only exports data for this event.
	Event string `json:"event"`
}

/*
 Adds sig to QueryOptions
 args = all query parameters going to be sent out with the request
 (e.g. api_key, unit, interval, expire, format, etc.) excluding sig.

 args_sorted = sort_args_alphabetically_by_key(args)

 args_concat = join(args_sorted)

 Output: api_key=ed0b8ff6cc3fbb37a521b40019915f18event=["pages"]
	  expire=1248499222format=jsoninterval=24unit=hour

 sig = md5(args_concat + api_secret)
*/
func (q *QueryOptions) AddSig() {
	keys, params := q.toMap()
	sort.Strings(keys)

	var buf bytes.Buffer
	for _, k := range keys {
		buf.WriteString(fmt.Sprintf("%s=%s", k, params[k]))
	}
	buf.WriteString(q.Secret)
	sig := md5.Sum(buf.Bytes())
	q.Sig = fmt.Sprintf("%x", sig)
}

// Returns map of structTagName -> value
func (q *QueryOptions) toMap() ([]string, map[string]string) {
	params := make(map[string]string)
	keys := make([]string, 0)
	vt := reflect.ValueOf(q).Elem()
	for i := 0; i < vt.NumField(); i++ {
		key := vt.Type().Field(i).Tag.Get("json")
		val := vt.Field(i).Interface()
		if val != "" && key != "api_secret" {
			params[key] = val.(string)
			keys = append(keys, key)
		}
	}
	return keys, params
}

type Client struct {
	config Config
}

func NewClient(config Config) *Client {
	return &Client{
		config: config,
	}
}

// mixpanel export API, prints data to STDOUT.
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
	start := time.Now()
	resp, err := c.config.HttpClient.Do(req)
	diff := time.Now().Sub(start)
	return diff, resp, err
}

// request is used to help build up a request
type request struct {
	config *Config
	method string
	url    *url.URL
	params url.Values
	body   io.Reader
}

// toHTTP converts the request to an HTTP request
func (r *request) toHTTP() (*http.Request, error) {
	// Encode the query parameters
	r.url.RawQuery = r.params.Encode()

	// Create the HTTP request
	req, err := http.NewRequest(r.method, r.url.RequestURI(), r.body)
	if err != nil {
		return nil, err
	}

	req.URL.Host = r.url.Host
	req.URL.Scheme = r.url.Scheme
	req.Host = r.url.Host

	return req, nil
}

// setQueryOptions is used to annotate the request with
// additional query options
func (r *request) setQueryOptions(q *QueryOptions) {
	if q == nil {
		return
	}
	q.AddSig()
	_, params := q.toMap()
	for k, v := range params {
		r.params.Set(k, v)
	}
}
