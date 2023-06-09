package model

import "time"

type Message struct {
	MessageID string    `json:"message_id"`
	Content   string    `json:"content"`
	SenderID  string    `json:"sender_id"`
	RoomID    string    `json:"room_id"`
	CreatedAt time.Time `json:"created_at"`
}
