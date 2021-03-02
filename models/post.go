package models

type Post struct {
	ID     uint64
	Time   string
	Title  string
	Text   string
	Author User

	Units    []PostUnit
	Comments []Comment
}
