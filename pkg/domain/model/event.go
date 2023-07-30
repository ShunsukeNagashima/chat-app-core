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
	MessageSent    EventType = "MessageSent"
	RoomUserChange EventType = "RoomUserChange"
)

type RoomUserDetails struct {
	RoomID string `json:"roomId"`
	UserID string `json:"userId"`
}
