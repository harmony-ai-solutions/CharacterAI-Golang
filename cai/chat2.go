package cai

import (
	"errors"
	"fmt"
	http "github.com/bogdanfinn/fhttp"
	"github.com/gorilla/websocket"
	"strconv"
)

type Chat2 struct {
	Token   string
	Session *Session
	WS      *websocket.Conn
}

func (c *Chat2) NextMessage(char, chatID, parentMsgUUID string) (map[string]interface{}, error) {
	message := map[string]interface{}{
		"command": "generate_turn_candidate",
		"payload": map[string]interface{}{
			"character_id": char,
			"turn_key": map[string]string{
				"turn_id": parentMsgUUID,
				"chat_id": chatID,
			},
		},
	}
	err := c.WS.WriteJSON(message)
	if err != nil {
		return nil, err
	}

	for {
		var response map[string]interface{}
		err := c.WS.ReadJSON(&response)
		if err != nil {
			return nil, err
		}

		turn, exists := response["turn"].(map[string]interface{})
		if !exists {
			return nil, errors.New(response["comment"].(string))
		}

		authorID := turn["author"].(map[string]interface{})["author_id"].(string)
		if _, err := strconv.Atoi(authorID); err != nil {
			if candidates, ok := turn["candidates"].([]interface{}); ok {
				candidate := candidates[0].(map[string]interface{})
				if candidate["is_final"].(bool) {
					return response, nil
				}
			}
		}
	}
}

func (c *Chat2) SendMessage(char, chatID, text string, author map[string]interface{}, turnID, customID, candidateID *string) (map[string]interface{}, error) {
	var turnKey map[string]string
	if customID != nil {
		turnKey = map[string]string{
			"turn_id": *customID,
			"chat_id": chatID,
		}
	} else {
		turnKey = map[string]string{"chat_id": chatID}
	}

	message := map[string]interface{}{
		"command": "create_and_generate_turn",
		"payload": map[string]interface{}{
			"character_id": char,
			"turn": map[string]interface{}{
				"turn_key": turnKey,
				"author":   author,
				"candidates": []map[string]string{
					{"raw_content": text},
				},
			},
		},
	}

	if turnID != nil && candidateID != nil {
		message["update_primary_candidate"] = map[string]interface{}{
			"candidate_id": *candidateID,
			"turn_key": map[string]string{
				"turn_id": *turnID,
				"chat_id": chatID,
			},
		}
	}

	err := c.WS.WriteJSON(message)
	if err != nil {
		return nil, err
	}

	for {
		var response map[string]interface{}
		err := c.WS.ReadJSON(&response)
		if err != nil {
			return nil, err
		}

		turn, exists := response["turn"].(map[string]interface{})
		if !exists {
			return nil, errors.New(response["comment"].(string))
		}

		authorID := turn["author"].(map[string]interface{})["author_id"].(string)
		if _, err := strconv.Atoi(authorID); err != nil {
			if candidates, ok := turn["candidates"].([]interface{}); ok {
				candidate := candidates[0].(map[string]interface{})
				if candidate["is_final"].(bool) {
					return response, nil
				}
			}
		}
	}
}

func (c *Chat2) NewChat(char, chatID, creatorID string, withGreeting ...bool) (map[string]interface{}, map[string]interface{}, error) {
	wGreeting := true
	if len(withGreeting) > 0 {
		wGreeting = withGreeting[0]
	}

	message := map[string]interface{}{
		"command": "create_chat",
		"payload": map[string]interface{}{
			"chat": map[string]interface{}{
				"chat_id":      chatID,
				"creator_id":   creatorID,
				"visibility":   "VISIBILITY_PRIVATE",
				"character_id": char,
				"type":         "TYPE_ONE_ON_ONE",
			},
			"with_greeting": wGreeting,
		},
	}

	err := c.WS.WriteJSON(message)
	if err != nil {
		return nil, nil, err
	}

	var response map[string]interface{}
	err = c.WS.ReadJSON(&response)
	if err != nil {
		return nil, nil, err
	}

	if _, exists := response["chat"]; !exists {
		return nil, nil, errors.New(response["comment"].(string))
	}

	var answer map[string]interface{}
	err = c.WS.ReadJSON(&answer)
	if err != nil {
		return nil, nil, err
	}

	return response, answer, nil
}

func (s *Session) GetHistories(char string, preview int, token string) (map[string]interface{}, error) {
	if char == "" {
		char = "default_char_value" // Replace with the default character value if any
	}
	if preview == 0 {
		preview = 2
	}

	url := fmt.Sprintf("chats/?character_ids=%s&num_preview_turns=%d", char, preview)
	return request(url, s, token, http.MethodGet, nil, false, true)
}

func (s *Session) GetChat(char string, token string) (map[string]interface{}, error) {
	if char == "" {
		char = "default_char_value" // Replace with the default character value if any
	}

	url := fmt.Sprintf("chats/recent/%s", char)
	return request(url, s, token, http.MethodGet, nil, false, true)
}

func (s *Session) GetHistory(chatID string, token string) (map[string]interface{}, error) {
	if chatID == "" {
		return nil, fmt.Errorf("chatID cannot be empty")
	}

	url := fmt.Sprintf("turns/%s/", chatID)
	return request(url, s, token, http.MethodGet, nil, false, true)
}

func (s *Session) Rate(rate int, chatID string, turnID string, candidateID string, token string) (map[string]interface{}, error) {
	data := map[string]interface{}{
		"turn_key": map[string]string{
			"chat_id": chatID,
			"turn_id": turnID,
		},
		"candidate_id": candidateID,
		"annotation": map[string]interface{}{
			"annotation_type":  "star",
			"annotation_value": rate,
		},
	}

	return request("annotation/create", s, token, http.MethodPost, data, false, true)
}

func (c *Chat2) DeleteMessage(chatID string, turnIDs []string) (map[string]interface{}, error) {
	// Set the data for the message to send over the WebSocket
	message := map[string]interface{}{
		"command": "remove_turns",
		"payload": map[string]interface{}{
			"chat_id":  chatID,
			"turn_ids": turnIDs,
		},
	}

	// Convert struct to JSON and send it over the WebSocket
	err := c.WS.WriteJSON(message)
	if err != nil {
		return nil, fmt.Errorf("failed to send message over WebSocket: %v", err)
	}

	// Wait for a response from the server
	var response map[string]interface{}
	err = c.WS.ReadJSON(&response)
	if err != nil {
		return nil, fmt.Errorf("failed to read message from WebSocket: %v", err)
	}

	return response, nil
}
