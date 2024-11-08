package cai

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

func (c *Client) SendMessage(characterID, chatID, text string) (*Turn, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	// Initialize WebSocket connection if not connected
	err := c.Requester.InitializeWebSocket()
	if err != nil {
		return nil, err
	}

	candidateID := generateUUID()
	turnID := generateUUID()
	requestID := generateUUID()

	// Construct the message
	message := WebSocketMessage{
		Command:   "create_and_generate_turn",
		OriginID:  "web-next",
		RequestID: requestID,
		Payload: CreateAndGenerateTurnPayload{
			CharacterID:         characterID,
			NumCandidates:       1,
			PreviousAnnotations: generatePreviousAnnotations(),
			SelectedLanguage:    "",
			TTSEnabled:          false,
			UserName:            "",
			Turn: TurnPayload{
				Author: AuthorPayload{
					AuthorID: c.UserAccountID,
					IsHuman:  true,
					Name:     "",
				},
				Candidates: []CandidatePayload{
					{
						CandidateID: candidateID,
						RawContent:  text,
					},
				},
				PrimaryCandidateID: candidateID,
				TurnKey: TurnKey{
					ChatID: chatID,
					TurnID: turnID,
				},
			},
		},
	}

	// Send the message
	err = c.Requester.SendWebSocketMessage(message)
	if err != nil {
		return nil, err
	}

	// Receive response
	for {
		responseBytes, err := c.Requester.ReceiveRawWebSocketMessage()
		if err != nil {
			return nil, err
		}

		var response WebSocketResponse
		err = json.Unmarshal(responseBytes, &response)
		if err != nil {
			return nil, err
		}

		switch response.Command {
		case "neo_error":
			return nil, errors.New(response.Comment)
		case "add_turn", "update_turn":
			var payload TurnResponsePayload
			err = json.Unmarshal(response.Payload, &payload)
			if err != nil {
				return nil, err
			}
			return &payload.Turn, nil
		}
	}
}

// CreateChat creates a new chat with a character
func (c *Client) CreateChat(characterID string, greeting bool) (*Chat, *Turn, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	// Initialize WebSocket connection if not connected
	err := c.Requester.InitializeWebSocket()
	if err != nil {
		return nil, nil, err
	}

	requestID := generateUUID()
	chatID := generateUUID()

	// Construct the message
	message := WebSocketMessage{
		Command:   "create_chat",
		RequestID: requestID,
		Payload: CreateChatPayload{
			Chat: ChatPayload{
				ChatID:      chatID,
				CreatorID:   c.UserAccountID,
				Visibility:  "VISIBILITY_PRIVATE",
				CharacterID: characterID,
				Type:        "TYPE_ONE_ON_ONE",
			},
			WithGreeting: greeting,
		},
	}

	// Send the message
	err = c.Requester.SendWebSocketMessage(message)
	if err != nil {
		return nil, nil, err
	}

	var newChat *Chat
	var greetingTurn *Turn

	// Receive response
	for {
		responseBytes, err := c.Requester.ReceiveRawWebSocketMessage()
		if err != nil {
			return nil, nil, err
		}

		var response WebSocketResponse
		err = json.Unmarshal(responseBytes, &response)
		if err != nil {
			return nil, nil, err
		}

		switch response.Command {
		case "neo_error":
			return nil, nil, errors.New(response.Comment)
		case "create_chat_response":
			var payload CreateChatResponsePayload
			err = json.Unmarshal(responseBytes, &payload)
			if err != nil {
				return nil, nil, err
			}
			newChat = &payload.Chat

			if !greeting {
				return newChat, nil, nil
			}
			// Continue to wait for greeting turn
		case "add_turn":
			var payload TurnResponsePayload
			err = json.Unmarshal(responseBytes, &payload)
			if err != nil {
				return nil, nil, err
			}
			greetingTurn = &payload.Turn
			return newChat, greetingTurn, nil
		}
	}
}

