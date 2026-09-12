package models

type Todo struct {
	ID uint64 `json:"id"`
	Task string `json:"task"`
	Done bool `json:"done"`
}
