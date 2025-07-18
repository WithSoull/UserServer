package model

import "time"

type Role int32

const (
	ROLE_USER  Role = 0
	ROLE_ADMIN Role = 1
)

type UserInfo struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  Role   `json:"role"`
}

type User struct {
	Id        int64     `json:"id"`
	UserInfo  UserInfo  `json:"user_info"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
