package model

import "time"

type User struct {
	UserID    string    `json:"userId"`
	Username  string    `json:"userName"`
	Email     string    `json:"email"`
	ImageURL  string    `json:"imageUrl"`
	CreatedAt time.Time `json:"createdAt"`
}
