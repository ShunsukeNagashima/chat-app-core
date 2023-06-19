package model

import "fmt"

type RoomType string

const (
	Public  RoomType = "public"
	Private RoomType = "private"
)

type Room struct {
	RoomID   string   `json:"room_id"`
	Name     string   `json:"name"`
	RoomType RoomType `json:"room_type"`
}

func ParseRoomType(s string) (RoomType, error) {
	switch s {
	case string(Private):
		return Private, nil
	case string(Public):
		return Public, nil
	default:
		return "", fmt.Errorf("invalid RoomType: %s", s)
	}
}
