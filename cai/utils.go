package cai

import (
	"encoding/json"
	"fmt"
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

	return resp.StatusCode == 200, nil
}

// GenerateImage generates images based on a prompt
func (c *Client) GenerateImage(prompt string, numCandidates int) ([]string, error) {
	urlStr := "https://plus.character.ai/chat/character/generate-avatar-options"
	headers := c.GetHeaders(false)

	payload := map[string]interface{}{
		"prompt":         prompt,
		"num_candidates": numCandidates,
		"model_version":  "v1",
	}
	bodyBytes, _ := json.Marshal(payload)

	resp, err := c.Requester.Post(urlStr, headers, bodyBytes)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("failed to generate image, status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	images := []string{}
	if res, ok := result["result"].([]interface{}); ok {
		for _, img := range res {
			imgMap, ok := img.(map[string]interface{})
			if ok {
				if urlStr, ok := imgMap["url"].(string); ok {
					images = append(images, urlStr)
				}
			}
		}
	}

	return images, nil
}
