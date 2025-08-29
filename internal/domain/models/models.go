package models

import "time"

type User struct {
	UID   string `json:"uid"`
	Name  string `json:"name" validate:"required"`
	Email string `json:"email" validate:"required,email"`
	Pass  string `json:"pass" validate:"required"`
}

type Book struct {
	BID        string    `json:"b_id"`
	Label      string    `json:"label" validate:"required"`
	Author     string    `json:"author" validate:"required"`
	Deleted    bool      `json:"delete"`
	User_UID   string    `json:"user_uid" validate:"required"`
	Created_at time.Time `json:"created_at"`
}
