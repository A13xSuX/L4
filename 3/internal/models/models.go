package models

import "time"

type Event struct {
	ID          int        `json:"id"`
	UserID      int        `json:"user_id"`
	Title       string     `json:"title"`
	EventTime   time.Time  `json:"date"`
	Description *string    `json:"description"`
	Priority    string     `json:"priority"`
	RemindAt    *time.Time `json:"remind_at"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at"`
}
