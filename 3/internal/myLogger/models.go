package myLogger

import "time"

type LogMsg struct {
	Time       time.Time
	Level      string
	Msg        string
	Err        error
	Method     string
	Path       string
	StatusCode int
	Duration   time.Duration
}
