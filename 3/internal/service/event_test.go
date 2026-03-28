package service

import (
	"4_3/internal/customErrs"
	"4_3/internal/dto"
	"4_3/internal/repository"
	"4_3/internal/worker"
	"errors"
	"testing"
	"time"
)

func TestCreateEventWithoutRemind(t *testing.T) {
	eventRepo := repository.NewEventRepository()
	remindCh := make(chan worker.ReminderTask, 100)
	eventService := NewEventService(eventRepo, remindCh)
	eventReq := dto.CreateEventRequest{
		UserID:      1,
		Title:       "test",
		EventTime:   time.Now().Add(2 * time.Hour),
		Description: nil,
		Priority:    "high",
	}
	id, err := eventService.Create(eventReq)
	if err != nil {
		t.Fatalf("Got error %v", err)
	}
	if id <= 0 {
		t.Errorf("Expected id > 0, but got %d", id)
	}
	if len(eventRepo.Events) != 1 {
		t.Errorf("Expected 1 events, but got %d", len(eventRepo.Events))
	}
	if len(remindCh) != 0 {
		t.Errorf("Expected 0 remind task, but got %d", len(remindCh))
	}
}

func TestCreateEventWithRemind(t *testing.T) {
	eventRepo := repository.NewEventRepository()
	remindCh := make(chan worker.ReminderTask, 100)
	eventService := NewEventService(eventRepo, remindCh)
	remindAt := time.Now().Add(15 * time.Minute)
	eventReq := dto.CreateEventRequest{
		UserID:      1,
		Title:       "test",
		EventTime:   time.Now().Add(2 * time.Hour),
		Description: nil,
		Priority:    "high",
		RemindAt:    &remindAt,
	}
	id, err := eventService.Create(eventReq)
	if err != nil {
		t.Fatalf("Got error %v", err)
	}
	if id <= 0 {
		t.Errorf("Expected id > 0, but got %d", id)
	}
	if len(eventRepo.Events) != 1 {
		t.Errorf("Expected 1 events, but got %d", len(eventRepo.Events))
	}
	if len(remindCh) != 1 {
		t.Errorf("Expected 1 remind task, but got %d", len(remindCh))
	}
	task := <-remindCh
	if task.ID != id {
		t.Errorf("Expected %d id, but got %d", id, task.ID)
	}
	if task.UserID != 1 {
		t.Errorf("Expected %d user_id, but got %d", 1, task.UserID)
	}
	if task.Title != "test" {
		t.Errorf("Expected title=test, but got %s", task.Title)
	}
}

func TestCreateEventWithInvalidUserID(t *testing.T) {
	eventRepo := repository.NewEventRepository()
	remindCh := make(chan worker.ReminderTask, 100)
	eventService := NewEventService(eventRepo, remindCh)

	eventReq := dto.CreateEventRequest{
		UserID:      0,
		Title:       "test",
		EventTime:   time.Now().Add(2 * time.Hour),
		Description: nil,
		Priority:    "high",
	}

	id, err := eventService.Create(eventReq)
	if !errors.Is(err, customErrs.InvalidUserIDErr) {
		t.Fatalf("expected InvalidUserIDErr, got %v", err)
	}
	if id != -1 {
		t.Fatalf("expected id = -1, got %d", id)
	}
	if len(eventRepo.Events) != 0 {
		t.Fatalf("expected 0 events, got %d", len(eventRepo.Events))
	}
	if len(remindCh) != 0 {
		t.Fatalf("expected 0 remind tasks, got %d", len(remindCh))
	}
}

func TestCreateEventWithPastEventTime(t *testing.T) {
	eventRepo := repository.NewEventRepository()
	remindCh := make(chan worker.ReminderTask, 100)
	eventService := NewEventService(eventRepo, remindCh)

	eventReq := dto.CreateEventRequest{
		UserID:      1,
		Title:       "test",
		EventTime:   time.Now().Add(-1 * time.Minute),
		Description: nil,
		Priority:    "high",
	}

	id, err := eventService.Create(eventReq)
	if !errors.Is(err, customErrs.EventPastTimeErr) {
		t.Fatalf("expected EventPastTimeErr, got %v", err)
	}
	if id != -1 {
		t.Fatalf("expected id = -1, got %d", id)
	}
	if len(eventRepo.Events) != 0 {
		t.Fatalf("expected 0 events, got %d", len(eventRepo.Events))
	}
	if len(remindCh) != 0 {
		t.Fatalf("expected 0 remind tasks, got %d", len(remindCh))
	}
}

