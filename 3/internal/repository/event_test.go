package repository

import (
	"4_3/internal/models"
	"testing"
	"time"
)

func TestMoveToArchive(t *testing.T) {
	repo := NewEventRepository()

	event := models.Event{
		ID:        1,
		UserID:    1,
		Title:     "archive-me",
		EventTime: time.Now().Add(-1 * time.Hour),
		Priority:  "high",
	}

	repo.CreateEvent(event)

	err := repo.MoveToArchive(1)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if _, ok := repo.Events[1]; ok {
		t.Fatalf("expected event to be removed from active events")
	}
	if _, ok := repo.ArchiveEvents[1]; !ok {
		t.Fatalf("expected event to be added to archive")
	}
}

func TestGetAllDoesNotReturnArchived(t *testing.T) {
	repo := NewEventRepository()

	event1 := models.Event{
		ID:        1,
		UserID:    1,
		Title:     "active",
		EventTime: time.Now().Add(1 * time.Hour),
		Priority:  "high",
	}
	event2 := models.Event{
		ID:        2,
		UserID:    1,
		Title:     "archived",
		EventTime: time.Now().Add(-1 * time.Hour),
		Priority:  "high",
	}

	repo.CreateEvent(event1)
	repo.CreateEvent(event2)

	err := repo.MoveToArchive(2)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	events := repo.GetAll()
	if len(events) != 1 {
		t.Fatalf("expected 1 active event, got %d", len(events))
	}
	if events[0].ID != 1 {
		t.Fatalf("expected active event id=1, got %d", events[0].ID)
	}
}
