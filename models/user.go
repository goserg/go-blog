package models

type User struct {
	ID           uint64
	UserName     string
	PasswordHash string
	Authorized   bool
}
