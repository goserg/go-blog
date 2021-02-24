package server

type User struct {
	ID           uint64
	UserName     string
	PasswordHash string
	Authorized   bool
}
