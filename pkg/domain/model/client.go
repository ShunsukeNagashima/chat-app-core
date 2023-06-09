package model

import (
	"log"

	"github.com/gorilla/websocket"
)

type Client struct {
	Conn *websocket.Conn
	Send chan *Message
	Hub  *Hub
}

func NewClient(ws *websocket.Conn, hub *Hub) *Client {
	return &Client{
		Conn: ws,
		Send: make(chan *Message),
		Hub:  hub,
	}
}

func (c *Client) Read() {
	defer func() {
		c.disconnect()
	}()

	for {
		var message Message
		err := c.Conn.ReadJSON(&message)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("unexpected close error: %v", err)
			}
			break
		}

		c.Hub.Broadcast <- &message
	}
}

func (c *Client) Write() {
	defer func() {
		c.disconnect()
	}()

	for {
		message, ok := <-c.Send
		if !ok {
			c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
			return
		}

		c.Conn.WriteJSON(message)
	}
}

func (c *Client) disconnect() {
	select {
	case c.Hub.Unregister <- c:
		log.Printf("Unregistered client: %v", c)
	default:
		log.Printf("Failed to unregister client: %v", c)
	}

	err := c.Conn.Close()
	if err != nil {
		log.Printf("Failed to close connection for client: %v, error: %v", c, err)
	}
}
