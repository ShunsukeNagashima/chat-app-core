package model

type User struct {
	UserID   string `json:"userId"`
	Username string `json:"userName"`
	Email    string `json:"email"`
}
