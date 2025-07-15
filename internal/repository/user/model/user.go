package model

import "time"

type Role int32

type UserInfo struct {
	Name  string `db:"name"`
	Email string `db:"email"`
	Role  string `db:"role"`
}

type User struct {
	Id        int64     `db:"id"`
	UserInfo  UserInfo  `db:""`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
