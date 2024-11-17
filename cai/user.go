package cai

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// FetchUser retrieves a public user's information
func (c *Client) FetchUser(username string) (*PublicUser, error) {
	urlStr := "https://plus.character.ai/chat/user/public/"
	headers := c.GetHeaders(false)

	payload := FetchUserRequest{
		Username: username,
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

	if resp.StatusCode == 500 {
		return nil, nil // User does not exist
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch user, status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result FetchUserResponse
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	if result.PublicUser == nil {
		return nil, fmt.Errorf("failed to fetch user, public user data is missing")
	}

	return result.PublicUser, nil
}

// FollowUser follows a user
func (c *Client) FollowUser(username string) error {
	urlStr := "https://plus.character.ai/chat/user/follow/"
	headers := c.GetHeaders(false)

	payload := FollowUserRequest{
		Username: username,
	}
	bodyBytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	resp, err := c.Requester.Post(urlStr, headers, bodyBytes)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to follow user, status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var result FollowUserResponse
	err = json.Unmarshal(body, &result)
	if err != nil {
		return err
	}

	if result.Status != "OK" {
		return fmt.Errorf("failed to follow user, error: %s", result.Error)
	}

	return nil
}

// UnfollowUser unfollows a user
func (c *Client) UnfollowUser(username string) error {
	urlStr := "https://plus.character.ai/chat/user/unfollow/"
	headers := c.GetHeaders(false)

	payload := FollowUserRequest{
		Username: username,
	}
	bodyBytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	resp, err := c.Requester.Post(urlStr, headers, bodyBytes)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to unfollow user, status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var result FollowUserResponse
	err = json.Unmarshal(body, &result)
	if err != nil {
		return err
	}

	if result.Status != "OK" {
		return fmt.Errorf("failed to unfollow user, error: %s", result.Error)
	}

	return nil
}

// FetchUserVoices retrieves the voices created by a public user.
func (c *Client) FetchUserVoices(username string) ([]*Voice, error) {
	urlStr := fmt.Sprintf("https://neo.character.ai/multimodal/api/v1/voices/search?creatorInfo.username=%s", url.QueryEscape(username))
	headers := c.GetHeaders(false)

	resp, err := c.Requester.Get(urlStr, headers)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch user voices, status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result SearchVoicesResponse
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	return result.Voices, nil
}
