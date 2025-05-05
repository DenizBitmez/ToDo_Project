package model

import "time"

type TodoStep struct {
	ID        int        `json:"id"`
	TODOID    int        `json:"todo_id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
	Content   string     `json:"content"`
	Status    int        `json:"status"`
}
