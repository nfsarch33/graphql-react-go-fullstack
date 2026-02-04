package models

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// Todo represents a task in the application
type Todo struct {
	ID          int            `gorm:"primaryKey" json:"id"`
	Title       string         `json:"title"`
	Description string         `json:"description"`
	Completed   bool           `json:"completed"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName specifies the table name for Todo model
func (Todo) TableName() string {
	return "todos"
}

// CreateTodoInput represents the input for creating a todo
type CreateTodoInput struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

// UpdateTodoInput represents the input for updating a todo
type UpdateTodoInput struct {
	Title       *string `json:"title,omitempty"`
	Description *string `json:"description,omitempty"`
	Completed   *bool   `json:"completed,omitempty"`
}

// TodoFilter represents filtering options for todos
type TodoFilter struct {
	Completed *bool  `json:"completed,omitempty"`
	Search    string `json:"search,omitempty"`
}
