package models

type User struct {
	Name  string `json: "name" validate:"required"`
	Emaim string `json: "email" validate:"required, email"`
	Pass  string `json: "pass" validate:"required"`
}

type Book struct {
	BID    string `json: "bid" validate:"required"`
	Lable  string `json: "label" validate:"required"`
	Author string `json: "author" validate:"required"`
}