// FetchHistories retrieves chat histories for a character
func (c *Client) FetchHistories(characterID string, amount int) ([]ChatHistory, error) {
	urlStr := "https://plus.character.ai/chat/character/histories/"
	headers := c.GetHeaders(false)

	payload := FetchHistoriesRequest{
		ExternalID: characterID,
		Number:     amount,
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
		return nil, fmt.Errorf("failed to fetch histories, status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result FetchHistoriesResponse
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	// Parse time strings to time.Time
	for i := range result.Histories {
		history := &result.Histories[i]
		history.CreateTime, err = time.Parse(time.RFC3339Nano, history.CreateTimeStr)
		if err != nil {
			return nil, err
		}
		history.LastInteraction, err = time.Parse(time.RFC3339Nano, history.LastInteractionStr)
		if err != nil {
			return nil, err
		}
	}

	return result.Histories, nil
}

// FetchChats retrieves chats for a character
func (c *Client) FetchChats(characterID string, numPreviewTurns int) ([]*Chat, error) {
	urlStr := fmt.Sprintf("https://neo.character.ai/chats/?character_ids=%s&num_preview_turns=%d", characterID, numPreviewTurns)
	headers := c.GetHeaders(false)

	resp, err := c.Requester.Get(urlStr, headers)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch chats, status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result FetchChatsResponse
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	for _, chat := range result.Chats {
		chat.CreateTime, err = time.Parse(time.RFC3339Nano, chat.CreateTimeStr)
		if err != nil {
			return nil, err
		}
		if chat.CharacterAvatarURI != "" {
			chat.CharacterAvatar = &Avatar{FileName: chat.CharacterAvatarURI}
		}
	}

	return result.Chats, nil
}

// FetchChat retrieves a chat by its ID
func (c *Client) FetchChat(chatID string) (*Chat, error) {
	urlStr := fmt.Sprintf("https://neo.character.ai/chat/%s/", chatID)
	headers := c.GetHeaders(false)

	resp, err := c.Requester.Get(urlStr, headers)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch chat, status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var response struct {
		Chat *Chat `json:"chat"`
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	// Process the chat data
	chat := response.Chat

	// Convert time strings to time.Time
	chat.CreateTime, err = time.Parse(time.RFC3339Nano, chat.CreateTimeStr)
	if err != nil {
		return nil, err
	}

	// Handle character avatar if available
	if chat.CharacterAvatarURI != "" {
		chat.CharacterAvatar = &Avatar{FileName: chat.CharacterAvatarURI}
	}

	return chat, nil
}

// FetchRecentChats retrieves recent chats for the user
func (c *Client) FetchRecentChats() ([]*Chat, error) {
	urlStr := "https://neo.character.ai/chats/recent/"
	headers := c.GetHeaders(false)

	resp, err := c.Requester.Get(urlStr, headers)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch recent chats, status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result FetchChatsResponse
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	for _, chat := range result.Chats {
		chat.CreateTime, err = time.Parse(time.RFC3339Nano, chat.CreateTimeStr)
		if err != nil {
			return nil, err
		}
		if chat.CharacterAvatarURI != "" {
			chat.CharacterAvatar = &Avatar{FileName: chat.CharacterAvatarURI}
		}
	}

	return result.Chats, nil
}

// FetchMessages retrieves messages from a chat
func (c *Client) FetchMessages(chatID string, pinnedOnly bool, nextToken string) ([]*Turn, string, error) {
	urlStr := fmt.Sprintf("https://neo.character.ai/turns/%s/", chatID)
	if nextToken != "" {
		urlStr += fmt.Sprintf("?next_token=%s", url.QueryEscape(nextToken))
	}
	headers := c.GetHeaders(false)

	resp, err := c.Requester.Get(urlStr, headers)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, "", fmt.Errorf("failed to fetch messages, status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, "", err
	}

	var result FetchMessagesResponse
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, "", err
	}

	var turns []*Turn
	for i := range result.Turns {
		turn := &result.Turns[i]
		if !pinnedOnly || turn.IsPinned {
			turns = append(turns, turn)
		}
	}

	return turns, result.Meta.NextToken, nil
}

// FetchAllMessages retrieves all messages from a chat
func (c *Client) FetchAllMessages(chatID string, pinnedOnly bool) ([]*Turn, error) {
	var allTurns []*Turn
	var nextToken string
	for {
		turns, token, err := c.FetchMessages(chatID, pinnedOnly, nextToken)
		if err != nil {
			return nil, err
		}
		allTurns = append(allTurns, turns...)
		if token == "" {
			break
		}
		nextToken = token
	}
	return allTurns, nil
}

// UpdateChatName updates the name of a chat
func (c *Client) UpdateChatName(chatID string, name string) error {
	urlStr := fmt.Sprintf("https://neo.character.ai/chat/%s/update_name", chatID)
	headers := c.GetHeaders(false)

	payload := UpdateChatNamePayload{
		Name: name,
	}
	bodyBytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	resp, err := c.Requester.DoRequest("PATCH", urlStr, headers, bodyBytes)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to update chat name, status code: %d", resp.StatusCode)
	}

	return nil
}

// ArchiveChat archives a chat
func (c *Client) ArchiveChat(chatID string) error {
	urlStr := fmt.Sprintf("https://neo.character.ai/chat/%s/archive", chatID)
	headers := c.GetHeaders(false)

	resp, err := c.Requester.DoRequest("PATCH", urlStr, headers, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to archive chat, status code: %d", resp.StatusCode)
	}

	return nil
}

// UnarchiveChat unarchives a chat
func (c *Client) UnarchiveChat(chatID string) error {
	urlStr := fmt.Sprintf("https://neo.character.ai/chat/%s/unarchive", chatID)
	headers := c.GetHeaders(false)

	resp, err := c.Requester.DoRequest("PATCH", urlStr, headers, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to unarchive chat, status code: %d", resp.StatusCode)
	}

	return nil
}

// CopyChat copies a chat up to a specific turn
func (c *Client) CopyChat(chatID string, endTurnID string) (string, error) {
	urlStr := fmt.Sprintf("https://neo.character.ai/chat/%s/copy", chatID)
	headers := c.GetHeaders(false)

	payload := CopyChatRequest{
		EndTurnID: endTurnID,
	}
	bodyBytes, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	resp, err := c.Requester.Post(urlStr, headers, bodyBytes)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to copy chat, status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var result CopyChatResponse
	err = json.Unmarshal(body, &result)
	if err != nil {
		return "", err
	}

	return result.NewChatID, nil
}

// UpdatePrimaryCandidate updates the primary candidate of a turn
func (c *Client) UpdatePrimaryCandidate(chatID string, turnID string, candidateID string) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	// Initialize WebSocket connection if not connected
	err := c.Requester.InitializeWebSocket()
	if err != nil {
		return err
	}

	requestID := generateUUID()

	// Construct the message
	message := WebSocketMessage{
		Command:   "update_primary_candidate",
		OriginID:  "web-next",
		RequestID: requestID,
		Payload: UpdatePrimaryCandidatePayload{
			CandidateID: candidateID,
			TurnKey: TurnKey{
				ChatID: chatID,
				TurnID: turnID,
			},
		},
	}

	// Send the message
	err = c.Requester.SendWebSocketMessage(message)
	if err != nil {
		return err
	}

	// Receive response
	for {
		responseBytes, err := c.Requester.ReceiveRawWebSocketMessage()
		if err != nil {
			return err
		}

		var response WebSocketResponse
		err = json.Unmarshal(responseBytes, &response)
		if err != nil {
			return err
		}

		switch response.Command {
		case "neo_error":
			return errors.New(response.Comment)
		case "ok":
			return nil
		}
	}
}

