package worker

import "time"

type ReminderTask struct {
	ID        int
	UserID    int
	Title     string
	EventTime time.Time
	RemindAt  time.Time
}
