package model

import (
	"time"
)

type Todo struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Completed bool      `json:"completed"`
	CreatedAt time.Time `json:"created_at"`
}

type DBHandler interface {
	GetTodos(sessionId string) []*Todo
	AddTodo(sessionId, name string) *Todo
	DeleteTodo(id int) bool
	CompleteTodo(id int, complete bool) bool
	Close()
}

func NewDBHandler(dbConn string) DBHandler {
	//handler = newMemoryHandler()
	//return newSqliteHandler(filePath)
	return newPqHandler(dbConn)
}
