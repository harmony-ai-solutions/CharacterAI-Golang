package cai

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// FetchCharacterInfo retrieves information about a character.
func (c *Client) FetchCharacterInfo(characterID string) (*Character, error) {
	urlStr := "https://plus.character.ai/chat/character/info/"
	headers := c.GetHeaders(false)

	payload := CharacterInfoPayload{
		ExternalID: characterID,
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
		return nil, fmt.Errorf("failed to fetch character info, status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result CharacterInfoResponse
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	if result.Status != "OK" {
		return nil, fmt.Errorf("error fetching character info: %s", result.Error)
	}

	return result.Character, nil
}

// FetchCharactersByCategory retrieves characters categorized by curated categories.
func (c *Client) FetchCharactersByCategory() (map[string][]*CharacterShort, error) {
	urlStr := "https://plus.character.ai/chat/curated_categories/characters/"
	headers := c.GetHeaders(false)

	resp, err := c.Requester.Get(urlStr, headers)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch characters by category, status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result struct {
		CharactersByCategory map[string][]*CharacterShort `json:"characters_by_curated_category"`
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	return result.CharactersByCategory, nil
}

// FetchRecommendedCharacters retrieves recommended characters for the user.
func (c *Client) FetchRecommendedCharacters() ([]*CharacterShort, error) {
	urlStr := "https://neo.character.ai/recommendation/v1/user"
	headers := c.GetHeaders(false)

	resp, err := c.Requester.Get(urlStr, headers)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch recommended characters, status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result struct {
		Characters []*CharacterShort `json:"characters"`
	}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	return result.Characters, nil
}

// FetchFeaturedCharacters retrieves featured characters.
func (c *Client) FetchFeaturedCharacters() ([]*CharacterShort, error) {
	urlStr := "https://plus.character.ai/chat/characters/featured_v2/"
	headers := c.GetHeaders(false)

	resp, err := c.Requester.Get(urlStr, headers)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch featured characters, status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result struct {
		Characters []*CharacterShort `json:"characters"`
	}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	return result.Characters, nil
}

// FetchSimilarCharacters retrieves characters similar to the given character.
func (c *Client) FetchSimilarCharacters(characterID string) ([]*CharacterShort, error) {
	urlStr := fmt.Sprintf("https://neo.character.ai/recommendation/v1/character/%s", characterID)
	headers := c.GetHeaders(false)

	resp, err := c.Requester.Get(urlStr, headers)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch similar characters, status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result struct {
		Characters []*CharacterShort `json:"characters"`
	}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	return result.Characters, nil
}

// SearchCharacters searches for characters by name.
func (c *Client) SearchCharacters(query string) ([]*CharacterSearchResult, error) {
	urlStr := fmt.Sprintf("https://plus.character.ai/chat/characters/search/?query=%s", url.QueryEscape(query))
	headers := c.GetHeaders(false)

	resp, err := c.Requester.Get(urlStr, headers)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to search characters, status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result struct {
		Characters []*CharacterSearchResult `json:"characters"`
	}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	return result.Characters, nil
}

// SearchCreators searches for creators by name.
func (c *Client) SearchCreators(query string) ([]Creator, error) {
	urlStr := fmt.Sprintf("https://plus.character.ai/chat/creators/search/?query=%s", url.QueryEscape(query))
	headers := c.GetHeaders(false)

	resp, err := c.Requester.Get(urlStr, headers)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to search creators, status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result SearchCreatorsResponse
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	if result.Status != "OK" {
		return nil, fmt.Errorf("error searching creators: %s", result.Error)
	}

	return result.Creators, nil
}

// CharacterVote casts a vote for a character. Use 'nil' for removing a vote.
func (c *Client) CharacterVote(characterID string, vote *bool) error {
	urlStr := "https://plus.character.ai/chat/character/vote/"
	headers := c.GetHeaders(false)

	payload := CharacterVotePayload{
		ExternalID: characterID,
		Vote:       vote,
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
		return fmt.Errorf("failed to vote on character, status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var result CharacterVoteResponse
	err = json.Unmarshal(body, &result)
	if err != nil {
		return err
	}

	if result.Status != "OK" {
		return fmt.Errorf("error voting on character: %s", result.Error)
	}

	return nil
}

// CreateCharacter creates a new character.
func (c *Client) CreateCharacter(name, greeting string, title, description, definition string, copyable bool, visibility, avatarRelPath, defaultVoiceID string) (*Character, error) {
	// Validate inputs
	if len(name) < 3 || len(name) > 20 {
		return nil, errors.New("name must be at least 3 characters and no more than 20")
	}
	if len(greeting) < 3 || len(greeting) > 2048 {
		return nil, errors.New("greeting must be at least 3 characters and no more than 2048")
	}
	visibility = strings.ToUpper(visibility)
	if visibility != "UNLISTED" && visibility != "PUBLIC" && visibility != "PRIVATE" {
		return nil, errors.New(`visibility must be "unlisted", "public", or "private"`)
	}
	if title != "" && (len(title) < 3 || len(title) > 50) {
		return nil, errors.New("title must be at least 3 characters and no more than 50")
	}
	if len(description) > 500 {
		return nil, errors.New("description must be no more than 500 characters")
	}
	if len(definition) > 32000 {
		return nil, errors.New("definition must be no more than 32000 characters")
	}

	urlStr := "https://plus.character.ai/chat/character/create/"
	headers := c.GetHeaders(false)

	payload := CreateCharacterPayload{
		AvatarRelPath:         avatarRelPath,
		BaseImgPrompt:         "",
		Categories:            []string{},
		Copyable:              copyable,
		DefaultVoiceID:        defaultVoiceID,
		Definition:            definition,
		Description:           description,
		Greeting:              greeting,
		Identifier:            fmt.Sprintf("id:%s", generateUUID()),
		ImgGenEnabled:         false,
		Name:                  name,
		StripImgPromptFromMsg: false,
		Title:                 title,
		Visibility:            visibility,
		VoiceID:               "",
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
		return nil, fmt.Errorf("failed to create character, status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result CreateCharacterResult
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	if result.Status != "OK" {
		return nil, fmt.Errorf("failed to create character, error: %s", result.Error)
	}

	if result.Character == nil {
		return nil, errors.New("character not returned in response")
	}

	return result.Character, nil
}

// EditCharacter edits an existing character.
func (c *Client) EditCharacter(characterID, name, greeting string, title, description, definition string, copyable bool, visibility, avatarRelPath, defaultVoiceID string) (*Character, error) {
	// Validate inputs
	if len(name) < 3 || len(name) > 20 {
		return nil, errors.New("name must be at least 3 characters and no more than 20")
	}
	if len(greeting) < 3 || len(greeting) > 2048 {
		return nil, errors.New("greeting must be at least 3 characters and no more than 2048")
	}
	visibility = strings.ToUpper(visibility)
	if visibility != "UNLISTED" && visibility != "PUBLIC" && visibility != "PRIVATE" {
		return nil, errors.New(`visibility must be "unlisted", "public", or "private"`)
	}
	if title != "" && (len(title) < 3 || len(title) > 50) {
		return nil, errors.New("title must be at least 3 characters and no more than 50")
	}
	if len(description) > 500 {
		return nil, errors.New("description must be no more than 500 characters")
	}
	if len(definition) > 32000 {
		return nil, errors.New("definition must be no more than 32000 characters")
	}

	urlStr := "https://plus.character.ai/chat/character/update/"
	headers := c.GetHeaders(false)

	payload := EditCharacterPayload{
		Archived:              false,
		AvatarRelPath:         avatarRelPath,
		BaseImgPrompt:         "",
		Categories:            []string{},
		Copyable:              copyable,
		DefaultVoiceID:        defaultVoiceID,
		Definition:            definition,
		Description:           description,
		ExternalID:            characterID,
		Greeting:              greeting,
		ImgGenEnabled:         false,
		Name:                  name,
		StripImgPromptFromMsg: false,
		Title:                 title,
		Visibility:            visibility,
		VoiceID:               "",
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
		return nil, fmt.Errorf("failed to edit character, status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result EditCharacterResult
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	if result.Status != "OK" {
		return nil, fmt.Errorf("failed to edit character, error: %s", result.Error)
	}

	if result.Character == nil {
		return nil, errors.New("character not returned in response")
	}

	return result.Character, nil
}
