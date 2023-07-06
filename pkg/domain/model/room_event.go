package model

import "fmt"

type EventType string

const (
	RoomCrated  EventType = "ROOM_CREATED"
	RoomDeleted EventType = "ROOM_DELETED"
)

type RoomEvent struct {
	RoomID    string    `json:"roomId"`
	EventType EventType `json:"eventType"`
}

func ParseEventType(s string) (EventType, error) {
	switch s {
	case string(RoomCrated):
		return RoomCrated, nil
	case string(RoomDeleted):
		return RoomDeleted, nil
	default:
		return "", fmt.Errorf("invalid EventType: %s", s)
	}
}
