package model

import "time"

type Message struct {
	MessageID string    `json:"messageId"`
	Content   string    `json:"content"`
	SenderID  string    `json:"senderId"`
	RoomID    string    `json:"roomId"`
	CreatedAt time.Time `json:"createdAt"`
}