// EditMessage edits a message in a turn
func (c *Client) EditMessage(chatID string, turnID string, candidateID string, text string) (*Turn, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	// Initialize WebSocket connection if not connected
	err := c.Requester.InitializeWebSocket()
	if err != nil {
		return nil, err
	}

	requestID := generateUUID()

	// Construct the message
	message := WebSocketMessage{
		Command:   "edit_turn_candidate",
		OriginID:  "web-next",
		RequestID: requestID,
		Payload: EditTurnCandidatePayload{
			TurnKey: TurnKey{
				ChatID: chatID,
				TurnID: turnID,
			},
			CurrentCandidateID:     candidateID,
			NewCandidateRawContent: text,
		},
	}

	// Send the message
	err = c.Requester.SendWebSocketMessage(message)
	if err != nil {
		return nil, err
	}

	// Receive response
	for {
		responseBytes, err := c.Requester.ReceiveRawWebSocketMessage()
		if err != nil {
			return nil, err
		}

		var response WebSocketResponse
		err = json.Unmarshal(responseBytes, &response)
		if err != nil {
			return nil, err
		}

		switch response.Command {
		case "neo_error":
			return nil, errors.New(response.Comment)
		case "update_turn":
			var payload TurnResponsePayload
			err = json.Unmarshal(response.Payload, &payload)
			if err != nil {
				return nil, err
			}
			return &payload.Turn, nil
		}
	}
}

// DeleteMessages deletes messages from a chat
func (c *Client) DeleteMessages(chatID string, turnIDs []string) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	// Initialize WebSocket connection if not connected
	err := c.Requester.InitializeWebSocket()
	if err != nil {
		return err
	}

	requestID := generateUUID()

	// Construct the message
	message := WebSocketMessage{
		Command:   "remove_turns",
		OriginID:  "web-next",
		RequestID: requestID,
		Payload: RemoveTurnsPayload{
			ChatID:  chatID,
			TurnIDs: turnIDs,
		},
	}

	// Send the message
	err = c.Requester.SendWebSocketMessage(message)
	if err != nil {
		return err
	}

	// Receive response
	for {
		responseBytes, err := c.Requester.ReceiveRawWebSocketMessage()
		if err != nil {
			return err
		}

		var response WebSocketResponse
		err = json.Unmarshal(responseBytes, &response)
		if err != nil {
			return err
		}

		switch response.Command {
		case "neo_error":
			return errors.New(response.Comment)
		case "remove_turns_response":
			return nil
		}
	}
}

