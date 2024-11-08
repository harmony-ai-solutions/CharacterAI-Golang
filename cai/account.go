// account.go

package cai

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
)

// FetchMe retrieves the account information.
func (c *Client) FetchMe() (*Account, error) {
	urlStr := "https://beta.character.ai/chat/user/"
	headers := c.GetHeaders(false)

	resp, err := c.Requester.Get(urlStr, headers)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch account info, status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result FetchMeResponse
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	account := &result.User
	if account == nil {
		return nil, errors.New("account information not found")
	}

	return account, nil
}

// GetSelf is an alias for FetchMe.
func (c *Client) GetSelf() (*Account, error) {
	return c.FetchMe()
}

// FetchMySettings retrieves the user's settings.
func (c *Client) FetchMySettings() (*Settings, error) {
	urlStr := "https://plus.character.ai/chat/user/settings/"
	headers := c.GetHeaders(false)

	resp, err := c.Requester.Get(urlStr, headers)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch settings, status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var settings Settings
	err = json.Unmarshal(body, &settings)
	if err != nil {
		return nil, err
	}

	return &settings, nil
}

// FetchMyFollowers retrieves the user's followers.
func (c *Client) FetchMyFollowers() ([]*PublicUser, error) {
	urlStr := "https://plus.character.ai/chat/user/followers/"
	headers := c.GetHeaders(false)

	resp, err := c.Requester.Get(urlStr, headers)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch followers, status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result struct {
		Followers []*PublicUser `json:"followers"`
	}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	return result.Followers, nil
}

// FetchMyFollowing retrieves the users that the user is following.
func (c *Client) FetchMyFollowing() ([]*PublicUser, error) {
	urlStr := "https://plus.character.ai/chat/user/following/"
	headers := c.GetHeaders(false)

	resp, err := c.Requester.Get(urlStr, headers)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch following, status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result struct {
		Following []*PublicUser `json:"following"`
	}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	return result.Following, nil
}

// FetchMyPersonas retrieves all the user's personas.
func (c *Client) FetchMyPersonas() ([]*Character, error) {
	urlStr := "https://plus.character.ai/chat/personas/?force_refresh=1"
	headers := c.GetHeaders(false)

	resp, err := c.Requester.Get(urlStr, headers)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch personas, status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result struct {
		Personas []*Character `json:"personas"`
	}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	return result.Personas, nil
}

