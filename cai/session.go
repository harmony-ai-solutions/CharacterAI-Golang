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
	"bytes"
	http "github.com/bogdanfinn/fhttp"
	tls_client "github.com/bogdanfinn/tls-client"
	"github.com/bogdanfinn/tls-client/profiles"
)

type Session struct {
	// Public Attributes
	URL   string
	Token string
	// Internal
	client tls_client.HttpClient
}

func NewSession(url, token string, clientIdentifier profiles.ClientProfile) (*Session, error) {
	jar := tls_client.NewCookieJar()
	options := []tls_client.HttpClientOption{
		tls_client.WithTimeoutSeconds(30),
		tls_client.WithClientProfile(clientIdentifier),
		tls_client.WithNotFollowRedirects(),
		tls_client.WithCookieJar(jar), // create cookieJar instance and pass it as argument
	}

	// Init TLS Client
	client, errClient := tls_client.NewHttpClient(tls_client.NewNoopLogger(), options...)
	if errClient != nil {
		return nil, errClient
	}

	return &Session{
		URL:    url,
		Token:  token,
		client: client,
	}, nil
}

func (s *Session) GET(url string, headers http.Header) (*http.Response, error) {
	req, errReq := http.NewRequest(http.MethodGet, url, nil)
	if errReq != nil {
		return nil, errReq
	}
	req.Header = headers
	resp, errExecute := s.client.Do(req)
	if errExecute != nil {
		return nil, errExecute
	}
	return resp, nil
}

func (s *Session) POST(url string, headers http.Header, data []byte) (*http.Response, error) {
	req, errReq := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(data))
	if errReq != nil {
		return nil, errReq
	}
	req.Header = headers
	resp, errExecute := s.client.Do(req)
	if errExecute != nil {
		return nil, errExecute
	}
	return resp, nil
}

func (s *Session) PUT(url string, headers http.Header, data []byte) (*http.Response, error) {
	req, errReq := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(data))
	if errReq != nil {
		return nil, errReq
	}
	req.Header = headers
	resp, errExecute := s.client.Do(req)
	if errExecute != nil {
		return nil, errExecute
	}
	return resp, nil
}
