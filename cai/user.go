package cai

import (
	"encoding/json"
	"errors"
	"fmt"
)

// FetchUser retrieves a public user's information
func (c *Client) FetchUser(username string) (*PublicUser, error) {
	urlStr := "https://plus.character.ai/chat/user/public/"
	headers := c.GetHeaders(false)

	payload := map[string]string{
		"username": username,
	}
	bodyBytes, _ := json.Marshal(payload)

	resp, err := c.Requester.Post(urlStr, headers, bodyBytes)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 500 {
		return nil, nil // User does not exist
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("failed to fetch user, status code: %d", resp.StatusCode)
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

	userMap, ok := result["public_user"].(map[string]interface{})
	if !ok {
		return nil, errors.New("invalid response structure")
	}

	user := &PublicUser{
		Username: userMap["username"].(string),
		Name:     userMap["name"].(string),
		Bio:      userMap["bio"].(string),
	}

	// Handle avatar
	if avatarFileName, ok := userMap["avatar_file_name"].(string); ok && avatarFileName != "" {
		user.Avatar = &Avatar{FileName: avatarFileName}
	}

	// Handle other fields
	if numFollowers, ok := userMap["num_followers"].(float64); ok {
		user.NumFollowers = int(numFollowers)
	}
	if numFollowing, ok := userMap["num_following"].(float64); ok {
		user.NumFollowing = int(numFollowing)
	}

	// Handle characters
	if chars, ok := userMap["characters"].([]interface{}); ok {
		user.Characters = make([]*CharacterShort, len(chars))
		for i, char := range chars {
			charMap, ok := char.(map[string]interface{})
			if ok {
				user.Characters[i] = parseCharacterShort(charMap)
			}
		}
	}

	return user, nil
}

// parseCharacterShort parses a CharacterShort from a map
func parseCharacterShort(charMap map[string]interface{}) *CharacterShort {
	char := &CharacterShort{
		CharacterID:     charMap["external_id"].(string),
		Name:            charMap["name"].(string),
		Title:           charMap["title"].(string),
		Greeting:        charMap["greeting"].(string),
		Description:     charMap["description"].(string),
		Definition:      charMap["definition"].(string),
		Visibility:      charMap["visibility"].(string),
		AuthorUsername:  charMap["user__username"].(string),
		NumInteractions: charMap["participant__num_interactions"].(string),
	}

	// Handle avatar
	if avatarFileName, ok := charMap["avatar_file_name"].(string); ok && avatarFileName != "" {
		char.Avatar = &Avatar{FileName: avatarFileName}
	}

	return char
}

// FollowUser follows a user
func (c *Client) FollowUser(username string) error {
	urlStr := "https://plus.character.ai/chat/user/follow/"
	headers := c.GetHeaders(false)

	payload := map[string]string{
		"username": username,
	}
	bodyBytes, _ := json.Marshal(payload)

	resp, err := c.Requester.Post(urlStr, headers, bodyBytes)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("failed to follow user, status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return err
	}

	if status, ok := result["status"].(string); !ok || status != "OK" {
		return errors.New("failed to follow user")
	}

	return nil
}

// UnfollowUser unfollows a user
func (c *Client) UnfollowUser(username string) error {
	urlStr := "https://plus.character.ai/chat/user/unfollow/"
	headers := c.GetHeaders(false)

	payload := map[string]string{
		"username": username,
	}
	bodyBytes, _ := json.Marshal(payload)

	resp, err := c.Requester.Post(urlStr, headers, bodyBytes)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("failed to unfollow user, status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return err
	}

	if status, ok := result["status"].(string); !ok || status != "OK" {
		return errors.New("failed to unfollow user")
	}

	return nil
}
