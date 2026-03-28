package handlers

import (
	"4_3/internal/dto"
	"4_3/internal/myLogger"
	"4_3/internal/repository"
	"4_3/internal/service"
	"4_3/internal/worker"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"testing"
	"time"
)

func newTestHandler() *EventHandler {
	eventRepo := repository.NewEventRepository()
	remindCh := make(chan worker.ReminderTask, 100)
	eventService := service.NewEventService(eventRepo, remindCh)
	logger := myLogger.NewLogger()
	return NewEventHandler(eventService, logger)
}

func TestStatus(t *testing.T) {
	h := newTestHandler()

	req := httptest.NewRequest(http.MethodGet, "/status", nil)
	rr := httptest.NewRecorder()

	h.Status(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rr.Code)
	}
	if !strings.Contains(rr.Body.String(), "server is running") {
		t.Fatalf("expected body to contain 'server is running', got %s", rr.Body.String())
	}
}

func TestCreateEventSuccess(t *testing.T) {
	h := newTestHandler()

	now := time.Now()
	eventTime := now.Add(2 * time.Hour).Format(time.RFC3339)
	remindAt := now.Add(1 * time.Hour).Format(time.RFC3339)

	form := url.Values{}
	form.Set("user_id", "1")
	form.Set("title", "test")
	form.Set("date", eventTime)
	form.Set("description", "manual test")
	form.Set("priority", "high")
	form.Set("remind_at", remindAt)

	req := httptest.NewRequest(http.MethodPost, "/create_event", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr := httptest.NewRecorder()

	h.CreateEvent(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d, body=%s", rr.Code, rr.Body.String())
	}
	if !strings.Contains(rr.Body.String(), "event created") {
		t.Fatalf("expected body to contain 'event created', got %s", rr.Body.String())
	}
	if !strings.Contains(rr.Body.String(), `"id"`) {
		t.Fatalf("expected body to contain id, got %s", rr.Body.String())
	}
}

func TestCreateEventInvalidDate(t *testing.T) {
	h := newTestHandler()

	form := url.Values{}
	form.Set("user_id", "1")
	form.Set("title", "test")
	form.Set("date", "bad-date")
	form.Set("description", "manual test")
	form.Set("priority", "high")

	req := httptest.NewRequest(http.MethodPost, "/create_event", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr := httptest.NewRecorder()

	h.CreateEvent(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d, body=%s", rr.Code, rr.Body.String())
	}
	if !strings.Contains(rr.Body.String(), "invalid date") {
		t.Fatalf("expected body to contain 'invalid date', got %s", rr.Body.String())
	}
}

func TestCreateEventMethodNotAllowed(t *testing.T) {
	h := newTestHandler()

	req := httptest.NewRequest(http.MethodGet, "/create_event", nil)
	rr := httptest.NewRecorder()

	h.CreateEvent(rr, req)

	if rr.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected status 405, got %d", rr.Code)
	}
}

func TestEventsForDaySuccess(t *testing.T) {
	h := newTestHandler()

	now := time.Now()
	eventTime := now.Add(2 * time.Hour)
	remindAt := now.Add(1 * time.Hour)

	createReq := dto.CreateEventRequest{
		UserID:      1,
		Title:       "test-day",
		EventTime:   eventTime,
		Description: nil,
		Priority:    "high",
		RemindAt:    &remindAt,
	}

	_, err := h.eventService.Create(createReq)
	if err != nil {
		t.Fatalf("failed to create event: %v", err)
	}

	q := url.Values{}
	q.Set("user_id", "1")
	q.Set("date", eventTime.Format(time.RFC3339))

	req := httptest.NewRequest(http.MethodGet, "/events_for_day?"+q.Encode(), nil)
	rr := httptest.NewRecorder()

	h.EventsForDay(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d, body=%s", rr.Code, rr.Body.String())
	}
	if !strings.Contains(rr.Body.String(), "test-day") {
		t.Fatalf("expected body to contain created event, got %s", rr.Body.String())
	}
}

func TestEventsForDayInvalidDate(t *testing.T) {
	h := newTestHandler()

	req := httptest.NewRequest(http.MethodGet, "/events_for_day?user_id=1&date=bad-date", nil)
	rr := httptest.NewRecorder()

	h.EventsForDay(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d, body=%s", rr.Code, rr.Body.String())
	}
	if !strings.Contains(rr.Body.String(), "invalid date") {
		t.Fatalf("expected body to contain 'invalid date', got %s", rr.Body.String())
	}
}

func TestDeleteSuccess(t *testing.T) {
	h := newTestHandler()

	now := time.Now()
	eventTime := now.Add(2 * time.Hour)

	createReq := dto.CreateEventRequest{
		UserID:      1,
		Title:       "to-delete",
		EventTime:   eventTime,
		Description: nil,
		Priority:    "high",
	}

	id, err := h.eventService.Create(createReq)
	if err != nil {
		t.Fatalf("failed to create event: %v", err)
	}

	form := url.Values{}
	form.Set("id", strconv.Itoa(id))

	req := httptest.NewRequest(http.MethodPost, "/delete_event", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr := httptest.NewRecorder()

	h.Delete(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d, body=%s", rr.Code, rr.Body.String())
	}
	if !strings.Contains(rr.Body.String(), "event deleted") {
		t.Fatalf("expected body to contain 'event deleted', got %s", rr.Body.String())
	}
}

func TestDeleteInvalidID(t *testing.T) {
	h := newTestHandler()

	form := url.Values{}
	form.Set("id", "bad-id")

	req := httptest.NewRequest(http.MethodPost, "/delete_event", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr := httptest.NewRecorder()

	h.Delete(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d, body=%s", rr.Code, rr.Body.String())
	}
	if !strings.Contains(rr.Body.String(), "id to int") {
		t.Fatalf("expected body to contain 'id to int', got %s", rr.Body.String())
	}
}

func TestUpdateEventSuccess(t *testing.T) {
	h := newTestHandler()

	now := time.Now()
	eventTime := now.Add(2 * time.Hour)

	createReq := dto.CreateEventRequest{
		UserID:      1,
		Title:       "old-title",
		EventTime:   eventTime,
		Description: nil,
		Priority:    "high",
	}

	id, err := h.eventService.Create(createReq)
	if err != nil {
		t.Fatalf("failed to create event: %v", err)
	}

	newEventTime := now.Add(3 * time.Hour).Format(time.RFC3339)
	newRemindAt := now.Add(150 * time.Minute).Format(time.RFC3339)

	form := url.Values{}
	form.Set("id", strconv.Itoa(id))
	form.Set("title", "new-title")
	form.Set("date", newEventTime)
	form.Set("description", "updated")
	form.Set("priority", "medium")
	form.Set("remind_at", newRemindAt)

	req := httptest.NewRequest(http.MethodPost, "/update_event", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr := httptest.NewRecorder()

	h.UpdateEvent(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d, body=%s", rr.Code, rr.Body.String())
	}
	if !strings.Contains(rr.Body.String(), "event updated") {
		t.Fatalf("expected body to contain 'event updated', got %s", rr.Body.String())
	}
}

func TestUpdateEventInvalidDate(t *testing.T) {
	h := newTestHandler()

	form := url.Values{}
	form.Set("id", "1")
	form.Set("title", "new-title")
	form.Set("date", "bad-date")
	form.Set("description", "updated")
	form.Set("priority", "medium")
	form.Set("remind_at", time.Now().Add(1*time.Hour).Format(time.RFC3339))

	req := httptest.NewRequest(http.MethodPost, "/update_event", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr := httptest.NewRecorder()

	h.UpdateEvent(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d, body=%s", rr.Code, rr.Body.String())
	}
	if !strings.Contains(rr.Body.String(), "invalid date") {
		t.Fatalf("expected body to contain 'invalid date', got %s", rr.Body.String())
	}
}

func TestEventsForWeekSuccess(t *testing.T) {
	h := newTestHandler()

	now := time.Now()
	eventTime := now.Add(2 * time.Hour)

	createReq := dto.CreateEventRequest{
		UserID:      1,
		Title:       "test-week",
		EventTime:   eventTime,
		Description: nil,
		Priority:    "high",
	}

	_, err := h.eventService.Create(createReq)
	if err != nil {
		t.Fatalf("failed to create event: %v", err)
	}

	q := url.Values{}
	q.Set("user_id", "1")
	q.Set("date", eventTime.Format(time.RFC3339))

	req := httptest.NewRequest(http.MethodGet, "/events_for_week?"+q.Encode(), nil)
	rr := httptest.NewRecorder()

	h.EventsForWeek(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d, body=%s", rr.Code, rr.Body.String())
	}
	if !strings.Contains(rr.Body.String(), "test-week") {
		t.Fatalf("expected body to contain created event, got %s", rr.Body.String())
	}
}

func TestEventsForMonthSuccess(t *testing.T) {
	h := newTestHandler()

	now := time.Now()
	eventTime := now.Add(2 * time.Hour)

	createReq := dto.CreateEventRequest{
		UserID:      1,
		Title:       "test-month",
		EventTime:   eventTime,
		Description: nil,
		Priority:    "high",
	}

	_, err := h.eventService.Create(createReq)
	if err != nil {
		t.Fatalf("failed to create event: %v", err)
	}

	q := url.Values{}
	q.Set("user_id", "1")
	q.Set("date", eventTime.Format(time.RFC3339))

	req := httptest.NewRequest(http.MethodGet, "/events_for_month?"+q.Encode(), nil)
	rr := httptest.NewRecorder()

	h.EventsForMonth(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d, body=%s", rr.Code, rr.Body.String())
	}
	if !strings.Contains(rr.Body.String(), "test-month") {
		t.Fatalf("expected body to contain created event, got %s", rr.Body.String())
	}
}

func TestDeleteMethodNotAllowed(t *testing.T) {
	h := newTestHandler()

	req := httptest.NewRequest(http.MethodGet, "/delete_event", nil)
	rr := httptest.NewRecorder()

	h.Delete(rr, req)

	if rr.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected status 405, got %d", rr.Code)
	}
}

func TestUpdateEventMethodNotAllowed(t *testing.T) {
	h := newTestHandler()

	req := httptest.NewRequest(http.MethodGet, "/update_event", nil)
	rr := httptest.NewRecorder()

	h.UpdateEvent(rr, req)

	if rr.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected status 405, got %d", rr.Code)
	}
}
