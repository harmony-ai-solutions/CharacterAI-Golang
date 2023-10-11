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
		"Authorization": {fmt.Sprintf("Token %v", token)},
		"Content-Type":  {"application/json"},
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

	// Check for error
	if errRequest != nil {
		return nil, errRequest
	}

	bodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	bodyString := string(bodyBytes)

	if split {
		splitted := strings.Split(bodyString, "\n")
		if len(splitted) >= 2 {
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