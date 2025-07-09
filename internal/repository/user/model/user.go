package model

import "time"

type Role int32

const (
	ROLE_USER Role = 0
	ROLE_ADMIN Role = 1
)

type User struct {
	Id int64
	Name string
	Email string
	Role Role
	CreatedAt time.Time
	UpdatedAt time.Time
}

