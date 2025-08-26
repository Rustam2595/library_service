package models

type User struct {
	UID   string `json:"uid"`
	Name  string `json:"name" validate:"required"`
	Email string `json:"email" validate:"required,email"`
	Pass  string `json:"pass" validate:"required"`
}

type Book struct {
	BID    string `json:"bid"`
	Label  string `json:"label" validate:"required"`
	Author string `json:"author" validate:"required"`
}
