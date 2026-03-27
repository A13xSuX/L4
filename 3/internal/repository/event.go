package repository

import (
	"4_3/internal/models"
	"errors"
	"sync"
)

type EventRepository struct {
	m             sync.RWMutex
	Events        map[int]models.Event
	ArchiveEvents map[int]models.Event
}

func NewEventRepository() *EventRepository {
	return &EventRepository{
		Events:        make(map[int]models.Event),
		ArchiveEvents: make(map[int]models.Event),
	}
}

func (r *EventRepository) CreateEvent(event models.Event) {
	r.m.Lock()
	defer r.m.Unlock()
	r.Events[event.ID] = event
}

func (r *EventRepository) GetByID(id int) (models.Event, error) {
	r.m.RLock()
	defer r.m.RUnlock()
	if _, ok := r.Events[id]; !ok {
		return models.Event{}, errors.New("event not found")
	}
	return r.Events[id], nil
}

func (r *EventRepository) GetAll() []models.Event {
	r.m.RLock()
	defer r.m.RUnlock()

	events := make([]models.Event, 0, len(r.Events))
	for k, v := range r.Events {
		if _, ok := r.ArchiveEvents[k]; !ok {
			events = append(events, v)
		}
	}
	return events
}

func (r *EventRepository) DeleteEvent(id int) error {
	r.m.Lock()
	defer r.m.Unlock()
	if _, ok := r.Events[id]; !ok {
		return errors.New("event not found")
	}
	delete(r.Events, id)
	return nil
}

func (r *EventRepository) UpdateEvent(id int, newEvent models.Event) error {
	r.m.Lock()
	defer r.m.Unlock()
	if _, ok := r.Events[id]; !ok {
		return errors.New("event not found")
	}
	r.Events[id] = newEvent
	return nil
}

func (r *EventRepository) MoveToArchive(id int) error {
	r.m.Lock()
	defer r.m.Unlock()

	event, ok := r.Events[id]
	if !ok {
		return errors.New("event not found")
	}

	r.ArchiveEvents[id] = event
	delete(r.Events, id)
	return nil
}
