package api

import (
	"errors"
	"net/http"
	"os"
)

// Config is used to configure the creation of a client
type Config struct {
	// Scheme is the URI scheme for the Mixpanel server.
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

// DefaultConfig returns a default configuration for the client
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
