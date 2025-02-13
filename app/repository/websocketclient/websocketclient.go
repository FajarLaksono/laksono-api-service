package websocketclient

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/gorilla/websocket"
)

type Client struct {
	conn *websocket.Conn
}

// NewClient creates a new WebSocket client and connects to the server
func NewClient(serverURL string) (*Client, error) {
	u, err := url.Parse(serverURL)
	fmt.Printf("Websocket URL: %+v\n", u)
	if err != nil {
		return nil, fmt.Errorf("invalid server URL: %v", err)
	}
	headers := http.Header{}

	conn, _, err := websocket.DefaultDialer.Dial(u.String(), headers)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to WebSocket server: %v", err)
	}

	return &Client{conn: conn}, nil
}

// SendMessage sends a message to the WebSocket server
func (c *Client) SendMessage(message string) error {
	err := c.conn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return fmt.Errorf("failed to send message: %v", err)
	}
	return nil
}

// Close closes the WebSocket connection
func (c *Client) Close() error {
	return c.conn.Close()
}
