package model

import (
	"encoding/json"
)

type Event interface{}

type EventType string

type RawEvent struct {
	Type EventType       `json:"type"`
	Data json.RawMessage `json:"data"`
}

const (
	MessageSent EventType = "MessageSent"
	UserJoined  EventType = "UserJoined"
)

type RoomUserDetails struct {
	Room   Room   `json:"room"`
	UserID string `json:"userId"`
}
