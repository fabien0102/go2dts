package client

import (
	"log"
	"net"
	"net/http"
	"time"
)

const (
	apiRoot = "/"
)

const sixtyseconds = 60 * time.Second

func debugPrintln(verbose bool, msg string) {
	if verbose {
		log.Println(msg)
	}
}

func debugPrintf(verbose bool, msg string, args ...interface{}) {
	if verbose {
		log.Printf(msg, args...)
	}
}

// Client is an http client that accepts a time and an authorization that will
// be used for requests.
type Client struct {
	http.Client
	Token      string
	AuthPrefix string
}

// NewClient returns a new client with the timeout and token configured
func NewClient(timeout time.Duration, token, tokenPrefix string) *Client {
	if tokenPrefix == "" {
		tokenPrefix = "Bearer"
	}

	client := &Client{Token: token, AuthPrefix: tokenPrefix}

	if timeout > 0 {
		client.Timeout = timeout
		client.Transport = &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout: timeout,
			}).DialContext,
			IdleConnTimeout:       120 * time.Millisecond,
			ExpectContinueTimeout: 1500 * time.Millisecond,
		}
	}

	return client
}

// Do sends an HTTP request and returns an HTTP response using the http.Client
// configured with the configured timeout and authtoken.
func (c *Client) Do(request *http.Request) (*http.Response, error) {
	if c.Token != "" {
		request.Header.Set("Authorization", c.AuthPrefix+" "+c.Token)
	}

	return c.Client.Do(request)
}
