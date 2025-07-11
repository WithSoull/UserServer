package model

import "time"

type Role int32

const (
	ROLE_USER Role = 0
	ROLE_ADMIN Role = 1
)

type UserInfo struct {
	Name string
	Email string
	Role Role
}

type User struct {
	Id int64
	UserInfo UserInfo
	CreatedAt time.Time
	UpdatedAt time.Time
}

