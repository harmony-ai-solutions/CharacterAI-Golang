// Package cai
/*
Copyright © 2023 Harmony AI Solutions & Contributors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cai

import (
	"encoding/json"
	"errors"
	"fmt"
	http "github.com/bogdanfinn/fhttp"
	"strings"
)

type Chat struct {
	Token   string
	Session *Session
}

type ChatData struct {
	Title           string             `json:"title"`
	Participants    []*ChatParticipant `json:"participants"`
	ExternalID      string             `json:"external_id"`
	Created         string             `json:"created"`
	LastInteraction string             `json:"last_interaction"`
	Type            string             `json:"type"`
	Description     string             `json:"description"`
}

type ChatUser struct {
	Username  string       `json:"username"`
	ID        uint64       `json:"id"`
	FirstName string       `json:"first_name"`
	IsStaff   bool         `json:"is_staff"`
	Account   *ChatAccount `json:"account"`
}

type ChatAccount struct {
	Name                     string `json:"name"`
	AvatarType               string `json:"avatar_type"`
	OnboardingComplete       bool   `json:"onboarding_complete"`
	AvatarFileName           string `json:"avatar_file_name"`
	MobileOnboardingComplete *bool  `json:"mobile_onboarding_complete"`
}

type ChatParticipant struct {
	User            *ChatUser `json:"user"`
	Name            string    `json:"name"`
	IsHuman         bool      `json:"is_human"`
	NumInteractions int       `json:"num_interactions"`
}

type ChatChar struct {
	Participant    *ChatParticipant `json:"participant"`
	AvatarFileName *string          `json:"avatar_file_name"`
}

type ChatHistory struct {
	Messages []*ChatHistoryMessage `json:"messages"`
	HasMore  bool                  `json:"has_more"`
	NextPage int                   `json:"next_page"`
}

type ChatHistoryMessage struct {
	Target                        string    `json:"tgt"`
	ImageRelPath                  string    `json:"image_rel_path"`
	ImagePromptText               string    `json:"image_prompt_text"`
	Deleted                       *bool     `json:"deleted"` // TODO: confirm bool
	IsAlternative                 bool      `json:"is_alternative"`
	SourceName                    string    `json:"src__name"`
	SourceUserUsername            string    `json:"src__user__username"`
	ResponsibleUserUsername       *string   `json:"responsible_user__username"`
	ID                            uint64    `json:"id"`
	UUID                          string    `json:"uuid"`
	Source                        string    `json:"src"`
	SourceIsHuman                 bool      `json:"src__is_human"`
	SourceCharacterAvatarFileName *string   `json:"src__character__avatar_file_name"`
	SourceChar                    *ChatChar `json:"src_char"`
	Text                          string    `json:"text"`
}

type ChatMessage struct {
	Replies         []*ChatMessageReply `json:"replies"`
	SourceChar      *ChatChar           `json:"src_char"`
	IsFinalChunk    bool                `json:"is_final_chunk"`
	LastUserMsgId   uint64              `json:"last_user_msg_id"`
	LastUserMsgUUID string              `json:"last_user_msg_uuid"`
}

type ChatMessageReply struct {
	Text string `json:"text"`
	UUID string `json:"uuid"`
	ID   uint64 `json:"id"`
}

func (c *Chat) CreateRoom(characters []string, name, topic string, extraArgs map[string]interface{}) (map[string]interface{}, error) {
	data := map[string]interface{}{
		"characters": characters,
		"name":       name,
		"topic":      topic,
		"visibility": "PRIVATE",
	}
	for key, val := range extraArgs {
		data[key] = val
	}
	return request("chat/room/create/", c.Session, c.Token, http.MethodPost, data, false, false)
}

func (c *Chat) Rate(rate int, historyID, messageID string, extraArgs map[string]interface{}) (map[string]interface{}, error) {
	var label []int
	switch rate {
	case 0:
		label = []int{234, 238, 241, 244}
	case 1:
		label = []int{235, 237, 241, 244}
	case 2:
		label = []int{235, 238, 240, 244}
	case 3:
		label = []int{235, 238, 241, 243}
	default:
		return nil, errors.New("Wrong Rate Value")
	}
	data := map[string]interface{}{
		"label_ids":           label,
		"history_external_id": historyID,
		"message_uuid":        messageID,
	}
	for key, val := range extraArgs {
		data[key] = val
	}
	return request("chat/annotations/label/", c.Session, c.Token, "PUT", data, false, false)
}

func (c *Chat) NextMessage(historyID, parentMsgUUID, tgt string, extraArgs map[string]interface{}) (map[string]interface{}, error) {
	data := map[string]interface{}{
		"history_external_id": historyID,
		"parent_msg_uuid":     parentMsgUUID,
		"tgt":                 tgt,
	}
	for key, val := range extraArgs {
		data[key] = val
	}
	return request("chat/streaming/", c.Session, c.Token, http.MethodPost, data, false, false)
}

func (c *Chat) GetHistories(char string, number int) (map[string]interface{}, error) {
	data := map[string]interface{}{
		"external_id": char,
		"number":      number,
	}
	return request("chat/character/histories_v2/", c.Session, c.Token, http.MethodPost, data, false, false)
}

func (c *Chat) GetHistory(historyID string) (*ChatHistory, error) {
	url := fmt.Sprintf("chat/history/msgs/user/?history_external_id=%s", historyID)
	historyResult, errHistory := baseRequest(url, c.Session, c.Token, http.MethodGet, nil, false)
	if errHistory != nil {
		return nil, errHistory
	}
	history := &ChatHistory{}
	if errUnmarshall := json.Unmarshal(historyResult, history); errUnmarshall != nil {
		return nil, errUnmarshall
	}
	return history, nil
}

func (c *Chat) GetChat(char string) (*ChatData, error) {
	data := map[string]interface{}{
		"character_external_id": char,
	}
	chatResult, errChat := baseRequest("chat/history/continue/", c.Session, c.Token, http.MethodPost, data, false)
	if errChat != nil {
		return nil, errChat
	}
	chatData := &ChatData{}
	if errUnmarshall := json.Unmarshal(chatResult, chatData); errUnmarshall != nil {
		return nil, errUnmarshall
	}
	return chatData, nil
}

func (c *Chat) SendMessage(historyID, tgt, text string, extraArgs map[string]interface{}) (*ChatMessage, error) {
	data := map[string]interface{}{
		"history_external_id": historyID,
		"tgt":                 tgt,
		"text":                text,
	}
	if extraArgs != nil {
		for key, val := range extraArgs {
			data[key] = val
		}
	}
	chatMessageResult, errChat := baseRequest("chat/streaming/", c.Session, c.Token, http.MethodPost, data, false)
	if errChat != nil {
		return nil, errChat
	}
	chatMessageString := string(chatMessageResult)
	chatMessageString = strings.ReplaceAll(chatMessageString, "\r\n", "\n")
	splitted := strings.Split(chatMessageString, "\n")
	if len(splitted) > 1 {
		chatMessageString = splitted[len(splitted)-2]
	}
	chatMessage := &ChatMessage{}
	if errUnmarshall := json.Unmarshal([]byte(chatMessageString), chatMessage); errUnmarshall != nil {
		return nil, errUnmarshall
	}
	return chatMessage, nil
}

func (c *Chat) DeleteMessage(historyID string, uuidsToDelete []string, extraArgs map[string]interface{}) (map[string]interface{}, error) {
	data := map[string]interface{}{
		"history_id":      historyID,
		"uuids_to_delete": uuidsToDelete,
	}
	if extraArgs != nil {
		for key, val := range extraArgs {
			data[key] = val
		}
	}
	return request("chat/history/msgs/delete/", c.Session, c.Token, http.MethodPost, data, false, false)
}

func (c *Chat) NewChat(char string) (*ChatData, error) {
	data := map[string]interface{}{
		"character_external_id": char,
	}
	chatResult, errChat := baseRequest("chat/history/create/", c.Session, c.Token, http.MethodPost, data, false)
	if errChat != nil {
		return nil, errChat
	}
	chatData := &ChatData{}
	if errUnmarshall := json.Unmarshal(chatResult, chatData); errUnmarshall != nil {
		return nil, errUnmarshall
	}
	return chatData, nil
}
