package worker

import (
	"4_3/internal/myLogger"
	"4_3/internal/repository"
	"strconv"
	"time"
)

type ArchiveWorker struct {
	repo     *repository.EventRepository
	interval time.Duration
	logger   *myLogger.Logger
}

func NewArchiveWorker(repo *repository.EventRepository, interval time.Duration, logger *myLogger.Logger) *ArchiveWorker {
	return &ArchiveWorker{
		repo:     repo,
		interval: interval,
		logger:   logger,
	}
}

func (w *ArchiveWorker) Start() {
	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	for range ticker.C {
		events := w.repo.GetAll()
		now := time.Now()
		for _, event := range events {
			if event.EventTime.Before(now) {
				err := w.repo.MoveToArchive(event.ID)
				if err != nil {
					w.logger.Error("failed move to archive", err)
				} else {
					msg := "event " + strconv.Itoa(event.ID) + " " + "moved to archive"
					w.logger.Info(msg)
				}
			}
		}
	}
}
