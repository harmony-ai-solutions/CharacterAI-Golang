package cai

import (
	"errors"
	"fmt"
	http "github.com/bogdanfinn/fhttp"
)

type Chat struct {
	Token   string
	Session *Session
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
	return request("../chat/room/create/", c.Session, c.Token, http.MethodPost, data, false, false)
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

func (c *Chat) GetHistory(historyID string) (map[string]interface{}, error) {
	url := fmt.Sprintf("chat/history/msgs/user/?history_external_id=%s", historyID)
	return request(url, c.Session, c.Token, http.MethodGet, nil, false, false)
}

func (c *Chat) GetChat(char string) (map[string]interface{}, error) {
	data := map[string]interface{}{
		"character_external_id": char,
	}
	return request("chat/history/continue/", c.Session, c.Token, http.MethodPost, data, false, false)
}

func (c *Chat) SendMessage(historyID, tgt, text string, extraArgs map[string]interface{}) (map[string]interface{}, error) {
	data := map[string]interface{}{
		"history_external_id": historyID,
		"tgt":                 tgt,
		"text":                text,
	}
	for key, val := range extraArgs {
		data[key] = val
	}
	return request("chat/streaming/", c.Session, c.Token, http.MethodPost, data, false, false)
}

func (c *Chat) DeleteMessage(historyID string, uuidsToDelete []string, extraArgs map[string]interface{}) (map[string]interface{}, error) {
	data := map[string]interface{}{
		"history_id":      historyID,
		"uuids_to_delete": uuidsToDelete,
	}
	for key, val := range extraArgs {
		data[key] = val
	}
	return request("chat/history/msgs/delete/", c.Session, c.Token, http.MethodPost, data, false, false)
}

func (c *Chat) NewChat(char string) (map[string]interface{}, error) {
	data := map[string]interface{}{
		"character_external_id": char,
	}
	return request("chat/history/create/", c.Session, c.Token, http.MethodPost, data, false, false)
}
