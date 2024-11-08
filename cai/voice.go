package cai

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"
)

// FetchVoice retrieves a voice by its ID.
func (c *Client) FetchVoice(voiceID string) (*Voice, error) {
	urlStr := fmt.Sprintf("https://neo.character.ai/multimodal/api/v1/voices/%s", voiceID)
	headers := c.GetHeaders(false)

	resp, err := c.Requester.Get(urlStr, headers)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch voice, status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result FetchVoiceResponse
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	if result.Voice == nil {
		return nil, errors.New("voice not found")
	}

	return result.Voice, nil
}

// SearchVoices searches for voices by name.
func (c *Client) SearchVoices(query string) ([]*Voice, error) {
	urlStr := fmt.Sprintf("https://neo.character.ai/multimodal/api/v1/voices/search?query=%s", url.QueryEscape(query))
	headers := c.GetHeaders(false)

	resp, err := c.Requester.Get(urlStr, headers)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to search voices, status code: %d", resp.StatusCode)
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

// UploadVoice uploads a new voice.
// The voiceData parameter should contain the voice file's data in bytes.
// The visibility parameter should be "public" or "private".
func (c *Client) UploadVoice(voiceData []byte, name string, description string, visibility string) (*Voice, error) {
	if len(name) < 3 || len(name) > 20 {
		return nil, errors.New("name must be at least 3 characters and no more than 20")
	}
	if len(description) > 120 {
		return nil, errors.New("description must be no more than 120 characters")
	}
	visibility = strings.ToLower(visibility)
	if visibility != "public" && visibility != "private" {
		return nil, errors.New("visibility must be 'public' or 'private'")
	}

	boundary := fmt.Sprintf("----WebKitFormBoundary%s", generateBoundary())
	headers := map[string]string{
		"Content-Type":  fmt.Sprintf("multipart/form-data; boundary=%s", boundary),
		"authorization": fmt.Sprintf("Token %s", c.Token),
	}

	var body bytes.Buffer

	writer := multipart.NewWriter(&body)
	writer.SetBoundary(boundary)

	// Part for the voice file
	part, err := writer.CreateFormFile("file", "input.mp3")
	if err != nil {
		return nil, err
	}
	_, err = part.Write(voiceData)
	if err != nil {
		return nil, err
	}

	// Part for the JSON metadata
	voiceMetadata := UploadVoiceMetadata{
		Voice: UploadVoiceInfo{
			Name:            name,
			Description:     description,
			Gender:          "neutral",
			Visibility:      visibility,
			PreviewText:     "Good day! Here to make life a little less complicated.",
			AudioSourceType: "file",
		},
	}

	metadataBytes, err := json.Marshal(voiceMetadata)
	if err != nil {
		return nil, err
	}
	err = writer.WriteField("json", string(metadataBytes))
	if err != nil {
		return nil, err
	}

	writer.Close()

	urlStr := "https://neo.character.ai/multimodal/api/v1/voices/"
	resp, err := c.Requester.Post(urlStr, headers, body.Bytes())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("failed to upload voice, status code: %d", resp.StatusCode)
	}

	bodyResp, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result UploadVoiceResponse
	err = json.Unmarshal(bodyResp, &result)
	if err != nil {
		return nil, err
	}

	if result.Voice == nil {
		return nil, errors.New("voice not returned in response")
	}

	// Optionally, call EditVoice to ensure metadata is updated
	return c.EditVoice(result.Voice.VoiceID, name, description, visibility)
}

// EditVoice edits an existing voice.
func (c *Client) EditVoice(voiceID string, name string, description string, visibility string) (*Voice, error) {
	if len(name) < 3 || len(name) > 20 {
		return nil, errors.New("name must be at least 3 characters and no more than 20")
	}
	if len(description) > 120 {
		return nil, errors.New("description must be no more than 120 characters")
	}
	visibility = strings.ToLower(visibility)
	if visibility != "public" && visibility != "private" {
		return nil, errors.New("visibility must be 'public' or 'private'")
	}

	// Fetch the existing voice
	voice, err := c.FetchVoice(voiceID)
	if err != nil {
		return nil, err
	}

	// Prepare the updated voice payload
	updatedVoice := VoiceUpdatePayload{
		Voice: VoiceInfo{
			AudioSourceType: "file",
			BackendID:       voice.VoiceID,
			BackendProvider: "cai",
			CreatorInfo:     voice.CreatorInfo,
			Description:     description,
			Gender:          voice.Gender,
			ID:              voice.VoiceID,
			InternalStatus:  "draft",
			LastUpdateTime:  "0001-01-01T00:00:00Z",
			Name:            name,
			PreviewAudioURI: voice.PreviewAudioURL,
			PreviewText:     voice.PreviewText,
			Visibility:      visibility,
		},
	}

	urlStr := fmt.Sprintf("https://neo.character.ai/multimodal/api/v1/voices/%s", voiceID)
	headers := c.GetHeaders(false)
	bodyBytes, err := json.Marshal(updatedVoice)
	if err != nil {
		return nil, err
	}

	resp, err := c.Requester.DoRequest("PUT", urlStr, headers, bodyBytes)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to edit voice, status code: %d", resp.StatusCode)
	}

	bodyResp, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result UploadVoiceResponse
	err = json.Unmarshal(bodyResp, &result)
	if err != nil {
		return nil, err
	}

	return result.Voice, nil
}

// DeleteVoice deletes a voice by its ID.
func (c *Client) DeleteVoice(voiceID string) error {
	urlStr := fmt.Sprintf("https://neo.character.ai/multimodal/api/v1/voices/%s", voiceID)
	headers := c.GetHeaders(false)

	resp, err := c.Requester.DoRequest("DELETE", urlStr, headers, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to delete voice, status code: %d", resp.StatusCode)
	}

	return nil
}

// GenerateSpeech generates speech audio for a turn using a specific voice.
// Returns the audio data as bytes.
func (c *Client) GenerateSpeech(chatID string, turnID string, candidateID string, voiceID string) ([]byte, error) {
	urlStr := "https://neo.character.ai/multimodal/api/v1/memo/replay"
	headers := c.GetHeaders(false)

	payload := GenerateSpeechPayload{
		CandidateID: candidateID,
		RoomID:      chatID,
		TurnID:      turnID,
		VoiceID:     voiceID,
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

	bodyResp, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		var errorResp ErrorResponse
		err = json.Unmarshal(bodyResp, &errorResp)
		if err != nil {
			return nil, fmt.Errorf("failed to generate speech, status code: %d", resp.StatusCode)
		}
		return nil, fmt.Errorf("failed to generate speech, error: %s", errorResp.Error.Message)
	}

	var result GenerateSpeechResponse
	err = json.Unmarshal(bodyResp, &result)
	if err != nil {
		return nil, err
	}

	audioURL := result.ReplayURL
	if audioURL == "" {
		return nil, errors.New("no audio URL returned")
	}

	// Fetch the audio data
	audioResp, err := c.Requester.Get(audioURL, nil)
	if err != nil {
		return nil, err
	}
	defer audioResp.Body.Close()

	if audioResp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch audio data, status code: %d", audioResp.StatusCode)
	}

	audioData, err := io.ReadAll(audioResp.Body)
	if err != nil {
		return nil, err
	}

	return audioData, nil
}