func TestCreateEventWithEmptyTitle(t *testing.T) {
	eventRepo := repository.NewEventRepository()
	remindCh := make(chan worker.ReminderTask, 100)
	eventService := NewEventService(eventRepo, remindCh)

	eventReq := dto.CreateEventRequest{
		UserID:      1,
		Title:       "",
		EventTime:   time.Now().Add(2 * time.Hour),
		Description: nil,
		Priority:    "high",
	}

	id, err := eventService.Create(eventReq)
	if !errors.Is(err, customErrs.TitleEmptyErr) {
		t.Fatalf("expected TitleEmptyErr, got %v", err)
	}
	if id != -1 {
		t.Fatalf("expected id = -1, got %d", id)
	}
	if len(eventRepo.Events) != 0 {
		t.Fatalf("expected 0 events, got %d", len(eventRepo.Events))
	}
	if len(remindCh) != 0 {
		t.Fatalf("expected 0 remind tasks, got %d", len(remindCh))
	}
}

func TestCreateEventWithPastRemindAt(t *testing.T) {
	eventRepo := repository.NewEventRepository()
	remindCh := make(chan worker.ReminderTask, 100)
	eventService := NewEventService(eventRepo, remindCh)

	remindAt := time.Now().Add(-1 * time.Minute)
	eventReq := dto.CreateEventRequest{
		UserID:      1,
		Title:       "test",
		EventTime:   time.Now().Add(2 * time.Hour),
		Description: nil,
		Priority:    "high",
		RemindAt:    &remindAt,
	}

	id, err := eventService.Create(eventReq)
	if !errors.Is(err, customErrs.RemindAtPastErr) {
		t.Fatalf("expected RemindAtPastErr, got %v", err)
	}
	if id != -1 {
		t.Fatalf("expected id = -1, got %d", id)
	}
	if len(eventRepo.Events) != 0 {
		t.Fatalf("expected 0 events, got %d", len(eventRepo.Events))
	}
	if len(remindCh) != 0 {
		t.Fatalf("expected 0 remind tasks, got %d", len(remindCh))
	}
}

func TestCreateEventWithRemindAfterEventTime(t *testing.T) {
	eventRepo := repository.NewEventRepository()
	remindCh := make(chan worker.ReminderTask, 100)
	eventService := NewEventService(eventRepo, remindCh)

	now := time.Now()
	eventTime := now.Add(1 * time.Hour)
	remindAt := now.Add(2 * time.Hour)

	eventReq := dto.CreateEventRequest{
		UserID:      1,
		Title:       "test",
		EventTime:   eventTime,
		Description: nil,
		Priority:    "high",
		RemindAt:    &remindAt,
	}

	id, err := eventService.Create(eventReq)
	if !errors.Is(err, customErrs.RemindAtAfterEventTimeErr) {
		t.Fatalf("expected RemindAtAfterEventTimeErr, got %v", err)
	}
	if id != -1 {
		t.Fatalf("expected id = -1, got %d", id)
	}
	if len(eventRepo.Events) != 0 {
		t.Fatalf("expected 0 events, got %d", len(eventRepo.Events))
	}
	if len(remindCh) != 0 {
		t.Fatalf("expected 0 remind tasks, got %d", len(remindCh))
	}
}

func TestCreateEventWithInvalidPriority(t *testing.T) {
	eventRepo := repository.NewEventRepository()
	remindCh := make(chan worker.ReminderTask, 100)
	eventService := NewEventService(eventRepo, remindCh)

	eventReq := dto.CreateEventRequest{
		UserID:      1,
		Title:       "test",
		EventTime:   time.Now().Add(2 * time.Hour),
		Description: nil,
		Priority:    "super-high",
	}

	id, err := eventService.Create(eventReq)
	if !errors.Is(err, customErrs.PriorityErr) {
		t.Fatalf("expected PriorityErr, got %v", err)
	}
	if id != -1 {
		t.Fatalf("expected id = -1, got %d", id)
	}
	if len(eventRepo.Events) != 0 {
		t.Fatalf("expected 0 events, got %d", len(eventRepo.Events))
	}
	if len(remindCh) != 0 {
		t.Fatalf("expected 0 remind tasks, got %d", len(remindCh))
	}
}
