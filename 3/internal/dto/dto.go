package dto

import (
	"time"
)

type CreateEventRequest struct {
	UserID      int        `json:"user_id"`
	Title       string     `json:"title"`
	EventTime   time.Time  `json:"date"`
	Description *string    `json:"description"`
	Priority    string     `json:"priority"`
	RemindAt    *time.Time `json:"remind_at"`
}
type UpdateEventRequest struct {
	Title       string     `json:"title"`
	EventTime   time.Time  `json:"date"`
	Description *string    `json:"description"`
	Priority    string     `json:"priority"`
	RemindAt    *time.Time `json:"remind_at"`
}
