package models

import "time"

type Task struct {
	ID          int64      `json:"id"`
	UserID      int64      `json:"user_id"`
	CategoryID  *int64     `json:"category_id,omitempty"`
	ParentID    *int64     `json:"parent_id,omitempty"`
	Title       string     `json:"title"`
	Description string     `json:"description,omitempty"`
	Status      string     `json:"status"`
	Priority    int        `json:"priority"`
	DueDate     *time.Time `json:"due_date,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}
