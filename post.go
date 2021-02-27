package server

type Post struct {
	ID     uint64
	Time   string
	Title  string
	Text   string
	Author User

	Units []PostUnit
}