// FetchMyCharacters retrieves all the user's characters.
func (c *Client) FetchMyCharacters() ([]*CharacterShort, error) {
	urlStr := "https://plus.character.ai/chat/characters/?scope=user"
	headers := c.GetHeaders(false)

	resp, err := c.Requester.Get(urlStr, headers)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch characters, status code: %d", resp.StatusCode)
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

// FetchMyUpvotedCharacters retrieves the characters the user has upvoted.
func (c *Client) FetchMyUpvotedCharacters() ([]*CharacterShort, error) {
	urlStr := "https://plus.character.ai/chat/user/characters/upvoted/"
	headers := c.GetHeaders(false)

	resp, err := c.Requester.Get(urlStr, headers)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch upvoted characters, status code: %d", resp.StatusCode)
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

// FetchMyVoices retrieves the user's voices.
func (c *Client) FetchMyVoices() ([]*Voice, error) {
	urlStr := "https://neo.character.ai/multimodal/api/v1/voices/user"
	headers := c.GetHeaders(false)

	resp, err := c.Requester.Get(urlStr, headers)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch voices, status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result struct {
		Voices []*Voice `json:"voices"`
	}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	return result.Voices, nil
}

// UpdateSettings updates the user's settings.
func (c *Client) UpdateSettings(newSettings *Settings) (*Settings, error) {
	urlStr := "https://plus.character.ai/chat/user/update_settings/"
	headers := c.GetHeaders(false)

	bodyBytes, err := json.Marshal(newSettings)
	if err != nil {
		return nil, err
	}

	resp, err := c.Requester.Post(urlStr, headers, bodyBytes)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to update settings, status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result struct {
		Success  bool      `json:"success"`
		Settings *Settings `json:"settings"`
	}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	if !result.Success {
		return nil, errors.New("failed to update settings")
	}

	return result.Settings, nil
}

// EditAccount edits the user's account information.
func (c *Client) EditAccount(name string, username string, bio string, avatarRelPath string) error {
	if len(username) < 2 || len(username) > 20 {
		return errors.New("username must be at least 2 characters and no more than 20")
	}

	if len(name) < 2 || len(name) > 50 {
		return errors.New("name must be at least 2 characters and no more than 50")
	}

	if len(bio) > 500 {
		return errors.New("bio must be no more than 500 characters")
	}

	urlStr := "https://plus.character.ai/chat/user/update/"
	headers := c.GetHeaders(false)

	newAccountInfo := UpdateProfilePayload{
		AvatarType: "DEFAULT",
		Bio:        bio,
		Name:       name,
		Username:   username,
	}

	if avatarRelPath != "" {
		newAccountInfo.AvatarRelPath = avatarRelPath
		newAccountInfo.AvatarType = "UPLOADED"
	}

	bodyBytes, err := json.Marshal(newAccountInfo)
	if err != nil {
		return err
	}

	resp, err := c.Requester.Post(urlStr, headers, bodyBytes)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to edit account, status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var result struct {
		Status string `json:"status"`
		Error  string `json:"error"`
	}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return err
	}

	if result.Status != "OK" {
		return fmt.Errorf("failed to edit account, error: %s", result.Error)
	}

	return nil
}

// FetchMyPersona retrieves a user's persona by ID.
func (c *Client) FetchMyPersona(personaID string) (*Persona, error) {
	urlStr := fmt.Sprintf("https://plus.character.ai/chat/persona/?id=%s", url.QueryEscape(personaID))
	headers := c.GetHeaders(false)

	resp, err := c.Requester.Get(urlStr, headers)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch persona, status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result struct {
		Persona *Persona `json:"persona"`
	}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	if result.Persona == nil {
		return nil, errors.New("persona not found")
	}

	return result.Persona, nil
}

// CreatePersona creates a new persona.
func (c *Client) CreatePersona(name string, definition string, avatarRelPath string) (*Persona, error) {
	if len(name) < 3 || len(name) > 20 {
		return nil, errors.New("name must be at least 3 characters and no more than 20")
	}

	if len(definition) > 728 {
		return nil, errors.New("definition must be no more than 728 characters")
	}

	urlStr := "https://plus.character.ai/chat/character/create/"
	headers := c.GetHeaders(false)

	payload := CreatePersonaPayload{
		Name:                  name,
		Title:                 name,
		Definition:            definition,
		Greeting:              "Hello! This is my persona.",
		Description:           "This is my persona.",
		Visibility:            "PRIVATE",
		AvatarFileName:        "",
		AvatarRelPath:         avatarRelPath,
		VoiceID:               "",
		Identifier:            fmt.Sprintf("id:%s", generateUUID()),
		Categories:            []string{},
		BaseImgPrompt:         "",
		ImgGenEnabled:         false,
		Copyable:              false,
		StripImgPromptFromMsg: false,
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
		return nil, fmt.Errorf("failed to create persona, status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result CreatePersonaResult
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	if result.Status != "OK" {
		return nil, fmt.Errorf("failed to create persona, error: %s", result.Error)
	}

	if result.Persona == nil {
		return nil, errors.New("persona not returned in response")
	}

	return result.Persona, nil
}

// EditPersona edits an existing persona.
func (c *Client) EditPersona(personaID string, name string, definition string, avatarRelPath string) (*Persona, error) {
	// Fetch the existing persona
	oldPersona, err := c.FetchMyPersona(personaID)
	if err != nil {
		return nil, err
	}

	// Prepare payload with existing values
	payload := EditPersonaPayload{
		ExternalID:            personaID,
		Name:                  oldPersona.Name,
		Definition:            oldPersona.Definition,
		Greeting:              oldPersona.Greeting,
		Description:           oldPersona.Description,
		Visibility:            oldPersona.Visibility,
		AvatarFileName:        oldPersona.AvatarFileName,
		AvatarRelPath:         oldPersona.AvatarFileName,
		VoiceID:               oldPersona.VoiceID,
		Identifier:            oldPersona.Identifier,
		Categories:            oldPersona.Categories,
		BaseImgPrompt:         oldPersona.BaseImgPrompt,
		ImgGenEnabled:         oldPersona.ImgGenEnabled,
		Copyable:              oldPersona.Copyable,
		StripImgPromptFromMsg: oldPersona.StripImgPromptFromMsg,
	}

	// Update fields if provided
	if name != "" {
		if len(name) < 3 || len(name) > 20 {
			return nil, errors.New("name must be at least 3 characters and no more than 20")
		}
		payload.Name = name
		payload.Title = name
	}

	if definition != "" {
		if len(definition) > 728 {
			return nil, errors.New("definition must be no more than 728 characters")
		}
		payload.Definition = definition
	}

	if avatarRelPath != "" {
		payload.AvatarFileName = avatarRelPath
		payload.AvatarRelPath = avatarRelPath
	}

	urlStr := "https://plus.character.ai/chat/character/update/"
	headers := c.GetHeaders(false)

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
		return nil, fmt.Errorf("failed to edit persona, status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result EditPersonaResult
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	if result.Status != "OK" {
		return nil, fmt.Errorf("failed to edit persona, error: %s", result.Error)
	}

	if result.Persona == nil {
		return nil, errors.New("persona not returned in response")
	}

	return result.Persona, nil
}

// DeletePersona deletes a persona by marking it as archived.
func (c *Client) DeletePersona(personaID string) error {
	// Fetch the existing persona
	oldPersona, err := c.FetchMyPersona(personaID)
	if err != nil {
		return err
	}

	// Prepare payload with existing values
	payload := EditPersonaPayload{
		ExternalID:            personaID,
		Name:                  oldPersona.Name,
		Definition:            oldPersona.Definition,
		Greeting:              oldPersona.Greeting,
		Description:           oldPersona.Description,
		Visibility:            oldPersona.Visibility,
		AvatarFileName:        oldPersona.AvatarFileName,
		AvatarRelPath:         oldPersona.AvatarFileName,
		VoiceID:               oldPersona.VoiceID,
		Identifier:            oldPersona.Identifier,
		Categories:            oldPersona.Categories,
		BaseImgPrompt:         oldPersona.BaseImgPrompt,
		ImgGenEnabled:         oldPersona.ImgGenEnabled,
		Copyable:              oldPersona.Copyable,
		StripImgPromptFromMsg: oldPersona.StripImgPromptFromMsg,
		Archived:              true,
	}

	urlStr := "https://plus.character.ai/chat/character/update/"
	headers := c.GetHeaders(false)

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
		return fmt.Errorf("failed to delete persona, status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var result DeletePersonaResult
	err = json.Unmarshal(body, &result)
	if err != nil {
		return err
	}

	if result.Status != "OK" {
		return fmt.Errorf("failed to delete persona, error: %s", result.Error)
	}

	return nil
}

// SetDefaultPersona sets the default persona for the user.
func (c *Client) SetDefaultPersona(personaID string) error {
	settings, err := c.FetchMySettings()
	if err != nil {
		return err
	}
	settings.DefaultPersonaID = personaID

	_, err = c.UpdateSettings(settings)
	if err != nil {
		return fmt.Errorf("failed to set default persona: %v", err)
	}

	return nil
}

// UnsetDefaultPersona unsets the default persona for the user.
func (c *Client) UnsetDefaultPersona() error {
	return c.SetDefaultPersona("")
}

// SetPersona sets the persona override for a character.
func (c *Client) SetPersona(characterID string, personaID string) error {
	settings, err := c.FetchMySettings()
	if err != nil {
		return err
	}

	if settings.PersonaOverrides == nil {
		settings.PersonaOverrides = make(map[string]string)
	}

	settings.PersonaOverrides[characterID] = personaID

	_, err = c.UpdateSettings(settings)
	if err != nil {
		return fmt.Errorf("failed to set persona override: %v", err)
	}

	return nil
}

// UnsetPersona unsets the persona override for a character.
func (c *Client) UnsetPersona(characterID string) error {
	settings, err := c.FetchMySettings()
	if err != nil {
		return err
	}

	if settings.PersonaOverrides == nil {
		return nil // Nothing to unset
	}

	delete(settings.PersonaOverrides, characterID)

	_, err = c.UpdateSettings(settings)
	if err != nil {
		return fmt.Errorf("failed to unset persona override: %v", err)
	}

	return nil
}

// SetVoice sets the voice override for a character.
func (c *Client) SetVoice(characterID string, voiceID string) error {
	urlStr := fmt.Sprintf("https://plus.character.ai/chat/character/%s/voice_override/update/", characterID)
	headers := c.GetHeaders(false)

	payload := SetVoicePayload{
		VoiceID: voiceID,
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
		return fmt.Errorf("failed to set voice, status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var result SetVoiceResponse
	err = json.Unmarshal(body, &result)
	if err != nil {
		return err
	}

	if !result.Success {
		return fmt.Errorf("failed to set voice, error: %s", result.Error)
	}

	return nil
}

// UnsetVoice unsets the voice override for a character.
func (c *Client) UnsetVoice(characterID string) error {
	urlStr := fmt.Sprintf("https://plus.character.ai/chat/character/%s/voice_override/delete/", characterID)
	headers := c.GetHeaders(false)

	resp, err := c.Requester.Post(urlStr, headers, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to unset voice, status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var result SetVoiceResponse
	err = json.Unmarshal(body, &result)
	if err != nil {
		return err
	}

	if !result.Success {
		return fmt.Errorf("failed to unset voice, error: %s", result.Error)
	}

	return nil
}

// UploadAvatar uploads a new avatar for the user.
// The imagePath can be a local file path or a URL.
// The checkImage parameter determines whether to verify the uploaded image.
func (c *Client) UploadAvatar(imagePath string, checkImage bool) (*Avatar, error) {
	var imageData []byte
	var err error

	if _, err = os.Stat(imagePath); err == nil {
		// It's a file path
		imageData, err = os.ReadFile(imagePath)
		if err != nil {
			return nil, err
		}
	} else {
		// Try to treat it as a URL
		resp, err := c.Requester.Get(imagePath, nil)
		if err != nil {
			return nil, errors.New("invalid image path or URL")
		}
		defer resp.Body.Close()
		imageData, err = io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
	}

	mimeType := http.DetectContentType(imageData)
	dataURI := fmt.Sprintf("data:%s;base64,%s", mimeType, base64.StdEncoding.EncodeToString(imageData))

	urlStr := "https://character.ai/api/trpc/user.uploadAvatar?batch=1"
	headers := c.GetHeaders(true)

	// Prepare the payload using the defined structs
	avatarPayload := UploadAvatarRequest{
		"0": AvatarRequestItem{
			JSON: AvatarJSON{
				ImageDataURL: dataURI,
			},
		},
	}
	bodyBytes, err := json.Marshal(avatarPayload)
	if err != nil {
		return nil, err
	}

	resp, err := c.Requester.Post(urlStr, headers, bodyBytes)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to upload avatar, status code: %d", resp.StatusCode)
	}

	bodyResp, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var response []UploadAvatarResponse
	err = json.Unmarshal(bodyResp, &response)
	if err != nil {
		return nil, err
	}

	if len(response) == 0 || response[0].Result.Data.JSON == "" {
		return nil, errors.New("failed to upload avatar, file name not returned")
	}

	fileName := response[0].Result.Data.JSON
	avatar := &Avatar{FileName: fileName}

	if checkImage {
		size := 150
		imageURL := avatar.GetURL(size, false)

		resp, err := c.Requester.Get(imageURL, nil)
		if err != nil || resp.StatusCode != http.StatusOK {
			return nil, errors.New("uploaded avatar did not pass the filter or is invalid")
		}
	}

	return avatar, nil
}
