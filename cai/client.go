package cai

import (
	"fmt"
	"strconv"
	"sync"
)

// Client is the main client structure
type Client struct {
	Token         string
	WebNextAuth   string
	UserAccountID string
	Requester     *Requester
	mutex         sync.Mutex
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
	if c.Token == "" {
		return fmt.Errorf("token not provided")
	}
	account, err := c.FetchMe()
	if err != nil {
		return err
	}
	c.UserAccountID = strconv.FormatInt(account.User.ID, 10)
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
