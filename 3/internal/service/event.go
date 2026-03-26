package service

import (
	"4_3/internal/dto"
	"4_3/internal/models"
	"4_3/internal/repository"
	"errors"
	"sync"
	"time"
)

type EventService struct {
	m    sync.Mutex
	ID   int
	repo *repository.EventRepository
}

func NewEventService(repo *repository.EventRepository) *EventService {
	return &EventService{
		ID:   0,
		repo: repo,
	}
}

func (s *EventService) Create(eventReq dto.CreateEventRequest) (int, error) {
	now := time.Now()
	if eventReq.UserID <= 0 {
		return -1, errors.New("invalid UserID")
	}
	if eventReq.EventTime.Before(now) {
		return -1, errors.New("eventTime must be in the future")
	}
	if eventReq.Title == "" {
		return -1, errors.New("title empty")
	}
	if err := validationPriority(eventReq.Priority); err != nil {
		return -1, err
	}
	if eventReq.RemindAt != nil && eventReq.RemindAt.Before(now) {
		return -1, errors.New("remindAt must be in the future")
	}
	if eventReq.RemindAt != nil && !eventReq.RemindAt.Before(eventReq.EventTime) {
		return -1, errors.New("remindAt must be before eventTime")
	}

	event := models.Event{
		ID:          s.generateID(),
		UserID:      eventReq.UserID,
		Title:       eventReq.Title,
		EventTime:   eventReq.EventTime,
		Description: eventReq.Description,
		Priority:    eventReq.Priority,
		RemindAt:    eventReq.RemindAt,
		CreatedAt:   now,
	}
	s.repo.CreateEvent(event)
	return event.ID, nil
}

func (s *EventService) Update(id int, eventReq dto.UpdateEventRequest) error {
	if id <= 0 {
		return errors.New("invalid id")
	}
	now := time.Now()
	if eventReq.EventTime.Before(now) {
		return errors.New("eventTime must be in the future")
	}
	if eventReq.Title == "" {
		return errors.New("title empty")
	}

	if err := validationPriority(eventReq.Priority); err != nil {
		return err
	}
	if eventReq.RemindAt != nil && eventReq.RemindAt.Before(now) {
		return errors.New("remindAt must be in the future")
	}
	if eventReq.RemindAt != nil && !eventReq.RemindAt.Before(eventReq.EventTime) {
		return errors.New("remindAt must be before eventTime")
	}

	oldEvent, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}

	newEvent := models.Event{
		ID:          oldEvent.ID,
		UserID:      oldEvent.UserID,
		Title:       eventReq.Title,
		EventTime:   eventReq.EventTime,
		Description: eventReq.Description,
		Priority:    eventReq.Priority,
		RemindAt:    eventReq.RemindAt,
		CreatedAt:   oldEvent.CreatedAt,
		UpdatedAt:   timePtr(now),
	}
	err = s.repo.UpdateEvent(id, newEvent)
	if err != nil {
		return err
	}
	return nil
}

func (s *EventService) Delete(id int) error {
	if id <= 0 {
		return errors.New("invalid id")
	}
	return s.repo.DeleteEvent(id)
}

func (s *EventService) GetEventsForDay(userID int, date time.Time) ([]models.Event, error) {
	if userID <= 0 {
		return nil, errors.New("invalid userID")
	}
	allEvents := s.repo.GetAll()

	year, month, day := date.Date()
	startOfDay := time.Date(year, month, day, 0, 0, 0, 0, date.Location())
	nextDay := startOfDay.AddDate(0, 0, 1)
	userEvents := []models.Event{}
	for _, event := range allEvents {
		if event.UserID == userID && !event.EventTime.Before(startOfDay) && event.EventTime.Before(nextDay) {
			userEvents = append(userEvents, event)
		}
	}
	return userEvents, nil
}

func (s *EventService) GetEventsForWeek(userID int, date time.Time) ([]models.Event, error) {
	if userID <= 0 {
		return nil, errors.New("invalid userID")
	}
	allEvents := s.repo.GetAll()

	year, month, day := date.Date()
	startOfDay := time.Date(year, month, day, 0, 0, 0, 0, date.Location())

	dayOfWeek := startOfDay.Weekday()
	daysForMonday := int(dayOfWeek) - 1
	if daysForMonday < 0 {
		daysForMonday = 6
	}
	startOfWeek := startOfDay.AddDate(0, 0, -daysForMonday)
	nextWeek := startOfWeek.AddDate(0, 0, 7)
	userEvents := []models.Event{}
	for _, event := range allEvents {
		if event.UserID == userID && !event.EventTime.Before(startOfWeek) && event.EventTime.Before(nextWeek) {
			userEvents = append(userEvents, event)
		}
	}
	return userEvents, nil
}

func (s *EventService) GetEventsForMonth(userID int, date time.Time) ([]models.Event, error) {
	if userID <= 0 {
		return nil, errors.New("invalid userID")
	}
	allEvents := s.repo.GetAll()

	year, month, _ := date.Date()
	startOfMonth := time.Date(year, month, 1, 0, 0, 0, 0, date.Location())
	nextMonth := startOfMonth.AddDate(0, 1, 0)
	userEvents := []models.Event{}
	for _, event := range allEvents {
		if event.UserID == userID && !event.EventTime.Before(startOfMonth) && event.EventTime.Before(nextMonth) {
			userEvents = append(userEvents, event)
		}
	}
	return userEvents, nil
}

func (s *EventService) generateID() int {
	s.m.Lock()
	defer s.m.Unlock()
	s.ID++
	return s.ID
}

func validationPriority(priority string) error {
	valid := []string{"high", "medium", "low"}
	for _, p := range valid {
		if p == priority {
			return nil
		}
	}
	return errors.New("priority not in list: high, medium, low")
}

func timePtr(t time.Time) *time.Time {
	return &t
}
