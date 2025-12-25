package models

import "time"

const (
	TODO_FILE = "todos.json"
)

type Todo struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	IsCompleted bool      `json:"is_completed"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func NewTodo(id int, title, description string) *Todo {
	return &Todo{ID: id,
		Title:       title,
		Description: description,
		IsCompleted: false,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now()}
}
