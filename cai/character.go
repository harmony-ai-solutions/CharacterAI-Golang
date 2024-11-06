package cai

import (
	"encoding/json"
	"errors"
	"fmt"
)

// FetchCharacterInfo retrieves information about a character
func (c *Client) FetchCharacterInfo(characterID string) (*Character, error) {
	url := "https://plus.character.ai/chat/character/info/"
	headers := c.GetHeaders(false)

	payload := map[string]string{
		"external_id": characterID,
	}
	bodyBytes, _ := json.Marshal(payload)

	resp, err := c.Requester.Post(url, headers, bodyBytes)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("failed to fetch character info, status code: %d", resp.StatusCode)
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

	status, _ := result["status"].(string)
	if status == "NOT_OK" {
		errorMsg, _ := result["error"].(string)
		return nil, fmt.Errorf("error fetching character info: %s", errorMsg)
	}

	characterMap, ok := result["character"].(map[string]interface{})
	if !ok {
		return nil, errors.New("invalid response structure")
	}

	character := &Character{
		CharacterID:     characterMap["external_id"].(string),
		Name:            characterMap["name"].(string),
		Title:           characterMap["title"].(string),
		Greeting:        characterMap["greeting"].(string),
		Description:     characterMap["description"].(string),
		Definition:      characterMap["definition"].(string),
		Visibility:      characterMap["visibility"].(string),
		AuthorUsername:  characterMap["user__username"].(string),
		NumInteractions: characterMap["participant__num_interactions"].(string),
		Copyable:        characterMap["copyable"].(bool),
		Identifier:      characterMap["identifier"].(string),
		ImgGenEnabled:   characterMap["img_gen_enabled"].(bool),
		BaseImgPrompt:   characterMap["base_img_prompt"].(string),
		ImgPromptRegex:  characterMap["img_prompt_regex"].(string),
		StripImgPrompt:  characterMap["strip_img_prompt_from_msg"].(bool),
	}

	// Handle avatar
	if avatarFileName, ok := characterMap["avatar_file_name"].(string); ok && avatarFileName != "" {
		character.Avatar = &Avatar{FileName: avatarFileName}
	}

	return character, nil
}
