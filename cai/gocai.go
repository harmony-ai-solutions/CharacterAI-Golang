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
	"github.com/bogdanfinn/tls-client/profiles"
	"io"
	"strings"
)

type GoCAI struct {
	Token     string
	Session   *Session // Assuming Session is a type from tls_client or its Go equivalent
	User      *User    // Assuming these types are defined elsewhere or will be
	Post      *Post
	Character *Character
	Chat      *Chat
	Chat2     *Chat2
}

func NewGoCAI(token string, plus bool) (*GoCAI, error) {
	// Determine subdomain to use
	subdomain := "beta"
	if plus {
		subdomain = "plus"
	}

	// Init TLSClient session
	session, errSession := NewSession(
		fmt.Sprintf("https://%v.character.ai/", subdomain),
		token,
		profiles.Chrome_112,
	)
	if errSession != nil {
		return nil, errSession
	}

	// Init Data Objects
	user := &User{Token: token, Session: session}
	post := &Post{Token: token, Session: session}
	character := &Character{Token: token, Session: session}
	chat := &Chat{Token: token, Session: session}
	chat2 := &Chat2{Token: token, Session: session}

	return &GoCAI{
		Token:     token,
		Session:   session,
		User:      user,
		Post:      post,
		Character: character,
		Chat:      chat,
		Chat2:     chat2,
	}, nil
}

func request(url string, session *Session, token string, method string, data map[string]interface{}, split bool, neo bool) (map[string]interface{}, error) {
	var link string
	if neo {
		link = "https://neo.character.ai/" + url
	} else {
		link = session.URL + url
	}

	if token == "" {
		token = session.Token
	}

	// Create Headers
	header := http.Header{
		"accept":          {"application/json", "text/plain", "*/*"},
		"accept-encoding": {"gzip", "deflate", "br"},
		"accept-language": {"de-DE,de;q=0.9,en-US;q=0.8,en;q=0.7"},
		"authorization":   {fmt.Sprintf("Token %v", token)},
		"connection":      {"keep-alive"},
		"content-type":    {"application/json"},
		"user-agent":      {"tls-client/0.2.2"},
		http.HeaderOrderKey: {
			"accept",
			"accept-encoding",
			"accept-language",
			"authorization",
			"connection",
			"content-type",
			"user-agent",
		},
	}

	// Send request
	var response *http.Response
	var errRequest error
	switch method {
	case http.MethodGet:
		response, errRequest = session.GET(link, header)
	case http.MethodPost:
		dataBytes, errEncode := json.Marshal(data)
		if errEncode != nil {
			return nil, errEncode
		}
		response, errRequest = session.POST(link, header, dataBytes)
	case http.MethodPut:
		dataBytes, errEncode := json.Marshal(data)
		if errEncode != nil {
			return nil, errEncode
		}
		response, errRequest = session.PUT(link, header, dataBytes)
	default:
		return nil, fmt.Errorf("method %v is currently not supported", method)
	}
	defer response.Body.Close()

	// Check for error
	if errRequest != nil {
		return nil, errRequest
	}

	// Check for Error codes
	if response.StatusCode != 200 {
		return nil, errors.New(response.Status)
	}

	bodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	bodyString := string(bodyBytes)

	// Set this to true if multiline output is to be expected.
	if split {
		// Ensure windows newlines are not an issue
		bodyString = strings.ReplaceAll(bodyString, "\r\n", "\n")
		splitted := strings.Split(bodyString, "\n")
		if len(splitted) > 1 {
			bodyString = splitted[len(splitted)-2]
		}
	}

	var responseData map[string]interface{}
	err = json.Unmarshal([]byte(bodyString), &responseData)
	if err != nil {
		return nil, err
	}

	if command, exists := responseData["command"]; exists && command == "neo_error" {
		return nil, errors.New(fmt.Sprintf("ServerError: %s", responseData["comment"]))
	}

	if detail, exists := responseData["detail"]; exists && detail == "Auth" {
		return nil, errors.New("invalid token")
	}

	if status, exists := responseData["status"]; exists && strings.HasPrefix(status.(string), "Error") {
		return nil, errors.New(fmt.Sprintf("ServerError: %s", status))
	}

	if errorDetails, exists := responseData["error"]; exists {
		return nil, errors.New(fmt.Sprintf("ServerError: %s", errorDetails))
	}

	return responseData, nil
}

func (g *GoCAI) Ping() (map[string]interface{}, error) {
	return request("ping", g.Session, g.Token, http.MethodGet, nil, false, true)
}
