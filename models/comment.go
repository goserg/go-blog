package models

type Comment struct {
	ID         int64
	Time       string
	Text       string
	Author     User
	RefPost    string
	RefComemnt string
}
