package cai

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// GenerateImage generates images based on a prompt.
func (c *Client) GenerateImage(prompt string, numCandidates int) ([]string, error) {
	urlStr := "https://plus.character.ai/chat/character/generate-avatar-options"
	headers := c.GetHeaders(false)

	payload := GenerateImageRequest{
		Prompt:        prompt,
		NumCandidates: numCandidates,
		ModelVersion:  "v1",
	}
	bodyBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	resp, err := c.Requester.Post(urlStr, headers, bodyBytes)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to generate image, status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result GenerateImageResponse
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	images := make([]string, len(result.Result))
	for i, img := range result.Result {
		images[i] = img.URL
	}

	return images, nil
}
