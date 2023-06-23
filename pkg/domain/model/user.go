package model

type User struct {
	UserID   string `json:"userID"`
	Username string `json:"userName"`
	Email    string `json:"email"`
}
