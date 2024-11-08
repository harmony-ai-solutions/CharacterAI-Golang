package cai

import (
	"fmt"
	"sync"
)

// Client is the main client structure
type Client struct {
	Token       string
	WebNextAuth string
	AccountID   string
	Requester   *Requester
	mutex       sync.Mutex
}

// NewClient creates a new Client instance
func NewClient(token string, webNextAuth string, proxy string) *Client {
	requester := NewRequester(token, proxy)
	return &Client{
		Token:       token,
		WebNextAuth: webNextAuth,
		Requester:   requester,
	}
}

// Authenticate retrieves the account ID
func (c *Client) Authenticate() error {
	account, err := c.FetchMe()
	if err != nil {
		return err
	}
	c.AccountID = account.AccountID
	return nil
}

// GetHeaders returns the headers for requests
func (c *Client) GetHeaders(includeWebNextAuth bool) map[string]string {
	headers := map[string]string{
		"authorization": fmt.Sprintf("Token %s", c.Token),
		"Content-Type":  "application/json",
	}
	if includeWebNextAuth {
		headers["cookie"] = c.WebNextAuth
	}
	return headers
}

// Close cleans up the client, closing any open connections
func (c *Client) Close() error {
	return c.Requester.CloseWebSocket()
}
