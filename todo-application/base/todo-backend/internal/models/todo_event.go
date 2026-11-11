package models

import (
	"crypto/rand"
	"encoding/hex"
	"time"

)

const (
	EventTypeCreated = "created"
	EventTypeCompleted = "completed"
	SubjectTodoCreated   = "todos.created"
	SubjectTodoCompleted = "todos.completed"
)

type TodoEvent struct {
	ID        string    `json:"id"`
	Type      string      `json:"type"` 
	Todo      Todo `json:"todo"`
	Timestamp time.Time   `json:"timestamp"`
}

func NewTodoEvent(eventType string, todo Todo) TodoEvent {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	return TodoEvent {
		ID: hex.EncodeToString(b),
		Type: eventType, 
		Todo: todo, 
		Timestamp: time.Now().UTC(),
	}
}

func (e TodoEvent) Subject() string {
	return "todos." + e.Type
}