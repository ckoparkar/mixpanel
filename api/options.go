// Package api implements the Mixpanel API
package api

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"time"
)

// QueryOptions describe the various options required for different requests
// the `json` struct tag value, decides name used in HTTP request
type QueryOptions struct {
	// Key is the API key for mixpanel access.
	Key string `json:"api_key"`

	// Secret is the API key for mixpanel access.
	Secret string `json:"api_secret"`

	// Expire is the UTC time in seconds;
	// used to expire an API request.
	Expire string `json:"expire"`

	// Sig is the auth signature for the method call, more documentation here
	// https://mixpanel.com/docs/api-documentation/data-export-api#auth-implementation .
	Sig string `json:"sig"`

	// FromDate sets the start date for export API. json/csv
	FromDate string `json:"from_date"`

	// ToDate sets the start date for export API.
	ToDate string `json:"to_date"`

	// Format describes the response format. json/csv
	Format string `json:"format"`

	// Event only exports data for this event.
	Event string `json:"event"`

	// SessionID used to maintain session in engage API
	SessionID string `json:"session_id"`

	// Page is page number to get results
	Page string `json:"page"`
}

func DefaultExportQueryOptions(config *Config) *QueryOptions {
	// expire api request 10 minutes from now
	yesterday := time.Now().Add(-24 * time.Hour).Format("2006-01-02")
	q := DefaultQueryOptions(config)
	q.FromDate = yesterday
	q.ToDate = yesterday
	return q
}

func DefaultQueryOptions(config *Config) *QueryOptions {
	expire := time.Now().Add(10 * time.Minute).Unix()
	return &QueryOptions{
		Key:    config.Key,
		Secret: config.Secret,
		Expire: strconv.Itoa(int(expire)),
		Format: "json",
	}
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

func MergeQueryOptions(a, b *QueryOptions) *QueryOptions {
	var result QueryOptions = *a
	if b.FromDate != "" {
		result.FromDate = b.FromDate
	}
	if b.ToDate != "" {
		result.ToDate = b.ToDate
	}
	if b.Format != "" {
		result.Format = b.Format
	}
	if b.Event != "" {
		result.Event = b.Event
	}
	return &result
}
