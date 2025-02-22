package domain

import (
	"time"
)

// User представляет модель пользователя
type User struct {
	ID        string     `json:"id" example:"1"`
	Name      string     `json:"name" example:"Иван Петров"`
	Email     string     `json:"email" example:"ivan@example.com"`
	CreatedAt time.Time  `json:"created_at" example:"2024-01-01T00:00:00Z"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

// Conditions представляет параметры пагинации
type Conditions struct {
	Limit  int `json:"limit" example:"10"`
	Offset int `json:"offset" example:"0"`
}
