package logger

import (
	"database/sql"
)

type Logger struct {
	db *sql.DB
}

//NewLogger - конструктор логгера
func NewLogger(db *sql.DB) *Logger {
	return &Logger{db}
}

//Log добавляет в базу данных лог с текущим unix временем
func (l *Logger) Log(text string, ip string) error {
	_, err := l.db.Exec(`insert into main_log ("text", "time", "ip") values($1, NOW(), $2)`, text, ip)
	return err
}
