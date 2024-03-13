package client

import (
	"net/http"
	"time"

	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/rl404/fairy/limit/atomic"
	"github.com/rl404/nagato"
)

// Client is client for nagato (mal).
type Client struct {
	client *nagato.Client
}

// New to create new nagato (mal) client.
func New(clientID, clientSecret string) *Client {
	n := nagato.New(clientID)
	n.SetLimiter(atomic.New(1, 3*time.Second))
	n.SetHttpClient(&http.Client{
		Timeout: 30 * time.Second,
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
	req.Header.Add("User-Agent", "Hibiki/0.4.13 (github.com/rl404/hibiki)")
	return c.transport.RoundTrip(req)
}
