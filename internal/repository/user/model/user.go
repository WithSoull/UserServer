package model

import "time"

type Role int32

type UserInfo struct {
	Name string
	Email string
	Role string
}

type User struct {
	Id int64
	UserInfo UserInfo
	CreatedAt time.Time
	UpdatedAt time.Time
}

