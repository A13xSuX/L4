package worker

import "time"

type RemindWorker struct {
	ReminderCh chan ReminderTask
}

func NewRemindWorker(reminderCh chan ReminderTask) *RemindWorker {
	return &RemindWorker{
		ReminderCh: reminderCh,
	}
}

func (w *RemindWorker) Start() {
	for {
		now := time.Now()
		task := <-w.ReminderCh
		if task.RemindAt.Before(now) {
			continue //напоминать не нужно - успешно
		}
		timeToWait := task.RemindAt.Sub(now)
		go func() {
			<-time.After(timeToWait)
			// send task
			// напишем в лог что сообщение отправлено когда время подойдет
		}()
	}
}
