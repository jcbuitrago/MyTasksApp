package models

type Task struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"` // e.g., "pending", "in progress", "completed"
	CategoryID  int    `json:"category_id"`
	UserID      int    `json:"user_id"`
}