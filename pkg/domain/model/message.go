package model

import "time"

type Message struct {
	MessageID string    `json:"messageId"`
	Content   string    `json:"content"`
	UserID    string    `json:"userId"`
	RoomID    string    `json:"roomId"`
	CreatedAt time.Time `json:"createdAt"`
}
