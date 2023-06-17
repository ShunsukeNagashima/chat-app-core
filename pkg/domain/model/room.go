package model

type RoomType string

const (
	Public  RoomType = "Public"
	Private RoomType = "Private"
)

type Room struct {
	RoomID   string   `json:"room_id"`
	Name     string   `json:"name"`
	RoomType RoomType `json:"room_type"`
}
