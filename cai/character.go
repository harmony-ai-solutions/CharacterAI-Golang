package cai

import (
	"fmt"
	http "github.com/bogdanfinn/fhttp"
)

type Character struct {
	Token   string
	Session *Session
}

func (c *Character) Create(greeting, identifier, name, avatarRelPath, baseImgPrompt, definition, description, title, visibility, token string, categories []string, copyable, imgGenEnabled bool, extraArgs map[string]interface{}) (map[string]interface{}, error) {
	data := map[string]interface{}{
		"greeting":        greeting,
		"identifier":      identifier,
		"name":            name,
		"avatar_rel_path": avatarRelPath,
		"base_img_prompt": baseImgPrompt,
		"categories":      categories,
		"copyable":        copyable,
		"definition":      definition,
		"description":     description,
		"img_gen_enabled": imgGenEnabled,
		"title":           title,
		"visibility":      visibility,
	}
	for key, val := range extraArgs {
		data[key] = val
	}
	return request("../chat/character/create/", c.Session, token, http.MethodPost, data, false, false)
}

func (c *Character) Update(externalID, greeting, identifier, name, title, definition, description, visibility, token string, categories []string, copyable bool, extraArgs map[string]interface{}) (map[string]interface{}, error) {
	data := map[string]interface{}{
		"external_id": externalID,
		"name":        name,
		"categories":  categories,
		"title":       title,
		"visibility":  visibility,
		"copyable":    copyable,
		"description": description,
		"greeting":    greeting,
		"definition":  definition,
	}
	for key, val := range extraArgs {
		data[key] = val
	}
	return request("../chat/character/update/", c.Session, token, http.MethodPost, data, false, false)
}

func (c *Character) Trending() (map[string]interface{}, error) {
	return request("chat/characters/trending/", c.Session, "", http.MethodGet, nil, false, false)
}

func (c *Character) Recommended(token string) (map[string]interface{}, error) {
	return request("chat/characters/recommended/", c.Session, token, http.MethodGet, nil, false, false)
}

func (c *Character) Categories() (map[string]interface{}, error) {
	return request("chat/character/categories/", c.Session, "", http.MethodGet, nil, false, false)
}

func (c *Character) Info(char, token string) (map[string]interface{}, error) {
	data := map[string]interface{}{"external_id": char}
	return request("chat/character/", c.Session, token, http.MethodPost, data, false, false)
}

func (c *Character) Search(query, token string) (map[string]interface{}, error) {
	url := fmt.Sprintf("chat/characters/search/?query=%s/", query)
	return request(url, c.Session, token, http.MethodGet, nil, false, false)
}

func (c *Character) Voices() (map[string]interface{}, error) {
	return request("chat/character/voices/", c.Session, "", http.MethodGet, nil, false, false)
}
