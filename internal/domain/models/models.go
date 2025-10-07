package models

import "time"

type User struct {
	UID         string `json:"uid"`
	Name        string `json:"name" validate:"required"`
	Email       string `json:"email" validate:"required,email"`
	Pass        string `json:"pass" validate:"required"`
	DeletedUser bool   `json:"deleted_user"`
}

type Book struct {
	BID       string    `json:"b_id"`
	Label     string    `json:"label" validate:"required"`
	Author    string    `json:"author" validate:"required"`
	Deleted   bool      `json:"delete"`
	UserUid   string    `json:"user_uid" validate:"required"`
	CreatedAt time.Time `json:"created_at"`
}
