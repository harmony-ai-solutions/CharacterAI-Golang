package cai

import (
	"github.com/google/uuid"
	"math/rand"
	"net/http"
)

// Ping checks if the service is reachable
func (c *Client) Ping() (bool, error) {
	urlStr := "https://neo.character.ai/ping/"
	headers := c.GetHeaders(false)

	resp, err := c.Requester.Get(urlStr, headers)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK, nil
}

// Helper function to generate a UUID
func generateUUID() string {
	return uuid.New().String()
}

func generateBoundary() string {
	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	b := make([]rune, 16)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
