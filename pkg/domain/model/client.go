package model

import (
	"encoding/json"
	"log"

	"github.com/gorilla/websocket"
)

type Client struct {
	Conn *websocket.Conn
	Send chan Event
	Hub  Hub
}

func NewClient(ws *websocket.Conn, hub Hub) *Client {
	return &Client{
		Conn: ws,
		Send: make(chan Event),
		Hub:  hub,
	}
}

func (c *Client) Read() {
	defer func() {
		c.disconnect()
	}()

	for {
		var rawEvent RawEvent
		err := c.Conn.ReadJSON(&rawEvent)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("unexpected close error: %v", err)
			}
			break
		}

		switch rawEvent.Type {
		case "message":
			var message Message
			err := json.Unmarshal(rawEvent.Data, &message)
			if err != nil {
				log.Printf("Failed to unmarshal message: %v", err)
				break
			}
			c.Hub.BroadcastEvent(&message)
		case "roomEvent":
			var roomEvent RoomEvent
			err := json.Unmarshal(rawEvent.Data, &roomEvent)
			if err != nil {
				log.Printf("Failed to unmarshal roomEvent: %v", err)
				break
			}
			c.Hub.BroadcastEvent(&roomEvent)
		default:
			log.Printf("Invalid event type: %s", rawEvent.Type)
		}
	}
}

func (c *Client) Write() {
	defer func() {
		c.disconnect()
	}()

	for {
		event, ok := <-c.Send
		if !ok {
			c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
			return
		}

		var eventType string
		var err error
		switch e := event.(type) {
		case *Message:
			eventType = "message"
		case *RoomEvent:
			eventType = "roomEvent"
		default:
			log.Printf("Invalid event type: %v", e)
			continue
		}

		data, err := json.Marshal(event)
		if err != nil {
			log.Printf("Failed to marshal event: %v", err)
			return
		}

		rawEvent := RawEvent{
			Type: eventType,
			Data: data,
		}

		err = c.Conn.WriteJSON(rawEvent)
		if err != nil {
			log.Printf("Failed to write event: %v", err)
			return
		}
	}
}

func (c *Client) disconnect() {
	c.Hub.UnregisterClient(c)
	log.Printf("Unregistered client: %v", c)

	err := c.Conn.Close()
	if err != nil {
		log.Printf("Failed to close connection for client: %v, error: %v", c, err)
	}
}
