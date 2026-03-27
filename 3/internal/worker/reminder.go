package worker

import (
	"4_3/internal/myLogger"
	"strconv"
	"time"
)

type RemindWorker struct {
	ReminderCh chan ReminderTask
	logger     *myLogger.Logger
}

func NewRemindWorker(reminderCh chan ReminderTask, logger *myLogger.Logger) *RemindWorker {
	return &RemindWorker{
		ReminderCh: reminderCh,
		logger:     logger,
	}
}

func (w *RemindWorker) Start() {
	for task := range w.ReminderCh {
		now := time.Now()
		if task.RemindAt.Before(now) {
			msg := "reminder skipped id=" + strconv.Itoa(task.ID) + " title=" + task.Title + " because remind time already passed"
			w.logger.Info(msg)
			continue //напоминать не нужно - успешно
		}
		timeToWait := task.RemindAt.Sub(now)
		go func() {
			<-time.After(timeToWait)
			msg := "reminder sent for event id=" + strconv.Itoa(task.ID) + " title=" + task.Title
			w.logger.Info(msg)
		}()
	}
}
