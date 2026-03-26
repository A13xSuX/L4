package worker

import (
	"4_3/internal/repository"
	"time"
)

type ArchiveWorker struct {
	repo     *repository.EventRepository
	interval time.Duration
}

func NewArchiveWorker(repo *repository.EventRepository, interval time.Duration) *ArchiveWorker {
	return &ArchiveWorker{
		repo:     repo,
		interval: interval,
	}
}

func (w *ArchiveWorker) Start() {
	ticker := time.NewTicker(w.interval)

	go func() {
		for {
			_ = <-ticker.C
			events := w.repo.GetAll()
			for _, event := range events {
				if event.EventTime.Before(time.Now()) {
					_ = w.repo.MoveToArchive(event.ID)
				}
			}
		}
	}()

}
