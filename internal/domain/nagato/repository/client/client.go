package client

import "github.com/rl404/nagato"

// Client is client for nagato (mal).
type Client struct {
	client *nagato.Client
}

// New to create new nagato (mal) client.
func New(clientID, clientSecret string) *Client {
	return &Client{
		client: nagato.New(clientID),
	}
}
