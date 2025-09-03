package model

import "time"

type UserInfo struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type User struct {
	Id        int64     `json:"id"`
	UserInfo  UserInfo  `json:"user_info"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
