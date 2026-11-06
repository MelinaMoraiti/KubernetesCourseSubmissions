package models

const SubjectTodoStatus = "todos.status"

// StatusMessage is the "Generic" webhook payload:
// {"user": "bot", "message": "A todo was created"}
type StatusMessage struct {
	User    string `json:"user"`
	Message string `json:"message"`
}
