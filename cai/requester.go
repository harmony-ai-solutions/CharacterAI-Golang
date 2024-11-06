package cai

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
	"net/url"
	"sync"
	"time"
)

// Requester handles HTTP and WebSocket requests
type Requester struct {
	client       *http.Client
	wsConn       *websocket.Conn
	wsMutex      sync.Mutex
	wsURL        url.URL
	wsHeaders    http.Header
	wsConnected  bool
	wsWriteMutex sync.Mutex
	wsReadMutex  sync.Mutex
	ctx          context.Context
	cancel       context.CancelFunc
}

// NewRequester creates a new Requester instance
func NewRequester(token string, proxy string) *Requester {
	transport := &http.Transport{}

	if proxy != "" {
		proxyURL, err := url.Parse(proxy)
		if err == nil {
			transport.Proxy = http.ProxyURL(proxyURL)
		}
	}

	ctx, cancel := context.WithCancel(context.Background())

	return &Requester{
		client: &http.Client{
			Transport: transport,
			Timeout:   30 * time.Second,
		},
		wsHeaders: http.Header{
			"User-Agent": []string{"Mozilla/5.0"},
			"Cookie":     []string{fmt.Sprintf(`HTTP_AUTHORIZATION="Token %s"`, token)},
		},
		wsURL:  url.URL{Scheme: "wss", Host: "neo.character.ai", Path: "/ws/"},
		ctx:    ctx,
		cancel: cancel,
	}
}

// DoRequest performs an HTTP request
func (r *Requester) DoRequest(method, urlStr string, headers map[string]string, body []byte) (*http.Response, error) {
	req, err := http.NewRequest(method, urlStr, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	return r.client.Do(req)
}

// Get performs a GET request
func (r *Requester) Get(urlStr string, headers map[string]string) (*http.Response, error) {
	return r.DoRequest("GET", urlStr, headers, nil)
}

// Post performs a POST request
func (r *Requester) Post(urlStr string, headers map[string]string, body []byte) (*http.Response, error) {
	return r.DoRequest("POST", urlStr, headers, body)
}

// InitializeWebSocket initializes the WebSocket connection
func (r *Requester) InitializeWebSocket() error {
	r.wsMutex.Lock()
	defer r.wsMutex.Unlock()

	if r.wsConnected {
		return nil
	}

	dialer := websocket.DefaultDialer

	conn, _, err := dialer.Dial(r.wsURL.String(), r.wsHeaders)
	if err != nil {
		return err
	}

	r.wsConn = conn
	r.wsConnected = true

	return nil
}

// CloseWebSocket closes the WebSocket connection
func (r *Requester) CloseWebSocket() error {
	r.wsMutex.Lock()
	defer r.wsMutex.Unlock()

	if !r.wsConnected {
		return nil
	}

	err := r.wsConn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	if err != nil {
		return err
	}

	r.cancel()

	err = r.wsConn.Close()
	if err != nil {
		return err
	}

	r.wsConnected = false
	return nil
}

// SendWebSocketMessage sends a message over the WebSocket connection
func (r *Requester) SendWebSocketMessage(message WebSocketMessage) error {
	r.wsWriteMutex.Lock()
	defer r.wsWriteMutex.Unlock()

	if !r.wsConnected {
		return errors.New("WebSocket not connected")
	}

	messageBytes, err := json.Marshal(message)
	if err != nil {
		return err
	}

	return r.wsConn.WriteMessage(websocket.TextMessage, messageBytes)
}

// ReceiveRawWebSocketMessage receives a raw message from the WebSocket connection
func (r *Requester) ReceiveRawWebSocketMessage() ([]byte, error) {
	r.wsReadMutex.Lock()
	defer r.wsReadMutex.Unlock()

	if !r.wsConnected {
		return nil, errors.New("WebSocket not connected")
	}

	_, messageBytes, err := r.wsConn.ReadMessage()
	if err != nil {
		return nil, err
	}

	return messageBytes, nil
}

// ListenWebSocketMessages listens for messages and sends them to a channel
func (r *Requester) ListenWebSocketMessages(messages chan<- []byte) {
	for {
		select {
		case <-r.ctx.Done():
			close(messages)
			return
		default:
			response, err := r.ReceiveRawWebSocketMessage()
			if err != nil {
				// Handle error (log it, close connection, etc.)
				close(messages)
				return
			}
			messages <- response
		}
	}
}
