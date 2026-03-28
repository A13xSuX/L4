package myLogger

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type Logger struct {
	LogCh chan LogMsg
}

func NewLogger() *Logger {
	return &Logger{
		LogCh: make(chan LogMsg, 100),
	}
}

func (l *Logger) Start() {
	for msg := range l.LogCh {
		log := formatLogMsg(msg)
		fmt.Println(log)
	}
}

func (l *Logger) Info(msg string) {
	logMsg := LogMsg{
		Time:  time.Now(),
		Level: "INFO",
		Msg:   msg,
	}
	l.LogCh <- logMsg
}

func (l *Logger) Error(msg string, err error) {
	logMsg := LogMsg{
		Time:  time.Now(),
		Level: "ERROR",
		Msg:   msg,
		Err:   err,
	}
	l.LogCh <- logMsg
}

func (l *Logger) LogHTTP(msg, method, path string, statusCode int, duration time.Duration) {
	logMsg := LogMsg{
		Time:       time.Now(),
		Level:      "INFO",
		Msg:        msg,
		Method:     method,
		Path:       path,
		StatusCode: statusCode,
		Duration:   duration,
	}
	l.LogCh <- logMsg
}
func formatLogMsg(msg LogMsg) string {
	parts := []string{}
	level := "[" + msg.Level + "]"
	parts = append(parts, msg.Time.Format(time.RFC3339), level, msg.Msg)
	if msg.Err != nil {
		err := "err=" + msg.Err.Error()
		parts = append(parts, err)
	}
	if msg.Method != "" {
		method := "method=" + msg.Method
		parts = append(parts, method)
	}
	if msg.Path != "" {
		path := "path=" + msg.Path
		parts = append(parts, path)
	}
	if msg.StatusCode != 0 {
		code := strconv.Itoa(msg.StatusCode)
		status := "status=" + code
		parts = append(parts, status)
	}
	if msg.Duration != 0 {
		duration := "duration=" + msg.Duration.String()
		parts = append(parts, duration)
	}
	log := strings.Join(parts, " ")
	return log
}
