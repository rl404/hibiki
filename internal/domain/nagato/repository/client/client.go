package client

import (
	"net/http"
	"time"

	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/rl404/nagato"
)

// Client is client for nagato (mal).
type Client struct {
	client *nagato.Client
}

// New to create new nagato (mal) client.
func New(clientID, clientSecret string) *Client {
	n := nagato.New(clientID)
	n.SetHttpClient(&http.Client{
		Timeout: 10 * time.Second,
		Transport: newrelic.NewRoundTripper(&clientIDTransport{
			clientID: clientID,
		}),
	})
	return &Client{
		client: n,
	}
}

type clientIDTransport struct {
	transport http.RoundTripper
	clientID  string
}

func (c *clientIDTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if c.transport == nil {
		c.transport = http.DefaultTransport
	}
	req.Header.Add("X-MAL-CLIENT-ID", c.clientID)
	return c.transport.RoundTrip(req)
}
