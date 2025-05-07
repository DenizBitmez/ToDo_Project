package model

import "time"

type TodoStep struct {
	ID        int        `json:"id"`
	TODOID    int        `json:"todo_id"`
	Title     string     `json:"title"`
	Content   string     `json:"content"`
	Status    int        `json:"status"`
	Username  string     `json:"username"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

func (s *TodoStep) IsCompleted() bool {
	return s.Status == 1
}