func (c *Client) DeleteMessage(chatID string, turnID string) error {
	return c.DeleteMessages(chatID, []string{turnID})
}

// PinMessage pins a message in a chat
func (c *Client) PinMessage(chatID string, turnID string) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	// Initialize WebSocket connection if not connected
	err := c.Requester.InitializeWebSocket()
	if err != nil {
		return err
	}

	requestID := generateUUID()

	// Construct the message
	message := WebSocketMessage{
		Command:   "set_turn_pin",
		OriginID:  "web-next",
		RequestID: requestID,
		Payload: SetTurnPinPayload{
			IsPinned: true,
			TurnKey: TurnKey{
				ChatID: chatID,
				TurnID: turnID,
			},
		},
	}

	// Send the message
	err = c.Requester.SendWebSocketMessage(message)
	if err != nil {
		return err
	}

	// Receive response
	for {
		responseBytes, err := c.Requester.ReceiveRawWebSocketMessage()
		if err != nil {
			return err
		}

		var response WebSocketResponse
		err = json.Unmarshal(responseBytes, &response)
		if err != nil {
			return err
		}

		switch response.Command {
		case "neo_error":
			return errors.New(response.Comment)
		case "update_turn":
			var payload TurnResponsePayload
			err = json.Unmarshal(response.Payload, &payload)
			if err != nil {
				return err
			}
			if payload.Turn.IsPinned {
				return nil
			}
			return errors.New("failed to pin message")
		}
	}
}

// UnpinMessage unpins a message in a chat
func (c *Client) UnpinMessage(chatID string, turnID string) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	// Initialize WebSocket connection if not connected
	err := c.Requester.InitializeWebSocket()
	if err != nil {
		return err
	}

	requestID := generateUUID()

	// Construct the message
	message := WebSocketMessage{
		Command:   "set_turn_pin",
		OriginID:  "web-next",
		RequestID: requestID,
		Payload: SetTurnPinPayload{
			IsPinned: false,
			TurnKey: TurnKey{
				ChatID: chatID,
				TurnID: turnID,
			},
		},
	}

	// Send the message
	err = c.Requester.SendWebSocketMessage(message)
	if err != nil {
		return err
	}

	// Receive response
	for {
		responseBytes, err := c.Requester.ReceiveRawWebSocketMessage()
		if err != nil {
			return err
		}

		var response WebSocketResponse
		err = json.Unmarshal(responseBytes, &response)
		if err != nil {
			return err
		}

		switch response.Command {
		case "neo_error":
			return errors.New(response.Comment)
		case "update_turn":
			var payload TurnResponsePayload
			err = json.Unmarshal(response.Payload, &payload)
			if err != nil {
				return err
			}
			if !payload.Turn.IsPinned {
				return nil
			}
			return errors.New("failed to unpin message")
		}
	}
}

func (c *Client) AnotherResponse(characterID, chatID, turnID string) (*Turn, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	// Initialize WebSocket connection if not connected
	err := c.Requester.InitializeWebSocket()
	if err != nil {
		return nil, err
	}

	requestID := generateUUID()

	// Construct the message
	message := WebSocketMessage{
		Command:   "generate_turn_candidate",
		OriginID:  "web-next",
		RequestID: requestID,
		Payload: GenerateTurnCandidatePayload{
			CharacterID:         characterID,
			TTSEnabled:          false,
			PreviousAnnotations: generatePreviousAnnotations(),
			SelectedLanguage:    "",
			UserName:            "",
			TurnKey: TurnKey{
				ChatID: chatID,
				TurnID: turnID,
			},
		},
	}

	// Send the message
	err = c.Requester.SendWebSocketMessage(message)
	if err != nil {
		return nil, err
	}

	// Receive response
	for {
		responseBytes, err := c.Requester.ReceiveRawWebSocketMessage()
		if err != nil {
			return nil, err
		}

		var response WebSocketResponse
		err = json.Unmarshal(responseBytes, &response)
		if err != nil {
			return nil, err
		}

		switch response.Command {
		case "neo_error":
			return nil, errors.New(response.Comment)
		case "update_turn":
			var payload TurnResponsePayload
			err = json.Unmarshal(response.Payload, &payload)
			if err != nil {
				return nil, err
			}
			return &payload.Turn, nil
		}
	}
}
